package ofunc

import (
	"fmt"
	"runtime"
)

type (
	StackTrace struct {
		pc   uintptr
		file string
		line int
	}
	StackTraceErr struct {
		StackTraces []StackTrace
		Err         error
	}
)

func (s StackTraceErr) String() string {
	out := fmt.Sprintf("%+v error stack traces\n", s.Err)
	for _, ss := range s.StackTraces {
		f := runtime.FuncForPC(ss.pc)
		out += fmt.Sprintf("%s\n\t%s:%d\n", f.Name(), ss.file, ss.line)
	}
	return out
}
