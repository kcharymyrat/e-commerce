package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kcharymyrat/e-commerce/internal/app"
)

func RegisterCustomEnTranslations(app *app.Application, trans ut.Translator) {
	app.Validator.RegisterTranslation("slug", trans, func(ut ut.Translator) error {
		return ut.Add("slug", "{0} must be a valid slug", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("slug", fe.Field())
		return t
	})

	app.Validator.RegisterTranslation("decimalpercent", trans, func(ut ut.Translator) error {
		return ut.Add("decimalpercent", "{0} must be between 0.00 and 100.00", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("decimalpercent", fe.Field())
		return t
	})

	// Register for other fields like decimalgtezero similarly
	app.Validator.RegisterTranslation("decimalgtezero", trans, func(ut ut.Translator) error {
		return ut.Add("decimalgtezero", "{0} must be greater than or equal to 0.00", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("decimalgtezero", fe.Field())
		return t
	})
}
