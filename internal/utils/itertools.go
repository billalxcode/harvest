package utils

import "io"

type Cycler[T any] struct {
	items []T
	index int
}

func NewCycler[T any](items []T) *Cycler[T] {
	if items == nil {
		items = []T{}
	}
	return &Cycler[T]{
		items: items,
		index: 0,
	}
}

func (c *Cycler[T]) Length() int {
	return len(c.items)
}

func (c *Cycler[T]) Next() (T, error) {
	if c.index >= len(c.items) {
		return *new(T), io.EOF
	}

	item := c.items[c.index]
	c.index = (c.index + 1) % len(c.items)
	return item, nil
}
