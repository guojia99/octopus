package ofunc

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const (
	IgnoredErrors = "github.com/guojia99/octopus/ofunc"
)

func Try(f func()) TryCatch { // skip 0
	c := &tryCatch{
		stack: make(StackTraces, 0),
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
	TryFn        func()
	CatchErrFn   func(err error)
	CatchTraceFn func(err error, stack StackTraces)

	TryCatch interface {
		Try(TryFn) TryCatch
		Catch(CatchErrFn) TryCatch
		CatchTrace(CatchTraceFn) TryCatch
		Finally(TryFn) TryCatch
	}
)

type tryCatch struct {
	err   error
	stack StackTraces
}

func (t *tryCatch) Try(fn TryFn) TryCatch {
	if t.stack == nil {
		t.stack = make(StackTraces, 0)
	}
	func() {
		defer func() {
			// get recover massage
			result := recover()
			if result == nil {
				return
			}

			// update error
			if t.err != nil {
				t.err = errors.Join(t.err, fmt.Errorf("%+v", result))
			} else {
				t.err = fmt.Errorf("%+v", result)
			}

			// get trace
			for i := 2; ; i++ {
				pc, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				st := StackTrace{
					idx:  i,
					pc:   pc,
					file: file,
					line: line,
				}

				if strings.Contains(st.String(), IgnoredErrors) {
					continue
				}
				t.stack = append(t.stack, st)
			}
		}()
		fn() // do run function, wait error into stack
	}()

	return t
}

func (t *tryCatch) Catch(fn CatchErrFn) TryCatch {
	if t.err == nil {
		return t
	}
	fn(t.err)
	return t
}

func (t *tryCatch) CatchTrace(fn CatchTraceFn) TryCatch {
	if t.err == nil {
		return t
	}
	fn(t.err, t.stack)
	return t
}

func (t *tryCatch) Finally(fn TryFn) TryCatch { fn(); return t }

type (
	StackTrace struct {
		idx  int
		pc   uintptr
		file string
		line int
	}
	StackTraces []StackTrace
)

func (s *StackTrace) String() string {
	f := runtime.FuncForPC(s.pc)
	return fmt.Sprintf("%s\n\t%s:%d\n", f.Name(), s.file, s.line)
}

func (s StackTraces) String() string {
	out := ""
	for _, ss := range s {
		out += ss.String()
	}
	return out
}
