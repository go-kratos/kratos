package dao

import (
	"context"
	"testing"
)

func TestAdvances(t *testing.T) {
	_, _, err := testDao.Advances(context.TODO(), 27515260, "all", "all", 1, 20)
	if err != nil {
		t.Error(err)
	}
}
