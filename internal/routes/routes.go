package routes

import (
	"github.com/TazmanS/smartcar-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/TazmanS/smartcar-backend/docs"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Cors())
	r.Use(middleware.Logger)

	r.Get("/swagger/*", httpSwagger.Handler())

	return r
}
