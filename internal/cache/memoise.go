package cache

import (
	"sync"
	"time"
)

type Memo[T any] struct {
	mu       sync.Mutex
	lifespan time.Duration
	expiry   time.Time
	fn       func() (*T, error)
	data     *T
}

func Memoise[T any](lifespan time.Duration, fn func() (*T, error)) *Memo[T] {
	c := &Memo[T]{
		lifespan: lifespan,
		expiry:   time.Now().Add(lifespan),
		fn:       fn,
	}

	return c
}

func (c *Memo[T]) Expire() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.expiry = time.Now().Add(-time.Second)
}

func (c *Memo[T]) Get() (*T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.data != nil && time.Now().Before(c.expiry) {
		return c.data, nil
	}

	data, err := c.fn()
	if err != nil {
		return nil, err
	}

	c.data = data
	c.expiry = time.Now().Add(c.lifespan)

	return c.data, nil
}
