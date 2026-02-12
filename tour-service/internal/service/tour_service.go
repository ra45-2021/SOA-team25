package service

import (
	"context"
	"errors"
	"tour-service/internal/client"
	"tour-service/internal/model"
	"tour-service/internal/store"
)

type TourService struct {
	store      *store.TourStore
	authClient *client.AuthClient
}

func NewTourService(s *store.TourStore, ac *client.AuthClient) *TourService {
	return &TourService{store: s, authClient: ac}
}

func (s *TourService) CreateTour(ctx context.Context, tour *model.Tour) error {
	isGuide, err := s.authClient.IsGuide(ctx, tour.AuthorID)
	if err != nil || !isGuide {
		return errors.New("user is not a guide or not found")
	}
	tour.Status = "DRAFT"
	tour.Price = 0
	return s.store.CreateTour(tour)
}

func (s *TourService) GetAuthorTours(authorID int64) ([]model.Tour, error) {
	return s.store.GetToursByAuthor(authorID)
}

func (s *TourService) AddCheckpoint(cp *model.Checkpoint) error {
	return s.store.CreateCheckpoint(cp)
}