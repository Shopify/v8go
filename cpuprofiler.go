// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

/*
#include <stdlib.h>
#include "v8go.h"
*/
import "C"
import (
	"unsafe"
)

type CPUProfiler struct {
	ptr *C.CPUProfiler
	iso *Isolate
}

// CPUProfiler is used to control CPU profiling.
func NewCPUProfiler(iso *Isolate) *CPUProfiler {
	return &CPUProfiler{
		ptr: C.NewCPUProfiler(iso.ptr),
		iso: iso,
	}
}

// Dispose will dispose the profiler.
func (c *CPUProfiler) Dispose() {
	if c.ptr == nil {
		return
	}

	C.CPUProfilerDispose(c.ptr)
	c.ptr = nil
}

// Changes default CPU profiler sampling interval to the specified number
// of microseconds. Default interval is 1000us. This method must be called
// when there are no profiles being recorded.
func (c *CPUProfiler) SetSamplingInterval(us int) {
	C.CPUProfilerSetSamplingInterval(c.ptr, C.int32_t(us))
}

// StartProfiling starts collecting a CPU profile. Title may be an empty string. Several
// profiles may be collected at once. Attempts to start collecting several
// profiles with the same title are silently ignored.
func (c *CPUProfiler) StartProfiling(title string) {
	if c.ptr == nil || c.iso.ptr == nil {
		panic("profiler or isolate are nil")
	}

	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))

	C.CPUProfilerStartProfiling(c.ptr, tstr)
}

// Synchronously collect current stack sample in all profilers attached to
// the isolate. The call does not affect number of ticks recorded for
// the current top node.
func (c *CPUProfiler) CollectSample() {
	C.CPUProfilerCollectSample(c.ptr)
}

// Stops collecting CPU profile with a given title and returns it.
// If the title given is empty, finishes the last profile started.
func (c *CPUProfiler) StopProfiling(title string) *CPUProfile {
	if c.ptr == nil || c.iso.ptr == nil {
		panic("profiler or isolate are nil")
	}

	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))

	profile := C.CPUProfilerStopProfiling(c.ptr, tstr)

	return &CPUProfile{
		p:         profile,
		title:     C.GoString(profile.title),
		root:      newCPUProfileNode(profile.root, nil),
		startTime: timeUnixMicro(-int64(profile.startTime)),
		endTime:   timeUnixMicro(-int64(profile.endTime)),
	}
}

func newCPUProfileNode(node *C.CPUProfileNode, parent *CPUProfileNode) *CPUProfileNode {
	n := &CPUProfileNode{
		scriptResourceName: C.GoString(node.scriptResourceName),
		functionName:       C.GoString(node.functionName),
		lineNumber:         int(node.lineNumber),
		columnNumber:       int(node.columnNumber),
		parent:             parent,
	}

	if node.childrenCount > 0 {
		for _, child := range (*[1 << 28]*C.CPUProfileNode)(unsafe.Pointer(node.children))[:node.childrenCount:node.childrenCount] {
			n.children = append(n.children, newCPUProfileNode(child, n))
		}
	}

	return n
}
