package jupiter

// A Index has several buckets of the same size (8192 bytes)
// Each bucket has a header and several entries of the same size (but possibly different among different buckets)
import (
	"fmt"
	"io"
)

type Index struct {
	filename          string
	scoreBytesInEntry int
	numBuckets        uint32
	maxBuckets        uint32 // 0 if there is no limit
	buckets           map[uint32]*Bucket
}

func OpenIndex(filename string) *Index {
	panic("not implemented")
	return &Index{filename: filename}
}

func NewIndex(filename string, scoreBytesInEntry int) *Index {
	return &Index{filename: filename, scoreBytesInEntry: scoreBytesInEntry}
}

func (in *Index) NumBuckets() uint32 {
	return in.numBuckets
}

func (in *Index) MaxBuckets() uint32 {
	return in.maxBuckets
}

func (in *Index) Bucket(n uint32) *Bucket {
	b := in.buckets[n]
	if b != nil {
		return b
	}
	panic("not implemented")
	return nil
}

func (in *Index) Sync() error {
	panic("not implemented")
	return nil
}

func (in *Index) Write(f io.Writer, numBlocks uint64) error {
	panic("not implemented")
	return nil
}

func (in *Index) NewBucket() (uint32, error) {
	if in.maxBuckets > 0 && in.numBuckets +1 >= in.maxBuckets {
		return 0, fmt.Errorf("Not enough size in index for a new bucket")
	}
	in.buckets[in.numBuckets] = newBucket(in.scoreBytesInEntry)
	in.numBuckets++
	return in.numBuckets-1, nil
}
