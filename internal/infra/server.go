package infra

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serverPort = ":8080"

func startHTTPServer(app *App) *http.Server {
	server := &http.Server{
		Addr:    serverPort,
		Handler: app.MainRouter,
	}

	go func() {
		log.Printf("server starting on %s", serverPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	}()

	return server
}

func waitForShutdownSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
}

func shutdown(server *http.Server, app *App) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
		return
	}

	if err := app.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown tracer: %v", err)
	}
}
