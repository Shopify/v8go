// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCreateSnapshot(t *testing.T) {
	// TODO: This is needed to run C.Init()
	iso := v8.NewIsolate()
	defer iso.Dispose()

	data := v8.CreateSnapshot("function run() { return 1 };", "script.js")

	iso2 := v8.NewIsolateWithCreateParams(v8.CreateParams{SnapshotBlob: data})
	ctx2 := v8.NewContext(iso2)
	val, err := ctx2.RunScript("run()", "script.js")
	fatalIf(t, err)
	if val.String() != "1" {
		t.Fatal("invalid val")
	}
}

func TestSnapshotCreator(t *testing.T) {
	t.Parallel()

	snapshotCreator := v8.NewSnapshotCreator()
	iso1 := snapshotCreator.GetIsolate()
	ctx1 := v8.NewContext(iso1)
	_, err := ctx1.RunScript("function run() { return 1 };", "script.js")
	fatalIf(t, err)

	data := snapshotCreator.CreateBlob(v8.FunctionCodeHandlingKeep)

	iso2 := v8.NewIsolateWithCreateParams(v8.CreateParams{SnapshotBlob: data})
	ctx2 := v8.NewContext(iso2)
	val, err := ctx2.RunScript("run()", "script.js")
	fatalIf(t, err)
	if val.String() != "1" {
		t.Fatal("invalid val")
	}
}
