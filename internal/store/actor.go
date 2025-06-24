package store

import (
	"time"
)

type actorMsg interface{}

type getItemsMsg struct {
	reply chan []Item
}
type addItemMsg struct {
	description string
	reply       chan Item
}
type updateItemMsg struct {
	id          int
	description string
	status      string
	reply       chan bool
}
type deleteItemMsg struct {
	id    int
	reply chan bool
}

type ToDoActor struct {
	inbox chan actorMsg
}

func NewToDoActor(initial []Item) *ToDoActor {
	inbox := make(chan actorMsg)
	items := initial

	go func() {
		for msg := range inbox {
			switch m := msg.(type) {
			case getItemsMsg:
				cp := make([]Item, len(items))
				copy(cp, items)
				m.reply <- cp
			case addItemMsg:
				newItem := Item{
					ID:          nextID(items),
					Description: m.description,
					Status:      StatusNotStarted,
					CreatedAt:   time.Now(),
				}
				items = append(items, newItem)
				m.reply <- newItem
			case updateItemMsg:
				updated := false
				for i := range items {
					if items[i].ID == m.id {
						if m.description != "" {
							items[i].Description = m.description
						}
						if m.status != "" {
							items[i].Status = m.status
						}
						updated = true
						break
					}
				}
				m.reply <- updated
			case deleteItemMsg:
				deleted := false
				for i := range items {
					if items[i].ID == m.id {
						items = append(items[:i], items[i+1:]...)
						deleted = true
						break
					}
				}
				m.reply <- deleted
			}
		}
	}()
	return &ToDoActor{inbox: inbox}
}

func (a *ToDoActor) GetItems() []Item {
	reply := make(chan []Item)
	a.inbox <- getItemsMsg{reply}
	return <-reply
}
func (a *ToDoActor) AddItem(description string) Item {
	reply := make(chan Item)
	a.inbox <- addItemMsg{description, reply}
	return <-reply
}
func (a *ToDoActor) UpdateItem(id int, description, status string) bool {
	reply := make(chan bool)
	a.inbox <- updateItemMsg{id, description, status, reply}
	return <-reply
}
func (a *ToDoActor) DeleteItem(id int) bool {
	reply := make(chan bool)
	a.inbox <- deleteItemMsg{id, reply}
	return <-reply
}
