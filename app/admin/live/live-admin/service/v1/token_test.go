package v1

import (
	"context"
	"flag"
	"go-common/app/admin/live/live-admin/dao"
	"testing"
	"time"

	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var tokenSrv *TokenService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	tokenSrv = NewTokenService(conf.Conf, dao.New(conf.Conf))
}

func TestUploadToken(t *testing.T) {
	Convey("TestTokenService", t, func() {
		req := &v1pb.NewTokenReq{
			Bucket:   "slive",
			Operator: "KC",
		}

		resp, err := tokenSrv.New(context.TODO(), req)
		So(err, ShouldBeNil)
		So(resp.Token, ShouldNotBeEmpty)

		ok := tokenSrv.dao.VerifyUploadToken(context.TODO(), resp.Token)
		So(ok, ShouldBeTrue)

		// Test token expiration.
		time.Sleep(time.Second * 11)
		ok = tokenSrv.dao.VerifyUploadToken(context.TODO(), resp.Token)
		So(ok, ShouldBeFalse)
	})
}
