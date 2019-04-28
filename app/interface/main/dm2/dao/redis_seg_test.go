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
	_, err := testDao.DMIDCache(c, tp, oid, cnt, num, 100)
	if err != nil {
		t.Fatalf("d.DMIDCache(%d %d %d %d) error(%v)", tp, oid, cnt, num, err)
	}
}

func TestIdxContentCache(t *testing.T) {
	var (
		tp    int32 = 1
		oid   int64 = 1508
		c           = context.TODO()
		dmids       = []int64{2355015081, 2356915089}
	)
	res, missed, err := testDao.IdxContentCache(c, tp, oid, dmids)
	if err != nil {
		t.Errorf("d.IdxContentCache(%d %d %v) error(%v)", tp, oid, dmids, err)
		t.FailNow()
	}
	t.Logf("res:%s, missed:%v", res, missed)
}

func TestXMLToElem(t *testing.T) {
	Convey("convert xml tag to elem struct", t, func() {
		s := []byte(`<d p="1,1,1,1,11,111,11,1,23123123">弹幕内容</d>`)
		elem, err := testDao.xmlToElem(s)
		So(err, ShouldBeNil)
		t.Logf("%+v", elem)
	})
}
