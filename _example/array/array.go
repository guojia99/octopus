package main

import (
	"fmt"

	"github.com/guojia99/octopus/oarray"
)

func main() {
	a := oarray.NewArray()

	a.Set("1", "2", "3", "4")
	fmt.Println(a.Len(), a.IsEmpty(), a.RandomOne())
	fmt.Println(a.Val())
	a.Reverse()
	fmt.Println(a.Val())

	a2 := a.Copy()
	fmt.Println(a2.Val())
	fmt.Println(a.GetByIndex(1))
	fmt.Println(a.GetSlice(1, 2))

	a2.Remove("1")
	fmt.Println(a2.Val())
	fmt.Println(a.Filter(func(value any) bool {
		if value == "1" || value == "2" {
			return true
		}
		return false
	}).Val())

	fmt.Println(a.HasValueAll("1"))
}
