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

type server struct {
	Cfg *config.Config
	Db  *db.Conn
	Rtr *chi.Mux
}

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

// @title Gdsi API
// @version 1.0.0
// @description Documentation for the gdsi API

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the Bearer token in the format: Bearer token
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

	// mount v1 api endpoints
	s.Rtr.Route("/v1", v1.MountHandlers(s.Cfg, s.Db))

	s.Rtr.Get("/swagger/*", httpSwagger.WrapHandler)
}

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
