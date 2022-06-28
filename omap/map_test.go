package omap

import (
	"testing"
)

func TestOmap_Clear(t *testing.T) {
	omap := NewMap()

	omap.Set("key1", "value1")
	omap.SetByMap(omap)

	omap.Clear()
	if omap.list.Len() > 0 {
		t.Errorf("clear error")
	}
}

func TestOmap_Copy(t *testing.T) {
	omap := NewMap()
	omap.Set("key1", "value1")
	omap2 := omap.Copy()

	if _, ok := omap2.Get("key1"); !ok {
		t.Errorf("copy error")
	}
}

func TestOmap_Filter(t *testing.T) {
	omap := NewMap()

	omap.Set("wantKey", "wantValue")
	omap.Set("notKey", "notValue")

	omap2 := omap.Filter(func(key any, val any) bool {
		return key == "wantKey"
	})

	if !omap2.Has("wantKey") {
		t.Errorf("filter error not appear wantKey ")
	}
	if omap2.Has("notKey") {
		t.Errorf("filter error appear notKey")
	}
}

func TestOmap_GetByValue(t *testing.T) {
	omap := NewMap()

	omap.Set("wantKey", "wantValue")

	if key, ok := omap.GetByValue("wantValue"); !ok || key != "wantKey" {
		t.Errorf("get by value error")
	}
}
