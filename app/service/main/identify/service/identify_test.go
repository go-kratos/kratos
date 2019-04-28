package service

import (
	"context"
	"testing"

	"go-common/app/service/main/identify/api/grpc"
	"go-common/app/service/main/identify/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SetCache(t *testing.T) {
	Convey("test set cache ", t, func() {
		once.Do(startService)
		var (
			accessKey = "11111111"
			res       = &model.IdentifyInfo{
				Mid:     1523454,
				Expires: 111111,
				Csrf:    "1111",
			}
			err error
			c   = context.Background()
		)
		err = s.SetCache(c, accessKey, res)
		So(err, ShouldBeNil)
	})
}

func TestService_AccessCookie(t *testing.T) {
	Convey("test GetCookieInfo ", t, func() {
		once.Do(startService)
		var (
			cookie = "SESSDATA=11111111;buvid3=fasr435"
			res    *model.IdentifyInfo
			err    error
			c      = context.Background()
		)
		res, err = s.GetCookieInfo(c, cookie)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_AccessToken(t *testing.T) {
	Convey("test GetTokenInfo", t, func() {
		once.Do(startService)
		var (
			token = &v1.GetTokenInfoReq{}
			c     = context.Background()
			res   *model.IdentifyInfo
			err   error
		)
		token.Token = "11111111"
		token.Buvid = "fert32434"
		res, err = s.GetTokenInfo(c, token)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_LoginLog(t *testing.T) {
	Convey("test login log", t, func() {

		once.Do(startService)
		var (
			mid    = 1244545
			ip     = "15.1.1.1"
			ipPort = "2434"
			buvid  = "2134435fderf"
		)
		s.loginLog(int64(mid), ip, ipPort, buvid)
	})
}

func TestService_DelCache(t *testing.T) {
	Convey("test del cache", t, func() {
		once.Do(startService)
		var (
			accessKey = "11111111"
			err       error
			c         = context.Background()
		)
		err = s.DelCache(c, accessKey)
		So(err, ShouldBeNil)
	})
}
