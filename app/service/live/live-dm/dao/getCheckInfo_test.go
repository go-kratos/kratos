package dao

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	fansMedalService "go-common/app/service/live/fans_medal/api/liverpc"
	"go-common/app/service/live/live-dm/conf"
	liveUserService "go-common/app/service/live/live_user/api/liverpc"
	roomService "go-common/app/service/live/room/api/liverpc"
	userextService "go-common/app/service/live/userext/api/liverpc"
	acctountService "go-common/app/service/main/account/api"
	filterService "go-common/app/service/main/filter/api/grpc/v1"
	spyService "go-common/app/service/main/spy/api"
)

func init() {
	dir, _ := filepath.Abs("../cmd/test.toml")
	flag.Set("conf", dir)
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	// InitAPI()
	// InitGrpc(conf.Conf)
	UserExtServiceClient = userextService.New(getConf("userext"))
	LiveUserServiceClient = liveUserService.New(getConf("liveUser"))
	FansMedalServiceClient = fansMedalService.New(getConf("fansMedal"))
	RoomServiceClient = roomService.New(getConf("room"))
	ac, err = acctountService.NewClient(conf.Conf.AccClient)
	if err != nil {
		panic(err)
	}
	vipCli, err = newVipService(conf.Conf.XuserClent)
	if err != nil {
		panic(err)
	}
	SpyClient, err = spyService.NewClient(conf.Conf.SpyClient)
	if err != nil {
		panic(err)
	}
	FilterClient, err = filterService.NewClient(conf.Conf.FilterClient)
	if err != nil {
		panic(err)
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDMConf_Get
func TestDMConf_Get(t *testing.T) {
	dc := &DMConf{}
	if err := dc.Get(context.TODO(), 111, 222, conf.Conf); err != nil {
		t.Error("获取弹幕配置失败", err)
	}
	if dc.Color == 0 && dc.Length == 0 && dc.Mode == 0 {
		t.Error("获取弹幕配置失败， 返回值错误")
	}
	fmt.Println("##### Mode: ", dc.Mode)
	fmt.Println("##### Color: ", dc.Color)
	fmt.Println("##### Length: ", dc.Length)
}

//TODO 未测试
//group=fat1 DEPLOY_ENV=uat go test -run TestUserInfo_Get
func TestUserInfo_Get(t *testing.T) {
	u := &UserInfo{}
	if err := u.Get(context.TODO(), 110000232); err != nil {
		t.Error(err)
	}
	if u.UserLever == 0 && u.UserScore == 0 {
		t.Error("返回值错误")
	}
	fmt.Println("#### UserLever: ", u.UserLever)
	fmt.Println("#### UserScore: ", u.UserScore)
	fmt.Println("### Usercolor: ", u.ULevelColor)
}

//DEPLOY_ENV=uat go test -run TestUserInfo_GetVipInfo
func TestUserInfo_GetVipInfo(t *testing.T) {
	u := &UserInfo{}
	if err := u.GetVipInfo(context.TODO(), 2); err != nil {
		t.Error("获取老爷失败: ", err)
	}
	fmt.Println("#### VIP: ", u.Vip)
	fmt.Println("### SVIP: ", u.Svip)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestUserInfo_GetPrivilegeType
func TestUserInfo_GetPrivilegeType(t *testing.T) {
	u := &UserInfo{}
	if err := u.GetPrivilegeType(context.TODO(), 10799340, 6810576); err != nil {
		t.Error("返回值错误: ", err)
	}
	fmt.Println("PrivilegeType", u.PrivilegeType)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestUserInfo_IsRoomAdmin
func TestUserInfo_IsRoomAdmin(t *testing.T) {
	u := &UserInfo{}
	if err := u.IsRoomAdmin(context.TODO(), 1877309, 5392); err != nil {
		t.Error("返回值错误: ", err)
	}
	fmt.Println("IsRoomAdmin->", u.RoomAdmin)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestUserInfo_GetFansMedal
func TestUserInfo_GetFansMedal(t *testing.T) {
	m := &FansMedalInfo{}
	if err := m.GetFansMedal(context.TODO(), 83940); err != nil {
		t.Error("获取粉丝勋章失败: ", err)
	}
	fmt.Println("#####RUID: ", m.RUID)
	fmt.Println("#####MedalLevel: ", m.MedalLevel)
	fmt.Println("#####MedalName: ", m.MedalName)
	fmt.Println("#####MColor: ", m.MColor)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestRoomConf_Get
func TestRoomConf_Get(t *testing.T) {
	r := &RoomConf{}
	if err := r.Get(context.TODO(), 1016); err != nil {
		t.Error("获取房间配置失败: ", err)
	}
	fmt.Println("RoomID->", r.RoomID)
	fmt.Println("UID->", r.UID)
	fmt.Println("RoomShield->", r.RoomShield)
	fmt.Println("Anchor->", r.Anchor)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestUserBindInfo_Get
func TestUserBindInfo_Get(t *testing.T) {
	u := &UserBindInfo{}
	if err := u.Get(context.TODO(), 222); err != nil {
		t.Error("获取用户绑定信息失败: ", err)
	}
	fmt.Println("Identification->", u.Identification)
	fmt.Println("MobileVerify->", u.MobileVerify)
	fmt.Println("Uname->", u.Uname)
	fmt.Println("URank->", u.URank)
}

//DEPLOY_ENV=uat go test -run TestGerUserScore
func TestGerUserScore(t *testing.T) {
	u := &UserScore{}
	if err := u.GetUserScore(context.TODO(), 111); err != nil {
		t.Error("获取用户真实分失败:", err)
	}
	fmt.Println("###### UserScore:", u.UserScore)
}

//缺少souce值
//DEPLOY_ENV=uat go test -run TestGetMsgScore
func TestGetMsgScore(t *testing.T) {
	u := &UserScore{}
	if err := u.GetMsgScore(context.TODO(), "fuck"); err != nil {
		t.Error("获取真实分失败:", err)
	}
	fmt.Println("MsgLeve->", u.MsgLevel)
	fmt.Println("MsgAI=>", u.MsgAI)
}
