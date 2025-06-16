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

func main() {
	// Define command-line flags
	filePath := flag.String("file", "todos.json", "where to load/save the to-do list")
	addText := flag.String("add", "", "add a new to-do item (e.g. -add=\"Buy milk\")")
	updateID := flag.Int("update-id", 0, "the ID of the item you want to update")
	updateText := flag.String("update-text", "", "the new description for the item")
	updateStatus := flag.String("update-status", "", "the new status for the item (not started, started, completed)")
	deleteID := flag.Int("delete-id", 0, "the ID of the item you want to delete")
	serveAPI := flag.Bool("start-server", false, "Start HTTP API server")
	flag.Parse()

	traceID := uuid.NewString()
	ctx := context.WithValue(context.Background(), store.TraceIDKey, traceID)

	// Load existing items from disk
	items, err := store.LoadItems(ctx, *filePath)
	if err != nil {
		slog.Error("Failed to load items", "file", *filePath, "error", err)
		os.Exit(1)
	}

	// Based on flags, perform add/update/delete, mutating “items”.
	switch {
	case *addText != "":
		items = store.AddItem(ctx, items, *addText)
		last := items[len(items)-1]
		fmt.Printf("Added: [%d] %s\n", last.ID, last.Description)

	case *updateID != 0 && *updateText != "":
		var found bool
		items, found = store.UpdateItem(ctx, items, *updateID, *updateText)
		if !found {
			slog.Error("No item to update", "id", *updateID)
			os.Exit(1)
		}
		fmt.Printf("Updated: [%d] %s\n", *updateID, *updateText)

	case *updateID != 0 && *updateStatus != "":
		var found bool
		items, found = store.UpdateItemStatus(ctx, items, *updateID, *updateStatus)
		if !found {
			slog.Error("No item to update", "id", *updateID)
			os.Exit(1)
		}
		fmt.Printf("Updated status: [%d] %s\n", *updateID, *updateStatus)

	case *deleteID != 0:
		var found bool
		items, found = store.DeleteItem(ctx, items, *deleteID)
		if !found {
			slog.Error("No item to delete", "id", *deleteID)
			os.Exit(1)
		}
		fmt.Printf("Deleted item %d\n", *deleteID)

	case *serveAPI:
		api := &store.API{Items: &items}
		mux := http.NewServeMux()

		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

		mux.HandleFunc("/create", api.Create)
		mux.HandleFunc("/get", api.Get)
		mux.HandleFunc("/update", api.Update)
		mux.HandleFunc("/delete", api.Delete)

		mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
			tmpl := template.Must(template.New("list").Parse(`
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
    `))
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := tmpl.Execute(w, items); err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
			}
		})

		server := &http.Server{Addr: ":8080", Handler: mux}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		go func() {
			slog.Info("Starting HTTP server on :8080", "traceID", traceID)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("HTTP server error", "error", err, "traceID", traceID)
				os.Exit(1)
			}
		}()

		<-sigChan // Wait for Ctrl+C

		slog.Info("Interrupt received, shutting down server and saving items...")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctxTimeout); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
		if err := store.SaveItems(ctx, *filePath, items); err != nil {
			slog.Error("Failed to save items on interrupt", "error", err)
		} else {
			slog.Info("Items saved successfully on interrupt")
		}
		os.Exit(0)
	default:
		// No Need to do anything, just print the current list
	}

	// Print the (possibly updated) to-do list
	store.PrintItems(ctx, store.ListItems(items))

	slog.Info("Saved items to disk", "file", *filePath, "count", len(items))
}
