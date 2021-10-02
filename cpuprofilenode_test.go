// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestCPUProfileNode(t *testing.T) {
	t.Parallel()

	ctx := v8go.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8go.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofilenodetest", false)

	_, err := ctx.RunScript(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global())
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofilenodetest")
	if cpuProfile == nil {
		t.Fatal("expected profile not to be nil")
	}
	defer cpuProfile.Delete()

	rootNode := cpuProfile.GetTopDownRoot()
	if rootNode == nil {
		t.Fatal("expected top down root not to be nil")
	}
	if rootNode.GetFunctionName() != "(root)" {
		t.Fatalf("expected (root), but got %v", rootNode.GetFunctionName())
	}
	checkChildren(t, rootNode, []string{"(program)", "start", "(garbage collector)"})

	invalidChild := rootNode.GetChild(4)
	if invalidChild != nil {
		t.Fatalf("expected nil child, but got %v", invalidChild.GetFunctionName())
	}

	startNode := rootNode.GetChild(1)
	if startNode.GetFunctionName() != "start" {
		t.Fatalf("expected start, but got %v", startNode.GetFunctionName())
	}
	checkChildren(t, startNode, []string{"foo"})
	checkPosition(t, startNode, 23, 15)

	parentName := startNode.GetParent().GetFunctionName()
	if parentName != "(root)" {
		t.Fatalf("expected (root), but got %v", parentName)
	}

	fooNode := startNode.GetChild(0)
	checkChildren(t, fooNode, []string{"delay", "bar", "baz"})
	checkPosition(t, fooNode, 15, 13)

	delayNode := fooNode.GetChild(0)
	checkChildren(t, delayNode, []string{"loop"})
	checkPosition(t, delayNode, 12, 15)

	barNode := fooNode.GetChild(1)
	checkChildren(t, barNode, []string{"delay"})

	bazNode := fooNode.GetChild(2)
	checkChildren(t, bazNode, []string{"delay"})
}

func checkChildren(t *testing.T, node *v8go.CPUProfileNode, names []string) {
	nodeName := node.GetFunctionName()
	if node.GetChildrenCount() != len(names) {
		t.Fatalf("expected child count for node %s to equal length of child names", nodeName)
	}
	for i, n := range names {
		if node.GetChild(i).GetFunctionName() != n {
			t.Fatalf("expected %s child %d to have name %s", nodeName, i, n)
		}
	}
}

func checkPosition(t *testing.T, node *v8go.CPUProfileNode, line, column int) {
	nodeName := node.GetFunctionName()
	if node.GetLineNumber() != line {
		t.Fatalf("expected node %s at line %d, but got %d", nodeName, line, node.GetLineNumber())
	}
	if node.GetColumnNumber() != column {
		t.Fatalf("expected node %s at column %d, but got %d", nodeName, column, node.GetColumnNumber())
	}
}
