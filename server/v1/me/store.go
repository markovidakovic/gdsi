package me

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/markovidakovic/gdsi/server/db"
	"github.com/markovidakovic/gdsi/server/failure"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type store struct {
	db      *db.Conn
	queries struct {
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
	findByIdBytes, err := sqlFiles.ReadFile("queries/find_by_id.sql")
	if err != nil {
		return fmt.Errorf("failed to load find_by_id.sql -> %w", err)
	}
	updateBytes, err := sqlFiles.ReadFile("queries/update.sql")
	if err != nil {
		return fmt.Errorf("failed to load update.sql -> %w", err)
	}

	s.queries.findById = string(findByIdBytes)
	s.queries.update = string(updateBytes)

	return nil
}

func (s *store) findMe(ctx context.Context, accountId string) (*MeModel, error) {
	var dest MeModel
	row := s.db.QueryRow(ctx, s.queries.findById, accountId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("account not found", err)
		}
		return nil, failure.New("unable to find account", err)
	}

	return &dest, nil
}

func (s *store) updateMe(ctx context.Context, tx pgx.Tx, accountId string, model UpdateMeRequestModel) (*MeModel, error) {
	var q db.Querier
	if tx != nil {
		q = tx
	} else {
		q = s.db
	}

	var dest MeModel
	row := q.QueryRow(ctx, s.queries.update, model.Name, accountId)
	err := dest.ScanRow(row)
	if err != nil {
		if errors.Is(err, failure.ErrNotFound) {
			return nil, failure.New("account for updating not found", err)
		}
		return nil, failure.New("unable to update account", err)
	}

	return &dest, nil
}
