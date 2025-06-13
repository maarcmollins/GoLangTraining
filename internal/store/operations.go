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

<<<<<<< Updated upstream
// AddItem appends a new Item to the slice and returns the new slice.
// It does NOT save to disk—that responsibility lies with the caller.
func AddItem(items []Item, description string) []Item {
=======
func AddItem(ctx context.Context, items []Item, description string) []Item {
	traceID, _ := ctx.Value(TraceIDKey).(string)
>>>>>>> Stashed changes
	id := nextID(items)
	newItem := Item{
		ID:          id,
		Description: description,
		CreatedAt:   time.Now(),
	}
<<<<<<< Updated upstream
	return append(items, newItem)
}

// UpdateItem searches “items” for an Item whose ID == targetID. If found, it replaces its Description.
// Returns the modified slice and a boolean “found==true”. If not found, returns the original slice and false.
func UpdateItem(items []Item, targetID int, newDescription string) ([]Item, bool) {
=======
	slog.Info("Added new item",
		"id", newItem.ID,
		"description", newItem.Description,
		"traceID", traceID,
	)
	return append(items, newItem)
}

func UpdateItem(ctx context.Context, items []Item, targetID int, newDescription string) ([]Item, bool) {
	traceID, _ := ctx.Value(TraceIDKey).(string)
>>>>>>> Stashed changes
	for i, it := range items {
		if it.ID == targetID {
			items[i].Description = newDescription
<<<<<<< Updated upstream
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
=======
			slog.Info("Updated item description",
				"id", targetID,
				"old_description", oldDesc,
				"new_description", newDescription,
				"traceID", traceID,
			)
			return items, true
		}
	}
	slog.Error("Item to update not found",
		"id", targetID,
		"traceID", traceID,
	)
	return items, false
}

func UpdateItemStatus(ctx context.Context, items []Item, targetID int, newStatus string) ([]Item, bool) {
	traceID, _ := ctx.Value(TraceIDKey).(string)
	switch newStatus {
	case StatusNotStarted, StatusStarted, StatusCompleted:
	default:
		slog.Error("Invalid status for update",
			"id", targetID,
			"status", newStatus,
			"traceID", traceID,
		)
		return items, false
	}
	for i, it := range items {
		if it.ID == targetID {
			oldStatus := items[i].Status
			items[i].Status = newStatus
			slog.Info("Updated item status",
				"id", targetID,
				"old_status", oldStatus,
				"new_status", newStatus,
				"traceID", traceID,
			)
			return items, true
		}
	}
	slog.Error("Item to update status not found",
		"id", targetID,
		"traceID", traceID,
	)
	return items, false
}

func DeleteItem(ctx context.Context, items []Item, targetID int) ([]Item, bool) {
	traceID, _ := ctx.Value(TraceIDKey).(string)
	for i, it := range items {
		if it.ID == targetID {
			slog.Info("Deleted item",
				"id", targetID,
				"traceID", traceID,
			)
			return append(items[:i], items[i+1:]...), true
		}
	}
	slog.Error("Item to delete not found",
		"id", targetID,
		"traceID", traceID,
	)
>>>>>>> Stashed changes
	return items, false
}

func ListItems(items []Item) []Item {
	return items
}

<<<<<<< Updated upstream
// PrintItems writes the current list of items to stdout in a human-readable format.
func PrintItems(items []Item) {
	if len(items) == 0 {
=======
func PrintItems(ctx context.Context, items []Item) {
	traceID, _ := ctx.Value(TraceIDKey).(string)
	if len(items) == 0 {
		slog.Info("No to-do items", "traceID", traceID)
>>>>>>> Stashed changes
		fmt.Println("No to-do items.")
		return
	}
	fmt.Println("Current to-do list:")
	for _, it := range items {
<<<<<<< Updated upstream
		fmt.Printf("  [%d] %s (created: %s)\n",
=======
		slog.Info("To-do item",
			"id", it.ID,
			"description", it.Description,
			"status", it.Status,
			"created_at", it.CreatedAt.Format(time.RFC3339),
			"traceID", traceID,
		)
		fmt.Printf("  [%d] %s (status: %s, created: %s)\n",
>>>>>>> Stashed changes
			it.ID,
			it.Description,
			it.CreatedAt.Format(time.RFC3339))
	}
}
