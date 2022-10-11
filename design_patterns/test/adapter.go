package main

import (
	"fmt"

	dp "github.com/guojia99/octopus/design_patterns"
)

type oldClient1 struct{}

func (c *oldClient1) Print(src interface{}) { fmt.Println("old1", src) }

type oldClient2 struct{}

func (c *oldClient2) PrintMulti(src1, src2, src3 interface{}) {
	fmt.Println(append([]interface{}{"old2"}, src1, src2, src3))
}

func main() {
	old1Cli, old2Cli := &oldClient1{}, &oldClient2{}
	adaptersCli := dp.NewAdapters()
	if err := adaptersCli.Register(old1Cli.Print, "old1"); err != nil {
		panic(err)
	}
	if err := adaptersCli.Register(old2Cli.PrintMulti, "old2"); err != nil {
		panic(err)
	}
	// register different functions through different targets to access different functions
	_, err := adaptersCli.Walk("old1", 1, 2, 3)
	if err != nil {
		panic(err)
	}
	_, err = adaptersCli.Walk("old2", 4, 5, 6)
	if err != nil {
		panic(err)
	}
}
