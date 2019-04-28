package service

import (
	"context"
	"go-common/app/interface/main/creative/model/tag"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/app/interface/main/videoup/model/mission"
	"go-common/app/interface/main/videoup/model/porder"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	_emptyUnicodeRegForTitle = []*regexp.Regexp{
		regexp.MustCompile(`[\x{202e}]+`),  // right-to-left override
		regexp.MustCompile(`[\x{200b}]+`),  // zeroWithChar
		regexp.MustCompile(`[\x{1f6ab}]+`), // no_entry_sign
		regexp.MustCompile(`[\n]+`),        // newline
		regexp.MustCompile(`[\r]+`),        // newline
	}
	_emptyUnicodeReg = []*regexp.Regexp{
		regexp.MustCompile(`[\x{202e}]+`),  // right-to-left override
		regexp.MustCompile(`[\x{200b}]+`),  // zeroWithChar
		regexp.MustCompile(`[\x{1f6ab}]+`), // no_entry_sign
	}
	_nocharReg = []*regexp.Regexp{
		regexp.MustCompile(`[\p{Hangul}]+`),  // kr
		regexp.MustCompile(`[\p{Tibetan}]+`), // tibe
		regexp.MustCompile(`[\p{Arabic}]+`),  // arabic
	}
	_filenameReg  = regexp.MustCompile(`^[A-Z0-9a-z]+$`)              // only letter digital.
	_staffNameReg = regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9]+$") //职能：数字、字母、中文
)

func (s *Service) checkMission(c context.Context, ap *archive.ArcParam) (err error) {
	if ap.MissionID <= 0 {
		log.Warn("MissionID(%d) error", ap.MissionID)
		ap.MissionID = 0
		return
	}
	missionID := ap.MissionID
	tid := ap.TypeID
	m, ok := s.missCache[missionID]
	if !ok || m.ID == 0 {
		err = ecode.VideoupMissionErr
		return
	}
	if m.ETime.Before(time.Now()) {
		log.Error("VideoupMissionEtimeInvalid err, tid(%d)|etime(%+v)", tid, m.ETime)
		err = ecode.VideoupMissionEtimeInvalid
		return
	}
	var missionTys map[int]*mission.Mission
	if missionTys, err = s.miss.MissionOnlineByTid(c, tid); err != nil {
		log.Error("MissionOnlineByTid err, s.tid(%+v)|err(%+v)", tid, err)
		err = nil
		return
	}
	// case: 这是活动全下架，还想参加活动
	if len(missionTys) == 0 {
		log.Error("missionTys already empty, tid(%d)|missionTys(%+v)", tid, missionTys)
		err = ecode.VideoupMissionNoMatch
		return
	}
	// 包含对分区无限制的活动
	if _, ok := missionTys[missionID]; !ok {
		log.Error("VideoupMissionNoMatch err, tid(%d)|missionID(%d)", tid, missionID)
		err = ecode.VideoupMissionNoMatch
		return
	}
	if ap.Copyright == archive.CopyrightCopy {
		log.Error("VideoupCopyForbidJoinMission err, copyright(%d)|missionID(%d)", ap.Copyright, ap.MissionID)
		err = ecode.VideoupCopyForbidJoinMission
		return
	}
	return
}

func (s *Service) checkMissionTag(srcTag string, missionID int) (dstTag string, err error) {
	var ctags = strings.Split(srcTag, ",")
	dstTag = srcTag
	// 交叉对比剔除掉当前活动已经使用的tag内容
	if len(s.missTagsCache) != 0 {
		var tags = make([]string, 0, len(ctags))
		for _, t := range ctags {
			if _, ok := s.missTagsCache[t]; !ok {
				tags = append(tags, t)
			}
		}
		dstTag = strings.Join(tags, ",")
	}
	// 校验
	// 两种情况报错提示用户: 1. 未参加活动,只提交一个tag，且是活动tag, 2. 参加活动，只提交了一个tag，且是其他活动的活动tag
	if dstTag == "" && missionID == 0 {
		log.Error("forbidMissionTagWithoutJoinMission srcTag(%s), MissionID(%d),s.missCache(%+v),s.missTagsCache(%+v)", srcTag, missionID, s.missCache, s.missTagsCache)
		err = ecode.VideoupTagForbidNotJoinMission
		return
	}
	// 未参加活动，就直接返回校验后的tag内容
	m, ok := s.missCache[missionID]
	if !ok {
		return
	}
	// 如果参加了当前有效的活动，就会把对应的活动第一个tag拼接在头部, 对于活动id和tag不匹配的会做校验
	var singleMissionTag string
	if m.Tags != "" {
		singleMissionTag = strings.Split(m.Tags, ",")[0]
	} else {
		singleMissionTag = m.Name
	}
	if len(dstTag) > 0 {
		dstTag = singleMissionTag + "," + dstTag
	} else {
		dstTag = singleMissionTag
	}
	return
}

func (s *Service) checkVideo(ap *archive.ArcParam) (err error) {
	vds := make([]*archive.VideoParam, 0)
	fnMap := make(map[string]int)
	var ok bool
	for i, v := range ap.Videos {
		if v == nil {
			continue
		}
		if v.Title, ok = s.checkTitle(v.Title); !ok {
			newErr := ecode.VideoupVideoTitleErr
			err = ecode.Errorf(newErr, newErr.Message(), i+1)
			log.Error("ap.Videos checkTitle err(%+v)|Title(%s)", err, v.Title)
			return
		}
		if v.Desc, ok = s.checkDesc(v.Desc); !ok {
			newErr := ecode.VideoupVideoDescErr
			err = ecode.Errorf(newErr, newErr.Message(), i+1)
			log.Error("ap.Videos checkDesc err(%+v)|Desc(%s)", err, v.Desc)
			return
		}
		if ok = _filenameReg.MatchString(v.Filename); !ok {
			newErr := ecode.VideoupVideoFilenameErr
			err = ecode.Errorf(newErr, newErr.Message(), i+1)
			log.Error("ap.Videos _filenameReg err(%+v)|filename(%s)", err, v.Filename)
			return
		}
		if v.Cid == 0 && v.Filename == "" { // NOTE: cid>0 means code mode
			newErr := ecode.VideoupVideoFilenameErr
			err = ecode.Errorf(newErr, newErr.Message(), i+1)
			log.Error("ap.Videos err(%+v)|Filename(%s)|Cid(%d)", err, v.Filename, v.Cid)
			return
		}
		if _, ok := fnMap[v.Filename]; ok {
			err = ecode.VideoupFilenameCanotRepeat
			log.Error("ecode.VideoupFilenameCanotRepeat err(%+v)|Filename(%s)|index(%d)", err, v.Filename, i)
			return
		}
		vds = append(vds, v)
		fnMap[v.Filename] = 1
	}
	ap.Videos = vds
	return
}

func (s *Service) checkCover(cover string) (cv string, ok bool) {
	if cover == "" {
		ok = true
		return
	}
	uri, err := url.Parse(cover)
	if err != nil {
		return
	}
	if strings.Contains(uri.Host, "hdslb.com") {
		cv = uri.Path
		ok = true
		return
	} else if strings.Contains(uri.Host, "acgvideo.com") {
		cv = cover
		ok = true
		return
	}
	return
}

func (s *Service) checkDynamicLen233(dynamic string) (dyn string, ok bool) {
	dyn = strings.TrimSpace(dynamic)
	var _emptyDynUnicodeReg = []*regexp.Regexp{
		regexp.MustCompile(`[\x{FFFC}]+`), // obj
	}
	for _, reg := range _emptyDynUnicodeReg {
		dyn = reg.ReplaceAllString(dyn, "")
	}
	if utf8.RuneCountInString(dyn) > 233 {
		return
	}
	ok = true
	return
}

func (s *Service) checkTitle(title string) (ct string, ok bool) {
	ct = strings.TrimSpace(title)
	if utf8.RuneCountInString(ct) > 80 {
		return
	}
	for _, reg := range _nocharReg {
		if reg.MatchString(ct) {
			return
		}
	}
	for _, reg := range _emptyUnicodeRegForTitle {
		ct = reg.ReplaceAllString(ct, "")
	}
	ok = true
	return
}

func (s *Service) checkDesc(desc string) (cd string, ok bool) {
	cd = strings.TrimSpace(desc)
	for _, reg := range _emptyUnicodeReg {
		cd = reg.ReplaceAllString(cd, "")
	}
	if utf8.RuneCountInString(cd) > 2000 {
		return
	}
	ok = true
	return
}

func (s *Service) checkAccount(c context.Context, mid int64, ip string) (p *accapi.Profile, err error) {
	if p, err = s.acc.Profile(c, mid, ip); err != nil {
		return
	}
	if p.Silence == 1 {
		err = ecode.UserDisabled
	} else if p.Level < 1 {
		err = ecode.UserLevelLow
	}
	if _, ok := s.exemptZeroLevelAndAnswerUps[mid]; ok && err == ecode.UserLevelLow {
		log.Info("s.exemptZeroLevelAndAnswerUps, (%s),(%d),(%+v)", ip, mid, err)
		err = nil
	}
	return
}

func (s *Service) checkOrderID(c context.Context, mid, orderID int64, ip string) (err error) {
	orderIDs, err := s.order.ExecuteOrders(c, mid, ip)
	if err != nil {
		log.Error("s.order.ExecuteOrders mid(%d) ip(%s) error(%v)", mid, err)
		err = ecode.VideoupOrderAPIErr
		return
	}
	if _, ok := orderIDs[orderID]; !ok {
		err = ecode.VideoupOrderIDNotAllow
	}
	return
}

func (s *Service) checkIdentify(c context.Context, mid int64, ip string) (err error) {
	if _, ok := s.exemptIDCheckUps[mid]; ok {
		log.Info("s.exemptIDCheckUps, (%s),(%d),(%+v)", ip, mid, err)
		return
	}
	// fault-tolerant for service interruption
	if err = s.acc.IdentifyInfo(c, ip, mid); err != nil {
		if err != ecode.UserCheckNoPhone && err != ecode.UserCheckInvalidPhone {
			log.Warn("s.accIdentifyInfo, account service maybe in interruption,(%s),(%d),(%+v)", ip, mid, err)
			return nil
		}
		log.Error("s.accIdentifyInfo, (%s),(%d),(%+v)", ip, mid, err)
		return
	}
	return
}

// checkPorderForAdd
func (s *Service) checkPorderForAdd(c context.Context, ap *archive.ArcParam, mid int64) (err error) {
	// 防止脏数据, 强制计算和校验Porder的数据
	ap.Porder.FlowID = 1
	log.Info("ap.Porder (%+v)", ap.Porder)
	// showType check
	if len(ap.Porder.ShowType) > 0 {
		var showTypes []int64
		if showTypes, err = xstr.SplitInts(ap.Porder.ShowType); err != nil {
			log.Error("SplitInts ShowType err, (%s),(%+v)", ap.Porder.ShowType, err)
			return
		}
		//广告的展现形式太多或者太少
		if len(showTypes) == 0 {
			err = ecode.VideoupAdShowTypeErr
			log.Error("check showTypes (%+v)|err(%+v)", ap.Porder, err)
			return
		}
		for _, showType := range showTypes {
			if showType > 0 {
				if _, ok := s.PorderCfgs[showType]; !ok {
					err = ecode.VideoupAdShowTypeErr
					log.Error("VideoupAdShowTypeErr Porder(%+v)|err(%+v)", ap.Porder.ShowType, err)
					return
				}
			}
		}
	}
	// Official check
	if ap.Porder.Official == 1 {
		if _, ok := porder.OfficialIndustryMaps[ap.Porder.IndustryID]; !ok {
			err = ecode.VideoupAdOfficialIndustryIDErr
			log.Error("VideoupAdOfficialIndustryIDErr Porder(%+v)|err(%+v)", ap.Porder, err)
			return
		}
		if ap.Porder.BrandID < 0 {
			err = ecode.VideoupAdBrandIDErr
			log.Error("VideoupAdBrandIDErr Porder(%+v)|err(%+v)", ap.Porder, err)
			return
		}
		// logic map to OfficialIndustryMaps waiting for add other official industry
		if _, ok := s.PorderGames[ap.Porder.BrandID]; !ok {
			err = ecode.VideoupAdBrandIDErr
			log.Error("VideoupAdBrandIDErr Porder(%+v)|err(%+v)", ap.Porder, err)
			return
		}
	}
	// Industry check
	if ap.Porder.IndustryID > 0 {
		if _, ok := s.PorderCfgs[ap.Porder.IndustryID]; !ok {
			err = ecode.VideoupAdIndustryIDErr
			log.Error("VideoupAdIndustryIDErr Porder(%+v)|err(%+v)", ap.Porder, err)
			return
		}
	}
	return
}

// checkDescForLength fn
func (s *Service) checkDescForLength(desc string, descFormatID int, typeID int16, copyright int8) (err error) {
	if descFormatID == 0 {
		if utf8.RuneCountInString(desc) > 250 {
			err = ecode.VideoupFmDesLenOverLimit
			log.Error("ecode.VideoupFmDesLenOverLimit, desc(%s),formatID(%d)", desc, descFormatID)
			return
		}
		return
	}
	if utf8.RuneCountInString(desc) > 2000 {
		err = ecode.VideoupFmDesLenOverLimit
		log.Error("ecode.VideoupFmDesLenOverLimit, desc(%s),formatID(%d)", desc, descFormatID)
		return
	}
	return
}

// checkVideos, 1: check len(video) == ? 0 , 2: check typeID in (145,146,147,83) and with multi Videos
func (s *Service) checkVideos(c context.Context, ap *archive.ArcParam) (err error) {
	if len(ap.Videos) == 0 {
		log.Error("checkVideos vds length 0")
		err = ecode.VideoupZeroVideos
		return
	}
	if len(ap.Videos) > 1 && ap.ForbidMultiVideoType() {
		log.Error("checkVideos vds ForbidMultiVideoType, len(%d), type(%d) ", len(ap.Videos), ap.TypeID)
		err = ecode.VideoupForbidMultiVideoForTypes
		return
	}
	return
}

// tagsCheck fn
func (s *Service) tagsCheck(c context.Context, mid int64, tagName, ip string) (err error) {
	var t *tag.Tag
	tags := strings.Split(tagName, ",")
	for i, tagStr := range tags {
		if t, err = s.tag.TagCheck(c, mid, tagStr); err != nil {
			log.Error("s.tag.TagCheck(%d, %+v, %s) error(%+v)", mid, t, ip, err)
			err = nil
			return
		}
		if t != nil && (t.State == tag.TagStateDel || t.State == tag.TagStateHide || t.Type == tag.OfficailActiveTag) {
			newErr := ecode.VideoupTagForbid
			err = ecode.Errorf(newErr, newErr.Message(), i+1)
			log.Error("s.tag.VideoupTagForbid (%d, %+v, %s) error(%+v)", mid, t, ip, err)
			return
		}
	}
	return
}

// checkVideosMaxLimitForEdit fn
func (s *Service) checkVideosMaxLimitForEdit(vs []*archive.Video, pvideos []*archive.VideoParam) (addCnt int) {
	fnMaps := make(map[string]string)
	for _, v := range vs {
		fnMaps[v.Filename] = v.Filename
	}
	for _, v := range pvideos {
		if _, exist := fnMaps[v.Filename]; len(v.Filename) > 0 && !exist {
			addCnt++
		}
	}
	return
}

// checkPay fn
func (s *Service) checkAddPay(c context.Context, ap *archive.ArcParam, ip string) (err error) {
	if err = s.checkPayProtocol(c, ap.Pay, ap.Mid); err != nil {
		log.Error("s.checkAddPayProtocol (ap %+v) error(%+v)", ap, err)
		return
	}
	if err = s.checkPayLimit(c, ap); err != nil {
		log.Error("s.checkAddPayLimit (ap %+v) error(%+v)", ap, err)
		return
	}
	if err = s.checkPayWithOrder(c, ap.Porder, ap.Pay, ap.OrderID, ap.Mid); err != nil {
		log.Error("s.checkAddPayWithOrder (ap %+v) error(%+v)", ap, err)
		return
	}
	return
}

// checkEditPay fn
// 关于RefuseUpdate
// 1.当ctime==ptime的时候，也就是尚未稿件一审，都是可以update付费模块的
// 2.开放过的，打回的稿件可以在任何时间点进行付费修改
// 3.开放过的，非打回的稿件想要修改只能等60天时间过了之后(时间锁是为了保护普通用户的收看权利)，否则自己申诉付费审核人员，要求强制打回
func (s *Service) checkEditPay(c context.Context, ap *archive.ArcParam, a *archive.Archive, ip string) (err error) {
	_, registed, _ := s.pay.Ass(c, a.Aid, ip)
	// 只要注册过付费信息都允许自主修改付费模块
	if registed {
		if err = s.checkPayProtocol(c, ap.Pay, ap.Mid); err != nil {
			log.Error("s.checkAddPayProtocol (ap %+v) error(%+v)", ap, err)
			return
		}
		if err = s.checkPayLimit(c, ap); err != nil {
			log.Error("s.checkAddPayLimit (ap %+v) error(%+v)", ap, err)
			return
		}
		if err = s.checkPayWithOrder(c, ap.Porder, ap.Pay, ap.OrderID, ap.Mid); err != nil {
			log.Error("s.checkAddPayWithOrder (ap %+v) error(%+v)", ap, err)
			return
		}
		//如果其他端都不传付费信息，那么就用现在有的来进行最后覆盖
		if ap.Pay == nil {
			ap.UgcPay = a.UgcPay
			return
		}
		updateDeadLine := xtime.Time(a.PTime.Time().AddDate(0, 0, s.c.UgcPayAllowEditDays).Unix())
		if ap.Pay != nil &&
			a.CTime != a.PTime &&
			a.State != archive.StateForbidRecicle &&
			xtime.Time(time.Now().Unix()) < updateDeadLine {
			log.Warn("checkEditPay ap.Pay.RefuseUpdate updateDeadLine (%+v)|(%+v)|(%+v)", a.Aid, a.CTime, a.PTime)
			ap.Pay.RefuseUpdate = true
		}
	}
	return
}
