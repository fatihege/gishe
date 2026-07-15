package authhttp

import (
	"errors"
	"log"
	"net/http"

	"github.com/fatihege/gishe/internal/auth"
	"github.com/fatihege/gishe/internal/httpx"
)

type Handler struct {
	service *auth.Service
}

func NewHandler(service *auth.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := httpx.ReadJSON(w, r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Register(ctx, auth.RegisterInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrEmailAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, "user with this email already exists")

		case errors.Is(err, auth.ErrWeakPassword):
			httpx.WriteError(w, http.StatusUnprocessableEntity, "password does not meet requirements")

		default:
			log.Println("register:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"user": auth.NewUserResponse(user),
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := httpx.ReadJSON(w, r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Login(ctx, auth.LoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		} else {
			log.Println("login:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"user": auth.NewUserResponse(user),
	})
}
