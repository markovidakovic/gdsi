package players

import (
	"context"

	"github.com/markovidakovic/gdsi/server/internal/db"
)

type store struct {
	db db.Conn
}

func (s *store) queryPlayers(ctx context.Context) (string, error) {
	return "queried plyers from store", nil
}

func newStore(db db.Conn) *store {
	var s = &store{
		db,
	}
	return s
}
