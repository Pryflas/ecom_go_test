package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTodoSuccess(t *testing.T) {
	store := NewStore()
	req := CreateTodoRequest{
		Title:       "Купить молоко",
		Description: "В магазине",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handleCreateTodo(w, httpReq, store)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var created Todo
	json.Unmarshal(w.Body.Bytes(), &created)

	if created.Title != req.Title {
		t.Errorf("Expected title %q, got %q", req.Title, created.Title)
	}

	if created.ID != 1 {
		t.Errorf("Expected ID 1, got %d", created.ID)
	}
}

func TestCreateTodoValidationError(t *testing.T) {
	store := NewStore()
	req := CreateTodoRequest{
		Title:       "",
		Description: "Описание",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handleCreateTodo(w, httpReq, store)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTodosEmpty(t *testing.T) {
	store := NewStore()
	httpReq := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	handleGetTodos(w, httpReq, store)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var todos []Todo
	json.Unmarshal(w.Body.Bytes(), &todos)

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestGetTodoByIDSuccess(t *testing.T) {
	store := NewStore()
	todo := Todo{
		Title:       "Тестовая задача",
		Description: "Описание",
		Completed:   false,
	}
	created := store.Create(todo)

	httpReq := httptest.NewRequest("GET", "/todos/1", nil)
	w := httptest.NewRecorder()

	handleGetTodoByID(w, httpReq, store, created.ID)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var retrieved Todo
	json.Unmarshal(w.Body.Bytes(), &retrieved)

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}
}

func TestGetTodoByIDNotFound(t *testing.T) {
	store := NewStore()
	httpReq := httptest.NewRequest("GET", "/todos/999", nil)
	w := httptest.NewRecorder()

	handleGetTodoByID(w, httpReq, store, 999)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateTodoSuccess(t *testing.T) {
	store := NewStore()
	todo := Todo{
		Title:       "Старое название",
		Description: "Описание",
		Completed:   false,
	}
	created := store.Create(todo)

	updateReq := UpdateTodoRequest{
		Title:       "Новое название",
		Description: "Новое описание",
		Completed:   true,
	}

	body, _ := json.Marshal(updateReq)
	httpReq := httptest.NewRequest("PUT", "/todos/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handleUpdateTodo(w, httpReq, store, created.ID)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var updated Todo
	json.Unmarshal(w.Body.Bytes(), &updated)

	if updated.Title != updateReq.Title {
		t.Errorf("Expected title %q, got %q", updateReq.Title, updated.Title)
	}

	if !updated.Completed {
		t.Errorf("Expected completed to be true")
	}
}

func TestUpdateTodoNotFound(t *testing.T) {
	store := NewStore()
	updateReq := UpdateTodoRequest{
		Title:       "Новое название",
		Description: "Описание",
	}

	body, _ := json.Marshal(updateReq)
	httpReq := httptest.NewRequest("PUT", "/todos/999", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handleUpdateTodo(w, httpReq, store, 999)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteTodoSuccess(t *testing.T) {
	store := NewStore()
	todo := Todo{
		Title:       "Удаляемая задача",
		Description: "Описание",
	}
	created := store.Create(todo)

	httpReq := httptest.NewRequest("DELETE", "/todos/1", nil)
	w := httptest.NewRecorder()

	handleDeleteTodo(w, httpReq, store, created.ID)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	_, found := store.GetByID(created.ID)
	if found {
		t.Errorf("Expected todo to be deleted")
	}
}

func TestDeleteTodoNotFound(t *testing.T) {
	store := NewStore()
	httpReq := httptest.NewRequest("DELETE", "/todos/999", nil)
	w := httptest.NewRecorder()

	handleDeleteTodo(w, httpReq, store, 999)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDuplicateIDPrevention(t *testing.T) {
	store := NewStore()

	todo1 := Todo{Title: "Первая", Description: "Описание 1"}
	todo2 := Todo{Title: "Вторая", Description: "Описание 2"}

	created1 := store.Create(todo1)
	created2 := store.Create(todo2)

	if created1.ID == created2.ID {
		t.Errorf("IDs should be unique: %d and %d", created1.ID, created2.ID)
	}

	all := store.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(all))
	}
}
