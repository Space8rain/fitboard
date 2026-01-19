package handlers

import (
    "fmt"
    "net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "pong")
}

func Users(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `[{"id":1,"name":"Nikita"}]`)
}

func Trainers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `[{"id":2,"name":"Trainer"}]`)
}
