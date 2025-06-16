package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
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
	defer app.Logger.Close() // Move this here to ensure it's always called

	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Channel to capture server errors
	serverErrors := make(chan error, 1)

	go func() {
		app.Logger.Info("Starting application on port http://localhost:%d", port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Wait for either shutdown signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			app.Logger.Error("Error starting server: %v", err)
			return
		}
	case <-shutdown.Done():
		app.Logger.Info("Shutdown signal received")
	}

	app.Logger.Info("Shutting down Server....")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.Logger.Error("Error shutting down server: %v", err)
		// Force close if graceful shutdown fails
		server.Close()
	} else {
		app.Logger.Info("Server shutdown gracefully")
	}
}
