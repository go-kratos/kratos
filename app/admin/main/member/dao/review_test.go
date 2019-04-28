package dao

import (
	"context"
	"testing"
	"time"

	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_MvArchivedFaceToPri(t *testing.T) {
	Convey("MvArchivedFaceToPri", t, func() {
		err := d.MvArchivedFaceToPriv(context.Background(), "/bfs/face/7e68723b9d3664ac3773c1f3c26d5e2bfabc0f23.jpg", "/bfs/face/7e68723b9d3664ac3773c1f3c26d5e2bfabc0f21.jpg", "sys", "")
		So(err, ShouldBeNil)
	})
}

func TestDao_IncrFaceReject(t *testing.T) {
	Convey("IncrFaceReject", t, func() {
		err := d.IncrFaceReject(context.Background(), 2)
		So(err, ShouldBeNil)
	})
}

func TestDao_IncrViolationCount(t *testing.T) {
	Convey("IncrViolationCount", t, func() {
		err := d.IncrViolationCount(context.Background(), 2)
		So(err, ShouldBeNil)
	})
}

func TestDao_FaceAutoPass(t *testing.T) {
	Convey("FaceAutoPass", t, func() {
		t := xtime.Time(time.Now().Unix())
		err := d.FaceAutoPass(context.Background(), []int64{1}, t)
		So(err, ShouldBeNil)
	})
}

func TestDao_prepareReviewRange(t *testing.T) {
	Convey("prepareReviewRange", t, func() {
		stime := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local).Unix()
		etime := time.Now().Unix()
		s, e, err := d.prepareReviewRange(context.Background(), xtime.Time(stime), xtime.Time(etime))
		So(err, ShouldBeNil)
		So(s, ShouldNotBeEmpty)
		So(e, ShouldNotBeEmpty)
	})
}

func TestDao_UpdateRemark(t *testing.T) {
	Convey("UpdateRemark", t, func() {
		err := d.UpdateRemark(context.Background(), 1, "12334")
		So(err, ShouldBeNil)
	})
}

func TestDao_QueuingFaceReviewsByTime(t *testing.T) {
	Convey("QueuingFaceReviewsByTime", t, func() {
		stime := time.Date(2018, 10, 1, 0, 0, 0, 0, time.Local).Unix()
		etime := time.Date(2018, 11, 1, 0, 0, 0, 0, time.Local).Unix()
		rws, err := d.QueuingFaceReviewsByTime(context.Background(), xtime.Time(stime), xtime.Time(etime))
		So(err, ShouldBeNil)
		// FIXME : UAT上全时间段查询不到该类数据
		So(rws, ShouldBeEmpty)
	})
}

func TestDao_ReviewByIDs(t *testing.T) {
	Convey("ReviewByIDs", t, func() {
		rws, err := d.ReviewByIDs(context.Background(), []int64{1, 2}, []int8{})
		So(err, ShouldBeNil)
		So(rws, ShouldNotBeEmpty)
	})
}

func TestDao_Reviews(t *testing.T) {
	Convey("Reviews", t, func() {
		rws, total, err := d.Reviews(context.Background(), 2231365, []int8{1}, []int8{0, 1, 2}, true, true, "", 1530542443, 1540910443, 1, 10)
		So(err, ShouldBeNil)
		So(total, ShouldNotBeNil)
		So(rws, ShouldNotBeNil)
	})
}

func TestDao_ReviewAudit(t *testing.T) {
	Convey("ReviewAudit", t, func() {
		err := d.ReviewAudit(context.Background(), []int64{2231365, 2231365}, 0, "test", "test")
		So(err, ShouldBeNil)
	})
}

func TestDao_Review(t *testing.T) {
	Convey("Review", t, func() {
		userPropertyReview, err := d.Review(context.Background(), 1)
		So(err, ShouldBeNil)
		So(userPropertyReview, ShouldNotBeNil)
	})
}

func TestDao_UpdateReviewFace(t *testing.T) {
	Convey("UpdateReviewFace", t, func() {
		err := d.UpdateReviewFace(context.Background(), 2231365, "face test")
		So(err, ShouldBeNil)
	})
}

func TestDao_AuditQueuingFace(t *testing.T) {
	Convey("Review", t, func() {
		err := d.AuditQueuingFace(context.Background(), 2231365, "face test", 0)
		So(err, ShouldBeNil)
	})
}
