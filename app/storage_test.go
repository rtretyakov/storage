package app

import (
	"testing"
	"time"
)

func TestNewStorage(t *testing.T) {
	storage := newStorage()

	if len(storage.items) != 0 {
		t.Errorf("Expected empty storage, but there is %d items", len(storage.items))
	}
}

func TestStorage_SetAndSuccessGet(t *testing.T) {
	storage := newStorage()
	key := "testkey"
	value := "testval"

	storage.Set(key, value, time.Minute)

	actualValue, err := storage.Get(key)
	if err != nil {
		t.Errorf("Error on get item value: %v", err)
	}

	if actualValue != value {
		t.Errorf("Expected %s, but got %v", value, actualValue)
	}
}

func TestStorage_GetExpired(t *testing.T) {
	storage := newStorage()
	key := "testkey"
	value := "testval"

	storage.Set(key, value, -1)

	actualValue, err := storage.Get(key)
	if actualValue != nil {
		t.Errorf("Value of expired item should be nil, but got %v", actualValue)
	}

	if err != errNotFound {
		t.Errorf("Error on get expired item's value should be errNotFound, but got: %v", err)
	}
}

func TestStorage_GetNotExist(t *testing.T) {
	storage := newStorage()
	key := "testkey"

	actualValue, err := storage.Get(key)
	if actualValue != nil {
		t.Errorf("Value of not existing item should be nil, but got %v", actualValue)
	}

	if err != errNotFound {
		t.Errorf("Error on get not existing item's value should be errNotFound, but got: %v", err)
	}
}

func TestStorage_Incr(t *testing.T) {
	storage := newStorage()
	key := "testkey"
	value := 1.3

	storage.Set(key, value, time.Minute)

	incrementedValue, err := storage.Incr("testkey")
	if err != nil {
		t.Errorf("Error on increment: %v", err)
	}

	if incrementedValue != value + 1 {
		t.Errorf("Incremented value should be %f, but got %v", value + 1, incrementedValue)
	}
}

func TestStorage_IncrWrongType(t *testing.T) {
	storage := newStorage()
	key := "testkey"
	value := "testvalue"

	storage.Set(key, value, time.Minute)

	incrementedValue, err := storage.Incr("testkey")
	if incrementedValue != nil {
		t.Errorf("Incremented wrong type value should be nil, but got %v", incrementedValue)
	}

	if err != errWrongType {
		t.Errorf("Error on increment wrong type value should be errWrongType, but got: %v", err)
	}
}

func TestStorage_IncrNotFound(t *testing.T) {
	storage := newStorage()

	incrementedValue, err := storage.Incr("testkey")
	if incrementedValue != nil {
		t.Errorf("Result of increment not existing item should be nil, but got %v", incrementedValue)
	}

	if err != errNotFound {
		t.Errorf("Error on increment not existing item's value should be errNotFound, but got: %v", err)
	}
}

func TestStorage_Clean(t *testing.T) {
	storage := newStorage()

	expiredKey := "testkey1"
	expiredValue := "testvalue1"

	key := "testkey2"
	value := "testvalue2"

	storage.Set(expiredKey, expiredValue, -1)
	storage.Set(key, value, time.Minute)

	storage.Clean()

	_, err := storage.Get(expiredKey)
	if err != errNotFound {
		t.Errorf("Expect not found error, but got %v", err)
	}

	actualValue, err := storage.Get(key)
	if err != nil {
		t.Errorf("Error on get value after cleanup: %v", err)
	}

	if actualValue != value {
		t.Errorf("Expect %s value, but got %v", value, actualValue)
	}

	if len(storage.items) != 1 {
		t.Errorf("Storage should has one item, but has %d items", len(storage.items))
	}
}