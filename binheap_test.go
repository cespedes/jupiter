package jupiter

import (
	"testing"
)

type newbucketer struct{}

func (_ newbucketer) NewBucket(numScoreCommonBits int, scoreCommonBytes []byte) (uint32, error) {
	return 0, nil
}

func TestBinHeap(t *testing.T) {
	bh, err := NewBinHeap(newbucketer{})
	if bh == nil || err != nil{
		t.Errorf("error calling NewBinHeap()")
	}
	s1 := ZeroScore
	s2 := GetScore([]byte{})
	k, v := bh.GetBucket(s1)
	t.Logf("GetBucket(s1) = (%d,%d)", k, v)
	k, v = bh.GetBucket(s2)
	t.Logf("GetBucket(s2) = (%d,%d)", k, v)
	err = bh.NewLeaf(0, 1)
	if err != nil {
		t.Errorf("NewLeaf(0,1): %v", err)
	}
	k, v = bh.GetBucket(s1)
	t.Logf("GetBucket(s1) = (%d,%d)", k, v)
	k, v = bh.GetBucket(s2)
	t.Logf("GetBucket(s2) = (%d,%d)", k, v)
	err = bh.NewLeaf(0, 1)
	if err == nil {
		t.Errorf("NewLeaf(0,1): should return error")
	}
	err = bh.NewLeaf(2, 2)
	if err != nil {
		t.Errorf("NewLeaf(1,2): %v", err)
	}
	err = bh.NewLeaf(6, 3)
	if err != nil {
		t.Errorf("NewLeaf(2,3): %v", err)
	}
	k, v = bh.GetBucket(s1)
	t.Logf("GetBucket(s1) = (%d,%d)", k, v)
	k, v = bh.GetBucket(s2)
	t.Logf("GetBucket(s2) = (%d,%d)", k, v)
	for i:=0; ; i++ {
		v, err := bh.Get(i)
		if err != nil {
			break
		}
		t.Logf("Get(%d) = %d", i, v)
	}
}
