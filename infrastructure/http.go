package infrastructure

import (
	"net/http"

	"github.com/NatBrian/Stockbit-Golang-Challenge/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ServeHTTP(app application.App) *chi.Mux {
	dep := application.SetupDependency(app)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	r.Post("/upload", dep.StockController.UploadTransaction)

	return r
}
