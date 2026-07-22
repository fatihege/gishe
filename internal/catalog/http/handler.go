package cataloghttp

import (
	"errors"
	"log"
	"net/http"

	"github.com/fatihege/gishe/internal/catalog"
	"github.com/fatihege/gishe/internal/httpx"
)

type Handler struct {
	service *catalog.Service
}

func NewHandler(service *catalog.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		City    string `json:"city"`
	}

	if err := httpx.ReadJSON(w, r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	venue, err := h.service.CreateVenue(ctx, catalog.CreateVenueInput{
		Name:    request.Name,
		Address: request.Address,
		City:    request.City,
	})
	if err != nil {
		log.Println("create venue:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"venue": venue,
	})
}

func (h *Handler) GetVenues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	venues, err := h.service.GetVenues(ctx)
	if errors.Is(err, catalog.ErrNoVenuesFound) {
		httpx.WriteError(w, http.StatusOK, "no venues were added yet")
		return
	}
	if err != nil {
		log.Println("create venue:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"venues": venues,
	})
}
