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
	snapshotCreatorIso := snapshotCreator.GetIsolate()
	snapshotCreatoCtx := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatoCtx.Close()

	snapshotCreatoCtx.RunScript(`const add = (a, b) => a + b`, "add.js")
	snapshotCreatoCtx.RunScript(`function run() { return add(3, 4); }`, "main.js")
	err := snapshotCreator.AddContext(snapshotCreatoCtx)
	fatalIf(t, err)

	data, err := snapshotCreator.Create(v8.FunctionCodeHandlingKlear)
	fatalIf(t, err)

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
	if val.String() != "7" {
		t.Fatal("invalid val")
	}
}

func TestCreateSnapshotErrorAfterSuccessfullCreate(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	snapshotCreatorIso := snapshotCreator.GetIsolate()
	snapshotCreatoCtx := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatoCtx.Close()

	snapshotCreatoCtx.RunScript(`const add = (a, b) => a + b`, "add.js")
	snapshotCreatoCtx.RunScript(`function run() { return add(3, 4); }`, "main.js")
	err := snapshotCreator.AddContext(snapshotCreatoCtx)
	fatalIf(t, err)

	_, err = snapshotCreator.Create(v8.FunctionCodeHandlingKlear)
	fatalIf(t, err)

	err = snapshotCreator.AddContext(snapshotCreatoCtx)
	if err == nil {
		t.Error("Adding context should have fail")
	}

	_, err = snapshotCreator.Create(v8.FunctionCodeHandlingKlear)
	if err == nil {
		t.Error("Creating snapshot should have fail")
	}
}
