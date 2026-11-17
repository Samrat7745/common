// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"errors"
	"fmt"
	"iter"
	"reflect"
)

// Result represents either a success value of type S or a failure reason of
// type F.
//
// Results are created with [Success], [Failure], [Of], [OfReason] and
// [OfMapped]. A Result is an immutable value; every operation returns a new
// one. Results can be compared with == if S and F are comparable.
//
// The zero value is a failure with the zero value of F as its reason.
type Result[S, F any] struct {
	value  S
	reason F
	ok     bool
}

// Success returns a successful Result holding the value.
func Success[S, F any](value S) Result[S, F] {
	return Result[S, F]{value: value, ok: true}
}

// Failure returns a failed Result holding the reason.
func Failure[S, F any](reason F) Result[S, F] {
	return Result[S, F]{reason: reason}
}

// Of calls the supplier and returns its value as a success, or a failure with
// the given reason if the supplier returns an error.
func Of[S, F any](supplier func() (S, error), reason F) Result[S, F] {
	requireNonNil(supplier == nil, "supplier")
	value, err := supplier()
	if err != nil {
		return Failure[S, F](reason)
	}
	return Success[S, F](value)
}

// OfReason calls the supplier and returns its value as a success. If the
// supplier returns an error, the result is a failure whose reason is the first
// error in the chain of that error that implements [FailureReason], or an
// [Unexpected] for the error if there is none.
func OfReason[S any](supplier func() (S, error)) Result[S, FailureReason] {
	requireNonNil(supplier == nil, "supplier")
	value, err := supplier()
	if err == nil {
		return Success[S, FailureReason](value)
	}
	var reason FailureReason
	if errors.As(err, &reason) {
		return Failure[S](reason)
	}
	return Failure[S, FailureReason](NewUnexpected(err))
}

// OfMapped calls the supplier and returns its value as a success. If the
// supplier returns an error, the result is a failure whose reason is looked up
// in reasons by the type of the error:
//
//  1. The dynamic type of the error's cause, as returned by [errors.Unwrap],
//     is used as the key.
//  2. The dynamic type of the error itself is used as the key.
//  3. A key that the error's type is assignable to is used. In practice these
//     are interface types the error implements, such as
//     reflect.TypeFor[error](), which matches every error. If there are
//     several, the one whose qualified name sorts first is chosen.
//  4. The missing reason is used.
//
// Keys are obtained with [reflect.TypeFor], for example
// reflect.TypeFor[*fs.PathError](). A nil map is treated as an empty one.
func OfMapped[S, F any](supplier func() (S, error), reasons map[reflect.Type]F, missing F) Result[S, F] {
	requireNonNil(supplier == nil, "supplier")
	value, err := supplier()
	if err == nil {
		return Success[S, F](value)
	}
	if cause := errors.Unwrap(err); cause != nil {
		if reason, ok := reasons[reflect.TypeOf(cause)]; ok {
			return Failure[S, F](reason)
		}
	}
	errType := reflect.TypeOf(err)
	if reason, ok := reasons[errType]; ok {
		return Failure[S, F](reason)
	}
	var match reflect.Type
	for key := range reasons {
		if key == nil || !errType.AssignableTo(key) {
			continue
		}
		if match == nil || qualifiedName(key) < qualifiedName(match) {
			match = key
		}
	}
	if match != nil {
		return Failure[S, F](reasons[match])
	}
	return Failure[S, F](missing)
}

// qualifiedName returns a name for the type that includes its package path.
func qualifiedName(t reflect.Type) string {
	return t.PkgPath() + "." + t.String()
}

// IsSuccess reports whether r is a success.
func (r Result[S, F]) IsSuccess() bool {
	return r.ok
}

// IsFailure reports whether r is a failure.
func (r Result[S, F]) IsFailure() bool {
	return !r.ok
}

// Value returns the success value and true if r is a success. Otherwise it
// returns the zero value of S and false.
func (r Result[S, F]) Value() (S, bool) {
	return r.value, r.ok
}

// Reason returns the failure reason and true if r is a failure. Otherwise it
// returns the zero value of F and false.
func (r Result[S, F]) Reason() (F, bool) {
	return r.reason, !r.ok
}

// All returns an iterator that yields the success value if r is a success and
// nothing otherwise.
func (r Result[S, F]) All() iter.Seq[S] {
	return func(yield func(S) bool) {
		if r.ok {
			yield(r.value)
		}
	}
}

// Peek calls the consumer with the success value if r is a success. It
// returns r.
func (r Result[S, F]) Peek(consumer func(S)) Result[S, F] {
	requireNonNil(consumer == nil, "consumer")
	if r.ok {
		consumer(r.value)
	}
	return r
}

// TryPeek calls the consumer with the success value if r is a success. It
// returns r, or a failure with the given reason if the consumer returns an
// error.
func (r Result[S, F]) TryPeek(consumer func(S) error, reason F) Result[S, F] {
	requireNonNil(consumer == nil, "consumer")
	if r.ok {
		if err := consumer(r.value); err != nil {
			return Failure[S, F](reason)
		}
	}
	return r
}

// Recover returns a success holding the result of applying the function to
// the failure reason if r is a failure. Otherwise it returns r.
func (r Result[S, F]) Recover(function func(F) S) Result[S, F] {
	requireNonNil(function == nil, "function")
	if r.ok {
		return r
	}
	return Success[S, F](function(r.reason))
}

// Filter returns r if r is a failure or if its success value matches the
// predicate. Otherwise it returns a failure with the given reason.
func (r Result[S, F]) Filter(predicate func(S) bool, reason F) Result[S, F] {
	requireNonNil(predicate == nil, "predicate")
	if !r.ok || predicate(r.value) {
		return r
	}
	return Failure[S, F](reason)
}

// Or returns r if r is a success, and the result of the supplier otherwise.
func (r Result[S, F]) Or(supplier func() Result[S, F]) Result[S, F] {
	requireNonNil(supplier == nil, "supplier")
	if r.ok {
		return r
	}
	return supplier()
}

// Swap turns a success into a failure with the success value as its reason,
// and a failure into a success with the failure reason as its value.
func (r Result[S, F]) Swap() Result[F, S] {
	if r.ok {
		return Failure[F, S](r.value)
	}
	return Success[F, S](r.reason)
}

// String returns "Success[value=...]" or "Failure[reason=...]".
func (r Result[S, F]) String() string {
	if r.ok {
		return fmt.Sprintf("Success[value=%v]", r.value)
	}
	return fmt.Sprintf("Failure[reason=%v]", r.reason)
}

// Map applies the function to the success value of r. A failure is returned
// with its reason unchanged.
func Map[S, F, N any](r Result[S, F], function func(S) N) Result[N, F] {
	requireNonNil(function == nil, "function")
	if r.ok {
		return Success[N, F](function(r.value))
	}
	return Failure[N, F](r.reason)
}

// MapFailure applies the function to the failure reason of r. A success is
// returned with its value unchanged.
func MapFailure[S, F, R any](r Result[S, F], function func(F) R) Result[S, R] {
	requireNonNil(function == nil, "function")
	if r.ok {
		return Success[S, R](r.value)
	}
	return Failure[S, R](function(r.reason))
}

// TryMap applies the function to the success value of r. If the function
// returns an error, the result is a failure with the given reason. A failure
// is returned with its reason unchanged.
func TryMap[S, F, N any](r Result[S, F], function func(S) (N, error), reason F) Result[N, F] {
	requireNonNil(function == nil, "function")
	if !r.ok {
		return Failure[N, F](r.reason)
	}
	value, err := function(r.value)
	if err != nil {
		return Failure[N, F](reason)
	}
	return Success[N, F](value)
}

// FlatMap applies the Result-bearing function to the success value of r and
// returns its result. A failure is returned with its reason unchanged.
func FlatMap[S, F, N any](r Result[S, F], function func(S) Result[N, F]) Result[N, F] {
	requireNonNil(function == nil, "function")
	if r.ok {
		return function(r.value)
	}
	return Failure[N, F](r.reason)
}

// Fold applies onSuccess to the success value or onFailure to the failure
// reason of r, and returns the outcome.
func Fold[S, F, N any](r Result[S, F], onSuccess func(S) N, onFailure func(F) N) N {
	requireNonNil(onSuccess == nil, "success function")
	requireNonNil(onFailure == nil, "failure function")
	if r.ok {
		return onSuccess(r.value)
	}
	return onFailure(r.reason)
}

// UnsafeFlatten removes one level of nesting: if the success value or the
// failure reason of r is itself a Result[V, R], that nested result is
// returned. Otherwise r is returned.
//
// The types V and R have to be given explicitly, for example
// UnsafeFlatten[int, string](r). The function panics if neither the nested
// result nor r is a Result[V, R].
func UnsafeFlatten[V, R, S, F any](r Result[S, F]) Result[V, R] {
	var content any = r.reason
	if r.ok {
		content = r.value
	}
	if nested, ok := content.(Result[V, R]); ok {
		return nested
	}
	if same, ok := any(r).(Result[V, R]); ok {
		return same
	}
	panic(fmt.Sprintf("common: cannot flatten %T to %T", r, Result[V, R]{}))
}
