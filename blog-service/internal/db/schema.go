package db

import "database/sql"

func MustInitSchema(db *sql.DB) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS blogs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description_markdown TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  images_json JSON NULL,
  author_user_id VARCHAR(24) NOT NULL
);`)
	if err != nil {
		panic(err)
	}
}
