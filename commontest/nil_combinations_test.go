// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package commontest_test

import (
	"reflect"
	"testing"

	"github.com/Samrat7745/common/commontest"
)

func TestNilCombinations(t *testing.T) {
	tests := []struct {
		name     string
		baseline []any
		want     [][]any
	}{
		{"no values", nil, nil},
		{"one value", []any{"a"}, [][]any{{nil}}},
		{"two values", []any{"a", 1}, [][]any{{"a", nil}, {nil, 1}, {nil, nil}}},
		{"three values", []any{"a", 1, true}, [][]any{
			{"a", 1, nil}, {"a", nil, true}, {"a", nil, nil},
			{nil, 1, true}, {nil, 1, nil}, {nil, nil, true}, {nil, nil, nil},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commontest.NilCombinations(tt.baseline...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilCombinations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilCombinations_DoesNotModifyBaseline(t *testing.T) {
	baseline := []any{"a", 1}
	commontest.NilCombinations(baseline...)
	if !reflect.DeepEqual(baseline, []any{"a", 1}) {
		t.Errorf("baseline = %v, want it unchanged", baseline)
	}
}

func TestNilCombinations_TypeAssertionYieldsNilFunction(t *testing.T) {
	for _, args := range commontest.NilCombinations(func() {}) {
		if f, _ := args[0].(func()); f != nil {
			t.Error("replaced function is not nil")
		}
	}
}
