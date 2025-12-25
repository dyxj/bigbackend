package app

import "time"

type HttpConfig interface {
	IsDebug() bool
	Host() string
	Port() int
	ReadHeaderTimeout() time.Duration
	ReadTimeout() time.Duration
	IdleTimeout() time.Duration
	HandlerTimeout() time.Duration

	ShutDownTimeout() time.Duration
	ShutDownHardTimeout() time.Duration
	ReadinessProbeDelay() time.Duration
}
