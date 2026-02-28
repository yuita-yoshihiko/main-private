package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"main-private/main/internal/interface/handler"
	itemuc "main-private/main/internal/usecase/item"
)

func New(logger *slog.Logger, itemUC itemuc.UseCase) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	healthHandler := handler.NewHealthHandler()
	itemHandler := handler.NewItemHandler(itemUC)
	statsHandler := handler.NewStatsHandler(itemUC)

	r.Get("/health", healthHandler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/stats", statsHandler.GetStats)
		r.Route("/items", func(r chi.Router) {
			r.Post("/", itemHandler.Create)
			r.Get("/", itemHandler.List)
			r.Get("/search", itemHandler.Search)
			r.Get("/latest", itemHandler.Latest)
			r.Get("/{id}", itemHandler.Get)
			r.Put("/{id}", itemHandler.Update)
			r.Patch("/{id}", itemHandler.Patch)
			r.Delete("/{id}", itemHandler.Delete)
		})
	})

	return r
}
