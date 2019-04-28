package dao

import (
	"testing"

	"go-common/app/service/live/userexp/model"
)

func TestLevelCache(t *testing.T) {
	once.Do(startService)
	rs, err := d.LevelCache(ctx, 10001)
	if err != nil {
		t.Error("d.LevelCache err:", err.Error())
	} else {
		t.Logf("LevelCache %v", rs)
	}
}

func TestSetLevelCache(t *testing.T) {
	once.Do(startService)
	err := d.SetLevelCache(ctx, &model.Level{Uid: 10001, Uexp: 1000, Rexp: 100, Ulevel: 2, Rlevel: 1, Color: 12345})
	if err != nil {
		t.Error("d.SetLevelCache err:", err.Error())
	} else {
		t.Logf("SetLevelCache Succ!")
	}
}

func TestDelLevelCache(t *testing.T) {
	once.Do(startService)
	err := d.DelLevelCache(ctx, 10001)
	if err != nil {
		t.Error("d.DelLevelCache err:", err.Error())
	} else {
		t.Logf("DelLevelCache Succ!")
	}
}

func TestMuitiLevelCache(t *testing.T) {
	once.Do(startService)
	rs, _, err := d.MultiLevelCache(ctx, []int64{10001, 10002})
	if err != nil {
		t.Error("d.LevelCache err:", err.Error())
	} else {
		t.Logf("LevelCache %v", rs)
	}
}
