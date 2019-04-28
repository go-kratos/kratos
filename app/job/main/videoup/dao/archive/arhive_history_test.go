package archive

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_ArchiveHistory(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("HistoryCount", t, WithDao(func(d *Dao) {
		_, err = d.HistoryCount(c, 23333)
		So(err, ShouldBeNil)
	}))
}

func Test_DelArcEditHistoryBefore(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		tm  time.Time
	)
	Convey("DelArcEditHistoryBefore", t, WithDao(func(d *Dao) {
		_, err = d.DelArcEditHistoryBefore(c, tm, 1)
		So(err, ShouldBeNil)
	}))
}

func Test_DelArcVideoEditHistoryBefore(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		tm  time.Time
	)
	Convey("DelArcVideoEditHistoryBefore", t, WithDao(func(d *Dao) {
		_, err = d.DelArcVideoEditHistoryBefore(c, tm, 1)
		So(err, ShouldBeNil)
	}))
}
