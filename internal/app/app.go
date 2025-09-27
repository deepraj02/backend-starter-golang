package app

import (
	"database/sql"

	"log"
	"net/http"

	"github.com/deepraj02/go-postgres-starter/internal/api"
	"github.com/deepraj02/go-postgres-starter/internal/middleware"
	"github.com/deepraj02/go-postgres-starter/internal/store"
	"github.com/deepraj02/go-postgres-starter/internal/utils/json"
	"github.com/deepraj02/go-postgres-starter/internal/utils/logger"
	"github.com/deepraj02/go-postgres-starter/migrations"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Logger      *logger.Logger
	DB          *sql.DB
	AuthHandler *api.AuthHandler
	Middleware  middleware.AuthMiddleware
	Redis       *redis.Client
}

func NewApplication() (*Application, error) {
	pgDB, redisClient, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// logger := log.New(os.Stdout, "app:", log.Ldate|log.Ltime|log.Lshortfile)
	logger := initializeLogger()
	authStore := store.NewPostgresAuthStore(pgDB)
	authHandler := api.NewAuthHandler(authStore, logger)
	authMiddleware := middleware.AuthMiddleware{AuthStore: authStore}
	app := &Application{
		Logger:      logger,
		DB:          pgDB,
		AuthHandler: authHandler,
		Middleware:  authMiddleware,
		Redis:       redisClient,
	}
	return app, nil
}

func initializeLogger() *logger.Logger {
	loggerInstance, err := logger.NewLogger("go-postgres-starter")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return loggerInstance
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := app.DB.Ping(); err != nil {
		app.Logger.Error("Health check failed: PostgreSQL unavailable", err)
		json.WriteJson(w, http.StatusInternalServerError, json.Envelope{
			"status": "Unhealthy",
			"error":  "Database connection failed",
		})
		return
	}
	ctx := r.Context()
	if err := app.Redis.Ping(ctx).Err(); err != nil {
		app.Logger.Error("Health check failed: Redis unavailable", err)
		json.WriteJson(w, http.StatusInternalServerError, json.Envelope{
			"status": "Unhealthy",
			"error":  "Redis connection failed",
		})
		return
	}

	app.Logger.Info("Health check passed")
	json.WriteJson(w, http.StatusOK, json.Envelope{
		"status": "Healthy",
		"services": map[string]string{
			"database": "connected",
			"Redis":    "connected",
		},
	})
}
