package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpWithdraw(t *testing.T) {
	Convey("UpWithdraw", t, WithService(func(s *Service) {
		var (
			mids        = []int64{int64(1101)}
			isDeleted   = 1
			from, limit = 0, 1000
		)
		res, _, err := s.UpWithdraw(context.Background(), mids, isDeleted, from, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_UpWithdrawExport(t *testing.T) {
	Convey("UpWithdrawExport", t, WithService(func(s *Service) {
		var (
			mids        = []int64{int64(1101)}
			isDeleted   = 1
			from, limit = 0, 1000
		)
		res, err := s.UpWithdrawExport(context.Background(), mids, isDeleted, from, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_GetUpWithdraw(t *testing.T) {
	Convey("GetUpWithdraw", t, WithService(func(s *Service) {
		var (
			query = ""
		)
		res, err := s.GetUpWithdraw(context.Background(), query)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_ListUpWithdraw(t *testing.T) {
	Convey("ListUpWithdraw", t, WithService(func(s *Service) {
		var (
			id    int64
			query = ""
			limit = 500
		)
		res, err := s.ListUpWithdraw(context.Background(), id, query, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_UpWithdrawStatis(t *testing.T) {
	Convey("UpWithdrawStatis", t, WithService(func(s *Service) {
		var (
			isDeleted = 0
		)
		_, err := s.UpWithdrawStatis(context.Background(), 1522512000, 1533052800, isDeleted)
		So(err, ShouldBeNil)
	}))
}

func Test_UpWithdrawDetail(t *testing.T) {
	Convey("UpWithdrawDetail", t, WithService(func(s *Service) {
		var (
			mid int64 = 50
		)
		res, err := s.UpWithdrawDetail(context.Background(), mid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_UpWithdrawDetailExport(t *testing.T) {
	Convey("UpWithdrawDetailExport", t, WithService(func(s *Service) {
		var (
			mid int64 = 50
		)
		res, err := s.UpWithdrawDetailExport(context.Background(), mid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
