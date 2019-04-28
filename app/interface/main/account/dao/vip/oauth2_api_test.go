package vip

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/account/model"

	. "github.com/smartystreets/goconvey/convey"
)

//  go test  -test.v -test.run TestDaoOAuth2
func TestDaoOAuth2(t *testing.T) {
	Convey("TestDaoOAuth2", t, func() {
		data, err := d.OAuth2ByCode(context.Background(), &model.ArgAuthCode{
			Code: "778b91c87a43e4fdab6c647010b77566",
		})
		if fmt.Sprintf("%v", err) != "dao oauth2 userinfo: -907" {
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)
		}
		So(fmt.Sprintf("%v", err) == "dao oauth2 userinfo: -907", ShouldBeTrue)
	})
}
