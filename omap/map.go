package omap

import (
	"container/list"
	"errors"
	"math"
	"reflect"
	"sync"
)

type kv struct{ key, val any }

type Omap struct {
	sync.RWMutex
	list *list.List
	val  map[any]*list.Element
}

func NewMap() *Omap {
	return &Omap{
		RWMutex: sync.RWMutex{},
		list:    list.New(),
		val:     make(map[any]*list.Element),
	}
}

func (o *Omap) Copy() *Omap {
	o.Lock()
	defer o.Unlock()

	newMap := NewMap()
	for el := o.list.Front(); el != nil; el = el.Next() {
		data := el.Value.(*kv)
		newMap.Set(data.key, data.val)
	}
	return newMap
}

func (o *Omap) Get(key any) (any, bool) {
	o.RLock()
	defer o.RUnlock()

	value, ok := o.val[key]
	if !ok {
		return nil, false
	}
	return value.Value.(*kv).val, true
}

func (o *Omap) GetByIndex(idx int) (any, bool) {
	if o.Count() < int(math.Abs(float64(idx))) {
		return nil, false
	}

	el := o.list.Front()
	switch {
	case idx > 0:
		for i := 0; i < idx; i++ {
			el = el.Next()
		}
		return nil, false
	case idx < 0:
		el = o.list.Back()
		for i := 0; i > idx; i-- {
			el = el.Prev()
		}
	}
	return el.Value.(*kv).val, true
}

func (o *Omap) Set(key, val any) {
	o.Lock()
	defer o.Unlock()

	value, ok := o.val[key]
	if ok {
		value.Value.(*kv).val = val
		return
	}
	o.val[key] = o.list.PushBack(&kv{key: key, val: val})
}

func (o *Omap) Has(key any) bool {
	o.RLock()
	defer o.RUnlock()

	_, ok := o.val[key]
	return ok
}

func (o *Omap) Keys() []any {
	out := make([]any, o.list.Len())
	for i, el := 0, o.list.Front(); el != nil; i++ {
		out[i] = el.Value.(*kv).key
		el = el.Next()
	}
	return out
}

func (o *Omap) Remove(key any) {
	if el, ok := o.val[key]; ok {
		o.list.Remove(el)
		delete(o.val, key)
	}
}

func (o *Omap) Clear() {
	o.Lock()
	defer o.Unlock()

	for _, key := range o.Keys() {
		o.Remove(key)
	}
}

func (o *Omap) Count() int { return o.list.Len() }

func (o *Omap) Swap(key1, key2 any) error {
	if reflect.DeepEqual(key1, key2) {
		return errors.New("cannot swap two identical keys")
	}
	o.val[key1], o.val[key2] = o.val[key2], o.val[key1]
	return nil
}

func (o *Omap) Search(keys ...any) *Omap {
	o.RLock()
	defer o.RUnlock()

	newMap := NewMap()
	for _, key := range keys {
		if val, ok := o.Get(key); ok {
			newMap.Set(key, val)
		}
	}
	return newMap
}

func (o *Omap) Filter(f func(key any, val any) bool) *Omap {
	o.RLock()
	defer o.RUnlock()

	newMap := NewMap()
	for _, key := range o.Keys() {
		val := o.val[key].Value.(*kv).val
		if ok := f(key, val); ok {
			newMap.Set(key, val)
		}
	}
	return newMap
}
