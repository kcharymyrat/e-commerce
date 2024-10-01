package handlers

import (
	"errors"
	"fmt"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func CreateLanguageHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		input := requests.CreateLanguageInput{}
		err := common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		language := mappers.CreateLanguageInputToLanguageMapper(&input)
		err = app.Validator.Struct(language)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.CreateLanguageService(app, language)
		if err != nil {
			HandlePGErrors(app.Logger, localizer, w, r, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("api/v1/%s", language.ID))

		languageResponse := mappers.LanguageToLanguageManagerResponseMapper(language)

		err = common.WriteJson(w, http.StatusCreated, common.Envelope{"language": languageResponse}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func GetLanguageHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		language, err := services.GetLanguageService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		languageResponse := mappers.LanguageToLanguageManagerResponseMapper(language)

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"language": languageResponse}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func ListLanguagesHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		input := requests.ListLanguagesInput{}

		qs := r.URL.Query()
		readLanguageQueryParameters(&input, qs)
		err := app.Validator.Struct(&input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		languages, metadata, err := services.ListLanguagesService(app, input)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		languagesResponse := make([]*responses.LanguageManagerResponse, len(languages))
		for _, language := range languages {
			res := mappers.LanguageToLanguageManagerResponseMapper(language)
			languagesResponse = append(languagesResponse, res)
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{
			"metadata": metadata,
			"results":  languagesResponse,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func UpdateLanguageHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		input := requests.UpdateLanguageInput{}
		err = common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			transErrs := make(map[string]string)
			for _, e := range errs {
				transErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, transErrs)
			return
		}

		language, err := services.GetLanguageService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrEditConflict):
				common.EditConflictResponse(app.Logger, localizer, w, r)
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}
		}

		err = services.UpdateLanguageService(app, &input, language)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"language": language}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func PartialUpdateLanguageHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
		}

		input := requests.PartialUpdateLanguageInput{}
		err = common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			transErrs := make(map[string]string)
			for _, e := range errs {
				transErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, transErrs)
			return
		}

		language, err := services.GetLanguageService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = services.PartialUpdateLanguageService(app, &input, language)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"language": language}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func DeleteLanguageHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = services.DeleteLanguageService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		// TODO: Needs localiztions
		err = common.WriteJson(w, http.StatusOK, common.Envelope{"message": "language successfully deleted"}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}
