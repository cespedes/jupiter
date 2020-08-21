package jupiter

import (
	"encoding/binary"
	"fmt"
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
	bktNumEntries         = 4 // Number of entries stored in bucket
	bktNumAddressBytes    = 6 // Number of bytes for addressing in each entry
	bktNumScoreBytes      = 7 // Total number of bytes used to address each entry
	bktNumScoreCommonBits = 8 // Number of bits for the score common to all the entries
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
	b[bktNumAddressBytes] = 1
	b[bktNumScoreBytes] = byte(numScoreBytes)
	b[bktNumScoreCommonBits] = byte(numScoreCommonBits)
	for i := 0; i < (numScoreCommonBits+7)/8; i++ {
		b[bktScoreCommonOffset+i] = scoreCommonBytes[i]
	}
	return b
}

// Init prepares a Bucket for its use.  A Bucket cannot be used until it has been Init'd
func (b *Bucket) Init(numScoreBytes int, numScoreCommonBits int, scoreCommonBytes []byte) {
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
	mask = b.numScoreCommonBits()
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
	copy(e.score.s[scoreBytes-entryBytes:scoreBytes], b[offset:])

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

func numBytesInUint64(a uint64) int {
	if a == 0 {
		return 0
	} else {
		return 1 + numBytesInUint64(a>>8)
	}
}

// Add adds an entry to a bucket, identified by a score, and points it to a given address.
// It returns true if successful, false if there is no space in the bucket.
func (b *Bucket) Add(s Score, a uint64) bool {
	fmt.Printf("DEBUG: Bucket.Add(%s, %d)\n", s, a)
	commonScore, mask := b.CommonScore()
	if !s.Match(commonScore, mask) {
		panic(fmt.Sprintf("Bucket.Add(): score %s outside of %s/%d", s, commonScore, mask))
	}

	numEntries := b.NumEntries()

	// Does "a" fit in "numAddressBytes"?
	if numBytesInUint64(a) <= b.numAddressBytes() {
		// Is this entry already added?
		for i := 0; i < numEntries; i++ {
			e := b.GetEntry(i)
			if s.Match(e.score, e.mask) && a == e.addr {
				return true
			}
		}
	} else {
		newEntrySize := numBytesInUint64(a) + b.numScoreBytes() - b.numScoreCommonBits()/8
		if (numEntries+1)*newEntrySize > BlockSize-b.entryOffset() {
			return false
		}
		// TODO: increment "numAddressBytes"
		fmt.Printf("numAddressBytes=%d, a=%d\n", b.numAddressBytes(), a)
		panic("Bucket.Add(): incrementing numAddressBytes: not implemented")
	}
	addressBytes := b.numAddressBytes()
	scoreBytes := b.numScoreBytes()
	commonBits := b.numScoreCommonBits()
	scoreEntryBytes := scoreBytes - commonBits/8

	entrySize := addressBytes + scoreEntryBytes
	maxEntries := (BlockSize - b.entryOffset()) / entrySize
	if maxEntries <= numEntries {
		// not enough space to add the new score
		return false
	}
	offset := b.entryOffset() + numEntries*entrySize

	// Add score to entry:
	copy(b[offset:offset+scoreEntryBytes], s.s[scoreBytes-scoreEntryBytes:])

	// Add addr to entry:
	a2 := make([]byte, 8)
	binary.BigEndian.PutUint64(a2, a)
	copy(b[offset+scoreEntryBytes:], a2[8-addressBytes:])

	// Increment NumEntries:
	binary.BigEndian.PutUint16(b[bktNumEntries:bktNumEntries+2], uint16(numEntries+1))
	return true
}

// Split divides the entries in a bucket into two, in order to add a new bucket to an index.
func (b *Bucket) Split(b2 *Bucket) error {
	panic("Bucket.Split(): not implemented")
	return nil
}
