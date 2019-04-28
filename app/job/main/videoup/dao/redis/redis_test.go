package redis

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_AddArcClick(t *testing.T) {
	var (
		err error
		c   = context.TODO()
	)
	Convey("AddArcClick", t, WithDao(func(d *Dao) {
		err = d.AddArcClick(c, 1, 1)
		So(err, ShouldBeNil)
	}))
}

//ArcClick
func Test_ArcClick(t *testing.T) {
	var (
		err error
		c   = context.TODO()
	)
	Convey("ArcClick", t, WithDao(func(d *Dao) {
		_, err = d.ArcClick(c, 1)
		So(err, ShouldBeNil)
	}))
}

//DelFilename
func Test_DelFilename(t *testing.T) {
	var (
		err error
		c   = context.TODO()
	)
	Convey("DelFilename", t, WithDao(func(d *Dao) {
		err = d.DelFilename(c, "1")
		So(err, ShouldBeNil)
	}))
}

//SetMonitorCache
func Test_SetMonitorCache(t *testing.T) {
	var (
		err error
		c   = context.TODO()
	)
	Convey("SetMonitorCache", t, WithDao(func(d *Dao) {
		_, err = d.SetMonitorCache(c, 1)
		So(err, ShouldBeNil)
	}))
}
