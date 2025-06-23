package store

import (
	"context"
	"testing"
)

func TestAddItem(t *testing.T) {
	ctx := context.Background()
	items := []Item{}
	items = AddItem(ctx, items, "Test item")
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Description != "Test item" {
		t.Errorf("expected description 'Test item', got '%s'", items[0].Description)
	}
}

func TestUpdateItem(t *testing.T) {
	ctx := context.Background()
	items := []Item{{ID: 1, Description: "Old"}}
	updated, found := UpdateItem(ctx, items, 1, "New")
	if !found {
		t.Fatal("expected item to be found")
	}
	if updated[0].Description != "New" {
		t.Errorf("expected description 'New', got '%s'", updated[0].Description)
	}
}

func TestUpdateItemStatus(t *testing.T) {
	ctx := context.Background()
	items := []Item{{ID: 1, Status: StatusNotStarted}}
	updated, found := UpdateItemStatus(ctx, items, 1, StatusCompleted)
	if !found {
		t.Fatal("expected item to be found")
	}
	if updated[0].Status != StatusCompleted {
		t.Errorf("expected status '%s', got '%s'", StatusCompleted, updated[0].Status)
	}
}

func TestDeleteItem(t *testing.T) {
	ctx := context.Background()
	items := []Item{{ID: 1, Description: "ToDelete"}}
	updated, found := DeleteItem(ctx, items, 1)
	if !found {
		t.Fatal("expected item to be found")
	}
	if len(updated) != 0 {
		t.Errorf("expected item to be deleted")
	}
}
