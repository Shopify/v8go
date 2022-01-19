// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
)

type FunctionCodeHandling int

const (
	FunctionCodeHandlingKlear FunctionCodeHandling = iota
	FunctionCodeHandlingKeep
)

type StartupData struct {
	data     *C.char
	raw_size C.int
	index    C.size_t
}

type SnapshotCreator struct {
	ptr   C.SnapshotCreatorPtr
	iso   *Isolate
	ctx   *Context
	index C.size_t
}

func NewSnapshotCreator() *SnapshotCreator {
	v8once.Do(func() {
		C.Init()
	})

	rtn := C.NewSnapshotCreator()

	return &SnapshotCreator{
		ptr: rtn.creator,
		iso: &Isolate{ptr: rtn.iso},
	}
}

func (s *SnapshotCreator) GetIsolate() *Isolate {
	return s.iso
}

func (s *SnapshotCreator) AddContext(ctx *Context) error {
	if s.ptr == nil {
		return errors.New("v8go: Cannot add context to snapshot creator after creating the blob")
	}

	s.index = C.AddContext(s.ptr, ctx.ptr)
	s.ctx = ctx

	return nil
}

func (s *SnapshotCreator) Create(functionCode FunctionCodeHandling) (*StartupData, error) {
	if s.ptr == nil {
		return nil, errors.New("v8go: Cannot use snapshot creator after creating the blob")
	}

	rtn := C.CreateBlob(s.ptr, s.ctx.ptr, C.int(functionCode))

	s.ptr = nil
	s.ctx.ptr = nil
	s.iso.ptr = nil

	raw_size := rtn.raw_size
	// data := C.GoBytes(unsafe.Pointer(rtn.data), raw_size)
	// defer C.free(unsafe.Pointer(rtn.data))

	return &StartupData{data: rtn.data, raw_size: raw_size, index: s.index}, nil
}

func (s *SnapshotCreator) Dispose() {
	if s.ptr != nil {
		C.DeleteSnapshotCreator(s.ptr)
	}
}
