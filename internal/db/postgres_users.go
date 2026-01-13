package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserStore struct {
	pool *pgxpool.Pool
}

func NewPostgresUserStore(pool *pgxpool.Pool) *PostgresUserStore {
	return &PostgresUserStore{pool: pool}
}

func (s *PostgresUserStore) CreateUser(ctx context.Context, email, passwordHash string) (User, error) {
	id, err := NewUUIDv4()
	if err != nil {
		return User{}, err
	}

	const q = `
INSERT INTO users (id, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email
`
	var u User
	if err := s.pool.QueryRow(ctx, q, id, email, passwordHash).Scan(&u.ID, &u.Email); err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrEmailConflict
		}
		return User{}, err
	}
	return u, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

