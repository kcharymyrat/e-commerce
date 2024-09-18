package routes

import (
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/handlers"
	"github.com/kcharymyrat/e-commerce/internal/middleware"
)

func Routes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.Recoverer(app))

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Route("/categories", func(r chi.Router) {
			r.Get("/", handlers.ListCategoriesHandler(app))
			r.Post("/", app.createCategoryHandler)
			r.Get("/{id}", app.getCategoryHandler)
			r.Put("/{id}", app.updateCategoryHandler)
			r.Patch("/{id}", app.partialUpdateCategoryHandler)
			r.Delete("/{id}", app.deleteCategoryHandler)
		})
	})

	// TODO: optimistic locking for the Products, Order tables

	return r
}
