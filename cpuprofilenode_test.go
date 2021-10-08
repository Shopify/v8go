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

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofilenodetest")

	_, _ = ctx.RunScript(profileScript, "")
	val, _ := ctx.Global().Get("start")
	fn, _ := val.AsFunction()
	_, _ = fn.Call(ctx.Global())

	cpuProfile := cpuProfiler.StopProfiling("cpuprofilenodetest")
	defer cpuProfile.Delete()

	node := cpuProfile.GetTopDownRoot()

	if node.GetFunctionName() != "(root)" {
		t.Fatalf("expected start but got %s", node.GetFunctionName())
	}

	if node.GetLineNumber() != 0 {
		t.Fatalf("expected 0 but got %d", node.GetLineNumber())
	}

	if node.GetColumnNumber() != 0 {
		t.Fatalf("expected 0 but got %d", node.GetColumnNumber())
	}

	if node.GetChildrenCount() < 2 {
		t.Fatalf("expected at least 2 children, but got %d", node.GetChildrenCount())
	}

	if node.GetChild(1).GetFunctionName() != "start" {
		t.Fatalf("expected child node with name `start` but got %s", node.GetChild(1).GetFunctionName())
	}
}
