package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// New: uygulama için tek bir DB bağlantısı oluşturur
func New() (*sql.DB, error) {
	// DSN'i sabit kullanıyoruz
	dsn := "postgres://approval_user:approval_pass@localhost:5433/approval_db?sslmode=disable"

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(time.Hour)

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
