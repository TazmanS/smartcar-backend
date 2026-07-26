package cars

import (
	"context"

	"github.com/TazmanS/smartcar-backend/internal/cars/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, car *models.Car) error {
	query := `
		INSERT INTO cars (name)
		VALUES ($1)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		car.Name,
	).Scan(
		&car.ID,
		&car.CreatedAt,
		&car.UpdatedAt,
	)
}
