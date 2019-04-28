package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testAccountMid int64 = 152
	testExp        int64 = 10
)

func TestSaveFigure(t *testing.T) {
	once.Do(startService)
	s.AccountReg(context.TODO(), 7593623)
}

// go test -test.v -test.run TestAccountExp
func TestAccountExp(t *testing.T) {
	Convey("TestSaveFigure put account exp", t, WithService(func(s *Service) {
		So(s.AccountReg(context.TODO(), 10), ShouldBeNil)
	}))
	Convey("TestAccountExp put account exp", t, WithService(func(s *Service) {
		So(s.AccountExp(context.TODO(), testAccountMid, testExp), ShouldBeNil)
	}))
	Convey("TestAccountViewVideo put account exp", t, WithService(func(s *Service) {
		So(s.AccountViewVideo(context.TODO(), testAccountMid), ShouldBeNil)
	}))
}
