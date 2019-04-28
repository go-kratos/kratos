package command

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCommandReadContributor(t *testing.T) {
	convey.Convey("readContributor", t, func(ctx convey.C) {
		content, _ := ioutil.ReadFile("../../CONTRIBUTORS.md")
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cc := readContributor(content)
			ctx.Convey("Then c should not be nil.", func(ctx convey.C) {
				ctx.So(fmt.Sprint(cc.Author), convey.ShouldEqual, `[muyang yubaihai wangweizhen wuwei]`)
				ctx.So(fmt.Sprint(cc.Owner), convey.ShouldEqual, `[muyang zhanglin]`)
				ctx.So(fmt.Sprint(cc.Reviewer), convey.ShouldEqual, `[muyang]`)
			})
		})
	})
}

func TestCommandHasBranch(t *testing.T) {
	convey.Convey("hasbranch", t, func(ctx convey.C) {
		var (
			branch  = "test"
			branchs = []string{"master", "test"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hasbranch(branch, branchs)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, true)
			})
		})
	})
}
