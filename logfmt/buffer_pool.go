package logfmt

import (
	"bytes"
	"sync"
)

// bufferPool provides buffers for writing logfmt messages.
// Reusing buffers is intended to put less stress on the GC
// in times when large numbers of messages are being logged.
var bufferPool sync.Pool
var bufferPoolAllocated = func() {}

func getBuffer() *bytes.Buffer {
	buf, ok := bufferPool.Get().(*bytes.Buffer)
	if !ok {
		buf = &bytes.Buffer{}
		bufferPoolAllocated()
	}
	return buf
}

func releaseBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}
