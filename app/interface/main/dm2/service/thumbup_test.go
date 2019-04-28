package service

import (
	"context"
	"testing"
)

func TestServiceAddAct(t *testing.T) {
	var (
		c              = context.TODO()
		err            error
		cid, dmid, mid int64 = 5, 719149462, 3078992
		op             int8  = 1
	)
	err = svr.ThumbupDM(c, cid, dmid, mid, op)
	if err != nil {
		t.Error(err)
	}
}

func TestServiceLikeList(t *testing.T) {
	var (
		c        = context.TODO()
		err      error
		cid, mid int64 = 5, 3078992
		dmids          = []int64{719149462, 719149463}
	)
	m, err := svr.ThumbupList(c, cid, mid, dmids)
	if err != nil {
		t.Error(err)
	}
	for k, v := range m {
		t.Logf("=====%d: %+v", k, v)
	}
	t.Logf("=====%+v", m)
}
