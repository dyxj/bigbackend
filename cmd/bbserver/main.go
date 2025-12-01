package main

import (
	// Set GOMAXPROCS to match Linux container CPU quota(limit, cpu.cfs_quota_us)
	// IMO: we should not set k8 cpu.limit rather only set cpu.request(cpu.shares).
	// This allows pod to utilize unused CPU while still making sure pods are guaranteed
	// their requested CPU.
	_ "go.uber.org/automaxprocs"
)

func main() {
}
