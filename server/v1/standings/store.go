package standings

import (
	"context"
	"embed"
	"fmt"

	"github.com/markovidakovic/gdsi/server/db"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type store struct {
	db      *db.Conn
	queries struct {
		list string
	}
}

func newStore(db *db.Conn) (*store, error) {
	s := &store{
		db: db,
	}
	if err := s.loadQueries(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *store) loadQueries() error {
	listBytes, err := sqlFiles.ReadFile("queries/list.sql")
	if err != nil {
		return fmt.Errorf("failed to read list.sql -> %v", err)
	}

	s.queries.list = string(listBytes)

	return nil
}

func (s *store) findStandings(ctx context.Context, seasonId, leagueId string) ([]StandingModel, error) {
	dest := []StandingModel{}

	rows, err := s.db.Query(ctx, s.queries.list, seasonId, leagueId)
	if err != nil {
		return nil, fmt.Errorf("quering standing rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sm StandingModel
		err := sm.ScanRows(rows)
		if err != nil {
			return nil, err
		}

		dest = append(dest, sm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating standing rows: %v", err)
	}

	return dest, nil
}
