package app

import (
	"net/http"
	"encoding/json"
	"time"
	"log"
	"github.com/go-chi/chi"
)

type app struct {
	storage *storage
}

func NewApp(cleanupInterval time.Duration) *app {
	app := new(app)
	app.storage = newStorage()
	cleaner := newCleaner(cleanupInterval)
	go cleaner.Start(app.storage)

	return app
}

func (app *app) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	rq := new(createItemRequest)
	err := json.NewDecoder(r.Body).Decode(rq)
	if err != nil {
		log.Printf("Cannot decode JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	app.storage.Set(chi.URLParam(r, "key"), rq.Value, time.Second*time.Duration(rq.Ttl))
	w.WriteHeader(http.StatusOK)
}

func (app *app) GetItemHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	value, err := app.storage.Get(key)
	if err != nil {
		if err == errNotFound {
			log.Printf("Item %s is not found", key)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Cannot get item %s: %v", key, err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	resp := new(getItemResponse)
	resp.Value = value
	respB, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Cannot marshal JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respB)
}

func (app *app) IncrItemHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	value, err := app.storage.Incr(key)
	if err != nil {
		if err == errNotFound {
			log.Printf("Item %s is not found", key)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Cannot increment item %s: %v", key, err)
			w.WriteHeader(500)
		}
		return
	}

	resp := new(incrItemResponse)
	resp.Value = value.(float64)
	respB, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Cannot marshal JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respB)
}
