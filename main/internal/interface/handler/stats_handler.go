package handler

import (
	"net/http"

	itemuc "main-private/main/internal/usecase/item"
)

type StatsHandler struct {
	uc itemuc.UseCase
}

func NewStatsHandler(uc itemuc.UseCase) *StatsHandler {
	return &StatsHandler{uc: uc}
}

type statsResponse struct {
	TotalItems int `json:"total_items"`
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.ListItems(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, statsResponse{TotalItems: len(items)})
}
