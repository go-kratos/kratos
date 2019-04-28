package lrucache

import (
	"sync"
	"testing"
	"time"
)

func Test_hashCode(t *testing.T) {
	/*if hashCode(-1) != 1 {
		t.Error("case 1 failed")
	}
	if hashCode(0) != 0 {
		t.Error("case 2 failed")
	}
	if hashCode(0x7FFFFFFF) != 0x7FFFFFFF {
		t.Error("case 3 failed")
	}*/
	if hashCode("12345") != 3421846044 {
		t.Error("case 4 failed")
	}
	if hashCode("abcdefghijklmnopqrstuvwxyz") != 1277644989 {
		t.Error("case 5 failed")
	}
	/*if hashCode(123.45) != 123 {
		t.Error("case 6 failed")
	}
	if hashCode(-15268.45) != 15268 {
		t.Error("case 7 failed")
	}*/
}

func Test_nextPowOf2(t *testing.T) {
	if nextPowOf2(0) != 2 {
		t.Error("case 1 failed")
	}
	if nextPowOf2(1) != 2 {
		t.Error("case 2 failed")
	}
	if nextPowOf2(2) != 2 {
		t.Error("case 3 failed")
	}
	if nextPowOf2(3) != 4 {
		t.Error("case 4 failed")
	}
	if nextPowOf2(123) != 128 {
		t.Error("case 5 failed")
	}
	if nextPowOf2(0x7FFFFFFF) != 0x80000000 {
		t.Error("case 6 failed")
	}
}

func Test_timeout(t *testing.T) {
	sc := NewSyncCache(1, 2, 2)
	sc.Put("1", "2")
	if v, ok := sc.Get("1"); !ok || v != "2" {
		t.Error("case 1 failed")
	}
	time.Sleep(2 * time.Second)
	if _, ok := sc.Get("1"); ok {
		t.Error("case 2 failed")
	}
}

func Test_concurrent(t *testing.T) {
	sc := NewSyncCache(1, 4, 2)
	var wg sync.WaitGroup
	for index := 0; index < 100000; index++ {
		wg.Add(3)
		go func() {
			sc.Put("1", "2")
			wg.Done()
		}()
		go func() {
			sc.Get("1")
			wg.Done()
		}()
		go func() {
			sc.Delete("1")
			wg.Done()
		}()
	}
	wg.Wait()
}
