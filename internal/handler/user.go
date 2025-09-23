package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chgrape/vaultpp/internal/service"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user service.UserValidator

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.RegisterUser(user, ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}
