package common

import "net/http"

func GetAcceptLanguage(r *http.Request) string {
	lang := r.Header.Get("Accept-Language")
	switch {
	case lang == "tk":
		return "tk_TM"
	case lang == "ru":
		return "ru_RU"
	default:
		return "en"
	}
}
