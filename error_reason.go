// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common

import "runtime"

// FailureError is an error that is also a [FailureReason].
type FailureError interface {
	error
	FailureReason
}

// ErrorReason adapts an error to a [FailureError]: its Synopsis is the message
// of the error.
//
// An error type of your own becomes a FailureError by implementing Synopsis
// itself; ErrorReason is for errors you do not control.
type ErrorReason struct {
	err   error
	stack stack
}

var (
	_ FailureError = (*ErrorReason)(nil)
	_ StackTracer  = (*ErrorReason)(nil)
)

// NewErrorReason creates an ErrorReason for the given error. It panics if err
// is nil.
//
// Go errors do not carry a stack trace, so the stack of the call to
// NewErrorReason is recorded instead.
func NewErrorReason(err error) *ErrorReason {
	requireNonNil(err == nil, "error")
	return &ErrorReason{err: err, stack: callers()}
}

// Error returns the message of the wrapped error.
func (e *ErrorReason) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error.
func (e *ErrorReason) Unwrap() error {
	return e.err
}

// Synopsis returns the message of the wrapped error.
func (e *ErrorReason) Synopsis() string {
	return e.err.Error()
}

// StackTrace returns the stack of the call that created e.
func (e *ErrorReason) StackTrace() []runtime.Frame {
	return e.stack.frames()
}
