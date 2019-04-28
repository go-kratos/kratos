package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewHuaweis(t *testing.T) {
	Convey("test new huawei", t, WithDao(func(d *Dao) {
		cs := d.newHuaweiClients(123, "")
		So(len(cs), ShouldEqual, 0)
	}))
}
func TestNewOppoClients(t *testing.T) {
	Convey("test new huawei", t, WithDao(func(d *Dao) {
		cs := d.newOppoClients(123, "")
		So(len(cs), ShouldEqual, 0)
	}))
}

func TestNewJpushClients(t *testing.T) {
	Convey("test new huawei", t, WithDao(func(d *Dao) {
		cs := d.newJpushClients("123", "")
		So(len(cs), ShouldBeGreaterThan, 0)
	}))
}
