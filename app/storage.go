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

func (s *storage) Get(key string) (value interface{}, err error) {
	s.mu.RLock()
	i, ok := s.items[key]
	s.mu.RUnlock()

	if ok && !i.IsExpired() {
		return i.value, nil
	} else {
		return nil, errNotFound
	}

	return
}

func (s *storage) Incr(key string) (value interface{}, err error) {
	s.mu.Lock()
	i, ok := s.items[key]

	if ok && !i.IsExpired() {
		if _, ok := i.value.(float64); ok {
			i.value = i.value.(float64) + 1
			s.items[key] = i
			s.mu.Unlock()
		} else {
			s.mu.Unlock()
			return nil, errWrongType
		}

		return i.value, nil
	} else {
		s.mu.Unlock()
		return nil, errNotFound
	}

	return
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