package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go-common/app/service/main/vip/dao"
	"go-common/app/service/main/vip/model"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceRegisterOpenID(t *testing.T) {
	Convey(" TestServiceRegisterOpenID ", t, func() {
		res, err := s.RegisterOpenID(c, &model.ArgRegisterOpenID{
			Mid:   4,
			AppID: 32,
		})
		fmt.Println("res", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServicecreateOpenID(t *testing.T) {
	Convey(" TestServicecreateOpenID ", t, func() {
		res := createOpenID(2233, 30)
		fmt.Println("res", res)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceUserInfoByOpenID(t *testing.T) {
	Convey(" TestServiceUserInfoByOpenID ", t, func() {
		res, err := s.UserInfoByOpenID(c, &model.ArgUserInfoByOpenID{
			AppID:  32,
			OpenID: "d9e1aba143454a20efa63cf48e2f5903",
			IP:     "127.0.0.1",
		})
		fmt.Println("res", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestOpenAuthCallBack(t *testing.T) {
	Convey(" TestOpenAuthCallBack Mock OpenAuthCallBack ", t, func() {
		monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "EleOauthGenerateAccessToken", func(_ *dao.Dao,
			_ context.Context, _ *model.ArgEleAccessToken) (*model.EleAccessTokenResp, error) {
			return &model.EleAccessTokenResp{
				OpenID: "aaaaaxxxxxxx",
			}, nil
		})
		monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "EleUnionUpdateOpenID", func(_ *dao.Dao,
			_ context.Context, _ *model.ArgEleUnionUpdateOpenID) (*model.EleUnionUpdateOpenIDResp, error) {
			return &model.EleUnionUpdateOpenIDResp{
				Status: 1,
			}, nil
		})
		err := s.OpenAuthCallBack(c, &model.ArgOpenAuthCallBack{
			Mid:       4,
			AppID:     model.EleAppID,
			ThirdCode: "xxx",
		})
		So(err, ShouldBeNil)
	})
}
