package store

import (
	"encoding/json"
	"os"
)

<<<<<<< Updated upstream
// LoadItems reads a JSON file at path “filename” and returns the slice of Items.
// If the file does not exist, it returns an empty slice and no error. Any other error is returned directly.
func LoadItems(filename string) ([]Item, error) {
=======
func LoadItems(ctx context.Context, filename string) ([]Item, error) {
	traceID, _ := ctx.Value(TraceIDKey).(string)
>>>>>>> Stashed changes
	f, err := os.Open(filename)
	if err != nil {
		// If the file simply doesn’t exist, we start with an empty list.
		if os.IsNotExist(err) {
			return []Item{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var items []Item
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

<<<<<<< Updated upstream
// SaveItems writes the slice of Items as JSON to “filename” (overwriting or creating it).
// Returns any error encountered while creating or encoding.
func SaveItems(filename string, items []Item) error {
=======
func SaveItems(ctx context.Context, filename string, items []Item) error {
	traceID, _ := ctx.Value(TraceIDKey).(string)
>>>>>>> Stashed changes
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") // pretty-print with two-space indentation
	return enc.Encode(items)
}
