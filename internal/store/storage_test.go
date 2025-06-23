package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadItems_FileDoesNotExist(t *testing.T) {
	ctx := context.WithValue(context.Background(), TraceIDKey, "test-trace-id")
	tmpDir := t.TempDir()
	nonExistent := filepath.Join(tmpDir, "doesnotexist.json")

	items, err := LoadItems(ctx, nonExistent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected empty items, got %v", items)
	}
}

func TestSaveAndLoadItems(t *testing.T) {
	ctx := context.WithValue(context.Background(), TraceIDKey, "test-trace-id")
	tmpFile := filepath.Join(t.TempDir(), "todos.json")

	original := []Item{
		{ID: 1, Description: "Test 1", Status: StatusNotStarted},
		{ID: 2, Description: "Test 2", Status: StatusCompleted},
	}

	if err := SaveItems(ctx, tmpFile, original); err != nil {
		t.Fatalf("SaveItems failed: %v", err)
	}

	loaded, err := LoadItems(ctx, tmpFile)
	if err != nil {
		t.Fatalf("LoadItems failed: %v", err)
	}
	if len(loaded) != len(original) {
		t.Fatalf("expected %d items, got %d", len(original), len(loaded))
	}
	for i := range original {
		if loaded[i].ID != original[i].ID ||
			loaded[i].Description != original[i].Description ||
			loaded[i].Status != original[i].Status {
			t.Errorf("item mismatch at %d: got %+v, want %+v", i, loaded[i], original[i])
		}
	}
}

func TestLoadItems_InvalidJSON(t *testing.T) {
	ctx := context.WithValue(context.Background(), TraceIDKey, "test-trace-id")
	tmpFile := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(tmpFile, []byte("{not valid json"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	_, err := LoadItems(ctx, tmpFile)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
