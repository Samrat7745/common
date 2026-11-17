// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"io"
	"slices"
	"testing"

	"github.com/Samrat7745/common"
	"github.com/Samrat7745/common/commontest"
)

const failureReason = 21

var failure = common.Failure[int, int](failureReason)

func TestFailure_IsFailure(t *testing.T) {
	if !failure.IsFailure() {
		t.Error("IsFailure() = false, want true")
	}
	commontest.AssertThat(t, failure).IsFailure()
	if failure.IsSuccess() {
		t.Error("IsSuccess() = true, want false")
	}
}

func TestFailure_Reason(t *testing.T) {
	commontest.AssertThat(t, failure).HasReason(failureReason)
	if r, ok := failure.Reason(); !ok || r != failureReason {
		t.Errorf("Reason() = %d, %v, want %d, true", r, ok, failureReason)
	}
}

func TestFailure_Map(t *testing.T) {
	commontest.AssertThat(t, common.Map(failure, double)).HasReason(failureReason)
}

func TestFailure_MapNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.Map[int, int, int](failure, nil) })
}

func TestFailure_MapToOtherType(t *testing.T) {
	commontest.AssertThat(t, common.Map(failure, describe)).HasReason(failureReason)
}

func TestFailure_MapFailure(t *testing.T) {
	commontest.AssertThat(t, common.MapFailure(failure, double)).HasReason(failureReason * 2)
}

func TestFailure_MapFailureNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.MapFailure[int, int, int](failure, nil) })
}

func TestFailure_MapFailureToOtherType(t *testing.T) {
	commontest.AssertThat(t, common.MapFailure(failure, describe)).HasReason("21 * 2 -> 42")
}

func TestFailure_TryMap(t *testing.T) {
	commontest.AssertThat(t, common.TryMap(failure, tryDouble, failureReason*2)).HasReason(failureReason)
}

func TestFailure_TryMapNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.TryMap[int, int, int](failure, nil, failureReason*2) })
}

func TestFailure_TryMapToOtherType(t *testing.T) {
	commontest.AssertThat(t, common.TryMap(failure, tryDescribe, failureReason*2)).HasReason(failureReason)
}

func TestFailure_TryMapFailingFunction(t *testing.T) {
	commontest.AssertThat(t, common.TryMap(failure, tryFail, failureReason*2)).HasReason(failureReason)
}

func TestFailure_Peek(t *testing.T) {
	output := -1
	got := failure.Peek(func(v int) { output = v })
	if output != -1 {
		t.Errorf("output = %d, want -1", output)
	}
	commontest.AssertThat(t, got).HasReason(failureReason)
}

func TestFailure_PeekNilConsumerPanics(t *testing.T) {
	assertPanics(t, func() { failure.Peek(nil) })
}

func TestFailure_TryPeek(t *testing.T) {
	output := -1
	got := failure.TryPeek(func(v int) error { output = v; return nil }, failureReason*2)
	if output != -1 {
		t.Errorf("output = %d, want -1", output)
	}
	commontest.AssertThat(t, got).HasReason(failureReason)
}

func TestFailure_TryPeekFailingConsumer(t *testing.T) {
	got := failure.TryPeek(func(int) error { return io.ErrUnexpectedEOF }, failureReason*2)
	commontest.AssertThat(t, got).HasReason(failureReason)
}

func TestFailure_TryPeekNilConsumerPanics(t *testing.T) {
	assertPanics(t, func() { failure.TryPeek(nil, failureReason*2) })
}

func TestFailure_Recover(t *testing.T) {
	commontest.AssertThat(t, failure.Recover(double)).HasValue(42)
}

func TestFailure_RecoverNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { failure.Recover(nil) })
}

func TestFailure_FlatMap(t *testing.T) {
	toSuccess := func(v int) intResult { return common.Success[int, int](v * 2) }
	toFailure := func(v int) intResult { return common.Failure[int, int](v * 2) }
	commontest.AssertThat(t, common.FlatMap(failure, toSuccess)).HasReason(failureReason)
	commontest.AssertThat(t, common.FlatMap(failure, toFailure)).HasReason(failureReason)
}

func TestFailure_FlatMapNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.FlatMap[int, int, int](failure, nil) })
}

func TestFailure_UnsafeFlattenNestedFailure(t *testing.T) {
	res := common.Failure[int, intResult](common.Failure[int, int](21))
	commontest.AssertThat(t, common.UnsafeFlatten[int, int](res)).HasReason(21)
}

func TestFailure_UnsafeFlattenNestedSuccess(t *testing.T) {
	res := common.Failure[int, intResult](common.Success[int, int](12))
	commontest.AssertThat(t, common.UnsafeFlatten[int, int](res)).HasValue(12)
}

func TestFailure_UnsafeFlattenNotNested(t *testing.T) {
	commontest.AssertThat(t, common.UnsafeFlatten[int, int](failure)).HasReason(failureReason)
}

func TestFailure_Fold(t *testing.T) {
	if got := common.Fold(failure, double, func(r int) int { return r + 9 }); got != 30 {
		t.Errorf("Fold() = %d, want 30", got)
	}
}

func TestFailure_FoldNilFunctionsPanic(t *testing.T) {
	testFoldNilFunctionsPanic(t, failure)
}

func TestFailure_Filter(t *testing.T) {
	called := false
	got := failure.Filter(func(v int) bool { called = true; return v == 12 }, 22)
	commontest.AssertThat(t, got).HasReason(failureReason)
	if called {
		t.Error("predicate was called for a failure")
	}
}

func TestFailure_FilterNilPredicatePanics(t *testing.T) {
	assertPanics(t, func() { failure.Filter(nil, 0) })
}

func TestFailure_Or(t *testing.T) {
	commontest.AssertThat(t, failure.Or(func() intResult { return common.Success[int, int](12) })).
		HasValue(12)
	commontest.AssertThat(t, failure.Or(func() intResult { return common.Failure[int, int](-11) })).
		HasReason(-11)
}

func TestFailure_OrNilSupplierPanics(t *testing.T) {
	assertPanics(t, func() { failure.Or(nil) })
}

func TestFailure_Swap(t *testing.T) {
	commontest.AssertThat(t, failure.Swap()).HasValue(failureReason)
}

func TestFailure_ValueIsAbsent(t *testing.T) {
	if v, ok := failure.Value(); ok || v != 0 {
		t.Errorf("Value() = %d, %v, want 0, false", v, ok)
	}
}

func TestFailure_AllYieldsNothing(t *testing.T) {
	if got := slices.Collect(failure.All()); len(got) != 0 {
		t.Errorf("All() = %v, want nothing", got)
	}
}
