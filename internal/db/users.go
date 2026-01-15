package db

import (
	"context"
	"errors"
)

type User struct {
	ID        string
	Email     string
	Username  string
	AvatarURL *string
}

var ErrEmailConflict = errors.New("email_conflict")
var ErrUserNotFound = errors.New("user_not_found")

type UserWithPassword struct {
	ID           string
	Email        string
	PasswordHash string
}

type UserStore interface {
	CreateUser(ctx context.Context, email, passwordHash, username string, avatarURL *string) (User, error)
	FindUserByEmail(ctx context.Context, email string) (UserWithPassword, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	UpdateUserProfile(ctx context.Context, id, username string, avatarURL *string) (User, error)
}
