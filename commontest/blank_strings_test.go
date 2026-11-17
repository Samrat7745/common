// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package commontest_test

import (
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Samrat7745/common/commontest"
)

// isJavaWhitespace reports whether Java's Character.isWhitespace is true for
// the character: a space, line or paragraph separator that is not a
// non-breaking space, or one of the listed control characters.
func isJavaWhitespace(r rune) bool {
	switch {
	case r >= '\t' && r <= '\r', r >= '\u001C' && r <= '\u001F':
		return true
	case r == ' ', r == ' ', r == ' ':
		return false
	case r == ' ', r == ' ', r >= ' ' && r <= ' ', r == ' ', r == '　':
		return true
	case r == ' ', r == ' ':
		return true
	}
	return false
}

func isBlank(s string) bool {
	return strings.TrimFunc(s, isJavaWhitespace) == ""
}

func TestBlankStrings_AreBlank(t *testing.T) {
	blanks := commontest.BlankStrings()
	if len(blanks) != 26 {
		t.Fatalf("got %d blank strings, want 26", len(blanks))
	}
	for _, blank := range blanks {
		t.Run(blank.Name, func(t *testing.T) {
			if !isBlank(blank.Value) {
				t.Errorf("%q is not blank", blank.Value)
			}
		})
	}
}

func TestBlankStrings_AreDistinct(t *testing.T) {
	names := map[string]bool{}
	values := map[string]bool{}
	for _, blank := range commontest.BlankStrings() {
		if names[blank.Name] || values[blank.Value] {
			t.Errorf("%s (%q) occurs twice", blank.Name, blank.Value)
		}
		names[blank.Name] = true
		values[blank.Value] = true
		if n := utf8.RuneCountInString(blank.Value); (blank.Name == "EMPTY") != (n == 0) || n > 1 {
			t.Errorf("%s has %d characters", blank.Name, n)
		}
	}
}

func TestBlankStrings_ReturnsFreshSlice(t *testing.T) {
	commontest.BlankStrings()[0].Value = "changed"
	if got := commontest.BlankStrings()[0]; got.Name != "EMPTY" || got.Value != "" {
		t.Errorf("first blank string = %+v, want the empty string", got)
	}
}

// The documented difference to Go's definition of white space.
func TestBlankStrings_NotBlankForTrimSpace(t *testing.T) {
	var got []string
	for _, blank := range commontest.BlankStrings() {
		if strings.TrimSpace(blank.Value) != "" {
			got = append(got, blank.Name)
		}
	}
	want := []string{"FILE SEPARATOR", "GROUP SEPARATOR", "RECORD SEPARATOR", "UNIT SEPARATOR"}
	if !slices.Equal(got, want) {
		t.Errorf("not blank for strings.TrimSpace: %q, want %q", got, want)
	}
}
