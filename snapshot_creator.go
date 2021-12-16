// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import "unsafe"

func CreateSnapshot(source, origin string) *StartupData {
	v8once.Do(func() {
		C.Init()
	})

	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cOrigin))

	ptr := C.CreateSnapshot(cSource, cOrigin)
	return &StartupData{ptr: ptr}
}

type SnapshotCreator struct {
	ptr C.SnapshotCreatorPtr
	iso *Isolate
}

// func NewSnapshotCreator() *SnapshotCreator {
// 	v8once.Do(func() {
// 		C.Init()
// 	})

// 	wrap := C.NewSnapshotCreator()

// 	iso := &Isolate{
// 		ptr: wrap.iso,
// 		cbs: make(map[int]FunctionCallback),
// 	}
// 	iso.null = newValueNull(iso)
// 	iso.undefined = newValueUndefined(iso)

// 	return &SnapshotCreator{
// 		ptr: wrap.ptr,
// 		iso: iso,
// 	}
// }

// TODO: Delete snapshot creator will delete associated iso too

// func (s *SnapshotCreator) GetIsolate() *Isolate {
// 	return s.iso
// }

type FunctionCodeHandling string

const (
	FunctionCodeHandlingKeep  FunctionCodeHandling = "kKeep"
	FunctionCodeHandlingClear FunctionCodeHandling = "kClear"
)

type StartupData struct {
	ptr C.StartupDataPtr
}

// func (s SnapshotCreator) CreateBlob(fch FunctionCodeHandling) *StartupData {
// 	ptr := C.SnapshotCreatorCreateBlob(s.ptr)
// 	return &StartupData{ptr: ptr}
// }
