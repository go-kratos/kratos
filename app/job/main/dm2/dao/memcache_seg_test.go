package dao

import (
	"testing"

	"go-common/app/job/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDelXMLSegCache(t *testing.T) {
	Convey("check delete segment xml cache,error should be nil", t, func() {
		err := testDao.DelXMLSegCache(c, 1, 1221, 1, 1)
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
