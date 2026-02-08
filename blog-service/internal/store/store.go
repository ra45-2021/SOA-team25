package store

import (
	"context"

	"blog-service/internal/model"
)

type BlogStore interface {
	GetAll(ctx context.Context) ([]model.BlogDto, error)
	GetByID(ctx context.Context, id int64) (model.BlogDto, bool, error) // bool = found
	Create(ctx context.Context, req model.CreateBlogReq, authorID int64) (model.BlogDto, error)
}

func NullableString(s string) interface{} {
	if len(s) == 0 {
		return nil
	}
	return s
}
