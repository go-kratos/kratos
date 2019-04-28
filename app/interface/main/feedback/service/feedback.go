package service

import (
	"context"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"go-common/app/interface/main/feedback/model"
	locmdl "go-common/app/service/main/location/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_androidFeedback  = "android-feedback"
	_androidFeedbackI = "android-fb-i"
	_iosFeedback      = "ios-feedback"
	_iosFeedbackI     = "ios-fb-i"
	_androidPlayer    = "android-player"
	_iosPlayer        = "ios-player"
	_androidCreative  = "android-creative"
	_iosCreative      = "ios-creative"
	_androidLivePink  = "android-live-pink"
	_iosLivePink      = "ios-live-pink"
	_tvYST            = "tv-yst"
	_creativeCenter   = "creative-center"
)

var (
	defaultTag = &model.Tag{
		ID:   314,
		Name: "其它问题",
		Type: 0,
	}
)

// AddReply add feedback raply and create session if session isn't exist
func (s *Service) AddReply(c context.Context, mid, tagID int64, buvid, system, version, mobiApp, content, imgURL, logURL, device, channel, entrance,
	netState, netOperator, agencyArea, platform, browser, qq, email string, now time.Time) (r *model.Reply, err error) {
	// check size
	if utf8.RuneCountInString(content) > s.c.Feedback.MaxContentSize {
		err = ecode.FeedbackBodyTooLarge
		return
	}
	var (
		ssn     *model.Session
		reply   *model.Reply
		cTime   xtime.Time
		replyID = buvid
		id      int64
		imgURLs []string
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	if entrance == "view" {
		if model.IsAndroid(model.Plat(mobiApp, device)) {
			platform = _androidPlayer
		} else if model.IsIOS(model.Plat(mobiApp, device)) {
			platform = _iosPlayer
		}
	} else if entrance == _tvYST || entrance == _creativeCenter {
		platform = entrance
	} else if platform == "android" {
		if entrance == _androidCreative {
			platform = _androidCreative
		} else if entrance == _androidLivePink {
			platform = _androidLivePink
		} else if model.Plat(mobiApp, device) == model.PlatAndroid {
			platform = _androidFeedback
		} else if model.Plat(mobiApp, device) == model.PlatAndroidI {
			platform = _androidFeedbackI
		}
	} else if platform == "ios" {
		if entrance == _iosCreative {
			platform = _iosCreative
		} else if entrance == _iosLivePink {
			platform = _iosLivePink
		} else if model.Plat(mobiApp, device) == model.PlatIPhone {
			platform = _iosFeedback
		} else if model.Plat(mobiApp, device) == model.PlatIPhoneI {
			platform = _iosFeedbackI
		}
	}
	if ssn, err = s.dao.Session(c, buvid, system, version, mid); err != nil {
		log.Error("s.dao.Session(%s,%s,%s,%d) error(%v)", buvid, system, version, mid, err)
		return
	}
	if mid != 0 {
		replyID = strconv.FormatInt(mid, 10)
	}
	cTime = xtime.Time(now.Unix())
	if ssn == nil {
		if content != "" || logURL != "" || imgURL != "" {
			// TODO delete ;
			imgURLs = strings.Split(imgURL, ";")
			var sid int64
			sid, err = s.session(c, mid, tagID, buvid, system, version, "", content, imgURLs[0], logURL, device, channel, ip, entrance, netState, netOperator, agencyArea, platform, browser, qq, email, now)
			if err == nil {
				for _, v := range imgURLs[1:] {
					if v != "" {
						reply = &model.Reply{
							SessionID: sid,
							ReplyID:   replyID,
							Content:   "",
							ImgURL:    v,
							LogURL:    "",
							CTime:     cTime,
							MTime:     cTime,
						}
						if id, err = s.dao.AddReply(c, reply); err != nil {
							log.Error("s.dao.AddReply error(%v)", err)
							continue
						}
					}
				}
			}
		}
	} else {
		if entrance == "view" || entrance == _androidLivePink || entrance == _iosLivePink || platform == _androidLivePink || platform == _iosLivePink || model.IsPlayerScreen(tagID) || entrance == _tvYST || platform == _tvYST || entrance == _creativeCenter || platform == _creativeCenter {
			_, err = s.session(c, mid, tagID, buvid, system, version, "", content, imgURL, logURL, device, channel, ip, entrance, netState, netOperator, agencyArea, platform, browser, qq, email, now)
			if err != nil {
				log.Error("s.session error (%v)", err)
				return
			}
		} else {
			stat := ssn.State
			if ssn.State == model.StateReplied {
				stat = model.StateRepeated
			}
			var ip32 uint32
			ipv := net.ParseIP(ip)
			if ip2 := ipv.To4(); ip2 != nil {
				ip32 = model.InetAtoN(ip)
			}
			ssn = &model.Session{
				ID:          ssn.ID,
				Device:      device,
				Channel:     channel,
				IP:          ip32,
				NetState:    netState,
				NetOperator: netOperator,
				AgencyArea:  agencyArea,
				Platform:    platform,
				Browser:     browser,
				QQ:          qq,
				Email:       email,
				State:       stat,
				LasterTime:  cTime,
				MTime:       cTime,
			}
			if _, err = s.dao.UpdateSession(c, ssn); err != nil {
				log.Error("s.dao.UpdateSession error(%v)", err)
				return
			}
			if content != "" || logURL != "" {
				reply = &model.Reply{
					SessionID: ssn.ID,
					ReplyID:   replyID,
					Content:   content,
					ImgURL:    "",
					LogURL:    logURL,
					CTime:     cTime,
					MTime:     cTime,
				}
				if id, err = s.dao.AddReply(c, reply); err != nil {
					log.Error("s.dao.AddReply error(%v)", err)
					return
				}
			}
			for _, v := range strings.Split(imgURL, ";") {
				if v != "" {
					cTime = xtime.Time(now.Unix())
					reply = &model.Reply{
						SessionID: ssn.ID,
						ReplyID:   replyID,
						Content:   "",
						ImgURL:    v,
						LogURL:    "",
						CTime:     cTime,
						MTime:     cTime,
					}
					if id, err = s.dao.AddReply(c, reply); err != nil {
						log.Error("s.dao.AddReply error(%v)", err)
						continue
					}
				}
			}
		}
	}
	r = &model.Reply{
		ID:      id,
		ReplyID: replyID,
		Type:    model.TypeCustomer,
		Content: content,
		ImgURL:  imgURL,
		LogURL:  logURL,
		CTime:   cTime,
	}
	return
}

// AddWebReply add web reply.
func (s *Service) AddWebReply(c context.Context, mid, sid, tagID int64, aid, content, imgURL, netState, netOperator, agencyArea, platform, version, buvid, browser, qq, email string, now time.Time) (r *model.Reply, err error) {
	var (
		cTime     xtime.Time
		reply     *model.Reply
		tx        *sql.Tx
		rid       int
		id        int64
		sessionID int64
		replyID   string
		ip        = metadata.String(c, metadata.RemoteIP)
	)
	if mid != 0 {
		replyID = strconv.FormatInt(mid, 10)
	}
	sessionID = sid
	cTime = xtime.Time(now.Unix())
	if buvid == "" {
		buvid = strconv.FormatInt(now.Unix(), 10)
	}
	rid, err = s.dao.JudgeSsnRecord(c, sessionID)
	if err != nil {
		log.Error("s.dao.JudgeSsnRecord error(%v)", err)
		return
	}
	if sessionID > 0 || rid > 0 {
		reply = &model.Reply{
			SessionID: sessionID,
			ReplyID:   replyID,
			Content:   content,
			ImgURL:    imgURL,
			CTime:     cTime,
			MTime:     cTime,
		}
		tx, err = s.dao.BeginTran(c)
		if err != nil {
			log.Error("s.dao.Begin error(%v)", err)
			return
		}
		if id, err = s.dao.TxAddReply(tx, reply); err != nil {
			log.Error("s.dao.TxAddReply error(%v)", err)
			tx.Rollback()
			return
		}
		if err = s.dao.TxUpdateSessionState(tx, model.StateNoReply, sessionID); err != nil {
			log.Error("s.dao.TxUpdateSessionState error(%v)", err)
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit(), error(%v)", err)
			return
		}
		r = &model.Reply{
			ID:      id,
			ReplyID: replyID,
			Type:    model.TypeCustomer,
			Content: content,
			ImgURL:  imgURL,
			CTime:   cTime,
		}
	} else {
		sessionID, err = s.session(c, mid, tagID, buvid, "", version, aid, content, imgURL, "", "", "", ip, "", netState, netOperator, agencyArea, platform, browser, qq, email, now)
		if err == nil {
			reply = &model.Reply{
				SessionID: sessionID,
				ReplyID:   replyID,
				Content:   content,
				ImgURL:    imgURL,
				CTime:     cTime,
				MTime:     cTime,
			}
			if _, err = s.dao.AddReply(c, reply); err != nil {
				log.Error("s.dao.AddReplyNoTx error(%v)", err)
				return
			}
		}
	}
	return
}

func (s *Service) session(c context.Context, mid, tagID int64, buvid, system, version string, aid, content, imgURL, logURL, device, channel, ip, entrance,
	netState, netOperator, agencyArea, platform, browser, qq, email string, now time.Time) (sid int64, err error) {
	if platform == "ugc" {
		var count int
		// 通过平台进行处理
		count, err = s.dao.SessionCount(c, mid)
		if err != nil {
			log.Error("s.dao.SessionCount error(%v)", err)
			return
		}
		if count >= 10 {
			err = ecode.FeedbackContentOver
			return
		}
	}
	cTime := xtime.Time(now.Unix())
	var ip32 uint32
	ipv := net.ParseIP(ip)
	if ip2 := ipv.To4(); ip2 != nil {
		ip32 = model.InetAtoN(ip)
	}
	ssn := &model.Session{
		Buvid:       buvid,
		System:      system,
		Version:     version,
		Mid:         mid,
		Aid:         aid,
		Content:     content,
		ImgURL:      imgURL,
		LogURL:      logURL,
		Device:      device,
		Channel:     channel,
		IP:          ip32,
		NetState:    netState,
		NetOperator: netOperator,
		AgencyArea:  agencyArea,
		Platform:    platform,
		Browser:     browser,
		QQ:          qq,
		Email:       email,
		State:       model.StateNoReply,
		LasterTime:  cTime,
		CTime:       cTime,
		MTime:       cTime,
	}
	if entrance == "view" || content == "播放器反馈日志" || entrance == _androidLivePink || entrance == _iosLivePink || platform == _androidLivePink || platform == _iosLivePink || model.IsPlayerScreen(tagID) || entrance == _tvYST || platform == _tvYST || entrance == _creativeCenter || platform == _creativeCenter {
		ssn.State = model.StateOther
	}
	var tx *sql.Tx
	tx, err = s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.Begin error(%v)", err)
		return
	}
	sid, err = s.dao.TxAddSession(tx, ssn)
	if err != nil {
		log.Error("s.dao.TxAddSession error(%v)", err)
		tx.Rollback()
		return
	}
	if err = s.dao.TxUpSsnMtime(tx, now, sid); err != nil {
		log.Error("s.dao.TxUpSsnMtime error(%v)", err)
		tx.Rollback()
		return
	}
	ssn.ID = sid
	if tagID > 0 {
		if _, err = s.dao.TxAddSessionTag(tx, ssn.ID, tagID, now); err != nil {
			log.Error("s.dao.TxAddSessionTag error(%v)", err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(), error(%v)", err)
	}
	return
}

// Replys show all feedback replays
func (s *Service) Replys(c context.Context, buvid, platform, mobiApp, device, system, version, entrance string, mid int64, pn, ps int) (rs []model.Reply, isEndReply bool, err error) {
	var (
		replyID = buvid
		ssns    []*model.Session
		ssn     *model.Session
		r       model.Reply
		tmp     []model.Reply
		offset  int
		limit   int
		start   int
		end     int
		rsl     int
	)
	offset = (pn - 1) * ps
	limit = ps
	if mid != 0 {
		if platform == "android" {
			if entrance == _androidCreative {
				platform = _androidCreative
			} else if entrance == _androidLivePink {
				platform = _androidLivePink
			} else if model.Plat(mobiApp, device) == model.PlatAndroid {
				platform = _androidFeedback
			} else if model.Plat(mobiApp, device) == model.PlatAndroidI {
				platform = _androidFeedbackI
			}
		} else if platform == "ios" {
			if entrance == _iosCreative {
				platform = _iosCreative
			} else if entrance == _iosLivePink {
				platform = _iosLivePink
			} else if model.Plat(mobiApp, device) == model.PlatIPhone {
				platform = _iosFeedback
			} else if model.Plat(mobiApp, device) == model.PlatIPhoneI {
				platform = _iosFeedbackI
			}
		}
		if ssns, err = s.dao.SessionByMid(c, mid, platform); err != nil {
			log.Error("s.dao.SessionByMid(%d) error(%v)", mid, err)
			return
		}
		if len(ssns) == 0 {
			err = ecode.FeedbackSsnNotExist
			return
		}
		replyID = strconv.FormatInt(mid, 10)
		if tmp, err = s.dao.ReplysByMid(c, mid, offset, limit); err != nil {
			log.Error("s.dao.ReplysByMid(%d) error(%v)", mid, err)
			return
		}
		if len(tmp) > 0 {
			rs = tmp
		}
		for _, ssn = range ssns {
			r = model.Reply{
				ID:      ssn.ID,
				ReplyID: replyID,
				Type:    model.TypeCustomer,
				Content: ssn.Content,
				ImgURL:  ssn.ImgURL,
				LogURL:  ssn.LogURL,
				CTime:   ssn.CTime,
			}
			rs = append(rs, r)
		}
		sort.Sort(model.Replys(rs))
		rsl = len(rs)
		end = rsl
		if limit > rsl {
			start = 0
			isEndReply = true
		} else {
			start = end - limit
		}
		rs = rs[start:end]
	} else {
		if ssn, err = s.dao.Session(c, buvid, system, version, mid); err != nil {
			log.Error("s.dao.Session(%s,%s,%s,%d) error(%v)", buvid, system, version, mid, err)
			return
		}
		if ssn == nil {
			err = ecode.FeedbackSsnNotExist
			return
		}
		if tmp, err = s.dao.Replys(c, ssn.ID, offset, limit); err != nil {
			log.Error("s.dao.Replays(%d) error(%v)", ssn.ID, err)
			return
		}
		if len(tmp) == limit {
			rs = tmp
		} else {
			r = model.Reply{
				ID:      ssn.ID,
				ReplyID: replyID,
				Type:    model.TypeCustomer,
				Content: ssn.Content,
				ImgURL:  ssn.ImgURL,
				LogURL:  ssn.LogURL,
				CTime:   ssn.CTime,
			}
			rs = append(tmp, r)
			isEndReply = true
		}
		sort.Sort(model.Replys(rs))
	}
	return
}

// Sessions sessions.
func (s *Service) Sessions(c context.Context, mid int64, state string, tagID, platform string, start, end time.Time, ps, pn int) (total int, wssns []*model.WebSession, err error) {
	var (
		sids, sidsTmp, tids, intersect, sidTmp, sidCut []int64
		limit, st, en, sls                             int
		ssnMap                                         map[int64]*model.Session
		ssns                                           []*model.Session
	)
	if mid > 0 {
		ssns, err = s.dao.SessionByMid(c, mid, platform)
		if err != nil {
			log.Error("s.dao.SessionByMid error(%v)", err)
			return
		}
		for _, v := range ssns {
			sids = append(sids, v.ID)
		}
		if tagID != "" {
			tids, err = xstr.SplitInts(tagID)
			if err != nil {
				log.Error("xstr.SplitInts error(%v)", err)
				return
			}
			if len(tids) > 0 {
				sidsTmp, err = s.dao.SessionIDByTagID(c, tids)
				if err != nil {
					log.Error("s.dao.SessionIDByTagID error(%v)", err)
					return
				}
			}
			for _, v := range sids {
				if contains(sidsTmp, v) {
					intersect = append(intersect, v)
				}
			}
		} else {
			intersect = sids
		}
	}
	if len(intersect) > 0 {
		if state == "" {
			ssns, err = s.dao.SSnBySsnIDAllSate(c, intersect, start, end)
			if err != nil {
				log.Error("s.dao.sSessionBySsnID error(%v)", err)
				return
			}
		} else {
			ssns, err = s.dao.SessionBySsnID(c, intersect, state, start, end)
			if err != nil {
				log.Error("s.dao.sSessionBySsnID error(%v)", err)
				return
			}
		}
		ssnMap = make(map[int64]*model.Session, len(ssns))
		for _, v := range ssns {
			sidTmp = append(sidTmp, v.ID)
			ssnMap[v.ID] = v
		}
		sls = len(sidTmp)
		total = sls
		limit = ps
		if limit > sls {
			st = 0
			en = sls
		} else {
			st = (pn - 1) * limit
			en = pn * limit
			if en > sls {
				en = st + (sls % limit)
			}
		}
	}
	sidCut = sidTmp[st:en]
	if len(sidCut) > 0 {
		var tagsMap map[int64][]*model.Tag
		tagsMap, err = s.dao.TagBySsnID(c, sidCut)
		if err != nil {
			log.Error("s.dao.TagBySsnID error(%v)", err)
			return
		}
		for _, v := range sidCut {
			wssn := &model.WebSession{}
			wssn.Session = ssnMap[v]
			tags := tagsMap[v]
			if len(tags) == 0 {
				wssn.Tag = defaultTag
			} else {
				wssn.Tag = tags[len(tags)-1]
			}
			wssns = append(wssns, wssn)
		}
	}
	return
}

// UpdateSessionState up session state.
func (s *Service) UpdateSessionState(c context.Context, state int, sid int64) (err error) {
	if err = s.dao.UpdateSessionState(c, state, sid); err != nil {
		log.Error("s.dao.UpdateSessionState error(%v)", err)
	}
	return
}

// Tags tags.
func (s *Service) Tags(c context.Context, mid int64, mold int, platform string) (tag *model.UGCTag, err error) {
	tags, err := s.dao.Tags(c, mold, platform)
	if err != nil {
		log.Error("s.dao.Tags error(%v)", err)
		return
	}
	cnt, err := s.dao.SessionCount(c, mid)
	if err != nil {
		log.Error("s.dao.SessionCount error(%v)", err)
		return
	}
	tag = &model.UGCTag{
		Tags:  tags,
		Limit: 10 - cnt,
	}
	return
}

// WebReplys web replys.
func (s *Service) WebReplys(c context.Context, sid, mid int64) (mp []*model.Reply, err error) {
	if mp, err = s.dao.WebReplys(c, sid, mid); err != nil {
		log.Error("s.dao.WebReplys(%d, %d) error(%v)", sid, mid, err)
		return
	}
	if len(mp) == 0 {
		err = ecode.NothingFound
	}
	return
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// PlayerCheck check for player.
func (s *Service) PlayerCheck(c context.Context, platform, ipChangeTimes int, mid, checkTime, aid, connectSpeed, ioSpeed int64, region, school, cdnip string) (err error) {
	var (
		ipinfo *locmdl.Info
		isp    int
		ipAddr = metadata.String(c, metadata.RemoteIP)
	)
	if ipinfo, err = s.locationRPC.Info(c, &locmdl.ArgIP{IP: ipAddr}); err != nil {
		log.Error("s.locationRPC.Info(%s) error(%v)", ipAddr, err)
	}
	if ipinfo != nil {
		if strings.Contains(ipinfo.ISP, "移动") {
			isp = 1
		} else if strings.Contains(ipinfo.ISP, "联通") {
			isp = 2
		} else if strings.Contains(ipinfo.ISP, "电信") {
			isp = 3
		} else {
			isp = 0
		}
	}
	if _, err = s.dao.InPlayCheck(c, platform, isp, ipChangeTimes, mid, checkTime, aid, connectSpeed, ioSpeed, region, school, ipAddr, cdnip); err != nil {
		log.Error("s.dao.InPlayCheck(%d, %d, %d, %d, %d, %d, %d, %d, %s, %s, %s, %s) error(%v)", platform, isp, ipChangeTimes, mid, checkTime, aid, connectSpeed, ioSpeed, region, school, ipAddr, cdnip, err)
	}
	return
}
