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

func (s Score) Equal(s2 Score) bool {
	return bytes.Equal(s.s[:], s2.s[:])
}

func (s Score) String() string {
	return fmt.Sprintf("%0x", s.s[:])
}

func CalcScore(b []byte) Score {
	var s Score
	h := sha512.New512_256()
	h.Write(b)
	copy(s.s[:], h.Sum(nil))
	return s
}
