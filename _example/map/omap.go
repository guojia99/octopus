package main

import (
	"bytes"
	"fmt"

	"github.com/guojia99/octopus/omap"
)

func main() {
	m := omap.NewOMap[float64, []byte]()
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

	buff, _ := m.Marshal()
	fmt.Println("json ->", string(buff))

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

	o3 := m.Filter(func(key float64, val []byte) bool {
		if key == 1 {
			return true
		}
		if bytes.Equal(val, []byte{3}) {
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
