package oarray

import (
	"sync"
)

type Array struct {
	lock sync.RWMutex
	val  []any
}

func NewArray(val ...any) *Array {
	a := &Array{
		lock: sync.RWMutex{},
		val:  make([]any, 0),
	}
	if len(val) > 0 {
		a.val = append(a.val, val)
	}
	return a
}

func (a *Array) Len() int { return len(a.val) }

func (a *Array) Reverse() {
	a.lock.Lock()
	defer a.lock.RUnlock()

	length := a.Len()
	for i := 0; i < length/2; i++ {
		li := length - i - 1
		a.val[i], a.val[li] = a.val[li], a.val[i]
	}
}
