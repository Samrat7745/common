// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"testing"

	"github.com/Samrat7745/common"
	"github.com/Samrat7745/common/commontest"
)

type intResult = common.Result[int, int]

func failing() (int, error) { return 0, io.ErrUnexpectedEOF }

func twelve() (int, error) { return 12, nil }

func wrapping() (int, error) { return 0, &wrappingError{cause: plainError{}} }

func noop(int) {}

func fold(r intResult) int { return common.Fold(r, identity, identity) }

func TestResult_ChainFromFailureToSuccess(t *testing.T) {
	value := 12
	other := common.Success[int, int](value)
	res := common.Of(failing, 21)
	res = common.Map(res, func(v int) int { return v * 2 }).
		Peek(noop).
		Recover(identity).
		Peek(noop)
	res = common.FlatMap(res, func(v int) intResult { return common.Failure[int, int](v - 10) }).
		Peek(noop).
		Filter(func(v int) bool { return v > 10 }, -1).
		Or(func() intResult { return other }).
		Or(func() intResult {
			return common.FlatMap(common.Failure[int, int](-11), func(v int) intResult {
				return common.Success[int, int](v)
			})
		}).
		Filter(func(v int) bool { return v > 10 }, -1)

	if !res.IsSuccess() || res.IsFailure() {
		t.Errorf("got %v, want a success", res)
	}
	if got := fold(res); got != value {
		t.Errorf("Fold() = %d, want %d", got, value)
	}
	if _, ok := res.Value(); !ok {
		t.Error("Value() reports no value")
	}
	if got := slices.Collect(res.All()); !slices.Equal(got, []int{value}) {
		t.Errorf("All() = %v, want [%d]", got, value)
	}
}

func TestResult_ChainFromSuccessToFailure(t *testing.T) {
	res := common.Of(func() (int, error) { return 21, nil }, 21)
	res = common.Map(res, func(v int) int { return v * 2 }).
		Peek(noop).
		Recover(identity).
		Peek(noop)
	res = common.FlatMap(res, func(v int) intResult { return common.Failure[int, int](v - 10) }).
		Peek(noop).
		Filter(func(v int) bool { return v < 10 }, -1).
		Or(func() intResult { return common.Failure[int, int](-11) })

	if res.IsSuccess() || !res.IsFailure() {
		t.Errorf("got %v, want a failure", res)
	}
	if got := fold(res); got != -11 {
		t.Errorf("Fold() = %d, want -11", got)
	}
	if _, ok := res.Value(); ok {
		t.Error("Value() reports a value")
	}
	if got := slices.Collect(res.All()); len(got) != 0 {
		t.Errorf("All() = %v, want nothing", got)
	}
}

func TestResult_ChainFlatMapRecoverFilter(t *testing.T) {
	res := common.FlatMap(common.Failure[int, int](21), common.Success[int, int]).
		Recover(identity).
		Filter(func(v int) bool { return v < 5 }, -21)

	if res.IsSuccess() || !res.IsFailure() {
		t.Errorf("got %v, want a failure", res)
	}
	if got := fold(res); got != -21 {
		t.Errorf("Fold() = %d, want -21", got)
	}
	if got := slices.Collect(res.All()); len(got) != 0 {
		t.Errorf("All() = %v, want nothing", got)
	}
}

func TestResult_ChainMapFailureAndFlatMap(t *testing.T) {
	widened := common.MapFailure(common.Success[int, int](12), func(r int) int64 { return int64(r) })
	failed := common.FlatMap(widened, func(v int) common.Result[int, int64] {
		return common.Failure[int, int64](int64(v) * 2)
	})
	res := common.MapFailure(failed, func(r int64) int { return int(r) })
	commontest.AssertThat(t, res).HasReason(24)
}

func TestResult_ChainTryPeek(t *testing.T) {
	calls := 0
	count := func(int) error { calls++; return nil }
	res := common.Success[int, int](12).
		TryPeek(count, -1).
		TryPeek(func(int) error { return io.ErrUnexpectedEOF }, -1).
		TryPeek(count, -1)
	commontest.AssertThat(t, res).HasReason(-1)
	if calls != 1 {
		t.Errorf("consumer called %d times, want 1", calls)
	}
}

func TestResult_ChainTryMap(t *testing.T) {
	double := func(v int) (int, error) { return v * 2, nil }
	res := common.TryMap(common.Success[int, int](12), double, -11)
	res = common.TryMap(res, func(int) (int, error) { return 0, errors.New("failed") }, -21)
	res = common.TryMap(res, double, -11)
	commontest.AssertThat(t, res).HasReason(-21)
}

func TestOf_Success(t *testing.T) {
	commontest.AssertThat(t, common.Of(twelve, "some")).HasValue(12)
}

func TestOf_Error(t *testing.T) {
	commontest.AssertThat(t, common.Of(failing, "some")).HasReason("some")
}

func TestOf_NilSupplierPanics(t *testing.T) {
	assertPanics(t, func() { common.Of[int](nil, "some") })
	assertPanics(t, func() { common.Of[int, *string](nil, nil) })
}

// Unlike in Java, the reason is not checked: nil is a reason like any other.
func TestOf_NilReasonIsAccepted(t *testing.T) {
	commontest.AssertThat(t, common.Of[int, *string](twelve, nil)).HasValue(12)
	commontest.AssertThat(t, common.Of[int, *string](failing, nil)).HasReason(nil)
}

func TestOfReason_Success(t *testing.T) {
	commontest.AssertThat(t, common.OfReason(twelve)).HasValue(12)
}

func TestOfReason_ErrorIsFailureReason(t *testing.T) {
	reason := newFailureError()
	res := common.OfReason(func() (int, error) { return 0, reason })
	commontest.AssertThat(t, res).HasReason(reason)
	if got, _ := res.Reason(); got != common.FailureReason(reason) {
		t.Errorf("Reason() = %v, want the returned error itself", got)
	}
}

func TestOfReason_WrappedFailureReason(t *testing.T) {
	reason := newFailureError()
	res := common.OfReason(func() (int, error) { return 0, fmt.Errorf("wrapped: %w", reason) })
	if got, _ := res.Reason(); got != common.FailureReason(reason) {
		t.Errorf("Reason() = %v, want the wrapped failure reason", got)
	}
}

func TestOfReason_OtherErrorIsUnexpected(t *testing.T) {
	res := common.OfReason(failing)
	reason, ok := res.Reason()
	if !ok {
		t.Fatalf("got %v, want a failure", res)
	}
	unexpected, ok := reason.(*common.Unexpected)
	if !ok {
		t.Fatalf("Reason() is a %T, want *common.Unexpected", reason)
	}
	if unexpected.Err() != io.ErrUnexpectedEOF {
		t.Errorf("Err() = %v, want %v", unexpected.Err(), io.ErrUnexpectedEOF)
	}
	if got := common.Description(reason); got != io.ErrUnexpectedEOF.Error() {
		t.Errorf("Description() = %q, want %q", got, io.ErrUnexpectedEOF.Error())
	}
}

func TestOfReason_NilSupplierPanics(t *testing.T) {
	assertPanics(t, func() { common.OfReason[int](nil) })
}

type reasons = map[reflect.Type]*failureError

func TestOfMapped_SuccessWithEmptyMapping(t *testing.T) {
	commontest.AssertThat(t, common.OfMapped(twelve, reasons{}, newFailureError())).HasValue(12)
}

func TestOfMapped_ErrorWithEmptyMapping(t *testing.T) {
	missing := newFailureError()
	commontest.AssertThat(t, common.OfMapped(failing, reasons{}, missing)).HasReason(missing)
}

func TestOfMapped_ExactErrorType(t *testing.T) {
	reason := &failureError{"IO", "IOException"}
	mapping := reasons{reflect.TypeFor[plainError](): reason}
	res := common.OfMapped(func() (int, error) { return 0, plainError{} }, mapping, newFailureError())
	commontest.AssertThat(t, res).HasReason(reason)
}

func TestOfMapped_CauseErrorType(t *testing.T) {
	reason := &failureError{"IO", "IOException"}
	mapping := reasons{reflect.TypeFor[plainError](): reason}
	commontest.AssertThat(t, common.OfMapped(wrapping, mapping, newFailureError())).HasReason(reason)
}

func TestOfMapped_CauseTakesPrecedenceOverError(t *testing.T) {
	reason := &failureError{"cause", "cause"}
	mapping := reasons{
		reflect.TypeFor[plainError]():     reason,
		reflect.TypeFor[*wrappingError](): {"error", "error"},
		reflect.TypeFor[error]():          {"any", "any"},
	}
	commontest.AssertThat(t, common.OfMapped(wrapping, mapping, newFailureError())).HasReason(reason)
}

func TestOfMapped_ErrorTakesPrecedenceOverAssignable(t *testing.T) {
	reason := &failureError{"error", "error"}
	mapping := reasons{
		reflect.TypeFor[*wrappingError](): reason,
		reflect.TypeFor[error]():          {"any", "any"},
	}
	commontest.AssertThat(t, common.OfMapped(wrapping, mapping, newFailureError())).HasReason(reason)
}

func TestOfMapped_AssignableErrorType(t *testing.T) {
	reason := &failureError{"IO", "IOException"}
	mapping := reasons{reflect.TypeFor[error](): reason}
	commontest.AssertThat(t, common.OfMapped(wrapping, mapping, newFailureError())).HasReason(reason)
}

type unwrapper interface{ Unwrap() error }

func TestOfMapped_SeveralAssignableErrorTypes(t *testing.T) {
	reason := &failureError{"any", "any"}
	mapping := reasons{
		nil:                           {"nil", "nil"},
		reflect.TypeFor[otherError](): {"other", "other"},
		reflect.TypeFor[unwrapper]():  {"unwrapper", "unwrapper"},
		reflect.TypeFor[error]():      reason,
		reflect.TypeFor[interface {
			error
			Unwrap() error
		}](): {"anonymous", "anonymous"},
	}
	// Map iteration order is random; the choice must not depend on it.
	for range 100 {
		commontest.AssertThat(t, common.OfMapped(wrapping, mapping, newFailureError())).HasReason(reason)
	}
}

func TestOfMapped_MissingMappingEntry(t *testing.T) {
	missing := newFailureError()
	mapping := reasons{reflect.TypeFor[otherError](): {"IO", "IOException"}}
	commontest.AssertThat(t, common.OfMapped(wrapping, mapping, missing)).HasReason(missing)
}

func TestOfMapped_NilSupplierPanics(t *testing.T) {
	assertPanics(t, func() { common.OfMapped[int](nil, reasons{}, newFailureError()) })
	assertPanics(t, func() { common.OfMapped[int](nil, reasons(nil), newFailureError()) })
}

// Unlike in Java, a nil mapping is accepted and treated as an empty one.
func TestOfMapped_NilMappingIsEmpty(t *testing.T) {
	missing := newFailureError()
	commontest.AssertThat(t, common.OfMapped(twelve, nil, missing)).HasValue(12)
	commontest.AssertThat(t, common.OfMapped(failing, nil, missing)).HasReason(missing)
}

func TestUnsafeFlatten(t *testing.T) {
	value, reason := 12, 21
	success := common.Success[int, int](value)
	failure := common.Failure[int, int](reason)

	tests := []struct {
		name string
		got  intResult
		want int
	}{
		{"success in success", common.UnsafeFlatten[int, int](common.Success[intResult, int](success)), value},
		{"failure in success", common.UnsafeFlatten[int, int](common.Success[intResult, int](failure)), reason},
		{"success", common.UnsafeFlatten[int, int](success), value},
		{"failure in failure", common.UnsafeFlatten[int, int](common.Failure[int, intResult](failure)), reason},
		{"success in failure", common.UnsafeFlatten[int, int](common.Failure[int, intResult](success)), value},
		{"failure", common.UnsafeFlatten[int, int](failure), reason},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fold(tt.got); got != tt.want {
				t.Errorf("Fold() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestUnsafeFlatten_TypeMismatchPanics(t *testing.T) {
	assertPanics(t, func() { common.UnsafeFlatten[string, string](common.Success[int, int](12)) })
	nested := common.Success[intResult, int](common.Success[int, int](12))
	assertPanics(t, func() { common.UnsafeFlatten[string, string](nested) })
}

func TestResult_Swap(t *testing.T) {
	value := 12
	res := common.Success[int, int](value)
	minusOne := func(int) int { return -1 }
	if got := common.Fold(res, identity, minusOne); got != value {
		t.Errorf("Fold() = %d, want %d", got, value)
	}
	if got := common.Fold(res.Swap(), identity, minusOne); got != -1 {
		t.Errorf("Fold() after Swap() = %d, want -1", got)
	}
	if got := common.Fold(res.Swap().Swap(), identity, minusOne); got != value {
		t.Errorf("Fold() after two Swap() = %d, want %d", got, value)
	}
}

func TestResult_String(t *testing.T) {
	if got, want := common.Success[int, string](12).String(), "Success[value=12]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if got, want := common.Failure[int, string]("some").String(), "Failure[reason=some]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if got, want := fmt.Sprint(common.Failure[int, int](21)), "Failure[reason=21]"; got != want {
		t.Errorf("Sprint() = %q, want %q", got, want)
	}
}

func TestResult_Equality(t *testing.T) {
	if common.Success[int, int](12) != common.Success[int, int](12) {
		t.Error("equal successes are not ==")
	}
	if common.Failure[int, int](21) != common.Failure[int, int](21) {
		t.Error("equal failures are not ==")
	}
	if common.Success[int, int](12) == common.Success[int, int](13) {
		t.Error("different successes are ==")
	}
	if common.Success[int, int](12) == common.Failure[int, int](12) {
		t.Error("a success and a failure are ==")
	}
	if common.Success[int, int](0) == common.Failure[int, int](0) {
		t.Error("a success and a failure of the zero value are ==")
	}
}

func TestResult_ZeroValueIsFailure(t *testing.T) {
	var res common.Result[int, string]
	commontest.AssertThat(t, res).HasReason("")
	if res != common.Failure[int, string]("") {
		t.Error("zero value is not == Failure of the zero reason")
	}
}

func TestResult_AllStopsWhenAsked(t *testing.T) {
	for range common.Success[int, int](12).All() {
		break
	}
}
