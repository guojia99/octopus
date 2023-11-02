package main

import (
	"errors"
	"fmt"

	"github.com/guojia99/octopus/ofunc"
)

func main() {
	ofunc.Try(func() {
		// todo code
		var err = errors.New("test error by panic")
		if err != nil {
			panic(err)
		}
	}).Try(func() {
		ofunc.Throw("test error by throw")
	}).Catch(func(err error) {
		// todo code by error
		fmt.Println(err)
	}).CatchTrace(func(err error, stack ofunc.StackTraces) {
		// todo code by detailed error
		fmt.Println(stack.String())
	}).Finally(func() {
		fmt.Println("finally")
	}).Finally(func() {
		fmt.Println("finally2")
	})
}
