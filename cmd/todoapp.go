package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoapp/internal/store"

	"github.com/google/uuid"
)

const templateHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>ToDo List</title>
    <style>
        body { font-family: Arial, sans-serif; background: #f7f7f7; }
        .container { max-width: 600px; margin: 40px auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); padding: 32px; }
        h1 { color: #3498db; }
        ul { padding-left: 1.2em; }
        li { margin-bottom: 0.5em; }
    </style>
</head>
<body>
    <div class="container">
        <h1>ToDo List</h1>
        <ul>
            {{range .}}
                <li>
                    <strong>[{{.ID}}]</strong> {{.Description}} 
                    <em>(Status: {{.Status}}, Created: {{.CreatedAt.Format "2006-01-02 15:04"}})</em>
                </li>
            {{else}}
                <li>No items found.</li>
            {{end}}
        </ul>
    </div>
</body>
</html>
`

func main() {

	filePath := flag.String("file", "todos.json", "where to load/save the to-do list")
	addText := flag.String("add", "", "add a new to-do item")
	updateID := flag.Int("update-id", 0, "the ID of the item you want to update")
	updateText := flag.String("update-text", "", "the new description for the item")
	updateStatus := flag.String("update-status", "", "the new status for the item")
	deleteID := flag.Int("delete-id", 0, "the ID of the item you want to delete")
	serveAPI := flag.Bool("start-server", false, "Start HTTP API server")
	flag.Parse()

	traceID := uuid.NewString()
	ctx := context.WithValue(context.Background(), store.TraceIDKey, traceID)

	items, err := store.LoadItems(ctx, *filePath)
	if err != nil {
		slog.Error("Failed to load items", "file", *filePath, "error", err, "traceID", traceID)
		os.Exit(1)
	}
	actor := store.NewToDoActor(items)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if !*serveAPI {
		handleCLI(actor, ctx, *filePath, *addText, *updateID, *updateText, *updateStatus, *deleteID, traceID)
		return
	}

	startAPIServer(actor, ctx, *filePath, traceID, sigChan)
}

func handleCLI(actor *store.ToDoActor, ctx context.Context, filePath, addText string, updateID int, updateText, updateStatus string, deleteID int, traceID string) {
	switch {
	case addText != "":
		item := actor.AddItem(addText)
		fmt.Printf("Added: [%d] %s\n", item.ID, item.Description)
	case updateID != 0 && updateText != "":
		if !actor.UpdateItem(updateID, updateText, "") {
			slog.Error("No item to update", "id", updateID, "traceID", traceID)
			os.Exit(1)
		}
		fmt.Printf("Updated: [%d] %s\n", updateID, updateText)
	case updateID != 0 && updateStatus != "":
		if !actor.UpdateItem(updateID, "", updateStatus) {
			slog.Error("No item to update", "id", updateID, "traceID", traceID)
			os.Exit(1)
		}
		fmt.Printf("Updated status: [%d] %s\n", updateID, updateStatus)
	case deleteID != 0:
		if !actor.DeleteItem(deleteID) {
			slog.Error("No item to delete", "id", deleteID, "traceID", traceID)
			os.Exit(1)
		}
		fmt.Printf("Deleted item %d\n", deleteID)
	default:
		items := actor.GetItems()
		store.PrintItems(ctx, items)
	}

	// Save items on exit
	items := actor.GetItems()
	if err := store.SaveItems(ctx, filePath, items); err != nil {
		slog.Error("Failed to save items", "file", filePath, "error", err, "traceID", traceID)
		os.Exit(1)
	}
	slog.Info("Saved items to disk", "file", filePath, "count", len(items), "traceID", traceID)
}

func startAPIServer(actor *store.ToDoActor, ctx context.Context, filePath, traceID string, sigChan chan os.Signal) {
	api := &store.API{Actor: actor}
	mux := http.NewServeMux()
	mux.HandleFunc("/create", api.Create)
	mux.HandleFunc("/get", api.Get)
	mux.HandleFunc("/update", api.Update)
	mux.HandleFunc("/delete", api.Delete)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		items := actor.GetItems()
		tmpl := template.Must(template.New("list").Parse(templateHTML))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, items); err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
	})

	handler := store.TraceIDMiddleware(mux)
	server := &http.Server{Addr: ":8080", Handler: handler}

	go func() {
		slog.Info("Starting HTTP server on :8080", "traceID", traceID)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err, "traceID", traceID)
			os.Exit(1)
		}
	}()

	<-sigChan
	slog.Info("Interrupt received, shutting down server and saving items...", "traceID", traceID)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxTimeout); err != nil {
		slog.Error("Server shutdown error", "error", err, "traceID", traceID)
	}
	items := actor.GetItems()
	if err := store.SaveItems(ctx, filePath, items); err != nil {
		slog.Error("Failed to save items on interrupt", "error", err, "traceID", traceID)
	} else {
		slog.Info("Items saved successfully on interrupt", "traceID", traceID)
	}
}
