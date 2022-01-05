// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"unsafe"
)

type FunctionCodeHandling int

const (
	FunctionCodeHandlingKlear FunctionCodeHandling = iota
	FunctionCodeHandlingKeep
)

type StartupData struct {
	ptr *C.SnapshotBlob
}

func CreateSnapshot(source, origin string, functionCode FunctionCodeHandling) *StartupData {
	v8once.Do(func() {
		C.Init()
	})

	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cOrigin))

	return &StartupData{
		ptr: C.CreateSnapshot(cSource, cOrigin, C.int(functionCode)),
	}
}

func (s *StartupData) Dispose() {
	C.SnapshotBlobDelete(s.ptr)
}
