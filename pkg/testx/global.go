package testx

import (
	"log"
	"sync"
)

var (
	global     *Environment
	globalOnce sync.Once

	globalBench     *Environment
	globalBenchOnce sync.Once
)

func RunGlobalEnv() (<-chan struct{}, <-chan error) {
	globalOnce.Do(func() {
		global = NewEnvironment("global", true, false)
	})
	return global.Run()
}

func GlobalEnv() *Environment {
	if global == nil {
		log.Panicf("global test environment is not initialized")
	}
	return global
}

func RunGlobalBenchEnv() (<-chan struct{}, <-chan error) {
	globalBenchOnce.Do(func() {
		globalBench = NewEnvironment("global-bench", true, true)
	})
	return globalBench.Run()
}

func GlobalBenchEnv() *Environment {
	if globalBench == nil {
		log.Panicf("global-bench test environment is not initialized")
	}
	return globalBench
}
