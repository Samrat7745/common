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

var errSome = errors.New("Some message.")

func TestErrorReason_NilErrorPanics(t *testing.T) {
	assertPanics(t, func() { common.NewErrorReason(nil) })
}

func TestErrorReason_SynopsisIsErrorMessage(t *testing.T) {
	var reason common.FailureError = common.NewErrorReason(errSome)
	if got := reason.Synopsis(); got != "Some message." {
		t.Errorf("Synopsis() = %q, want %q", got, "Some message.")
	}
	if got := reason.Error(); got != "Some message." {
		t.Errorf("Error() = %q, want %q", got, "Some message.")
	}
	if got := common.Description(reason); got != "Some message." {
		t.Errorf("Description() = %q, want %q", got, "Some message.")
	}
	if common.ReportIssue(reason) {
		t.Error("ReportIssue() = true, want false")
	}
}

func TestErrorReason_StackTrace(t *testing.T) {
	frames, ok := common.StackTrace(common.NewErrorReason(errSome))
	if !ok || len(frames) == 0 {
		t.Fatalf("StackTrace() = %v, %v, want frames", frames, ok)
	}
	if !strings.HasSuffix(frames[0].Function, ".TestErrorReason_StackTrace") {
		t.Errorf("innermost frame = %q, want the creating function", frames[0].Function)
	}
}

func TestErrorReason_Unwrap(t *testing.T) {
	reason := common.NewErrorReason(errSome)
	if !errors.Is(reason, errSome) {
		t.Error("errors.Is(reason, errSome) = false, want true")
	}
}

// A FailureError without a stack trace of its own has none.
func TestFailureError_WithoutStackTrace(t *testing.T) {
	var reason common.FailureError = newFailureError()
	if frames, ok := common.StackTrace(reason); ok || frames != nil {
		t.Errorf("StackTrace() = %v, %v, want none", frames, ok)
	}
}
