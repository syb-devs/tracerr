package tracerr

import (
	"fmt"
	"runtime"
)

type stackFrame struct {
	pc   uintptr
	file string
	line int
}

type stackTrace []stackFrame

func newStackTrace(skip, limit int) stackTrace {
	var st stackTrace
	if skip < 0 {
		return st
	}
	var n int
	for {
		if limit > 0 && n == limit {
			return st
		}
		pc, file, line, b := runtime.Caller(skip + 1)
		if !b {
			return st
		}
		st = append(st, stackFrame{
			pc:   pc,
			file: file,
			line: line,
		})
		skip++
		n++
	}
}

// String returns a string representation of the stackTrace
func (st stackTrace) String() string {
	var str string
	for _, sf := range st {
		str += fmt.Sprintf("%s:%d\n", sf.file, sf.line)
	}
	return str
}

// WrappedError is used to combine multiple errors into one containing stack trace information
type WrappedError struct {
	msg   string
	err   error
	trace stackTrace
}

func (we *WrappedError) Error() string {
	return we.msg
}

// TraceString returns a string representing the stack trace of the error
func (we *WrappedError) TraceString() string {
	if we2, ok := we.err.(*WrappedError); ok {
		return fmt.Sprintf(
			"wrapper error: %s\n%s\n%s\n",
			we.msg,
			we.trace.String(),
			we2.TraceString(),
		)
	}
	return fmt.Sprintf("wrapper error: %s\n%s\ntrigger error: %s\n", we.msg, we.trace.String(), we.err.Error())
}

// WrapError returns a new error that wraps the given one adding stack trace info to it
func WrapError(err error, message string) *WrappedError {
	we := &WrappedError{
		err: err,
		msg: message,
	}
	var limit int
	if _, ok := err.(*WrappedError); ok {
		limit = 1
	} else {
		limit = 0
	}
	we.trace = newStackTrace(1, limit)
	return we
}
