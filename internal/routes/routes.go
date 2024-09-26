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

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestLogger(app))
	r.Use(middleware.GeneralRateLimiter(app))
	r.Use(middleware.IPBasedRateLimiter(app))
	r.Use(middleware.Recoverer(app))

	r.NotFound(middleware.NotFound(app.Logger))
	r.MethodNotAllowed(middleware.MethodNotAllowed(app.Logger))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", handlers.HealthcheckHandler(app))

		r.Route("/categories", func(r chi.Router) {
			r.Get("/", handlers.ListCategoriesPublicHandler(app))
			r.Get("/{id}", handlers.GetCategoryPublicHandler(app))
		})

		r.Route("/manager", func(r chi.Router) {
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", handlers.ListCategoriesManagerHandler(app))
				r.Post("/", handlers.CreateCategoryManagerHandler(app))
				r.Get("/{id}", handlers.GetCategoryManagerHandler(app))
				r.Put("/{id}", handlers.UpdateCategoryManagerHandler(app))
				r.Patch("/{id}", handlers.PartialUpdateCategoryManagerHandler(app))
				r.Delete("/{id}", handlers.DeleteCategoryManagerHandler(app))
			})
		})
	})

	// TODO: optimistic locking for the Products, Order tables

	return r
}
