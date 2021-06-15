package handler

import (
    "net/http"
    "log"
    "encoding/json"

    "github.com/go-chi/render"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello world"))
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")

    if q == "" {
        w.Write([]byte("no query"))
        return
    }

    w.Write([]byte(q))
}

func JsonGetHandler(w http.ResponseWriter, r *http.Request) {
    // Map
    data := map[string]string{
        "message": "hello",
    }

    render.JSON(w, r, data)
}

func JsonPostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Content-Type") != "application/json" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    type Input struct {
        Name string `json:"name"`
    }
    var input Input
    err := json.NewDecoder(r.Body).Decode(&input)

    if err != nil {
        log.Printf("%s", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    data := map[string]string{
        "result": input.Name,
    }

    render.JSON(w, r, data)
}
