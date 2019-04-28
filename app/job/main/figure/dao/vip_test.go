package dao

import (
	"context"
	"testing"
)

func TestUpdateVipStatus(t *testing.T) {
	once.Do(startService)
	d.UpdateVipStatus(context.TODO(), 7593623, 1)
}
