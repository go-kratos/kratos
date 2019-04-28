package service

import (
	"image"
	"math"
)

type rotate struct {
	dx   float64
	dy   float64
	sin  float64
	cos  float64
	neww float64
	newh float64
	src  *image.RGBA
}

func (r *rotate) rotate(angle float64, src *image.RGBA) *rotate {
	r.src = src
	srsize := src.Bounds().Size()
	width, height := srsize.X, srsize.Y

	// 源图四个角的坐标（以图像中心为坐标系原点）
	// 左下角,右下角,左上角,右上角
	srcwp, srchp := float64(width)*0.5, float64(height)*0.5
	srcx1, srcy1 := -srcwp, srchp
	srcx2, srcy2 := srcwp, srchp
	srcx3, srcy3 := -srcwp, -srchp
	srcx4, srcy4 := srcwp, -srchp

	r.sin, r.cos = math.Sincos(radian(angle))
	// 旋转后的四角坐标
	desx1, desy1 := r.cos*srcx1+r.sin*srcy1, -r.sin*srcx1+r.cos*srcy1
	desx2, desy2 := r.cos*srcx2+r.sin*srcy2, -r.sin*srcx2+r.cos*srcy2
	desx3, desy3 := r.cos*srcx3+r.sin*srcy3, -r.sin*srcx3+r.cos*srcy3
	desx4, desy4 := r.cos*srcx4+r.sin*srcy4, -r.sin*srcx4+r.cos*srcy4

	// 新的高度很宽度
	r.neww = math.Max(math.Abs(desx4-desx1), math.Abs(desx3-desx2)) + 0.5
	r.newh = math.Max(math.Abs(desy4-desy1), math.Abs(desy3-desy2)) + 0.5
	r.dx = -0.5*r.neww*r.cos - 0.5*r.newh*r.sin + srcwp
	r.dy = 0.5*r.neww*r.sin - 0.5*r.newh*r.cos + srchp
	return r
}

func radian(angle float64) float64 {
	return angle * math.Pi / 180.0
}

func (r *rotate) transformRGBA() image.Image {

	srcb := r.src.Bounds()
	b := image.Rect(0, 0, int(r.neww), int(r.newh))
	dst := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sx, sy := r.pt(x, y)
			if inBounds(srcb, sx, sy) {
				// 消除锯齿填色
				c := bili.RGBA(r.src, sx, sy)
				off := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
				dst.Pix[off+0] = c.R
				dst.Pix[off+1] = c.G
				dst.Pix[off+2] = c.B
				dst.Pix[off+3] = c.A
			}
		}
	}
	return dst
}

func (r *rotate) pt(x, y int) (float64, float64) {
	return float64(-y)*r.sin + float64(x)*r.cos + r.dy,
		float64(y)*r.cos + float64(x)*r.sin + r.dx
}

func inBounds(b image.Rectangle, x, y float64) bool {
	if x < float64(b.Min.X) || x >= float64(b.Max.X) {
		return false
	}
	if y < float64(b.Min.Y) || y >= float64(b.Max.Y) {
		return false
	}
	return true
}

func offRGBA(src *image.RGBA, x, y int) int {
	return (y-src.Rect.Min.Y)*src.Stride + (x-src.Rect.Min.X)*4
}
