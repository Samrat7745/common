// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

// Package common provides the types shared by the Carbyne Stack services,
// clients and other components.
//
// Its centre is [Result], which holds either a success value or a failure
// reason:
//
//	res := common.Of(func() (int, error) { return strconv.Atoi(input) }, "not a number")
//	doubled := common.Map(res, func(v int) int { return v * 2 })
//	if v, ok := doubled.Value(); ok {
//		fmt.Println(v)
//	}
//
// Go methods cannot declare type parameters of their own. The operations that
// change the success or the failure type of a result ([Map], [MapFailure],
// [TryMap], [FlatMap], [Fold], [UnsafeFlatten]) are therefore package-level
// functions that take the result as their first argument; all other operations
// are methods of [Result].
//
// Failure reasons can be values of any type. Reasons that describe themselves
// implement [FailureReason]; [Unexpected] and [ErrorReason] adapt an error to
// that interface.
//
// A nil function passed to any operation is a programming error and causes a
// panic, whether or not the operation would have called the function.
//
// The package commontest contains helpers for tests that work with results.
package common
