// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"fmt"
	"io"
	"slices"
	"testing"

	"github.com/Samrat7745/common"
	"github.com/Samrat7745/common/commontest"
)

const successValue = 12

var success = common.Success[int, int](successValue)

func double(v int) int { return v * 2 }

func describe(v int) string { return fmt.Sprintf("%d * 2 -> %d", v, v*2) }

func tryDouble(v int) (int, error) { return v * 2, nil }

func tryDescribe(v int) (string, error) { return describe(v), nil }

func tryFail(int) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestSuccess_IsSuccess(t *testing.T) {
	if !success.IsSuccess() {
		t.Error("IsSuccess() = false, want true")
	}
	commontest.AssertThat(t, success).IsSuccess()
	if success.IsFailure() {
		t.Error("IsFailure() = true, want false")
	}
}

func TestSuccess_Value(t *testing.T) {
	commontest.AssertThat(t, success).HasValue(successValue)
	if v, ok := success.Value(); !ok || v != successValue {
		t.Errorf("Value() = %d, %v, want %d, true", v, ok, successValue)
	}
	if r, ok := success.Reason(); ok || r != 0 {
		t.Errorf("Reason() = %d, %v, want 0, false", r, ok)
	}
}

func TestSuccess_Map(t *testing.T) {
	commontest.AssertThat(t, common.Map(success, double)).HasValue(24)
}

func TestSuccess_MapNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.Map[int, int, int](success, nil) })
}

func TestSuccess_MapToOtherType(t *testing.T) {
	commontest.AssertThat(t, common.Map(success, describe)).HasValue("12 * 2 -> 24")
}

func TestSuccess_MapFailure(t *testing.T) {
	commontest.AssertThat(t, common.MapFailure(success, double)).HasValue(successValue)
}

func TestSuccess_MapFailureNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.MapFailure[int, int, int](success, nil) })
}

func TestSuccess_MapFailureToOtherType(t *testing.T) {
	commontest.AssertThat(t, common.MapFailure(success, describe)).HasValue(successValue)
}

func TestSuccess_TryMap(t *testing.T) {
	commontest.AssertThat(t, common.TryMap(success, tryDouble, -11)).HasValue(24)
}

func TestSuccess_TryMapNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.TryMap[int, int, int](success, nil, -11) })
}

func TestSuccess_TryMapToOtherType(t *testing.T) {
	commontest.AssertThat(t, common.TryMap(success, tryDescribe, -11)).HasValue("12 * 2 -> 24")
}

func TestSuccess_TryMapFailingFunction(t *testing.T) {
	commontest.AssertThat(t, common.TryMap(success, tryFail, -11)).HasReason(-11)
}

func TestSuccess_Peek(t *testing.T) {
	output := -1
	got := success.Peek(func(v int) { output = v })
	if output != successValue {
		t.Errorf("output = %d, want %d", output, successValue)
	}
	commontest.AssertThat(t, got).HasValue(successValue)
}

func TestSuccess_PeekNilConsumerPanics(t *testing.T) {
	assertPanics(t, func() { success.Peek(nil) })
}

func TestSuccess_TryPeek(t *testing.T) {
	output := -1
	got := success.TryPeek(func(v int) error { output = v; return nil }, -11)
	if output != successValue {
		t.Errorf("output = %d, want %d", output, successValue)
	}
	commontest.AssertThat(t, got).HasValue(successValue)
}

func TestSuccess_TryPeekFailingConsumer(t *testing.T) {
	got := success.TryPeek(func(int) error { return io.ErrUnexpectedEOF }, -11)
	commontest.AssertThat(t, got).HasReason(-11)
}

func TestSuccess_TryPeekNilConsumerPanics(t *testing.T) {
	assertPanics(t, func() { success.TryPeek(nil, -11) })
}

func TestSuccess_Recover(t *testing.T) {
	commontest.AssertThat(t, success.Recover(identity)).HasValue(successValue)
}

func TestSuccess_RecoverNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { success.Recover(nil) })
}

func TestSuccess_FlatMap(t *testing.T) {
	toSuccess := func(v int) intResult { return common.Success[int, int](v * 2) }
	toFailure := func(v int) intResult { return common.Failure[int, int](v + 9) }
	commontest.AssertThat(t, common.FlatMap(success, toSuccess)).HasValue(24)
	commontest.AssertThat(t, common.FlatMap(success, toFailure)).HasReason(21)
}

func TestSuccess_FlatMapNilFunctionPanics(t *testing.T) {
	assertPanics(t, func() { common.FlatMap[int, int, int](success, nil) })
}

func TestSuccess_UnsafeFlattenNestedSuccess(t *testing.T) {
	res := common.Success[intResult, int](common.Success[int, int](12))
	commontest.AssertThat(t, common.UnsafeFlatten[int, int](res)).HasValue(12)
}

func TestSuccess_UnsafeFlattenNestedFailure(t *testing.T) {
	res := common.Success[intResult, int](common.Failure[int, int](21))
	commontest.AssertThat(t, common.UnsafeFlatten[int, int](res)).HasReason(21)
}

func TestSuccess_UnsafeFlattenNotNested(t *testing.T) {
	commontest.AssertThat(t, common.UnsafeFlatten[int, int](success)).HasValue(successValue)
}

func TestSuccess_Fold(t *testing.T) {
	if got := common.Fold(success, double, func(r int) int { return r + 9 }); got != 24 {
		t.Errorf("Fold() = %d, want 24", got)
	}
}

func TestSuccess_FoldNilFunctionsPanic(t *testing.T) {
	testFoldNilFunctionsPanic(t, success)
}

func testFoldNilFunctionsPanic(t *testing.T, result intResult) {
	combinations := commontest.NilCombinations(identity, identity)
	if len(combinations) != 3 {
		t.Fatalf("got %d combinations, want 3", len(combinations))
	}
	for i, args := range combinations {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			onSuccess, _ := args[0].(func(int) int)
			onFailure, _ := args[1].(func(int) int)
			assertPanics(t, func() { common.Fold(result, onSuccess, onFailure) })
		})
	}
}

func TestSuccess_Filter(t *testing.T) {
	commontest.AssertThat(t, success.Filter(func(v int) bool { return v == successValue }, 21)).
		HasValue(successValue)
	commontest.AssertThat(t, success.Filter(func(v int) bool { return v != successValue }, 21)).
		HasReason(21)
}

func TestSuccess_FilterNilPredicatePanics(t *testing.T) {
	assertPanics(t, func() { success.Filter(nil, 0) })
}

func TestSuccess_Or(t *testing.T) {
	called := false
	got := success.Or(func() intResult { called = true; return common.Success[int, int](24) })
	commontest.AssertThat(t, got).HasValue(successValue)
	if called {
		t.Error("supplier was called for a success")
	}
}

func TestSuccess_OrNilSupplierPanics(t *testing.T) {
	assertPanics(t, func() { success.Or(nil) })
}

func TestSuccess_Swap(t *testing.T) {
	commontest.AssertThat(t, success.Swap()).HasReason(successValue)
}

func TestSuccess_All(t *testing.T) {
	if got := slices.Collect(success.All()); !slices.Equal(got, []int{successValue}) {
		t.Errorf("All() = %v, want [%d]", got, successValue)
	}
}

// Unlike Java's toOptional, Value reports a nil success value as present.
func TestSuccess_NilValueIsPresent(t *testing.T) {
	res := common.Success[*int, int](nil)
	if v, ok := res.Value(); !ok || v != nil {
		t.Errorf("Value() = %v, %v, want nil, true", v, ok)
	}
	if got := slices.Collect(res.All()); len(got) != 1 {
		t.Errorf("All() yields %d values, want 1", len(got))
	}
}
