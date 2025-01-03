package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
)

// server is a struct that holds configuration, database connection, and router
// necessary to run the API server
type server struct {
	Cfg *config.Config // Server configuration
	Rtr *chi.Mux       // Router for handling HTTP requests
	Db  db.Conn        // Database connection
}

// Shutdown gracefully shuts down the API server by closing relevant services and connections
// or resources. It is called when the server is shutting down or terminating, ensuring the services
// and connections are properly closed to prevent "hanging" or resource leaks.
func (s *server) Shutdown(ctx context.Context) error {
	// Close the db connection
	if s.Db != nil {
		if err := db.Disconnect(ctx, s.Db); err != nil {
			return fmt.Errorf("error closing the database connection: %v", err)
		}
	}

	log.Println("api server shutdown completed")

	return nil
}

// NewServer initializes a new API server by loading the configuration, connecting
// to the database, setting up the router, and attaching service endpoints. This function
// returns a pointer to the server object and any error encountered during initialization.
func NewServer() (*server, error) {
	var srv = &server{}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Connect the database
	db, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize router
	rtr := chi.NewRouter()

	srv.Cfg = cfg
	srv.Db = db
	srv.Rtr = rtr

	srv.Rtr.Use(middleware.Logger)
	srv.Rtr.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from gdsi api\n"))
	})

	return srv, nil
}
