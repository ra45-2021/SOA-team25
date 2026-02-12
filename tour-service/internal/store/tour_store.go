package store

import (
	"tour-service/internal/model"
	"gorm.io/gorm"
)

type TourStore struct {
	db *gorm.DB
}

func NewTourStore(db *gorm.DB) *TourStore {
	return &TourStore{db: db}
}

func (s *TourStore) CreateTour(tour *model.Tour) error {
	return s.db.Create(tour).Error
}

func (s *TourStore) GetToursByAuthor(authorID int64) ([]model.Tour, error) {
	var tours []model.Tour
	return tours, s.db.Preload("Checkpoints").Where("author_id = ?", authorID).Find(&tours).Error
}

func (s *TourStore) CreateCheckpoint(cp *model.Checkpoint) error {
	return s.db.Create(cp).Error
}

func (s *TourStore) GetAll() ([]model.Tour, error) {
    var tours []model.Tour
    
    result := s.db.Where("status <> ?", 2).Find(&tours)
    
    return tours, result.Error
}