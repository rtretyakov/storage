package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rtretyakov/storage/app"
	"log"
	"net/http"
	"time"
)

func main() {
	cleanupInterval := flag.Duration("i", time.Minute, "cleanup interval")
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	a := app.NewApp(*cleanupInterval)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Put("/items/{key}", a.CreateItemHandler)
	r.Post("/items/{key}/incr", a.IncrItemHandler)
	r.Get("/items/{key}", a.GetItemHandler)
	r.Delete("/items/{key}", a.DeleteItemHandler)

	log.Printf("Starting HTTP Server on %d port...", *port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}
