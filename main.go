package main

import (
	// Set GOMAXPROCS to match Linux container CPU quota(limit, cpu.cfs_quota_us)
	// IMO: we should not set k8 cpu.limit rather only set cpu.request(cpu.shares).
	// This allows pod to utilize unused CPU while still making sure pods are guaranteed
	// their requested CPU.

	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/dyxj/bigbackend/internal/config"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	_ "go.uber.org/automaxprocs"
)

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
			log.Println(err)
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
			log.Println(err)
		}
	}()

	// Run database migrations
	err = sqldb.RunMigration(dbConn, nil)
	if err != nil {
		log.Panicf("failed to run database migrations: %v", err)
	}
}
