package drawimg

import (
	"image"
	"image/draw"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func TestDrawimgParseFont(t *testing.T) {
	convey.Convey("ParseFont", t, func(ctx convey.C) {
		var (
			b = []byte("")
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := ParseFont(b)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgPt(t *testing.T) {
	convey.Convey("Pt", t, func(ctx convey.C) {
		var (
			x = int(0)
			y = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := Pt(x, y)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgPointToFixed(t *testing.T) {
	convey.Convey("PointToFixed", t, func(ctx convey.C) {
		var (
			x = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := di.c.PointToFixed(x)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgdrawContour(t *testing.T) {
	convey.Convey("drawContour", t, func(ctx convey.C) {
		var (
			ps = []truetype.Point{
				{
					X:     2,
					Y:     2,
					Flags: 1,
				},
				{
					X:     2,
					Y:     2,
					Flags: 1,
				},
			}
			dx fixed.Int26_6
			dy fixed.Int26_6
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.drawContour(ps, dx, dy)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgrecalc(t *testing.T) {
	convey.Convey("recalc", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.recalc()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgSetDPI(t *testing.T) {
	convey.Convey("SetDPI", t, func(ctx convey.C) {
		var (
			dpi = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.SetDPI(dpi)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgSetFont(t *testing.T) {
	convey.Convey("SetFont", t, func(ctx convey.C) {
		var (
			f = &truetype.Font{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.SetFont(f)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgSetFontSize(t *testing.T) {
	convey.Convey("SetFontSize", t, func(ctx convey.C) {
		var (
			fontSize = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.SetFontSize(fontSize)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgSetHinting(t *testing.T) {
	convey.Convey("SetHinting", t, func(ctx convey.C) {
		var (
			hinting font.Hinting
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.SetHinting(hinting)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgSetDst(t *testing.T) {
	convey.Convey("SetDst", t, func(ctx convey.C) {
		var (
			dst draw.Image
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.SetDst(dst)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgSetSrc(t *testing.T) {
	convey.Convey("SetSrc", t, func(ctx convey.C) {
		var (
			src image.Image
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.SetSrc(src)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgSetClip(t *testing.T) {
	convey.Convey("SetClip", t, func(ctx convey.C) {
		var (
			clip image.Rectangle
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c.SetClip(clip)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgNewContext(t *testing.T) {
	convey.Convey("NewContext", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := NewContext()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgInt26_6ToString(t *testing.T) {
	convey.Convey("Int26_6ToString", t, func(ctx convey.C) {
		var (
			x fixed.Int26_6
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := Int26_6ToString(x)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgDrawString(t *testing.T) {
	convey.Convey("DrawString", t, func(ctx convey.C) {
		var (
			s = "mock"
			p fixed.Point26_6
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			di.c = &c
			monkeyLoad()
			monkeyClear()
			monkeyRasterizer()
			p1, err := di.c.DrawString(s, p)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
