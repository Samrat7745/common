// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common

import "runtime"

// Describable is implemented by values that have a description and a short
// summary of it.
type Describable interface {
	// Synopsis returns a short summary of the description.
	Synopsis() string
	// Description returns the full description.
	Description() string
}

// FailureReason is the base type of the Carbyne Stack failure reasons.
//
// Only Synopsis is required. A reason may additionally implement
// [Describable], [StackTracer] and [IssueReporter]; use [Description],
// [StackTrace] and [ReportIssue] to read these properties with their defaults
// for reasons that do not.
type FailureReason interface {
	// Synopsis returns a short summary of the failure reason.
	Synopsis() string
}

// StackTracer is implemented by failure reasons that carry a stack trace.
type StackTracer interface {
	// StackTrace returns the stack trace of the failure, innermost frame
	// first.
	StackTrace() []runtime.Frame
}

// IssueReporter is implemented by failure reasons that decide whether the
// application should prompt the user to report the issue.
type IssueReporter interface {
	// ReportIssue reports whether the user should be prompted to report the
	// issue.
	ReportIssue() bool
}

// Description returns the description of the reason. It is the reason's
// Description if the reason implements [Describable], and its Synopsis
// otherwise.
func Description(reason FailureReason) string {
	if d, ok := reason.(Describable); ok {
		return d.Description()
	}
	return reason.Synopsis()
}

// StackTrace returns the stack trace of the reason. The boolean is false if
// the reason does not implement [StackTracer] or has no stack trace.
func StackTrace(reason FailureReason) ([]runtime.Frame, bool) {
	if s, ok := reason.(StackTracer); ok {
		if frames := s.StackTrace(); frames != nil {
			return frames, true
		}
	}
	return nil, false
}

// ReportIssue reports whether the failure should make the application prompt
// the user to report the issue. It is false unless the reason implements
// [IssueReporter] and says otherwise.
func ReportIssue(reason FailureReason) bool {
	if r, ok := reason.(IssueReporter); ok {
		return r.ReportIssue()
	}
	return false
}

// ToFailure converts the reason into a failed [Result] with the success type
// S.
func ToFailure[S any, F FailureReason](reason F) Result[S, F] {
	return Failure[S, F](reason)
}

// stack is a captured call stack.
type stack []uintptr

// callers captures the stack of the caller of the function that calls it.
func callers() stack {
	pcs := make([]uintptr, 64)
	return pcs[:runtime.Callers(3, pcs)]
}

// frames resolves the captured program counters.
func (s stack) frames() []runtime.Frame {
	if len(s) == 0 {
		return nil
	}
	frames := make([]runtime.Frame, 0, len(s))
	iter := runtime.CallersFrames(s)
	for {
		frame, more := iter.Next()
		frames = append(frames, frame)
		if !more {
			return frames
		}
	}
}

// requireNonNil panics if isNil is true. It is the counterpart of a
// NullPointerException for a missing function argument.
func requireNonNil(isNil bool, name string) {
	if isNil {
		panic("common: " + name + " must not be nil")
	}
}
