package weeklyhonor

import (
	"context"
	"fmt"
	"testing"

	whmdl "go-common/app/interface/main/creative/model/weeklyhonor"
	upgrpc "go-common/app/service/main/up/api/v1"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

var c = context.Background()

func TestWeeklyhonorUpCount(t *testing.T) {
	convey.Convey("UpCount", t, func(ctx convey.C) {
		var (
			mid = int64(1627855)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.UpCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWeeklyhonorUpActiveLists(t *testing.T) {
	convey.Convey("UpActivesList", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", WithMock(t, func(mockCtrl *gomock.Controller) {
			// mock
			mockUpClient := upgrpc.NewMockUpClient(mockCtrl)
			d.upClient = mockUpClient
			mockReq := upgrpc.UpListByLastIDReq{
				LastID: 0,
				Ps:     100,
			}
			mockReply := upgrpc.UpActivityListReply{
				UpActivitys: []*upgrpc.UpActivity{
					{Mid: 1},
					{Activity: 2},
				},
				LastID: 1,
			}
			mockUpClient.EXPECT().UpInfoActivitys(gomock.Any(), &mockReq).Return(&mockReply, nil)
			// test
			upActives, newId, err := d.UpActivesList(c, 0, 100)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upActives[0], convey.ShouldEqual, mockReply.UpActivitys[0])
				ctx.So(newId, convey.ShouldEqual, mockReply.LastID)
			})
		}))
	})
}

func TestWeeklyhonorGetUpSwitch(t *testing.T) {
	convey.Convey("GetUpSwitch", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.Background()
			mid = int64(75379101)
		)
		convey.Convey("When everything gose positive", func(ctx convey.C) {
			mockUpClient := upgrpc.NewMockUpClient(mockCtrl)
			d.upClient = mockUpClient
			mockReq := upgrpc.UpSwitchReq{
				Mid:  mid,
				From: fromWeeklyHonor,
			}
			mockState := whmdl.HonorUnSub
			mockUpClient.EXPECT().UpSwitch(gomock.Any(), &mockReq).Return(&upgrpc.UpSwitchReply{State: mockState}, nil)
			state, err := d.GetUpSwitch(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldEqual, mockState)
			})
		})
		convey.Convey("When return err", func(ctx convey.C) {
			mockUpClient := upgrpc.NewMockUpClient(mockCtrl)
			d.upClient = mockUpClient
			mockReq := upgrpc.UpSwitchReq{
				Mid:  mid,
				From: fromWeeklyHonor,
			}
			mockUpClient.EXPECT().UpSwitch(gomock.Any(), &mockReq).Return(nil, fmt.Errorf("mock err"))
			_, err := d.GetUpSwitch(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	}))
}
