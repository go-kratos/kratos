package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/dm2/model"
)

var (
	c  = context.TODO()
	dm = &model.DM{
		ID:       719150142,
		Oid:      1221,
		Type:     1,
		Mid:      478046,
		Progress: 0,
		State:    0,
		Content: &model.Content{
			ID:       719150142,
			FontSize: 24,
			Mode:     1,
			Msg:      "aaa",
		}}
)

func TestAddDMCache(t *testing.T) {
	if err := testDao.AddDMCache(context.TODO(), dm); err != nil {
		t.Error(err)
	}
}

func TestSetDMCache(t *testing.T) {
	if err := testDao.SetDMCache(c, dm.Type, dm.Oid, []*model.DM{dm, dm}); err != nil {
		t.Error(err)
	}
}

func TestDelDMCache(t *testing.T) {
	if err := testDao.DelDMCache(context.TODO(), 1, 1221); err != nil {
		t.Error(err)
	}
}

func TestExpireDMCache(t *testing.T) {
	ok, err := testDao.ExpireDMCache(context.TODO(), 1, 1221)
	if err != nil {
		t.Error(err)
	}
	t.Log(ok)
}

func TestDMCache(t *testing.T) {
	values, err := testDao.DMCache(context.TODO(), dm.Type, dm.Oid)
	if err != nil {
		t.Error(err)
	}
	for _, value := range values {
		dmCache := &model.DM{}
		if err = dmCache.Unmarshal(value); err != nil {
			t.Errorf("Unmarshal(%s) error(%v)", value, err)
		}
		t.Log(dmCache)
	}
}

func TestTrimDMCache(t *testing.T) {
	if err := testDao.TrimDMCache(context.TODO(), 1, dm.Oid, 1); err != nil {
		t.Error(err)
	}
}
