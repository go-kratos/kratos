package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/member/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRealnameList(t *testing.T) {
	Convey("Realname list", t, func() {
		var (
			mids     = []int64{46333}
			cardType = 0
			country  = 0
			opName   = ""
			tsFrom   = int64(1528023084)
			tsTo     = time.Now().Unix()
			state    = 2
			pn, ps   = 1, 20
		)
		list, total, err := d.RealnameMainList(context.Background(), mids, cardType, country, opName, tsFrom, tsTo, state, pn, ps, false)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
		So(total, ShouldNotBeNil)
	})
}

func TestRealnameReason(t *testing.T) {
	Convey("Realname reason", t, func() {
		var (
			list = []string{
				"+1s",
				"蛤",
				"苟利国家",
				"19260817",
			}
		)

		err := d.UpdateRealnameReason(context.Background(), list)
		So(err, ShouldBeNil)

		list2, total, err := d.RealnameReasonList(context.Background())
		So(err, ShouldBeNil)
		So(total, ShouldEqual, len(list))
		So(list2, ShouldResemble, list)
	})
}

func TestRealnameApplyCount(t *testing.T) {
	Convey("Realname apply count", t, func() {
		var (
			mid = int64(1)
		)
		count, err := d.RealnameApplyCount(context.Background(), mid)
		So(err, ShouldBeNil)
		So(count, ShouldNotBeNil)
	})
}

func TestRealnameApply(t *testing.T) {
	Convey("Realname apply", t, func() {
		var (
			id = int64(1)
		)
		apply, err := d.RealnameMainApply(context.Background(), id)
		So(err, ShouldBeNil)
		So(apply, ShouldNotBeNil)
	})
}

func TestRealnameApplyUpdate(t *testing.T) {
	Convey("Realname apply update", t, func() {
		var (
			id     = 1
			state  = 2
			opname = "ut"
			opid   = int64(233)
			optime = time.Now()
			remark = "ut_reason"
		)
		err := d.UpdateRealnameMainApply(context.Background(), id, state, opname, opid, optime, remark)
		So(err, ShouldBeNil)
	})
}

func TestRealnameAlipayApply(t *testing.T) {
	Convey("Realname alipay apply", t, func() {
		var (
			id = int64(1)
		)
		apply, err := d.RealnameAlipayApply(context.Background(), id)
		So(err, ShouldBeNil)
		t.Log(apply)
	})
}

func TestRealnameUpdateAlipayApply(t *testing.T) {
	Convey("Realname update alipay apply", t, func() {
		var (
			id = int64(1)
		)
		err := d.UpdateRealnameAlipayApply(context.Background(), id, 1, "someone", 2, "ut")
		So(err, ShouldBeNil)
	})
}

func TestUpdateRealnameInfo(t *testing.T) {
	Convey("Realname update realname info", t, func() {
		var (
			id = int64(1)
		)
		err := d.UpdateRealnameInfo(context.Background(), id, 2, "ut")
		So(err, ShouldBeNil)
	})
}

func TestAddRealnameIMG(t *testing.T) {
	Convey("AddRealnameIMG", t, func() {
		err := d.AddRealnameIMG(context.Background(), &model.DBRealnameApplyIMG{IMGData: "testing"})
		So(err, ShouldBeNil)
	})
}

func TestAddRealnameApply(t *testing.T) {
	Convey("AddRealnameApply", t, func() {
		err := d.AddRealnameApply(context.Background(), &model.DBRealnameApply{MID: 1})
		So(err, ShouldBeNil)
	})
}

func TestBatchRealnameInfo(t *testing.T) {
	Convey("BatchRealnameInfo", t, func() {
		res, err := d.BatchRealnameInfo(context.Background(), []int64{1, 2, 3})
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestRejectRealnameMainApply(t *testing.T) {
	Convey("RejectRealnameMainApply", t, func() {
		err := d.RejectRealnameMainApply(context.Background(), 1, "admin", 1, "test")
		So(err, ShouldBeNil)
	})
}

func TestRejectRealnameAlipayApply(t *testing.T) {
	Convey("RejectRealnameAlipayApply", t, func() {
		err := d.RejectRealnameAlipayApply(context.Background(), 1, "admin", 1, "test")
		So(err, ShouldBeNil)
	})
}

func TestAddOldRealnameIMG(t *testing.T) {
	Convey("AddOldRealnameIMG", t, func() {
		err := d.AddOldRealnameIMG(context.Background(), &model.DeDeIdentificationCardApplyImg{IMGData: "/test"})
		So(err, ShouldBeNil)
	})
}

func TestAddOldRealnameApply(t *testing.T) {
	Convey("AddOldRealnameApply", t, func() {
		err := d.AddOldRealnameApply(context.Background(), &model.DeDeIdentificationCardApply{MID: 1})
		So(err, ShouldBeNil)
	})
}
