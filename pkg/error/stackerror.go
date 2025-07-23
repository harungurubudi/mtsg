package error

import (
	"fmt"
	"runtime"
	"strings"
)

type StackError struct {
	Msg      string    // Error message
	Cause    error     // Wrapped cause error (optional)
	Stack    []uintptr // Captured stack trace
	Function string    // Caller function name
	File     string    // Caller file
	Line     int       // Caller line
	Code     string    // (Optional) Error code for structured handling
}

// Error implements the error interface.
func (e *StackError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

// Unwrap returns the wrapped cause error.
func (e *StackError) Unwrap() error {
	return e.Cause
}

// StackTrace returns a formatted stack trace as a string.
func (e *StackError) StackTrace() string {
	var b strings.Builder
	frames := runtime.CallersFrames(e.Stack)
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&b, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return b.String()
}

// NewStackError creates a new StackError with stack trace and caller info.
func NewStackError(msg string, cause error) *StackError {
	const depth = 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(2, pcs)
	frame, _ := runtime.CallersFrames(pcs[:n]).Next()
	return &StackError{
		Msg:      msg,
		Cause:    cause,
		Stack:    pcs[:n],
		Function: frame.Function,
		File:     frame.File,
		Line:     frame.Line,
	}
}

// NewStackErrorf creates a new StackError with formatted message.
func NewStackErrorf(format string, args ...interface{}) *StackError {
	return NewStackError(fmt.Sprintf(format, args...), nil)
}
