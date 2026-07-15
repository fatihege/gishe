package authhttp

import (
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
		log.Println("readjson:", err)
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Register(ctx, auth.RegisterInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "an error occurred while registering the user")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"user": auth.NewUserResponse(user),
	})
}
