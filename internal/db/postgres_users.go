package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserStore struct {
	pool *pgxpool.Pool
}

func NewPostgresUserStore(pool *pgxpool.Pool) *PostgresUserStore {
	return &PostgresUserStore{pool: pool}
}

func (s *PostgresUserStore) CreateUser(ctx context.Context, email, passwordHash, username string, avatarURL *string) (User, error) {
	id, err := NewUUIDv4()
	if err != nil {
		return User{}, err
	}

	const q = `
INSERT INTO users (id, email, password_hash, username, avatar_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, username, avatar_url
`
	var u User
	if err := s.pool.QueryRow(ctx, q, id, email, passwordHash, username, avatarURL).Scan(&u.ID, &u.Email, &u.Username, &u.AvatarURL); err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrEmailConflict
		}
		return User{}, err
	}
	return u, nil
}

func (s *PostgresUserStore) FindUserByEmail(ctx context.Context, email string) (UserWithPassword, error) {
	const q = `
SELECT id, email, password_hash
FROM users
WHERE email = $1
`
	var u UserWithPassword
	if err := s.pool.QueryRow(ctx, q, email).Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserWithPassword{}, ErrUserNotFound
		}
		return UserWithPassword{}, err
	}
	return u, nil
}

func (s *PostgresUserStore) GetUserByID(ctx context.Context, id string) (User, error) {
	const q = `
SELECT id, email, username, avatar_url
FROM users
WHERE id = $1
`
	var u User
	if err := s.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Email, &u.Username, &u.AvatarURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

func (s *PostgresUserStore) UpdateUserProfile(ctx context.Context, id, username string, avatarURL *string) (User, error) {
	const q = `
UPDATE users
SET username = $2,
    avatar_url = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, username, avatar_url
`
	var u User
	if err := s.pool.QueryRow(ctx, q, id, username, avatarURL).Scan(&u.ID, &u.Email, &u.Username, &u.AvatarURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
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
