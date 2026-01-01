package idempotency

import "errors"

var ErrInProgress = errors.New("key in progress")
