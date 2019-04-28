package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	emailmdl "go-common/app/job/main/videoup-report/model/email"
	"go-common/app/job/main/videoup-report/model/manager"
	account "go-common/app/service/main/account/api"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

//一审备注中有"私单报备"， 则发送邮件
func (s *Service) sendVideoPrivateEmail(c context.Context, a *archive.Archive, v *archive.Video) (err error) {
	if a == nil || v == nil {
		return
	}
	note, err := s.arc.VideoAuditNote(c, v.ID)
	if err != nil {
		log.Error("sendVideoPrivateEmail s.arc.VideoAuditNote error(%v)", err)
		return
	}
	if !s.needPrivateEmail(a.TypeID, note) {
		log.Info("sendVideoPrivateEmail no need to send email: vid(%d), aid(%d), note(%s), typeId(%d)", v.ID, v.Aid, note, a.TypeID)
		return
	}

	log.Info("start to sendVideoPrivateEmail: note(%s), typeId(%d), aid(%d), video(%v), archive(%v)", note, a.TypeID, a.ID, v, a)
	//up主信息
	pfl, err := s.profile(c, a.Mid)
	if err != nil || pfl == nil {
		log.Error("sendVideoPrivateEmail s.profile error(%v) or nil, mid(%d), vid(%d), aid(%d)", err, a.Mid, v.ID, a.ID)
		return
	}

	//审核人员信息
	mngUID, err := s.arc.LastVideoOperUID(c, v.ID)
	if err != nil || mngUID <= 0 {
		log.Error("sendVideoPrivateEmail s.arc.LastVideoOperUID error(%v) or zero(%d) vid(%d) aid(%d)", err, mngUID, v.ID, v.Aid)
		return
	}
	mngUser, err := s.mng.User(c, mngUID)
	if err != nil || mngUser == nil || mngUser.ID != mngUID {
		log.Error("sendVideoPrivateEmail s.mng.User(%d) error(%v) or not found(%+v) uid(%d)", mngUID, err, mngUser)
		return
	}

	//禁止项状态
	noRankAttr, noDynamicAttr, noRecommend := "关", "关", "关"
	if v.AttrVal(archive.AttrBitNoRank) == archive.AttrYes {
		noRankAttr = "开"
	}
	if v.AttrVal(archive.AttrBitNoDynamic) == archive.AttrYes {
		noDynamicAttr = "开"
	}
	if v.AttrVal(archive.AttrBitNoRecommend) == archive.AttrYes {
		noRecommend = "开"
	}

	//组合邮件参数
	params := map[string]string{
		"upName":          pfl.Profile.Name,
		"aid":             strconv.FormatInt(v.Aid, 10),
		"arcTitle":        a.Title,
		"arcState":        archive.StateMean[a.State],
		"noRankAttr":      noRankAttr,
		"noDynamicAttr":   noDynamicAttr,
		"noRecommendAttr": noRecommend,
		"upFans":          strconv.FormatInt(pfl.Follower, 10),
		"uid":             strconv.FormatInt(mngUID, 10),
		"mngName":         mngUser.Username,
		"mngDepartment":   mngUser.Department,
		"note":            note,
		"typeId":          strconv.Itoa(int(s.topType(a.TypeID))),
		"emailType":       emailmdl.EmailPrivateVideo,
	}
	tpl := s.email.PrivateEmailTemplate(params)
	s.email.PushToRedis(c, tpl)
	return
}

//稿件备注中有"私单报备"， 则发送邮件
func (s *Service) sendArchivePrivateEmail(c context.Context, a *archive.Archive) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.sendArchivePrivateEmail panic(%v)", pErr)
		}
	}()
	if a == nil {
		return
	}
	note, err := s.arc.ArchiveNote(c, a.ID)
	if err != nil {
		log.Error("sendArchivePrivateEmail s.arc.ArchiveNote error(%v) aid(%d)", err, a.ID)
		return
	}
	if !s.needPrivateEmail(a.TypeID, note) {
		log.Info("sendArchivePrivateEmail: no need to send email: aid(%d), note(%s), typeId(%d)", a.ID, note, a.TypeID)
		return
	}

	log.Info("start to sendArchivePrivateEmail: note(%s), typeId(%d), aid(%d), archive(%v)", note, a.TypeID, a.ID, a)
	//up主信息
	pfl, err := s.profile(c, a.Mid)
	if err != nil || pfl == nil {
		log.Error("sendArchivePrivateEmail s.profile error(%v) or nil, mid(%d), aid(%d)", err, a.Mid, a.ID)
		return
	}

	//审核人员信息
	arcOper, err := s.arc.LastArcOper(c, a.ID)
	if err != nil || arcOper == nil || arcOper.AID <= 0 {
		log.Error("sendArchivePrivateEmail s.arc.LastArcOper(%d), error(%v) or not found(%+v)", a.ID, err, arcOper)
		return
	}
	mngUser, err := s.mng.User(c, arcOper.UID)
	if err != nil || mngUser == nil || mngUser.ID != arcOper.UID {
		log.Error("sendArchivePrivateEmail s.mng.User error(%v) or not found(%+v), aid(%d), uid(%d)", err, mngUser, a.ID, arcOper.UID)
		return
	}

	//禁止项状态
	noRankAttr, noDynamicAttr, noRecommend := "关", "关", "关"
	if a.AttrVal(archive.AttrBitNoRank) == archive.AttrYes {
		noRankAttr = "开"
	}
	if a.AttrVal(archive.AttrBitNoDynamic) == archive.AttrYes {
		noDynamicAttr = "开"
	}
	if a.AttrVal(archive.AttrBitNoRecommend) == archive.AttrYes {
		noRecommend = "开"
	}

	//组合邮件参数
	params := map[string]string{
		"upName":          pfl.Profile.Name,
		"aid":             strconv.FormatInt(a.ID, 10),
		"arcTitle":        a.Title,
		"arcState":        archive.StateMean[a.State],
		"noRankAttr":      noRankAttr,
		"noDynamicAttr":   noDynamicAttr,
		"noRecommendAttr": noRecommend,
		"upFans":          strconv.FormatInt(pfl.Follower, 10),
		"uid":             strconv.FormatInt(mngUser.ID, 10),
		"mngName":         mngUser.Username,
		"mngDepartment":   mngUser.Department,
		"note":            note,
		"typeId":          strconv.Itoa(int(s.topType(a.TypeID))),
		"emailType":       emailmdl.EmailPrivateArchive,
	}
	tpl := s.email.PrivateEmailTemplate(params)
	s.email.PushToRedis(c, tpl)
	return
}

//触发私单报备邮件的条件
func (s *Service) needPrivateEmail(typeID int16, note string) (matched bool) {
	typeI := int(s.topType(typeID))
	_, exist := s.email.PrivateAddr[strconv.Itoa(typeI)]
	matched = exist && strings.Contains(note, "私单报备")
	return
}

func (s *Service) sendMail(c context.Context, a *archive.Archive, v *archive.Video) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.sendMail panic(%v)", pErr)
		}
	}()
	var (
		condition, additionCondition, content []string
		mngUser                               *manager.User
		pfl                                   *account.ProfileStatReply
		remark, additionType                  string
		mngUID                                int64
		videoOper                             *archive.VideoOper
		arcOper                               *archive.Oper
	)
	//触发条件
	if s.isSigned(a.Mid) {
		additionCondition = append(additionCondition, "签约UP主")
		additionType = "signed"
	}
	if s.isWhite(a.Mid) {
		condition = append(condition, "优质UP主")
	}
	if s.isPolitices(a.Mid) {
		condition = append(condition, "时政UP主")
	}
	if s.isEnterprise(a.Mid) {
		condition = append(condition, "企业机构实名认证UP主")
	}
	if pfl, _ = s.profile(c, a.Mid); pfl != nil && pfl.Follower > 100000 {
		condition = append(condition, fmt.Sprintf("10万以上粉丝（%d个粉丝）", pfl.Follower))
	}
	if len(condition)+len(additionCondition) <= 0 {
		log.Info("sendMail aid(%d) mid(%d) is not signed uper & whith uper && shizheng uper && qiye uper && funs < 10W", a.ID, a.Mid)
		return
	}

	content = append(content, fmt.Sprintf("[审核状态]: %s", archive.StateMean[a.State]))
	fromVideo := v != nil
	if fromVideo {
		//视频审核操作
		videoOper, err = s.arc.LastVideoOper(c, v.ID)
		if err != nil || videoOper == nil || videoOper.VID <= 0 {
			log.Error("sendMail s.arc.LastVideoOper(vid(%d)) error(%v) or not found", v.ID, err)
			return
		}

		mngUID = videoOper.UID
		if len(videoOper.Content) != 0 {
			content = append(content, videoOper.Content)
		}
		remark = strings.TrimSpace(videoOper.Remark)
		if len(remark) != 0 {
			content = append(content, fmt.Sprintf("备注：%s", strings.Replace(remark, "\n", ",", -1)))
		}
	} else {
		//稿件审核操作
		if arcOper, err = s.arc.LastArcOper(c, a.ID); err != nil || arcOper == nil || arcOper.AID <= 0 {
			log.Error("sendMail s.arc.LastArcOper(aid(%d)) error(%v) or not found", a.ID, err)
			return
		}
		mngUID = arcOper.UID
		if len(arcOper.Content) != 0 {
			content = append(content, arcOper.Content)
		}
		remark = strings.TrimSpace(arcOper.Remark)
		if len(remark) != 0 {
			content = append(content, fmt.Sprintf("备注: %s", strings.Replace(remark, "\n", ",", -1)))
		}
	}

	//审核人员信息
	if mngUser, err = s.mng.User(c, mngUID); err != nil || mngUser == nil || mngUser.ID != mngUID {
		log.Error("s.mng.User(%d) error(%v) or not found(%d)", arcOper.UID, err, mngUser)
		return
	}

	if len(condition) > 0 {
		params := map[string]string{
			"aid":        strconv.FormatInt(a.ID, 10),
			"title":      a.Title,
			"upName":     pfl.Profile.Name,
			"condition":  strings.Join(condition, "/"),
			"change":     strings.Join(content, " ,"),
			"uid":        strconv.FormatInt(mngUser.ID, 10),
			"username":   mngUser.Username,
			"department": mngUser.Department,
			"typeId":     strconv.Itoa(int(s.topType(a.TypeID))),
			"fromVideo":  strconv.FormatBool(fromVideo),
		}
		tpl := s.email.NotifyEmailTemplate(params)
		s.email.PushToRedis(c, tpl)
	}
	if len(additionCondition) > 0 {
		params := map[string]string{
			"aid":        strconv.FormatInt(a.ID, 10),
			"title":      a.Title,
			"upName":     pfl.Profile.Name,
			"condition":  strings.Join(additionCondition, "/"),
			"change":     strings.Join(content, " ,"),
			"uid":        strconv.FormatInt(mngUser.ID, 10),
			"username":   mngUser.Username,
			"department": mngUser.Department,
			"typeId":     additionType,
			"fromVideo":  strconv.FormatBool(fromVideo),
		}
		tpl := s.email.NotifyEmailTemplate(params)
		s.email.PushToRedis(c, tpl)
	}
	return
}

func (s *Service) emailProc() {
	defer s.waiter.Done()

	for {
		if s.closed {
			return
		}

		s.email.Start(emailmdl.MailKey)
		time.Sleep(200 * time.Millisecond)
	}
}

func (s *Service) emailFastProc() {
	defer s.waiter.Done()
	var (
		err error
	)

	for {
		if s.closed {
			return
		}

		<-s.email.FastChan()
		log.Info("emailFastProc start to handle")

		for {
			if s.closed {
				return
			}

			err = s.email.Start(emailmdl.MailFastKey)
			if err == redis.ErrNil {
				log.Info("emailFastProc start to rest")
				break
			}
			time.Sleep(200 * time.Millisecond)
		}

	}
}
