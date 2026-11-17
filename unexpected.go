// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common

import "runtime"

// Unexpected is a [FailureReason] for an error that was not anticipated.
type Unexpected struct {
	err   error
	stack stack
}

var (
	_ FailureReason = (*Unexpected)(nil)
	_ Describable   = (*Unexpected)(nil)
	_ StackTracer   = (*Unexpected)(nil)
	_ IssueReporter = (*Unexpected)(nil)
)

// NewUnexpected creates an Unexpected for the given error. It panics if err is
// nil.
//
// Go errors do not carry a stack trace, so the stack of the call to
// NewUnexpected is recorded instead.
func NewUnexpected(err error) *Unexpected {
	requireNonNil(err == nil, "error")
	return &Unexpected{err: err, stack: callers()}
}

// Synopsis returns a fixed summary.
func (u *Unexpected) Synopsis() string {
	return "An unknown exception has occurred."
}

// Description returns the message of the error.
func (u *Unexpected) Description() string {
	return u.err.Error()
}

// StackTrace returns the stack of the call that created u.
func (u *Unexpected) StackTrace() []runtime.Frame {
	return u.stack.frames()
}

// ReportIssue returns true: the user should be asked to report unexpected
// failures.
func (u *Unexpected) ReportIssue() bool {
	return true
}

// Err returns the error u was created for.
func (u *Unexpected) Err() error {
	return u.err
}
