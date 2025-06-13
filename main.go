package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/deepraj02/go-postgres-starter/internal/app"
	"github.com/deepraj02/go-postgres-starter/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	defer app.Logger.Close()
	app.Logger.Info("Starting application on port http://localhost:%d", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Error("Error starting server: %v", err)
		return
	}
}
