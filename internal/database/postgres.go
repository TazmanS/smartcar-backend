package database

import (
	"context"
	"fmt"
	"log"

	"github.com/TazmanS/smartcar-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	log.Println("Connected to PostgreSQL")

	return pool, nil
}
