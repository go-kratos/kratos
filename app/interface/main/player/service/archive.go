package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/player/dao"
	"go-common/app/interface/main/player/model"
	accmdl "go-common/app/service/main/account/api"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_maxLevel  = 6
	_hasUGCPay = 1
)

// View get view info
func (s *Service) View(c context.Context, aid int64) (view *model.View, err error) {
	var viewReply *arcmdl.ViewReply
	if viewReply, err = s.arcClient.View(c, &arcmdl.ViewRequest{Aid: aid}); err != nil {
		dao.PromError("View接口错误", "s.arcClientView3(%d) error(%v)", aid, err)
		return
	}
	view = &model.View{Arc: viewReply.Arc, Pages: viewReply.Pages}
	return
}

// Matsuri get matsuri info
func (s *Service) Matsuri(c context.Context, now time.Time) (view *model.View) {
	if now.Unix() < s.matTime.Unix() {
		return s.pastView
	}
	if s.matOn || len(s.matView.Pages) < 1 {
		return s.matView
	}
	view = new(model.View)
	*view = *s.matView
	view.Pages = view.Pages[0 : len(view.Pages)-1]
	return
}

// PageList many p video pages
func (s *Service) PageList(c context.Context, aid int64) (rs []*arcmdl.Page, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if rs, err = s.arc.Page3(c, &archive.ArgAid2{Aid: aid, RealIP: ip}); err != nil {
		dao.PromError("Page3 接口错误", "s.arc.Page3(%d) error(%v)", aid, err)
	}
	return
}

// VideoShot get archive video shot data
func (s *Service) VideoShot(c context.Context, aid, cid int64, index bool) (res *model.Videoshot, err error) {
	var (
		viewReply *arcmdl.ViewReply
		ip        = metadata.String(c, metadata.RemoteIP)
	)
	if viewReply, err = s.arcClient.View(c, &arcmdl.ViewRequest{Aid: aid}); err != nil {
		log.Error("VideoShot s.arcClient.View(%d) error(%v)", aid, err)
		return
	}
	if !viewReply.Arc.IsNormal() || viewReply.Arc.Rights.UGCPay == _hasUGCPay {
		log.Warn("VideoShot warn arc(%d) state(%d) or ugcpay(%d)", aid, viewReply.Arc.State, viewReply.Arc.Rights.UGCPay)
		err = ecode.NothingFound
		return
	}
	if cid == 0 {
		if len(viewReply.Pages) == 0 {
			err = ecode.NothingFound
			return
		}
		cid = viewReply.Pages[0].Cid
	}
	res = &model.Videoshot{}
	if res.Videoshot, err = s.arc.Videoshot2(c, &archive.ArgCid2{Aid: aid, Cid: cid, RealIP: ip}); err != nil {
		log.Error("s.arc.Videoshot2(%d,%d) err(%v)", aid, cid, err)
		return
	}
	if index && res.PvData != "" {
		if pv, e := s.dao.PvData(c, res.PvData); e != nil {
			log.Error("s.dao.PvData(aid:%d,cid:%d) err(%+v)", aid, cid, e)
		} else if len(pv) > 0 {
			var (
				v   uint16
				pvs []uint16
				buf = bytes.NewReader(pv)
			)
			for {
				if e := binary.Read(buf, binary.BigEndian, &v); e != nil {
					if e != io.EOF {
						log.Warn("binary.Read pvdata(%s) err(%v)", res.PvData, e)
					}
					break
				}
				pvs = append(pvs, v)
			}
			res.Index = pvs
		}
	}
	fmtVideshot(res)
	return
}

func fmtVideshot(res *model.Videoshot) {
	if res.PvData != "" {
		res.PvData = strings.Replace(res.PvData, "http://", "//", 1)
	}
	for i, v := range res.Image {
		res.Image[i] = strings.Replace(v, "http://", "//", 1)
	}
}

// PlayURLToken get playurl token
func (s *Service) PlayURLToken(c context.Context, mid, aid, cid int64) (res *model.PlayURLToken, err error) {
	var (
		arcReply    *arcmdl.ArcReply
		ui          *accmdl.CardReply
		owner, svip int
		vip         int32
	)
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil {
		dao.PromError("Arc接口错误", "s.arcClient.Arc(%d) error(%v)", aid, err)
		err = ecode.NothingFound
		return
	}
	if !arcReply.Arc.IsNormal() {
		err = ecode.NothingFound
		return
	}
	if mid == arcReply.Arc.Author.Mid {
		owner = 1
	}
	if ui, err = s.accClient.Card3(c, &accmdl.MidReq{Mid: mid}); err != nil {
		dao.PromError("Card3接口错误", "s.accClient.Card3(%d) error(%v)", mid, err)
		err = ecode.AccessDenied
		return
	}
	if vip = ui.Card.Level; vip > _maxLevel {
		vip = _maxLevel
	}
	if ui.Card.Vip.Type != 0 && ui.Card.Vip.Status == 1 {
		svip = 1
	}
	res = &model.PlayURLToken{
		From:  "pc",
		Ts:    time.Now().Unix(),
		Aid:   aid,
		Cid:   cid,
		Mid:   mid,
		Owner: owner,
		VIP:   int(vip),
		SVIP:  svip,
	}
	params := url.Values{}
	params.Set("from", res.From)
	params.Set("ts", strconv.FormatInt(res.Ts, 10))
	params.Set("aid", strconv.FormatInt(res.Aid, 10))
	params.Set("cid", strconv.FormatInt(res.Cid, 10))
	params.Set("mid", strconv.FormatInt(res.Mid, 10))
	params.Set("vip", strconv.Itoa(res.VIP))
	params.Set("svip", strconv.Itoa(res.SVIP))
	params.Set("owner", strconv.Itoa(res.Owner))
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(strings.ToLower(tmp) + s.c.PlayURLToken.Secret))
	res.Fcs = hex.EncodeToString(mh[:])
	res.Token = base64.StdEncoding.EncodeToString([]byte(tmp + "&fcs=" + res.Fcs))
	return
}
