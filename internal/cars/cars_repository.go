package cars

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/TazmanS/smartcar-backend/internal/cars/dto"
	"github.com/TazmanS/smartcar-backend/internal/cars/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		INSERT INTO app.cars (id, name)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		car.ID,
		car.Name,
	)

	return err
}

func (r *Repository) GetCarsList(ctx context.Context, req *dto.GetCarsListRequest) ([]models.Car, int, error) {
	offset := (req.Page - 1) * req.PerPage

	var totalItems int

	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM app.cars
	`).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, name, created_at, updated_at
		FROM app.cars
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, req.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	cars := make([]models.Car, 0, req.PerPage)

	for rows.Next() {
		var car models.Car

		err := rows.Scan(
			&car.ID,
			&car.Name,
			&car.CreatedAt,
			&car.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return cars, totalItems, nil
}

func (r *Repository) GetCarInfo(ctx context.Context, id uuid.UUID) (models.Car, error) {
	const query = `
		SELECT
			id,
			name,
			last_seen,
			created_at,
			updated_at
		FROM app.cars
		WHERE id = $1;
	`

	var car models.Car

	err := r.db.QueryRow(ctx, query, id).Scan(
		&car.ID,
		&car.Name,
		&car.LastSeen,
		&car.CreatedAt,
		&car.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Car{}, fmt.Errorf("car not found")
		}

		return models.Car{}, fmt.Errorf("get car info: %w", err)
	}

	return car, nil
}

func (r *Repository) CarHeartbeat(ctx context.Context, req dto.CarsHeartbeatRequest) error {
	const query = `
		UPDATE app.cars
		SET last_seen = NOW()
		WHERE id = $1;
	`

	tag, err := r.db.Exec(ctx, query, req.SessionID)
	if err != nil {
		return fmt.Errorf("update heartbeat: %w", err)
	}

	if tag.RowsAffected() == 0 {
		log.Printf("Heartbeat from unknown car id=%d", req.SessionID)
	}

	return nil
}

func (r *Repository) DeleteInactiveCars(ctx context.Context) error {
	const query = `
		DELETE FROM app.cars
		WHERE last_seen < NOW() - INTERVAL '60 seconds';
	`

	_, err := r.db.Exec(ctx, query)
	return err
}
