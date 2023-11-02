package main

import (
	"errors"
	"fmt"

	"github.com/guojia99/octopus/ofunc"
)

func main() {

	//// The normal processing process is to only process one
	//ofunc.Try(func() {
	//	// todo code
	//	// Empty errors will not be handled, you can easily place your errors directly in Throw
	//	var err error
	//	ofunc.Throw(err)
	//}).Catch(func(err ofunc.StackTraceErr) {
	//	fmt.Println(err)
	//})
	//
	//// Multiple function chain execution will not interfere with each other Try
	//// When you have many Try, please use CatchAll to catch errors
	//ofunc.Try(func() {
	//	// todo code
	//	// Empty errors will not be handled, you can easily place your errors directly in Throw
	//	var err error
	//	ofunc.Throw(err)
	//}).Try(func() {
	//	// todo code
	//	panic(errors.New("test error by panic"))
	//}).Try(func() {
	//	// todo code
	//	// The Throw function will determine whether error is empty. If it is empty,
	//	// it will interrupt the subsequent execution of the function.
	//	ofunc.Throw("test error by throw")
	//	fmt.Println("this code will not be executed ")
	//}).Try(func() {
	//	// todo code
	//	a := 0
	//	b := 0
	//	fmt.Println(a / b)
	//}).CatchAll(func(errs []ofunc.StackTraceErr) {
	//	for i := 0; i < len(errs); i++ {
	//		fmt.Println(errs[i])
	//	}
	//})

	// If you only need to catch certain errors, you can use interceptors
	var DataError = errors.New("test data error")
	ofunc.Try(func() {
		ofunc.Throw(DataError)
	}).Try(func() {
		ofunc.Throw("todo error")
	}).CatchWithErr(DataError, func(err ofunc.StackTraceErr) {
		fmt.Println(err)
	})
}
