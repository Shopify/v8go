// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCPUProfiler_Dispose(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	cpuProfiler := v8.NewCPUProfiler(iso)

	cpuProfiler.Dispose()
	// noop when called multiple times
	cpuProfiler.Dispose()

	// verify panics when profiler disposed
	if recoverPanic(func() { cpuProfiler.StartProfiling("") }) == nil {
		t.Error("expected panic")
	}

	if recoverPanic(func() { cpuProfiler.StopProfiling("") }) == nil {
		t.Error("expected panic")
	}

	cpuProfiler = v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()
	iso.Dispose()

	// verify panics when isolate disposed
	if recoverPanic(func() { cpuProfiler.StartProfiling("") }) == nil {
		t.Error("expected panic")
	}

	if recoverPanic(func() { cpuProfiler.StopProfiling("") }) == nil {
		t.Error("expected panic")
	}
}

func TestCPUProfiler_Sampling(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	// Force sampling interval to be large so we know the sample is from CollectSample call
	cpuProfiler.SetSamplingInterval(10000000)

	title := "cpuprofilersamplingtest"
	cpuProfiler.StartProfiling(title)

	foo := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		cpuProfiler.CollectSample()
		return nil
	})
	global := v8.NewObjectTemplate(iso)
	global.Set("foo", foo)
	ctx := v8.NewContext(iso, global)
	ctx.RunScript("function one() { foo() };", "")
	val, err := ctx.Global().Get("one")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global())
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling(title)
	defer cpuProfile.Delete()

	root := cpuProfile.GetTopDownRoot()
	if root == nil {
		t.Fatal("expected profile top down root not to be nil")
	}
	if findChild(t, root, "one") == nil {
		t.Fatal("expected sample to capture node from script")
	}
}

func TestCPUProfiler(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	title := "cpuprofilertest"
	cpuProfiler.StartProfiling(title)

	_, err := ctx.RunScript("function noop() {}", "script.js")
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling(title)
	defer cpuProfile.Delete()

	if cpuProfile.GetTitle() != title {
		t.Errorf("expected %s, but got %s", title, cpuProfile.GetTitle())
	}
	if cpuProfile.GetTopDownRoot() == nil {
		t.Fatal("expected profile top down root not to be nil")
	}
}

const profileScript = `function loop(timeout) {
  this.mmm = 0;
  var start = Date.now();
  while (Date.now() - start < timeout) {
    var n = 10;
    while(n > 1) {
      n--;
      this.mmm += n * n * n;
    }
  }
}
function delay() { try { loop(10); } catch(e) { } }
function bar() { delay(); }
function baz() { delay(); }
function foo() {
    try {
       delay();
       bar();
       delay();
       baz();
    } catch (e) { }
}
function start(timeout) {
  var start = Date.now();
  do {
    foo();
    var duration = Date.now() - start;
  } while (duration < timeout);
  return duration;
};`
