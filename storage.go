package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// storage struct
type Storage struct {
	mu       sync.RWMutex
	todos    []Todo
	lastID   int
	filename string
}

// storage constructor
func NewStorage(filename string) *Storage {
	return &Storage{
		todos:    make([]Todo, 0),
		lastID:   0,
		filename: filename,
	}
}

// getall method
func (s *Storage) GetAll() []Todo {
	s.mu.RLock()         // lock storage
	defer s.mu.RUnlock() // unlock storage

	todosCopy := make([]Todo, len(s.todos)) //create a slice of todos
	copy(todosCopy, s.todos)
	return todosCopy // return
}

// create method
func (s *Storage) Create(title string) Todo {
	s.mu.Lock()         // lock
	defer s.mu.Unlock() // unlock

	s.lastID++ // increment last id

	task := Todo{
		ID:        s.lastID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	} // new task

	s.todos = append(s.todos, task) // add new task to existing ones
	err := s.SaveToFile(s.filename)
	if err != nil {
		log.Printf("warning: failed to save after create: %v", err)
	}
	return task
}

// delete method
func (s *Storage) Delete(id int) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskIndex := -1 // make taskIndex -1 before actually changing its value
	var deletedTask Todo

	// search for task
	for i, task := range s.todos {
		if task.ID == id {
			taskIndex = i
			deletedTask = task
			break
		}
	}

	// task was not found
	if taskIndex == -1 {
		return Todo{}, fmt.Errorf("task with id %d not found", id)
	}

	// delete if found
	s.todos = append(s.todos[:taskIndex], s.todos[taskIndex+1:]...)

	err := s.SaveToFile(s.filename)
	if err != nil {
		log.Printf("warning: failed to save after deleting: %v", err)
	}

	// return
	return deletedTask, nil
}

func (s *Storage) Rename(id int, title string) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskIndex := -1
	var updatedTask Todo

	for i, task := range s.todos {
		if task.ID == id {
			s.todos[i].Title = title
			taskIndex = i
			updatedTask = s.todos[i]
			break
		}
	}

	if taskIndex == -1 {
		return Todo{}, fmt.Errorf("task with id %d not found", id)
	}

	if title == "" {
		return Todo{}, fmt.Errorf("title cannot be empty")
	}

	err := s.SaveToFile(s.filename)
	if err != nil {
		log.Printf("warning: failed to save after renaming: %v", err)
	}
	return updatedTask, nil
}

// toggle method
func (s *Storage) Toggle(id int) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	taskIndex := -1
	var updatedTask Todo

	for i, task := range s.todos {
		if task.ID == id {
			s.todos[i].Completed = !task.Completed
			taskIndex = i
			updatedTask = s.todos[i]
			break
		}
	}
	if taskIndex == -1 {
		return Todo{}, fmt.Errorf("task with id %d not found", id)
	}

	err := s.SaveToFile(s.filename)
	if err != nil {
		log.Printf("warning: failed to save after deleting: %v", err)
	}
	return updatedTask, nil
}

func (s *Storage) SaveToFile(filename string) error {

	tasks := make([]Todo, len(s.todos))
	copy(tasks, s.todos)
	// 1. Создаем файл
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to save to file: %w", err)
	}
	defer file.Close()

	// 2. Кодируем JSON прямо в файл
	err = json.NewEncoder(file).Encode(tasks)
	if err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}

func (s *Storage) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil //
		}
		return err
	}

	var tasks []Todo
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// check if file is empty
	if len(tasks) == 0 {
		s.todos = make([]Todo, 0) // empty slice
		s.lastID = 0
		return nil
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	s.todos = tasks
	s.lastID = maxID
	return nil
}
