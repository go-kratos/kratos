package service

import (
	"context"
	"testing"

	"go-common/app/service/main/msm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScope(t *testing.T) {
	var (
		serviceTreeID = int64(22)
		appTreeID     = int64(33)
		appAuth       = "qwertyuiop"
		sign          = "123455667"
		c             = context.TODO()
	)
	scope := &model.Scope{
		AppTreeID:   appTreeID,
		RPCMethods:  append([]string{"RPC.Info", "RPC.Scope"}),
		HTTPMethods: append([]string{"/test/info", "/test/add"}),
		Quota:       5000,
		Sign:        signature(appTreeID, serviceTreeID, appAuth),
	}
	real := make(map[int64]*model.Scope)
	real[serviceTreeID] = scope
	Convey("err should return nil and id greater than zero", t, func() {
		res, err := svr.ServiceScopes(c, serviceTreeID)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, real)
	})
	Convey("compare sign is equal", t, func() {
		b := svr.CheckSign(14212, sign)
		So(b, ShouldBeFalse)
	})
}
