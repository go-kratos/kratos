package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueConf(t *testing.T) {
	Convey("QueConf true", t, func() {
		d.QueConf(true)
	})
	Convey("QueConf false", t, func() {
		d.QueConf(false)
	})
}

func TestHeight(t *testing.T) {
	Convey("Height true", t, func() {
		textImgConf := d.QueConf(true)
		d.Height(textImgConf, "这里错了", 1)
	})
	Convey("Height false", t, func() {
		textImgConf := d.QueConf(false)
		d.Height(textImgConf, "这里对了", 1)
	})
}
func TestDaoBoard(t *testing.T) {
	Convey("Board", t, func() {
		res := d.Board(4)
		So(res, ShouldNotBeNil)
	})
}

func TestContext(t *testing.T) {
	Convey("Context", t, func() {
		board := d.Board(4)
		d.Context(board, "/data/conf/yahei.ttf")
	})
}

// func TestDrawQue(t *testing.T) {
// 	Convey("DrawQue", t, func() {
// 		quec := d.QueConf(true)
// 		imgh := d.Height(quec, "这里对了", 1)
// 		board := d.Board(imgh)
// 		imgc := d.Context(board, "/data/conf/yahei.ttf")
// 		pt := freetype.Pt(0, int(quec.Fontsize))
// 		d.DrawQue(imgc, "这里对了", quec, &pt)
// 		d.DrawAns(imgc, quec, [4]string{"A", "B", "C", "D"}, &pt)
// 	})
// }
