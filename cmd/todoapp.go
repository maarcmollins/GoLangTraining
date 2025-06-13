package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
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

	default:
		// No Need to do anything, just print the current list
	}

	// Print the (possibly updated) to-do list
	store.PrintItems(ctx, store.ListItems(items))

	// Save the (possibly modified) list back to disk
	if err := store.SaveItems(ctx, *filePath, items); err != nil {
		slog.Error("Failed to save items", "file", *filePath, "error", err)
		os.Exit(1)
	}
	slog.Info("Saved items to disk", "file", *filePath, "count", len(items))
}
