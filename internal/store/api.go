package store

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type API struct {
	Items *[]Item
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body for create", "error", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	*api.Items = AddItem(r.Context(), *api.Items, req.Description)
	slog.Info("Created new item", "description", req.Description)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode((*api.Items)[len(*api.Items)-1])
}

func (api *API) Get(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get all items")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*api.Items)
}

func (api *API) Update(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body for update", "error", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var found bool
	if req.Description != "" {
		*api.Items, found = UpdateItem(r.Context(), *api.Items, req.ID, req.Description)
		slog.Info("Update item description", "id", req.ID, "description", req.Description)
	} else if req.Status != "" {
		*api.Items, found = UpdateItemStatus(r.Context(), *api.Items, req.ID, req.Status)
		slog.Info("Update item status", "id", req.ID, "status", req.Status)
	}
	if !found {
		slog.Error("Item not found for update", "id", req.ID)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*api.Items)
}

func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body for delete", "error", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var found bool
	*api.Items, found = DeleteItem(r.Context(), *api.Items, req.ID)
	if !found {
		slog.Error("Item not found for delete", "id", req.ID)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	slog.Info("Deleted item", "id", req.ID)
	w.WriteHeader(http.StatusNoContent)
}
