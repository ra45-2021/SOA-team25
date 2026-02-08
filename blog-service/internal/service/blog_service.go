package service

import (
	"context"

	"blog-service/internal/model"
	"blog-service/internal/store"
)

type BlogService struct {
	store store.BlogStore
}

func NewBlogService(store store.BlogStore) *BlogService {
	return &BlogService{store: store}
}

func (s *BlogService) GetAll(ctx context.Context) ([]model.BlogDto, error) {
	return s.store.GetAll(ctx)
}

func (s *BlogService) GetByID(ctx context.Context, id int64) (model.BlogDto, bool, error) {
	return s.store.GetByID(ctx, id)
}

func (s *BlogService) Create(ctx context.Context, req model.CreateBlogReq, authorID int64) (model.BlogDto, error) {
	return s.store.Create(ctx, req, authorID)
}
