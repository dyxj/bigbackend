package main

import (
	// Set GOMAXPROCS to match Linux container CPU quota(limit, cpu.cfs_quota_us)
	// IMO: we should not set k8 cpu.limit rather only set cpu.request(cpu.shares).
	// This allows pod to utilize unused CPU while still making sure pods are guaranteed
	// their requested CPU.

	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dyxj/bigbackend/internal/config"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	// Set default logger to stdout
	log.SetOutput(os.Stdout)

	_, err := maxprocs.Set(maxprocs.Logger(log.Printf))
	if err != nil {
		log.Panicf("failed to set GOMAXPROCS: %v", err)
	}
}

func main() {
	// listen to interrupt and termination signals
	mainCtx, mainStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// ensures stop function is called on exit to avoid unintended diversion of signals to context
	defer mainStop()

	// used to terminate program in case of initialization failures
	//errSign := make(chan struct{})

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
}
