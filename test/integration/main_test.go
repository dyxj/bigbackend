//go:build integration

package integration

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

	ready, errors := testx.RunGlobalEnv()

	select {
	case <-ready:
		log.Printf("test environment is ready")
	case err := <-errors:
		log.Panicf("failed to start test environment: %v", err)
	}
	defer testx.GlobalEnv().Close()

	code = m.Run()
}
