package main

import (
	"fmt"

	"github.com/guojia99/octopus/omap"
)

func main() {
	m := omap.NewMap()
	for i := 0; i < 5; i++ {
		m.Set(float64(i), []byte{
			uint8(i),
		})
	}
	fmt.Println("key -->", m.Keys())
	fmt.Println("val -->", m.Val())
	fmt.Println("len -->", m.Len())
	m.Remove(float64(4))
	m.Swap(float64(1), float64(2))
	fmt.Println("key -->", m.Keys())
	fmt.Println("val -->", m.Val())
	fmt.Println("len -->", m.Len())

	getData, _ := m.Get(float64(1))
	fmt.Println("get -->", getData)
	hasData := m.Has(float64(2))
	fmt.Println("has -->", hasData)

	val := []byte{1}
	byValueData, _ := m.GetByValue(val)
	fmt.Println("get by value ->", byValueData)

	byIndexData, _ := m.GetByIndex(1)
	fmt.Println("get by index ->", byIndexData)

	o2 := m.Search(float64(1), float64(2))
	fmt.Println("search ->", o2.Val())

	o3 := m.Filter(func(key any, val any) bool {
		if key.(float64) == 1 {
			return true
		}
		return false
	})
	fmt.Println("filter ->", o3.Val())

	m.Clear()
	fmt.Println("key -->", m.Keys())
	fmt.Println("val -->", m.Val())
	fmt.Println("len -->", m.Len())
}
