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
	err := snapshotCreator.AddContext("function run() { return 1 };", "script.js")
	fatalIf(t, err)

	data, err := snapshotCreator.Create(v8.FunctionCodeHandlingKlear)
	fatalIf(t, err)
	defer data.Dispose()

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()

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

	err := snapshotCreator.AddContext("function run() { return 1 };", "script.js")
	fatalIf(t, err)

	data, err := snapshotCreator.Create(v8.FunctionCodeHandlingKlear)
	fatalIf(t, err)
	defer data.Dispose()

	_, err = snapshotCreator.Create(v8.FunctionCodeHandlingKlear)
	if err == nil {
		t.Error("Creating snapshot should have fail")
	}
}

func TestAddContextFail(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	defer snapshotCreator.Dispose()

	err := snapshotCreator.AddContext("feuihyvfeuyfeu", "script.js")
	if err == nil {
		t.Error("add context should have fail")
	}
}

func TestCreateSnapshotFailAndReuse(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	err := snapshotCreator.AddContext("feuihyvfeuyfeu", "script.js")
	if err == nil {
		t.Error("add context should have fail")
	}
	err = snapshotCreator.AddContext("function run() { return 1 };", "script.js")
	data, err := snapshotCreator.Create(v8.FunctionCodeHandlingKlear)
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
