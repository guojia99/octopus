package dp

import (
	"fmt"
	"reflect"
)

func NewAdapters() Adapters {
	return &adapters{srcFuncs: make(map[string]reflect.Value)}
}

type Adapters interface {
	Register(srcFunc interface{}, target string) error
	Walk(target string, args ...interface{}) (val []interface{}, err error)
}

type adapters struct {
	srcFuncs map[string]reflect.Value
}

func (c *adapters) Register(srcFunc interface{}, target string) error {
	funcType := reflect.TypeOf(srcFunc)
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("adapter register target `%s` to srcFunc kind will be func", target)
	}
	c.srcFuncs[target] = reflect.ValueOf(srcFunc)
	return nil
}

func (c *adapters) Walk(target string, args ...interface{}) (val []interface{}, err error) {
	defer func() {
		if result := recover(); result != nil {
			err = fmt.Errorf("walk recover error %+v", result)
		}
	}()
	if _, ok := c.srcFuncs[target]; !ok {
		err = fmt.Errorf("target `%s` not register", target)
		return
	}
	paramNum := c.srcFuncs[target].Type().NumIn()
	var paramList []reflect.Value
	for idx, arg := range args {
		if idx+1 > paramNum {
			break
		}
		paramList = append(paramList, reflect.ValueOf(arg))
	}
	retList := c.srcFuncs[target].Call(paramList)
	for _, ret := range retList {
		val = append(val, ret.Interface())
	}
	return
}
