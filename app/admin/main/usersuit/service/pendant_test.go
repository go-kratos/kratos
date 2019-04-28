package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/usersuit/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PendantInfoList(t *testing.T) {
	Convey("return sth", t, func() {
		arg := &model.ArgPendantGroupList{}
		res, pager, err := s.PendantInfoList(context.Background(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}

func Test_PendantInfoID(t *testing.T) {
	Convey("return sth", t, func() {
		pi, err := s.PendantInfoID(context.Background(), 11, 22)
		So(err, ShouldBeNil)
		So(pi, ShouldNotBeNil)
	})
}

func Test_PendantGroupID(t *testing.T) {
	Convey("return sth", t, func() {
		pg, err := s.PendantGroupID(context.Background(), 11)
		So(err, ShouldBeNil)
		So(pg, ShouldNotBeNil)
	})
}

func Test_PendantGroupList(t *testing.T) {
	Convey("return sth", t, func() {
		arg := &model.ArgPendantGroupList{}
		pgs, pager, err := s.PendantGroupList(context.Background(), arg)
		So(err, ShouldBeNil)
		So(pgs, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}

func Test_PendantGroupAll(t *testing.T) {
	Convey("return sth", t, func() {
		pgs, err := s.PendantGroupAll(context.Background())
		So(err, ShouldBeNil)
		So(pgs, ShouldNotBeNil)
	})
}

func Test_PendantInfoAllOnSale(t *testing.T) {
	Convey("return sth", t, func() {
		pis, err := s.PendantInfoAllNoPage(context.Background())
		So(err, ShouldBeNil)
		So(pis, ShouldNotBeNil)
	})
}

func Test_AddPendantInfo(t *testing.T) {
	Convey("return sth", t, func() {
		arg := &model.ArgPendantInfo{
			PID:           111,
			GID:           22,
			Name:          "222",
			Image:         "sdada",
			ImageModel:    "dasdada",
			Rank:          11,
			Status:        1,
			IntegralPrice: 22,
		}
		err := s.AddPendantInfo(context.Background(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_UpPendantInfo(t *testing.T) {
	Convey("return sth", t, func() {
		arg := &model.ArgPendantInfo{
			PID:           111,
			GID:           22,
			Name:          "222",
			Image:         "sdada",
			ImageModel:    "dasdada",
			Rank:          11,
			Status:        1,
			IntegralPrice: 22,
		}
		err := s.UpPendantInfo(context.Background(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_UpPendantGroupStatus(t *testing.T) {
	Convey("return sth", t, func() {
		err := s.UpPendantGroupStatus(context.Background(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_UpPendantInfoStatus(t *testing.T) {
	Convey("return sth", t, func() {
		err := s.UpPendantInfoStatus(context.Background(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_AddPendantGroup(t *testing.T) {
	Convey("return sth", t, func() {
		arg := &model.ArgPendantGroup{
			GID:    1,
			Name:   "2121",
			Rank:   11,
			Status: 1,
		}
		err := s.AddPendantGroup(context.Background(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_UpPendantGroup(t *testing.T) {
	Convey("return sth", t, func() {
		arg := &model.ArgPendantGroup{
			GID:    1,
			Name:   "2121",
			Rank:   11,
			Status: 1,
		}
		err := s.UpPendantGroup(context.Background(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_PendantOrders(t *testing.T) {
	Convey("return sth", t, func() {
		arg := &model.ArgPendantOrder{}
		pos, pager, err := s.PendantOrders(context.Background(), arg)
		So(err, ShouldBeNil)
		So(pos, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}
