package jupiter

import (
	"fmt"
	"io"
)

const (
	BHNotLeaf = ^uint32(0)
)

type BinHeap struct {
	filename    string
	table       []uint32
	newBucketer NewBucketer
}

type NewBucketer interface {
	NewBucket(numScoreCommonBits int, scoreCommonBytes []byte) (uint32, error)
}

type Bucketer interface {
	Bucket(n uint32) *Bucket
}

// OpenBinHeap opens a BinHeap stored in disk
func OpenBinHeap(filename string) *BinHeap {
	panic("not implemented")
	return nil
}

// NewBinHeap creates a new binary heap, to be used as a hash function combined with an Index
func NewBinHeap(f NewBucketer) (*BinHeap, error) {
	firstBucket, err := f.NewBucket(0, []byte{})
	if err != nil {
		return nil, err
	}
	bh := new(BinHeap)
	bh.table = []uint32{firstBucket}
	return &BinHeap{table: []uint32{firstBucket}}, nil
}

func isBitSet(b []byte, n int) bool {
	if n >= 8 {
		return isBitSet(b[n/8:], n%8)
	}
	return (b[0]&(1<<(7-n)) != 0)
}

func (bh *BinHeap) Sync() error {
	panic("not implemented")
	return nil
}

func (bh *BinHeap) Close() error {
	panic("not implemented")
	return nil
}

// Get returns one entry in the binary heap
func (bh *BinHeap) Get(k int) (uint32, error) {
	if k < 0 || k >= len(bh.table) {
		return 0, fmt.Errorf("BinHeap.Get: key=%d out of bounds (should be between 0 and %d)", k, len(bh.table)-1)
	}
	return bh.table[k], nil
}

// GetBucket returns the entry in the binary heap table and the bucket number
func (bh *BinHeap) GetBucket(s Score) (int, uint32) {
	i := 0
	if len(bh.table) == 0 {
		panic("len(bh.table)=nil (should not happen")
	}
	for b := 0; ; b++ {
	        if bh.table[i] != BHNotLeaf {
			return i, bh.table[i]
	        }
	        if !isBitSet(s.s[:], b) {
			i = i*2 + 1
	        } else {
			i = i*2 + 2
	        }
	}
}

// Write stores a BinHeap into disk
func (bh *BinHeap) Write(f io.Writer) error {
	panic("not implemented")
	return nil
}

// Set sets a new value 'v' for entry 'k' in the binary heap
func (bh *BinHeap) Set(k int, v uint32) error {
	if k < 0 || k >= len(bh.table) {
		return fmt.Errorf("BinHeap.Set: key=%d out of bounds (should be between 0 and %d)", k, len(bh.table)-1)
	}
	if bh.table[k] == BHNotLeaf {
		return fmt.Errorf("BinHeap.Set: key=%d: invalid argument (this is not a leaf)", k)
	}
	bh.table[k] = v
	return nil
}

// NewLeaf replaces one leaf with two leaves: the left one will have the old value, the new one will have v
func (bh *BinHeap) NewLeaf(k int, v uint32) error {
	if k < 0 || k >= len(bh.table) {
		return fmt.Errorf("BinHeap.Set: invalid argument for key")
	}
	if bh.table[k] == BHNotLeaf {
		return fmt.Errorf("BinHeap.Set: key=%d: invalid argument (this is not a leaf)", k)
	}
	if len(bh.table) <= 2*k+2 {
		bh.table = append(bh.table, make([]uint32, 2*k+3-len(bh.table))...)
	}
	bh.table[2*k+1] = bh.table[k]
	bh.table[2*k+2] = v
	bh.table[k] = BHNotLeaf
	return nil
}

// Split replaces one leaf with two leaves
func (bh *BinHeap) Split(k int, f interface{NewBucketer; Bucketer}) error {
	if k < 0 || k >= len(bh.table) {
		return fmt.Errorf("BinHeap.Set: invalid argument for key")
	}
	if bh.table[k] == BHNotLeaf {
		return fmt.Errorf("BinHeap.Set: key=%d: invalid argument (this is not a leaf)", k)
	}
	if len(bh.table) <= 2*k+2 {
		bh.table = append(bh.table, make([]uint32, 2*k+3-len(bh.table))...)
	}
	// numScoreCommonBits := 
	// scoreCommonBits := 
	// nextBucket, err := f.NewBucket(numScoreCommonBits, scoreCommonBytes)
	// if err != nil {
	// 	return nil, err
	// }
	bh.table[2*k+1] = bh.table[k]
	// bh.table[2*k+2] = nextBucket
	bh.table[k] = BHNotLeaf
	b1 := f.Bucket(k)
	b2 := f.Bucket(nextBucket)
	return nil
}
