package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ModifyByPGC(t *testing.T) {
	var (
		c = context.TODO()
	)
	attrs := make(map[uint]int32, 7)
	attrs[archive.AttrBitJumpURL] = 1
	attrs[archive.AttrBitAllowBp] = 1
	attrs[archive.AttrBitIsBangumi] = 1
	attrs[archive.AttrBitIsMovie] = 1
	attrs[archive.AttrBitBadgepay] = 1
	attrs[archive.AttrBitIsPGC] = 1
	attrs[archive.AttrBitLimitArea] = 0
	Convey("ModifyByPGC", t, WithService(func(s *Service) {
		err := svr.ModifyByPGC(c, 1, 1, attrs, "")
		So(err, ShouldBeNil)
	}))
}

func TestService_LockByPGC(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("LockByPGC", t, WithService(func(s *Service) {
		err := svr.LockByPGC(c, 12345)
		So(err, ShouldBeNil)
	}))
}
