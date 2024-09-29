package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GeneralRateLimiter(app *app.Application) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)
			key := "project:general"

			res, err := app.Limiter.Allow(r.Context(), key, redis_rate.PerMinute(10_000))
			if err != nil {
				app.Logger.Error().
					Err(err).
					Str("url", r.URL.String()).
					Str("method", r.Method).
					Str("remote_addr", r.RemoteAddr).
					Str("rate_limit_key", key).
					Msg("Rate limiting error")
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}

			// // Log remaining requests for the rate limit period
			// app.Logger.Debug().
			// 	Str("url", r.URL.String()).
			// 	Str("method", r.Method).
			// 	Str("remote_addr", r.RemoteAddr).
			// 	Str("rate_limit_key", key).
			// 	Int("rate_limit_remaining", res.Remaining).
			// 	Msg("Rate limit check")

			w.Header().Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))

			if res.Allowed == 0 {
				seconds := int(res.RetryAfter / time.Second)
				w.Header().Set("RateLimit-RetryAfter", strconv.Itoa(seconds))

				app.Logger.Warn().
					Str("url", r.URL.String()).
					Str("method", r.Method).
					Str("remote_addr", r.RemoteAddr).
					Str("rate_limit_key", key).
					Int("retry_after_seconds", seconds).
					Msg("Rate limit exceeded")

				common.RateLimitExceedResponse(app.Logger, localizer, w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func IPBasedRateLimiter(app *app.Application) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.Logger.Warn().
					Str("url", r.URL.String()).
					Str("method", r.Method).
					Str("remote_addr", r.RemoteAddr).
					Msg("Can not SplitHostPort and will use r.RemoteAddr as fallback")
				ip = r.RemoteAddr
			}

			key := fmt.Sprintf("rate_limit:%s", ip)
			res, err := app.Limiter.Allow(r.Context(), key, redis_rate.PerMinute(100))
			if err != nil {
				app.Logger.Error().
					Err(err).
					Str("url", r.URL.String()).
					Str("method", r.Method).
					Str("remote_addr", r.RemoteAddr).
					Str("rate_limit_key", key).
					Msg("Rate limiting error")
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}

			// app.Logger.Debug().
			// 	Str("url", r.URL.String()).
			// 	Str("method", r.Method).
			// 	Str("remote_addr", r.RemoteAddr).
			// 	Str("rate_limit_key", key).
			// 	Int("rate_limit_remaining", res.Remaining).
			// 	Msg("Rate limit check")

			w.Header().Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))

			if res.Allowed == 0 {
				// We are rate limited.
				seconds := int(res.RetryAfter / time.Second)
				w.Header().Set("RateLimit-RetryAfter", strconv.Itoa(seconds))

				app.Logger.Warn().
					Str("url", r.URL.String()).
					Str("method", r.Method).
					Str("remote_addr", r.RemoteAddr).
					Str("rate_limit_key", key).
					Int("retry_after_seconds", seconds).
					Msg("Rate limit exceeded")

				common.RateLimitExceedResponse(app.Logger, localizer, w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
