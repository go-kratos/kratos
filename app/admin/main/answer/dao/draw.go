package dao

import (
	"image"
	"image/draw"
	"io/ioutil"
	"math"
	"unicode/utf8"

	"go-common/library/log"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// QueConf question img conf.
func (d *Dao) QueConf(mobile bool) (c *TextImgConf) {
	if mobile {
		// Mobile
		c = &TextImgConf{
			Fontsize:    16, // mobile font size in points
			Length:      11, // mobile question length
			Ansfontsize: 12, // mobile ans font size in points
		}
	} else {
		// PC
		c = &TextImgConf{
			Fontsize:    12, //font size in points
			Length:      36, //question length
			Ansfontsize: 10, //ans font size in points
		}
	}
	c.Spacing = 2    // line spacing (e.g. 2 means double spaced)
	c.Ansspacing = 2 // line ansspacing (e.g. 2 means double spaced)
	return
}

// DrawQue draw question title.
func (d *Dao) DrawQue(c *freetype.Context, s string, conf *TextImgConf, pt *fixed.Point26_6) {
	c.SetFontSize(float64(conf.Fontsize))
	srune := []rune(s)
	var end = conf.Length
	for len(srune) > 0 {
		if conf.Length > len(srune) {
			end = len(srune)
		}
		d.text(c, string(srune[:end]), pt, conf.Fontsize, conf.Spacing)
		srune = srune[end:]
	}
}

// DrawAns draw ans
func (d *Dao) DrawAns(c *freetype.Context, conf *TextImgConf, anss [4]string, pt *fixed.Point26_6) {
	c.SetFontSize(float64(conf.Ansfontsize))
	arr := [4]string{"A.", "B.", "C.", "D."}
	for i, a := range anss {
		d.text(c, arr[i]+a, pt, conf.Ansfontsize, conf.Ansspacing)
	}
}

//Height get img height
func (d *Dao) Height(c *TextImgConf, que string, anslen int) (h int) {
	len := utf8.RuneCountInString(que)
	line := math.Ceil(float64(len) / float64(c.Length))
	h = int(math.Ceil(c.Spacing*line*float64(c.Fontsize))) + int(math.Ceil(c.Ansspacing*float64(anslen)*float64(c.Ansfontsize)))
	return
}

// text Draw text.
func (d *Dao) text(c *freetype.Context, s string, pt *fixed.Point26_6, size int, spacing float64) (err error) {
	_, err = c.DrawString(s, *pt)
	if err != nil {
		return
	}
	pt.Y += fixed.Int26_6(int(float64(size)*spacing) << 6)
	return
}

var (
	dpi     = float64(72) // screen resolution in Dots Per Inch
	hinting = "none"      // none | full
)

// Board init draw board.
func (d *Dao) Board(h int) (r *image.Gray) {
	bg := image.White
	r = image.NewGray(image.Rect(0, 0, 600, h))
	draw.Draw(r, r.Bounds(), bg, image.ZP, draw.Src)
	return
}

// Context freetype init context.
func (d *Dao) Context(r *image.Gray, fileStr string) (c *freetype.Context) {
	fg := image.Black
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fileStr)
	if err != nil {
		log.Error("ioutil.ReadFile(),err:%+v", err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Error("freetype.ParseFont(),err:%+v", err)
		return
	}
	c = freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetClip(r.Bounds())
	c.SetDst(r)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}
	return
}
