// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCreateSnapshot(t *testing.T) {
	data := v8.CreateSnapshot("function run() { return 1 };", "script.js")

	iso2 := v8.NewIsolateWithCreateParams(v8.CreateParams{SnapshotBlob: data})
	ctx2 := v8.NewContext(iso2)
	val, err := ctx2.RunScript("run()", "script.js")
	fatalIf(t, err)
	if val.String() != "1" {
		t.Fatal("invalid val")
	}
}
