// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
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

func (s *StartupData) Dispose() {
	if s.ptr != nil {
		C.SnapshotBlobDelete(s.ptr)
	}
}

type SnapshotCreator struct {
	ptr C.SnapshotCreatorPtr
}

func NewSnapshotCreator() *SnapshotCreator {
	v8once.Do(func() {
		C.Init()
	})

	return &SnapshotCreator{
		ptr: C.NewSnapshotCreator(),
	}
}

func (s *SnapshotCreator) Create(source, origin string, functionCode FunctionCodeHandling) (*StartupData, error) {
	if s.ptr == nil {
		return nil, errors.New("v8go: Cannot use snapshot creator after creating the blob")
	}

	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cOrigin))

	rtn := C.CreateBlob(s.ptr, cSource, cOrigin, C.int(functionCode))

	if rtn.blob == nil {
		return nil, newJSError(rtn.error)
	}

	s.ptr = nil

	return &StartupData{ptr: rtn.blob}, nil
}

func (s *SnapshotCreator) Dispose() {
	if s.ptr != nil {
		C.DeleteSnapshotCreator(s.ptr)
	}
}
