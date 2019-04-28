package account

import (
	"context"
	"reflect"
	"testing"

	relaMdl "go-common/app/service/main/relation/model"
	"go-common/app/service/main/relation/rpc/client"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"

	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"

	"github.com/golang/mock/gomock"
)

func WithMock(t *testing.T, f func(mock *gomock.Controller)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		f(mockCtrl)
	}
}

func TestAccountProfile(t *testing.T) {
	Convey("1", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			ip  = "127.0.0.1"
			err error
			p   *accapi.Profile
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidReq{
			Mid: mid,
		}
		mock.EXPECT().Profile3(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		p, err = d.Profile(c, mid, ip)
		So(err, ShouldNotBeNil)
		So(p, ShouldBeNil)
	}))
}

func TestDao_Cards(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2089809)
		ip  = "127.0.0.1"
		err error
	)
	Convey("Cards", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidsReq{
			Mids: []int64{mid},
		}
		res := &accapi.CardsReply{}
		mock.EXPECT().Cards3(gomock.Any(), mockReq).Return(res, nil)
		_, err = d.Cards(c, []int64{mid}, ip)
		So(err, ShouldBeNil)
	}))
}

func TestDao_Infos(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2089809)
		ip  = "127.0.0.1"
		err error
	)
	Convey("Infos", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidsReq{
			Mids: []int64{mid},
		}
		res := &accapi.InfosReply{}
		mock.EXPECT().Infos3(gomock.Any(), mockReq).Return(res, nil)
		_, err = d.Infos(c, []int64{mid}, ip)
		So(err, ShouldBeNil)
	}))
}

func TestDao_Relations(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2089809)
		ip  = "127.0.0.1"
		err error
	)
	Convey("Relations", t, func(ctx C) {
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.rela), "Relations",
			func(_ *relation.Service, _ context.Context, _ *relaMdl.ArgRelations) (res map[int64]*relaMdl.Following, err error) {
				res = make(map[int64]*relaMdl.Following)
				res[2089809] = &relaMdl.Following{}
				return res, nil
			})
		defer mock.Unpatch()
		_, err = d.Relations(c, mid, []int64{mid}, ip)
		So(err, ShouldBeNil)
	})
}
