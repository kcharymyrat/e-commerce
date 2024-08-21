package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type envelope map[string]interface{}

func (app *application) readUUIDParam(r *http.Request) (uuid.UUID, error) {
	idStr := chi.URLParamFromCtx(r.Context(), "id") // eg: c303282d-f2e6-46ca-a04a-35d3d873712d

	idUUID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return idUUID, nil
}

func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// js, err := json.Marshal(data)
	js, err := json.MarshalIndent(data, "", "  ") // FIXME:Change this latter
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(js)

	return nil
}
