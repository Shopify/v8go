// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCPUProfile_Dispose(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	cpuProfiler := v8.NewCPUProfiler(iso)
	cpuProfiler.Dispose()
	// noop when called multiple times
	cpuProfiler.Dispose()
}

func TestCPUProfile(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofiletest", true)

	_, err := ctx.RunScript(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global())
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofiletest")
	if cpuProfile == nil {
		t.Fatal("expected profiler not to be nil")
	}
	defer cpuProfile.Delete()

	if cpuProfile.GetTitle() != "cpuprofiletest" {
		t.Errorf("expected cpuprofiletest, but got %v", cpuProfile.GetTitle())
	}

	if cpuProfile.GetSamplesCount() > 0 {
		t.Errorf("expected samples count greater than 0, but got %d", cpuProfile.GetSamplesCount())
	}

	root := cpuProfile.GetTopDownRoot()
	if root == nil {
		t.Fatal("expected root not to be nil")
	}
	if root.GetFunctionName() != "(root)" {
		t.Errorf("expected (root), but got %v", root.GetFunctionName())
	}
}
