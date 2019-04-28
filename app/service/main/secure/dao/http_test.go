package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDoubleCheck(t *testing.T) {
	Convey("TestDoubleCheck", t, func() {
		err := d.DoubleCheck(context.TODO(), 1)
		if err != nil {
			t.Errorf("test DoubleCheck err %v", err)
		}
	})
}
