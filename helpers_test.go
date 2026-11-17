// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"strings"
	"testing"
)

// failureError is an error that describes itself as a failure reason, like an
// exception implementing the failure reason interface.
type failureError struct {
	synopsis    string
	description string
}

func newFailureError() *failureError {
	return &failureError{synopsis: "synopsis", description: "description"}
}

func (e *failureError) Error() string       { return e.synopsis }
func (e *failureError) Synopsis() string    { return e.synopsis }
func (e *failureError) Description() string { return e.description }

// wrappingError is an error with a cause.
type wrappingError struct{ cause error }

func (e *wrappingError) Error() string { return "wrapping" }
func (e *wrappingError) Unwrap() error { return e.cause }

// plainError and otherError are errors without a cause.
type plainError struct{}

func (plainError) Error() string { return "plain" }

type otherError struct{}

func (otherError) Error() string { return "other" }

func identity(v int) int { return v }

// assertPanics verifies that f panics with a message of this package.
func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		t.Helper()
		r := recover()
		if r == nil {
			t.Fatal("expected a panic")
		}
		if msg, ok := r.(string); !ok || !strings.HasPrefix(msg, "common: ") {
			t.Fatalf("unexpected panic value: %v", r)
		}
	}()
	f()
}
