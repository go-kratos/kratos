package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/admin/main/workflow/model/search"

	"github.com/smartystreets/goconvey/convey"
)

func TestSetGroupResult(t *testing.T) {
	convey.Convey("SetGroupResult", t, func() {
		g := new(model.Group)
		err := s.dao.ORM.Table("workflow_group").Where("business=1").Find(g).Error
		convey.So(err, convey.ShouldBeNil)

		err = s.SetGroupResult(context.Background(), &param.GroupResParam{Oid: g.Oid, Business: g.Business, State: int8(1), AdminID: int64(1), Reason: "hhh"})
		convey.So(err, convey.ShouldBeNil)
		ng := new(model.Group)

		err = s.dao.ORM.Table("workflow_group").Where("business=1 and oid=?", g.Oid).Find(ng).Error
		convey.So(err, convey.ShouldBeNil)
		convey.So(ng.State, convey.ShouldEqual, int8(1))
	})
}

func TestChallListCommon(t *testing.T) {
	convey.Convey("ChallListCommon", t, func() {
		cPage, err := s.ChallListCommon(context.Background(), &search.ChallSearchCommonCond{
			Order: "ctime", Sort: "desc", PN: 1, PS: 50})
		time.Sleep(3 * time.Second)
		convey.So(err, convey.ShouldBeNil)
		convey.So(cPage.Page.Total, convey.ShouldBeGreaterThanOrEqualTo, int32(1))
	})
}

func TestChallList(t *testing.T) {
	convey.Convey("ChallList", t, func() {
		cPage, err := s.ChallList(context.Background(), &search.ChallSearchCommonCond{IDs: []int64{1},
			PN: 1, PS: 50, Order: "id", Sort: "desc"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(cPage.Page.Total, convey.ShouldBeGreaterThanOrEqualTo, int32(1))
	})
}

func TestChallDetail(t *testing.T) {
	convey.Convey("ChallDetail", t, func() {
		chall, err := s.ChallDetail(context.Background(), int64(1))
		convey.So(err, convey.ShouldBeNil)
		convey.So(chall.Cid, convey.ShouldEqual, int32(1))
	})
}
