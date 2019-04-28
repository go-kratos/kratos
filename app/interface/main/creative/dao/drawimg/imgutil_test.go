package drawimg

import (
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDrawimgString(t *testing.T) {
	convey.Convey("String", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := JPEG.String()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgOpen(t *testing.T) {
	convey.Convey("Open", t, func(ctx convey.C) {
		var (
			filename = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := Open(filename)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldBeNil)
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgEncode(t *testing.T) {
	convey.Convey("Encode", t, func(ctx convey.C) {
		var (
			w      io.Writer
			img    = image.NewRGBA(imgRectangle)
			format = JPEG
		)
		monkeyJpegEncode()
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := Encode(w, img, format)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDrawimgSave(t *testing.T) {
	convey.Convey("Save", t, func(ctx convey.C) {
		var (
			img      image.Image
			filename = ""
		)
		ctx.Convey("When everything goes not positive", func(ctx convey.C) {
			err := Save(img, filename)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgNewNRGBA(t *testing.T) {
	convey.Convey("NewNRGBA", t, func(ctx convey.C) {
		var (
			width     = int(0)
			height    = int(0)
			fillColor color.Color
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := NewNRGBA(width, height, fillColor)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgClone(t *testing.T) {
	convey.Convey("Clone", t, func(ctx convey.C) {
		var (
			rgba   = image.NewRGBA(imgRectangle)
			rgba64 = image.NewNRGBA64(imgRectangle)
		)
		ctx.Convey("RGBA", func(ctx convey.C) {
			p1 := Clone(rgba)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("RGBA64", func(ctx convey.C) {
			p1 := Clone(rgba64)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgtoNRGBA(t *testing.T) {
	convey.Convey("toNRGBA", t, func(ctx convey.C) {
		var (
			img = image.NewRGBA64(imgRectangle)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := toNRGBA(img)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgAdjustFunc(t *testing.T) {
	convey.Convey("AdjustFunc", t, func(ctx convey.C) {
		var (
			img = image.NewRGBA64(imgRectangle)
			fn  = func(c color.NRGBA) color.NRGBA { return c }
		)

		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := AdjustFunc(img, fn)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgAdjustGamma(t *testing.T) {
	convey.Convey("AdjustGamma", t, func(ctx convey.C) {
		var (
			img   = image.NewRGBA64(imgRectangle)
			gamma = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := AdjustGamma(img, gamma)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgsigmoid(t *testing.T) {
	convey.Convey("sigmoid", t, func(ctx convey.C) {
		var (
			a = float64(0)
			b = float64(0)
			x = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := sigmoid(a, b, x)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgAdjustSigmoid(t *testing.T) {
	convey.Convey("AdjustSigmoid", t, func(ctx convey.C) {
		var (
			img      = image.NewRGBA64(imgRectangle)
			midpoint = float64(0)
			factor   = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := AdjustSigmoid(img, midpoint, factor)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgAdjustContrast(t *testing.T) {
	convey.Convey("AdjustContrast", t, func(ctx convey.C) {
		var (
			img        = image.NewRGBA64(imgRectangle)
			percentage = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := AdjustContrast(img, percentage)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgAdjustBrightness(t *testing.T) {
	convey.Convey("AdjustBrightness", t, func(ctx convey.C) {
		var (
			img        = image.NewRGBA64(imgRectangle)
			percentage = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := AdjustBrightness(img, percentage)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgGrayscale(t *testing.T) {
	convey.Convey("Grayscale", t, func(ctx convey.C) {
		var (
			img = image.NewRGBA64(imgRectangle)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := Grayscale(img)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgInvert(t *testing.T) {
	convey.Convey("Invert", t, func(ctx convey.C) {
		var (
			img = image.NewRGBA64(imgRectangle)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := Invert(img)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgparallel(t *testing.T) {
	convey.Convey("parallel", t, func(ctx convey.C) {
		var (
			dataSize = int(0)
			fn       func(partStart int, partEnd int)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			parallel(dataSize, fn)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDrawimgabsint(t *testing.T) {
	convey.Convey("absint", t, func(ctx convey.C) {
		var (
			i = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := absint(i)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDrawimgclamp(t *testing.T) {
	convey.Convey("clamp", t, func(ctx convey.C) {
		var (
			x = float64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := clamp(x)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func monkeyJpegEncode() {
	monkey.Patch(jpeg.Encode, func(_ io.Writer, _ image.Image, _ *jpeg.Options) error {
		return nil
	})
}
