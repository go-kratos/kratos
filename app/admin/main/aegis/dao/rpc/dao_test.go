package rpc

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"go-common/app/admin/main/aegis/conf"
	api "go-common/app/service/main/account/api"
	relmod "go-common/app/service/main/relation/model"
	uprpc "go-common/app/service/main/up/api/v1"

	"github.com/golang/mock/gomock"
)

var (
	d    *Dao
	cntx context.Context
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.aegis-admin")
		flag.Set("conf_token", "cad913269be022e1eb8c45a8d5408d78")
		flag.Set("tree_id", "60977")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/aegis-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	d = New(conf.Conf)
	cntx = context.TODO()
	os.Exit(m.Run())
}

var (
	_Mid  = int64(0)
	_Mids = []int64{0, 1, 2}
)

func WithDao(t *testing.T, f func(d *Dao)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)

		accMock := NewMockAccRPC(mockCtrl)
		d.AccountClient = accMock
		accMock.EXPECT().Cards3(gomock.Any(), &api.MidsReq{Mids: _Mids}).Return(&api.CardsReply{Cards: map[int64]*api.Card{
			10086: {
				Mid: _Mid,
			},
		}}, nil).AnyTimes()

		relMock := NewMockRelationRPC(mockCtrl)
		d.relRPC = relMock
		relMock.EXPECT().Stats(gomock.Any(), &relmod.ArgMids{Mids: _Mids}).Return(map[int64]*relmod.Stat{
			10086: {
				Mid: 10086,
			},
		}, nil)

		f(d)
		mockCtrl.Finish()
	}
}

func WithMockAccount(t *testing.T, f func(d *Dao)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		accMock := NewMockAccRPC(mockCtrl)
		d.AccountClient = accMock
		accMock.EXPECT().Info3(gomock.Any(), &api.MidReq{Mid: _Mid}).Return(&api.InfoReply{Info: &api.Info{
			Mid: 10086,
		}}, nil).AnyTimes()
		accMock.EXPECT().Cards3(gomock.Any(), &api.MidsReq{Mids: _Mids}).Return(&api.CardsReply{Cards: map[int64]*api.Card{
			10086: {
				Mid: _Mid,
			},
		}}, nil).AnyTimes()
		accMock.EXPECT().ProfileWithStat3(gomock.Any(), &api.MidReq{Mid: _Mid}).Return(&api.ProfileStatReply{
			Profile: &api.Profile{
				Mid: _Mid,
			},
		}, nil).AnyTimes()

		f(d)
	}
}

func WithMockRelation(t *testing.T, f func(d *Dao)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		relMock := NewMockRelationRPC(mockCtrl)
		d.relRPC = relMock
		relMock.EXPECT().Stats(gomock.Any(), &relmod.ArgMids{Mids: _Mids}).Return(map[int64]*relmod.Stat{
			10086: {
				Mid: 10086,
			},
		}, nil)

		f(d)
	}
}

func WithMockUp(t *testing.T, f func(d *Dao)) func() {
	return func() {
		upspecial := &uprpc.UpSpecial{GroupIDs: []int64{1}}
		upspecialreply := &uprpc.UpSpecialReply{UpSpecial: upspecial}
		upspecialsreply := &uprpc.UpsSpecialReply{UpSpecials: map[int64]*uprpc.UpSpecial{0: upspecial}}
		upgroupsreply := &uprpc.UpGroupsReply{UpGroups: map[int64]*uprpc.UpGroup{0: {ID: 0}}}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		upMock := NewMockUpRPC(mockCtrl)
		d.UpClient = upMock
		upMock.EXPECT().UpSpecial(gomock.Any(), &uprpc.UpSpecialReq{Mid: _Mid}).Return(upspecialreply, nil)
		upMock.EXPECT().UpsSpecial(gomock.Any(), &uprpc.UpsSpecialReq{Mids: _Mids}).Return(upspecialsreply, nil)
		upMock.EXPECT().UpGroups(gomock.Any(), &uprpc.NoArgReq{}).Return(upgroupsreply, nil)
		f(d)
	}
}
