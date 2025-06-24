package store

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type API struct {
	Actor *ToDoActor
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID, _ := ctx.Value(TraceIDKey).(string)
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body for create", "error", err, "traceID", traceID)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	item := api.Actor.AddItem(req.Description)
	slog.Info("Created new item", "description", req.Description, "traceID", traceID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (api *API) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID, _ := ctx.Value(TraceIDKey).(string)
	slog.Info("Get all items", "traceID", traceID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.Actor.GetItems())
}

func (api *API) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID, _ := ctx.Value(TraceIDKey).(string)
	var req struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body for update", "error", err, "traceID", traceID)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	found := api.Actor.UpdateItem(req.ID, req.Description, req.Status)
	if !found {
		slog.Error("Item not found for update", "id", req.ID, "traceID", traceID)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	slog.Info("Updated item", "id", req.ID, "traceID", traceID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.Actor.GetItems())
}

func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID, _ := ctx.Value(TraceIDKey).(string)
	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body for delete", "error", err, "traceID", traceID)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	found := api.Actor.DeleteItem(req.ID)
	if !found {
		slog.Error("Item not found for delete", "id", req.ID, "traceID", traceID)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	slog.Info("Deleted item", "id", req.ID, "traceID", traceID)
	w.WriteHeader(http.StatusNoContent)
}
