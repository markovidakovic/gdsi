package rest

import (
	"context"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	gochimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/middleware"
	v1 "github.com/markovidakovic/gdsi/server/v1"
	httpSwagger "github.com/swaggo/http-swagger"
)

type server struct {
	Cfg            *config.Config
	Db             *db.Conn
	Rtr            *chi.Mux
	swaggerEnabled bool
}

type serverOption func(*server) error

func NewServer() (*server, error) {
	var srv = &server{}

	opts := []serverOption{
		withConfig(),
		withDatabase(),
		withRouter(),
		withSwagger(),
	}

	for _, opt := range opts {
		if err := opt(srv); err != nil {
			return nil, err
		}
	}

	return srv, nil
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
func (s *server) MountRouters() {
	s.setupMiddleware()

	// mount v1 api endpoints
	s.Rtr.Route("/v1", v1.MountRouter(s.Cfg, s.Db))

	if s.swaggerEnabled {
		s.Rtr.Get("/swagger/*", httpSwagger.WrapHandler)
	}
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

func (s *server) setupMiddleware() {
	s.Rtr.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))
	// s.Rtr.Use(gochimiddleware.Logger)
	s.Rtr.Use(middleware.Logger)
	s.Rtr.Use(gochimiddleware.AllowContentType("application/json"))
	s.Rtr.Use(gochimiddleware.CleanPath)
	s.Rtr.Use(gochimiddleware.NoCache)
	s.Rtr.Use(gochimiddleware.StripSlashes)
	s.Rtr.Use(gochimiddleware.Heartbeat("/"))
}

func withConfig() serverOption {
	return func(s *server) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		s.Cfg = cfg
		return nil
	}
}

func withDatabase() serverOption {
	return func(s *server) error {
		if s.Cfg == nil {
			return fmt.Errorf("config must be initialized before database")
		}
		db, err := db.Connect(s.Cfg)
		if err != nil {
			return err
		}
		s.Db = db
		return nil
	}
}

func withRouter() serverOption {
	return func(s *server) error {
		s.Rtr = chi.NewRouter()
		return nil
	}
}

func withSwagger() serverOption {
	return func(s *server) error {
		s.swaggerEnabled = true
		return nil
	}
}
