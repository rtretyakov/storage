package app

import (
	"testing"
	"time"
)

func TestNewItem(t *testing.T) {
	value := "testval"
	item := newItem(value, time.Minute)

	if item.value != value {
		t.Errorf("value is not %s", value)
	}

	if item.expireAt.IsZero() {
		t.Errorf("TTL is zero")
	}
}

func TestItem_IsExpired(t *testing.T) {
	value := "testval"
	tests := []struct{
		ttl time.Duration
		expectedResult bool
	}{
		{time.Minute, false},
		{time.Hour, false},
		{-1, true},
	}

	assert := func(expectedResult bool, actualResult bool) {
		if expectedResult != actualResult {
			t.Errorf("Excepted %t result of IsExpired, but got %t", expectedResult, actualResult)
		}
	}

	for _, test := range tests {
		item := newItem(value, test.ttl)
		assert(test.expectedResult, item.IsExpired())
	}

	// Expired
	item := newItem(value, time.Nanosecond)
	time.Sleep(time.Microsecond)
	assert(true, item.IsExpired())
}
