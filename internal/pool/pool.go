// Package pool provides wrapper functions of sync.Pool for efficiency.
package pool

import (
	"bytes"
	"fmt"
	"sync"
)

var bytesBuffer = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBytesBuffer returns an item from the *bytes.Buffer pool.
// Put for *bytes.Buffer calls (*bytes.Buffer).Reset before
// putting to the pool so the result of Get can be used without
// resetting.
func GetBytesBuffer() *bytes.Buffer {
	buf := bytesBuffer.Get().(*bytes.Buffer)
	return buf
}

// Put adds x to the pool for the type.
// This function causes panic if the type of the value is unknown.
func Put(x interface{}) {
	switch v := x.(type) {
	case *bytes.Buffer:
		putBytesBuffer(v)
	default:
		panic(fmt.Sprintf("unknown type: %T", x))
	}
}

func putBytesBuffer(x *bytes.Buffer) {
	x.Reset()
	bytesBuffer.Put(x)
}
