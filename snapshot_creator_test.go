// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"fmt"
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCreateSnapshot(t *testing.T) {
	data := v8.CreateSnapshot("function run() { return 1 };", "script.js")

	iso := v8.NewIsolateWithCreateParams(v8.CreateParams{SnapshotBlob: data})
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()

	runVal, err := ctx.Global().Get("run")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", runVal)
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
