package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/markovidakovic/gdsi/server/internal/rest"
)

// main is the entry point for the application. It initializes the rest server,
// starts it in a separate goroutine, and listens for system termination signals.
// Upon receiving a signal, it gracefully shuts down the server and any associated resources.
func main() {
	// Create a new rest server
	srv, err := rest.NewServer()
	if err != nil {
		log.Fatalf("api server failed to start -> %v", err)
	}

	go func() {
		// Run the server in a separate goroutine so that the http.ListenAndServe
		// function does not block the rest of the execution in the main goroutine
		log.Printf("api server started on port %s\n", srv.Cfg.ApiPort)
		err = http.ListenAndServe(":"+srv.Cfg.ApiPort, srv.Rtr)
		if err != nil {
			log.Fatalf("api server failed to start -> %v", err)
		}
	}()

	// System interrupt signals channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block for signal
	<-stop
	log.Println("termination signal received, server shutting down...")

	// Graceful shutdown cancellation context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Server gracefull shutdown
	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("error during api server shutdown -> %v", err)
	}
}
