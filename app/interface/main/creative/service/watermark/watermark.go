package watermark

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/model/newcomer"
	"go-common/app/interface/main/creative/model/watermark"
	"go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	wmTipFormat = "原创水印将展示在视频%s角，以防他人盗用"
	wmTipMap    = map[int]string{
		1: "左上",
		2: "右上", // default
		3: "左下",
		4: "右下",
	}
)

// WaterMark get user watermark.
func (s *Service) WaterMark(c context.Context, mid int64) (w *watermark.Watermark, err error) {
	if w, err = s.wm.WaterMark(c, mid); err != nil {
		log.Error("s.wm.WaterMark(%d) error(%v)", mid, err)
		return
	}
	if w == nil {
		w = &watermark.Watermark{
			MID: mid,
			Pos: 2,
			Ty:  1,
		}
	}
	w.Tip = fmt.Sprintf(wmTipFormat, wmTipMap[int(w.Pos)])
	var pf *model.Profile
	ip := metadata.String(c, metadata.RemoteIP)
	if pf, err = s.acc.Profile(c, mid, ip); err != nil {
		log.Error("s.acc.Profile(%d) error(%v)", mid, err)
		return
	}
	if pf == nil {
		return
	}
	w.Uname = pf.Name
	return
}

// WaterMarkSet set watermark for user.
func (s *Service) WaterMarkSet(c context.Context, wp *watermark.WatermarkParam) (wm *watermark.Watermark, err error) {
	var (
		w   *watermark.Watermark
		wmm *watermark.Watermark
		pf  *model.Profile
	)
	mid, st, ty, pos, sync, ip := wp.MID, wp.State, wp.Ty, wp.Pos, wp.Sync, wp.IP
	if !watermark.IsState(st) {
		err = ecode.CreativeWaterMarkWrongState
		return
	}
	if !watermark.IsType(ty) {
		err = ecode.CreativeWaterMarkWrongType
		return
	}
	if !watermark.IsPos(pos) {
		err = ecode.CreativeWaterMarkWrongPosition // 位置参数错误
		return
	}
	if w, err = s.wm.WaterMark(c, mid); err != nil {
		log.Error("s.wm.Watermark(%d) error(%v)", mid, err)
		return
	}
	wm = &watermark.Watermark{
		MID:   mid,
		State: st,
		Ty:    ty,
		Pos:   pos,
		MTime: time.Now(),
	}
	if w != nil {
		wm.ID = w.ID
		wm.Uname = w.Uname
		wm.URL = w.URL
		wm.MD5 = w.MD5
		wm.Info = w.Info
		wm.CTime = w.CTime
	}
	if st == watermark.StatOpen || st == watermark.StatPreview { //开启、预览
		var (
			text   string
			isName bool
		)
		if ty == watermark.TypeName || ty == watermark.TypeNewName { //增加用户名在logo下方的水印
			if pf, err = s.acc.Profile(c, mid, ip); err != nil {
				log.Error("s.acc.Profile(%d) error(%v)", mid, err)
				return
			}
			if pf == nil {
				return
			}
			if w != nil && w.State == watermark.StatOpen && ty == w.Ty && pos == w.Pos && w.Uname == pf.Name && sync == 0 {
				log.Info("repeat uname watermark")
				return
			}
			text = pf.Name
			wm.Uname = text
			isName = true
		} else if ty == watermark.TypeUID {
			if w != nil && w.State == watermark.StatOpen && ty == w.Ty && pos == w.Pos {
				log.Info("repeat uid watermark")
				return
			}
			text = strconv.FormatInt(mid, 10)
			isName = false
		}
		if ty == watermark.TypeName || ty == watermark.TypeUID { //old get wm
			if wmm, err = s.draw(c, mid, text, isName); err != nil {
				log.Error("s.draw error(%v)", err)
				err = ecode.CreativeWaterMarkCreateFailed
				return
			}
			if wmm == nil {
				return
			}
			wm.Info, wm.URL, wm.MD5 = wmm.Info, wmm.URL, wmm.MD5
		} else if ty == watermark.TypeNewName { //new get wm
			var gm *watermark.Watermark
			gm, err = s.GenWm(c, mid, wm.Uname, ip)
			if err != nil || gm == nil {
				return
			}
			wm.Info, wm.URL, wm.MD5 = gm.Info, gm.URL, gm.MD5
		}
	}
	if st == watermark.StatPreview { //预览不更新db
		return
	}
	if w == nil {
		wm.CTime = time.Now()
		_, err = s.wm.AddWaterMark(c, wm)
	} else {
		_, err = s.wm.UpWaterMark(c, wm)
	}

	res, _ := s.WaterMark(c, mid)
	if res != nil && res.State == 1 && res.URL != "" {
		s.p.TaskPub(mid, newcomer.MsgForWaterMark, newcomer.MsgFinishedCount)
	}
	return
}

func (s *Service) userInfoConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.userInfoSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.userInfoSub.Messages closed")
			return
		}
		msg.Commit()
		s.userInfoMo++
		u := &watermark.Msg{}
		if err = json.Unmarshal(msg.Value, u); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if u == nil || u.Action != "update" {
			continue
		}
		s.update(c, u)
		log.Info("userInfoConsumer key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
	}
}

func (s *Service) update(c context.Context, u *watermark.Msg) (err error) {
	if u.Old.Uname == u.New.Uname {
		return
	}
	var w, wm *watermark.Watermark
	if w, err = s.wm.WaterMark(c, u.New.MID); err != nil {
		log.Error("s.wm.Watermark(%d) error(%v)", u.New.MID, err)
		return
	}
	if w == nil {
		return
	}
	log.Info("user mid(%d) origin data(%+v)", w.MID, w)
	uname := u.New.Uname
	if w.Ty == watermark.TypeName {
		wm, err = s.draw(c, w.MID, uname, true)
		if err != nil {
			log.Error("s.draw error(%v)", err)
			err = ecode.CreativeWaterMarkCreateFailed
			return
		}
		if wm == nil {
			return
		}
		if wm.MD5 == "" {
			log.Error("md5Sum fail")
			err = ecode.CreativeWaterMarkCreateFailed
			return
		}
		w.Info, w.URL, w.MD5 = wm.Info, wm.URL, wm.MD5
	} else if w.Ty == watermark.TypeNewName { //new get wm
		var gm *watermark.Watermark
		gm, err = s.GenWm(c, w.MID, uname, "")
		if err != nil || gm == nil {
			return
		}
		w.Info, w.URL, w.MD5 = gm.Info, gm.URL, gm.MD5
	}
	w.Uname = uname
	w.MTime = time.Now()
	_, err = s.wm.UpWaterMark(c, w)
	log.Info("user mid(%d) uname from (%s) to (%s) update data(%+v)", u.New.MID, u.Old.Uname, u.New.Uname, w)
	return
}

func (s *Service) draw(c context.Context, mid int64, text string, isUname bool) (w *watermark.Watermark, err error) {
	dw, err := s.drawimg.Make(c, mid, text, isUname)
	if err != nil {
		log.Error("s.drawimg.Make error(%v)", err)
		return
	}
	if dw == nil {
		return
	}
	file := dw.File
	defer os.Remove(file)
	url, err := s.bfs.UploadByFile(c, file)
	if err != nil {
		log.Error("s.bfs.UploadByFile error(%v)", err)
		return
	}
	info, err := ImageInfo(dw.CanvasWidth, dw.CanvasHeight)
	if err != nil {
		return
	}
	w = &watermark.Watermark{}
	w.URL, w.Info, w.MD5 = url, info, MD5Sum(file)
	return
}

// MD5Sum calculate file md5.
func MD5Sum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		log.Error("md5Sum os.Open error(%v)", err)
		return ""
	}
	defer f.Close()
	r := bufio.NewReader(f)
	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		log.Error("md5Sum io.Copy error(%v)", err)
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

//GenWm for new genwm api.
func (s *Service) GenWm(c context.Context, mid int64, uname, ip string) (wm *watermark.Watermark, err error) {
	var genwm *watermark.GenWatermark
	genwm, err = s.wm.GenWm(c, mid, uname, ip)
	if err != nil {
		log.Error("s.wm.GenWm error(%v)", err)
		return
	}
	if genwm == nil {
		return
	}
	wm = &watermark.Watermark{}
	info, err := ImageInfo(genwm.Width, genwm.Height)
	if err != nil {
		return
	}
	wm.URL, wm.MD5, wm.Info = genwm.Location, genwm.MD5, info
	return
}

//ImageInfo for image info json.
func ImageInfo(width, height int) (info string, err error) {
	imgInfo := &watermark.Image{
		Width:  width,
		Height: height,
	}
	var bs []byte
	if bs, err = json.Marshal(&imgInfo); err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}
	info = string(bs)
	return
}

// AsyncWaterMarkSet fn
func (s *Service) AsyncWaterMarkSet(wp *watermark.WatermarkParam) {
	if s.closed {
		log.Warn("AsyncWaterMarkSet chan is closed")
		return
	}
	select {
	case s.wmChan <- wp:
	default:
		log.Error("AsyncWaterMarkSet chan is full data(%+v)", wp)
	}
}

func (s *Service) asyncWmSetProc() {
	c := context.Background()
	for {
		v, ok := <-s.wmChan
		if ok {
			log.Info("watermark set by async with data(%+v)", v)
			if _, err := s.WaterMarkSet(c, v); err != nil {
				log.Error("s.WaterMarkSet watermark err (%+v)", err)
			}
		}
	}
}
