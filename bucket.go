package jupiter

import (
	"encoding/binary"
)

const (
	BlockSize   = 8192
	DirtyBucket = 1 << 31 // this bucket has not been synced to disk
)

//type Bucket struct {
//	scoreBytes      int
//	addressBytes    int
//	scoreCommonBits int
//	commonScore     []byte
//	entries [EntriesInBucket]Entry
//}

// A Bucket is a byte array of a given size with a header and a list of entries.

type Bucket [BlockSize]byte

// Contents of the header:
// * The first 4 bytes will be "Jbkt" (magic number)
// * The next 2 bytes will have the number of used entries in this bucket
// * The next byte will have the number of bytes used to address
//   each entry in the data log (addressBytes)
// * The next byte will have the number of bytes for the score used    <--- is this really needed?
//   to identify each entry (scoreBytes)
// * The next byte will have the number of bits of the score which are
//   common to all the entries in this bucket (scoreCommonBits)
// * The next ((scoreCommonBits+7)/8) bytes will have the value of these bits
// * The rest of the array will have entries of size
//   (addressBytes + scoreBytes - scoreCommonBits/8)
// Each entry has 2 parts:
// * (scoreBytes - scoreCommonBits/8) bytes will be part of the score (right-aligned)
// * (addressBytes) bytes will be the address of the block in the data log

const (
	bktNumEntries         = 4
	bktNumAddressBytes    = 6
	bktNumScoreBytes      = 7
	bktNumScoreCommonBits = 8
	bktScoreCommonOffset  = 9
)

// Examples (blockSize=8192)
// -------------------------
// addressBytes scoreBytes scoreCommonBits header-size entry-size max-entries
// ------------ ---------- --------------- ----------- ---------- -----------
//           2          4               5          10          6        1363
//           6         10              10          11         15         545
// value of the initial bits in the score
// - the number of bytes in the score used to identify each entry
// - the number and value of the initial bits in the score which are
//   common to all the entries in this bucket

func newBucket(numScoreBytes int, numScoreCommonBits int, scoreCommonBytes []byte) *Bucket {
	b := new(Bucket)
	b[0] = 'J'
	b[1] = 'b'
	b[2] = 'k'
	b[3] = 't'
	b[bktNumAddressBytes] = 0
	b[bktNumScoreBytes] = byte(numScoreBytes)
	b[bktNumScoreCommonBits] = byte(numScoreCommonBits)
	for i := 0; i < (numScoreCommonBits+7)/8; i++ {
		b[bktScoreCommonOffset+i] = scoreCommonBytes[i]
	}
	return b
}

func (b *Bucket) numAddressBytes() int {
	return int(b[bktNumAddressBytes])
}

func (b *Bucket) numScoreBytes() int {
	return int(b[bktNumScoreBytes])
}

func (b *Bucket) numScoreCommonBits() int {
	return int(b[bktNumScoreCommonBits])
}

func (b *Bucket) CommonScore() (score Score, mask int) {
	mask := b.numScoreCommonBits()
	commonBytes := (mask + 7) / 8
	copy(score.s[:commonBytes], b[bktScoreCommonOffset:])

	return score, mask
}

func (b *Bucket) entrySize() int {
	return b.numAddressBytes() + b.numScoreBytes() - b.numScoreCommonBits()/8
}

func (b *Bucket) entryOffset() int {
	return bktNumScoreCommonBits + 1 + (b.numScoreCommonBits()+7)/8
}

func (b *Bucket) NumEntries() int {
	r := binary.BigEndian.Uint16(b[bktNumEntries : bktNumEntries+2])
	return int(r)
}

type Entry struct {
	score Score
	mask  int
	addr  uint64
}

func (b *Bucket) GetEntry(i int) *Entry {
	var e Entry
	scoreBytes := b.numScoreBytes()
	commonBits := b.numScoreCommonBits()
	commonBytes := (commonBits + 7) / 8
	copy(e.score.s[:commonBytes], b[bktScoreCommonOffset:])

	entryBytes := scoreBytes - commonBits/8
	offset := b.entryOffset() + i*b.entrySize()
	copy(e.score.s[scoreBytes-entryBytes:scoreBytes+1], b[offset:])

	e.mask = scoreBytes * 8

	addrBytes := b.numAddressBytes()
	j := offset + entryBytes
	for i := 0; i < addrBytes; i++ {
		e.addr *= 256
		e.addr += uint64(b[j])
		j++
	}
	return &e
}

// GetAddress returns the possible addresses for a given Score, if found in the bucket
func (b *Bucket) GetAddress(s Score) []uint64 {
	var result []uint64
	for i := 0; i < b.NumEntries(); i++ {
		e := b.GetEntry(i)
		if s.Match(e.score, e.mask) {
			result = append(result, e.addr)
		}
	}
	return result
}

// Add adds an entry to a bucket, identified by a score, and points it to a given address.
// It returns true if successful, false if there is no space in the bucket.
func (b *Bucket) Add(s Score, a uint64) bool {
	commonScore, mask = b.CommonScore()
	if !s.Match(commonScore, mask) {
		panic(fmt.Sprintf("Bucket.Add: score %s outsoide of %s/%d", s, commonScore, mask))
	}
	// TODO:
	//	if "a" fits in "numAddressBytes" {
	//		if (s, a) is already in the bucket {
	//			return true
	//		}
	//	} else {
	//		if there is not enough space to increment "numAddressBytes" *and* add the new score {
	//			return false
	//		}
	//		increment "numAddressBytes"
	//	}
	//	if there is not enough space to add the new score {
	//		return false
	//	}
	//	add new score
	panic("not implemented")
	return false
}

// Split divides the entries in a bucket into two, in order to add a new bucket to an index.
func (b *Bucket) Split(b2 *Bucket) error {
	panic("not implemented")
	return nil
}
