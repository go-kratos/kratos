package service

import (
	"context"
	"strings"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// CompareInfo compare user_info with sys_info
func (s *Service) CompareInfo(c context.Context, rid int64) (err error) {
	var (
		games             []*model.Game
		sysiInfo          = &model.SysInfo{}
		emailphone        = &model.CheckEmailPhone{}
		sysRegCheck       string
		sysUnamesCheck    string
		sysSafeCheck      string
		sysCardCheck      string
		sysPwdsCheck      string
		sysLoginAddrCheck string
		userType          int64
	)
	uinfo, err := s.d.GetUnCheckInfo(c, rid)
	if err != nil || uinfo == nil {
		return
	}
	if err == nil && uinfo.Mid == 0 {
		return
	}
	eg, _ := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		games, e = s.d.GetUserType(c, uinfo.Mid)
		return
	})
	eg.Go(func() (e error) {
		//前端只能以"台湾地区"展示给用户，所以在这里做一下转换
		if uinfo.RegAddr == "台湾地区" {
			uinfo.RegAddr = "台湾"
		}
		sysRegCheck, e = s.checkRegInfo(c, uinfo)
		return
	})
	eg.Go(func() (e error) {
		sysUnamesCheck, e = s.checkUnames(c, uinfo)
		return
	})
	eg.Go(func() (e error) {
		phoneAddr := strings.Split(uinfo.Phones, ";")
		//手机去掉第一个0，进行对比
		phones := strings.Split(string(phoneAddr[0]), ",")
		phoneDeal := ""
		for _, p := range phones {
			if strings.Index(p, "0") == 0 && len(p) > 1 {
				p = p[1:]
			}
			phoneDeal += "," + p
		}
		uinfo.Phones = strings.TrimLeft(phoneDeal, ",")
		emailphone, e = s.checkPhonesAndEmails(c, uinfo)
		return
	})
	eg.Go(func() (e error) {
		sysSafeCheck, e = s.checkSafe(c, uinfo)
		return
	})
	eg.Go(func() (e error) {
		sysCardCheck, e = s.checkCard(c, uinfo)
		return
	})
	eg.Go(func() (e error) {
		sysPwdsCheck, e = s.checkPwds(c, uinfo)
		return
	})
	eg.Go(func() (e error) {
		sysLoginAddrCheck, e = s.getLoginAddrs(c, uinfo) //调用java得到ip再去通过rpc调用得到地址
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("CompareInfo err(%v)", err)
		return
	}

	if len(games) == 0 {
		userType = 0
	} else {
		userType = 1
	}
	sysiInfo.SysReg = sysRegCheck
	sysiInfo.SysUNames = sysUnamesCheck
	sysiInfo.SysPhones = emailphone.PhonesCheck
	sysiInfo.SysEmails = emailphone.EmailCheck
	sysiInfo.SysSafe = sysSafeCheck
	sysiInfo.SysCard = sysCardCheck
	sysiInfo.SysPwds = sysPwdsCheck
	sysiInfo.SysLoginAddrs = sysLoginAddrCheck

	err = s.d.UpdateSysInfo(c, sysiInfo, userType, rid)
	if err != nil {
		return
	}
	s.AddMailch(func() {
		s.deal(context.Background(), rid)
	})
	return
}

// checkRegInfo 注册信息校验
func (s *Service) checkRegInfo(c context.Context, uinfo *model.UserInfoReq) (sysRegCheck string, err error) {
	reg, err := s.d.CheckReg(c, uinfo.Mid, uinfo.RegTime.Time().Unix(), uinfo.RegType, uinfo.RegAddr)
	if err != nil {
		log.Error("checkRegInfo error(%v)", err)
		return
	}
	log.Info("checkRegInfo req uinfo:%+v, reg: %+v", uinfo, reg)
	sysRegCheck = reg.CheckInfo
	return
}

// checkUnames 昵称列表信息校验
func (s *Service) checkUnames(c context.Context, uinfo *model.UserInfoReq) (sysUnamesCheck string, err error) {
	var req = &model.NickNameReq{Mid: uinfo.Mid, Size: 100}
	var res *model.NickNameLogRes
	res, err = s.d.NickNameLog(c, req)
	if err != nil {
		return
	}
	info, err1 := s.d.Info3(c, uinfo.Mid)
	if err1 != nil {
		return "", err1
	}

	names := make(map[string]bool)
	if info != nil {
		names[info.Name] = true
	}
	if len(res.Result) != 0 {
		for _, r := range res.Result {
			names[r.OldName] = true
			names[r.NewName] = true
		}
	}
	if uinfo.Unames == "" {
		sysUnamesCheck = "0"
		if len(names) == 0 {
			sysUnamesCheck = "1"
		}
		return
	}
	unames := strings.Split(uinfo.Unames, ",")
	unameCheck := make([]string, len(unames))
	for i, uname := range unames {
		unameCheck[i] = "0"
		if names[uname] {
			unameCheck[i] = "1"
		}
	}
	log.Info("checkUnames req uinfo:%+v, NickNameLog: %+v, Info: %+v, sysUnamesCheck:%+v", uinfo, res, info, sysUnamesCheck)
	sysUnamesCheck = strings.Join(unameCheck, ",")
	return
}

// checkPhonesAndEmails 手机and邮箱列表信息校验
func (s *Service) checkPhonesAndEmails(c context.Context, uinfo *model.UserInfoReq) (result *model.CheckEmailPhone, err error) {
	var resPhones = &model.UserBindLogRes{}
	var resEmails = &model.UserBindLogRes{}

	phoneReq := &model.UserBindLogReq{
		Action: "telBindLog",
		Mid:    uinfo.Mid,
		Size:   100,
	}
	if resPhones, err = s.d.UserBindLog(c, phoneReq); err != nil {
		return
	}
	emailReq := &model.UserBindLogReq{
		Action: "emailBindLog",
		Mid:    uinfo.Mid,
		Size:   100,
	}
	if resEmails, err = s.d.UserBindLog(c, emailReq); err != nil {
		return
	}

	var userInfo *model.UserInfo
	if userInfo, err = s.d.GetUserInfo(c, uinfo.Mid); err != nil {
		return
	}

	tels := make(map[string]bool)
	emails := make(map[string]bool)
	if userInfo.Phone != "" {
		tels[userInfo.Phone] = true
	}
	if userInfo.Email != "" {
		emails[userInfo.Email] = true
	}
	for _, item := range resPhones.Result {
		tels[item.Phone] = true
	}
	for _, item := range resEmails.Result {
		emails[item.Email] = true
	}
	result = new(model.CheckEmailPhone)
	result.PhonesCheck = checkTel(uinfo.Phones, tels)
	result.EmailCheck = checkEmail(uinfo.Emails, emails)
	log.Info("checkPhonesAndEmails req uinfo:%+v, resPhones: %+v, resEmails: %+v, userInfo: %+v,result: %+v", uinfo, resPhones, resEmails, userInfo, result)
	return
}

func checkEmail(ues string, emails map[string]bool) (result string) {
	if ues == "" {
		result = "0"
		if len(emails) == 0 {
			result = "1"
		}
		return
	}
	es := strings.Split(ues, ",")
	emailCheck := make([]string, len(es))
	for i, e := range es {
		emailCheck[i] = "0"
		if emails[e] {
			emailCheck[i] = "1"
		}
	}
	result = strings.Join(emailCheck, ",")
	return
}

func checkTel(phones string, tels map[string]bool) (result string) {
	if phones == "" {
		result = "0"
		if len(tels) == 0 {
			result = "1"
		}
		return
	}
	ps := strings.Split(phones, ",")
	phoneCheck := make([]string, len(ps))
	for i, phone := range ps {
		phoneCheck[i] = "0"
		if tels[phone] {
			phoneCheck[i] = "1"
		}
	}
	result = strings.Join(phoneCheck, ",")
	return
}

// checkSafe  密保信息校验
func (s *Service) checkSafe(c context.Context, uinfo *model.UserInfoReq) (sysSafeCheck string, err error) {
	//没有绑定密保
	if uinfo.SafeQuestion == 99 {
		uinfo.SafeQuestion = 0
	}
	safe, err := s.d.CheckSafe(c, uinfo.Mid, uinfo.SafeQuestion, uinfo.SafeAnswer)
	if err != nil {
		log.Error("dao method CheckSafe error(%v)", err)
	}
	sysSafeCheck = safe.CheckInfo
	return
}

// checkCard  校验证件信息
func (s *Service) checkCard(c context.Context, uinfo *model.UserInfoReq) (sysCardCheck string, err error) {
	sysCardCheck = "0"
	//（1）如果，用户在系统里，账号本身没有设置没记录，则不论用户找回时填写如否，系统资料对比列，显示记为 null。-->请求status查看是否认证
	//（2）如果，用户在系统里，账号有设置有记录，则用户在找回时，填写错误或没填写，则比对标红记录 错；填写比对正确，则为 对。-->请求check，查看是否输入正确
	//me: 调用接口失败，统一判断为错误
	status, err := s.d.CheckRealnameStatus(c, uinfo.Mid)
	if err != nil {
		return
	}

	//只要系统没有记录认证信息，均返回null
	if status == 0 {
		return "null", nil
	}

	//用户没有填写信息，但是系统有记录认证信息
	if uinfo.CardID == "" || uinfo.CardType == 99 {
		return "0", nil
	}

	//用户填写了信息，且已经认证，检查认证信息
	flag, err := s.d.CheckCard(c, uinfo.Mid, uinfo.CardType, uinfo.CardID)
	if err != nil {
		return
	}
	if flag {
		sysCardCheck = "1"
	}
	return
}

// checkPwds  校验历史密码信息
func (s *Service) checkPwds(c context.Context, uinfo *model.UserInfoReq) (sysPwdsCheck string, err error) {
	if sysPwdsCheck, err = s.d.CheckPwds(c, uinfo.Mid, uinfo.Pwds); err != nil {
		return
	}
	return
}

// getLoginAddrs 获取登陆地
func (s *Service) getLoginAddrs(c context.Context, uinfo *model.UserInfoReq) (addrs string, err error) {
	addrs, err = s.d.GetAddrByIP(c, uinfo.Mid, 100)
	if err != nil {
		log.Error("GetAddrByIP error(%v)", err)
		return
	}
	return
}
