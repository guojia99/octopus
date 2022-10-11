package main

import (
	"context"
	"fmt"

	dp "github.com/guojia99/octopus/design_patterns"
)

func main() {
	intIter := dp.NewIterator([]interface{}{1, 2, 3, 4, 5}...)
	var v interface{}

	// 1. ordinary
	intIter.Reset()
	for intIter.HasNext() {
		v = intIter.Next()
		fmt.Println("ordinary", v)
	}

	// 2. chancel
	intIter.Reset()
	ctx, cancel := context.WithCancel(context.Background())
	ch := intIter.IteratorChan(ctx)
	for {
		if !intIter.HasNext() {
			cancel()
			break
		}
		v = <-ch
		fmt.Println("chan", v)
	}

	// 3. do function
	intIter.Reset()
	intIter.Do(func(idx int, val interface{}) bool {
		fmt.Println("do", idx, val)
		return true
	})
}
