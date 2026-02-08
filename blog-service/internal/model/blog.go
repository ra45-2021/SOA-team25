package model

import "time"

type CreateBlogReq struct {
	Title               string   `json:"title" binding:"required"`
	DescriptionMarkdown string   `json:"descriptionMarkdown" binding:"required"`
	Images              []string `json:"images"`
}

type BlogDto struct {
	ID                  int64     `json:"id"`
	Title               string    `json:"title"`
	DescriptionMarkdown string    `json:"descriptionMarkdown"`
	CreatedAt           time.Time `json:"createdAt"`
	Images              []string  `json:"images,omitempty"`
	AuthorUserID        string    `json:"authorUserId"`
}
