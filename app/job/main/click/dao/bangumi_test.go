package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_LoadAllBangumi(t *testing.T) {
	Convey("LoadAllBangumi", t, func() {
		etam, err := d.LoadAllBangumi(context.TODO())
		So(err, ShouldBeNil)
		So(etam, ShouldNotBeNil)
		for epid, aid := range etam {
			Printf("epid:%d,aid:%d\n", epid, aid)
		}
	})
}

func Test_IsLegal(t *testing.T) {
	Convey("IsLegal", t, func() {
		isLegal, err := d.IsLegal(context.TODO(), 11696747, 157927, 1)
		So(err, ShouldBeNil)
		Println(isLegal)
	})
}
