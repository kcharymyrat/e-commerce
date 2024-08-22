package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Get("/v1/healthcheck", app.healthcheckHandler)
	router.Post("/v1/categories", app.createCategoryHandler)
	router.Get("/v1/categories/{id}", app.getCategoryHandler)

	return router
}
