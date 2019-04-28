package dao

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	activityService "go-common/app/service/live/activity/api/liverpc"
	"go-common/app/service/live/live-dm/conf"
	rankdbService "go-common/app/service/live/rankdb/api/liverpc"
	rcService "go-common/app/service/live/rc/api/liverpc"
	userextService "go-common/app/service/live/userext/api/liverpc"
	acctountService "go-common/app/service/main/account/api"
	"go-common/library/net/metadata"
)

func init() {
	dir, _ := filepath.Abs("../cmd/test.toml")
	flag.Set("conf", dir)
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	UserExtServiceClient = userextService.New(getConf("userext"))
	ActivityServiceClient = activityService.New(getConf("activity"))
	RankdbServiceClient = rankdbService.New(getConf("rankdbService"))
	RcServiceClient = rcService.New(getConf("rc"))

	ac, err = acctountService.NewClient(conf.Conf.AccClient)
	if err != nil {
		panic(err)
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestUserInfo_GetUnameColor
func TestUserInfo_GetUnameColor(t *testing.T) {
	u := &UserInfo{}
	if err := u.GetUnameColor(context.TODO(), 28272030, 10004); err != nil {
		t.Error("获取用户昵称颜色失败: ", err)
	}
	fmt.Println("UnameColor->", u.UnameColor)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestUserInfo_GetSpeicalMedal
func TestUserInfo_GetSpeicalMedal(t *testing.T) {
	m := &FansMedalInfo{}
	if err := m.GetSpeicalMedal(context.TODO(), 111, 222); err != nil {
		t.Error("获取特殊勋章信息失败:", err)
	}
	fmt.Println("SpecialMedal->", m.SpecialMedal)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestUserInfo_GetUserLevelRank
func TestUserInfo_GetUserLevelRank(t *testing.T) {
	u := &UserInfo{}
	if err := u.GetUserLevelRank(context.TODO(), 111); err != nil {
		t.Error("获取用户等级RANK失败:", err)
	}
	fmt.Println("ULevelRank->", u.ULevelRank)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestCommentTitle_GetCommentTitle
func TestCommentTitle_GetCommentTitle(t *testing.T) {
	c := &CommentTitle{}

	ctx1 := metadata.NewContext(context.TODO(), metadata.MD{})
	if md, ok := metadata.FromContext(ctx1); ok {
		md[metadata.Mid] = 5200
	}
	if err := c.GetCommentTitle(ctx1); err != nil {
		t.Error("获取用户头衔失败:", err)
	}
	fmt.Println("OldTitle->", c.OldTitle)
	fmt.Println("Title->", c.Title)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestFansMedalInfo_GetMedalanchorName
func TestFansMedalInfo_GetMedalanchorName(t *testing.T) {
	f := &FansMedalInfo{}
	if err := f.GetMedalanchorName(context.TODO(), 222); err != nil {
		t.Error("获取勋章对应主播昵称错误:", err)
	}
	fmt.Println("RUName->", f.RUName)
}

//group=fat1 DEPLOY_ENV=uat go test -run TestUserInof_GetUserBubble
func TestUserInof_GetUserBubble(t *testing.T) {
	u := &UserInfo{}
	if err := u.GetUserBubble(context.TODO(), 1, 1, 1, 1); err != nil {
		t.Error("GetUserBubble调用失败")
	}
	if u.Bubble != 1 {
		t.Error("判断气泡失败 uid 1 roomid 1 bubble 1: bubble: ", u.Bubble)
	}
	fmt.Println("Bubble1->", u.Bubble)
	if err := u.GetUserBubble(context.TODO(), 1, 2, 1, 1); err != nil {
		t.Error("GetUserBubble调用失败")
	}
	if u.Bubble != 0 {
		t.Error("判断气泡失败 uid 1 roomid 2 bubble 1: bubble: ", u.Bubble)
	}
	fmt.Println("Bubble2->", u.Bubble)

}

// //group=qa01 DEPLOY_ENV=uat go test -run TestUserInfo_GetUserLevelColor
// func TestUserInfo_GetUserLevelColor(t *testing.T) {
// 	u := &UserInfo{}
// 	if err := u.GetUserLevelColor(52); err != nil {
// 		t.Error("返回值错误: ", err)
// 	}

// 	if u.ULevelColor != 16752445 {
// 		t.Error("51级以上颜色错误 16752445 ->", u.ULevelColor)
// 	}

// 	if err := u.GetUserLevelColor(42); err != nil {
// 		t.Error("返回值错误: ", err)
// 	}
// 	if u.ULevelColor != 16746162 {
// 		t.Error("51-41级颜色错误 16752445 ->", u.ULevelColor)
// 	}

// 	if err := u.GetUserLevelColor(32); err != nil {
// 		t.Error("返回值错误: ", err)
// 	}
// 	if u.ULevelColor != 10512625 {
// 		t.Error("41-31级颜色错误 10512625 ->", u.ULevelColor)
// 	}

// 	if err := u.GetUserLevelColor(22); err != nil {
// 		t.Error("返回值错误: ", err)
// 	}
// 	if u.ULevelColor != 5805790 {
// 		t.Error("31-21级颜色错误 16752445 ->", u.ULevelColor)
// 	}

// 	if err := u.GetUserLevelColor(12); err != nil {
// 		t.Error("返回值错误: ", err)
// 	}
// 	if u.ULevelColor != 6406234 {
// 		t.Error("21-11级颜色错误 16752445 ->", u.ULevelColor)
// 	}

// 	if err := u.GetUserLevelColor(2); err != nil {
// 		t.Error("返回值错误: ", err)
// 	}
// 	if u.ULevelColor != 9868950 {
// 		t.Error("0-11级颜色错误 16752445 ->", u.ULevelColor)
// 	}
// }
