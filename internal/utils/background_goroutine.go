package utils

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

func BackgroundGoroutine(logger *zerolog.Logger, wg *sync.WaitGroup, fn func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				logger.Err(fmt.Errorf("%v", err)).Msg("panic")
			}
		}()

		fn()
	}()
}
