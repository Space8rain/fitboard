package handlers

import (
	"encoding/json"
	"fitboard/backend/internal/db"
	"fmt"
	"net/http"
)

type Handlers struct{}

func New() *Handlers {
    return &Handlers{}
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "pong")
}

func (h *Handlers) Users(w http.ResponseWriter, r *http.Request) {
    users, err := db.Repo.GetAll()
    if err != nil {
        fmt.Println("DB ERROR:", err) // ← добавляем лог
        http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}



func (h *Handlers) Trainers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `[{"id":2,"name":"Trainer"}]`)
}
