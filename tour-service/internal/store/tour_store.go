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

func (s *TourStore) GetAll() ([]model.Tour, error) {
    var tours []model.Tour
    err := s.db.Preload("Checkpoints").Where("status = ?", 1).Find(&tours).Error
    return tours, err
}

func (s *TourStore) GetToursByAuthor(authorID int64) ([]model.Tour, error) {
    var tours []model.Tour
    err := s.db.Preload("Checkpoints").Where("author_id = ?", authorID).Find(&tours).Error
    return tours, err
}

func (s *TourStore) GetByID(id uint) (*model.Tour, error) {
    var tour model.Tour
    err := s.db.Preload("Checkpoints").First(&tour, id).Error
    if err != nil {
        return nil, err
    }
    return &tour, nil
}

func (s *TourStore) CreateCheckpoint(cp *model.Checkpoint) error {
    return s.db.Create(cp).Error
}

func (s *TourStore) Update(tour *model.Tour) error {
    return s.db.Save(tour).Error
}

func (s *TourStore) UpdateStatus(id uint, status int) error {
    return s.db.Model(&model.Tour{}).Where("id = ?", id).Update("status", status).Error
}