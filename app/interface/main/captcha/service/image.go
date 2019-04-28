package service

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// Image Image.
type Image struct {
	*image.RGBA
}

// newImage create a new image
func newImage(width, length int) *Image {
	img := &Image{image.NewRGBA(image.Rect(0, 0, width, length))}
	return img
}

func (img *Image) fillBkg(c image.Image) {
	draw.Draw(img, img.Bounds(), c, image.ZP, draw.Over)
}

// drawCircle draw circle.
func (img *Image) drawCircle(cx, cy, radius int, isFillColor bool, color color.Color) {
	point := img.Bounds().Size()
	// 如果圆在图片可见区域外，直接退出
	if cx+radius < 0 || cx-radius >= point.X || cy+radius < 0 || cy-radius >= point.Y {
		return
	}
	x, y, d := 0, radius, 3-2*radius
	for x <= y {
		if isFillColor {
			for yi := x; yi <= y; yi++ {
				img.drawCircle8(cx, cy, x, yi, color)
			}
		} else {
			img.drawCircle8(cx, cy, x, y, color)
		}
		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}

// drawLine .
// Bresenham算法(https://zh.wikipedia.org/zh-cn/布雷森漢姆直線演算法).
// startX,startY 起点 endX,endY终点.
func (img *Image) drawLine(startX, startY, endX, endY int, color color.Color) {
	dx, dy, flag := int(math.Abs(float64(startY-startX))), int(math.Abs(float64(endY-startY))), false
	if dy > dx {
		flag = true
		startX, startY = startY, startX
		endX, endY = endY, endX
		dx, dy = dy, dx
	}
	ix, iy := sign(endX-startX), sign(endY-startY)
	n2dy := dy * 2
	n2dydx := (dy - dx) * 2
	d := n2dy - dx
	for startX != endX {
		if d < 0 {
			d += n2dy
		} else {
			startY += iy
			d += n2dydx
		}
		if flag {
			img.Set(startY, startX, color)
		} else {
			img.Set(startX, startY, color)
		}
		startX += ix
	}
}

func (img *Image) drawCircle8(xc, yc, x, y int, color color.Color) {
	img.Set(xc+x, yc+y, color)
	img.Set(xc-x, yc+y, color)
	img.Set(xc+x, yc-y, color)
	img.Set(xc-x, yc-y, color)
	img.Set(xc+y, yc+x, color)
	img.Set(xc-y, yc+x, color)
	img.Set(xc+y, yc-x, color)
	img.Set(xc-y, yc-x, color)
}

// drawString image draw string.
func (img *Image) drawString(font *truetype.Font, color color.Color, str string, fontsize float64) {
	ctx := freetype.NewContext()
	// default 72dpi
	ctx.SetDst(img)
	ctx.SetClip(img.Bounds())
	ctx.SetSrc(image.NewUniform(color))
	ctx.SetFontSize(fontsize)
	ctx.SetFont(font)
	// 写入文字的位置
	pt := freetype.Pt(0, int(-fontsize/6)+ctx.PointToFixed(fontsize).Ceil())
	ctx.DrawString(str, pt)
}

// 水波纹, amplude=振幅, period=周期
// copy from https://github.com/dchest/captcha/blob/master/image.go
func (img *Image) distortTo(amplude float64, period float64) {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	oldm := img.RGBA
	dx := 1.4 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			rgba := oldm.RGBAAt(x+int(xo), y+int(yo))
			if rgba.A > 0 {
				oldm.SetRGBA(x, y, rgba)
			}
		}
	}
}

// Rotate 旋转
func (img *Image) rotate(angle float64) image.Image {
	return new(rotate).rotate(angle, img.RGBA).transformRGBA()
}
