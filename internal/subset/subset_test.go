package subset

import (
	"reflect"
	"testing"
)

type testMember string

func (m testMember) String() string { return string(m) }

func TestSubsetKeepsSmallInput(t *testing.T) {
	in := []testMember{"a", "b"}
	got := Subset("key", in, 2)
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("Subset() = %v, want %v", got, in)
	}
}

func TestSubsetIsDeterministic(t *testing.T) {
	in := []testMember{"a", "b", "c", "d"}
	got := Subset("key", in, 2)
	again := Subset("key", in, 2)
	if !reflect.DeepEqual(got, again) {
		t.Fatalf("Subset() = %v, want %v", got, again)
	}
	if len(got) != 2 {
		t.Fatalf("len(Subset()) = %d, want 2", len(got))
	}
}
