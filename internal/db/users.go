package db

import (
	"context"
	"errors"
)

type User struct {
	ID    string
	Email string
}

var ErrEmailConflict = errors.New("email_conflict")

type UserStore interface {
	CreateUser(ctx context.Context, email, passwordHash string) (User, error)
}

