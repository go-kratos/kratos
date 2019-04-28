package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_QueryFromUpInfo(t *testing.T) {
	var (
		accType      = 3
		states       = []int64{int64(1)}
		mid          = int64(1011)
		category     = 1
		signType     = 1
		nickname     = "hello"
		lower, upper = 0, 100
		from, limit  = 0, 1000
		sort         = "ctime"
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, _, err := s.QueryFromUpInfo(context.Background(), 0, accType, states, mid, category, signType, nickname, lower, upper, from, limit, sort)
		So(err, ShouldBeNil)
	}))
}

func Test_Reject(t *testing.T) {
	var (
		mids   = []int64{int64(1101)}
		reason = "reject"
		days   = 1
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.Reject(context.Background(), 0, mids, reason, days)
		So(err, ShouldBeNil)
	}))
}

func Test_Pass(t *testing.T) {
	var (
		mids = []int64{int64(1101)}
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.Pass(context.Background(), mids, 0)
		So(err, ShouldBeNil)
	}))
}

func Test_Dismiss(t *testing.T) {
	var (
		operator = "user"
		mid      = int64(1101)
		reason   = "dismiss"
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.Dismiss(context.Background(), operator, 0, 3, mid, reason)
		So(err, ShouldBeNil)
	}))
}

func Test_Forbid(t *testing.T) {
	var (
		operator = "user"
		mid      = int64(1101)
		reason   = "dismiss"
		days     = 5
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.Forbid(context.Background(), operator, 0, 3, mid, reason, days, 100)
		So(err, ShouldBeNil)
	}))
}

func Test_Recovery(t *testing.T) {
	var (
		mid = int64(1101)
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.Recovery(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func Test_UpdateUpAccountState(t *testing.T) {
	var (
		mid   = int64(1101)
		state = 3
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.UpdateUpAccountState(context.Background(), "up_info_video", mid, state)
		So(err, ShouldBeNil)
	}))
}

func Test_DeleteUp(t *testing.T) {
	var (
		mid = int64(1101)
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.DeleteUp(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func Test_Block(t *testing.T) {
	var (
		mid = int64(1101)
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.Block(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func Test_QueryFromBlocked(t *testing.T) {
	var (
		mid          = int64(1011)
		category     = 1
		nickname     = "hello"
		lower, upper = 0, 100
		from, limit  = 0, 1000
		sort         = "ctime"
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, _, err := s.QueryFromBlocked(context.Background(), mid, category, nickname, lower, upper, from, limit, sort)
		So(err, ShouldBeNil)
	}))
}

func Test_DeleteFromBlocked(t *testing.T) {
	var (
		mid = int64(1101)
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.DeleteFromBlocked(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func Test_DelUpAccount(t *testing.T) {
	var (
		mid = int64(1101)
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.DelUpAccount(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func Test_CreditRecords(t *testing.T) {
	var (
		mid = int64(1101)
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, err := s.CreditRecords(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func Test_RecoverCreditScore(t *testing.T) {
	var (
		id, mid int64 = 1, 1011
	)
	Convey("admins", t, WithService(func(s *Service) {
		err := s.RecoverCreditScore(context.Background(), 0, id, mid)
		So(err, ShouldBeNil)
	}))
}

func Test_ExportUps(t *testing.T) {
	var (
		accType      = 3
		states       = []int64{int64(1)}
		mid          = int64(1011)
		category     = 1
		signType     = 1
		nickname     = "hello"
		lower, upper = 0, 100
		from, limit  = 0, 1000
		sort         = "ctime"
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, err := s.ExportUps(context.Background(), 0, accType, states, mid, category, signType, nickname, lower, upper, from, limit, sort)
		So(err, ShouldBeNil)
	}))
}
