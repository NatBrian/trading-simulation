package infrastructure

import (
	"github.com/NatBrian/Stockbit-Golang-Challenge/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ServeHTTP(app application.App) *chi.Mux {
	dep := application.SetupDependency(app)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/ping", dep.HealthContoller.Ping)
	r.Post("/upload", dep.StockController.UploadTransaction)
	r.Get("/summary", dep.StockController.GetSummary)

	return r
}
