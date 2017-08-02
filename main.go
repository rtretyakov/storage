package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"github.com/rtretyakov/storage/app"
	"log"
	"github.com/go-chi/chi/middleware"
	"time"
	"flag"
)

func main() {
	cleanupInterval := flag.Duration("i", time.Minute, "cleanup interval")
	flag.Parse()

	a := app.NewApp(*cleanupInterval)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Put("/items/{key}", a.CreateItemHandler)
	r.Post("/items/{key}/incr", a.IncrItemHandler)
	r.Get("/items/{key}", a.GetItemHandler)

	log.Println("Starting HTTP Server...")
	http.ListenAndServe(":8080", r)
}
