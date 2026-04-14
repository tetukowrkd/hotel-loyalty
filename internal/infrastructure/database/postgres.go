package database

import (
	"database/sql"
	"e-commerce/internal/config"
	"e-commerce/internal/infrastructure/logger"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.ErrorLogger.Println("Failed to open DB:", err)
		return nil, err
	}

	// connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		logger.ErrorLogger.Println("Failed to ping DB:", err)
		return nil, err
	}

	logger.InfoLogger.Println("Database connected successfully")

	return db, nil
}
