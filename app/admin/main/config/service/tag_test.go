package service

import (
	"context"
	"go-common/app/admin/main/config/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_TagByID(t *testing.T) {
	svr := svr(t)
	Convey("should tag by id", t, func() {
		res, err := svr.Tag(2)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestService_TagsByAppID(t *testing.T) {
	svr := svr(t)
	Convey("should get tags by app id", t, func() {
		res, err := svr.TagsByBuild("main.common-arch.msm-service", "dev", "shd", "server-1", 2, 1, 2888)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestService_UpdateTag(t *testing.T) {
	svr := svr(t)
	Convey("should update tag", t, func() {
		err := svr.UpdateTag(context.Background(), 2888, "dev", "shd", "server-1", &model.Tag{ConfigIDs: "212,260", Mark: "test", Operator: "zjx"})
		So(err, ShouldBeNil)
	})
}
