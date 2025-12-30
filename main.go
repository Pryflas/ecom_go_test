package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func main() {
	store := NewStore()

	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		if r.Method == http.MethodPost {
			handleCreateTodo(w, r, store)
		} else if r.Method == http.MethodGet {
			handleGetTodos(w, r, store)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/todos/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			handleGetTodoByID(w, r, store, id)
		} else if r.Method == http.MethodPut {
			handleUpdateTodo(w, r, store, id)
		} else if r.Method == http.MethodDelete {
			handleDeleteTodo(w, r, store, id)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request, store *Store) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var req CreateTodoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" || strings.TrimSpace(req.Title) == "" {
		http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}
	todo := Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	created := store.Create(todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func handleGetTodos(w http.ResponseWriter, r *http.Request, store *Store) {
	todos := store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func handleGetTodoByID(w http.ResponseWriter, r *http.Request, store *Store, id int) {
	todo, found := store.GetByID(id)
	if !found {
		http.Error(w, `{"error":"Todo not found"}`, http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func handleUpdateTodo(w http.ResponseWriter, r *http.Request, store *Store, id int) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var req UpdateTodoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" || strings.TrimSpace(req.Title) == "" {
		http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	updated, found := store.Update(id, Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	})

	if !found {
		http.Error(w, `{"error":"Todo not found"}`, http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

func handleDeleteTodo(w http.ResponseWriter, r *http.Request, store *Store, id int) {
	found := store.Delete(id)
	if !found {
		http.Error(w, `{"error":"Todo not found"}`, http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func logRequest(r *http.Request) {
	log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
}
