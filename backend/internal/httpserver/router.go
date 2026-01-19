package httpserver

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    "fitboard/backend/internal/handlers"
)

func NewRouter(h *handlers.Handlers) *chi.Mux {
    r := chi.NewRouter()

    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    r.Get("/api/ping", h.Ping)
    r.Get("/api/users", h.Users)
    r.Get("/api/trainers", h.Trainers)

    return r
}
