package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// storage initialization
var storage *Storage

// helpers
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// get all tasks handler
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := storage.GetAll()

	respondWithJSON(w, http.StatusOK, tasks)
}

// create task handler
func CreateTask(w http.ResponseWriter, r *http.Request) {
	// инициализируем запрос
	var req CreateTaskRequest

	// читаем json из тела запроса
	err := json.NewDecoder(r.Body).Decode(&req)
	// возвращаем ошибку
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	// закрываем тело запроса
	defer r.Body.Close()

	//проверяем, не пустой ли таск
	if req.Title == "" {
		respondWithError(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	for _, t := range storage.todos {
		if req.Title == t.Title {
			respondWithError(w, http.StatusBadRequest, "Task already exists")
			return
		}
	}
	// создаем таск
	task := storage.Create(req.Title)

	// возвращаем
	respondWithJSON(w, http.StatusCreated, task)

}

// delete task handler
func DeleteTask(w http.ResponseWriter, r *http.Request) {

	// get id from url

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// deleting task
	deletedTask, err := storage.Delete(id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, deletedTask)
}

func RenameTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req CreateTaskRequest

	// read request body
	error := json.NewDecoder(r.Body).Decode(&req)
	// return error
	if error != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	//close request body
	defer r.Body.Close()

	//check if title is not empty
	if req.Title == "" {
		respondWithError(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	for _, t := range storage.todos {
		if t.Title == req.Title {
			respondWithError(w, http.StatusBadRequest, "Task already exists")
			return
		}
	}
	// rename task
	renamedTask, err := storage.Rename(id, req.Title)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// return renamedTask
	respondWithJSON(w, http.StatusOK, renamedTask)
}

// toggle "completed" state handler
func ToggleTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	// handle invalid id
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	// toggle task
	toggledTask, err := storage.Toggle(id)

	//handle error
	if err != nil {
		respondWithError(w, http.StatusNotFound, "task was not found")
	}
	// return
	respondWithJSON(w, http.StatusOK, toggledTask)
}
