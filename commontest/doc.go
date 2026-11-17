// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

// Package commontest provides helpers for tests of code that uses the package
// common: assertions on results and generators of test inputs.
//
//	func TestParse(t *testing.T) {
//		commontest.AssertThat(t, Parse("12")).HasValue(12)
//
//		for _, blank := range commontest.BlankStrings() {
//			t.Run(blank.Name, func(t *testing.T) {
//				commontest.AssertThat(t, Parse(blank.Value)).IsFailure()
//			})
//		}
//	}
package commontest
