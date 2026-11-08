package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oleksandr-vlasov/task-manager/internal/database"
	"github.com/oleksandr-vlasov/task-manager/internal/handler"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("CRITICAL ERROR: DATABASE_URL environment variable is not set!")
	}

	db, err := database.NewDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Conn.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	taskHandler := handler.NewTaskHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok","message":"Task Manager API is running"}`)); err != nil {
			log.Printf("write response: %v", err)
		}
	})

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		handleTasks(w, r, taskHandler)
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		handleTasks(w, r, taskHandler)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Task Manager API starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handleTasks(w http.ResponseWriter, r *http.Request, h *handler.TaskHandler) {
	path := r.URL.Path

	if path == "/tasks" || path == "/tasks/" {
		switch r.Method {
		case http.MethodGet:
			h.GetTasks(w, r)
		case http.MethodPost:
			h.CreateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if strings.HasPrefix(path, "/tasks/") {
		switch r.Method {
		case http.MethodPut:
			h.UpdateTask(w, r)
		case http.MethodDelete:
			h.DeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
}
