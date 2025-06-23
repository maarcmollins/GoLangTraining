package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestToDoActor_ConcurrentAddAndGet(t *testing.T) {
	t.Parallel()
	actor := NewToDoActor([]Item{})
	const numGoroutines = 10
	const itemsPerGoroutine = 5

	var wg sync.WaitGroup

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < itemsPerGoroutine; i++ {
				actor.AddItem(fmt.Sprintf("Goroutine %d - Item %d", gid, i))
			}
		}(g)
	}

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = actor.GetItems()
		}()
	}

	wg.Wait()

	items := actor.GetItems()
	expected := numGoroutines * itemsPerGoroutine
	if len(items) != expected {
		t.Errorf("expected %d items, got %d", expected, len(items))
	}
}

func TestToDoActor_ConcurrentUpdateAndDelete(t *testing.T) {
	t.Parallel()
	initial := make([]Item, 20)
	for i := range initial {
		initial[i] = Item{ID: i + 1, Description: "desc"}
	}
	actor := NewToDoActor(initial)

	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			actor.UpdateItem(id, fmt.Sprintf("updated-%d", id), StatusStarted)
		}(i)
	}

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			actor.DeleteItem(id)
		}(i)
	}

	wg.Wait()

	items := actor.GetItems()
	if len(items) != 10 {
		t.Errorf("expected 10 items left, got %d", len(items))
	}
	for _, it := range items {
		if it.Status != StatusStarted || it.Description != fmt.Sprintf("updated-%d", it.ID) {
			t.Errorf("unexpected item: %+v", it)
		}
	}
}
