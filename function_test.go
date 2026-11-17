// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"io"
	"testing"

	"github.com/Samrat7745/common"
	"github.com/Samrat7745/common/commontest"
)

// The Java function package (ThrowingSupplier, ThrowingConsumer,
// ThrowingFunction and their AnyThrowing variants) has no Go counterpart: a
// function returning an error takes its place. These tests cover the three
// shapes as the Result API accepts them.

func TestSupplierReturningError(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		supplier := func() (int, error) { return 12, nil }
		commontest.AssertThat(t, common.Of(supplier, "failed")).HasValue(12)
	})
	t.Run("error", func(t *testing.T) {
		supplier := func() (int, error) { return 0, io.ErrUnexpectedEOF }
		if _, err := supplier(); err != io.ErrUnexpectedEOF {
			t.Errorf("supplier() error = %v, want %v", err, io.ErrUnexpectedEOF)
		}
		commontest.AssertThat(t, common.Of(supplier, "failed")).HasReason("failed")
	})
}

func TestConsumerReturningError(t *testing.T) {
	result := common.Success[int, string](12)
	t.Run("value", func(t *testing.T) {
		output := -1
		consumer := func(v int) error { output = v; return nil }
		commontest.AssertThat(t, result.TryPeek(consumer, "failed")).HasValue(12)
		if output != 12 {
			t.Errorf("output = %d, want 12", output)
		}
	})
	t.Run("error", func(t *testing.T) {
		consumer := func(int) error { return io.ErrUnexpectedEOF }
		if err := consumer(12); err != io.ErrUnexpectedEOF {
			t.Errorf("consumer() error = %v, want %v", err, io.ErrUnexpectedEOF)
		}
		commontest.AssertThat(t, result.TryPeek(consumer, "failed")).HasReason("failed")
	})
}

func TestFunctionReturningError(t *testing.T) {
	result := common.Success[int, string](12)
	t.Run("value", func(t *testing.T) {
		function := func(v int) (int, error) { return v, nil }
		commontest.AssertThat(t, common.TryMap(result, function, "failed")).HasValue(12)
	})
	t.Run("error", func(t *testing.T) {
		function := func(int) (int, error) { return 0, io.ErrUnexpectedEOF }
		if _, err := function(12); err != io.ErrUnexpectedEOF {
			t.Errorf("function() error = %v, want %v", err, io.ErrUnexpectedEOF)
		}
		commontest.AssertThat(t, common.TryMap(result, function, "failed")).HasReason("failed")
	})
}
