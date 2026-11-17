// Copyright (c) 2021-2026 - for information on the respective copyright owner
// see the NOTICE file and/or the repository https://github.com/carbynestack/common.
//
// SPDX-License-Identifier: Apache-2.0

package common

import "testing"

func TestStack_EmptyHasNoFrames(t *testing.T) {
	if frames := stack(nil).frames(); frames != nil {
		t.Errorf("frames() = %v, want nil", frames)
	}
	if _, ok := StackTrace(&Unexpected{}); ok {
		t.Error("StackTrace() of a zero Unexpected reports a stack trace")
	}
}
