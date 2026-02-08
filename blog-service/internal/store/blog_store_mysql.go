package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"blog-service/internal/model"
)

type MySQLBlogStore struct {
	db *sql.DB
}

func NewMySQLBlogStore(db *sql.DB) *MySQLBlogStore {
	return &MySQLBlogStore{db: db}
}

func (s *MySQLBlogStore) GetAll(ctx context.Context) ([]model.BlogDto, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT 
			id, title, description_markdown, created_at, images_json, author_user_id
		FROM blogs
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.BlogDto

	for rows.Next() {
		var b model.BlogDto
		var imagesJSON sql.NullString

		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.DescriptionMarkdown,
			&b.CreatedAt,
			&imagesJSON,
			&b.AuthorUserID,
		)
		if err != nil {
			continue
		}

		if imagesJSON.Valid && imagesJSON.String != "" {
			_ = json.Unmarshal([]byte(imagesJSON.String), &b.Images)
		}

		out = append(out, b)
	}

	return out, nil
}

func (s *MySQLBlogStore) GetByID(ctx context.Context, id int64) (model.BlogDto, bool, error) {
	var b model.BlogDto
	var imagesJSON sql.NullString

	err := s.db.QueryRowContext(ctx, `
		SELECT 
			id, title, description_markdown, created_at, images_json, author_user_id
		FROM blogs
		WHERE id = ?`, id).
		Scan(
			&b.ID,
			&b.Title,
			&b.DescriptionMarkdown,
			&b.CreatedAt,
			&imagesJSON,
			&b.AuthorUserID,
		)

	if err == sql.ErrNoRows {
		return model.BlogDto{}, false, nil
	}
	if err != nil {
		return model.BlogDto{}, false, err
	}

	if imagesJSON.Valid && imagesJSON.String != "" {
		_ = json.Unmarshal([]byte(imagesJSON.String), &b.Images)
	}

	return b, true, nil
}

func (s *MySQLBlogStore) Create(ctx context.Context, req model.CreateBlogReq, authorID int64) (model.BlogDto, error) {
	now := time.Now()
	authorIDStr := strconv.FormatInt(authorID, 10)

	imagesStr := ""
	if len(req.Images) > 0 {
		imagesBytes, _ := json.Marshal(req.Images)
		imagesStr = string(imagesBytes)
	}

	res, err := s.db.ExecContext(ctx, `
		INSERT INTO blogs(title, description_markdown, created_at, images_json, author_user_id)
		VALUES(?,?,?,?,?)`,
		req.Title,
		req.DescriptionMarkdown,
		now,
		NullableString(imagesStr),
		authorIDStr,
	)
	if err != nil {
		return model.BlogDto{}, err
	}

	newID, _ := res.LastInsertId()

	return model.BlogDto{
		ID:                  newID,
		Title:               req.Title,
		DescriptionMarkdown: req.DescriptionMarkdown,
		CreatedAt:           now,
		Images:              req.Images,
		AuthorUserID:        authorIDStr,
	}, nil
}
