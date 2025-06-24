package store

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
)

// LoadItems reads a JSON file at path “filename” and returns the slice of Items.
// If the file does not exist, it returns an empty slice and no error. Any other error is returned directly.
func LoadItems(ctx context.Context, filename string) ([]Item, error) {
	traceID, _ := ctx.Value(TraceIDKey).(string)
	f, err := os.Open(filename)
	if err != nil {
		// If the file simply doesn’t exist, we start with an empty list.
		if os.IsNotExist(err) {
			slog.Info("File does not exist, starting with empty list",
				"file", filename,
				"traceID", traceID,
			)
			return []Item{}, nil
		}
		slog.Error("Failed to open file",
			"file", filename,
			"error", err,
			"traceID", traceID,
		)
		return nil, err
	}
	defer f.Close()

	var items []Item
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		slog.Error("Failed to decode items from file",
			"file", filename,
			"error", err,
			"traceID", traceID,
		)
		return nil, err
	}
	slog.Info("Loaded items from file",
		"file", filename,
		"count", len(items),
		"traceID", traceID,
	)
	return items, nil
}

// SaveItems writes the slice of Items as JSON to “filename” (overwriting or creating it).
// Returns any error encountered while creating or encoding.
func SaveItems(ctx context.Context, filename string, items []Item) error {
	traceID, _ := ctx.Value(TraceIDKey).(string)
	f, err := os.Create(filename)
	if err != nil {
		slog.Error("Failed to create file for saving items",
			"file", filename,
			"error", err,
			"traceID", traceID,
		)
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") // pretty-print with two-space indentation
	if err := enc.Encode(items); err != nil {
		slog.Error("Failed to encode items to file",
			"file", filename,
			"error", err,
			"traceID", traceID,
		)
		return err
	}
	slog.Info("Saved items to file",
		"file", filename,
		"count", len(items),
		"traceID", traceID,
	)
	return nil

}
