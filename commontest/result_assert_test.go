// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package commontest_test

import (
	"fmt"
	"testing"

	"github.com/Samrat7745/common"
	"github.com/Samrat7745/common/commontest"
)

// recorder is a testing.TB that records fatal failures instead of ending the
// test. Any other reporting method panics, as the embedded TB is nil.
type recorder struct {
	testing.TB
	failures []string
	helpers  int
}

func (r *recorder) Helper() { r.helpers++ }

func (r *recorder) Fatalf(format string, args ...any) {
	r.failures = append(r.failures, fmt.Sprintf(format, args...))
}

// expectPass verifies that the assertion reports no failure.
func expectPass(t *testing.T, assertion func(tb testing.TB)) {
	t.Helper()
	rec := &recorder{}
	assertion(rec)
	if len(rec.failures) != 0 {
		t.Errorf("unexpected failures: %q", rec.failures)
	}
	if rec.helpers == 0 {
		t.Error("the assertion did not mark itself as a helper")
	}
}

// expectFailure verifies that the assertion reports exactly the given failure.
func expectFailure(t *testing.T, want string, assertion func(tb testing.TB)) {
	t.Helper()
	rec := &recorder{}
	assertion(rec)
	if len(rec.failures) != 1 || rec.failures[0] != want {
		t.Errorf("failures = %q, want [%q]", rec.failures, want)
	}
}

const (
	value  = 12
	reason = 21
	some   = "some"
	none   = "none"
)

var (
	success = common.Success[int, int](value)
	failure = common.Failure[int, int](reason)
)

func stringSuccess(s string) common.Result[string, string] {
	return common.Success[string, string](s)
}

func stringFailure(s string) common.Result[string, string] {
	return common.Failure[string, string](s)
}

func TestResultAssert_IsSuccess(t *testing.T) {
	commontest.AssertThat(t, success).IsSuccess()
	expectPass(t, func(tb testing.TB) { commontest.AssertThat(tb, success).IsSuccess() })
	expectFailure(t, "Expecting a success result but was: Failure[reason=21]", func(tb testing.TB) {
		commontest.AssertThat(tb, failure).IsSuccess()
	})
}

func TestResultAssert_IsFailure(t *testing.T) {
	commontest.AssertThat(t, failure).IsFailure()
	expectPass(t, func(tb testing.TB) { commontest.AssertThat(tb, failure).IsFailure() })
	expectFailure(t, "Expecting a failure result but was: Success[value=12]", func(tb testing.TB) {
		commontest.AssertThat(tb, success).IsFailure()
	})
}

func TestResultAssert_HasValue(t *testing.T) {
	commontest.AssertThat(t, success).HasValue(value)
	expectPass(t, func(tb testing.TB) { commontest.AssertThat(tb, success).HasValue(value) })
	expectFailure(t, "Expecting a success result but was: Failure[reason=21]", func(tb testing.TB) {
		commontest.AssertThat(tb, failure).HasValue(value)
	})
	expectFailure(t, "Expecting result success value to equal 21 but was: 12", func(tb testing.TB) {
		commontest.AssertThat(tb, success).HasValue(reason)
	})
}

func TestResultAssert_HasStringValue(t *testing.T) {
	expectFailure(t, "Expecting result success value to equal 'none' but was: 'some'", func(tb testing.TB) {
		commontest.AssertThat(tb, stringSuccess(some)).HasValue(none)
	})
}

func TestResultAssert_HasStringValueWithSingleQuotes(t *testing.T) {
	expectFailure(t, "Expecting result success value to equal 'second' but was: ''first''", func(tb testing.TB) {
		commontest.AssertThat(tb, stringSuccess("'first'")).HasValue("second")
	})
}

func TestResultAssert_HasStringValueWithLineBreaks(t *testing.T) {
	expectFailure(t, `Expecting result success value to equal 'third' but was: 'first\nsecond'`, func(tb testing.TB) {
		commontest.AssertThat(tb, stringSuccess("first\nsecond")).HasValue("third")
	})
}

func TestResultAssert_HasReason(t *testing.T) {
	commontest.AssertThat(t, failure).HasReason(reason)
	expectPass(t, func(tb testing.TB) { commontest.AssertThat(tb, failure).HasReason(reason) })
	expectFailure(t, "Expecting a failure result but was: Success[value=12]", func(tb testing.TB) {
		commontest.AssertThat(tb, success).HasReason(reason)
	})
	expectFailure(t, "Expecting result failure reason to equal 12 but was: 21", func(tb testing.TB) {
		commontest.AssertThat(tb, failure).HasReason(value)
	})
}

func TestResultAssert_HasStringReason(t *testing.T) {
	expectFailure(t, "Expecting result failure reason to equal 'none' but was: 'some'", func(tb testing.TB) {
		commontest.AssertThat(tb, stringFailure(some)).HasReason(none)
	})
}

func TestResultAssert_HasStringReasonWithSingleQuotes(t *testing.T) {
	expectFailure(t, "Expecting result failure reason to equal 'second' but was: ''first''", func(tb testing.TB) {
		commontest.AssertThat(tb, stringFailure("'first'")).HasReason("second")
	})
}

func TestResultAssert_HasStringReasonWithLineBreaks(t *testing.T) {
	expectFailure(t, `Expecting result failure reason to equal 'third' but was: 'first\nsecond'`, func(tb testing.TB) {
		commontest.AssertThat(tb, stringFailure("first\nsecond")).HasReason("third")
	})
}

func TestResultAssert_Chaining(t *testing.T) {
	commontest.AssertThat(t, success).IsSuccess().HasValue(value)
	commontest.AssertThat(t, failure).IsFailure().HasReason(reason)
}

func TestResultAssert_ComparesDeeply(t *testing.T) {
	res := common.Success[[]int, string]([]int{1, 2})
	commontest.AssertThat(t, res).HasValue([]int{1, 2})
	expectFailure(t, "Expecting result success value to equal [1 3] but was: [1 2]", func(tb testing.TB) {
		commontest.AssertThat(tb, res).HasValue([]int{1, 3})
	})
}
