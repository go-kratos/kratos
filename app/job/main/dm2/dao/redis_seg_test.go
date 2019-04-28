package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDMIDCache(t *testing.T) {
	var (
		tp  int32 = 1
		oid int64 = 1508
		cnt int64 = 26
		num int64 = 1
		c         = context.TODO()
	)
	Convey("", t, func() {
		_, err := testDao.DMIDCache(c, tp, oid, cnt, num, 100)
		So(err, ShouldBeNil)
	})
}

func TestAddDMIDCache(t *testing.T) {
	Convey("", t, func() {
		err := testDao.AddDMIDCache(c, 1, 1508, 26, 1, 1233333333333)
		So(err, ShouldBeNil)
	})
}

func TestIdxContentCacheV2(t *testing.T) {
	var (
		tp    int32 = 1
		oid   int64 = 1508
		c           = context.TODO()
		dmids       = []int64{2355015081, 2356915089}
	)
	Convey("", t, func() {
		elems, missed, err := testDao.IdxContentCacheV2(c, tp, oid, dmids)
		So(err, ShouldBeNil)
		t.Logf("missed dmid:%v", missed)
		t.Logf("elems:%+v", elems)
	})
}

func TestXMLToElem(t *testing.T) {
	Convey("convert xml tag to elem struct", t, func() {
		s := []byte(`<d p="1,1,1,1,11,111,11,1,23123123">弹幕内容</d>`)
		elem, err := testDao.xmlToElem(s)
		So(err, ShouldBeNil)
		t.Logf("%+v", elem)
	})
}
