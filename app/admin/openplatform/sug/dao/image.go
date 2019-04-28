package dao

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go-common/app/admin/openplatform/sug/model"
	"go-common/library/log"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"code.google.com/p/graphics-go/graphics"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	_Width        = 1125
	_Height       = 234
	_HeaderLength = 168
)

var (
	_RGBA       = []int{244, 245, 247, 255}
	_HeadArea   = []int{36, 33, 204, 201}
	_LabelArea  = []int{234, 156, 348, 201}
	_BriefText  = []int{234, 72}
	_NameText   = []int{234, 129}
	_BriefSize  = 44
	_NameSize   = 36
	_BriefRGBA  = []int{33, 33, 33, 255}
	_NameRGBA   = []int{153, 153, 153, 255}
	_BriefLimit = 1000
	_NameLimit  = 1000
)

// CreateItemPNG make a pic for sug
func (d *Dao) CreateItemPNG(item model.Item) (location string, err error) {
	var r *image.NRGBA64
	if r, err = d.makeBoard(item.Img); err != nil {
		log.Error("Create picture board error(%v)", err)
		return
	}
	if item.Brief == "" {
		item.Brief = item.Name
	}
	r = d.drawText(r, item.Brief, item.Name)
	buf := new(bytes.Buffer)
	png.Encode(buf, r)
	bufReader := bufio.NewReader(buf)
	if location, err = d.Upload(context.TODO(), "image/png", fmt.Sprintf("season_sug_%s/%d.png", d.c.Env, item.ItemsID), bufReader); err != nil {
		log.Error("Upload pic png error (%v)", err)
		return
	}
	reg, _ := regexp.CompilePOSIX(`//(.*)+`)
	location = reg.FindString(location)
	return
}

func (d *Dao) makeBoard(headerURL string) (board *image.NRGBA64, err error) {
	var header image.Image
	radius, _ := os.Open(d.c.SourcePath + "radius.png")
	label, _ := os.Open(d.c.SourcePath + "label.png")
	defer radius.Close()
	defer label.Close()
	resp, err := http.Get("http:" + headerURL)
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("ioutil.ReadAll img error(%v)", err)
		return
	}
	imgType := http.DetectContentType(bs)
	buf := bytes.NewBuffer(bs)
	resp.Body = ioutil.NopCloser(buf)
	if err != nil {
		log.Error("http download head img error(%v)", err)
		return
	}
	defer resp.Body.Close()
	switch imgType {
	case "image/png":
		header, err = png.Decode(resp.Body)
		if err != nil {
			log.Error("png picture decode error(%v)", err)
			return
		}
	case "image/jpeg", "image/jpg":
		header, err = jpeg.Decode(resp.Body)
		if err != nil {
			log.Error("jpg picture decode err(%v)", err)
			return
		}
	default:
		log.Error("invaild picture type (%s)", headerURL)
		return
	}
	board = image.NewNRGBA64(image.Rect(0, 0, _Width, _Height))
	draw.Draw(board, board.Bounds(), image.White, image.ZP, draw.Src)
	radiusPNG, _ := png.Decode(radius)
	cover := image.NewRGBA64(image.Rect(0, 0, _HeaderLength, _HeaderLength))
	border := image.NewNRGBA64(image.Rect(0, 0, _HeaderLength, _HeaderLength))
	headBoard := image.NewRGBA64(image.Rect(0, 0, _HeaderLength, _HeaderLength))
	draw.Draw(headBoard, headBoard.Bounds(), image.NewUniform(color.NRGBA{uint8(_RGBA[0]), uint8(_RGBA[1]), uint8(_RGBA[2]), uint8(_RGBA[3])}), image.ZP, draw.Over)
	graphics.Thumbnail(border, radiusPNG)
	graphics.Thumbnail(cover, header)
	labelPNG, _ := png.Decode(label)
	draw.Draw(headBoard, headBoard.Bounds(), cover, image.ZP, draw.Over)
	draw.Draw(headBoard, headBoard.Bounds(), border, image.ZP, draw.Over)
	draw.Draw(board, image.Rect(_HeadArea[0], _HeadArea[1], _HeadArea[2], _HeadArea[3]), headBoard, image.ZP, draw.Over)
	draw.Draw(board, image.Rect(_LabelArea[0], _LabelArea[1], _LabelArea[2], _LabelArea[3]), labelPNG, image.ZP, draw.Over)
	return
}

func (d *Dao) drawText(r *image.NRGBA64, brief string, name string) *image.NRGBA64 {
	ptBrief := fixed.P(_BriefText[0], _BriefText[1])
	ptName := fixed.P(_NameText[0], _NameText[1])
	freetypeC := d.Context(r)
	d.Text(freetypeC, brief, &ptBrief, _BriefSize, image.NewUniform(color.NRGBA{uint8(_BriefRGBA[0]), uint8(_BriefRGBA[1]), uint8(_BriefRGBA[2]), uint8(_BriefRGBA[3])}), fixed.I(_BriefLimit))
	d.Text(freetypeC, name, &ptName, _NameSize, image.NewUniform(color.NRGBA{uint8(_NameRGBA[0]), uint8(_NameRGBA[1]), uint8(_NameRGBA[2]), uint8(_RGBA[3])}), fixed.I(_NameLimit))
	return r
}

// Context get font for drawing
func (d *Dao) Context(r *image.NRGBA64) (c *freetype.Context) {
	c = freetype.NewContext()
	c.SetClip(r.Bounds())
	c.SetDst(r)
	c.SetHinting(font.HintingNone)
	return
}

// Text draw letters on pic
func (d *Dao) Text(c *freetype.Context, s string, pt *fixed.Point26_6, size int, color image.Image, length fixed.Int26_6) (err error) {
	c.SetFontSize(float64(size))
	c.SetSrc(color)
	ttf, _ := ioutil.ReadFile(d.c.SourcePath + "font.ttf")
	font, _ := freetype.ParseFont(ttf)
	c.SetFont(font)
	for _, r := range s {
		c.DrawString(string(r), *pt)
		pt.X += font.HMetric(fixed.Int26_6(size*64), font.Index(r)).AdvanceWidth
		if pt.X > length {
			break
		}
	}
	return
}
