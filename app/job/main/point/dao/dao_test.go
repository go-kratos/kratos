package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/point/conf"
	"go-common/app/job/main/point/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	mid int64 = 7593623
	d   *Dao
	c   = context.TODO()
)

func init() {
	dir, _ := filepath.Abs("../cmd/point-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func TestAddPoint(t *testing.T) {
	Convey("TestAddPoint", t, func() {
		p := &model.VipPoint{
			Mid:          111,
			PointBalance: 1,
			Ver:          1,
		}
		_, err := d.AddPoint(c, p)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdatePoint
func TestUpdatePoint(t *testing.T) {
	Convey("TestUpdatePoint", t, func() {
		p := &model.VipPoint{
			Mid:          111,
			PointBalance: 2,
			Ver:          1,
		}
		oldp := &model.VipPoint{
			Ver: 1,
		}
		_, err := d.UpdatePoint(c, p, oldp.Ver)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestChangeHistory
func TestChangeHistory(t *testing.T) {
	Convey("TestChangeHistory", t, func() {
		h := &model.VipPointChangeHistory{
			Mid:          111,
			Point:        2,
			OrderID:      "wqwqe22112",
			ChangeType:   10,
			ChangeTime:   xtime.Time(time.Now().Unix()),
			RelationID:   "11",
			PointBalance: 111,
		}
		_, err := d.AddPointHistory(c, h)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDeletePointCache
func TestDeletePointCache(t *testing.T) {
	Convey("TestDeletePointCache", t, func() {
		err := d.DelPointInfoCache(c, mid)
		So(err, ShouldBeNil)
	})
}
