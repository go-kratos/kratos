package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetBankBind(t *testing.T) {
	Convey("TestGetBankBind: ", t, func() {
		res, err := d.GetBankBind(context.TODO(), 1, 1, []string{"1111"}, true)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetBindBank(t *testing.T) {
	Convey("TestGetBindBank: ", t, func() {
		res, err := d.GetBindBank(context.TODO(), 1, 1, []string{"1111"})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestCountBindItem(t *testing.T) {
	Convey("TestCountBindItem: ", t, func() {
		res, err := d.CountBindItem(context.TODO(), _qbid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetBindItem(t *testing.T) {
	Convey("TestCountBindItem: ", t, func() {
		res, cnt, err := d.GetBindItem(context.TODO(), _qbid, _page, _pagesize)
		So(err, ShouldBeNil)
		So(cnt, ShouldNotBeNil)
		So(res, ShouldNotBeNil)
	})
}
