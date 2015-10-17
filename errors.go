package modules

import (
	"bytes"
	"fmt"
)

// An AnnotatedError holds a message and wraps another error.
type AnnotatedError struct {
	msg   string
	cause error
}

func (e *AnnotatedError) Error() string {
	return e.msg + " caused by: " + e.cause.Error()
}

// A BindingError indicates failure during binding.
// Holds one or more errors which prevented binding.
type BindingError struct {
	errs []error
}

func (e *BindingError) Error() string {
	errMsg := bytes.NewBufferString("binding failed:")
	for _, err := range e.errs {
		fmt.Fprintf(errMsg, "\t%s\n", err.Error())
	}
	return errMsg.String()
}
