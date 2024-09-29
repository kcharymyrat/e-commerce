package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kcharymyrat/e-commerce/internal/app"
)

func RegisterCustomTkTranslations(app *app.Application, trans ut.Translator) {
	app.Validator.RegisterTranslation("slug", trans, func(ut ut.Translator) error {
		return ut.Add("slug", "{0} dogry slug bolmaly", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("slug", fe.Field())
		return t
	})

	app.Validator.RegisterTranslation("decimalpercent", trans, func(ut ut.Translator) error {
		return ut.Add("decimalpercent", "{0} 0.00 bilen 100.00 aralygynda bolmaly", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("decimalpercent", fe.Field())
		return t
	})

	app.Validator.RegisterTranslation("decimalgtezero", trans, func(ut ut.Translator) error {
		return ut.Add("decimalgtezero", "{0} 0.00-den uly ýa-da deň bolmaly", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("decimalgtezero", fe.Field())
		return t
	})
}
