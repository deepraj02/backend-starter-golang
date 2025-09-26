package routes

import (
	"github.com/deepraj02/go-postgres-starter/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)
	r.Post("/auth/register", app.AuthHandler.Register)
	r.Post("/auth/login", app.AuthHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)
		r.Get("/auth/profile", app.AuthHandler.Profile)
	})

	return r
}
