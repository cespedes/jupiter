package jupiter

import (
	"bytes"
	"crypto/sha512"
	"fmt"
)

const (
	ScoreSize = 32
)

type Score struct {
	s [ScoreSize]byte // should have 256 bits (32 bytes)
}

var ZeroScore = Score{}

func GetScore(b []byte) Score {
	var s Score
	h := sha512.New512_256()
	h.Write(b)
	copy(s.s[:], h.Sum(nil))
	return s
}

func isBitEqual(b1 []byte, b2 []byte, n int) bool {
	if n >= 8 {
		return isBitEqual(b1[n/8:], b2[n/8:], n%8)
	}
	return (b1[0]&(1<<(7-n)) == b2[0]&(1<<(7-n)))
}

func (s Score) Equal(s2 Score) bool {
	return bytes.Equal(s.s[:], s2.s[:])
}

func (s Score) String() string {
	return fmt.Sprintf("%0x", s.s[:])
}

func (s Score) Match(s2 Score, mask int) bool {
	for i := 0; i < mask; i++ {
		if !isBitEqual(s.s[:], s2.s[:], i) {
			return false
		}
	}
	return true
}
