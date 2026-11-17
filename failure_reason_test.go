// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common_test

import (
	"runtime"
	"testing"

	"github.com/Samrat7745/common"
	"github.com/Samrat7745/common/commontest"
)

// testReason implements only the required method of a failure reason.
type testReason struct{ synopsis string }

func (r testReason) Synopsis() string { return r.synopsis }

// fullReason implements all optional interfaces of a failure reason.
type fullReason struct{ frames []runtime.Frame }

func (fullReason) Synopsis() string              { return "synopsis" }
func (fullReason) Description() string           { return "description" }
func (r fullReason) StackTrace() []runtime.Frame { return r.frames }
func (fullReason) ReportIssue() bool             { return true }

func TestFailureReason_DescriptionDefaultsToSynopsis(t *testing.T) {
	if got := common.Description(testReason{"synopsis"}); got != "synopsis" {
		t.Errorf("Description() = %q, want %q", got, "synopsis")
	}
}

func TestFailureReason_StackTraceDefaultsToNone(t *testing.T) {
	if frames, ok := common.StackTrace(testReason{"synopsis"}); ok || frames != nil {
		t.Errorf("StackTrace() = %v, %v, want none", frames, ok)
	}
}

func TestFailureReason_ReportIssueDefaultsToFalse(t *testing.T) {
	if common.ReportIssue(testReason{"synopsis"}) {
		t.Error("ReportIssue() = true, want false")
	}
}

func TestFailureReason_ToFailure(t *testing.T) {
	reason := testReason{"synopsis"}
	commontest.AssertThat(t, common.ToFailure[int](reason)).IsFailure().HasReason(reason)

	described := newFailureError()
	commontest.AssertThat(t, common.ToFailure[string](described)).HasReason(described)
}

func TestFailureReason_OptionalInterfacesAreUsed(t *testing.T) {
	reason := fullReason{frames: []runtime.Frame{{Function: "f"}}}
	var _ common.Describable = reason
	if got := common.Description(reason); got != "description" {
		t.Errorf("Description() = %q, want %q", got, "description")
	}
	if frames, ok := common.StackTrace(reason); !ok || len(frames) != 1 {
		t.Errorf("StackTrace() = %v, %v, want one frame", frames, ok)
	}
	if !common.ReportIssue(reason) {
		t.Error("ReportIssue() = false, want true")
	}
}

func TestFailureReason_NilStackTraceIsNone(t *testing.T) {
	if frames, ok := common.StackTrace(fullReason{}); ok || frames != nil {
		t.Errorf("StackTrace() = %v, %v, want none", frames, ok)
	}
}
