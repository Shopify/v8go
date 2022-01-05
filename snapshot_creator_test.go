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
	data := v8.CreateSnapshot("function run() { return 1 };", "script.js", v8.FunctionCodeHandlingKlear)

	iso := v8.NewIsolateWithCreateParams(v8.CreateParams{StartupData: data})
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

	// Create another context from the same iso to validate it works again

	ctx2 := v8.NewContext(iso)
	defer ctx2.Close()

	runVal2, err := ctx2.Global().Get("run")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", runVal)
	fn2, err := runVal2.AsFunction()
	if err != nil {
		panic(err)
	}
	val2, err := fn2.Call(v8.Undefined(iso))
	if err != nil {
		panic(err)
	}
	if val2.String() != "1" {
		t.Fatal("invalid val")
	}
}
