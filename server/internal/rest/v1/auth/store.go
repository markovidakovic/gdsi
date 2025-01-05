package auth

import (
	"context"

	"github.com/markovidakovic/gdsi/server/internal/db"
)

type store struct {
	db db.Conn
}

func (s *store) createNewAccount(ctx context.Context) (string, error) {
	return "new account created", nil
}

func newStore(db db.Conn) *store {
	var s = &store{
		db,
	}
	return s
}
