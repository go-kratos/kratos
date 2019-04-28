package dao

import (
	"context"
	"testing"
)

func TestAddRecentDM(t *testing.T) {
	count, err := testDao.AddRecentDM(context.TODO(), 123, dm)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("count:%d", count)
}

func TestZRemRecentDM(t *testing.T) {
	if err := testDao.ZRemRecentDM(context.TODO(), 123, 719150142); err != nil {
		t.Fatal(err)
	}
}

func TestTrimRecentDM(t *testing.T) {
	if err := testDao.TrimRecentDM(context.TODO(), 123, 1); err != nil {
		t.Fatal(err)
	}
}
