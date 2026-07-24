package routes

import (
	"github.com/TazmanS/smartcar-backend/internal/app"
	"github.com/TazmanS/smartcar-backend/internal/cars"
	"github.com/TazmanS/smartcar-backend/internal/health"
	"github.com/TazmanS/smartcar-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/TazmanS/smartcar-backend/docs"
)

func NewRouter(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Cors())
	r.Use(middleware.Logger)

	r.Get("/swagger/*", httpSwagger.Handler())

	r.Route("/api", func(api chi.Router) {
		health.RegisterHealthRoutes(api, app)

		cars.RegisterCarRoutes(api, app)
	})

	return r
}
