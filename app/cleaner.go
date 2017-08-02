package app

import "time"

type cleaner struct {
	interval time.Duration
}

func newCleaner(interval time.Duration) *cleaner {
	c := new(cleaner)
	c.interval = interval
	return c
}

func (c *cleaner) Start(storage *storage) {
	ticker := time.NewTicker(c.interval)
	for {
		<-ticker.C
		storage.Clean()
	}
}
