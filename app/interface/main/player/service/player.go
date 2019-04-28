package service

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"html/template"
	"strconv"
	"strings"
	"time"

	dm2 "go-common/app/interface/main/dm2/model"
	history "go-common/app/interface/main/history/model"
	"go-common/app/interface/main/player/dao"
	"go-common/app/interface/main/player/model"
	tagmdl "go-common/app/interface/main/tag/model"
	accmdl "go-common/app/service/main/account/api"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/assist/model/assist"
	locmdl "go-common/app/service/main/location/model"
	resmdl "go-common/app/service/main/resource/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_content = `<a href="%s" target="_blank"><font color="#FFFFFF">%s</font></a>`
	_china   = "中国"
	_local   = "局域网"

	_accBanNor     = 0 // no block
	_accBanSta     = 1 // block MSpacesta
	_accBlockSta   = 1
	_dmMaskPlatWeb = 0
	_mockBlockTime = 100
)

var (
	_copyRightMap = map[int32]string{
		0: "Nnknown",
		1: "Original",
		2: "Copy",
	}

	// if typeid in this xml add bottom = 1
	_bottomMap = map[int32]struct{}{
		// 番剧
		33:  {},
		32:  {},
		153: {},
		// 电影
		82:  {},
		85:  {},
		145: {},
		146: {},
		147: {},
		83:  {},
		// note 电视剧存在三级分区
		15:  {},
		34:  {},
		86:  {},
		128: {},
		// 三级分区
		110: {},
		111: {},
		112: {},
		113: {},
		87:  {},
		88:  {},
		89:  {},
		90:  {},
		91:  {},
		92:  {},
		73:  {},
	}
	iconTagIDs = map[int64]struct{}{
		516:     {},
		374306:  {},
		16054:   {},
		18612:   {},
		2611047: {},
		1008087: {},
		50:      {},
		2513658: {},
		56:      {},
		2512304: {},
		6977:    {},
		8035683: {},
		1060128: {},
	}
)

// Carousel return carousel items.
func (s *Service) Carousel(c context.Context) (items []*model.Item, err error) {
	items = s.caItems
	return
}

// Player return player info.
func (s *Service) Player(c context.Context, mid, aid, cid int64, cdnIP, refer string, now time.Time) (res []byte, err error) {
	var (
		ip     = metadata.String(c, metadata.RemoteIP)
		vi     *arcmdl.ViewReply
		cuPage *arcmdl.Page
		pi     = &model.Player{
			IP:          ip,
			Login:       mid > 0,
			Time:        now.Unix(),
			ZoneIP:      cdnIP,
			Upermission: "1000,1001",
		}
		withU bool
	)
	if vi, err = s.view(c, aid); err != nil {
		dao.PromError("View接口错误", "s.arcClientView3(%d) error(%v)", aid, err)
		return
	} else if vi == nil || vi.Arc == nil {
		log.Error("vi(%v) is nill || vi.Archive is nil", vi)
		return
	} else if len(vi.Pages) == 0 {
		log.Error("len(vi.Pages) == 0 aid(%d)", aid)
		return
	}
	for _, page := range vi.Pages {
		if cid == page.Cid {
			cuPage = page
			break
		}
	}
	if cuPage == nil {
		log.Warn("cuPage is nil aid(%d) cid(%d) refer(%s)", aid, cid, refer)
	}
	s.fillArc(c, cid, pi, vi, cuPage, ip, now)
	withU = s.fillAcc(c, pi, vi, mid, cid, ip, now)
	// template
	var doc = bytes.NewBuffer(nil)
	if withU {
		s.tWithU.Execute(doc, pi)
	} else {
		s.tNoU.Execute(doc, pi)
	}
	if s.params != "" {
		doc.WriteString(s.params)
	}
	res = doc.Bytes()
	return
}

func (s *Service) fillAcc(c context.Context, pi *model.Player, vi *arcmdl.ViewReply, mid, cid int64, ip string, now time.Time) (withU bool) {
	if mid == 0 {
		return
	}
	var (
		proReply *accmdl.ProfileStatReply
		pro      map[int64]*history.History
		err      error
	)
	if proReply, err = s.accClient.ProfileWithStat3(c, &accmdl.MidReq{Mid: mid}); err != nil {
		dao.PromError("UserInfo接口错误", "s.acc.UserInfo(%v) error(%v)", mid, err)
		return
	}
	if proReply != nil {
		withU = true
		var nameBu = bytes.NewBuffer(nil)
		if err = xml.EscapeText(nameBu, []byte(proReply.Profile.Name)); err != nil {
			log.Error("xml.EscapeText(%s) error(%v)", proReply.Profile.Name, err)
		} else {
			pi.Name = nameBu.String()
		}
		pi.User = proReply.Profile.Mid
		pi.UserHash = midCrc(proReply.Profile.Mid)
		pi.Money = fmt.Sprintf("%.2f", proReply.Coins)
		pi.Face = strings.Replace(proReply.Profile.Face, "http://", "//", 1)
		var bs []byte
		if bs, err = json.Marshal(proReply.LevelInfo); err != nil {
			log.Error("json.Marshal(%v) error(%v)", proReply.LevelInfo, err)
		} else {
			pi.LevelInfo = template.HTML(bs)
		}
		vip := model.VIPInfo{Type: proReply.Profile.Vip.Type, DueDate: proReply.Profile.Vip.DueDate, VipStatus: proReply.Profile.Vip.Status}
		if bs, err = json.Marshal(vip); err != nil {
			log.Error("json.Marshal(%v) error(%v)", vip, err)
		} else {
			pi.Vip = template.HTML(bs)
		}
		off := &model.Official{Type: -1}
		if proReply.Profile.Official.Role != 0 {
			if proReply.Profile.Official.Role <= 2 {
				off.Type = 0
			} else {
				off.Type = 1
			}
			off.Desc = proReply.Profile.Official.Title
		}
		if bs, err = json.Marshal(off); err != nil {
			log.Error("json.Marshal(%v) error(%v)", off, err)
		} else {
			pi.OfficialVerify = template.HTML(bs)
		}
		group, errCtx := errgroup.WithContext(c)
		if vi.Arc != nil {
			pi.Upermission = userPermission(vi.Arc, proReply)
			// NOTE: if vInfo==nil, no admin
			if mid == vi.Arc.Author.Mid {
				pi.IsAdmin = true
			}
			group.Go(func() error {
				arg := &history.ArgPro{Mid: mid, RealIP: ip, Aids: []int64{vi.Arc.Aid}}
				if pro, err = s.his.Progress(errCtx, arg); err != nil {
					dao.PromError("Progress接口错误", "s.his.Progress(%d,%d) error(%v)", mid, vi.Arc.Aid, err)
				} else if progress, ok := pro[vi.Arc.Aid]; ok && progress != nil && progress.Cid > 0 && progress.Cid == cid {
					if progress.Pro >= 0 {
						pi.LastPlayTime = 1000 * progress.Pro
						pi.LastCid = progress.Cid
					} else if len(vi.Pages) != 0 {
						for _, page := range vi.Pages {
							if page.Cid == progress.Cid {
								pi.LastPlayTime = 1000 * page.Duration
								pi.LastCid = progress.Cid
								break
							}
						}
					}
				}
				return nil
			})
		}
		if s.c.Rule.NoAssistMid != vi.Arc.Author.Mid {
			group.Go(func() error {
				if assist, err := s.ass.Assist(errCtx, &assist.ArgAssist{Mid: vi.Arc.Author.Mid, AssistMid: proReply.Profile.Mid, Type: assist.TypeDm, RealIP: ip}); err != nil {
					dao.PromError("Assist接口错误", "s.ass.Assist(%d,%d) error(%v)", vi.Arc.Author.Mid, proReply.Profile.Mid, err)
				} else {
					pi.Role = strconv.FormatInt(assist.Assist, 10)
				}
				return nil
			})
		}
		if proReply.Profile.Silence == _accBanSta {
			group.Go(func() error {
				if blockTime, err := s.dao.BlockTime(errCtx, mid); err != nil {
					dao.PromError("BlockTime接口错误", "s.dao.BlockTime(%d) error(%v)", mid, err)
				} else if blockTime != nil {
					if blockTime.BlockStatus == _accBlockSta {
						pi.BlockTime = blockTime.BlockedEnd - now.Unix()
						if blockTime.BlockedForever || blockTime.BlockedEnd == 0 {
							pi.BlockTime = _mockBlockTime
						}
					}
				}
				return nil
			})
		}
		group.Wait()
	}
	return
}

func (s *Service) fillArc(c context.Context, cid int64, pi *model.Player, vi *arcmdl.ViewReply, page *arcmdl.Page, ip string, now time.Time) {
	// 稿件和其弹幕信息
	pi.Aid = vi.Arc.Aid
	pi.Typeid = vi.Arc.TypeID
	if page != nil {
		if page.From != "sina" {
			pi.Vtype = page.From
		} else {
			pi.Vtype = ""
		}
		pi.Maxlimit = dmLimit(page.Duration)
		pi.Chatid = page.Cid
		pi.Oriurl = oriURL(page.From, page.Vid)
		pi.Pid = int64(page.Page)
	} else {
		pi.Chatid = cid
		pi.Maxlimit = 1500
		pi.Pid = 1
	}
	pi.Arctype = _copyRightMap[vi.Arc.Copyright]
	pi.SuggestComment = false
	pi.Click = int(vi.Arc.Stat.View)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if click, err := s.arc.Click3(errCtx, &archive.ArgAid2{Aid: vi.Arc.Aid}); err != nil {
			dao.PromError("Click接口错误", "s.arc.Click2(%d) error(%v)", vi.Arc.Aid, err)
		} else if click != nil {
			pi.FwClick = click.H5 + click.Outter
		}
		return nil
	})
	pi.OnlineCount = 1
	group.Go(func() error {
		if onlineCount, err := s.dao.OnlineCount(errCtx, pi.Aid, cid); err == nil && onlineCount > 1 {
			pi.OnlineCount = onlineCount
		}
		return nil
	})
	group.Go(func() error {
		pi.MaskNew = s.dmMask(errCtx, cid)
		return nil
	})
	group.Go(func() error {
		pi.Subtitle = s.dmSubtitle(errCtx, pi.Aid, cid)
		return nil
	})
	group.Go(func() error {
		if ipInfo, e := s.loc.Info(errCtx, &locmdl.ArgIP{IP: ip}); e != nil {
			log.Error("fillArc s.loc.Info(%s) error(%v)", ip, e)
		} else if ipInfo != nil {
			pi.Zoneid = ipInfo.ZoneID
			pi.Country = ipInfo.Country
			pi.Acceptaccel = ipInfo.Country != _china && ipInfo.Country != _local
			pi.Cache = ipInfo.Country != _china && ipInfo.Country != _local
		}
		return nil
	})
	if vi.Arc.AttrVal(archive.AttrBitHasViewpoint) == archive.AttrYes {
		group.Go(func() error {
			pi.ViewPoints = s.viewPoints(errCtx, pi.Aid, cid)
			return nil
		})
	}
	group.Go(func() error {
		pi.PlayerIcon = s.tagPlayerIcon(errCtx, pi.Aid, ip)
		return nil
	})
	group.Wait()
	pi.Duration = formatDuration(vi.Arc.Duration)
	pi.AllowBp = vi.Arc.AttrVal(archive.AttrBitAllowBp) == 1
	if _, ok := _bottomMap[vi.Arc.TypeID]; ok {
		pi.Bottom = 1
	}
	pi.Acceptguest = false
	if s.BrBegin.Unix() <= now.Unix() && now.Unix() <= s.BrEnd.Unix() {
		pi.BrTCP = s.c.Broadcast.TCPAddr
		pi.BrWs = s.c.Broadcast.WsAddr
		pi.BrWss = s.c.Broadcast.WssAddr
	}
	for index, pa := range vi.Pages {
		if pa != nil && cid == pa.Cid && index+1 < len(vi.Pages) {
			pi.HasNext = 1
		}
	}
}

func isAdmin(uRank int32) (b bool) {
	// 32000 -> admin
	// 31300 -> 评论管理员
	if uRank == 31300 || uRank == 32000 {
		b = true
		return
	}
	return
}

func userPermission(a *arcmdl.Arc, u *accmdl.ProfileStatReply) (permission string) {
	if u.Profile.Silence == _accBanNor || isAdmin(u.Profile.Rank) {
		permission = strings.Join(append([]string{strconv.FormatInt(int64(u.Profile.Rank), 10), "1001"}), ",")
	} else {
		permission = "0"
	}
	// if a.AttrVal(archive.AttrBitNoMission) == 0 && a.Author.Mid == u.Mid {
	// 	permission = strings.Join([]string{permission, "20000"}, ",")
	// }
	return
}

func oriURL(dmType, dmIndex string) (url string) {
	switch dmType {
	case "sina":
		url = "http://p.you.video.sina.com.cn/swf/bokePlayer20131203_V4_1_42_33.swf?vid=" + dmIndex
	case "youku":
		url = "http://v.youku.com/v_show/id_" + dmIndex + ".html"
	case "qq":
		if len(dmIndex) >= 3 {
			url = "http://v.qq.com/page/" + dmIndex[0:1] + "/" + dmIndex[1:2] + "/" + dmIndex[2:3] + "/" + dmIndex + ".html"
		}
	default:
		url = ""
	}
	return
}

func formatDuration(duration int64) (du string) {
	if duration == 0 {
		du = "00:00"
	} else {
		var duFen, duMiao string
		duFen = strconv.Itoa(int(duration / 60))
		if int(duration%60) < 10 {
			duMiao = "0" + strconv.Itoa(int(duration%60))
		} else {
			duMiao = strconv.Itoa(int(duration % 60))
		}
		du = duFen + ":" + duMiao
	}
	return
}

func midCrc(mid int64) string {
	midStr := strconv.FormatInt(mid, 10)
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(midStr)))
}

func dmLimit(duration int64) (limit int) {
	switch {
	case duration > 3600:
		limit = 8000
	case duration > 2400:
		limit = 6000
	case duration > 900:
		limit = 3000
	case duration > 600:
		limit = 1500
	case duration > 150:
		limit = 1000
	case duration > 60:
		limit = 500
	case duration > 30:
		limit = 300
	case duration <= 30:
		limit = 100
	default:
		limit = 1500
	}
	return
}

func (s *Service) dmMask(c context.Context, cid int64) (mask template.HTML) {
	if dmMask, err := s.dm2.Mask(c, &dm2.ArgMask{Cid: cid, Plat: _dmMaskPlatWeb}); err != nil {
		dao.PromError("MaskList 错误", "s.dm2.MaskList cid(%d) error(%v)", cid, err)
	} else if dmMask != nil && dmMask.MaskURL != "" {
		dmMask.MaskURL = strings.Replace(dmMask.MaskURL, "http://", "//", 1)
		if bs, err := json.Marshal(dmMask); err != nil {
			log.Error("dmMask json.Marshal(%+v) error(%v)", dmMask, err)
		} else {
			mask = template.HTML(bs)
		}
	}
	return
}

func (s *Service) dmSubtitle(c context.Context, aid, cid int64) (subtitle template.HTML) {
	if dmSub, err := s.dm2.SubtitleGet(c, &dm2.ArgSubtitleGet{Aid: aid, Oid: cid, Type: dm2.SubTypeVideo}); err != nil {
		log.Error("s.dm2.SubtitleGet aid(%d) cid(%d) error(%v)", aid, cid, err)
	} else {
		if dmSub != nil {
			if len(dmSub.Subtitles) == 0 {
				dmSub.Subtitles = make([]*dm2.VideoSubtitle, 0)
			}
			for _, v := range dmSub.Subtitles {
				v.SubtitleURL = strings.Replace(v.SubtitleURL, "http://", "//", 1)
			}
			if bs, err := json.Marshal(dmSub); err != nil {
				log.Error("dmSubject json.Marshal(%v) error(%v)", dmSub, err)
			} else {
				subtitle = template.HTML(bs)
			}
		}
	}
	return
}

func (s *Service) tagPlayerIcon(c context.Context, aid int64, ip string) (icon template.HTML) {
	icon = s.icon
	now := time.Now()
	tags, err := s.tag.ArcTags(c, &tagmdl.ArgAid{Aid: aid, RealIP: ip})
	if err != nil {
		log.Error("tagPlayerIcon s.tag.ArcTags aid(%d) error(%v)", aid, err)
		return
	}
	// TODO delete tmp logic
	if now.Unix() >= s.c.Icon.Start.Unix() && now.Unix() <= s.c.Icon.End.Unix() {
		for _, vt := range tags {
			if _, ok := iconTagIDs[vt.ID]; ok {
				playerIcon := &resmdl.PlayerIcon{
					URL1:  s.c.Icon.URL1,
					Hash1: s.c.Icon.Hash1,
					URL2:  s.c.Icon.URL2,
					Hash2: s.c.Icon.Hash2,
				}
				bs, err := json.Marshal(playerIcon)
				if err != nil {
					log.Error("tagPlayerIcon json.Marshal(%v) error(%v)", playerIcon, err)
					continue
				}
				icon = template.HTML(bs)
				break
			}
		}
	}
	return
}

func (s *Service) viewPoints(c context.Context, aid, cid int64) (points template.HTML) {
	if data, err := s.dao.ViewPoints(c, aid, cid); err != nil {
		log.Error("s.dao.ViewPoints aid(%d) cid(%d) error(%v)", aid, cid, err)
	} else if len(data) > 0 {
		if bs, err := json.Marshal(data); err != nil {
			log.Error("viewPoints json.Marshal(%v) error(%v)", data, err)
		} else {
			points = template.HTML(bs)
		}
	}
	return
}

func (s *Service) view(c context.Context, aid int64) (data *arcmdl.ViewReply, err error) {
	if view, ok := s.bnj2019ViewMap[aid]; ok && view != nil {
		data = view
		return
	}
	return s.arcClient.View(c, &arcmdl.ViewRequest{Aid: aid})
}
