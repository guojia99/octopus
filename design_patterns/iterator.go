package dp

import (
	"context"
)

type selfIteratorFunc func(idx int) (ok bool, val interface{})
type interatorFunc func() (val interface{}, ok bool)

func NewIterator(data ...interface{}) Iterator {
	return &iterator{data: data}
}

func NewSelfIterator(f selfIteratorFunc, maxIndex int) Iterator {
	return &selfIterator{f: f, maxIndex: maxIndex}
}

type Iterator interface {
	Reset()
	HasNext() bool
	Next() interface{}
	Do(func(idx int, val interface{}) bool)
	Iterator() interatorFunc
	IteratorChan(ctx context.Context) <-chan interface{}
}

type iterator struct {
	index int
	data  []interface{}
}

func (c *iterator) Reset() { c.index = 0 }

func (c *iterator) HasNext() bool { return c.index < len(c.data) }

func (c *iterator) Next() interface{} {
	v := c.data[c.index]
	c.index++
	return v
}

func (c *iterator) Iterator() interatorFunc {
	return func() (val interface{}, ok bool) {
		if !c.HasNext() {
			return
		}
		val, ok = c.Next(), true
		return
	}
}

func (c *iterator) IteratorChan(ctx context.Context) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for {
			if !c.HasNext() {
				close(ch)
				return
			}
			v := c.Next()
			select {
			case <-ctx.Done():
				return
			case ch <- v:
			}
		}
	}()
	return ch
}

func (c *iterator) Do(fn func(idx int, data interface{}) bool) {
	for {
		if !c.HasNext() {
			return
		}
		v := c.Next()
		if ok := fn(c.index-1, v); !ok {
			return
		}
	}
}

type selfIterator struct {
	f        selfIteratorFunc
	index    int
	maxIndex int
}

func (c *selfIterator) Reset() { c.index = 0 }

func (c *selfIterator) HasNext() bool { return c.index < c.maxIndex }

func (c *selfIterator) Next() interface{} {
	ok, v := c.f(c.index)
	if !ok {
		return nil
	}
	c.index++
	return v
}

func (c *selfIterator) Do(fn func(idx int, val interface{}) bool) {
	for {
		if !c.HasNext() {
			return
		}
		v := c.Next()
		if v == nil {
			return
		}
		if ok := fn(c.index, v); !ok {
			return
		}
	}
}

func (c *selfIterator) Iterator() interatorFunc {
	return func() (val interface{}, ok bool) {
		if !c.HasNext() {
			return
		}
		if val = c.Next(); val == nil {
			ok = true
		}
		return
	}
}

func (c *selfIterator) IteratorChan(ctx context.Context) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for {
			if !c.HasNext() {
				close(ch)
				return
			}
			v := c.Next()
			select {
			case <-ctx.Done():
				return
			case ch <- v:
			}
		}
	}()
	return ch
}
