package feed

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-feed/model"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
	"go-common/library/net/metadata"
)

type infoc struct {
	typ         string
	mid         string
	client      string
	build       string
	buvid       string
	disid       string
	ip          string
	style       string
	api         string
	now         string
	isRcmd      string
	pull        string
	userFeature json.RawMessage
	code        string
	items       []*ai.Item
	zoneID      string
	adResponse  string
	deviceID    string
	network     string
	newUser     string
	flush       string
	autoPlay    string
	deviceType  string
}

func (s *Service) IndexInfoc(c context.Context, mid int64, plat int8, build int, buvid, disid, api string, userFeature json.RawMessage, style, code int, items []*ai.Item, isRcmd, pull, newUser bool, now time.Time, adResponse, deviceID, network string, flush int, autoPlay string, deviceType int) {
	if items == nil {
		return
	}
	var (
		isRc      = "0"
		isPull    = "0"
		isNewUser = "0"
		zoneID    int64
	)
	if isRcmd {
		isRc = "1"
	}
	if pull {
		isPull = "1"
	}
	if newUser {
		isNewUser = "1"
	}
	ip := metadata.String(c, metadata.RemoteIP)
	info, err := s.loc.Info(c, ip)
	if err != nil {
		log.Warn(" s.loc.Info(%v) error(%v)", ip, err)
		err = nil
	}
	if info != nil {
		zoneID = info.ZoneID
	}
	s.infoc(infoc{"综合推荐", strconv.FormatInt(mid, 10), strconv.Itoa(int(plat)), strconv.Itoa(build), buvid, disid, ip, strconv.Itoa(style), api, strconv.FormatInt(now.Unix(), 10), isRc, isPull, userFeature, strconv.Itoa(code), items, strconv.FormatInt(zoneID, 10), adResponse, deviceID, network, isNewUser, strconv.Itoa(flush), autoPlay, strconv.Itoa(deviceType)})
}

func (s *Service) infoc(i interface{}) {
	select {
	case s.logCh <- i:
	default:
		log.Warn("infocproc chan full")
	}
}

// writeInfoc
func (s *Service) infocproc() {
	const (
		// infoc format {"section":{"id":"%s推荐","pos":1,"style":%d,"items":[{"id":%s,"pos":%d,"type":1,"url":""}]}}
		noItem1 = `{"section":{"id":"`
		noItem2 = `{","pos":1,"style":`
		noItem3 = `,"items":[]}}`
	)
	// is_ad_loc, resource_id，source_id, creative_id,
	var (
		msg1    = []byte(`{"section":{"id":"`)
		msg2    = []byte(`","pos":1,"style":`)
		msg3    = []byte(`,"items":[`)
		msg4    = []byte(`{"id":`)
		msg5    = []byte(`,"pos":`)
		msg6    = []byte(`,"type":`)
		msg7    = []byte(`,"source":"`)
		msg8    = []byte(`","tid":`)
		msg9    = []byte(`,"av_feature":`)
		msg10   = []byte(`,"url":"`)
		msg11   = []byte(`","rcmd_reason":"`)
		msg12   = []byte(`","is_ad_loc":`)
		msg13   = []byte(`,"resource_id":`)
		msg14   = []byte(`,"source_id":`)
		msg15   = []byte(`,"creative_id":`)
		msg16   = []byte(`},`)
		msg17   = []byte(`"},`)
		msg18   = []byte(`]}}`)
		showInf = binfoc.New(s.c.ShowInfoc2)
		buf     bytes.Buffer
		list    string
		trackID string
	)
	for {
		i, ok := <-s.logCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		switch l := i.(type) {
		case infoc:
			if len(l.items) > 0 {
				buf.Write(msg1)
				buf.WriteString(l.typ)
				buf.Write(msg2)
				buf.WriteString(l.style)
				buf.Write(msg3)
				for i, v := range l.items {
					if v == nil {
						continue
					}
					if v.TrackID != "" {
						trackID = v.TrackID
					}
					buf.Write(msg4)
					buf.WriteString(strconv.FormatInt(v.ID, 10))
					buf.Write(msg5)
					buf.WriteString(strconv.Itoa(i + 1))
					buf.Write(msg6)
					buf.WriteString(gotoMapID(v.Goto))
					buf.Write(msg7)
					buf.WriteString(v.Source)
					buf.Write(msg8)
					buf.WriteString(strconv.FormatInt(v.Tid, 10))
					buf.Write(msg9)
					if v.AvFeature != nil {
						buf.Write(v.AvFeature)
					} else {
						buf.Write([]byte(`""`))
					}
					buf.Write(msg10)
					buf.WriteString("")
					buf.Write(msg11)
					if v.RcmdReason != nil {
						buf.WriteString(v.RcmdReason.Content)
					}
					if v.Ad != nil {
						buf.Write(msg12)
						buf.WriteString(strconv.FormatBool(v.Ad.IsAdLoc))
						buf.Write(msg13)
						buf.WriteString(strconv.FormatInt(v.Ad.Resource, 10))
						buf.Write(msg14)
						buf.WriteString(strconv.Itoa(v.Ad.Source))
						buf.Write(msg15)
						buf.WriteString(strconv.FormatInt(v.Ad.CreativeID, 10))
						buf.Write(msg16)
					} else {
						buf.Write(msg17)
					}
				}
				buf.Truncate(buf.Len() - 1)
				buf.Write(msg18)
				list = buf.String()
				buf.Reset()
			} else {
				list = noItem1 + l.typ + noItem2 + l.style + noItem3
			}
			showInf.Info(l.ip, l.now, l.api, l.buvid, l.mid, l.client, l.pull, list, l.disid, l.isRcmd, l.build, l.code, string(l.userFeature), l.zoneID, l.adResponse, l.deviceID, l.network, l.newUser, l.flush, l.autoPlay, trackID, l.deviceType)
			log.Info("infocproc %s param(mid:%s,buvid:%s,plat:%s,build:%s,isRcmd:%s,code:%s,zone_id:%s,user_feature:%s,ad_response:%s,device_id:%s,network:%s,new_user:%s,flush:%s,autoplay_card:%s,trackid:%s,device_type:%s) response(%s)", l.api, l.mid, l.buvid, l.client, l.build, l.isRcmd, l.code, l.zoneID, l.userFeature, l.adResponse, l.deviceID, l.network, l.newUser, l.flush, l.autoPlay, trackID, l.deviceType, list)
		}
	}
}

func gotoMapID(gt string) (id string) {
	if gt == model.GotoAv {
		id = "1"
	} else if gt == model.GotoBangumi {
		id = "2"
	} else if gt == model.GotoLive {
		id = "3"
	} else if gt == model.GotoRank {
		id = "6"
	} else if gt == model.GotoAdAv {
		id = "8"
	} else if gt == model.GotoAdWeb {
		id = "9"
	} else if gt == model.GotoBangumiRcmd {
		id = "10"
	} else if gt == model.GotoLogin {
		id = "11"
	} else if gt == model.GotoUpBangumi {
		id = "12"
	} else if gt == model.GotoBanner {
		id = "13"
	} else if gt == model.GotoAdWebS {
		id = "14"
	} else if gt == model.GotoUpArticle {
		id = "15"
	} else if gt == model.GotoConverge {
		id = "17"
	} else if gt == model.GotoSpecial {
		id = "18"
	} else if gt == model.GotoArticleS {
		id = "20"
	} else if gt == model.GotoGameDownloadS {
		id = "21"
	} else if gt == model.GotoShoppingS {
		id = "22"
	} else if gt == model.GotoAudio {
		id = "23"
	} else if gt == model.GotoPlayer {
		id = "24"
	} else if gt == model.GotoSpecialS {
		id = "25"
	} else if gt == model.GotoAdLarge {
		id = "26"
	} else if gt == model.GotoPlayerLive {
		id = "27"
	} else if gt == model.GotoSong {
		id = "28"
	} else if gt == model.GotoLiveUpRcmd {
		id = "29"
	} else if gt == model.GotoUpRcmdAv {
		id = "30"
	} else if gt == model.GotoSubscribe {
		id = "31"
	} else if gt == model.GotoChannelRcmd {
		id = "32"
	} else if gt == model.GotoMoe {
		id = "33"
	} else if gt == model.GotoPGC {
		id = "34"
	} else if gt == model.GotoSearchSubscribe {
		id = "35"
	} else if gt == model.GotoPicture {
		id = "36"
	} else if gt == model.GotoInterest {
		id = "37"
	} else {
		id = "-1"
	}
	return
}
