package model

import "time"

type TourDifficulty int

const (
	Easy TourDifficulty = iota
	Medium
	Hard
)

type TourStatus int

const (
    Draft TourStatus = iota
    Published
    Archived
)

type TransportType int

const (
    Walking TransportType = iota
    Bicycle
    Car
)

type TourDuration struct {
    Minutes       int           `json:"minutes"`
    TransportType TransportType `json:"transportType"`
}


type Tour struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AuthorID    int64          `json:"author_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Difficulty  TourDifficulty `json:"difficulty"`
	Tags        string         `json:"tags"`
	Status      TourStatus     `json:"status"`
	Distance    float64        `json:"distance"` 
    PublishedDateTime *time.Time `json:"publishedDateTime"` 
	Price       float64        `json:"price"`  // uvek 0
	CreatedAt   time.Time      `json:"created_at"`
	Checkpoints []Checkpoint   `json:"checkpoints" gorm:"foreignKey:TourID"`
	Durations   []TourDuration `json:"durations" gorm:"serializer:json"`
}

type Checkpoint struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	TourID      uint    `json:"tour_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageURL    string  `json:"image_url"`
}

