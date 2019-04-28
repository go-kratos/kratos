package service

import (
	"context"
	"testing"

	thumbupMdl "go-common/app/service/main/thumbup/model"
	thumbup "go-common/app/service/main/thumbup/rpc/client"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_isLike(t *testing.T) {
	var (
		mid   = int64(1)
		aid   = int64(2)
		state = int8(1)
	)
	Convey("get data", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		mock := thumbup.NewMockThumbupRPC(mockCtrl)
		s.thumbupRPC = mock
		arg := &thumbupMdl.ArgHasLike{Business: "article", MessageIDs: []int64{aid}, Mid: mid}
		mock.EXPECT().HasLike(gomock.Any(), arg).Return(map[int64]int8{aid: state}, nil)
		res, err := s.isLike(context.TODO(), mid, aid)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, state)
	}))
}

func Test_HadLikesByMid(t *testing.T) {
	var (
		mid   = int64(1)
		aids  = []int64{2, 2000}
		state = map[int64]int8{2: 1, 2000: 0}
	)
	Convey("get data", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		mock := thumbup.NewMockThumbupRPC(mockCtrl)
		s.thumbupRPC = mock
		arg := &thumbupMdl.ArgHasLike{Business: "article", MessageIDs: aids, Mid: mid}
		mock.EXPECT().HasLike(gomock.Any(), arg).Return(state, nil)
		res, err := s.HadLikesByMid(context.TODO(), mid, aids)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, state)
	}))
}
