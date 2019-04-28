package account

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	accapi "go-common/app/service/main/account/api"
	relaMdl "go-common/app/service/main/relation/model"
	relation "go-common/app/service/main/relation/rpc/client"
	"go-common/library/ecode"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative")
		flag.Set("conf_token", "96b6a6c10bb311e894c14a552f48fef8")
		flag.Set("tree_id", "2305")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	d.fastClient.SetTransport(gock.DefaultTransport)
	return r
}

func TestIdentifyInfo(t *testing.T) {
	Convey("IdentifyInfo", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.TODO()
			err error
			mid = int64(27515256)
			ip  = "127.0.0.1"
			ret int
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		arg := &accapi.MidReq{
			Mid: mid,
		}
		mock.EXPECT().Profile3(gomock.Any(), arg).Return(nil, ecode.CreativeAccServiceErr)
		ret, err = d.IdentifyInfo(c, mid, 1, ip)
		So(err, ShouldNotBeNil)
		So(ret, ShouldNotBeNil)
	}))
}

func TestMidByName(t *testing.T) {
	Convey("IdentifyInfo", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c    = context.TODO()
			err  error
			ret  int64
			name = "iamname"
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		arg := &accapi.NamesReq{
			Names: []string{name},
		}
		mock.EXPECT().InfosByName3(gomock.Any(), arg).Return(nil, ecode.CreativeAccServiceErr)
		ret, err = d.MidByName(c, name)
		So(err, ShouldNotBeNil)
		So(ret, ShouldBeZeroValue)
	}))
}

func TestInfos(t *testing.T) {
	Convey("IdentifyInfo", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c    = context.TODO()
			err  error
			mids = []int64{2089809}
			ip   = "127.0.0.1"
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidsReq{
			Mids: mids,
		}
		mock.EXPECT().Infos3(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		_, err = d.Infos(c, mids, ip)
		So(err, ShouldNotBeNil)
	}))
}

func TestProfile(t *testing.T) {
	Convey("Profile", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.TODO()
			err error
			mid = int64(27515256)
			ip  = "127.0.0.1"
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
func TestProfileWithStat(t *testing.T) {
	Convey("ProfileWithStat", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.TODO()
			err error
			mid = int64(27515256)
			p   *accapi.ProfileStatReply
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidReq{
			Mid: mid,
		}
		mock.EXPECT().ProfileWithStat3(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		p, err = d.ProfileWithStat(c, mid)
		So(err, ShouldNotBeNil)
		So(p, ShouldBeNil)
	}))
}
func TestCard(t *testing.T) {
	Convey("Card", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.TODO()
			err error
			mid = int64(27515256)
			ip  = "127.0.0.1"
			ret *accapi.Card
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		mockReq := &accapi.MidReq{
			Mid: mid,
		}
		mock.EXPECT().Card3(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		ret, err = d.Card(c, mid, ip)
		So(err, ShouldNotBeNil)
		So(ret, ShouldBeNil)
	}))
}

func TestRichRelation(t *testing.T) {
	Convey("RichRelation", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c     = context.TODO()
			err   error
			owner = int64(27515256)
			mids  = []int64{2089809}
			ip    = "127.0.0.1"
			ret   map[int64]int32
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		arg := &accapi.RichRelationReq{
			Owner: owner,
			Mids:  mids,
		}
		mock.EXPECT().RichRelations3(gomock.Any(), arg).Return(nil, ecode.CreativeAccServiceErr)
		ret, err = d.RichRelation(c, owner, mids, ip)
		So(err, ShouldNotBeNil)
		So(ret, ShouldBeNil)
	}))
}

func TestRelationFollowers(t *testing.T) {
	Convey("RelationFollowers", t, func(ctx C) {
		var (
			c   = context.TODO()
			err error
			mid = int64(2089809)
			ip  = "127.0.0.1"
			ret map[int64]int32
		)
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.rela), "Followers",
			func(_ *relation.Service, _ context.Context, _ *relaMdl.ArgMid) (res []*relaMdl.Following, err error) {
				return nil, ecode.CreativeAccServiceErr
			})
		defer mock.Unpatch()
		ret, err = d.RelationFollowers(c, mid, ip)
		ctx.Convey("RelationFollowers", func(ctx C) {
			ctx.So(err, ShouldNotBeNil)
			ctx.So(ret, ShouldBeNil)
		})
	})
}

func TestShouldFollow(t *testing.T) {
	Convey("ShouldFollow", t, func(ctx C) {
		var (
			c    = context.TODO()
			err  error
			mid  = int64(2089809)
			fids = []int64{2089809}
			ip   = "127.0.0.1"
			ret  []int64
		)
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.rela), "Relations",
			func(_ *relation.Service, _ context.Context, _ *relaMdl.ArgRelations) (res map[int64]*relaMdl.Following, err error) {
				res = make(map[int64]*relaMdl.Following)
				res[2089809] = &relaMdl.Following{
					Attribute: 0,
				}
				return res, nil
			})
		defer mock.Unpatch()
		ret, err = d.ShouldFollow(c, mid, fids, ip)
		ctx.Convey("RelationFollowers", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(ret, ShouldNotBeNil)
		})
	})
}

func TestSwitchPhoneRet(t *testing.T) {
	var (
		new, old, identify int
		err                error
	)
	new = 1
	Convey("switchPhoneRet", t, func(ctx C) {
		old = d.switchPhoneRet(new)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx C) {
			ctx.So(old, ShouldBeZeroValue)
		})
	})
	identify = 1
	Convey("CheckIdentify", t, func(ctx C) {
		err = d.CheckIdentify(identify)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx C) {
			ctx.So(err, ShouldEqual, ecode.UserCheckInvalidPhone)
		})
	})
}

func TestDao_Followers(t *testing.T) {
	Convey("Followers", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.TODO()
			err error
			mid = int64(2089809)
			fid = int64(2089809)
			ip  = "127.0.0.1"
		)
		mock := accapi.NewMockAccountClient(mockCtrl)
		d.acc = mock
		arg := &accapi.RelationReq{
			Owner: mid,
			Mid:   fid,
		}
		mock.EXPECT().Relation3(gomock.Any(), arg).Return(nil, ecode.CreativeAccServiceErr)
		_, err = d.Followers(c, mid, []int64{fid}, ip)
		Convey("Followers", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	}))
}

func TestDao_Relations(t *testing.T) {
	Convey("Relations", t, func(ctx C) {
		var (
			c    = context.TODO()
			err  error
			mid  = int64(2089809)
			fids = []int64{2089809}
			ip   = "127.0.0.1"
			ret  map[int64]int
		)
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.rela), "Relations",
			func(_ *relation.Service, _ context.Context, _ *relaMdl.ArgRelations) (res map[int64]*relaMdl.Following, err error) {
				return nil, ecode.CreativeAccServiceErr
			})
		defer mock.Unpatch()
		ret, err = d.Relations(c, mid, fids, ip)
		ctx.Convey("RelationFollowers", func(ctx C) {
			ctx.So(err, ShouldNotBeNil)
			ctx.So(ret, ShouldBeNil)
		})
	})
}

func TestDao_Relations2(t *testing.T) {
	Convey("Relations2", t, func(ctx C) {
		var (
			c    = context.TODO()
			err  error
			mid  = int64(2089809)
			fids = []int64{2089809}
			ip   = "127.0.0.1"
			ret  map[int64]int
		)
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.rela), "Relations",
			func(_ *relation.Service, _ context.Context, _ *relaMdl.ArgRelations) (res map[int64]*relaMdl.Following, err error) {
				return nil, ecode.CreativeAccServiceErr
			})
		defer mock.Unpatch()
		ret, err = d.Relations2(c, mid, fids, ip)
		ctx.Convey("RelationFollowers", func(ctx C) {
			ctx.So(err, ShouldNotBeNil)
			ctx.So(ret, ShouldBeNil)
		})
	})
}
func WithMock(t *testing.T, f func(mock *gomock.Controller)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		f(mockCtrl)
	}
}
