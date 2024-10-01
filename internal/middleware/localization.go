package middleware

import (
	"context"
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func LocalizationMiddleware(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Detect the language from the Accept-Language header
			lang := common.GetAcceptLanguage(r) // This function detects language from header
			localizer := i18n.NewLocalizer(app.I18nBundle, lang)
			valTrans, _ := app.ValUniTrans.GetTranslator(lang)

			// Attach the localizer to the context
			ctx := context.WithValue(r.Context(), constants.LocalizerKey, localizer)
			ctx = context.WithValue(ctx, constants.ValTransKey, valTrans)
			r = r.WithContext(ctx)

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
