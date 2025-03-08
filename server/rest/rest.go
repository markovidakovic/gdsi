package rest

import (
	"context"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markovidakovic/gdsi/server/config"
	"github.com/markovidakovic/gdsi/server/db"
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
func (s *server) MountRouters() error {
	s.setupMiddleware()

	v1, err := v1.New(s.Cfg, s.Db)
	if err != nil {
		return err
	}

	// mount v1
	s.Rtr.Route("/v1", v1.Mount)

	if s.swaggerEnabled {
		s.Rtr.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	return nil
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
		MaxAge:         300, // maximum value not ignored by any of major browsers
	}))
	s.Rtr.Use(chimiddleware.Logger)
	s.Rtr.Use(chimiddleware.AllowContentType("application/json"))
	s.Rtr.Use(chimiddleware.CleanPath)
	s.Rtr.Use(chimiddleware.NoCache)
	s.Rtr.Use(chimiddleware.StripSlashes)
	s.Rtr.Use(chimiddleware.Heartbeat("/"))
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
