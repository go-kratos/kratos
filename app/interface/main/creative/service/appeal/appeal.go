package appeal

import (
	"context"
	xtime "go-common/library/time"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/creative/model/appeal"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

// List get appeal List.
func (s *Service) List(c context.Context, mid int64, pn, ps int, tp, ip string) (all, open, closed int, res []*appeal.Appeal, err error) {
	list, err := s.ap.AppealList(c, mid, appeal.Business, ip)
	if err != nil {
		log.Error("s.ap.Appeals error(%v)", err)
		return
	}
	if len(list) == 0 {
		return
	}
	var (
		start, end int
		oaps       []*appeal.Appeal // open appeal
		caps       []*appeal.Appeal // close appeal
		apTmp      []*appeal.Appeal // tmp appeal .
	)
	aps := make([]*appeal.Appeal, 0, len(list))
	for _, v := range list {
		ap := &appeal.Appeal{}
		ap.ID = v.ID
		ap.Oid = v.Oid
		ap.Cid = v.ID
		ap.Mid = v.Mid
		ap.State = v.BusinessState
		ap.Content = v.Desc
		ap.Description = v.Desc
		ap.CTime = v.CTime
		ap.MTime = v.MTime
		aps = append(aps, ap)
	}
	for _, v := range aps {
		if appeal.IsOpen(v.State) {
			oaps = append(oaps, v)
		} else {
			caps = append(caps, v)
		}
	}
	all = len(aps)
	open = len(oaps)
	closed = len(caps)
	if pn > 1 {
		start = (pn - 1) * ps
	} else {
		start = 0
	}
	end = pn * ps
	if tp == "open" {
		apTmp = oaps
	} else if tp == "closed" {
		apTmp = caps
	} else {
		apTmp = aps
	}
	total := len(apTmp)
	if total == 0 {
		return
	}
	if total <= start {
		res = make([]*appeal.Appeal, 0)
	} else if total <= end {
		res = apTmp[start:total]
	} else {
		res = apTmp[start:end]
	}
	if len(res) > 0 {
		var (
			aids []int64
			aMap map[int64]*archive.ArcVideo
		)
		for _, v := range res {
			aids = append(aids, v.Oid)
		}
		if len(aids) > 0 {
			if aMap, err = s.arc.Views(c, mid, aids, ip); err != nil {
				log.Error("s.arc.Archives aids(%v), ip(%s) err(%v)", aids, ip, err)
				return
			}
		}
		for _, v := range res {
			if _, ok := aMap[v.Oid]; !ok {
				continue
			}
			arc := &api.Arc{
				Aid:   aMap[v.Oid].Archive.Aid,
				Pic:   coverURL(aMap[v.Oid].Archive.Cover),
				Title: aMap[v.Oid].Archive.Title,
				State: int32(aMap[v.Oid].Archive.State),
			}
			v.ID = v.Cid
			v.Title = arc.Title
			v.Pics = arc.Pic
			v.Archive = arc
		}
	}
	return
}

// Detail get a appeal and events.
func (s *Service) Detail(c context.Context, mid, cid int64, ip string) (ap *appeal.Appeal, err error) {
	var apmeta *appeal.AppealMeta
	if apmeta, err = s.ap.AppealDetail(c, mid, cid, appeal.Business, ip); err != nil {
		log.Error("s.ap.AppealDetail error(%v)", err)
		return
	}
	if apmeta == nil {
		err = ecode.AppealNotExist
		return
	}
	var (
		av                *archive.ArcVideo
		arc               *api.Arc
		starStr, etimeStr string
		star, etime       int64
	)
	ap = &appeal.Appeal{}
	ap.ID = apmeta.ID
	ap.Oid = apmeta.Oid
	ap.Cid = apmeta.ID
	ap.Mid = apmeta.Mid
	ap.State = apmeta.BusinessState
	ap.Content = apmeta.Desc
	ap.Description = apmeta.Desc
	ap.CTime = apmeta.CTime
	ap.MTime = apmeta.MTime
	ap.Attachments = apmeta.Attachments
	if starStr, etimeStr, err = s.ap.AppealStarInfo(c, mid, cid, appeal.Business, ip); err != nil {
		log.Error("s.ap.AppealStarInfo(%d,%d) error(%v)", mid, cid, err)
		err = nil
	}
	if starStr != "" {
		star, err = strconv.ParseInt(starStr, 10, 64)
		if err != nil {
			log.Error("strconv.Atoi(%s) applealid(%d) error(%v)", starStr, cid, err)
			err = nil
		}
		ap.Star = int8(star)
	}
	if etimeStr != "" {
		etime, err = strconv.ParseInt(etimeStr, 10, 64)
		if err != nil {
			log.Error("strconv.Atoi(%s) applealid(%d) error(%v)", etimeStr, cid, err)
			err = nil
		}
		ap.MTime = xtime.Time(etime)
	}
	if av, err = s.arc.View(c, mid, ap.Oid, ip, 0, 0); err != nil {
		log.Error("s.arc.View(%d,%d) error(%v)", mid, ap.Oid, err)
		err = ecode.ArchiveNotExist
		return
	}
	if av != nil && av.Archive != nil {
		arc = &api.Arc{
			Aid:   av.Archive.Aid,
			Pic:   coverURL(av.Archive.Cover),
			Title: av.Archive.Title,
			State: int32(av.Archive.State),
		}
	}
	if arc == nil {
		arc = &api.Arc{}
	}
	ap.Title = arc.Title
	ap.Aid = arc.Aid
	ap.Archive = arc
	eventTmp := make([]*appeal.Event, 0, len(apmeta.Events))
	for _, v := range apmeta.Events { // Attachments 过滤管理员备注.
		if v.Event == 2 {
			continue
		}
		apev := &appeal.Event{}
		apev.ID = v.ID
		apev.AdminID = v.Adminid
		apev.Content = v.Content
		apev.ApID = v.Cid
		apev.Pics = v.Attachments
		apev.Event = v.Event
		apev.Attachments = v.Attachments
		apev.CTime = v.CTime
		apev.MTime = v.MTime
		eventTmp = append(eventTmp, apev)
	}
	ap.Events = eventTmp
	var strTmp []string
	for _, v := range ap.Attachments {
		strTmp = append(strTmp, v.Path)
	}
	ap.Pics = strings.Join(strTmp, ";")
	if ap.State == appeal.StateNoRead {
		if err = s.ap.AppealState(c, mid, cid, appeal.Business, appeal.StateReply, ip); err != nil {
			log.Error("s.ap.AppealState error(%v)", err)
			err = nil
		}
	}
	var (
		pf *model.Profile
	)
	if pf, err = s.acc.Profile(c, mid, ip); err != nil {
		log.Error("s.acc.Profile(%d) mid(%d)|ip(%s)|error(%v)", mid, ip, err)
		return
	}
	if pf != nil {
		ap.UserInfo = &appeal.UserInfo{
			MID:   pf.Mid,
			Name:  pf.Name,
			Sex:   pf.Sex,
			Face:  pf.Face,
			Rank:  pf.Rank,
			Level: pf.Level,
		}
	}
	return
}

// State shutdown an appeal.
func (s *Service) State(c context.Context, mid, cid, state int64, ip string) (err error) {
	if err = s.ap.AppealState(c, mid, cid, appeal.Business, state, ip); err != nil {
		log.Error("s.ap.AppealState error(%v)", err)
	}
	return
}

// Add create an appeal.
func (s *Service) Add(c context.Context, mid, aid int64, qq, phone, email, desc, attachments, ip string, ap *appeal.BusinessAppeal) (apID int64, err error) {
	arc, err := s.arc.View(c, mid, aid, ip, 0, 0)
	if err != nil {
		log.Error("s.arc.Archive error(%v)", err)
		err = ecode.CreativeArcServiceErr
		return
	}
	if arc == nil {
		log.Error("archive not exist")
		err = ecode.ArchiveNotExist
		return
	}
	if arc.Archive.Mid != mid {
		log.Error("login mid(%d) and archive mid(%d) are different ", mid, arc.Archive.Mid)
		err = ecode.AppealOwner
		return
	}
	if !appeal.Allow(arc.Archive.State) {
		log.Error("archive aid(%d) mid(%d) state(%d)  ", aid, mid, arc.Archive.State)
		err = ecode.AppealLimit
		return
	}
	appeals, err := s.ap.AppealList(c, mid, appeal.Business, ip)
	if err != nil {
		log.Error("s.ap.AppealList error (%v)", err)
		return
	}
	for _, k := range appeals {
		if aid == k.Oid && appeal.IsOpen(k.BusinessState) {
			err = ecode.AppealOpen
			return
		}
	}
	var tid int64
	tid, err = s.tag.AppealTag(c, aid, ip)
	if err != nil {
		log.Error("s.tag.AppealTag error(%v)", err)
		return
	}
	if tid == 0 && s.appealTag != 0 {
		tid = s.appealTag
	}
	if tid == 0 {
		log.Error("s.tag.AppealTag tid(%d)", tid)
		return
	}
	if apID, err = s.ap.AddAppeal(c, tid, aid, mid, appeal.Business, qq, phone, email, desc, attachments, ip, ap); err != nil {
		log.Error("s.ap.AddAppeal error(%v)", err)
	}
	return
}

// Reply add reply an appeal.
func (s *Service) Reply(c context.Context, mid, cid, event int64, content, attachments, ip string) (err error) {
	if err = s.ap.AddReply(c, cid, event, content, attachments, ip); err != nil {
		log.Error("s.ap.AddReply error(%v)", err)
		return
	}
	if err = s.ap.AppealState(c, mid, cid, appeal.Business, appeal.StateCreate, ip); err != nil {
		log.Error("user add reply s.ap.AppealState error(%v)", err)
		err = nil
	}
	return
}

// PhoneEmail  get user phone & email
func (s *Service) PhoneEmail(c context.Context, ck, ip string) (ct *appeal.Contact, err error) {
	if ct, err = s.acc.PhoneEmail(c, ck, ip); err != nil {
		log.Error("s.acc.PhoneEmail error(%v)", err)
	}
	if ct == nil {
		err = ecode.NothingFound
	}
	return
}

// Star give star to appeal.
func (s *Service) Star(c context.Context, mid, cid, star int64, ip string) (err error) {
	if err = s.ap.AppealExtra(c, mid, cid, appeal.Business, star, "star", ip); err != nil {
		log.Error("s.ap.AppealExtra error(%v)", err)
		return
	}
	if err = s.ap.AppealExtra(c, mid, cid, appeal.Business, time.Now().Unix(), "etime", ip); err != nil {
		log.Error("s.ap.AppealExtra error(%v)", err)
		return
	}
	if s.ap.AppealState(c, mid, cid, appeal.Business, appeal.StateUserFinished, ip); err != nil {
		log.Error("star change stats s.ap.AppealState error(%v)", err)
		err = nil
	}
	return
}

// coverURL convert cover url to full url.
func coverURL(uri string) (cover string) {
	if uri == "" {
		return
	}
	cover = uri
	if strings.Index(uri, "http://") == 0 {
		return
	}
	if len(uri) >= 10 && uri[:10] == "/templets/" {
		return
	}
	if strings.HasPrefix(uri, "group1") {
		cover = "http://i0.hdslb.com/" + uri
		return
	}
	if pos := strings.Index(uri, "/uploads/"); pos != -1 && (pos == 0 || pos == 3) {
		cover = uri[pos+8:]
	}
	cover = strings.Replace(cover, "{IMG}", "", -1)
	cover = "http://i" + strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(cover)))%3, 10) + ".hdslb.com" + cover
	return
}
