//go:build integration

package integration

import (
	"log"
	"os"
	"testing"

	"github.com/dyxj/bigbackend/internal/user/profile"
	"github.com/dyxj/bigbackend/pkg/testx"
	"go.uber.org/zap"
)

var (
	logger                 *zap.Logger
	userProfileCreatorRepo profile.CreatorRepo
)

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()

	ready, errChan := testx.RunGlobalEnv()

	select {
	case <-ready:
		log.Printf("test environment is ready")
	case err := <-errChan:
		log.Panicf("failed to start test environment: %v", err)
	}
	defer testx.GlobalEnv().Close()

	var err error
	logger, err = zap.NewDevelopment(zap.Fields(
		zap.String("logger", "test"),
	))
	if err != nil {
		log.Panicf("failed to test initialize logger: %v", err)
	}
	userProfileCreatorRepo = profile.NewCreatorSQLDB(logger)

	code = m.Run()
}
