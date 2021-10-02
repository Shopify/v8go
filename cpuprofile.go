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
}

// Returns the root node of the top down call tree.
func (c *CPUProfile) GetTopDownRoot() *CPUProfileNode {
	ptr := C.CpuProfileGetTopDownRoot(c.ptr)
	return &CPUProfileNode{ptr: ptr, iso: c.iso}
}

// Returns CPU profile title.
func (c *CPUProfile) GetTitle() string {
	str := C.CpuProfileGetTitle(c.iso.ptr, c.ptr)
	return C.GoString(str)
}

// Returns number of samples recorded. The samples are not recorded unless
// recordSamples parameter of CpuProfiler.StartCpuProfiling is true.
func (c *CPUProfile) GetSamplesCount() int {
	i := C.CpuProfileGetSamplesCount(c.ptr)
	return int(i)
}

func (c *CPUProfile) Delete() {
	if c.ptr == nil {
		return
	}
	C.CpuProfileDelete(c.ptr)
	c.ptr = nil
}
