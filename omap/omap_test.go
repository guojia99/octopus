package omap

import (
	"fmt"
	"testing"
)

func TestNewOMap(t *testing.T) {

	t.Run("int2intMap", func(t *testing.T) {
		t.Log(NewOMap[int, int]())
	})

	t.Run("int2stringMap", func(t *testing.T) {
		t.Log(NewOMap[int, string]())
	})

	t.Run("bool2stringMap", func(t *testing.T) {
		t.Log(NewOMap[bool, string]())
	})

	t.Run("float2intMap", func(t *testing.T) {
		t.Log(NewOMap[float64, int]())
	})
}

func TestMap_SetGet(t *testing.T) {
	c := NewOMap[float64, []byte]()
	c.Set(1, []byte{1, 2, 3})
	fmt.Println(c.Get(1))
}

func TestMap_SetByMap(t *testing.T) {
	c := NewOMap[float64, []byte]()
	c.Set(1, []byte{1, 2, 3})
	c2 := NewOMap[float64, []byte]()
	c2.Set(2, []byte{4, 5, 6})
	c.SetByOMap(c2)
	fmt.Println(c.Val())
}
