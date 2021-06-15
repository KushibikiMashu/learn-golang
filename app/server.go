package main

import (
    "net/http"

    "example.com/backend/handler"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/httplog"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Logger)

    l := httplog.NewLogger("app", httplog.Options{
        JSON: true,
    })
    r.Use(httplog.RequestLogger(l))

    r.Get("/", handler.IndexHandler)
    r.Get("/echo", handler.EchoHandler)
    r.Route("/json", func(r chi.Router) {
        r.Get("/", handler.JsonGetHandler)
        r.Post("/", handler.JsonPostHandler)
    })

    http.ListenAndServe(":3333", r)
}
