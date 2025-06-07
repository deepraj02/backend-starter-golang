package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/deepraj02/go-postgres-starter/internal/store"
	"github.com/deepraj02/go-postgres-starter/internal/utils"
	"github.com/deepraj02/go-postgres-starter/migrations"
)

type Application struct {
	Logger *log.Logger
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
	logger := log.New(os.Stdout, "app:", log.Ldate|log.Ltime|log.Lshortfile)
	app := &Application{
		Logger: logger,
		DB:     pgDB,
	}
	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := app.DB.Ping(); err != nil {
		app.Logger.Printf("Health check failed: %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}
	app.Logger.Println("Health check passed")
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"status": "Healthy"})
}
