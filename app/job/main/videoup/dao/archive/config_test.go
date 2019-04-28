package archive

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_RoundEndConf(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("RoundEndConf", t, WithDao(func(d *Dao) {
		_, err = d.RoundEndConf(c)
		So(err, ShouldBeNil)
	}))
}

func Test_FansConf(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int64
	)
	Convey("FansConf", t, WithDao(func(d *Dao) {
		sub, err = d.FansConf(c)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeEmpty)
	}))
}

func Test_RoundTypeConf(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub map[int16]struct{}
	)
	Convey("RoundTypeConf", t, WithDao(func(d *Dao) {
		sub, err = d.RoundTypeConf(c)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeNil)
	}))
}

func Test_AuditTypesConf(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub map[int16]struct{}
	)
	Convey("AuditTypesConf", t, WithDao(func(d *Dao) {
		sub, err = d.AuditTypesConf(c)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeNil)
	}))
}
func Test_ThresholdConf(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub map[int16]int
	)
	Convey("ThresholdConf", t, WithDao(func(d *Dao) {
		sub, err = d.ThresholdConf(c)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeNil)
	}))
}
