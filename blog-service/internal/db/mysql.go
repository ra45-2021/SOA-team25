package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func MustMySQL(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)

	for i := 0; i < 20; i++ {
		if err := db.Ping(); err == nil {
			return db
		}
		time.Sleep(1 * time.Second)
	}

	panic("mysql not ready")
}
