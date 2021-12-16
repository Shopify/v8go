// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func CreateSnapshot(source, origin string) *StartupData {
	v8once.Do(func() {
		C.Init()
	})

	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cOrigin))

	sd := &StartupData{
		ptr: C.CreateSnapshot(cSource, cOrigin),
	}
	fmt.Printf("%+v\n", sd)
	fmt.Printf("%#v\n", sd.ptr)
	// fmt.Println(sd.ptr.data)
	// fmt.Println(sd.ptr.raw_size)
	return sd
}

type FunctionCodeHandling string

const (
	FunctionCodeHandlingKeep  FunctionCodeHandling = "kKeep"
	FunctionCodeHandlingClear FunctionCodeHandling = "kClear"
)

type StartupData struct {
	ptr C.StartupDataPtr
}
