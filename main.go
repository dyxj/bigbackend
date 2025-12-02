package main

import (
	// Set GOMAXPROCS to match Linux container CPU quota(limit, cpu.cfs_quota_us)
	// IMO: we should not set k8 cpu.limit rather only set cpu.request(cpu.shares).
	// This allows pod to utilize unused CPU while still making sure pods are guaranteed
	// their requested CPU.

	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dyxj/bigbackend/internal/config"
	"github.com/dyxj/bigbackend/internal/userprofile"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func init() {
	// Set default logger to stdout
	log.SetOutput(os.Stdout)

	_, err := maxprocs.Set(maxprocs.Logger(log.Printf))
	if err != nil {
		log.Panicf("failed to set GOMAXPROCS: %v", err)
	}
}

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

func main() {
	// listen to interrupt and termination signals
	mainCtx, mainStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// ensures stop function is called on exit to avoid unintended diversion of signals to context
	defer mainStop()

	// used to terminate program in case of initialization failures
	errSign := make(chan struct{})

	// Parse environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logx.InitLogger()
	if err != nil {
		log.Panicf("failed to init logger: %v", err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			// Will need to dig into details "sync /dev/stdout: bad file descriptor"
			if errors.Is(err, syscall.EBADF) {
				return
			}
			log.Printf("failed to perform log sync: %v", err)
		}
	}()

	// Initialize database connection
	dbConn, err := sqldb.NewDBConn(mainCtx, cfg.DBConfig)
	if err != nil {
		log.Panicf("failed to connect to database: %v", err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Printf("failed to close db conn: %v", err)
		}
	}()

	// Run database migrations
	err = sqldb.RunMigration(dbConn, nil)
	if err != nil {
		log.Panicf("failed to run database migrations: %v", err)
	}

	server, serverForceStop := setupHTTPServer(mainCtx, cfg.HTTPServerConfig, logger)

	go func() {
		logger.Info("starting server")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", zap.Error(err))
			errSign <- struct{}{}
		}
		logger.Info("server stopped")
	}()

	select {
	case <-errSign:
		logger.Error("unexpected error occurred while starting up server")
	case <-mainCtx.Done():
	}
	mainStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	serverForceStop()
	if err != nil {
		logger.Error("Failed to wait for ongoing requests to finish, waiting for forced cancellation")
		time.Sleep(_shutdownHardPeriod)
		logger.Error("server shut down ungracefully")
		return
	}

	logger.Info("server shut down gracefully")
}

func setupHTTPServer(ctx context.Context, serverConfig *config.HTTPServerConfig, logger *zap.Logger) (*http.Server, context.CancelFunc) {

	userProfileGetterHandler := userprofile.NewGetterHandler(logger)

	router := http.NewServeMux()

	router.HandleFunc("GET /user/{id}", userProfileGetterHandler.ServeHTTP)

	ongoingCtx, forceStopOngoingCtx := context.WithCancel(ctx)
	server := &http.Server{
		Addr:              fmt.Sprintf("%v:%v", serverConfig.Host(), serverConfig.Port()),
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Second,
		IdleTimeout:       time.Second,
		// TODO better error body
		Handler: http.TimeoutHandler(router, time.Second, ""),
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	return server, forceStopOngoingCtx
}
