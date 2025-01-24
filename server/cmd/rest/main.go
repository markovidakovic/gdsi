package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/markovidakovic/gdsi/server/docs"
	"github.com/markovidakovic/gdsi/server/rest"
)

func main() {
	srv, err := rest.NewServer()
	if err != nil {
		log.Fatalf("api server failed to start -> %v", err)
	}
	srv.MountRouters()

	go func() {
		// Run server in a separate goroutine
		log.Printf("api server started on port %s\n", srv.Cfg.ApiPort)
		err = http.ListenAndServe(":"+srv.Cfg.ApiPort, srv.Rtr)
		if err != nil {
			log.Fatalf("api server failed to start -> %v", err)
		}
	}()

	// System interrupt signals channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Block for signal
	<-stop
	log.Println("termination signal received, server shutting down...")

	// Graceful shutdown cancellation context
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Server gracefull shutdown
	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during api server shutdown -> %v", err)
	}
}
