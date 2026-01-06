package safe

import (
	"log"
	"runtime/debug"
)

func Go(fn func()) {
	GoWithRecover(fn, nil)
}

func GoWithRecover(
	fn func(),
	recoverFn func(err interface{}, debugStack []byte),
) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if recoverFn != nil {
					recoverFn(err, debug.Stack())
				} else {
					log.Printf("Recovered from goroutine panic: %v\n%s", err, debug.Stack())
				}
			}
		}()
		fn()
	}()
}
