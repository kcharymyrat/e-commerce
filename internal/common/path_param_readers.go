package common

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/constants"
)

func ReadUUIDParam(r *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParam(r, "id") // eg: c303282d-f2e6-46ca-a04a-35d3d873712d

	idUUID, err := uuid.Parse(idStr)
	if err != nil {
		// FIXME: Maybe need localiation errors ??
		return uuid.Nil, err
	}

	return idUUID, nil
}

func ReadSlugParam(r *http.Request) (string, error) {
	slug := chi.URLParam(r, "slug")

	slugRegex, err := regexp.Compile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	if err != nil {
		return "", err
	}

	if !slugRegex.MatchString(slug) {
		return "", errors.New(constants.InvalidIDErrMsg)
	}

	return slug, nil
}
