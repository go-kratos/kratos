package drawimg

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"strconv"
	"strings"

	"go-common/library/log"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/math/fixed"
)

// DrawImg create img info.
type DrawImg struct {
	size         int
	CanvasWidth  int
	CanvasHeight int
	File         string
	txtWidth     int
	srcImg       image.Image
	Canvas       *image.NRGBA
	c            *Context
	f            *truetype.Font
}

// NewDrawImg create new img.
func NewDrawImg(fontfile string, size int) (w *DrawImg) {
	w = &DrawImg{
		size: size,
	}
	f, err := w.readFont(fontfile, size)
	if err != nil {
		return
	}
	if f == nil {
		return
	}
	w.f = f
	return
}

//ReadSrcImg  read an image
func (w *DrawImg) ReadSrcImg(path string) (img image.Image, err error) {
	if img, err = Open(path); err != nil {
		log.Error("readSrcImg error(%v)", err)
	}
	return
}

func (w *DrawImg) readFont(path string, size int) (f *truetype.Font, err error) {
	fbs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("readFont error(%v)", err)
		return
	}
	f, err = ParseFont(fbs)
	if err != nil {
		log.Error("error(%v)", err)
	}
	cxt := NewContext()
	if cxt == nil {
		return
	}
	w.c = cxt
	w.c.SetFont(f)
	w.c.SetFontSize(float64(size))
	return
}

func (w *DrawImg) imgWidth() int {
	return w.srcImg.Bounds().Max.X
}

func (w *DrawImg) imgHeight() int {
	return w.srcImg.Bounds().Max.Y
}

func (w *DrawImg) newImgWidth() int {
	return w.txtWidth + w.imgWidth()
}

func (w *DrawImg) newCanvas(width, height int) *image.NRGBA {
	return NewNRGBA(width, height, color.RGBA{255, 0, 0, 0})
}

func (w *DrawImg) fillColor(r, g, b, a int32) image.Image {
	return &image.Uniform{color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}}
}

func (w *DrawImg) textWidth(text string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("draw set textWidth text(%s)|error(%v)", text, err)
		}
	}()
	box, err := w.c.FontBox(text, w.pt(3, 5))
	if err != nil {
		log.Error("set textWidth text(%s)|error(%v)", text, err)
		return
	}
	widthStr := strings.Split(Int26_6ToString(box.X), ":")[0]
	wid, err := strconv.ParseInt(widthStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt widthStr(%+v)|error(%v)", widthStr, err)
		return
	}
	w.txtWidth = int(wid)
	return
}

func (w *DrawImg) pt(x, y int) fixed.Point26_6 {
	return Pt(x, y+int(w.c.PointToFixed(float64(w.size))>>6))
}

func (w *DrawImg) setFont(text string, dstRgba *image.NRGBA, fsrc image.Image, pt fixed.Point26_6) (err error) {
	w.c.SetClip(w.Canvas.Bounds())
	w.c.SetDst(dstRgba)
	w.c.SetSrc(fsrc)
	_, err = w.c.DrawString(text, pt)
	if err != nil {
		log.Error("setFont error(%v)", err)
	}
	return
}

func (w *DrawImg) composite(dstCanvas *image.NRGBA, src image.Image, isLeft bool) {
	var p image.Point
	if isLeft {
		p = image.Point{-int(w.txtWidth), 0}
	} else {
		p = image.ZP
	}
	draw.Draw(dstCanvas, image.Rect(0, 0, w.newImgWidth(), w.imgHeight()), src, p, draw.Over)
}

// Draw  write text to the left or right of img.
func (w *DrawImg) Draw(text, savepath string, isLeft bool) (err error) {
	if text == "" {
		err = errors.New("draw: DrawText called with a empty text")
		return
	}
	w.textWidth(text)
	w.CanvasWidth = w.newImgWidth()
	w.CanvasHeight = w.imgHeight()
	w.Canvas = w.newCanvas(w.CanvasWidth, w.CanvasHeight)
	if w.c == nil || w.Canvas == nil {
		err = errors.New("draw: DrawText w.c or w.Canvas is nil")
		return
	}
	draw.Draw(w.Canvas, w.Canvas.Bounds(), w.Canvas, image.ZP, draw.Src)
	var pt fixed.Point26_6
	if isLeft {
		pt = w.pt(3, 5)
	} else {
		pt = w.pt(w.imgWidth(), 8)
	}
	black := w.fillColor(0, 0, 0, 125)
	w.setFont(text, w.Canvas, black, pt)
	blurRgba := Blur(w.Canvas, 6, 3.5)
	white := w.fillColor(255, 255, 255, 180)
	w.setFont(text, blurRgba, white, pt)
	w.composite(blurRgba, w.srcImg, isLeft)
	Save(blurRgba, savepath)
	return
}
