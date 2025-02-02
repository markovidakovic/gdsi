package seasons

import (
	"context"

	"github.com/markovidakovic/gdsi/server/db"
)

type store struct {
	db *db.Conn
}

func newStore(db *db.Conn) *store {
	return &store{
		db,
	}
}

func insertSeason(ctx context.Context, input CreateSeasonRequestModel) {

}
