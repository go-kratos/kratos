package member

import (
	"context"

	"go-common/app/interface/main/account/model"
	accmdl "go-common/app/service/main/account/model"
	arcmdl "go-common/app/service/main/archive/model/archive"
	memmdl "go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_maxMonthlyOfficialSubmitTimes = 3
)

// IdentifyInfo get user identify info.
func (s *Service) IdentifyInfo(c context.Context, mid int64, ip string) (res *model.Identification, err error) {
	var rid *model.IdentifyInfo
	if rid, err = s.accDao.IdentifyInfo(c, mid, ip); err != nil {
		log.Error("s.memRPC.IdentifyInfo(%d) err(%+v)", mid, err)
		return
	}
	res = &model.Identification{}
	switch rid.Identify {
	case model.APIIdentifyOk:
		res.Identification = model.IdentifyOK
	case model.APIIdentifyNoInfo:
		res.Identification = model.IdentifyNotOK
	default:
		log.Error("unknow mid(%d) identify(%d) status", mid, rid.Identify)
	}
	return
}

// SubmitOfficial is.
func (s *Service) SubmitOfficial(c context.Context, mid int64, apply *model.OfficialApply) error {
	//ip := metadata.String(c, metadata.RemoteIP)
	if apply.Role == memmdl.OfficialRoleUp || apply.Role == memmdl.OfficialRoleIdentify {
		cons, err := s.OfficialConditions(c, mid)
		if err != nil {
			return err
		}
		if !cons.AllPass() {
			log.Warn("Unexpected official apply submited: mid: %d conditons: %+v apply: %+v", mid, cons, apply)
			return ecode.RequestErr
		}
	}

	// 是否超出本月提交次数限制
	times, err := s.accDao.GetMonthlyOfficialSubmittedTimes(c, mid)
	if err != nil {
		log.Error("Faield to get monthly official submitted times with mid %d: %+v", mid, err)
	}
	if times >= _maxMonthlyOfficialSubmitTimes {
		log.Warn("User %d is exceed max monthly official submitted times")
		return ecode.LimitExceed
	}

	ood, err := s.memRPC.OfficialDoc(c, &memmdl.ArgMid{Mid: mid})
	// 是否已经存在审核中的申请
	if err == nil && ood != nil && ood.State == memmdl.OfficialStateWait {
		return nil
	}

	if apply.Telephone != "" {
		if apply.TelVerifyCode == 0 {
			log.Error("Invalid tel verify code: mid: %d code:%d", mid, apply.TelVerifyCode)
			return ecode.RequestErr
		}
		vcode, verr := s.accDao.GetVerifyCode(c, mid, apply.Telephone)
		if verr != nil {
			log.Error("Failed to get verify code: %d, %s: %+v", mid, apply.Telephone, verr)
			return ecode.CaptchaErr
		}
		if apply.TelVerifyCode != vcode {
			log.Error("Failed to verify telephone verification code: %s, %d, %d", apply.Telephone, apply.TelVerifyCode, vcode)
			return ecode.CaptchaErr
		}
	}

	arg := &memmdl.ArgOfficialDoc{
		Mid:   mid,
		Name:  apply.Name,
		Role:  apply.Role,
		Title: apply.Title,
		Desc:  apply.Desc,

		Operator:          apply.Operator,
		Telephone:         apply.Telephone,
		Email:             apply.Email,
		Address:           apply.Address,
		Company:           apply.Company,
		CreditCode:        apply.CreditCode,
		Organization:      apply.Organization,
		OrganizationType:  apply.OrganizationType,
		BusinessLicense:   apply.BusinessLicense,
		BusinessScale:     apply.BusinessScale,
		BusinessLevel:     apply.BusinessLevel,
		BusinessAuth:      apply.BusinessAuth,
		Supplement:        apply.Supplement,
		Professional:      apply.Professional,
		Identification:    apply.Identification,
		OfficialSite:      apply.OfficialSite,
		RegisteredCapital: apply.RegisteredCapital,

		SubmitSource: "user", // 来自 account-interface 的全部为 user
	}

	pros, err := s.accRPC.ProfileWithStat3(c, &accmdl.ArgMid{Mid: mid})
	if err != nil {
		log.Error("Failed to call ProfileWithStat3(%d): %+v", mid, err)
		return err
	}
	arg.Realname = int8(pros.Identification)

	if err := s.accDao.DelVerifyCode(c, mid, apply.Telephone); err != nil {
		log.Error("Failed to delete verify code: mid: %d: mobile: %s: %+v", mid, apply.Telephone, err)
	}
	if _, err = s.accDao.IncreaseMonthlyOfficialSubmittedTimes(c, mid); err != nil {
		log.Error("Failed to increase monthly official submitted times with mid: %d: %+v", mid, err)
	}
	return s.memRPC.SetOfficialDoc(c, arg)
}

// OfficialConditions is.
func (s *Service) OfficialConditions(c context.Context, mid int64) (*model.OfficialConditions, error) {
	con := new(model.OfficialConditions)
	pros, err := s.accRPC.ProfileWithStat3(c, &accmdl.ArgMid{Mid: mid})
	if err != nil {
		log.Error("Failed to call ProfileWithStat3(%d): %+v", mid, err)
		return nil, err
	}
	if pros.Rank >= 10000 {
		con.IsFormal = true
	}
	// 1 正常号码，2 虚拟号码
	if pros.TelStatus >= 1 {
		con.BindTel = true
	}
	if pros.Identification == 1 {
		con.Realname = true
	}
	if pros.Follower >= 100000 {
		con.FollowerCount = true
	}

	arcCount, err := s.arcRPC.UpCount2(c, &arcmdl.ArgUpCount2{Mid: mid})
	if err != nil {
		log.Error("Failed to call s.arcRPC.UpCount2(%d): %+v", mid, err)
		// return nil, err
	}
	if err == nil && arcCount >= 1 {
		con.ArchiveCount = true
	}

	// 累计播放数
	// upStat, err := s.upRPC.UpStatBase(c, &upmdl.ArgMidWithDate{Mid: mid})
	// if err != nil {
	// 	log.Error("Failed to call s.upRPC.UpStatBase(%d): %+v", mid, err)
	// }
	// if err == nil && upStat != nil && upStat.View >= 1000000 {
	// 	con.ViewCount = true
	// }

	return con, nil
}

// UploadImage article upload cover.
func (s *Service) UploadImage(c context.Context, fileType string, body []byte) (url string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		return
	}
	if len(body) > s.c.BFS.MaxFileSize {
		err = ecode.FileTooLarge
		return
	}
	url, err = s.accDao.UploadImage(c, fileType, body, s.c.BFS)
	if err != nil {
		log.Error("account-interface: s.bfs.Upload error(%v)", err)
		return
	}
	return
}

// MobileVerify is.
func (s *Service) MobileVerify(c context.Context, mid int64, mobile string, country int64) error {
	ip := metadata.String(c, metadata.RemoteIP)
	vcode, err := s.accDao.GenVerifyCode(c, mid, mobile)
	if err != nil {
		log.Error("Failed to generate verify code: %+v", err)
		return err
	}
	return s.accDao.SendMobileVerify(c, vcode, country, mobile, ip)
}

// OfficialDoc is.
func (s *Service) OfficialDoc(c context.Context, mid int64) (*memmdl.OfficialDoc, error) {
	ip := metadata.String(c, metadata.RemoteIP)
	od, err := s.memRPC.OfficialDoc(c, &memmdl.ArgMid{Mid: mid, RealIP: ip})
	if err != nil {
		return nil, err
	}
	return od, nil
}

// MonthlyOfficialSubmittedTimes is
func (s *Service) MonthlyOfficialSubmittedTimes(c context.Context, mid int64) *model.OfficialSubmittedTimes {
	result := &model.OfficialSubmittedTimes{
		Submitted: 0,
		Remain:    _maxMonthlyOfficialSubmitTimes,
	}
	times, err := s.accDao.GetMonthlyOfficialSubmittedTimes(c, mid)
	if err != nil {
		log.Warn("Failed to get monthly official submitted times with mid: %d: %+v", mid, err)
		return result
	}
	result.Submitted = times
	if result.Submitted > _maxMonthlyOfficialSubmitTimes {
		result.Submitted = _maxMonthlyOfficialSubmitTimes
	}
	result.Remain = _maxMonthlyOfficialSubmitTimes - result.Submitted
	return result
}

// OfficialAutoFillDoc is
func (s *Service) OfficialAutoFillDoc(ctx context.Context, mid int64) (*memmdl.OfficialDoc, error) {
	res := &memmdl.OfficialDoc{
		Mid: mid,
	}
	// default name
	info, err := s.accRPC.Info3(ctx, &accmdl.ArgMid{Mid: mid})
	if err != nil {
		return nil, err
	}
	res.Name = info.Name
	// default from cm api
	func() {
		cminfo, err := s.accDao.BusinessAccountInfo(ctx, mid)
		if err != nil {
			log.Error("Failed to get cm business account info with mid: %d: %+v", mid, err)
			return
		}
		if cminfo.Nickname != "" {
			res.Name = cminfo.Nickname
		}
		if cminfo.CertificationTitle != "" {
			res.Title = cminfo.CertificationTitle
		}
		if cminfo.CreditCode != "" {
			res.CreditCode = cminfo.CreditCode
		}
		if cminfo.CompanyName != "" {
			res.Company = cminfo.CompanyName
		}
		if cminfo.Organization != "" {
			res.Organization = cminfo.Organization
		}
		if cminfo.OrganizationType != "" {
			res.OrganizationType = cminfo.OrganizationType
		}
	}()
	return res, nil
}
