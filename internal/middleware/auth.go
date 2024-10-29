package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/auth"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func AuthMiddleware(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

			accessClaims, err := verifyClaimsFromAuthHeader(app, localizer, w, r)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Error parsing JWT")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsActive || accessClaims.IsBanned {
				app.Logger.Error().Err(err).Msg("user is not active or banned")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			ctx := context.WithValue(r.Context(), types.UserClaimsKey{}, accessClaims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func StaffAuthMiddleware(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

			accessClaims, err := verifyClaimsFromAuthHeader(app, localizer, w, r)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Error parsing JWT")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsActive || accessClaims.IsBanned {
				app.Logger.Error().Err(err).Msg("user is not active or banned")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsStaff {
				app.Logger.Error().Err(err).Msg("user is not staff")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			ctx := context.WithValue(r.Context(), types.UserClaimsKey{}, accessClaims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func AdminAuthMiddleware(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

			accessClaims, err := verifyClaimsFromAuthHeader(app, localizer, w, r)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Error parsing JWT")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsActive || accessClaims.IsBanned {
				app.Logger.Error().Err(err).Msg("user is not active or banned")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsAdmin {
				app.Logger.Error().Err(err).Msg("user is not admin")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			ctx := context.WithValue(r.Context(), types.UserClaimsKey{}, accessClaims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func SuperuserAuthMiddleware(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

			accessClaims, err := verifyClaimsFromAuthHeader(app, localizer, w, r)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Error parsing JWT")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsActive || accessClaims.IsBanned {
				app.Logger.Error().Err(err).Msg("user is not active or banned")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsSuperuser {
				app.Logger.Error().Err(err).Msg("user is not superuser")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			ctx := context.WithValue(r.Context(), types.UserClaimsKey{}, accessClaims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func SelfAuthMiddleware(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

			accessClaims, err := verifyClaimsFromAuthHeader(app, localizer, w, r)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Error parsing JWT")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			if !accessClaims.IsActive || accessClaims.IsBanned {
				app.Logger.Error().Err(err).Msg("user is not active or banned")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			id, err := common.ReadUUIDParam(r)
			if err != nil {
				common.BadRequestResponse(app.Logger, localizer, w, r, err)
				return
			}

			if accessClaims.UserID != id {
				app.Logger.Error().Err(err).Msg("user is not self")
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
				return
			}

			ctx := context.WithValue(r.Context(), types.UserClaimsKey{}, accessClaims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func verifyClaimsFromAuthHeader(
	app *app.Application, localizer *i18n.Localizer, w http.ResponseWriter, r *http.Request,
) (*auth.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		app.Logger.Warn().Str("url", r.URL.String()).Msg("Missing Authorization header")
		common.UnauthorizedResponse(app.Logger, localizer, w, r)
		return nil, fmt.Errorf("missing Authorization header")
	}
	access_token := strings.TrimPrefix(authHeader, "Bearer ")

	return auth.ParseJWT(access_token, app.Config.SecretKey, app.Logger)

}
