package jupiter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrorNotFound = errors.New("Not Found")

type ReadWriteSeekCloser interface {
	io.ReadWriteSeeker
	io.Closer
}

type DataLog struct {
	filename string
	fp       ReadWriteSeekCloser
}

// Global header: version

// Each block is prefixed by a header that describes the contents of the
// block.

// The header contains the score and the uncompressed size.

// In the future, the header may also contain the type, a reference count,
// the compression algorithm, the compressed size, and a checksum.

func OpenDataLog(filename string) *DataLog {
	panic("OpenDataLog() not implemented")
	return nil
}

// NewDataLog creates a new DataLog without a physical back-up (data is stored in memory)
func NewDataLog() *DataLog {
	d := new(DataLog)
	d.fp = newSeekableBuffer()
	return d
}

// WriteChunk stores a block of data, compressing it, and returns its address in the data log.
func (d *DataLog) WriteChunk(score Score, t Type, b []byte) (addr uint64, err error) {
	fmt.Printf("DEBUG: DataLog.WriteChunk(score=%s, type=%d, b=%q)\n", score, t, b)
	position, err := d.fp.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	if len(b) >= 0xFFFF {
		return 0, fmt.Errorf("WriteChunk(): tried to write more than 64K of data")
	}
	defer func() {
		if err != nil {
			if x, ok := d.fp.(interface{ Truncate(size int64) (err error) }); ok {
				x.Truncate(position)
			}
		}
	}()
	_, err = d.fp.Write(score.s[:])
	if err != nil {
		return 0, err
	}
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(len(b)))
	_, err = d.fp.Write(l)
	if err != nil {
		return 0, err
	}
	_, err = d.fp.Write(b)
	if err != nil {
		return 0, err
	}
	return uint64(position), nil
}

// PeekChunk is used to check if a given block is stored at an address
func (d *DataLog) PeekChunk(score Score, addr uint64) (t Type, err error) {
	var s Score
	_, err = d.fp.Seek(int64(addr), os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	_, err = io.ReadFull(d.fp, s.s[:])
	if err != nil {
		return 0, err
	}
	if s.Equal(score) {
		return 0, nil
	}
	return 0, ErrorNotFound
}

// GetChunk returns the block with a given score stored at an address
func (d *DataLog) ReadChunk(score Score, addr uint64) (t Type, b []byte, err error) {
	fmt.Printf("DEBUG: DataLog.ReadChunk(score=%s, addr=%d)\n", score, addr)
	_, err = d.fp.Seek(int64(addr), os.SEEK_SET)
	if err != nil {
		return 0, nil, err
	}

	buf := make([]byte, ScoreSize+2)
	_, err = io.ReadFull(d.fp, buf)
	if err != nil {
		return 0, nil, err
	}

	var s Score
	copy(s.s[:], buf)
	if !s.Equal(score) {
		return 0, nil, ErrorNotFound
	}
	size := binary.BigEndian.Uint16(buf[ScoreSize:])
	buf = make([]byte, size)
	_, err = io.ReadFull(d.fp, buf)
	if err != nil {
		return 0, nil, err
	}
	return 0, buf, nil
}

// func (j *Jupiter) Write(t Type, b []byte) (Score, error) {
