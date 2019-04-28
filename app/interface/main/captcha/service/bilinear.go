package service

import (
	"image"
	"image/color"
	"math"
)

var bili = Bilinear{}

// Bilinear Bilinear.
type Bilinear struct{}

// BilinearSrc BilinearSrc.
type BilinearSrc struct {
	// Top-left and bottom-right interpolation sources
	low, high image.Point
	// Fraction of each pixel to take. The 0 suffix indicates
	// top/left, and the 1 suffix indicates bottom/right.
	frac00, frac01, frac10, frac11 float64
}

// RGBA RGBA.
func (Bilinear) RGBA(src *image.RGBA, x, y float64) color.RGBA {
	p := findLinearSrc(src.Bounds(), x, y)

	// Array offsets for the surrounding pixels.
	off00 := offRGBA(src, p.low.X, p.low.Y)
	off01 := offRGBA(src, p.high.X, p.low.Y)
	off10 := offRGBA(src, p.low.X, p.high.Y)
	off11 := offRGBA(src, p.high.X, p.high.Y)

	var fr, fg, fb, fa float64

	fr += float64(src.Pix[off00+0]) * p.frac00
	fg += float64(src.Pix[off00+1]) * p.frac00
	fb += float64(src.Pix[off00+2]) * p.frac00
	fa += float64(src.Pix[off00+3]) * p.frac00

	fr += float64(src.Pix[off01+0]) * p.frac01
	fg += float64(src.Pix[off01+1]) * p.frac01
	fb += float64(src.Pix[off01+2]) * p.frac01
	fa += float64(src.Pix[off01+3]) * p.frac01

	fr += float64(src.Pix[off10+0]) * p.frac10
	fg += float64(src.Pix[off10+1]) * p.frac10
	fb += float64(src.Pix[off10+2]) * p.frac10
	fa += float64(src.Pix[off10+3]) * p.frac10

	fr += float64(src.Pix[off11+0]) * p.frac11
	fg += float64(src.Pix[off11+1]) * p.frac11
	fb += float64(src.Pix[off11+2]) * p.frac11
	fa += float64(src.Pix[off11+3]) * p.frac11

	var c color.RGBA
	c.R = uint8(fr + 0.5)
	c.G = uint8(fg + 0.5)
	c.B = uint8(fb + 0.5)
	c.A = uint8(fa + 0.5)
	return c
}

func findLinearSrc(b image.Rectangle, sx, sy float64) BilinearSrc {
	maxX := float64(b.Max.X)
	maxY := float64(b.Max.Y)
	minX := float64(b.Min.X)
	minY := float64(b.Min.Y)
	lowX := math.Floor(sx - 0.5)
	lowY := math.Floor(sy - 0.5)
	if lowX < minX {
		lowX = minX
	}
	if lowY < minY {
		lowY = minY
	}
	highX := math.Ceil(sx - 0.5)
	highY := math.Ceil(sy - 0.5)
	if highX >= maxX {
		highX = maxX - 1
	}
	if highY >= maxY {
		highY = maxY - 1
	}
	// In the variables below, the 0 suffix indicates top/left, and the
	// 1 suffix indicates bottom/right.
	// Center of each surrounding pixel.
	x00 := lowX + 0.5
	y00 := lowY + 0.5
	x01 := highX + 0.5
	y01 := lowY + 0.5
	x10 := lowX + 0.5
	y10 := highY + 0.5
	x11 := highX + 0.5
	y11 := highY + 0.5
	p := BilinearSrc{
		low:  image.Pt(int(lowX), int(lowY)),
		high: image.Pt(int(highX), int(highY)),
	}
	// Literally, edge cases. If we are close enough to the edge of
	// the image, curtail the interpolation sources.
	if lowX == highX && lowY == highY {
		p.frac00 = 1.0
	} else if sy-minY <= 0.5 && sx-minX <= 0.5 {
		p.frac00 = 1.0
	} else if maxY-sy <= 0.5 && maxX-sx <= 0.5 {
		p.frac11 = 1.0
	} else if sy-minY <= 0.5 || lowY == highY {
		p.frac00 = x01 - sx
		p.frac01 = sx - x00
	} else if sx-minX <= 0.5 || lowX == highX {
		p.frac00 = y10 - sy
		p.frac10 = sy - y00
	} else if maxY-sy <= 0.5 {
		p.frac10 = x11 - sx
		p.frac11 = sx - x10
	} else if maxX-sx <= 0.5 {
		p.frac01 = y11 - sy
		p.frac11 = sy - y01
	} else {
		p.frac00 = (x01 - sx) * (y10 - sy)
		p.frac01 = (sx - x00) * (y11 - sy)
		p.frac10 = (x11 - sx) * (sy - y00)
		p.frac11 = (sx - x10) * (sy - y01)
	}
	return p
}
