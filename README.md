# ToDoApp

A simple concurrent-safe ToDo application written in Go.  
Supports both command-line and HTTP API/web frontend usage, using the Actor/CSP pattern for safe concurrent access.

---

## Getting Started

### 1. **Build the Application**

```sh
go build -o todoapp ./cmd/todoapp.go
```

### 2. **Command-Line Usage**

You can manage your to-do list directly from the command line.

#### **Flags**

- `-file`  
  Path to the JSON file for storing your to-do list.  
  **Default:** `todos.json`

- `-add`  
  Add a new to-do item.  
  **Example:** `-add="Buy milk"`

- `-update-id`  
  The ID of the item you want to update.

- `-update-text`  
  The new description for the item.

- `-update-status`  
  The new status for the item (`not started`, `started`, `completed`).

- `-delete-id`  
  The ID of the item you want to delete.

- `-start-server`  
  Start the HTTP API/web server.  
  **Example:** `-start-server`

#### **Examples**

Add a new item:
```sh
./todoapp -add="Buy milk"
```

Update an item:
```sh
./todoapp -update-id=1 -update-text="Buy oat milk"
```

Update an item's status:
```sh
./todoapp -update-id=1 -update-status="completed"
```

Delete an item:
```sh
./todoapp -delete-id=1
```

Show the current list:
```sh
./todoapp
```

---

### 3. **HTTP API & Web Frontend Usage**

Start the server:
```sh
./todoapp -start-server
```

The server will listen on [http://localhost:8080](http://localhost:8080).

#### **API Endpoints**

- `POST /create`  
  Create a new item.  
  **Body:** `{"description": "Task description"}`

- `GET /get`  
  Get all items.

- `POST /update`  
  Update an item.  
  **Body:** `{"id": 1, "description": "New desc", "status": "completed"}`

- `POST /delete`  
  Delete an item.  
  **Body:** `{"id": 1}`

#### **Web Frontend**

- `/static/create.html` — Create a new item
- `/static/get.html` — View all items
- `/static/update.html` — Update an item
- `/static/delete.html` — Delete an item
- `/static/about.html` — About page
- `/list` — Dynamic HTML list of all items

---

### 4. **Graceful Shutdown**

When running the server, press `Ctrl+C` to gracefully shut down and save your to-do list to disk.

---

## Concurrency

This app uses the Actor/CSP pattern (via channels and a dedicated goroutine) to ensure all reads and writes to the to-do list are safe and race-free, even under heavy concurrent access.

---

## Development & Testing

Run all tests (including concurrency and load tests):

```sh
go test -race ./internal/store/...
```

---

## License

MIT