package http

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"go-common/app/admin/main/apm/conf"
	bm "go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserRuleStates(t *testing.T) {
	fmt.Println("========UserRuleStates tests========")
	Convey("test userAuth", t, func() {
		var params = url.Values{}
		ctx := userRequest("GET", params)
		userRuleStates(&ctx)
		data, err := json.Marshal(ctx)
		fmt.Println(string(data))
		So(err, ShouldBeNil)
	})
}
func TestUserApplies(t *testing.T) {
	fmt.Println("========UserAppliestests========")
	Convey("test userApplies", t, func() {
		var params = url.Values{}
		ctx := userRequest("GET", params)
		userApplies(&ctx)
		data, err := json.Marshal(ctx)
		fmt.Println(string(data))
		So(err, ShouldBeNil)
	})
}

// func TestUserApply(t *testing.T) {
// 	fmt.Println("========UserApply tests========")
// 	Convey("test userApply", t, func() {
//		var params = url.Values{}
// 		params.Set("item", "APP_AUTH_VIEW,APP_EDIT")
// 		ctx := userRequest("POST", params)
// 		userApply(&ctx)
// 		data, err := json.Marshal(ctx)
// 		fmt.Println(string(data))
// 		So(err, ShouldBeNil)
// 	})
// }

// func TestUserAudit(t *testing.T) {
// 	fmt.Println("========UserAudit tests========")
// 	Convey("test userAudit", t, func() {
//		var params = url.Values{}
// 		params.Set("id", "2")
// 		params.Set("status", "2")
// 		ctx := userRequest("POST", params)
// 		userAudit(&ctx)
// 		data, err := json.Marshal(ctx)
// 		fmt.Println(string(data))
// 		So(err, ShouldBeNil)
// 	})
// }

func userRequest(method string, params url.Values) (ctx bm.Context) {
	flag.Set("conf", "../cmd/apm-admin-test.toml")
	conf.Init()
	//apmSvc = service.New(conf.Conf)
	path := "/"
	req, _ := http.NewRequest(method, path, strings.NewReader(params.Encode()))
	req.Header.Set("X-BACKEND-BILI-REAL-IP", "xxx")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "_AJSESSIONID=0eac92f0621b11e8b877522233007f8a; username=chengxing; sven-apm=fa7128be35363eb623e4dd5d611daf1f3273316cc09a66655db0e79f616d4f53")

	req.ParseForm()
	ctx = bm.Context{
		Context: context.TODO(),
		Request: req,
	}
	ctx.Set("username", "chengxing")
	ctx.Request.ParseForm()
	return
}
