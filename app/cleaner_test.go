package app

import (
	"testing"
	"time"
)

func TestNewCleaner(t *testing.T) {
	c := newCleaner(time.Minute)

	if c.interval != time.Minute {
		t.Errorf("Expect minute cleanup interval, but got %v", c.interval)
	}
}

func TestCleaner_Start(t *testing.T) {
	c := newCleaner(time.Nanosecond)

	storage := newStorage()

	expiredKey := "testkey1"
	expiredValue := "testvalue1"

	key := "testkey2"
	value := "testvalue2"

	storage.Set(expiredKey, expiredValue, -1)
	storage.Set(key, value, time.Minute)

	time.Sleep(time.Microsecond)

	go c.Start(storage)

	time.Sleep(time.Second)

	_, err := storage.Get(expiredKey)
	if err != errNotFound {
		t.Errorf("Expect not found error, but got %v", err)
	}

	actualItem, err := storage.Get(key)
	if err != nil {
		t.Errorf("Error on get value after cleanup: %v", err)
	}

	if actualItem.value != value {
		t.Errorf("Expect %s value, but got %v", value, actualItem.value)
	}

	if len(storage.items) != 1 {
		t.Errorf("Storage should has one item, but has %d items", len(storage.items))
	}
}
