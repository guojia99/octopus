# Omap

```go
package main

import (
	"github.com/guojia99/octopus/omap"
)

func main() {
	data := omap.NewMap()
	data.Set("key1", "value")
	data.SetByOmap(nil)
	_ = data.Has("key1")
	data = data.Copy()
	data.GetByValue("value")
}
```

```go
func NewMap() *Omap
func (o *Omap) Clear()
func (o *Omap) Copy() *Omap
func (o *Omap) Filter(f func(key any, val any) bool) *Omap
func (o *Omap) Get(key any) (any, bool)
func (o *Omap) GetByIndex(idx int) (any, bool)
func (o *Omap) GetByValue(val any) (key any, ok bool)
func (o *Omap) Has(key any) bool
func (o *Omap) Keys() []any
func (o *Omap) Len() int
func (o *Omap) Remove(key any)
func (o *Omap) Search(keys ...any) *Omap
func (o *Omap) Set(key, val any)
func (o *Omap) SetByOmap(o2 *Omap)
func (o *Omap) Swap(key1, key2 any) (err error)
```

