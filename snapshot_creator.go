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

type snapshotCreatorOptions struct {
	iso         *Isolate
	exitingBlob *StartupData
}

type creatorOptions func(*snapshotCreatorOptions)

func WithIsolate(iso *Isolate) creatorOptions {
	return func(options *snapshotCreatorOptions) {
		options.iso = iso
	}
}

type SnapshotCreator struct {
	ptr C.SnapshotCreatorPtr
	*StartupData
	*snapshotCreatorOptions
}

func NewSnapshotCreator(opts ...creatorOptions) *SnapshotCreator {
	v8once.Do(func() {
		C.Init()
	})

	options := &snapshotCreatorOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var cOptions C.SnapshotCreatorOptions

	if options.iso != nil {
		cOptions.iso = options.iso.ptr
	}

	return &SnapshotCreator{
		ptr:                    C.NewSnapshotCreator(cOptions),
		snapshotCreatorOptions: options,
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

	rtn := C.CreateSnapshot(s.ptr, cSource, cOrigin, C.int(functionCode))

	if rtn.blob == nil {
		return nil, newJSError(rtn.error)
	}

	s.ptr = nil

	if s.snapshotCreatorOptions.iso != nil {
		s.snapshotCreatorOptions.iso.ptr = nil
	}

	startupData := &StartupData{ptr: rtn.blob}
	s.StartupData = startupData

	return startupData, nil
}

func (s *SnapshotCreator) Dispose() {
	if s.ptr != nil {
		C.DeleteSnapshotCreator(s.ptr)
	}
	if s.StartupData != nil {
		C.SnapshotBlobDelete(s.StartupData.ptr)
	}
}
