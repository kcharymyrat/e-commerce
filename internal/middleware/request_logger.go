package middleware

import (
	"net/http"
	"time"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/rs/zerolog"
)

func RequestLogger(app *app.Application) func(http.Handler) http.Handler {
	// Create a sampled logger using zerolog's LevelSampler.
	sampledLogger := app.Logger.Sample(&zerolog.LevelSampler{
		DebugSampler: &zerolog.BurstSampler{
			Burst:       5, // Allow 5 debug logs within a second before sampling kicks in
			Period:      1 * time.Second,
			NextSampler: &zerolog.BasicSampler{N: 100}, // Log 1 out of every 100 after burst
		},
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			sampledLogger.Info().
				Str("method", r.Method).
				Str("url", r.URL.Path).
				Dur("duration", duration).
				Msg("request completed")
		})
	}
}
