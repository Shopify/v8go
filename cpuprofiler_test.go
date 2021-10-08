// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestCPUProfilerDispose(t *testing.T) {
	t.Parallel()

	iso := v8go.NewIsolate()
	defer iso.Dispose()
	cpuProfiler := v8go.NewCPUProfiler(iso)

	cpuProfiler.Dispose()
	// noop when called multiple times
	cpuProfiler.Dispose()

	// verify does not panic once disposed
	cpuProfiler.StartProfiling("", false)
	cpuProfiler.StopProfiling("")

	cpuProfiler = v8go.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()
	iso.Dispose()
	// verify does not panic once isolate disposed
	cpuProfiler.StartProfiling("", false)
	cpuProfiler.StopProfiling("")
}

func TestCPUProfiler(t *testing.T) {
	t.Parallel()

	ctx := v8go.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8go.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofilertest", true)

	_, err := ctx.RunScript(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global())
	fatalIf(t, err)

	// time.Sleep(100 * time.Millisecond)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofilertest")
	if cpuProfile == nil {
		t.Fatal("expected profiler not to be nil")
	}
	defer cpuProfile.Delete()

	if cpuProfile.GetTitle() != "cpuprofilertest" {
		t.Fatalf("expected cpu profile to be %s, but got %s", "cpuprofilertest", cpuProfile.GetTitle())
	}
	// if cpuProfile.GetSamplesCount() == 0 {
	// 	t.Fatalf("expected cpu profile to have samples count greater than 0, but got %d", cpuProfile.GetSamplesCount())
	// }

	root := cpuProfile.GetTopDownRoot()
	if root == nil {
		t.Fatal("expected root not to be nil")
	}
	if root.GetFunctionName() != "(root)" {
		t.Fatalf("expected (root), but got %v", root.GetFunctionName())
	}
	checkNode(t, root, "(root)", 0, 0)
	checkChildren(t, root, []string{"(program)", "start", "(garbage collector)"})

	start := root.GetChild(1)
	checkNode(t, start, "start", 23, 15)
	checkChildren(t, start, []string{"foo"})

	foo := start.GetChild(0)
	checkNode(t, foo, "foo", 15, 13)
	checkChildren(t, foo, []string{"delay", "bar", "baz"})

	baz := foo.GetChild(2)
	checkNode(t, baz, "baz", 14, 13)
	checkChildren(t, baz, []string{"delay"})

	delay := baz.GetChild(0)
	checkNode(t, delay, "delay", 12, 15)
	checkChildren(t, delay, []string{"loop"})
}

func checkChildren(t *testing.T, node *v8go.CPUProfileNode, names []string) {
	nodeName := node.GetFunctionName()
	if node.GetChildrenCount() != len(names) {
		present := []string{}
		for i := 0; i < node.GetChildrenCount(); i++ {
			present = append(present, node.GetChild(i).GetFunctionName())
		}
		t.Fatalf("child count for node %s should be %d but was %d: %v", nodeName, len(names), node.GetChildrenCount(), present)
	}
	for i, n := range names {
		if node.GetChild(i).GetFunctionName() != n {
			t.Fatalf("expected %s child %d to have name %s", nodeName, i, n)
		}
	}
}

func checkNode(t *testing.T, node *v8go.CPUProfileNode, name string, line, column int) {
	if node.GetFunctionName() != name {
		t.Fatalf("expected node to have function name `%s` but had `%s`", name, node.GetFunctionName())
	}
	if node.GetLineNumber() != line {
		t.Fatalf("expected node %s at line %d, but got %d", name, line, node.GetLineNumber())
	}
	if node.GetColumnNumber() != column {
		t.Fatalf("expected node %s at column %d, but got %d", name, column, node.GetColumnNumber())
	}
}

// const profileTree = `
// [Top down]:
//  1062     0   (root) [-1]
//  1054     0    start [-1]
//  1054     1      foo [-1]
//   265     0        baz [-1]
//   265     1          delay [-1]
//   264   264            loop [-1]
//   525     3        delay [-1]
//   522   522          loop [-1]
//   263     0        bar [-1]
//   263     1          delay [-1]
//   262   262            loop [-1]
//     2     2    (program) [-1]
//     6     6    (garbage collector) [-1]
// `

const profileScript = `function loop(timeout) {
  this.mmm = 0;
  var start = Date.now();
  while (Date.now() - start < timeout) {
    var n = 100;
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
