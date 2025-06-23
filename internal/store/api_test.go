package store

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testCtx() context.Context {
	return context.WithValue(context.Background(), TraceIDKey, "test-trace-id")
}

func TestAPI_Create(t *testing.T) {
	actor := NewToDoActor([]Item{})
	api := &API{Actor: actor}

	body := bytes.NewBufferString(`{"description":"Test task"}`)
	req := httptest.NewRequest(http.MethodPost, "/create", body).WithContext(testCtx())
	w := httptest.NewRecorder()

	api.Create(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var item Item
	if err := json.NewDecoder(w.Body).Decode(&item); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if item.Description != "Test task" {
		t.Errorf("expected description 'Test task', got '%s'", item.Description)
	}
}

func TestAPI_Get(t *testing.T) {
	actor := NewToDoActor([]Item{{ID: 1, Description: "Task"}})
	api := &API{Actor: actor}

	req := httptest.NewRequest(http.MethodGet, "/get", nil).WithContext(testCtx())
	w := httptest.NewRecorder()

	api.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var got []Item
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(got) != 1 || got[0].Description != "Task" {
		t.Errorf("unexpected items: %+v", got)
	}
}

func TestAPI_Update(t *testing.T) {
	actor := NewToDoActor([]Item{{ID: 1, Description: "Old"}})
	api := &API{Actor: actor}

	body := bytes.NewBufferString(`{"id":1,"description":"New"}`)
	req := httptest.NewRequest(http.MethodPost, "/update", body).WithContext(testCtx())
	w := httptest.NewRecorder()

	api.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	items := actor.GetItems()
	if items[0].Description != "New" {
		t.Errorf("expected updated description, got %s", items[0].Description)
	}
}

func TestAPI_Delete(t *testing.T) {
	actor := NewToDoActor([]Item{{ID: 1, Description: "ToDelete"}})
	api := &API{Actor: actor}

	body := bytes.NewBufferString(`{"id":1}`)
	req := httptest.NewRequest(http.MethodPost, "/delete", body).WithContext(testCtx())
	w := httptest.NewRecorder()

	api.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
	items := actor.GetItems()
	if len(items) != 0 {
		t.Errorf("expected item to be deleted")
	}
}
