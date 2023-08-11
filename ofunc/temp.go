/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/3/3 下午2:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

package ofunc

import (
	"errors"
	"reflect"
)

type TempsFunc interface {
	Call(args ...interface{}) (out []interface{})
	Update(fn interface{}) error
	Delete(fn interface{})
}

type tempFunc struct {
	callMap map[string]reflect.Value
}

func (t *tempFunc) checkFn(fn interface{}) (string, reflect.Value, error) {
	val := reflect.ValueOf(fn)
	typ := reflect.TypeOf(fn)

	if typ.Kind() != reflect.Func {
		return "", reflect.Value{}, errors.New("input fn value is not func")
	}
	typ.
}
