package common

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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
