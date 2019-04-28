package service

import (
	"context"
	"testing"
)

func TestVideoDuration(t *testing.T) {
	var (
		aid int64 = 10097265
		oid int64 = 1508
		c         = context.TODO()
	)
	d, err := svr.videoDuration(c, aid, oid)
	if err != nil {
		t.Errorf("s.videoDuration(%d %d) error(%v)", aid, oid, err)
		t.FailNow()
	}
	t.Logf("oid:%d duration:%d", oid, d)
}
