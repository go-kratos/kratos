package drawimg

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/library/log"
	"os"
	"strconv"
	"time"
)

// Dao  define
type Dao struct {
	// conf
	c *conf.Config
	// watermark
	dw *DrawImg
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	if !isExist(c.WaterMark.FontFile) {
		log.Error("font file not exist")
		return
	}
	if !isExist(c.WaterMark.UnameMark) {
		log.Error("uname image file not exist")
		return
	}
	if !isExist(c.WaterMark.UIDMark) {
		log.Error("uid image file not exist")
		return
	}
	d = &Dao{
		c:  c,
		dw: NewDrawImg(c.WaterMark.FontFile, c.WaterMark.FontSize),
	}
	return
}

// Make create watermark.
func (d *Dao) Make(c context.Context, mid int64, text string, isUname bool) (dw *DrawImg, err error) {
	var src string
	if isUname {
		src = d.c.WaterMark.UnameMark
	} else {
		src = d.c.WaterMark.UIDMark
	}
	img, err := d.dw.ReadSrcImg(src)
	if err != nil {
		return
	}
	if img == nil {
		return
	}
	d.dw.srcImg = img
	midStr := strconv.FormatInt(mid, 10)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	file := d.c.WaterMark.SaveImg + midStr + "-" + timestamp + ".png"
	if err = d.dw.Draw(text, file, isUname); err != nil {
		log.Error("d.dw.Draw error(%v)", err)
		return
	}
	dw = &DrawImg{
		CanvasWidth:  d.dw.CanvasWidth,
		CanvasHeight: d.dw.CanvasHeight,
		File:         file,
	}
	return
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
