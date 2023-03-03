/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/3/3 下午2:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

package ofunc

type TempsFunc interface {
	Call(args ...interface{}) (out []interface{})
	Update(fn interface{}) error
	Delete(fn interface{})
}
