package store

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestLoad_CreateAndGet(t *testing.T) {
	t.Parallel()
	const (
		baseURL     = "http://localhost:8080"
		numRequests = 100
		concurrency = 10
	)
	var wg sync.WaitGroup

	resp, err := http.Get(baseURL + "/get")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Skip("Server not running on :8080, skipping load test")
	}

	wg.Add(concurrency)
	for c := 0; c < concurrency; c++ {
		go func(cid int) {
			defer wg.Done()
			for i := 0; i < numRequests/concurrency; i++ {
				body, _ := json.Marshal(map[string]string{"description": "LoadTest"})
				resp, err := http.Post(baseURL+"/create", "application/json", bytes.NewReader(body))
				if err != nil {
					t.Errorf("create request failed: %v", err)
					continue
				}
				resp.Body.Close()
			}
		}(c)
	}
	wg.Wait()

	time.Sleep(500 * time.Millisecond)

	resp, err = http.Get(baseURL + "/get")
	if err != nil {
		t.Fatalf("failed to get items: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(items) < numRequests {
		t.Errorf("expected at least %d items, got %d", numRequests, len(items))
	}
}
