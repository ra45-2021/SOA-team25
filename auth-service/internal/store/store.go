package store

import (
	"context"

	"auth-service/internal/model"
)

type UserStore interface {
	EnsureIndexes(ctx context.Context) error
	InsertUser(ctx context.Context, u model.UserDoc) error
	FindByUsername(ctx context.Context, username string) (model.UserDoc, bool, error)
	NextUserID(ctx context.Context) (int64, error)
	FindByID(ctx context.Context, id int64) (model.UserDoc, bool, error)
}
