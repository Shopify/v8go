package v8go

/*
#include "v8go.h"
*/
import "C"

type CPUProfile struct {
	p *C.CPUProfile

	// The CPU profile title.
	title string

	// root is the root node of the top down call tree.
	root *CPUProfileNode
}

// Returns CPU profile title.
func (c *CPUProfile) GetTitle() string {
	return c.title
}

// Returns the root node of the top down call tree.
func (c *CPUProfile) GetTopDownRoot() *CPUProfileNode {
	return c.root
}

// Deletes the profile and removes it from CpuProfiler's list.
// All pointers to nodes previously returned become invalid.
func (c *CPUProfile) Delete() {
	if c.p == nil {
		return
	}
	C.CPUProfileDelete(c.p)
	c.p = nil
}
