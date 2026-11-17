// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package commontest

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Samrat7745/common"
)

// ResultAssert holds assertions on a [common.Result]. Create one with
// [AssertThat].
//
// A failed assertion ends the test through [testing.TB.Fatalf]. Like Fatalf,
// the assertions must be called from the goroutine running the test.
type ResultAssert[S, F any] struct {
	t      testing.TB
	actual common.Result[S, F]
}

// AssertThat returns assertions on the result that report to t.
func AssertThat[S, F any](t testing.TB, actual common.Result[S, F]) ResultAssert[S, F] {
	return ResultAssert[S, F]{t: t, actual: actual}
}

// IsSuccess verifies that the result is a success.
func (a ResultAssert[S, F]) IsSuccess() ResultAssert[S, F] {
	a.t.Helper()
	a.isSuccess()
	return a
}

// IsFailure verifies that the result is a failure.
func (a ResultAssert[S, F]) IsFailure() ResultAssert[S, F] {
	a.t.Helper()
	a.isFailure()
	return a
}

// HasValue verifies that the result is a success and that its value is equal
// to the given one. Values are compared with [reflect.DeepEqual].
func (a ResultAssert[S, F]) HasValue(value S) ResultAssert[S, F] {
	a.t.Helper()
	if !a.isSuccess() {
		return a
	}
	if actual, _ := a.actual.Value(); !reflect.DeepEqual(actual, value) {
		a.t.Fatalf("Expecting result success value to equal %s but was: %s", display(value), display(actual))
	}
	return a
}

// HasReason verifies that the result is a failure and that its reason is equal
// to the given one. Reasons are compared with [reflect.DeepEqual].
func (a ResultAssert[S, F]) HasReason(reason F) ResultAssert[S, F] {
	a.t.Helper()
	if !a.isFailure() {
		return a
	}
	if actual, _ := a.actual.Reason(); !reflect.DeepEqual(actual, reason) {
		a.t.Fatalf("Expecting result failure reason to equal %s but was: %s", display(reason), display(actual))
	}
	return a
}

func (a ResultAssert[S, F]) isSuccess() bool {
	a.t.Helper()
	if a.actual.IsFailure() {
		a.t.Fatalf("Expecting a success result but was: %v", a.actual)
		return false
	}
	return true
}

func (a ResultAssert[S, F]) isFailure() bool {
	a.t.Helper()
	if a.actual.IsSuccess() {
		a.t.Fatalf("Expecting a failure result but was: %v", a.actual)
		return false
	}
	return true
}

// display formats a value for a failure message. Strings are put in single
// quotes and their line breaks are escaped to keep the message on one line.
func display(v any) string {
	if s, ok := v.(string); ok {
		return "'" + strings.ReplaceAll(s, "\n", `\n`) + "'"
	}
	return fmt.Sprint(v)
}
