package models

import "time"

type Car struct {
	ID        string
	Name      string
	LastSeen  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
