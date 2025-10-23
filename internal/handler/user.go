package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chgrape/vaultpp/internal/middleware"
	"github.com/chgrape/vaultpp/internal/service"
)

type UserHandler struct {
	Service     *service.UserService
	JWTProvider *middleware.VaultProvider
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user service.UserValidator

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.Register(user, ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user service.LoginForm

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(user, ctx, h.JWTProvider.JwtKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error logging in: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *UserHandler) Details(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsCtxKey).(service.Claims)
	if !ok {
		http.Error(w, "No claims found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}
