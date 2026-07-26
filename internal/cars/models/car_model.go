package models

import "time"

type Car struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
