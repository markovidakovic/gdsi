package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/internal/config"
)

// Conn is a wrapper arrount pgx.Conn type to enable method promotion.
// By embedding pgx.Conn, we can directly access its methods on the Conn type.
type Conn struct {
	*pgx.Conn
}

// Connect establishes a connection to the PostgreSQL database using the pgx library.
// It constructs a connection string using the provided configuration and attempts to connect to the database.
// If successful, it pings the database to check if the connection is alive and returns the connection object.
func Connect(c *config.Config) (*Conn, error) {
	// Construct the db connection string
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", c.DbDriver, c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbSslMode)

	// Background context
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

// Disconnect gracefully closes the database connection. It checks if the connection is not nil
// and if so, attempts to close it using the provided context.
func Disconnect(ctx context.Context, conn *Conn) error {
	if conn != nil {
		return conn.Close(ctx)
	}
	return nil
}
