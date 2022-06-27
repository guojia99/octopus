package omap

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
)

// Omap map with read-write locks, sequential keys
type Omap struct {
	lock sync.RWMutex
	list *list.List
	val  map[any]*list.Element
}

// NewMap Create a new Omap
func NewMap() *Omap {
	return &Omap{
		lock: sync.RWMutex{},
		list: list.New(),
		val:  make(map[any]*list.Element),
	}
}

// Copy Completely copy a new Omap
func (o *Omap) Copy() *Omap {
	o.lock.Lock()
	defer o.lock.Unlock()

	newMap := NewMap()
	for el := o.list.Front(); el != nil; el = el.Next() {
		data := el.Value.(*kv)
		newMap.Set(data.key, data.val)
	}
	return newMap
}

// Set the value according to key and val
func (o *Omap) Set(key, val any) {
	o.lock.Lock()
	defer o.lock.Unlock()

	el, ok := o.val[key]
	if ok {
		el.Value.(*kv).val = val
		return
	}
	o.val[key] = o.list.PushBack(&kv{key: key, val: val})
}

// SetByOmap insert and overwrite update the kv value of another Omap
func (o *Omap) SetByOmap(o2 *Omap) {
	if o2 == nil {
		return
	}

	if o == o2 {
		return
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	o2.lock.RLock()
	defer o2.lock.RUnlock()

	for _, key := range o2.Keys() {
		val := o2.val[key].Value.(*kv).val
		if el, ok := o.val[key]; ok {
			el.Value.(*kv).val = val
			continue
		}
		o.val[key] = o.list.PushBack(&kv{key: key, val: val})
	}
}

// Get value by key
func (o *Omap) Get(key any) (any, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	el, ok := o.val[key]
	if !ok {
		return nil, false
	}
	return el.Value.(*kv).val, true
}

// Has Check if a key value exists
func (o *Omap) Has(key any) bool {
	o.lock.RLock()
	defer o.lock.RUnlock()

	_, ok := o.val[key]
	return ok
}

// GetByValue find out the key value based on the value
func (o *Omap) GetByValue(val any) (key any, ok bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	for el := o.list.Front(); el != nil; el = el.Next() {
		data := el.Value.(*kv)
		if reflect.DeepEqual(data.val, val) {
			return data.key, true
		}
	}
	return nil, false
}

// GetByIndex the key according to the first value of the map
func (o *Omap) GetByIndex(idx int) (any, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	if o.Len() < int(math.Abs(float64(idx))) {
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

// Remove key value
func (o *Omap) Remove(key any) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if el, ok := o.val[key]; ok {
		o.list.Remove(el)
		delete(o.val, key)
	}
}

// Clear Omap
func (o *Omap) Clear() {
	keys := o.Keys()

	o.lock.Lock()
	defer o.lock.Unlock()

	for _, key := range keys {
		el, _ := o.val[key]
		o.list.Remove(el)
		delete(o.val, key)
	}
}

// Keys extract the list of keys
func (o *Omap) Keys() []any {
	o.lock.RLock()
	defer o.lock.RUnlock()

	out := make([]any, o.list.Len())
	for i, el := 0, o.list.Front(); el != nil; i++ {
		out[i] = el.Value.(*kv).key
		el = el.Next()
	}
	return out
}

// Len Omap length
func (o *Omap) Len() int { return o.list.Len() }

// Swap two key
func (o *Omap) Swap(key1, key2 any) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if reflect.DeepEqual(key1, key2) {
		return errors.New("cannot swap two identical keys")
	}
	if _, ok := o.val[key1]; !ok {
		return fmt.Errorf("key1 %v does not exist", key1)
	}
	if _, ok := o.val[key2]; !ok {
		return fmt.Errorf("key2 %v does not exist", key2)
	}
	o.val[key1], o.val[key2] = o.val[key2], o.val[key1]
	return
}

// Search find a set that matches the keys
func (o *Omap) Search(keys ...any) *Omap {
	o.lock.RLock()
	defer o.lock.RUnlock()

	newMap := NewMap()
	for _, key := range keys {
		if el, ok := o.val[key]; ok {
			newMap.Set(key, el.Value.(*kv).val)
		}
	}
	return newMap
}

// Filter a collection that meets the requirements
func (o *Omap) Filter(f func(key any, val any) bool) *Omap {
	o.lock.RLock()
	defer o.lock.RUnlock()

	newMap := NewMap()
	for _, key := range o.Keys() {
		val := o.val[key].Value.(*kv).val
		if ok := f(key, val); ok {
			newMap.Set(key, val)
		}
	}
	return newMap
}

//func (o *Omap) Union(o2 *Omap) *Omap {
//	return nil
//}
//
//func (o *Omap) Intersection(o2 *Omap) *Omap {
//	return nil
//}
