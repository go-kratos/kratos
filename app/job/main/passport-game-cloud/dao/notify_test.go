package dao

import (
	"context"
	"testing"
)

func TestDao_NotifyGame(t *testing.T) {
	once.Do(startDao)
	mid := int64(4780461)
	ak := "3cf80530cccafaa9ed675d8a493c1a89#tx"
	action := "changePwd"
	if err := d.NotifyGame(context.TODO(), mid, ak, action); err != nil {
		t.Errorf("dao.NotifyGame(%d, %s, %s) error(%v)", mid, ak, action, err)
		t.FailNow()
	}
}
