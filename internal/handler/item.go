package handler

import (
	"encoding/json"
	"net/http"

	"github.com/chgrape/vaultpp/internal/repository"
	"github.com/chgrape/vaultpp/internal/service"
)

type ItemHandler struct {
	Service *service.ItemService
}

func (h *ItemHandler) HandleGetItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.Service.ListItems(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) HandleCreateItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var item repository.Item

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusInternalServerError)
		return
	}

	id, err := h.Service.AddItem(item, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}
