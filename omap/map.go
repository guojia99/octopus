package omap

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"

	json "github.com/json-iterator/go"
)

type kv struct {
	key any
	val any
}

// Map with read-write locks, sequential keys
type Map struct {
	lock sync.RWMutex
	list *list.List
	val  map[any]*list.Element
}

// NewMap Create a new Map
func NewMap() *Map {
	return &Map{
		lock: sync.RWMutex{},
		list: list.New(),
		val:  make(map[any]*list.Element),
	}
}

// Val get all metadata
func (o *Map) Val() map[any]any {
	o.lock.RLock()
	defer o.lock.RUnlock()

	var out = make(map[any]any)

	for key := range o.val {
		out[key] = o.val[key].Value.(*kv).val
	}
	return out
}

// Copy Completely copy a new Map
func (o *Map) Copy() *Map {
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
func (o *Map) Set(key, val any) {
	o.lock.Lock()
	defer o.lock.Unlock()

	el, ok := o.val[key]
	if ok {
		el.Value.(*kv).val = val
		return
	}
	o.val[key] = o.list.PushBack(&kv{key: key, val: val})
}

// SetByMap insert and overwrite update the kv value of another Map
func (o *Map) SetByMap(o2 *Map) {
	if o2 == nil {
		return
	}

	if o == o2 {
		return
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	for _, newKey := range o2.Keys() {
		newVal := o2.val[newKey].Value.(*kv).val
		if el, ok := o.val[newKey]; ok {
			el.Value.(*kv).val = newVal
			continue
		}
		o.val[newKey] = o.list.PushBack(&kv{key: newKey, val: newVal})
	}
}

// Get value by key
func (o *Map) Get(key any) (any, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	el, ok := o.val[key]
	if !ok {
		return nil, false
	}
	return el.Value.(*kv).val, true
}

// Has Check if a key value exists
func (o *Map) Has(key any) bool {
	o.lock.RLock()
	defer o.lock.RUnlock()

	_, ok := o.val[key]
	return ok
}

// GetByValue find out the key value based on the value
func (o *Map) GetByValue(val any) (key any, ok bool) {
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
func (o *Map) GetByIndex(idx int) (any, bool) {
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
func (o *Map) Remove(key any) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if el, ok := o.val[key]; ok {
		o.list.Remove(el)
		delete(o.val, key)
	}
}

// Clear Map
func (o *Map) Clear() {
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
func (o *Map) Keys() []any {
	o.lock.RLock()
	defer o.lock.RUnlock()

	out := make([]any, o.list.Len())
	for i, el := 0, o.list.Front(); el != nil; i++ {
		out[i] = el.Value.(*kv).key
		el = el.Next()
	}
	return out
}

// Len Map length
func (o *Map) Len() int { return o.list.Len() }

// Swap two key
func (o *Map) Swap(key1, key2 any) (err error) {
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
func (o *Map) Search(keys ...any) *Map {
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
func (o *Map) Filter(f func(key any, val any) bool) *Map {
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

func (o *Map) Marshal() ([]byte, error) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	items := make([]string, 0, o.Len())

	for i, el := 0, o.list.Front(); el != nil; i++ {
		data := el.Value.(*kv)
		b, err := json.Marshal(map[any]any{data.key: data.val})
		if err != nil {
			return nil, err
		}
		el = el.Next()
		items = append(items, string(b))
	}
	return []byte(fmt.Sprintf("{%s}", strings.Join(items, ","))), nil
}
