package drawimg

import (
	"image"
	"image/draw"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
	"github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	di = DrawImg{
		size:         100,
		CanvasWidth:  20,
		CanvasHeight: 20,
		File:         "",
		txtWidth:     20,
		srcImg:       image.NewAlpha(imgRectangle),
		Canvas:       &image.NRGBA{},
		c:            &c,
		f:            &truetype.Font{},
	}
	c = Context{
		r:        nil,
		f:        &truetype.Font{},
		glyphBuf: truetype.GlyphBuf{},
		clip:     image.Rectangle{},
		dst:      nil,
		src:      nil,
		fontSize: 0,
		dpi:      0,
		scale:    0,
		hinting:  0,
		cache:    [1024]cacheEntry{},
	}
	imgRectangle = image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: 1,
			Y: 1,
		},
	}
)

func TestDrawimgNewDrawImg(t *testing.T) {
	convey.Convey("NewDrawImg", t, func(ctx convey.C) {
		var (
			fontfile = ""
			size     = int(10)
		)
		monkeyReadFile([]byte{}, nil)
		monkeyTrueTypeParser(&truetype.Font{}, nil)
		monkeyFreeTypeSetFont()
		monkeySetFontSize()
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			w := NewDrawImg(fontfile, size)
			ctx.Convey("Then w should not be nil.", func(ctx convey.C) {
				ctx.So(w, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgreadFont(t *testing.T) {
	convey.Convey("readFont", t, func(ctx convey.C) {
		var (
			path = ""
			size = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			monkeyReadFile([]byte{}, nil)
			monkeyTrueTypeParser(&truetype.Font{}, nil)
			monkeyFreeTypeSetFont()
			monkeySetFontSize()
			f, err := di.readFont(path, size)
			ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(f, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgnewCanvas(t *testing.T) {
	convey.Convey("newCanvas", t, func(ctx convey.C) {
		var (
			width  = int(0)
			height = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := di.newCanvas(width, height)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgfillColor(t *testing.T) {
	convey.Convey("fillColor", t, func(ctx convey.C) {
		var (
			r = int32(0)
			g = int32(0)
			b = int32(0)
			a = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := di.fillColor(r, g, b, a)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgtextWidth(t *testing.T) {
	convey.Convey("textWidth", t, func(ctx convey.C) {
		var (
			text = "12"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			monkeyFontBox(fixed.Point26_6{}, nil)
			err := di.textWidth(text)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDrawimgpt(t *testing.T) {
	convey.Convey("pt", t, func(ctx convey.C) {
		var (
			x = int(0)
			y = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := di.pt(x, y)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgsetFont(t *testing.T) {
	convey.Convey("setFont", t, func(ctx convey.C) {
		var (
			text    = ""
			dstRgba = &image.NRGBA{}
			fsrc    image.Image
			pt      fixed.Point26_6
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			monkeyDrawString(fixed.Point26_6{}, nil)
			err := di.setFont(text, dstRgba, fsrc, pt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDrawimgcomposite(t *testing.T) {
	convey.Convey("composite", t, func(ctx convey.C) {
		var (
			dstCanvas = &image.NRGBA{}
			src       image.Image
			isLeft    bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			monkeyDraw()
			//monkeybounds(imgRectangle)
			di.composite(dstCanvas, src, isLeft)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgDraw(t *testing.T) {
	convey.Convey("Draw", t, func(ctx convey.C) {
		var (
			text     = "123"
			savepath = ""
			isLeft   bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			//monkeybounds(imgRectangle)
			monkeyFreeTypeSetFont()
			err := di.Draw(text, savepath, isLeft)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func monkeyReadFile(file []byte, err error) {
	monkey.Patch(ioutil.ReadFile, func(_ string) ([]byte, error) {
		return file, err
	})
}

func monkeyTrueTypeParser(font *truetype.Font, err error) {
	monkey.Patch(truetype.Parse, func(_ []byte) (*truetype.Font, error) {
		return font, err
	})
}

func monkeyFreeTypeSetFont() {
	monkey.PatchInstanceMethod(reflect.TypeOf(di.c), "SetFont", func(_ *Context, _ *truetype.Font) {})
}

func monkeySetFontSize() {
	monkey.PatchInstanceMethod(reflect.TypeOf(di.c), "SetFontSize", func(_ *Context, _ float64) {})
}

func monkeyDrawString(p fixed.Point26_6, err error) {
	monkey.PatchInstanceMethod(reflect.TypeOf(di.c), "DrawString", func(_ *Context, _ string, _ fixed.Point26_6) (fixed.Point26_6, error) {
		return p, err
	})
}

func monkeyDraw() {
	monkey.Patch(draw.Draw, func(_ draw.Image, _ image.Rectangle, _ image.Image, _ image.Point, _ draw.Op) {})
}

func monkeyFontBox(p fixed.Point26_6, err error) {
	monkey.PatchInstanceMethod(reflect.TypeOf(di.c), "FontBox", func(_ *Context, _ string, _ fixed.Point26_6) (fixed.Point26_6, error) {
		return p, err
	})
}

func monkeyLoad() {
	monkey.PatchInstanceMethod(reflect.TypeOf(&di.c.glyphBuf), "Load", func(_ *truetype.GlyphBuf, _ *truetype.Font, _ fixed.Int26_6, _ truetype.Index, _ font.Hinting) error {
		return nil
	})
}

func monkeyClear() {
	monkey.PatchInstanceMethod(reflect.TypeOf(di.c.r), "Clear", func(_ *raster.Rasterizer) {})
}

func monkeyRasterizer() {
	monkey.PatchInstanceMethod(reflect.TypeOf(di.c.r), "Rasterize", func(_ *raster.Rasterizer, _ raster.Painter) {})
}
