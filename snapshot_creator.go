// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import "unsafe"

type FunctionCodeHandling int

const (
	FunctionCodeHandlingKlear FunctionCodeHandling = iota
	FunctionCodeHandlingKeep
)

type StartupData struct {
	ptr *C.SnapshotBlob
}

// func CreateSnapshot(source, origin string, functionCode FunctionCodeHandling) (*StartupData, error) {
// 	v8once.Do(func() {
// 		C.Init()
// 	})

// 	cSource := C.CString(source)
// 	cOrigin := C.CString(origin)
// 	defer C.free(unsafe.Pointer(cSource))
// 	defer C.free(unsafe.Pointer(cOrigin))

// 	rtn := C.CreateSnapshot(cSource, cOrigin, C.int(functionCode))

// 	if rtn.blob == nil {
// 		return nil, newJSError(rtn.error)
// 	}

// 	return &StartupData{
// 		ptr: rtn.blob,
// 	}, nil
// }

// func (s *StartupData) Dispose(iso *Isolate) {
// 	C.SnapshotBlobDelete(iso.ptr, s.ptr)
// }

type SnapshotCreator struct {
	ptr     C.SnapshotCreatorPtr
	dataPtr *C.SnapshotBlob
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
		panic("Cannot use snapshot creator after creating the blob")
	}

	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cOrigin))

	rtn := C.CreateSnapshotV2(s.ptr, cSource, cOrigin, C.int(functionCode))
	s.ptr = nil

	if rtn.blob == nil {
		return nil, newJSError(rtn.error)
	}

	s.dataPtr = rtn.blob

	return &StartupData{
		ptr: rtn.blob,
	}, nil
}

func (s *SnapshotCreator) CreateV2(scripts []string, functionCode FunctionCodeHandling) (*StartupData, error) {
	if s.ptr == nil {
		panic("Cannot use snapshot creator after creating the blob")
	}

	charArray := make([]*C.char, len(scripts))
	for i, s := range scripts {
		charArray[i] = C.CString(s)
	}
	cOrigin := C.CString("<embedded>")

	rtn := C.CreateSnapshotV3(s.ptr, &charArray[0], C.int(len(scripts)), cOrigin, C.int(functionCode))
	s.ptr = nil

	for _, s := range charArray {
		C.free(unsafe.Pointer(s))
	}
	defer C.free(unsafe.Pointer(cOrigin))

	if rtn.blob == nil {
		return nil, newJSError(rtn.error)
	}

	s.dataPtr = rtn.blob

	return &StartupData{
		ptr: rtn.blob,
	}, nil
}

func (s *SnapshotCreator) Dispose(iso *Isolate) {
	C.SnapshotBlobDelete(iso.ptr, s.dataPtr)
}
