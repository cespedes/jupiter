package jupiter

import (
	"testing"
)

func TestScore(t *testing.T) {
	var s string
	t.Log("Testing ZeroScore...")
	s = ZeroScore.String()
	if s != `0000000000000000000000000000000000000000000000000000000000000000` {
		t.Errorf("ZeroScore is %s (should be 0)", s)
	}
	t.Log(`Testing Score("")...`)
	s = GetScore([]byte{}).String()
	if s != `c672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a` {
		t.Errorf(`Score("") is %s (should be c672b8d1...)`, s)
	}
	t.Log(`Testing Score(Score(""))...`)
	s = GetScore([]byte(s)).String()
	if s != `a8dfa748acde437f7b36261eae56f6e6078de914de28f87295d4f7aac568b16e` {
		t.Errorf(`Score(Score("")) is %s (should be a8dfa748...)`, s)
	}
}
