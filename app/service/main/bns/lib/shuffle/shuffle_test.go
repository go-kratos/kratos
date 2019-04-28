package shuffle

import (
	"strings"
	"testing"
)

type List []string

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func TestShuffle(t *testing.T) {
	l := List{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	old := strings.Join(l, "")
	Shuffle(l)
	new := strings.Join(l, "")
	if old == new {
		t.Errorf("shuffle error, %s == %s", old, new)
	}
	t.Log(new)
}
