package member

import (
	"context"
	"fmt"
	"net/url"
	"time"
	"unicode/utf8"

	"go-common/app/interface/main/account/model"
	accModel "go-common/app/service/main/account/model"
	coModel "go-common/app/service/main/coin/model"
	ftModel "go-common/app/service/main/filter/model/rpc"
	locModel "go-common/app/service/main/location/model"
	meModel "go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"

	"github.com/pkg/errors"
)

var monitoredIPCountry = map[string]string{
	"台湾":  "台湾 IP",
	"香港":  "香港 IP",
	"美国":  "美国 IP",
	"加拿大": "加拿大 IP",
}

var monitoredTelCountry = map[int64]string{
	1: "美国/加拿大手机号",
}

var upNameCostCoins = 6.0

// Account get Account info.
func (s *Service) Account(c context.Context, mid int64, ip string) (acc *model.Account, err error) {
	var (
		nickFree *model.NickFree
		mb       *meModel.Member
	)
	marg := &meModel.ArgMemberMid{Mid: mid, RemoteIP: ip}
	if mb, err = s.memRPC.Member(c, marg); err != nil {
		log.Error("service.memberRPC.MyInfo(%v) error(%v)", marg, err)
		return
	}
	if nickFree, err = s.NickFree(c, mid); err != nil {
		return
	}
	acc = &model.Account{}
	acc.Mid = mid
	acc.Birthday = mb.Birthday.Time().Format("2006-01-02")
	acc.Uname = mb.Name
	acc.Face = mb.Face
	acc.Sign = mb.Sign
	acc.Sex = int8(mb.Sex)
	acc.NickFree = nickFree.NickFree
	return
}

// UpdateFace Update Face
func (s *Service) UpdateFace(c context.Context, mid int64, faceFile []byte, ftype string) (string, error) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check
	profile, err := s.accRPC.Profile3(c, &accModel.ArgMid{Mid: mid})
	if err != nil {
		return "", errors.WithStack(err)
	}
	//判断是否绑定手机号
	if !s.validateTelStatus(profile.TelStatus) {
		return "", ecode.MemberPhoneRequired
	}
	if profile.Silence != 0 {
		return "", ecode.MemberBlocked
	}
	//Upload bfs
	faceURL, err := s.accDao.UploadImage(c, ftype, faceFile, s.c.FaceBFS)
	if err != nil {
		log.Error("s.bfsDao.Upload(%d) error(%v)", mid, err)
		return "", errors.WithStack(err)
	}
	URL, err := url.Parse(faceURL)
	if err != nil {
		return "", errors.WithStack(err)
	}
	inMonitor := s.ensureMonitor(c, mid, ip)
	arg := &meModel.ArgAddPropertyReview{
		Mid:      mid,
		New:      URL.Path,
		State:    meModel.ReviewStateQueuing,
		Property: meModel.ReviewPropertyFace,
	}
	if inMonitor {
		arg.State = meModel.ReviewStateWait
		return profile.Face, s.memRPC.AddPropertyReview(c, arg)
	}
	if err := s.memRPC.AddPropertyReview(c, arg); err != nil {
		log.Error("s.memRPC.AddPropertyReview(%d,%s) error(%v)", mid, faceFile, err)
		return "", errors.WithStack(err)
	}
	if err := s.memRPC.SetFace(c, &meModel.ArgUpdateFace{Mid: mid, Face: URL.Path}); err != nil {
		log.Error("s.memRPC.SetFace(%d,%s) error(%v)", mid, faceURL, err)
		return "", errors.WithStack(err)
	}
	return faceURL, nil
}

// UpdateName .
func (s *Service) UpdateName(c context.Context, mid int64, name, appkey string) error {
	ip := metadata.String(c, metadata.RemoteIP)
	_, inWhiteList := s.nickFreeAppKeys[appkey]
	if inWhiteList {
		return s.updateNameWithinWhiteList(c, mid, name, ip)
	}
	return s.updateName(c, mid, name, ip)
}

// updateNameWithinWhiteList 白名单 appkey 不扣硬币
func (s *Service) updateNameWithinWhiteList(c context.Context, mid int64, name, ip string) error {
	if err := s.nameIsValid(c, mid, name, ip); err != nil {
		return err
	}
	profile, err := s.accRPC.Profile3(c, &accModel.ArgMid{Mid: mid})
	if err != nil {
		return errors.WithStack(err)
	}
	if err := s.permitName(c, profile, ip); err != nil {
		return err
	}
	if profile.Name == name {
		log.Info("Update name is same to origin: mid: %d, name: %s, origin: %s", mid, name, profile.Name)
		return nil
	}
	inMonitor := s.ensureMonitor(c, mid, ip)
	remark := "appkey白名单修改昵称"
	// 在监控列表里就加入添加审核列表
	if inMonitor {
		saveUpNameLog(mid, profile.Name, name, remark, inMonitor, ip)
		return errors.WithStack(s.memRPC.AddPropertyReview(c, &meModel.ArgAddPropertyReview{
			Mid:      mid,
			New:      name,
			State:    meModel.ReviewStateWait,
			Property: meModel.ReviewPropertyName,
			Extra:    map[string]interface{}{"nick_free": true},
		}))
	}
	//修改昵称
	if err := s.passDao.UpdateName(c, mid, name, ip); err != nil {
		return errors.WithStack(err)
	}
	saveUpNameLog(mid, profile.Name, name, remark, inMonitor, ip)
	return nil
}

//UpdateName update name.
func (s *Service) updateName(c context.Context, mid int64, name, ip string) error {
	if err := s.nameIsValid(c, mid, name, ip); err != nil {
		return err
	}
	profile, err := s.accRPC.Profile3(c, &accModel.ArgMid{Mid: mid})
	if err != nil {
		return errors.WithStack(err)
	}
	if err = s.permitName(c, profile, ip); err != nil {
		return err
	}
	if profile.Name == name {
		log.Info("Update name is same to origin: mid: %d, name: %s, origin: %s", mid, name, profile.Name)
		return nil
	}
	// 判断是否改昵称免费
	nickFree, err := s.NickFree(c, mid)
	if err != nil {
		return err
	}
	remark := "快速注册修改昵称"
	if !nickFree.NickFree {
		coins, coinErr := s.coinRPC.UserCoins(c, &coModel.ArgCoinInfo{Mid: mid, RealIP: ip})
		if coinErr != nil {
			return errors.WithStack(coinErr)
		}
		if coins < upNameCostCoins {
			return ecode.UpdateUnameMoneyIsNot
		}
		remark = "修改昵称"
	}
	inMonitor := s.ensureMonitor(c, mid, ip)
	// 在监控列表里就加入添加审核列表
	if inMonitor {
		saveUpNameLog(mid, profile.Name, name, remark, inMonitor, ip)
		return errors.WithStack(s.memRPC.AddPropertyReview(c, &meModel.ArgAddPropertyReview{
			Mid:      mid,
			New:      name,
			State:    meModel.ReviewStateWait,
			Property: meModel.ReviewPropertyName,
			Extra:    map[string]interface{}{"nick_free": nickFree.NickFree},
		}))
	}
	//修改昵称
	if err = s.passDao.UpdateName(c, mid, name, ip); err != nil {
		return errors.WithStack(err)
	}
	saveUpNameLog(mid, profile.Name, name, remark, inMonitor, ip)
	if nickFree.NickFree {
		return errors.WithStack(s.memRPC.SetNickUpdated(c, &meModel.ArgMemberMid{Mid: mid}))
	}
	//扣除硬币
	if _, err = s.coinRPC.ModifyCoin(c, &coModel.ArgModifyCoin{
		Mid:    mid,
		Count:  -upNameCostCoins,
		Reason: fmt.Sprintf("UPDATE:NICK:%s=>%s", profile.Name, name),
		IP:     ip,
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Service) monitorByIP(ctx context.Context, mid int64, ip string) (bool, string) {
	IP, err := s.locRPC.Info(ctx, &locModel.ArgIP{IP: ip})
	if err != nil || IP == nil {
		log.Error("Failed to get ip info with ip: %s: %+v", ip, err)
		return false, ""
	}
	descr, shouldMonitor := monitoredIPCountry[IP.Country]
	if !shouldMonitor {
		return false, ""
	}
	return true, descr
}

func (s *Service) monitorByTel(ctx context.Context, mid int64, ip string) (bool, string) {
	p, err := s.passDao.QueryByMid(ctx, mid, ip)
	if err != nil {
		log.Error("Failed to query by mid form pasport: mid: %d: %+v", mid, err)
		return false, ""
	}
	descr, shouldMonitor := monitoredTelCountry[p.CountryCode]
	if !shouldMonitor {
		return false, ""
	}
	return true, descr
}

func (s *Service) shouldMonitor(ctx context.Context, mid int64, ip string) (bool, string) {
	should, descr := s.monitorByIP(ctx, mid, ip)
	if should {
		return true, descr
	}

	should, descr = s.monitorByTel(ctx, mid, ip)
	if should {
		return true, descr
	}

	return false, ""
}

func (s *Service) ensureMonitor(ctx context.Context, mid int64, ip string) bool {
	inMonitor, _ := s.memRPC.IsInMonitor(ctx, &meModel.ArgMid{Mid: mid})
	if inMonitor {
		return true
	}

	should, descr := s.shouldMonitor(ctx, mid, ip)
	if !should {
		return false
	}
	if err := s.memRPC.AddUserMonitor(ctx, &meModel.ArgAddUserMonitor{
		Mid:      mid,
		Operator: "system",
		Remark:   fmt.Sprintf("系统自动导入-%s", descr),
	}); err != nil {
		log.Error("Failed to add user moniter: mid: %d: %+v", mid, err)
	}
	return true
}

// UpdateSex update sex.
func (s *Service) UpdateSex(c context.Context, mid, sex int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	return s.accDao.UpdateSex(c, mid, sex, ip)
}

//UpdateSign update sign.
func (s *Service) UpdateSign(c context.Context, mid int64, sign string) error {
	ip := metadata.String(c, metadata.RemoteIP)
	// 签名最长 70 个字符
	if utf8.RuneCountInString(sign) > 70 {
		return ecode.MemberSignOverLimit
	}
	// 签名不能包含 emoji
	if model.HasEmoji(sign) {
		return ecode.MemberSignHasEmoji
	}

	// 过滤敏感词
	res, err := s.filterRPC.Filter(c, &ftModel.ArgFilter{Area: "sign", Message: sign})
	if err != nil {
		return err
	}
	// 大于 20 认为包含敏感词
	if res.Level >= 20 {
		return ecode.MemberSignSensitive
	}

	// 检查是否绑定手机
	profile, err := s.accRPC.Profile3(c, &accModel.ArgMid{Mid: mid})
	if err != nil {
		return errors.WithStack(err)
	}
	if !s.validateTelStatus(profile.TelStatus) {
		return ecode.MemberPhoneRequired
	}
	// 检查是否被禁言
	if profile.Silence != 0 {
		return ecode.MemberBlocked
	}
	// 如果和老的一模一样就没必要更新了
	if profile.Sign == sign {
		log.Info("Update sign is same to origin: mid: %d, sign: %s, origin: %s", mid, sign, profile.Sign)
		return nil
	}

	inMonitor := s.ensureMonitor(c, mid, ip)
	// 不在监控列表里就直接更新
	if !inMonitor {
		return errors.WithStack(s.memRPC.SetSign(c, &meModel.ArgUpdateSign{
			Mid:      mid,
			Sign:     sign,
			RemoteIP: ip,
		}))
	}
	// 否则就加入监控列表
	return errors.WithStack(s.memRPC.AddPropertyReview(c, &meModel.ArgAddPropertyReview{
		Mid:      mid,
		New:      sign,
		State:    meModel.ReviewStateWait,
		Property: meModel.ReviewPropertySign,
	}))
}

// UpdateBirthday update birthday.
func (s *Service) UpdateBirthday(c context.Context, mid int64, birthday string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	return s.accDao.UpdateBirthday(c, mid, ip, birthday)
}

// NickFree .
func (s *Service) NickFree(c context.Context, mid int64) (nickFree *model.NickFree, err error) {
	var (
		isRegFast   bool
		nickUpdated bool
		ip          = metadata.String(c, metadata.RemoteIP)
	)
	sarg := &meModel.ArgMemberMid{Mid: mid}
	if nickUpdated, err = s.memRPC.NickUpdated(c, sarg); err != nil {
		log.Error("s.memRPC.IsUpNickFree(%v) error (%v)", sarg, err)
		return
	}
	nickFree = &model.NickFree{}
	if nickUpdated {
		return
	}
	if isRegFast, err = s.passDao.FastReg(c, mid, ip); err != nil {
		return
	}
	if isRegFast {
		nickFree.NickFree = true
	}
	return
}

func saveUpNameLog(mid int64, oName, nName, remark string, isMonitor bool, ip string) {
	report.User(&report.UserInfo{
		Mid:      mid,
		Business: model.UpNameLogID,
		Action:   model.UpNameAction,
		IP:       ip,
		Ctime:    time.Now(),
		Index:    []interface{}{0, 0, 0, oName, nName, remark},
		Content: map[string]interface{}{
			"is_monitor": isMonitor,
			"old_name":   oName,
			"new_name":   nName,
			"reason":     fmt.Sprintf("修改昵称（原昵称：%s 新昵称：%s)", oName, nName),
			"remark":     remark,
		},
	})
}

func (s *Service) permitName(c context.Context, profile *accModel.Profile, ip string) error {
	if !s.validateTelStatus(profile.TelStatus) {
		return ecode.MemberPhoneRequired
	}
	// 检查是否被禁言
	if profile.Silence != 0 {
		return ecode.MemberBlocked
	}
	//昵称锁定,是否官方认证
	if profile.Official.Role != 0 {
		log.Info("update name fail, name is official, mid: %d", profile.Mid)
		return ecode.UpdateUnameHadOfficial
	}
	pProfile, err := s.passDao.QueryByMid(c, profile.Mid, ip)
	if err != nil {
		return err
	}
	if pProfile.NickLock == 1 {
		log.Info("update name fail, name is locked, mid: %d", profile.Mid)
		return ecode.UpdateUnameHadLocked
	}
	return nil
}

func (s *Service) nameIsValid(c context.Context, mid int64, name, ip string) error {
	if len(name) > 30 || utf8.RuneCountInString(name) > 16 {
		return ecode.UpdateUnameTooLong
	}
	if utf8.RuneCountInString(name) < 3 {
		return ecode.UpdateUnameTooShort
	}
	if !model.ValidName(name) {
		return ecode.UpdateUnameFormat
	}
	// 判断昵称是否重复
	if err := s.passDao.TestUserName(c, name, mid, ip); err != nil {
		return err
	}
	// 过滤敏感词
	res, err := s.filterRPC.Filter(c, &ftModel.ArgFilter{Area: "member", Message: name})
	if err != nil {
		return err
	}
	// 大于 20 认为包含敏感词
	if res.Level >= 20 {
		return ecode.UpdateUnameSensitive
	}
	return nil
}

func (s *Service) validateTelStatus(status int32) bool {
	if s.c.Switch.UpdatePropertyPhoneRequired && status == 0 {
		return false
	}
	return true
}
