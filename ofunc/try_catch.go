package ofunc

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

func Try(f func()) TryCatch { // skip 0
	c := &tryCatch{
		stack: make([]StackTraceErr, 0),
	}
	c.Try(f)
	return c
}

func Throw(err any) {
	if err != nil {
		panic(err)
	}
}

type (
	TryFn         func()
	CatchErrFn    func(err StackTraceErr)
	CatchAllErrFn func(errs []StackTraceErr)
	TryCatch      interface {
		Try(TryFn) TryCatch                             // 执行一次尝试操作
		Catch(CatchErrFn) TryCatch                      // 取出一个错误进行处理
		CatchAll(CatchAllErrFn) TryCatch                // 取出所有错误进行处理
		CatchWithErr(err error, fn CatchErrFn) TryCatch // 取出指定的错误进行处理
		Finally(TryFn) TryCatch                         // 立即执行一次操作
	}
)

type tryCatch struct {
	lock  sync.Mutex
	stack []StackTraceErr
}

const (
	IgnoredErrors = "github.com/guojia99/octopus/ofunc"
)

func (t *tryCatch) Try(fn TryFn) TryCatch {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.stack == nil {
		t.stack = make([]StackTraceErr, 0)
	}
	func() {
		defer func() {
			// get recover massage
			result := recover()
			if result == nil {
				return
			}

			// get trace
			var ts = StackTraceErr{
				Err:         fmt.Errorf("%+v", result),
				StackTraces: make([]StackTrace, 0),
			}
			for i := 2; ; i++ {
				pc, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				st := StackTrace{
					pc:   pc,
					file: file,
					line: line,
				}
				f := runtime.FuncForPC(st.pc)
				if strings.Contains(f.Name(), IgnoredErrors) {
					continue
				}
				ts.StackTraces = append(ts.StackTraces, st)
			}
			t.stack = append(t.stack, ts)
		}()
		fn() // do run function, wait error into stack
	}()

	return t
}

func (t *tryCatch) Catch(fn CatchErrFn) TryCatch {
	t.lock.Lock()
	defer t.lock.Unlock()

	if len(t.stack) == 0 {
		return t
	}

	fn(t.stack[0])
	t.stack = t.stack[:1]
	return t
}

func (t *tryCatch) CatchAll(fn CatchAllErrFn) TryCatch {
	t.lock.Lock()
	defer t.lock.Unlock()

	if len(t.stack) == 0 {
		return t
	}
	fn(t.stack)
	t.stack = make([]StackTraceErr, 0)
	return t
}

func (t *tryCatch) CatchWithErr(err error, fn CatchErrFn) TryCatch {
	t.lock.Lock()
	defer t.lock.Unlock()
	if err == nil {
		return t
	}
	for i := 0; i < len(t.stack); i++ {
		if errors.Is(err, t.stack[i].Err) ||
			strings.Contains(err.Error(), t.stack[i].Err.Error()) ||
			err.Error() == t.stack[i].Err.Error() {
			fn(t.stack[i])
			t.stack = append(t.stack[:i], t.stack[i+1:]...)
		}
	}
	return t
}

func (t *tryCatch) Finally(fn TryFn) TryCatch {
	t.lock.Lock()
	defer t.lock.Unlock()
	fn()
	return t
}
