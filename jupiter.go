package jupiter

import (
	"fmt"
)

const ScoreBytesInEntry = 10

type Jupiter struct {
	config   *Config
	binheap  *BinHeap
	// indexes  []*Index   // Just one index and one datalog for now
	// datalogs []*DataLog
	index   *Index
	datalog *DataLog
}

type Type byte

func Open(c *Config) *Jupiter {
	return nil
}

func New() (*Jupiter, error) {
	var j Jupiter
	var err error
	index := NewIndex(ScoreBytesInEntry)
	j.index = index
	j.binheap, err = NewBinHeap(j.index)
	if err != nil {
		return nil, err
	}
	j.datalog = NewDataLog()
	return &j, nil
}

func (j *Jupiter) Read(score Score) (Type, []byte, error) {
	_, buckn := j.binheap.GetBucket(score)
	bucket := j.index.Bucket(buckn)
	addrs := bucket.GetAddress(score)
	for _, addr := range addrs {
		t, b, err := j.datalog.GetChunk(score, addr)
		if err == nil {
			return t, b, nil
		}
	}
	return 0, nil, fmt.Errorf("jupiter: Read(%s): not found", score)
}

func (j *Jupiter) Write(t Type, b []byte) (Score, error) {
	score := GetScore(b)
	_, buckn := j.binheap.GetBucket(score)
	bucket := j.index.Bucket(buckn)
	addrs := bucket.GetAddress(score)
	for _, addr := range addrs {
		tt, err := j.datalog.PeekChunk(score, addr)
		if t==tt && err == nil {
			return score, nil
		}
		if t != tt && err == nil {
			return ZeroScore, fmt.Errorf("Jupiter.Write: block already written with different type")
		}
	}
	addr, err := j.datalog.NewChunk(score, t, b)
	if err != nil {
		return ZeroScore, err
	}
	if bucket.Add(score, addr) {
		return score, nil
	}
	// There is no room in bucket, we need another one
	// j.index.NewBucket(numScoreCommonBits int, scoreCommonBytes []byte) (uint32, error) {
	// bucket.Split(b2 *Bucket) error {
	// j.binheap.NewLeaf(k int, v uint32) error {
	panic("not implemented")
	return ZeroScore, nil
}
