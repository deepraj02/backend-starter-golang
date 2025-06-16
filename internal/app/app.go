package app

import (
	"database/sql"
	"time"

	"log"
	"net/http"

	"github.com/deepraj02/go-postgres-starter/internal/store"
	"github.com/deepraj02/go-postgres-starter/internal/utils/json"
	"github.com/deepraj02/go-postgres-starter/internal/utils/logger"
	"github.com/deepraj02/go-postgres-starter/migrations"
)

type Application struct {
	Logger *logger.Logger
	DB     *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// logger := log.New(os.Stdout, "app:", log.Ldate|log.Ltime|log.Lshortfile)
	logger := initializeLogger()

	app := &Application{
		Logger: logger,
		DB:     pgDB,
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
		app.Logger.Error("Health check failed: %v", err)
		json.WriteJson(w, http.StatusInternalServerError, json.Envelope{"error": err.Error()})
		return
	}
	app.Logger.Info("Health check passed")
	time.Sleep(2 * time.Second)
	json.WriteJson(w, http.StatusOK, json.Envelope{"status": "Healthy"})
}
