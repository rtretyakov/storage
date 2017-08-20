package app

import (
	"time"
	"sync"
)

type storage struct {
	items map[string]item
	mu    sync.RWMutex
}

func newStorage() *storage {
	s := new(storage)
	s.items = make(map[string]item)

	return s
}

func (s *storage) Set(key string, value interface{}, ttl time.Duration) {
	s.mu.Lock()
	s.items[key] = newItem(value, ttl)
	s.mu.Unlock()
}

func (s *storage) Get(key string) (item item, err error) {
	s.mu.RLock()
	item, ok := s.items[key]
	s.mu.RUnlock()

	if ok && !item.IsExpired() {
		return item, nil
	} else {
		return item, errNotFound
	}

	return
}

func (s *storage) Incr(key string) (item item, err error) {
	s.mu.Lock()
	item, ok := s.items[key]

	if ok && !item.IsExpired() {
		if _, ok := item.value.(float64); ok {
			item.value = item.value.(float64) + 1
			s.items[key] = item
			s.mu.Unlock()
		} else {
			s.mu.Unlock()
			return item, errWrongType
		}

		return item, nil
	} else {
		s.mu.Unlock()
		return item, errNotFound
	}

	return
}

func (s *storage) Delete(key string) {
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
}

func (s *storage) Clean() {
	s.mu.Lock()
	for key, item := range s.items {
		if item.IsExpired() {
			delete(s.items, key)
		}
	}
	s.mu.Unlock()
}