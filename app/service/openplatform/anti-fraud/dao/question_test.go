package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_qid      = 1527480165645
	_qbid     = 1527233672941
	_page     = 0
	_pagesize = 10
)

func TestGetQusInfo(t *testing.T) {
	Convey("TestBankSearch: ", t, func() {
		res, err := d.GetQusInfo(context.TODO(), _qid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetQusList(t *testing.T) {
	Convey("TestBankSearch: ", t, func() {
		res, err := d.GetQusList(context.TODO(), _page, _pagesize, _qbid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetQusIds(t *testing.T) {
	Convey("TestGetQusIds: ", t, func() {
		res, err := d.GetQusIds(context.TODO(), _qbid)
		length := len(res)
		So(err, ShouldBeNil)
		So(length, ShouldBeGreaterThan, 0)
	})
}

func TestGetQusCount(t *testing.T) {
	Convey("TestGetQusCount: ", t, func() {
		cnt, err := d.GetQusCount(context.TODO(), _qbid)
		So(err, ShouldBeNil)
		So(cnt, ShouldBeGreaterThan, 0)
	})
}

func TestGetAnswerList(t *testing.T) {
	Convey("TestGetAnswerList: ", t, func() {
		cnt, err := d.GetAnswerList(context.TODO(), _qid)
		So(err, ShouldBeNil)
		So(cnt, ShouldNotBeNil)
	})
}

func TestGDelQus(t *testing.T) {
	Convey("TestGetAnswerList: ", t, func() {
		cnt, err := d.DelQus(context.TODO(), 1527241107344)
		So(err, ShouldBeNil)
		So(cnt, ShouldNotBeNil)
	})
}
