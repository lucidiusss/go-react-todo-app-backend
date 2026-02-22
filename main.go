package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("APP_ENV") == "production" {
		godotenv.Load("/etc/myapp/.env")
	} else {
		godotenv.Load()
	}

	filename := os.Getenv("TASKS_FILE")

	storage = NewStorage(filename)

	if filename == "" {
		filename = "tasks.json"
	}

	storage.LoadFromFile(filename)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGIN"), "https://todo.lucidiusss.lol", "http://todo.lucidiusss.lol", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", GetAllTasks)
			r.Post("/", CreateTask)

			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/", DeleteTask)
				r.Put("/", RenameTask)
				r.Post("/toggle", ToggleTask)
			})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
