package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/shopspring/decimal"
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

func WriteJson(
	w http.ResponseWriter,
	status int,
	data types.Envelope,
	headers http.Header,
) error {
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

// TODO: Tranlate the errors
func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)

	if err != nil {

		fmt.Printf("%v, %T\n", err.Error(), err)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: uknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err

		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func ReadQueryStr(qs url.Values, key string) *string {
	s := qs.Get(key)
	if s == "" {
		return nil
	}
	return &s
}

func ReadQueryCSStrs(qs url.Values, key string) []string {
	s := qs.Get(key)
	if s == "" {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(strings.ToLower(s)), ",")
}

func ReadQueryUUID(qs url.Values, key string) *uuid.UUID {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	qsUUID, err := uuid.Parse(s)
	if err != nil {
		return nil
	}

	return &qsUUID
}

func ReadQueryCSUUIDs(qs url.Values, key string) []uuid.UUID {
	s := qs.Get(key)
	if s == "" {
		return []uuid.UUID{}
	}

	uuids := []uuid.UUID{}

	for _, s := range strings.Split(s, ",") {
		trimmedS := strings.TrimSpace(s)
		qsUUID, err := uuid.Parse(trimmedS)
		if err != nil {
			break
		}
		uuids = append(uuids, qsUUID)
	}

	return uuids
}

func ReadQueryInt(qs url.Values, key string) *int {
	s := qs.Get(key)

	if s == "" {
		return nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}

	return &i
}

func ReadQueryDecimal(qs url.Values, key string) *decimal.Decimal {
	s := qs.Get(key)

	if s == "" {
		return nil
	}

	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}

	return &d
}

func ReadQueryTime(qs url.Values, key string) *time.Time {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	qsTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}

	return &qsTime
}

func ReadQueryBool(qs url.Values, key string) *bool {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}

	return &b
}
