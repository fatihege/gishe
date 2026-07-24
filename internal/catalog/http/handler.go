package cataloghttp

import (
	"errors"
	"log"
	"net/http"

	"github.com/fatihege/gishe/internal/catalog"
	"github.com/fatihege/gishe/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *catalog.Service
}

func NewHandler(service *catalog.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input catalog.CreateVenueInput

	if err := httpx.ReadJSON(r, &input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	venue, err := h.service.CreateVenue(ctx, input)
	if err != nil {
		if errors.Is(err, catalog.ErrVenueFieldsRequired) {
			httpx.WriteError(w, http.StatusBadRequest, "name, address, and city fields are required")

		} else {
			log.Println("create venue:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"venue": venue,
	})
}

func (h *Handler) GetVenues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	venues, err := h.service.GetVenues(ctx)
	if err != nil {
		log.Println("get venues:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"venues": venues,
	})
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input catalog.CreateSessionInput

	if err := httpx.ReadJSON(r, &input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	session, err := h.service.CreateSession(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, catalog.ErrInvalidVenue):
			httpx.WriteError(w, http.StatusBadRequest, "a valid venue should be provided")

		case errors.Is(err, catalog.ErrSessionTitleRequired):
			httpx.WriteError(w, http.StatusBadRequest, "title is required for sessions")

		case errors.Is(err, catalog.ErrSessionTimesRequired):
			httpx.WriteError(w, http.StatusBadRequest, "session sales open time and start time are required")

		case errors.Is(err, catalog.ErrInvalidSessionSchedule):
			httpx.WriteError(w, http.StatusBadRequest, "there is a conflict between the session times")

		default:
			log.Println("create session:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"session": session,
	})
}

func (h *Handler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawID := chi.URLParam(r, "id")

	id, err := uuid.Parse(rawID)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	session, err := h.service.GetSessionById(ctx, id)
	if err != nil {
		if errors.Is(err, catalog.ErrSessionNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "session does not exist")
		} else {
			log.Println("get session by id:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"session": session,
	})
}

func (h *Handler) GetSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessions, err := h.service.GetSessions(ctx)
	if err != nil {
		log.Println("get sessions:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"sessions": sessions,
	})
}

func (h *Handler) GetSessionsByVenueID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawVenueID := chi.URLParam(r, "id")

	venueID, err := uuid.Parse(rawVenueID)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid venue id")
		return
	}

	sessions, err := h.service.GetSessionsByVenueID(ctx, venueID)
	if err != nil {
		if errors.Is(err, catalog.ErrInvalidVenue) {
			httpx.WriteError(w, http.StatusBadRequest, "venue does not exist")
		} else {
			log.Println("get sessions by venue id:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"sessions": sessions,
	})
}
