package main

import (
	"sync"
)

type Store struct {
	mu     sync.RWMutex
	todos  map[int]Todo
	nextID int
}

func NewStore() *Store {
	return &Store{
		todos:  make(map[int]Todo),
		nextID: 1,
	}
}

func (s *Store) Create(todo Todo) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo.ID = s.nextID
	s.todos[todo.ID] = todo
	s.nextID++

	return todo
}

func (s *Store) GetAll() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	return todos
}

func (s *Store) GetByID(id int) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, found := s.todos[id]
	return todo, found
}

func (s *Store) Update(id int, updated Todo) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.todos[id]; !found {
		return Todo{}, false
	}

	updated.ID = id
	s.todos[id] = updated

	return updated, true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.todos[id]; !found {
		return false
	}

	delete(s.todos, id)
	return true
}
