// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Samrat7745/common"
)

var errTest = errors.New("test")

func TestUnexpected_NilErrorPanics(t *testing.T) {
	assertPanics(t, func() { common.NewUnexpected(nil) })
}

func TestUnexpected_Synopsis(t *testing.T) {
	unknown := common.NewUnexpected(errTest)
	if got, want := unknown.Synopsis(), "An unknown exception has occurred."; got != want {
		t.Errorf("Synopsis() = %q, want %q", got, want)
	}
}

func TestUnexpected_Description(t *testing.T) {
	unknown := common.NewUnexpected(errTest)
	if got := unknown.Description(); got != "test" {
		t.Errorf("Description() = %q, want %q", got, "test")
	}
	if got := common.Description(unknown); got != "test" {
		t.Errorf("common.Description() = %q, want %q", got, "test")
	}
}

func TestUnexpected_StackTrace(t *testing.T) {
	unknown := common.NewUnexpected(errTest)
	frames, ok := common.StackTrace(unknown)
	if !ok || len(frames) == 0 {
		t.Fatalf("StackTrace() = %v, %v, want frames", frames, ok)
	}
	if !strings.HasSuffix(frames[0].Function, ".TestUnexpected_StackTrace") {
		t.Errorf("innermost frame = %q, want the creating function", frames[0].Function)
	}
}

func TestUnexpected_ReportIssue(t *testing.T) {
	unknown := common.NewUnexpected(errTest)
	if !unknown.ReportIssue() || !common.ReportIssue(unknown) {
		t.Error("ReportIssue() = false, want true")
	}
}

func TestUnexpected_Err(t *testing.T) {
	if got := common.NewUnexpected(errTest).Err(); got != errTest {
		t.Errorf("Err() = %v, want %v", got, errTest)
	}
}
