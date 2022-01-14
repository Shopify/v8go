// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

// func TestCreateSnapshot(t *testing.T) {
// 	data, err := v8.CreateSnapshot("function run() { return 1 };", "script.js", v8.FunctionCodeHandlingKlear)
// 	fatalIf(t, err)

// 	iso := v8.NewIsolate(v8.WithStartupData(data))
// 	defer iso.Dispose()
// 	defer data.Dispose(iso)
// 	ctx := v8.NewContext(iso)
// 	defer ctx.Close()

// 	runVal, err := ctx.Global().Get("run")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fn, err := runVal.AsFunction()
// 	if err != nil {
// 		panic(err)
// 	}
// 	val, err := fn.Call(v8.Undefined(iso))
// 	if err != nil {
// 		panic(err)
// 	}
// 	if val.String() != "1" {
// 		t.Fatal("invalid val")
// 	}

// 	// Create another context from the same iso to validate it works again

// 	ctx2 := v8.NewContext(iso)
// 	defer ctx2.Close()

// 	runVal2, err := ctx2.Global().Get("run")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fn2, err := runVal2.AsFunction()
// 	if err != nil {
// 		panic(err)
// 	}
// 	val2, err := fn2.Call(v8.Undefined(iso))
// 	if err != nil {
// 		panic(err)
// 	}
// 	if val2.String() != "1" {
// 		t.Fatal("invalid val")
// 	}
// }

// func TestCreateSnapshotFail(t *testing.T) {
// 	_, err := v8.CreateSnapshot("uidygwuiwgduw", "script.js", v8.FunctionCodeHandlingKlear)
// 	if err == nil {
// 		t.Error("Creating snapshot should have fail")
// 	}
// }

func TestCreateSnapshotV1(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()

	data, err := snapshotCreator.Create("function run() { return 1 };", "script.js", v8.FunctionCodeHandlingKlear)

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()
	// defer snapshotCreator.Dispose(iso)

	ctx := v8.NewContext(iso)
	defer ctx.Close()

	runVal, err := ctx.Global().Get("run")
	if err != nil {
		panic(err)
	}

	fn, err := runVal.AsFunction()
	if err != nil {
		panic(err)
	}
	val, err := fn.Call(v8.Undefined(iso))
	if err != nil {
		panic(err)
	}
	if val.String() != "1" {
		t.Fatal("invalid val")
	}
}

func TestCreateSnapshotV2(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()

	scripts := []string{
		"function run() { return 1 };",
		"function run2() { return 2 };",
	}

	data, err := snapshotCreator.CreateV2(scripts, v8.FunctionCodeHandlingKlear)

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()
	// defer snapshotCreator.Dispose(iso)

	ctx := v8.NewContext(iso)
	defer ctx.Close()

	runVal, err := ctx.Global().Get("run")
	if err != nil {
		panic(err)
	}

	fn, err := runVal.AsFunction()
	if err != nil {
		panic(err)
	}
	val, err := fn.Call(v8.Undefined(iso))
	if err != nil {
		panic(err)
	}
	if val.String() != "1" {
		t.Fatal("invalid val")
	}

	runVal2, err := ctx.Global().Get("run2")
	if err != nil {
		panic(err)
	}

	fn2, err := runVal2.AsFunction()
	if err != nil {
		panic(err)
	}
	val, err = fn2.Call(v8.Undefined(iso))
	if err != nil {
		panic(err)
	}
	if val.String() != "2" {
		t.Fatal("invalid val")
	}
}
