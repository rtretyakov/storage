package app

import "time"

type item struct {
	value    interface{}
	expireAt time.Time
}

func newItem(value interface{}, ttl time.Duration) (i item) {
	i.value = value
	i.expireAt = time.Now().Add(ttl)
	return
}

func (i item) IsExpired() bool {
	return i.expireAt.Before(time.Now())
}
