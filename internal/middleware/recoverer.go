package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
)

// Recoverer returns a middleware function with injected app.Application.
func Recoverer(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					app.Logger.Error().
						Interface("panic", err).
						Bytes("stack", debug.Stack()).
						Str("method", r.Method).
						Str("url", r.URL.Path).
						Str("remote_addr", r.RemoteAddr).
						Str("user_agent", r.UserAgent()).
						Msg("panic occurred during request")
					common.ErrorResponse(app.Logger, w, r, http.StatusInternalServerError, err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
