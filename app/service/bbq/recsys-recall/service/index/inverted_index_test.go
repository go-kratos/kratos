package index_test

import (
	"testing"

	"go-common/app/service/bbq/recsys-recall/service/index"
)

func TestInvertedIndex(t *testing.T) {
	src := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	srcII := &index.InvertedIndex{
		Data: src,
	}
	raw := srcII.Serialize()
	dstII := &index.InvertedIndex{}
	dstII.Load(raw)
	for i := range src {
		if src[i] != dstII.Data[i] {
			t.Error("incorrect data")
		}
	}
}
