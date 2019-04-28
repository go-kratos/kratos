package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetXMLSegCache(t *testing.T) {
	Convey("", t, func() {
		err := testDao.SetXMLSegCache(context.TODO(), model.SubTypeVideo, 1221, 1, 1, []byte("test"))
		So(err, ShouldBeNil)
	})
}

func TestXMLSegCache(t *testing.T) {
	Convey("", t, func() {
		_, err := testDao.XMLSegCache(context.TODO(), model.SubTypeVideo, 1221, 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestDurationCache(t *testing.T) {
	var (
		oid int64 = 1508
		c         = context.TODO()
	)
	Convey("", t, func() {
		_, err := testDao.DurationCache(c, oid)
		So(err, ShouldBeNil)
	})
}

func TestSetDurationCache(t *testing.T) {
	var (
		oid      int64 = 1508
		duration int64 = 9031 * 1000
		c              = context.TODO()
	)
	Convey("", t, func() {
		err := testDao.SetDurationCache(c, oid, duration)
		So(err, ShouldBeNil)
	})
}

func TestSetDMSegCache(t *testing.T) {
	Convey("set dm segment cache, error should be nil", t, func() {
		dmseg := new(model.DMSeg)
		dmseg.Elems = append(dmseg.Elems, &model.Elem{Content: "dm msg"})
		err := testDao.SetDMSegCache(c, 1, 1221, 1, 1, dmseg)
		So(err, ShouldBeNil)
	})
}

func TestDMSegCache(t *testing.T) {
	Convey("get dm segment cache", t, func() {
		dmseg, err := testDao.DMSegCache(c, 1, 1221, 1, 1)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v", dmseg)
	})
}
