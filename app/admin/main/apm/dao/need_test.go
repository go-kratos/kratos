package dao

import (
	"reflect"
	"testing"

	"go-common/app/admin/main/apm/model/need"

	"github.com/bouk/monkey"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoNeedInfoAdd(t *testing.T) {
	convey.Convey("NeedInfoAdd", t, func() {
		arg := &need.NAddReq{
			Title:   "wwe",
			Content: "sds",
		}
		err := d.NeedInfoAdd(arg, "fengshanshan")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoNeedInfoList(t *testing.T) {
	convey.Convey("NeedInfoList", t, func() {
		arg := &need.NListReq{
			Status: 1,
			Ps:     10,
			Pn:     1,
		}
		res, err := d.NeedInfoList(arg)
		t.Logf("res:%+v", res)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoNeedInfoCount(t *testing.T) {
	convey.Convey("NeedInfoCount", t, func() {
		arg := &need.NListReq{
			Status: 1,
		}
		count, err := d.NeedInfoCount(arg)
		t.Logf("count:%+v", count)
		convey.So(err, convey.ShouldBeNil)
		convey.So(count, convey.ShouldNotBeNil)
	})
}

func TestDaoneedInfoCondition(t *testing.T) {
	convey.Convey("needInfoCondition", t, func() {
		arg := &need.NListReq{}
		p1 := d.needInfoCondition(arg)
		t.Logf("condition:%+v", p1)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestDaoGetNeedInfo(t *testing.T) {
	convey.Convey("GetNeedInfo", t, func() {
		r, err := d.GetNeedInfo(97)
		t.Logf("GetNeedInfo:%+v", r)
		convey.So(err, convey.ShouldBeNil)
		convey.So(r, convey.ShouldNotBeNil)
	})
}

func TestDaoNeedInfoEdit(t *testing.T) {
	convey.Convey("NeedInfoEdit", t, func() {
		arg := &need.NEditReq{
			Content: "dsada",
			Title:   "fsd",
			ID:      28,
		}
		err := d.NeedInfoEdit(arg)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoNeedVerify(t *testing.T) {
	convey.Convey("NeedVerify", t, func() {
		v := &need.NVerifyReq{
			ID:     28,
			Status: 2,
		}
		err := d.NeedVerify(v)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoLikeCountsAdd(t *testing.T) {
	convey.Convey("LikeCountsAdd", t, func() {
		v := &need.Likereq{
			ReqID:    148,
			LikeType: 1,
		}
		err := d.LikeCountsStats(v, 1, 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoGetVoteInfo(t *testing.T) {
	convey.Convey("GetVoteInfo", t, func() {
		var (
			db = &gorm.DB{
				Error: nil,
			}
			v = &need.Likereq{
				ReqID:    148,
				LikeType: 1,
			}
		)
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.DB), "Find", func(_ *gorm.DB, _ interface{}, _ ...interface{}) *gorm.DB {
			return db
		})
		defer guard.Unpatch()
		res, err := d.GetVoteInfo(v, "fengshanshan")
		t.Logf("res:%+v", res)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoUpdateVoteInfo(t *testing.T) {
	convey.Convey("UpdateVoteInfo", t, func() {
		v := &need.Likereq{
			ReqID:    30,
			LikeType: 2,
		}
		err := d.UpdateVoteInfo(v, "fengshanshan")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoVoteInfoList(t *testing.T) {
	convey.Convey("VoteInfoList", t, func() {
		arg := &need.Likereq{
			ReqID:    11,
			LikeType: 2,
		}
		res, err := d.VoteInfoList(arg)
		t.Logf("res:%+v", res)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoVoteInfoCounts(t *testing.T) {
	convey.Convey("VoteInfoCounts", t, func() {
		arg := &need.Likereq{
			ReqID:    11,
			LikeType: 1,
		}
		count, err := d.VoteInfoCounts(arg)
		t.Logf("count:%+v", count)
		convey.So(err, convey.ShouldBeNil)
		convey.So(count, convey.ShouldNotBeNil)
	})
}
