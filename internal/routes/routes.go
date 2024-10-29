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
	r.Use(middleware.LocalizationMiddleware(app))
	r.Use(middleware.GeneralRateLimiter(app))
	r.Use(middleware.IPBasedRateLimiter(app))
	r.Use(middleware.Recoverer(app))

	r.NotFound(middleware.NotFound(app.Logger))
	r.MethodNotAllowed(middleware.MethodNotAllowed(app.Logger))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", handlers.HealthcheckHandler(app))

		r.Route("/categories", func(r chi.Router) {
			r.Get("/", handlers.ListCategoriesPublicHandler(app))
			r.Get("/{slug}", handlers.GetCategoryPublicHandler(app))
		})

		r.Route("/languages", func(r chi.Router) {
			r.Get("/", handlers.ListLanguagesPublicHandler(app))
			r.Get("/{id}", handlers.GetLanguagePublicHandler(app))
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", handlers.GetUsePublicHandler(app))
		})

		r.Route("/admin", func(r chi.Router) {
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", handlers.ListCategoriesManagerHandler(app))
				r.Post("/", handlers.CreateCategoryManagerHandler(app))
				r.Get("/{slug}", handlers.GetCategoryManagerHandler(app))
				r.Put("/{slug}", handlers.UpdateCategoryManagerHandler(app))
				r.Patch("/{slug}", handlers.PartialUpdateCategoryManagerHandler(app))
				r.Delete("/{slug}", handlers.DeleteCategoryManagerHandler(app))
			})

			r.Route("/languages", func(r chi.Router) {
				r.Get("/", handlers.ListLanguagesManagerHandler(app))
				r.Post("/", handlers.CreateLanguageManagerHandler(app))
				r.Get("/{id}", handlers.GetLanguageManagerHandler(app))
				r.Put("/{id}", handlers.UpdateLanguageManagerHandler(app))
				r.Patch("/{id}", handlers.PartialUpdateLanguageManagerHandler(app))
				r.Delete("/{id}", handlers.DeleteLanguageManagerHandler(app))
			})

			r.Route("/translations", func(r chi.Router) {
				r.Get("/", handlers.ListTranslationsHandler(app))
				r.Post("/", handlers.CreateTranslationMangerHandler(app))
				r.Get("/{id}", handlers.GetTranslationHandler(app))
				r.Put("/{id}", handlers.UpdateTranslationHandler(app))
				r.Patch("/{id}", handlers.PartialUpdateTranslationHandler(app))
				r.Delete("/{id}", handlers.DeleteTranslationHandler(app))
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/", handlers.ListUsersAdminHandler(app))
				r.Post("/", handlers.CreateUserAdminHandler(app))
				r.Get("/{id}", handlers.GetUsersAdminHandler(app))
				r.Put("/{id}", handlers.UpdateUserAdminHandler(app))
				r.Patch("/{id}", handlers.PartialUpdateUserAdminHandler(app))
				r.Delete("/{id}", handlers.DeleteUserAdminHandler(app))
				r.Post("/login", handlers.LoginUserHandler(app))
				r.Post("/logout", handlers.LogoutUserHandler(app))
			})

			r.Route("/tokens", func(r chi.Router) {
				r.Post("/renew", handlers.RenewAccessTokenReqHandler(app))
				r.Post("/revoke", handlers.RevokeSessionByIDHandler(app))
			})
		})

		r.Route("/me", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", handlers.GetUserSelfHandler(app))
				r.Put("/{id}", handlers.UpdateUserSelfHandler(app))
				r.Patch("/{id}", handlers.PartialUpdateUserSelfHandler(app))
				r.Delete("/{id}", handlers.DeleteUserSelfHandler(app))
			})
		})

	})

	// TODO: optimistic locking for the Products, Order tables

	return r
}
