package routes

import (
	"net/http"
	"path/filepath"

	"github.com/deepraj02/go-postgres-starter/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.AuthHandler.Register)
		r.Post("/login", app.AuthHandler.Login)
		r.Post("/forgot-password", app.AuthHandler.ForgotPassword)
		r.Post("/reset-password", app.AuthHandler.ResetPassword)

		r.Group(func(r chi.Router) {
			r.Use(app.Middleware.Authenticate)
			r.Get("/profile", app.AuthHandler.Profile)
		})
	})

	workDir, _ := filepath.Abs("./")
	staticDir := http.Dir(filepath.Join(workDir, "static"))
	r.Handle("/*", http.StripPrefix("/", http.FileServer(staticDir)))

	return r
}
