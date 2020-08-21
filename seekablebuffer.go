package jupiter

/*
Copyright 2019 Random Ingenuity InformationWorks

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"fmt"
	"io"
	"os"
)

// seekableBuffer is a simple memory structure that satisfies
// `io.ReadWriteSeeker` and `io.Closer`.
type seekableBuffer struct {
	data     []byte
	position int64
}

// newSeekableBuffer is a factory that returns a `*SeekableBuffer`.
func newSeekableBuffer() *seekableBuffer {
	data := make([]byte, 0)

	return &seekableBuffer{
		data: data,
	}
}

// newSeekableBufferWithBytes is a factory that returns a `*SeekableBuffer`.
func newSeekableBufferWithBytes(originalData []byte) *seekableBuffer {
	data := make([]byte, len(originalData))
	copy(data, originalData)

	return &seekableBuffer{
		data: data,
	}
}

func len64(data []byte) int64 {
	return int64(len(data))
}

// Bytes returns the underlying slice.
func (sb *seekableBuffer) Bytes() []byte {
	return sb.data
}

// Len returns the number of bytes currently stored.
func (sb *seekableBuffer) Len() int {
	return len(sb.data)
}

// Write does a standard write to the internal slice.
func (sb *seekableBuffer) Write(p []byte) (n int, err error) {
	// The current position we're already at is past the end of the data we
	// actually have. Extend our buffer up to our current position.
	if sb.position > len64(sb.data) {
		extra := make([]byte, sb.position-len64(sb.data))
		sb.data = append(sb.data, extra...)
	}

	positionFromEnd := len64(sb.data) - sb.position
	tailCount := positionFromEnd - len64(p)

	var tailBytes []byte
	if tailCount > 0 {
		tailBytes = sb.data[len64(sb.data)-tailCount:]
		sb.data = append(sb.data[:sb.position], p...)
	} else {
		sb.data = append(sb.data[:sb.position], p...)
	}

	if tailBytes != nil {
		sb.data = append(sb.data, tailBytes...)
	}

	dataSize := len64(p)
	sb.position += dataSize

	return int(dataSize), nil
}

// Read does a standard read against the internal slice.
func (sb *seekableBuffer) Read(p []byte) (n int, err error) {
	if sb.position >= len64(sb.data) {
		return 0, io.EOF
	}

	n = copy(p, sb.data[sb.position:])
	sb.position += int64(n)

	return n, nil
}

// Truncate either chops or extends the internal buffer.
func (sb *seekableBuffer) Truncate(size int64) (err error) {
	sizeInt := int(size)
	if sizeInt < len(sb.data)-1 {
		sb.data = sb.data[:sizeInt]
	} else {
		new := make([]byte, sizeInt-len(sb.data))
		sb.data = append(sb.data, new...)
	}

	return nil
}

// Seek does a standard seek on the internal slice.
func (sb *seekableBuffer) Seek(offset int64, whence int) (n int64, err error) {
	if whence == os.SEEK_SET {
		sb.position = offset
	} else if whence == os.SEEK_END {
		sb.position = len64(sb.data) + offset
	} else if whence == os.SEEK_CUR {
		sb.position += offset
	} else {
		panic(fmt.Sprintf("seek whence is not valid: (%d)", whence))
	}

	if sb.position < 0 {
		sb.position = 0
	}

	return sb.position, nil
}

// Close does nothing.
func (sb *seekableBuffer) Close() error {
	return nil
}
