package service

import (
	"context"
	coinm "go-common/app/service/main/coin/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testCoinMid int64 = 130
)

// go test -test.v -test.run TestPutCoinInfo
func TestPutCoinInfo(t *testing.T) {
	Convey("TestPutCoinInfo put coin info", t, WithService(func(s *Service) {
		So(s.PutCoinInfo(context.TODO(), &coinm.DataBus{
			Mid: testCoinMid,
		}), ShouldBeNil)
	}))
}
