package service

import (
	"context"
	"go-common/app/interface/main/videoup/model/archive"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"html"
	"strings"
	"time"
)

func (s *Service) preMust(c context.Context, mid int64, ap *archive.ArcParam, ip string, upFrom int8) (err error) {
	// title xss filter
	originTitleLen := len(ap.Title)
	ap.Title = html.UnescapeString(xssFilter(ap.Title))
	if len(ap.Title) != originTitleLen {
		log.Warn("ap.Title inject by xss:mid(%d)|ip(%d)", mid, ip)
	}
	// check videos
	if upFrom != archive.UpFromAPP {
		if err = s.checkVideos(c, ap); err != nil {
			log.Error("s.checkVideos mid(%d) ap.TypeID(%d) err(%+v)", mid, ap.TypeID, err)
			return
		}
	}
	// check archive
	if !s.allowType(ap.TypeID) {
		log.Error("s.allowType mid(%d) ap.TypeID(%d) typeid not exists", mid, ap.TypeID)
		err = ecode.VideoupTypeidErr
		return
	}
	if !s.allowCopyright(ap.Copyright) {
		log.Error("s.allowCopyright mid(%d) ap.Copyright(%d) no legal copyright", mid, ap.Copyright)
		err = ecode.VideoupCopyrightErr
		return
	}
	ap.Tag = s.removeDupTag(ap.Tag)
	if !s.allowTag(ap.Tag) {
		log.Error("s.allowTag mid(%d) ap.Tag(%s) tag name or number too large or Empty", mid, ap.Tag)
		err = ecode.VideoupTagErr
		return
	}
	if err = s.checkVideo(ap); err != nil {
		log.Error("s.checkVideo mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	var ok bool
	if ap.Cover, ok = s.checkCover(ap.Cover); !ok {
		log.Error("s.checkCover mid(%d) ap.Cover(%s) cover no legal", mid, ap.Cover)
		err = ecode.VideoupCoverErr
		return
	}
	if ap.Title, ok = s.checkTitle(ap.Title); !ok || ap.Title == "" {
		log.Error("s.checkTitle mid(%d) ap.Title(%s) title contains legal char or is empty", err, mid, ap.Title)
		err = ecode.VideoupTitleErr
		return
	}
	if ap.Dynamic, ok = s.checkDynamicLen233(ap.Dynamic); !ok {
		log.Error("s.checkDynamic err(%+v) mid(%d) ap.Dynamic(%s) contains length larger 233", err, mid, ap.Dynamic)
		err = ecode.VideoupDynamicErr
		return
	}
	if ap.Desc, ok = s.checkDesc(ap.Desc); !ok {
		log.Error("s.checkDesc mid(%d) ap.Desc(%s) desc contains legal char or is empty", mid, ap.Desc)
		err = ecode.VideoupDescErr
		return
	}
	var p *accapi.Profile
	if p, err = s.checkAccount(c, mid, ip); err != nil {
		log.Error("s.checkAccount mid(%d) error(%v)", mid, err)
		return
	}
	ap.Author = p.Name
	if ap.Copyright == archive.CopyrightCopy {
		ap.NoReprint = 0
	}
	// DisableVideoDesc except UpFromWindows, step 1 for all
	if upFrom != archive.UpFromWindows {
		for _, v := range ap.Videos {
			v.Desc = ""
		}
	}
	// 防止脏数据
	if ap.Vote != nil && ap.Vote.VoteID == 0 {
		ap.Vote = nil
	}
	return
}

func (s *Service) preOrder(c context.Context, ap *archive.ArcParam, a *archive.Archive, ip string) (err error) {
	if ap.Porder != nil && ap.Porder.FlowID > 0 && ap.OrderID > 0 {
		err = ecode.VideoupPvodForbidOrderAlready
		return
	}
	if ap.OrderID < 0 {
		err = ecode.VideoupOrderIDNotAllow
		return
	}
	if ap.Aid == 0 && ap.OrderID == 0 { // NOTE: add no orderid
		return
	}
	if ap.Aid > 0 && ap.OrderID == 0 && a.OrderID == 0 { // NOTE: edit always no orderid
		return
	}
	if ap.Aid > 0 { // NOTE: edit had order id, not allow change
		ap.OrderID = a.OrderID
		ap.DTime = a.DTime
		return
	}
	if !s.allowOrderUps(ap.Mid) {
		log.Error("s.allowOrderUps mid(%d) error(%v)", ap.Mid, err)
		err = ecode.VideoupUperIDNotAllow
		return
	}
	if err = s.checkOrderID(c, ap.Mid, ap.OrderID, ip); err != nil {
		return
	}
	var ptime xtime.Time
	if ptime, err = s.order.PubTime(c, ap.Mid, ap.OrderID, ip); err != nil {
		err = ecode.VideoupOrderAPIErr
		return
	}
	if ap.Aid == 0 && int64(ptime) < time.Now().Add(2*time.Hour).Unix() {
		err = ecode.VideoupLaunchTimeIllegal
		return
	}
	ap.DTime = ptime
	return
}

func (s *Service) preAdd(c context.Context, mid int64, ap *archive.ArcParam, ip string, upFrom int8) (err error) {
	if ap.ForbidAddVideoType() {
		err = ecode.VideoupTypeidErr
		log.Error("ap.ForbidAddVideoType VideoupTypeidErr mid(%d),type(%d),err(%v) ", mid, ap.TypeID, err)
		return
	}
	if ap.ForbidCopyrightAndTypes() {
		err = ecode.VideoupCopyrightErr
		log.Error("ap.ForbidCopyrightAndTypes VideoupCopyrightErr mid(%d),copyright(%d),type(%d),err(%v) ", mid, ap.Copyright, ap.TypeID, err)
		return
	}
	if len(ap.Videos) > s.c.MaxAddVsCnt {
		err = ecode.VideoupVideosMaxLimit
		log.Error("ap.VideoupVideosMaxLimit current(%d), max(%d),err(%v) ", len(ap.Videos), s.c.MaxAddVsCnt, err)
		return
	}
	if !s.allowSource(ap.Copyright, ap.Source) {
		err = ecode.VideoupSourceErr
		return
	}
	originDesc := ap.Desc
	//join source and desc for CopyrightCopy with \n
	if ap.Copyright == archive.CopyrightCopy && len(strings.TrimSpace(ap.Source)) > 0 {
		ap.Desc = ap.Source + "\n" + ap.Desc
	}
	// App端允许在添加和编辑稿件的时候简介为空,但是需要区分是操作系统平台
	ap.Desc = s.switchDesc(upFrom, ap.Desc)
	// preMust method must be first
	if err = s.preMust(c, mid, ap, ip, upFrom); err != nil {
		log.Error("s.preMust mid(%d), err(%v) ", mid, err)
		return
	}
	if !s.allowRepeat(c, mid, ap.Title) {
		err = ecode.VideoupCanotRepeat
		return
	}
	if !s.allowDelayTime(ap.DTime) {
		err = ecode.VideoupDelayTimeErr
		return
	}
	if err = s.checkMission(c, ap); err != nil {
		log.Error("s.checkMission mid(%d) ap.MissionID(%d)|TypeID(%d) missionId not exists", mid, ap.MissionID, ap.TypeID)
		return
	}
	if ap.Tag, err = s.checkMissionTag(ap.Tag, ap.MissionID); err != nil {
		log.Error("s.checkMissionTag mid(%d) ap.tag(%s) ap.MissionID(%d) missionId not exists", mid, ap.Tag, ap.MissionID)
		return
	}
	if err = s.checkDescForLength(originDesc, ap.DescFormatID, ap.TypeID, ap.Copyright); err != nil {
		log.Error("s.checkDescForLength mid(%d) ap.Source(%s), apDesc(%s),ap.DescFormatID(%d) ap.Lang(%d) err(%v)", mid, ap.Source, originDesc, ap.DescFormatID, ap.Lang, err)
		return
	}
	if err = s.preOrder(c, ap, nil, ip); err != nil {
		log.Error("s.preOrder mid(%d) ap(%v), err(%v) ", mid, ap, err)
		return
	}
	//checkPorderForAdd
	if ap.Porder != nil && ap.Porder.IndustryID > 0 {
		if err = s.checkPorderForAdd(c, ap, mid); err != nil {
			log.Error("s.checkPorderForAdd mid(%d) ap(%v) |err(%+v)", mid, ap, err)
			return
		}
	}
	ap.NilPoiObj()
	return
}

func (s *Service) switchDesc(upFrom int8, originDesc string) (resDesc string) {
	resDesc = originDesc
	if (upFrom == archive.UpFromAPP ||
		upFrom == archive.UpFromAPPAndroid ||
		upFrom == archive.UpFromIpad ||
		upFrom == archive.UpFromAPPiOS) &&
		len(originDesc) == 0 {
		resDesc = "-"
	}
	return
}

func (s *Service) preEdit(c context.Context, mid int64, a *archive.Archive, vs []*archive.Video, ap *archive.ArcParam, ip string, upFrom int8) (err error) {
	//检查联合投稿移区和修改转载类型
	if err = s.checkStaffMoveType(c, ap, a, ip); err != nil {
		return
	}
	if len(ap.Videos) > s.c.MaxAllVsCnt {
		newErr := ecode.VideoupMaxAllVsCntLimit
		err = ecode.Errorf(newErr, newErr.Message(), s.c.MaxAllVsCnt)
		log.Error("MaxAllVsCnt err(%+v)|MaxAllVsCnt(%d)|mid(%d)|aid(%d)", err, s.c.MaxAllVsCnt, mid, a.Aid)
		return
	}
	// App端允许在添加和编辑稿件的时候简介为空,但是需要区分是操作系统平台
	if ap.ForbidCopyrightAndTypes() {
		err = ecode.VideoupCopyrightErr
		log.Error("ap.ForbidCopyrightAndTypes VideoupCopyrightErr mid(%d),copyright(%d),type(%d),err(%v) ", mid, ap.Copyright, ap.TypeID, err)
		return
	}
	ap.Desc = s.switchDesc(upFrom, ap.Desc)
	if err = s.preMust(c, mid, ap, ip, upFrom); err != nil {
		log.Error("s.preMust mid(%d), err(%v) ", mid, err)
		return
	}
	// DisableVideoDesc, except UpFromWindows step 2 for edit
	if upFrom != archive.UpFromWindows {
		ap.DisableVideoDesc(vs)
	}
	ap.TypeID, ap.Copyright, ap.Tag, ap.MissionID, ap.DescFormatID = s.protectFieldForEdit(ap, a)
	//not in cache or not StateForbidRecicle
	_, ok := s.missCache[ap.MissionID]
	if (!ok && ap.MissionID > 0) || (a.State != archive.StateForbidRecicle) {
		ap.MissionID = a.MissionID
	}
	if a.State == archive.StateForbidRecicle {
		if err = s.checkMission(c, ap); err != nil {
			log.Error("s.checkMission mid(%d) ap.MissionID(%d)|TypeID(%d) missionId not exists", mid, ap.MissionID, ap.TypeID)
			return
		}
	}
	if ap.Tag, err = s.checkMissionTag(ap.Tag, ap.MissionID); err != nil {
		log.Error("s.checkMissionTag mid(%d) ap.tag(%s) ap.MissionID(%d) missionId not exists", mid, ap.Tag, ap.MissionID)
		return
	}
	// mid check
	if a.Mid != mid {
		log.Error("mid(%d) is not author(%d)", mid, a.Mid)
		err = ecode.ArchiveOwnerErr
		return
	}
	// state check
	if a.NotAllowUp() {
		err = ecode.ArchiveBlocked
		return
	}
	// web和新发粉版允许修改创作类型，其他的都不允许
	if upFrom != archive.UpFromWeb &&
		upFrom != archive.UpFromAPPiOS &&
		upFrom != archive.UpFromIpad &&
		upFrom != archive.UpFromAPPAndroid {
		ap.NoReprint = a.NoReprint
		log.Info("upfrom forbid change np, np(%d)|upfrom(%+v)", a.NoReprint, upFrom)
	}
	// NoReprint check
	if a.NoReprint == 0 && ap.NoReprint == 1 {
		log.Error("notAllow set NoReprint = 1 after now's Noreprint is 0 mid(%d) ap.NoReprint(%d) a.NoReprint(%d)", mid, ap.NoReprint, a.NoReprint)
		err = ecode.VideoupForbidNoreprint
		return
	}
	// allowDelayTime check for archive which state nq -40
	if a.State == archive.StateForbidUserDelay {
		ap.DTime = a.DTime
	} else if a.State != archive.StateOpen {
		if a.DTime != ap.DTime && !s.allowDelayTime(ap.DTime) {
			log.Error("s.allowDelayTime err(%+v) mid(%d) ap.Dtime(%d) must between 4h 15d", err, mid, ap.DTime)
			err = ecode.VideoupDelayTimeErr
			return
		}
	}
	if err = s.preOrder(c, ap, a, ip); err != nil {
		log.Error("s.preOrder mid(%d) ap(%v), err(%v) ", mid, ap, err)
		return
	}
	// checkDescForLength
	if err = s.checkDescForLength(ap.Desc, ap.DescFormatID, ap.TypeID, ap.Copyright); err != nil {
		log.Error("s.checkDescForLength mid(%d) ap.Desc(%s) ap.DescFormatID(%d) err(%v)", mid, ap.Desc, ap.DescFormatID, err)
		return
	}
	// 手动暴力禁止编辑的时候进行修改poi地理位置信息
	ap.PoiObj = nil
	// checkEditPay
	if err = s.checkEditPay(c, ap, a, ip); err != nil {
		log.Error("s.checkEditPay mid(%d) ap(%+v) a(%+v) err(%v)", mid, ap, a, err)
		return
	}
	return
}

// protectFieldForEdit only StateForbidRecicle allow change typeID and Copyright
// 简介模板的ID暂时近期内不允许修改
func (s *Service) protectFieldForEdit(ap *archive.ArcParam, a *archive.Archive) (typeID int16, copyright int8, tag string, missionID, descFormatID int) {
	if a.State == archive.StateForbidRecicle ||
		a.State == archive.StateForbidSubmit ||
		a.State == archive.StateForbidFixed ||
		a.State == archive.StateOrange ||
		a.State == archive.StateOpen {
		return ap.TypeID, ap.Copyright, ap.Tag, ap.MissionID, ap.DescFormatID
	}
	return a.TypeID, a.Copyright, a.Tag, a.MissionID, a.DescFormatID
}

func (s *Service) removeDupTag(tagStr string) string {
	result := []string{}
	elements := strings.Split(tagStr, ",")
	for i := 0; i < len(elements); i++ {
		exists := false
		for v := 0; v < i; v++ {
			if elements[v] == elements[i] {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, elements[i])
		}
	}
	return strings.Join(result, ",")
}
