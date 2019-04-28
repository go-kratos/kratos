package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/dm/conf"
)

func TestDuration(t *testing.T) {
	var (
		err      error
		res      int64
		cid      = int64(1)
		duration = int64(100)
		c        = context.TODO()
	)
	d := New(conf.Conf)
	if err = d.SetDurationCache(c, cid, duration); err != nil {
		t.Errorf("d.SetDurationCache(%d,%d) error(%v)", cid, duration, err)
	}
	if res, err = d.DurationCache(c, cid); err != nil {
		t.Errorf("d.DurationCache(%d) error(%v)", cid, err)
	}
	if res != duration {
		t.Errorf("not expect value %d,%d", duration, res)
	}
}
