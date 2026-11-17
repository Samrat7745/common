// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package commontest

// NamedString is a string with a name that describes it, for use as the name
// of a subtest.
type NamedString struct {
	Name  string
	Value string
}

// BlankStrings returns the empty string and one string for each character
// that may occur in a blank string, each with the Unicode name of the
// character.
//
// The characters are those for which Java's Character.isWhitespace is true.
// This is not the set of [unicode.IsSpace]: the four separators FILE, GROUP,
// RECORD and UNIT SEPARATOR (U+001C to U+001F) are in the list although
// unicode.IsSpace and [strings.TrimSpace] do not treat them as white space.
// Skip them by name if the code under test follows the Go definition.
func BlankStrings() []NamedString {
	return []NamedString{
		{"EMPTY", ""},
		{"SPACE", " "},
		{"EN QUAD", " "},
		{"EM QUAD", " "},
		{"EN SPACE", " "},
		{"EM SPACE", " "},
		{"THREE-PER-EM SPACE", " "},
		{"FOUR-PER-EM SPACE", " "},
		{"SIX-PER-EM SPACE", " "},
		{"PUNCTUATION SPACE", " "},
		{"THIN SPACE", " "},
		{"HAIR SPACE", " "},
		{"MEDIUM MATHEMATICAL SPACE", " "},
		{"IDEOGRAPHIC SPACE", "　"},
		{"OGHAM SPACE MARK", " "},
		{"HORIZONTAL TABULATION", "\t"},
		{"LINE FEED", "\n"},
		{"VERTICAL TABULATION", "\u000B"},
		{"FORM FEED", "\f"},
		{"CARRIAGE RETURN", "\r"},
		{"FILE SEPARATOR", "\u001C"},
		{"GROUP SEPARATOR", "\u001D"},
		{"RECORD SEPARATOR", "\u001E"},
		{"UNIT SEPARATOR", "\u001F"},
		{"LINE SEPARATOR", " "},
		{"PARAGRAPH SEPARATOR", " "},
	}
}
