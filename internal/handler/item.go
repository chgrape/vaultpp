package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (h *ItemHandler) HandleGetItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	item, err := h.Service.ListItem(id, ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get item: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)

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

func (h *ItemHandler) HandleUpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var item repository.Item
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Couldn't read body", http.StatusInternalServerError)
		return
	}

	id, err = h.Service.EditItem(item, id, ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update item: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(id)
}

func (h *ItemHandler) HandleDeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	id, err = h.Service.RemoveItem(id, ctx)
	if err != nil {
		http.Error(w, "Couldn't remove item", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(id)
}
