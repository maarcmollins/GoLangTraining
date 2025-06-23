package store

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Create(t *testing.T) {
	items := []Item{}
	api := &API{Items: &items}

	body := bytes.NewBufferString(`{"description":"Test task"}`)
	req := httptest.NewRequest(http.MethodPost, "/create", body)
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
	items := []Item{{ID: 1, Description: "Task"}}
	api := &API{Items: &items}

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
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
	items := []Item{{ID: 1, Description: "Old"}}
	api := &API{Items: &items}

	body := bytes.NewBufferString(`{"id":1,"description":"New"}`)
	req := httptest.NewRequest(http.MethodPost, "/update", body)
	w := httptest.NewRecorder()

	api.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if (*api.Items)[0].Description != "New" {
		t.Errorf("expected updated description, got %s", (*api.Items)[0].Description)
	}
}

func TestAPI_Delete(t *testing.T) {
	items := []Item{{ID: 1, Description: "ToDelete"}}
	api := &API{Items: &items}

	body := bytes.NewBufferString(`{"id":1}`)
	req := httptest.NewRequest(http.MethodPost, "/delete", body)
	w := httptest.NewRecorder()

	api.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
	if len(*api.Items) != 0 {
		t.Errorf("expected item to be deleted")
	}
}
