package main

import (
	"flag"
	"fmt"
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
	flag.Parse()

	// Load existing items from disk
	items, err := store.LoadItems(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", *filePath, err)
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
			fmt.Fprintf(os.Stderr, "No item with ID %d to update\n", *updateID)
			os.Exit(1)
		}
		fmt.Printf("Updated: [%d] %s\n", *updateID, *updateText)

	case *deleteID != 0:
		var found bool
		items, found = store.DeleteItem(items, *deleteID)
		if !found {
			fmt.Fprintf(os.Stderr, "No item with ID %d to delete\n", *deleteID)
			os.Exit(1)
		}
		fmt.Printf("Deleted item %d\n", *deleteID)

	default:
		// No Need to do anything, just print the current list
	}

	// Print the (possibly updated) to-do list
	store.PrintItems(store.ListItems(items))

	// Save the (possibly modified) list back to disk
	if err := store.SaveItems(*filePath, items); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving to %s: %v\n", *filePath, err)
		os.Exit(1)
	}
}
