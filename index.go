package jupiter

import (
	"fmt"
	"io"
)

// An Index has several buckets of the same size (8192 bytes)
// Each bucket has a header and several entries of the same size (but possibly different among different buckets)
type Index struct {
	filename          string
	scoreBytesInEntry int
	// A "firstBucket" could be used to support several indexes
	// TODO firstBucket       uint32
	numBuckets uint32
	maxBuckets uint32 // 0 if there is no limit
	buckets    map[uint32]*Bucket
}

// OpenIndex opens a file used as Index, and reads its contents into memory
func OpenIndex(filename string) *Index {
	panic("OpenIndex() not implemented")
	return &Index{filename: filename}
}

// NexIndex creates a new index from scratch
//func NewIndex(scoreBytesInEntry int, firstBucket uint32) *Index {
//	return &Index{scoreBytesInEntry: scoreBytesInEntry, firstBucket: firstBucket}
//}
func NewIndex(scoreBytesInEntry int) *Index {
	return &Index{scoreBytesInEntry: scoreBytesInEntry, buckets: make(map[uint32]*Bucket)}
}

// NumBuckets returns the number of buckets in an Index
func (in *Index) NumBuckets() uint32 {
	return in.numBuckets
}

// MaxBuckets returns the maximum number of buckets this Index can have
func (in *Index) MaxBuckets() uint32 {
	return in.maxBuckets
}

// Bucket returns a Bucket given its position
func (in *Index) Bucket(n uint32) *Bucket {
	b := in.buckets[n]
	if b != nil {
		return b
	}
	panic("Index.Bucket(): read from file: not implemented")
	return nil
}

// Sync writes the dirty buckets into disk
func (in *Index) Sync() error {
	panic("Index.Sync(): not implemented")
	return nil
}

// Write writes the whole index into disk
func (in *Index) Write(f io.Writer, numBlocks uint64) error {
	panic("Index.Write(): not implemented")
	return nil
}

// NewBucket adds a new bucket to the Index.
func (in *Index) NewBucket(numScoreCommonBits int, scoreCommonBytes []byte) (uint32, error) {
	if in.maxBuckets > 0 && in.numBuckets+1 >= in.maxBuckets {
		return 0, fmt.Errorf("Not enough size in index for a new bucket")
	}
	in.buckets[in.numBuckets] = newBucket(in.scoreBytesInEntry, numScoreCommonBits, scoreCommonBytes)
	in.numBuckets++
	return in.numBuckets - 1, nil
}
