package players

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
	"github.com/markovidakovic/gdsi/server/params"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type store struct {
	db      *db.Conn
	queries struct {
		list     string
		findById string
		update   string
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
		return fmt.Errorf("failed to load list.sql -> %w", err)
	}
	findByIdBytes, err := sqlFiles.ReadFile("queries/find_by_id.sql")
	if err != nil {
		return fmt.Errorf("failed to load find_by_id.sql -> %w", err)
	}
	updateBytes, err := sqlFiles.ReadFile("queries/update.sql")
	if err != nil {
		return fmt.Errorf("failed to load update.sql -> %w", err)
	}

	s.queries.list = string(listBytes)
	s.queries.findById = string(findByIdBytes)
	s.queries.update = string(updateBytes)

	return nil
}

var allowedSortFields = map[string]string{
	"created_at": "player.created_at",
}

func (s *store) findPlayers(ctx context.Context, limit, offset int, sort *params.OrderBy) ([]PlayerModel, error) {
	if sort != nil && sort.IsValid(allowedSortFields) {
		s.queries.list += fmt.Sprintf("order by %s %s\n", allowedSortFields[sort.Field], sort.Direction)
	} else {
		s.queries.list += fmt.Sprintln("order by player.created_at desc")
	}

	var err error
	var rows pgx.Rows
	if limit >= 0 {
		s.queries.list += `limit $1 offset $2`
		rows, err = s.db.Query(ctx, s.queries.list, limit, offset)
	} else {
		rows, err = s.db.Query(ctx, s.queries.list)
	}

	if err != nil {
		return nil, failure.New("unable to find players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	defer rows.Close()

	var dest = []PlayerModel{}
	for rows.Next() {
		var pm PlayerModel
		err := pm.ScanRows(rows)
		if err != nil {
			return nil, failure.New("unable to find players", err)
		}

		dest = append(dest, pm)
	}

	if err = rows.Err(); err != nil {
		return nil, failure.New("unable to find players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return dest, nil
}

func (s *store) countPlayers(ctx context.Context) (int, error) {
	var count int
	sql := `select count(*) from player`
	err := s.db.QueryRow(ctx, sql).Scan(&count)
	if err != nil {
		return 0, failure.New("unable to count players", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}
	return count, nil
}

func (s *store) findPlayer(ctx context.Context, playerId string) (*PlayerModel, error) {
	var dest PlayerModel
	row := s.db.QueryRow(ctx, s.queries.findById, playerId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("player not found", err)
		}
		return nil, failure.New("unable to find player", err)
	}

	return &dest, nil
}

func (s *store) updatePlayer(ctx context.Context, tx pgx.Tx, playerId string, model UpdatePlayerRequestModel) (*PlayerModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest PlayerModel
	row := q.QueryRow(ctx, s.queries.update, model.Height, model.Weight, model.Handedness, model.Racket, playerId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("player for update not found", err)
		}
		return nil, failure.New("unable to update player", err)
	}

	return &dest, nil
}

// helper
func (s *store) checkPlayerOwnership(ctx context.Context, playerId, accountId string) (bool, error) {
	sql1 := `
		select exists (
			select 1 from player where id = $1 and account_id = $2
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, sql1, playerId, accountId).Scan(&exists)
	if err != nil {
		return false, failure.New("unable to check player ownership", fmt.Errorf("%w -> %v", failure.ErrInternal, err))
	}

	return exists, nil
}
