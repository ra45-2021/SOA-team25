package service

import (
	"context"
	"errors"
	"time"
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
	tour.Status = 0
	tour.Price = 0
	return s.store.CreateTour(tour)
}

func (s *TourService) GetAuthorTours(authorID int64) ([]model.Tour, error) {
	return s.store.GetToursByAuthor(authorID)
}

func (s *TourService) AddCheckpoint(cp *model.Checkpoint) error {
	return s.store.CreateCheckpoint(cp)
}

func (s *TourService) GetAll() ([]model.Tour, error) {
    return s.store.GetAll()
}

func (s *TourService) PublishTour(ctx context.Context, tourID uint, distance float64, durations []model.TourDuration) error {
    tour, err := s.store.GetByID(tourID)
    if err != nil {
        return err
    }

    if tour.Name == "" || tour.Description == "" || len(tour.Tags) == 0 {
        return errors.New("Osnovni podaci (ime, opis, tagovi) nisu popunjeni!")
    }

    if len(tour.Checkpoints) < 2 {
        return errors.New("Tura mora imati bar dve ključne tačke!")
    }

    if len(durations) == 0 {
        return errors.New("Mora biti definisano bar jedno vreme prevoza!")
    }

    now := time.Now()
    tour.Status = 1
    tour.PublishedDateTime = &now
    tour.Distance = distance
    tour.Durations = durations

    return s.store.Update(tour)
}

func (s *TourService) ArchiveTour(ctx context.Context, tourID uint) error {
    tour, err := s.store.GetByID(tourID)
    if err != nil {
        return err
    }

    if tour.Status != 1 {
        return errors.New("samo objavljene ture se mogu arhivirati")
    }

    tour.Status = 2 
    return s.store.Update(tour)
}

func (s *TourService) ReactivateTour(ctx context.Context, tourID uint) error {
    tour, err := s.store.GetByID(tourID)
    if err != nil {
        return err
    }

    if tour.Status != 2 {
        return errors.New("samo arhivirane ture se mogu ponovo aktivirati")
    }

    tour.Status = 1 
    return s.store.Update(tour)
}