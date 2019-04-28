package service

import (
	"context"
	"testing"

	"go-common/app/service/openplatform/abtest/model"

	. "github.com/smartystreets/goconvey/convey"
)

var ab = model.AB{ID: 1828, Group: 1, Status: 1, Name: "n@m3"}

func TestSetGroupCache(t *testing.T) {
	var src = []*model.AB{&ab}
	Convey("TestSetGroupCache: ", t, func() {
		err := svr.setGroupCache(context.TODO(), 1, src, 123)
		So(err, ShouldBeNil)
		d, ok := svr.readGroupCache(context.TODO(), 1)
		So(ok, ShouldBeTrue)
		So(d[1828], ShouldEqual, &ab)
	})
}

func TestVersionIDListCache(t *testing.T) {
	var src = []*model.AB{&ab}
	svr.setGroupCache(context.TODO(), 1, src, 123)
	Convey("TestVersionIDListCache: ", t, func() {
		verList, err := svr.VersionIDListCache(context.TODO())
		So(err, ShouldBeNil)
		_, ok := verList[1]
		So(ok, ShouldBeTrue)
	})
}
