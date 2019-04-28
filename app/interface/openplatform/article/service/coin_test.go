package service

import (
	"context"
	"testing"

	coin "go-common/app/service/main/coin/api/gorpc"
	coinMdl "go-common/app/service/main/coin/model"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Coin(t *testing.T) {
	mid := int64(1)
	aid := int64(1)
	Convey("get data", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		mock := coin.NewMockRPC(mockCtrl)
		s.coinRPC = mock
		arg := &coinMdl.ArgCoinInfo{Mid: mid, Aid: aid, AvType: 2}
		mock.EXPECT().ArchiveUserCoins(gomock.Any(), arg).Return(&coinMdl.ArchiveUserCoins{Multiply: 10}, nil)
		res, err := s.Coin(context.TODO(), mid, aid, "")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 10)
	}))
}
