package rest

import (
	"context"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markovidakovic/gdsi/server/internal/config"
	"github.com/markovidakovic/gdsi/server/internal/db"
	v1 "github.com/markovidakovic/gdsi/server/internal/rest/v1"
	httpSwagger "github.com/swaggo/http-swagger"
)

// server is a struct that holds configuration, database connection, and router
// necessary to run the API server
type server struct {
	Cfg *config.Config // Server configuration
	Db  db.Conn        // Database connection
	Rtr *chi.Mux       // Router for handling HTTP requests
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

// MountHandlers configures the API endpoints and applies middleware to the router.
func (s *server) MountHandlers() {
	// Middleware
	s.Rtr.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))
	s.Rtr.Use(middleware.Logger)
	s.Rtr.Use(middleware.AllowContentType("application/json"))
	s.Rtr.Use(middleware.CleanPath)
	s.Rtr.Use(middleware.NoCache)
	s.Rtr.Use(middleware.StripSlashes)
	s.Rtr.Use(middleware.Heartbeat("/"))

	// Mount all application handlers
	v1.MountHandlers(s.Cfg, s.Db, s.Rtr)

	s.Rtr.Get("/swagger/*", httpSwagger.WrapHandler)
}

// NewServer initialized a new instance of the server, loading its configuration, database connection,
// and initializing the router.
func NewServer() (*server, error) {
	var err error
	var srv = &server{}

	// Load config
	srv.Cfg, err = config.Load()
	if err != nil {
		return nil, err
	}

	// Connect the database
	srv.Db, err = db.Connect(srv.Cfg)
	if err != nil {
		return nil, err
	}

	// Initialize router
	srv.Rtr = chi.NewRouter()

	return srv, nil
}
