package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/config"
)

type Conn struct {
	*pgx.Conn
}

func Connect(c *config.Config) (*Conn, error) {
	// Construct the db connection string
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", c.DbDriver, c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbSslMode)

	ctx := context.Background()

	// Connect to the db using pgx
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	// Check if the connection is alive
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not ping the database: %v", err)
	}

	log.Println("database connected")

	return &Conn{Conn: conn}, nil
}

func Disconnect(ctx context.Context, conn *Conn) error {
	if conn != nil {
		return conn.Close(ctx)
	}
	return nil
}
