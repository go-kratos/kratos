package service

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkTsSecond2Hash(b *testing.B) {
	ts := time.Now().Unix()
	for i := 0; i < b.N; i++ {
		ts += int64(rand.Intn(100))
		hash := TsSeconds2Hash(ts)
		if ats, err := Hash2TsSeconds(hash); err != nil {
			b.Errorf("failed to parse hash, error(%v)", err)
			b.FailNow()
		} else if ats != ts {
			b.Errorf("hash: %s, expect %d but got %d", hash, ts, ats)
			b.FailNow()
		}
	}
}

func TestDecompress(t *testing.T) {
	val := table[5] // 0xabd28536
	exp := []uint64{6, 3, 5, 8, 2, 13, 11, 10}
	poses := decompress(val)
	if len(poses) != len(exp) {
		t.Errorf("res is not correct, expected poses length equal to exp length, but not, poses length: %d, exp length: %d", len(poses), len(exp))
		t.FailNow()
	}
	for i, item := range exp {
		if poses[i] != item {
			t.Errorf("failed to decompress, expected %v but got %v", exp, poses)
			t.FailNow()
		}
	}
}
