package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	itemuc "main-private/main/internal/usecase/item"
	statsuc "main-private/main/internal/usecase/stats"
)

type StatsHandler struct {
	itemUC  itemuc.UseCase
	statsUC statsuc.UseCase
}

func NewStatsHandler(itemUC itemuc.UseCase, statsUC statsuc.UseCase) *StatsHandler {
	return &StatsHandler{itemUC: itemUC, statsUC: statsUC}
}

type statsResponse struct {
	TotalItems int `json:"total_items"`
}

type createStatsResponse struct {
	ID         uuid.UUID `json:"id"`
	TotalItems int       `json:"total_items"`
	CreatedAt  time.Time `json:"created_at"`
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	items, err := h.itemUC.ListItems(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, statsResponse{TotalItems: len(items)})
}

func (h *StatsHandler) CreateStats(w http.ResponseWriter, r *http.Request) {
	items, err := h.itemUC.ListItems(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	stats, err := h.statsUC.CreateStats(r.Context(), len(items))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createStatsResponse{
		ID:         stats.ID,
		TotalItems: stats.TotalItems,
		CreatedAt:  stats.CreatedAt,
	})
}
