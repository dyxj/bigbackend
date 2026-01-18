//go:build integration

package bench

import (
	"log"
	"os"
	"testing"

	"github.com/dyxj/bigbackend/pkg/testx"
)

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()

	ready, errChan := testx.RunGlobalBenchEnv()

	select {
	case <-ready:
		log.Printf("test environment is ready")
	case err := <-errChan:
		log.Panicf("failed to start test environment: %v", err)
	}
	defer testx.GlobalBenchEnv().Close()

	code = m.Run()
}
