package players

import "github.com/markovidakovic/gdsi/server/internal/db"

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}
