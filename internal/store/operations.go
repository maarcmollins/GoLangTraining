package store

import (
	"fmt"
	"time"
)

// nextID returns one larger than the maximum ID in “items”. If the slice is empty, returns 1.
func nextID(items []Item) int {
	max := 0
	for _, item := range items {
		if item.ID > max {
			max = item.ID
		}
	}
	return max + 1
}

// AddItem appends a new Item to the slice and returns the new slice.
// It does NOT save to disk—that responsibility lies with the caller.
func AddItem(items []Item, description string) []Item {
	id := nextID(items)
	newItem := Item{
		ID:          id,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return append(items, newItem)
}

// UpdateItem searches “items” for an Item whose ID == targetID. If found, it replaces its Description.
// Returns the modified slice and a boolean “found==true”. If not found, returns the original slice and false.
func UpdateItem(items []Item, targetID int, newDescription string) ([]Item, bool) {
	for i, it := range items {
		if it.ID == targetID {
			items[i].Description = newDescription
			return items, true
		}
	}
	return items, false
}

// DeleteItem removes the Item whose ID == targetID from the slice (if present).
// Returns the new slice and a boolean “found==true”. If not found, returns the original slice and false.
func DeleteItem(items []Item, targetID int) ([]Item, bool) {
	for i, it := range items {
		if it.ID == targetID {
			// Remove index i from slice
			return append(items[:i], items[i+1:]...), true
		}
	}
	return items, false
}

// ListItems simply returns the slice as-is. (We return it so main.go can range over it.)
func ListItems(items []Item) []Item {
	return items
}

// PrintItems writes the current list of items to stdout in a human-readable format.
func PrintItems(items []Item) {
	if len(items) == 0 {
		fmt.Println("No to-do items.")
		return
	}
	fmt.Println("Current to-do list:")
	for _, it := range items {
		fmt.Printf("  [%d] %s (created: %s)\n",
			it.ID,
			it.Description,
			it.CreatedAt.Format(time.RFC3339))
	}
}
