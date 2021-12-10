// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

type SnapshotCreator struct {
	ptr C.SnapshotCreatorPtr
	iso *Isolate
}

func NewSnapshotCreator(iso *Isolate) *SnapshotCreator {
	ptr := C.NewSnapshotCreator()
	return &SnapshotCreator{
		ptr: ptr,
		iso: iso,
	}
}

// TODO: Delete snapshot creator will delete associated iso too

// func (s *SnapshotCreator) GetIsolate() *Isolate {
// 	isoptr :=	C.SnapshotCreatorGetIsolate(s.ptr)
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

func (s SnapshotCreator) CreateBlob(fch FunctionCodeHandling) *StartupData {
	ptr := C.SnapshotCreatorCreateBlob(s.ptr)
	return &StartupData{ptr: ptr}
}
