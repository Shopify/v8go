// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

// CPUProfile contains a CPU profile in a form of top-down call tree
// (from main() down to functions that do all the work).
type CPUProfile struct {
	ptr C.CpuProfilePtr
	iso *Isolate

	title       string
	topDownRoot *CPUProfileNode
}

// Returns the root node of the top down call tree.
func (c *CPUProfile) GetTopDownRoot() *CPUProfileNode {
	return c.topDownRoot
}

// Returns CPU profile title.
func (c *CPUProfile) GetTitle() string {
	return c.title
}

func (c *CPUProfile) Delete() {
	if c.ptr == nil {
		return
	}
	C.CpuProfileDelete(c.ptr)
	c.ptr = nil
}
