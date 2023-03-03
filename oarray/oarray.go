/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/3/3 下午2:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

package oarray

import (
	"fmt"
	"math/rand"
	"sync"
)

type OArray[T comparable] struct {
	lock sync.RWMutex
	val  []T
}

// NewOArray create a new array
func NewOArray[T comparable](val ...T) *OArray[T] {
	a := &OArray[T]{
		lock: sync.RWMutex{},
		val:  make([]T, 0),
	}
	if len(val) > 0 {
		a.val = append(a.val, val...)
	}
	return a
}

// Copy a new array without interfering with each other
func (a *OArray[T]) Copy() *OArray[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	newArray := NewOArray[T]()
	for _, val := range a.val {
		newArray.Set(val)
	}
	return newArray
}

// Len with array length
func (a *OArray[T]) Len() int { return len(a.val) }

// IsEmpty is the array not empty
func (a *OArray[T]) IsEmpty() bool { return a.Len() == 0 }

// Val returns the current raw data
func (a *OArray[T]) Val() []T { return a.Copy().val }

// Set insert 1-n data
func (a *OArray[T]) Set(value ...T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.val = append(a.val, value...)
}

// GetByIndex the data of a subscript
func (a *OArray[T]) GetByIndex(idx int) (val T, err error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if idx > a.Len() {
		err = fmt.Errorf("subscript out of range")
		return
	}
	return a.val[idx], nil
}

// GetSlice get a slice
func (a *OArray[T]) GetSlice(start, end int) ([]T, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if start > end {
		start, end = end, start
	}

	if end > a.Len() {
		return nil, fmt.Errorf("subscript out of range")
	}
	return a.val[start:end], nil
}

// HasValueFirst find the first index with the same value
func (a *OArray[T]) HasValueFirst(value T) (int, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for idx, val := range a.val {
		if val == value {
			return idx, true
		}
	}
	return 0, false
}

// HasValueAll find all index with the same value
func (a *OArray[T]) HasValueAll(value T) ([]int, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	var out = make([]int, 0)
	for idx, val := range a.val {
		if val == value {
			out = append(out, idx)
		}
	}
	if len(out) > 0 {
		return out, true
	}
	return nil, false
}

// Reverse flip the entire array
func (a *OArray[T]) Reverse() {
	a.lock.Lock()
	defer a.lock.Unlock()

	length := a.Len()
	for i := 0; i < length/2; i++ {
		li := length - i - 1
		a.val[i], a.val[li] = a.val[li], a.val[i]
	}
}

// RandomOne output a random
func (a *OArray[T]) RandomOne() T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.val[rand.Intn(len(a.val))]
}

// Filter a new array with func
func (a *OArray[T]) Filter(f func(value T) bool) *OArray[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	newArray := NewOArray[T]()
	for _, val := range a.val {
		if f(val) {
			newArray.Set(val)
		}
	}
	return newArray
}

// Remove delete the required value
func (a *OArray[T]) Remove(value T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	newVal := make([]T, 0)
	for _, val := range a.val {
		if val != value {
			newVal = append(newVal, val)
		}
	}
	a.val = newVal
}
