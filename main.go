package main

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"os"
)

const (
	BlockSize        = 8192
	DirtyBucket      = 1 << 31 // this bucket has not been synced to disk
	NotLeaf          = ^uint32(0)
	TableInitialSize = 1 << 4 // should be a power of 2
	EntriesInBucket  = 512
)

type Score [32]byte // should have 256 bits (32 bytes)

func (s Score) String() string {
	return fmt.Sprintf("%0x", s[:])
}

var zeroScore = Score{}

//type Bucket [BlockSize]byte

//type Bucket int
type Bucket [EntriesInBucket]Score

// 0: value is not a leaf
// TableIsDirty
var buckets []Bucket

func (b Bucket) size() int {
	var i int
	for i = 0; i < EntriesInBucket; i++ {
		if bytes.Equal(b[i][:], zeroScore[:]) {
			break
		}
	}
	return i
}

var table []uint32 // this is a binary heap

func log2(n uint) int {
	i := 0
	for n > 1<<i {
		i++
	}
	return i
}

func bit(b []byte, n int) bool {
	if n >= 8 {
		return bit(b[n/8:], n%8)
	}
	return (b[0]&(1<<(7-n)) != 0)
}

func init() {
	if 1<<log2(TableInitialSize) != TableInitialSize {
		panic(fmt.Sprintf("TableInitialSize is %d, should be a power of 2", TableInitialSize))
	}
	table = make([]uint32, 2*TableInitialSize-1)
	for i := 0; i < TableInitialSize-1; i++ {
		table[i] = NotLeaf
	}
	for i := 0; i < TableInitialSize; i++ {
		table[i+TableInitialSize-1] = uint32(i)
	}
	buckets = make([]Bucket, TableInitialSize)
}

// getBucket returns the entry in the binary heap table and the bucket number
func getBucket(s Score) (tableno int, bucket uint32) {
	i := 0
	for b := 0; ; b++ {
		if !bit(s[:], b) {
			i = i*2 + 1
		} else {
			i = i*2 + 2
		}
		if table[i] != NotLeaf {
			return i, table[i]
		}
	}
}

func insertScore(s Score) {
	t, b := getBucket(s)
	fmt.Printf("score=%v, table=%d, bucket=%d\n", s, t, b)
	for i := 0; i < EntriesInBucket; i++ {
		if bytes.Equal(buckets[b][i][:], zeroScore[:]) {
			buckets[b][i] = s
			return
		}
	}
	buckets = append(buckets, Bucket{})
	fmt.Printf("No more room in bucket %d; created new bucket %d and moving scores...\n", b, len(buckets)-1)
	if len(table) <= 2*t+2 {
		table = append(table, make([]uint32, 2*t+2-len(table)+1)...)
	}
	table[2*t+1] = table[t]
	table[2*t+2] = uint32(len(buckets) - 1)
	table[t] = NotLeaf
	var bb Bucket = buckets[b]
	for i := 0; i < EntriesInBucket; i++ {
		buckets[b][i] = zeroScore
	}
	for i := 0; i < EntriesInBucket; i++ {
		insertScore(bb[i])
	}
	fmt.Printf("%d scores moved.\n", EntriesInBucket)
}

func printTable2(list []int, pos int) {
	if table[pos] != NotLeaf {
		fmt.Printf("score=%v bucket=%d size=%d\n", list, table[pos], buckets[table[pos]].size())
	} else {
		printTable2(append(list, 0), pos*2+1)
		printTable2(append(list, 1), pos*2+2)
	}
}

func printTable() {
	printTable2([]int{0}, 1)
	printTable2([]int{1}, 2)
}

func main() {
	fmt.Println("This is Jupiter")
	//for i := 0; i < 1100; i++ {
	//	fmt.Printf("log2(%d)=%d\n", i, log2(uint(i)))
	//}
	fmt.Printf("table size is %d\n", len(table))
	for i := 1; i < len(os.Args); i++ {
		h := sha512.New512_256()
		h.Write([]byte(os.Args[i]))
		var s Score
		copy(s[:], h.Sum(nil))
		insertScore(s)
	}
	for i := 0; i < len(buckets); i++ {
		// fmt.Printf("bucket[%d]=%v\n", i, buckets[i])
		fmt.Printf("bucket %d: size=%d\n", i, buckets[i].size())
	}
	printTable()
}
