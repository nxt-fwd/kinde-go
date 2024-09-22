package testutil

import (
	"fmt"
	"sync"
)

type CallCount struct {
	mu     sync.Mutex
	record map[string]int
}

func (*CallCount) key(method, path string) string {
	return fmt.Sprintf("%s - %s", method, path)
}

func (c *CallCount) Inc(method, path string) int {
	key := c.key(method, path)
	count := c.Get(method, path)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.record[key] = count + 1
	return c.record[key]
}

func (c *CallCount) Get(method, path string) int {
	key := c.key(method, path)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.record == nil {
		c.record = map[string]int{}
	}

	if _, ok := c.record[key]; !ok {
		c.record[key] = 0
	}

	return c.record[key]
}
