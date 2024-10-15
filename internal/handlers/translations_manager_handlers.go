package handlers

import (
	"errors"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func CreateTranslationMangerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		var input requests.CreateTranslationInput
		err := common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		tr := mappers.CreateTranslationInputToTranslationMapper(&input)
		err = app.Validator.Struct(tr)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.CreateTranslationService(app, tr)
		if err != nil {
			HandlePGErrors(app.Logger, localizer, w, r, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", "api/v1/translations/%s")

		trResponse := mappers.TranslationToTranslationManagerResponseMappper(tr)

		err = common.WriteJson(w, http.StatusCreated, types.Envelope{"translation": trResponse}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func GetTranslationHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		tr, err := services.GetTranslationService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.BadRequestResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		trResponse := mappers.TranslationToTranslationManagerResponseMappper(tr)
		err = common.WriteJson(w, http.StatusOK, types.Envelope{"translation": trResponse}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}
	}
}

func ListTranslationsHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		input := requests.ListTranslationsInput{}
		readTranslationQueryParameters(&input, r.URL.Query())
		err := app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		trList, metadata, err := services.ListTranslationsService(app, &input)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		trListRes := make([]*responses.TranslationManagerResponse, len(trList))
		for _, tr := range trList {
			trRes := mappers.TranslationToTranslationManagerResponseMappper(tr)
			trListRes = append(trListRes, trRes)
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{
			"metadata": metadata,
			"results":  trListRes,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}
	}
}

func UpdateTranslationHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		input := requests.UpdateTranslationInput{}
		err = common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		tr, err := services.GetTranslationService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
				return
			}
		}

		tr.LanguageCode = input.LanguageCode
		tr.EntityID = input.EntityID
		tr.TableName = input.TableName
		tr.FieldName = input.FieldName
		tr.TranslatedFieldName = input.TranslatedFieldName
		tr.TranslatedValue = input.TranslatedValue
		tr.UpdatedByID = input.UpdatedByID

		err = app.Validator.Struct(tr)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.UpdateTranslationService(app, tr)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{"translation": tr}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func PartialUpdateTranslationHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		input := requests.PartialUpdateTranslationInput{}
		err = common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		tr, err := services.GetTranslationService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
				return
			}
		}

		if input.LanguageCode != nil {
			tr.LanguageCode = *input.LanguageCode
		}
		if input.EntityID != nil {
			tr.EntityID = *input.EntityID
		}
		if input.TableName != nil {
			tr.TableName = *input.TableName
		}
		if input.FieldName != nil {
			tr.FieldName = *input.FieldName
		}
		if input.TranslatedFieldName != nil {
			tr.TranslatedFieldName = *input.TranslatedFieldName
		}
		if input.TranslatedValue != nil {
			tr.TranslatedValue = *input.TranslatedValue
		}
		tr.UpdatedByID = input.UpdatedByID

		err = app.Validator.Struct(tr)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.UpdateTranslationService(app, tr)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{"translation": tr}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func DeleteTranslationMapper(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = services.DeleteTranslationService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{"message": "translation successfully deleted"}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}
