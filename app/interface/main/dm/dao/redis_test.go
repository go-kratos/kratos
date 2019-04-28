package dao

import (
	"context"
	"testing"
)

func TestUptUsrPaCnt(t *testing.T) {
	err := testDao.UptUsrPaCnt(context.TODO(), 1234356789, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPaUsrCnt(t *testing.T) {
	cnt, err := testDao.PaUsrCnt(context.TODO(), 1234356789)
	if err != nil {
		t.Fatal(err)
	}
	if cnt == 0 {
		t.Fatal("pa usr cnt err")
	}
}

func TestUsrDMAccCnt(t *testing.T) {
	n, err := testDao.UsrDMAccCnt(context.TODO(), 1234356789, 1)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("usr dm acc cnt err")
	}
}

func TestUptRecallCnt(t *testing.T) {
	err := testDao.UptRecallCnt(context.TODO(), 1234356789)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRecallCnt(t *testing.T) {
	n, err := testDao.RecallCnt(context.TODO(), 1234356789)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("usr dm acc cnt err")
	}
}
func TestPaLock(t *testing.T) {
	incr, err := testDao.PaLock(context.TODO(), "go-test")
	if err != nil {
		t.Fatal(err)
	}
	if incr < 1 {
		t.Fatal("PaLock err")
	}
}
