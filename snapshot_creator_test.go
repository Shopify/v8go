// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCreateSnapshot(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()

	data, err := snapshotCreator.Create("function run() { return 1 };", "script.js", v8.FunctionCodeHandlingKlear)
	fatalIf(t, err)

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()
	defer data.Dispose()

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

func TestCreateSnapshotErrorAfterSuccessfullCreate(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()

	data, err := snapshotCreator.Create("function run() { return 1 };", "script.js", v8.FunctionCodeHandlingKlear)
	defer data.Dispose()
	fatalIf(t, err)

	_, err = snapshotCreator.Create("function run2() { return 2 };", "script2.js", v8.FunctionCodeHandlingKlear)
	if err == nil {
		t.Error("Creating snapshot should have fail")
	}
}

func TestCreateSnapshotFail(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	defer snapshotCreator.Dispose()

	_, err := snapshotCreator.Create("uidygwuiwgduw", "script.js", v8.FunctionCodeHandlingKlear)
	if err == nil {
		t.Error("Creating snapshot should have fail")
	}
}

func TestCreateSnapshotFailAndReuse(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	_, err := snapshotCreator.Create("uidygwuiwgduw", "script.js", v8.FunctionCodeHandlingKlear)
	if err == nil {
		t.Error("Creating snapshot should have fail")
	}

	data, err := snapshotCreator.Create("function run() { return 1 };", "script.js", v8.FunctionCodeHandlingKlear)
	fatalIf(t, err)

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()
	defer data.Dispose()

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

func TestCreateSnapshotWithIsolateOption(t *testing.T) {
	iso1 := v8.NewIsolate()
	defer iso1.Dispose()
	snapshotCreator := v8.NewSnapshotCreator(v8.WithIsolate(iso1))

	data, err := snapshotCreator.Create("function run() { return 1 };", "script.js", v8.FunctionCodeHandlingKlear)
	fatalIf(t, err)

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()
	defer data.Dispose()

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
