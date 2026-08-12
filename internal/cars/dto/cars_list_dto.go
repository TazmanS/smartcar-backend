package dto

import "github.com/TazmanS/smartcar-backend/internal/cars/models"

type GetCarsListRequest struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Search  string `json:"search"`
	SortBy  string `json:"sort_by"`
	Order   string `json:"order"` // asc | desc
}

type PageInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type GetCarsListResponse struct {
	Data []models.Car `json:"data"`
	Page PageInfo     `json:"page"`
}
