package dao

import (
	"context"
	"testing"
)

func TestPubStat(t *testing.T) {
	once.Do(startService)
	if err := d.UpdateAccountExp(context.TODO(), 7593623, 120.00); err != nil {
		t.Errorf("d.UpdateAccountExp(%d) error(%v)", 7593623, err)
	}
}

func TestIncArchiveViews(t *testing.T) {
	once.Do(startService)
	if err := d.IncArchiveViews(context.TODO(), 7593623); err != nil {
		t.Errorf("error(%v)", err)
	}
}
