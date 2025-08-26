package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, dsn string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Store{db: pool}, nil
}

func (s *Store) Close(ctx context.Context) {
	s.db.Close()
}

func (s *Store) CanUserAccess(ctx context.Context, docId uuid.UUID, userId string) (bool, string, error) {
	var role string

	err := s.db.QueryRow(ctx, "SELECT owner_id from documents where id=$1", docId).Scan(&role)
	if err == nil {
		if role == userId {
			return true, "owner", nil
		}
	} else {
		return false, "", err
	}

	var r string
	q := `select role from document_collaborators where doc_id=$1 and user_id=$2`
	err = s.db.QueryRow(ctx, q, docId, userId).Scan(&r)
	if err == nil {
		return true, r, nil
	}
	if err == sql.ErrNoRows {
		return false, "", nil
	}
	return false, "", err
}

func (s *Store) SaveSnapshot(ctx context.Context, docID uuid.UUID, version int64, state []byte) error {
	_, err := s.db.Exec(ctx, `insert into yjs_snapshots(doc_id, version, state) values($1,$2,$3)`, docID, version, state)
	return err
}
