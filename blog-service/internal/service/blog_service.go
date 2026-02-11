package service

import (
	"context"
	"mime/multipart"
	"strconv"

	"blog-service/internal/model"
	"blog-service/internal/store"
)

type BlogService struct {
	store    store.BlogStore
	auth     *AuthClient
	uploader ImageUploader
}

type ImageUploader interface {
	UploadMany(ctx context.Context, files []*multipart.FileHeader) ([]string, error)
}

func NewBlogService(store store.BlogStore, auth *AuthClient, uploader ImageUploader) *BlogService {
	return &BlogService{store: store, auth: auth, uploader: uploader}
}

func (s *BlogService) GetAll(ctx context.Context) ([]model.BlogDto, error) {
	blogs, err := s.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if s.auth == nil {
		return blogs, nil
	}

	for i := range blogs {
		uid, err := strconv.ParseInt(blogs[i].AuthorUserID, 10, 64)
		if err != nil || uid <= 0 {
			continue
		}
		username, err := s.auth.GetUsernameByID(ctx, uid)
		if err != nil {
			continue
		}
		blogs[i].AuthorUsername = username
	}

	return blogs, nil
}

func (s *BlogService) GetByID(ctx context.Context, id int64) (model.BlogDto, bool, error) {
	b, found, err := s.store.GetByID(ctx, id)
	if err != nil || !found {
		return b, found, err
	}

	if s.auth != nil {
		uid, err := strconv.ParseInt(b.AuthorUserID, 10, 64)
		if err == nil && uid > 0 {
			username, err := s.auth.GetUsernameByID(ctx, uid)
			if err == nil {
				b.AuthorUsername = username
			}
		}
	}

	return b, true, nil
}

func (s *BlogService) Create(ctx context.Context, req model.CreateBlogReq, authorID int64) (model.BlogDto, error) {
	out, err := s.store.Create(ctx, req, authorID)
	if err != nil {
		return model.BlogDto{}, err
	}

	if s.auth != nil {
		username, err := s.auth.GetUsernameByID(ctx, authorID)
		if err == nil {
			out.AuthorUsername = username
		}
	}

	return out, nil
}

func (s *BlogService) UploadImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	if s.uploader == nil {
		return []string{}, nil
	}
	return s.uploader.UploadMany(ctx, files)
}
