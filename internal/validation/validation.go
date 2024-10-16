package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru_RU"
	"github.com/go-playground/locales/tk_TM"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/shopspring/decimal"
)

var (
	EmailRX       = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PhoneNumberRX = regexp.MustCompile("^+[1-9][0-9]{7,14}$")
)

func NewValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("slug", validateSlug)
	validate.RegisterValidation("decimalpercent", validateDecimalPercent)
	validate.RegisterValidation("decimalgtezero", validateDecimalGTE)
	validate.RegisterValidation("password", validatePlainPassword)

	return validate
}

func NewUniversalTranslator() *ut.UniversalTranslator {
	// Create instances of different languages
	enLocale := en.New()    // English
	ruLocale := ru_RU.New() // Turkmen
	tkLocale := tk_TM.New() // Russian

	// Initialize the Universal Translator with available languages
	return ut.New(enLocale, enLocale, ruLocale, tkLocale)
}

func GetTranslator(uni *ut.UniversalTranslator, lang string) ut.Translator {
	trans, _ := uni.GetTranslator(lang)
	return trans
}

func RegisterTranslations(app *app.Application, trans ut.Translator, lang string) error {
	switch lang {
	case "ru_RU":
		return ru_translations.RegisterDefaultTranslations(app.Validator, trans)
	case "tk_TM":
		return RegisterDefaultTurkmenTranslations(app.Validator, trans)
	default:
		return en_translations.RegisterDefaultTranslations(app.Validator, trans)
	}
}

func validateSlug(fl validator.FieldLevel) bool {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return slugRegex.MatchString(fl.Field().String())
}

func validateDecimalPercent(fl validator.FieldLevel) bool {
	val := fl.Field().Interface().(decimal.Decimal)

	min := decimal.NewFromInt(0)
	max := decimal.NewFromInt(100)

	return val.GreaterThanOrEqual(min) && val.LessThanOrEqual(max)
}

func validateDecimalGTE(fl validator.FieldLevel) bool {
	val := fl.Field().Interface().(decimal.Decimal)

	min := decimal.NewFromInt(0)

	return val.GreaterThanOrEqual(min)
}

func validatePlainPassword(fl validator.FieldLevel) bool {
	passwordPlaintext := fl.Field().String()

	if len(passwordPlaintext) < 8 {
		return false
	}

	if len(passwordPlaintext) > 72 {
		return false
	}

	hasUpper, hasLower, hasNumber, isAscii := false, false, false, true
	for _, char := range passwordPlaintext {
		if char > unicode.MaxASCII {
			isAscii = false
			break
		}

		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		}
	}

	return isAscii && hasUpper && hasLower && hasNumber
}
