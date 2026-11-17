// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package commontest

// NilCombinations returns every combination of the baseline values in which
// at least one value is replaced by nil. It is meant for tests that check how
// a function with several arguments reacts to each of them being nil.
//
// For n baseline values the result has 2^n-1 entries of n values each. The
// combination with only the last value replaced comes first and the one with
// all values replaced comes last: for the baseline (a, b) the result is
// (a, nil), (nil, b), (nil, nil). Without baseline values the result is empty.
//
// The values are returned as any. Convert them back with the two-value form of
// a type assertion, which yields the nil value of the type for a replaced
// value:
//
//	for _, args := range commontest.NilCombinations(parse, format) {
//		parse, _ := args[0].(func(string) int)
//		format, _ := args[1].(func(int) string)
//		...
//	}
func NilCombinations(baseline ...any) [][]any {
	n := len(baseline)
	if n == 0 {
		return nil
	}
	combinations := make([][]any, 0, 1<<n-1)
	for mask := 1; mask < 1<<n; mask++ {
		combination := make([]any, n)
		for i, value := range baseline {
			if mask&(1<<(n-1-i)) == 0 {
				combination[i] = value
			}
		}
		combinations = append(combinations, combination)
	}
	return combinations
}
