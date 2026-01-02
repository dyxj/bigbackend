package idempotency

import "time"

type Config struct {
	extractor     KeyExtractor
	errRespWriter ErrorResponseWriter
	lockOptions   []LockOption
	cacheExpiry   time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		extractor:     DefaultKeyExtractor,
		errRespWriter: DefaultErrorResponseWriter,
		lockOptions:   []LockOption{},
		cacheExpiry:   24 * time.Hour,
	}
}

type Option func(*Config)

func WithKeyExtractor(extractor KeyExtractor) Option {
	return func(c *Config) {
		c.extractor = extractor
	}
}

func WithErrorResponseWriter(writer ErrorResponseWriter) Option {
	return func(c *Config) {
		c.errRespWriter = writer
	}
}

func WithLockOptions(options ...LockOption) Option {
	return func(c *Config) {
		c.lockOptions = options
	}
}

func WithCacheExpiry(expiry time.Duration) Option {
	return func(c *Config) {
		c.cacheExpiry = expiry
	}
}
