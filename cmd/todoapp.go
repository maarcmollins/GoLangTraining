package main

import (
	"flag"
	"fmt"
<<<<<<< Updated upstream
=======
	"log/slog"
	"net/http"
>>>>>>> Stashed changes
	"os"
	"todoapp/internal/store"
)

func main() {
	// Define command-line flags
	filePath := flag.String("file", "todos.json", "where to load/save the to-do list")
	addText := flag.String("add", "", "add a new to-do item (e.g. -add=\"Buy milk\")")
	updateID := flag.Int("update-id", 0, "the ID of the item you want to update")
	updateText := flag.String("update-text", "", "the new description for the item")
	deleteID := flag.Int("delete-id", 0, "the ID of the item you want to delete")
	serveAPI := flag.Bool("start-server", false, "Start HTTP API server")
	flag.Parse()

	// Load existing items from disk
	items, err := store.LoadItems(*filePath)
	if err != nil {
<<<<<<< Updated upstream
		fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", *filePath, err)
=======
		slog.Error("Failed to load items", "file", *filePath, "error", err, "traceID", traceID)
>>>>>>> Stashed changes
		os.Exit(1)
	}

	// Based on flags, perform add/update/delete, mutating “items”.
	switch {
	case *addText != "":
		items = store.AddItem(items, *addText)
		last := items[len(items)-1]
		fmt.Printf("Added: [%d] %s\n", last.ID, last.Description)

	case *updateID != 0 && *updateText != "":
		var found bool
		items, found = store.UpdateItem(items, *updateID, *updateText)
		if !found {
<<<<<<< Updated upstream
			fmt.Fprintf(os.Stderr, "No item with ID %d to update\n", *updateID)
=======
			slog.Error("No item to update", "id", *updateID, "traceID", traceID)
>>>>>>> Stashed changes
			os.Exit(1)
		}
		fmt.Printf("Updated: [%d] %s\n", *updateID, *updateText)

<<<<<<< Updated upstream
=======
	case *updateID != 0 && *updateStatus != "":
		var found bool
		items, found = store.UpdateItemStatus(ctx, items, *updateID, *updateStatus)
		if !found {
			slog.Error("No item to update", "id", *updateID, "traceID", traceID)
			os.Exit(1)
		}
		fmt.Printf("Updated status: [%d] %s\n", *updateID, *updateStatus)

>>>>>>> Stashed changes
	case *deleteID != 0:
		var found bool
		items, found = store.DeleteItem(items, *deleteID)
		if !found {
<<<<<<< Updated upstream
			fmt.Fprintf(os.Stderr, "No item with ID %d to delete\n", *deleteID)
=======
			slog.Error("No item to delete", "id", *deleteID, "traceID", traceID)
>>>>>>> Stashed changes
			os.Exit(1)
		}
		fmt.Printf("Deleted item %d\n", *deleteID)

	case *serveAPI:
		api := &store.API{Items: &items}
		mux := http.NewServeMux()
		mux.HandleFunc("/create", api.Create)
		mux.HandleFunc("/get", api.Get)
		mux.HandleFunc("/update", api.Update)
		mux.HandleFunc("/delete", api.Delete)

		handler := store.TraceIDMiddleware(mux)

		slog.Info("Starting HTTP server on :8080", "traceID", traceID)
		if err := http.ListenAndServe(":8080", handler); err != nil {
			slog.Error("HTTP server error", "error", err, "traceID", traceID)
			os.Exit(1)
		}

	default:
		// No Need to do anything, just print the current list
	}

<<<<<<< Updated upstream
	// Print the (possibly updated) to-do list
	store.PrintItems(store.ListItems(items))

	// Save the (possibly modified) list back to disk
	if err := store.SaveItems(*filePath, items); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving to %s: %v\n", *filePath, err)
		os.Exit(1)
	}
=======
	store.PrintItems(ctx, store.ListItems(items))

	if err := store.SaveItems(ctx, *filePath, items); err != nil {
		slog.Error("Failed to save items", "file", *filePath, "error", err, "traceID", traceID)
		os.Exit(1)
	}
	slog.Info("Saved items to disk", "file", *filePath, "count", len(items), "traceID", traceID)
>>>>>>> Stashed changes
}
