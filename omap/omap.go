/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/3/3 下午2:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

/*
	this file is use go 1.18 generics
*/
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

type okv[k comparable, v any] struct {
	key k
	val v
}

// OMap with read-write locks, sequential keys
type OMap[k comparable, v any] struct {
	lock sync.RWMutex
	list *list.List
	// [WARN] If the idea reminds, please ignore the error
	// `Invalid OMap key type: comparison operators == and != must be fully defined for the key type`
	// your can look https://youtrack.jetbrains.com/issue/GO-12615
	val map[k]*list.Element
}

// NewOMap Create a new OMap
func NewOMap[k comparable, v any]() *OMap[k, v] {
	return &OMap[k, v]{
		lock: sync.RWMutex{},
		list: list.New(),
		val:  make(map[k]*list.Element),
	}
}

// Val get all metadata
func (o *OMap[k, v]) Val() (out map[k]v) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	out = make(map[k]v)
	for key := range o.val {
		out[key] = o.val[key].Value.(*okv[k, v]).val
	}
	return out
}

// Copy Completely copy a new OMap
func (o *OMap[k, v]) Copy() (out *OMap[k, v]) {
	o.lock.Lock()
	defer o.lock.Unlock()

	out = NewOMap[k, v]()
	for el := o.list.Front(); el != nil; el = el.Next() {
		data := el.Value.(*okv[k, v])
		out.Set(data.key, data.val)
	}
	return
}

// Set the value according to key and val
func (o *OMap[k, v]) Set(key k, val v) {
	o.lock.Lock()
	defer o.lock.Unlock()

	el, ok := o.val[key]
	if ok {
		el.Value.(*okv[k, v]).val = val
		return
	}
	o.val[key] = o.list.PushBack(&okv[k, v]{key: key, val: val})
}

// SetByOMap insert and overwrite update the ookv value of another OMap
func (o *OMap[k, v]) SetByOMap(o2 *OMap[k, v]) {
	if o2 == nil {
		return
	}
	if o == o2 {
		return
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	for _, newKey := range o2.Keys() {
		newVal := o2.val[newKey].Value.(*okv[k, v]).val
		if el, ok := o.val[newKey]; ok {
			el.Value.(*okv[k, v]).val = newVal
			continue
		}
		o.val[newKey] = o.list.PushBack(&okv[k, v]{key: newKey, val: newVal})
	}
}

// Get value by key
func (o *OMap[k, v]) Get(key k) (val v, ok bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	el, ok := o.val[key]
	if !ok {
		return
	}
	return el.Value.(*okv[k, v]).val, true
}

// Has Check if a key value exists
func (o *OMap[k, v]) Has(key k) bool {
	o.lock.RLock()
	defer o.lock.RUnlock()

	_, ok := o.val[key]
	return ok
}

// GetByValue find out the key value based on the value
func (o *OMap[k, v]) GetByValue(val v) (key k, ok bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	for el := o.list.Front(); el != nil; el = el.Next() {
		data := el.Value.(*okv[k, v])
		if reflect.DeepEqual(data.val, val) {
			return data.key, true
		}
	}
	return
}

// GetByIndex the key according to the first value of the OMap
func (o *OMap[k, v]) GetByIndex(idx int) (val v, ok bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	if o.Len() < int(math.Abs(float64(idx))) {
		return
	}

	el := o.list.Front()
	switch {
	case idx > 0:
		for i := 0; i < idx; i++ {
			el = el.Next()
		}
		return
	case idx < 0:
		el = o.list.Back()
		for i := 0; i > idx; i-- {
			el = el.Prev()
		}
	}
	return el.Value.(*okv[k, v]).val, true
}

// Remove key value
func (o *OMap[k, v]) Remove(key k) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if el, ok := o.val[key]; ok {
		o.list.Remove(el)
		delete(o.val, key)
	}
}

// Clear OMap
func (o *OMap[k, v]) Clear() {
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
func (o *OMap[k, v]) Keys() []k {
	o.lock.RLock()
	defer o.lock.RUnlock()

	out := make([]k, o.list.Len())
	for i, el := 0, o.list.Front(); el != nil; i++ {
		out[i] = el.Value.(*okv[k, v]).key
		el = el.Next()
	}
	return out
}

// Len OMap length
func (o *OMap[k, v]) Len() int { return o.list.Len() }

// Swap two key
func (o *OMap[k, v]) Swap(key1, key2 k) (err error) {
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
func (o *OMap[k, v]) Search(keys ...k) *OMap[k, v] {
	o.lock.RLock()
	defer o.lock.RUnlock()

	newOMap := NewOMap[k, v]()
	for _, key := range keys {
		if el, ok := o.val[key]; ok {
			newOMap.Set(key, el.Value.(*okv[k, v]).val)
		}
	}
	return newOMap
}

// Filter a collection that meets the requirements
func (o *OMap[k, v]) Filter(f func(key k, val v) bool) *OMap[k, v] {
	o.lock.RLock()
	defer o.lock.RUnlock()

	newOMap := NewOMap[k, v]()
	for _, key := range o.Keys() {
		val := o.val[key].Value.(*okv[k, v]).val
		if ok := f(key, val); ok {
			newOMap.Set(key, val)
		}
	}
	return newOMap
}

func (o *OMap[k, v]) Marshal() ([]byte, error) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	items := make([]string, 0, o.Len())

	for i, el := 0, o.list.Front(); el != nil; i++ {
		data := el.Value.(*okv[k, v])
		b, err := json.Marshal(map[any]any{data.key: data.val})
		if err != nil {
			return nil, err
		}
		el = el.Next()
		items = append(items, string(b))
	}
	return []byte(fmt.Sprintf("{%s}", strings.Join(items, ","))), nil
}
