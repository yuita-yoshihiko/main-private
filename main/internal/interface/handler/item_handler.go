package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	domain "main-private/main/internal/domain/item"
	itemuc "main-private/main/internal/usecase/item"
)

type ItemHandler struct {
	uc itemuc.UseCase
}

func NewItemHandler(uc itemuc.UseCase) *ItemHandler {
	return &ItemHandler{uc: uc}
}

type createItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type itemResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func toItemResponse(item domain.Item) itemResponse {
	return itemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	item, err := h.uc.CreateItem(r.Context(), req.Name, req.Description)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toItemResponse(item))
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}

	item, err := h.uc.GetItem(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toItemResponse(item))
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.ListItems(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]itemResponse, len(items))
	for i, item := range items {
		resp[i] = toItemResponse(item)
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}

	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	item, err := h.uc.UpdateItem(r.Context(), id, req.Name, req.Description)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toItemResponse(item))
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}

	if err := h.uc.DeleteItem(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemHandler) Search(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "name query parameter is required"})
		return
	}

	items, err := h.uc.SearchItems(r.Context(), name)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]itemResponse, len(items))
	for i, item := range items {
		resp[i] = toItemResponse(item)
	}

	writeJSON(w, http.StatusOK, resp)
}
