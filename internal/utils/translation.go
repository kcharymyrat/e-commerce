package utils

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/app"
)

func GetTranslationMap(
	app *app.Application,
	entityID uuid.UUID,
	languageCode, fieldName string,
) (string, error) {
	tr, err := app.Repositories.Translations.GetByEntityIDLangCodeFieldName(
		entityID, languageCode, fieldName,
	)

	if err != nil {
		return "", err
	}

	return tr.TranslatedValue, nil
}
