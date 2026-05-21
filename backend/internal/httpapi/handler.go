package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"who-can-search-ygo/backend/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.healthz)
	mux.HandleFunc("GET /docs", h.docs)
	mux.HandleFunc("GET /openapi.json", h.openAPI)
	mux.HandleFunc("GET /cards", h.searchCards)
	mux.HandleFunc("GET /cards/{id}", h.getCard)
	mux.HandleFunc("GET /cards/{id}/searchers", h.getCardSearchers)
	return mux
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) searchCards(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		writeError(w, http.StatusBadRequest, "invalid_query", "query is required")
		return
	}

	limit := 20
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 || parsed > 50 {
			writeError(w, http.StatusBadRequest, "invalid_limit", "limit must be an integer between 1 and 50")
			return
		}
		limit = parsed
	}

	cards, err := h.service.SearchCards(r.Context(), query, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not search cards")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": cards})
}

func (h *Handler) getCard(w http.ResponseWriter, r *http.Request) {
	card, err := h.service.GetCard(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "card_not_found", "card not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "could not load card")
		return
	}
	writeJSON(w, http.StatusOK, newCardResponse(card))
}

func (h *Handler) getCardSearchers(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetAddDeckToHandSearchers(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "card_not_found", "card not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "could not load searchers")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
