package jupiter

import (
	"bytes"
	"fmt"
)

const (
	ScoreSize = 32
)

type Score [ScoreSize]byte // should have 256 bits (32 bytes)

func (s Score) Equal(s2 Score) bool {
	return bytes.Equal(s[:], s2[:])
}

func (s Score) String() string {
	return fmt.Sprintf("%0x", s[:])
}

var ZeroScore = Score{}
