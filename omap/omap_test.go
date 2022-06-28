package omap

import (
	"fmt"
	"testing"
)

func TestNewMap(t *testing.T) {

	t.Run("int2intMap", func(t *testing.T) {
		t.Log(NewMap[int, int]())
	})

	t.Run("int2stringMap", func(t *testing.T) {
		t.Log(NewMap[int, string]())
	})

	t.Run("bool2stringMap", func(t *testing.T) {
		t.Log(NewMap[bool, string]())
	})

	t.Run("float2intMap", func(t *testing.T) {
		t.Log(NewMap[float64, int]())
	})
}

func TestMap_SetGet(t *testing.T) {
	c := NewMap[float64, []byte]()
	c.Set(1, []byte{1, 2, 3})
	fmt.Println(c.Get(1))
}

func TestMap_SetByMap(t *testing.T) {
	c := NewMap[float64, []byte]()
	c.Set(1, []byte{1, 2, 3})
	c2 := NewMap[float64, []byte]()
	c2.Set(2, []byte{4, 5, 6})
	c.SetByMap(c2)
	fmt.Println(c.Val())
}
