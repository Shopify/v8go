// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCPUProfileNode(t *testing.T) {
	t.Parallel()

	child := v8.NewCPUProfileNode("child", 5, 6, []*v8.CPUProfileNode{})
	node := v8.NewCPUProfileNode("start", 1, 2, []*v8.CPUProfileNode{child})

	if node.GetFunctionName() != "start" {
		t.Fatalf("expected start but got %s", node.GetFunctionName())
	}

	if node.GetLineNumber() != 1 {
		t.Fatalf("expected 1 but got %d", node.GetLineNumber())
	}

	if node.GetColumnNumber() != 2 {
		t.Fatalf("expected 2 but got %d", node.GetColumnNumber())
	}

	if node.GetChildrenCount() != 1 {
		t.Fatalf("expected 1 but got %d", node.GetChildrenCount())
	}

	if node.GetChild(0) != child {
		t.Fatalf("expected child but got %v", node.GetChild(0))
	}
}
