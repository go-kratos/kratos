package account

import (
	"context"
	"flag"
	accapi "go-common/app/service/main/account/api"
	"go-common/app/service/main/assist/conf"
	"go-common/library/ecode"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.assist-service")
		flag.Set("conf_token", "6e0dae2c95d90ff8d0da53460ef11ae8")
		flag.Set("tree_id", "2084")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/assist-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func WithMock(t *testing.T, f func(mock *gomock.Controller)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		f(mockCtrl)
	}
}

func TestIdentifyInfo(t *testing.T) {
	convey.Convey("1", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			ip  = "127.0.0.1"
			err error
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidReq{
			Mid: mid,
		}
		mock.EXPECT().Profile3(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		err = d.IdentifyInfo(c, mid, ip)
		convey.So(err, convey.ShouldNotBeNil)
	}))
	convey.Convey("2", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			ip  = "127.0.0.1"
			err error
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidReq{
			Mid: mid,
		}
		rpcRes := &accapi.ProfileReply{
			Profile: &accapi.Profile{
				Identification: 0,
				TelStatus:      2,
			},
		}
		mock.EXPECT().Profile3(gomock.Any(), mockReq).Return(rpcRes, nil)
		err = d.IdentifyInfo(c, mid, ip)
		convey.So(err, convey.ShouldNotBeNil)
	}))
}

func TestIsFollow(t *testing.T) {
	convey.Convey("1", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c         = context.Background()
			mid       = int64(2089809)
			assistMid = int64(11)
			err       error
			follow    bool
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.RelationReq{
			Mid:   assistMid,
			Owner: mid,
		}
		mock.EXPECT().Relation3(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		follow, err = d.IsFollow(c, mid, assistMid)
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(follow, convey.ShouldBeFalse)
	}))
}
func TestCard(t *testing.T) {
	convey.Convey("TestCard", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			err error
			res *accapi.Card
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidReq{
			Mid: mid,
		}
		mock.EXPECT().Card3(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		res, err = d.Card(c, mid, "")
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(res, convey.ShouldBeNil)
	}))
}
