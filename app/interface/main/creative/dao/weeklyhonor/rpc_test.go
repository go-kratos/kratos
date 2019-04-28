package weeklyhonor

import (
	"context"
	"fmt"
	"testing"

	whmdl "go-common/app/interface/main/creative/model/weeklyhonor"
	up "go-common/app/service/main/up/api/v1"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

func TestWeeklyhonorChangeUpSwitch(t *testing.T) {
	convey.Convey("ChangeUpSwitch", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c     = context.Background()
			mid   = int64(75379101)
			state = whmdl.HonorSub
		)
		convey.Convey("When everything gose positive", func(ctx convey.C) {
			mockUpClient := up.NewMockUpClient(mockCtrl)
			d.upClient = mockUpClient
			mockReq := up.UpSwitchReq{
				Mid:   mid,
				From:  fromWeeklyHonor,
				State: state,
			}
			mockUpClient.EXPECT().SetUpSwitch(gomock.Any(), &mockReq).Return(&up.NoReply{}, nil)
			err := d.ChangeUpSwitch(c, mid, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		convey.Convey("When return err", func(ctx convey.C) {
			mockUpClient := up.NewMockUpClient(mockCtrl)
			d.upClient = mockUpClient
			mockReq := up.UpSwitchReq{
				Mid:   mid,
				From:  fromWeeklyHonor,
				State: state,
			}
			mockUpClient.EXPECT().SetUpSwitch(gomock.Any(), &mockReq).Return(nil, fmt.Errorf("mock err"))
			err := d.ChangeUpSwitch(c, mid, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	}))
}

func TestWeeklyhonorGetUpSwitch(t *testing.T) {
	convey.Convey("GetUpSwitch", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.Background()
			mid = int64(75379101)
		)

		convey.Convey("When everything gose positive", func(ctx convey.C) {
			mockUpClient := up.NewMockUpClient(mockCtrl)
			d.upClient = mockUpClient
			mockReq := up.UpSwitchReq{
				Mid:  mid,
				From: fromWeeklyHonor,
			}
			mockState := whmdl.HonorUnSub
			mockUpClient.EXPECT().UpSwitch(gomock.Any(), &mockReq).Return(&up.UpSwitchReply{State: mockState}, nil)
			state, err := d.GetUpSwitch(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldEqual, mockState)
			})
		})
		convey.Convey("When return err", func(ctx convey.C) {
			mockUpClient := up.NewMockUpClient(mockCtrl)
			d.upClient = mockUpClient
			mockReq := up.UpSwitchReq{
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
