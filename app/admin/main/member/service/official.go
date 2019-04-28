package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/app/service/main/member/model/block"
	spymodel "go-common/app/service/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_logActionAudit    = "official_doc_audit"
	_logActionEdit     = "official_doc_edit"
	_logActionEditName = "official_doc_edit_name"
)

func i64Toi8(in []int64) []int8 {
	out := make([]int8, 0, len(in))
	for _, i := range in {
		out = append(out, int8(i))
	}
	return out
}

func i8Toi64(in []int8) []int64 {
	out := make([]int64, 0, len(in))
	for _, i := range in {
		out = append(out, int64(i))
	}
	return out
}

func actOr(act ...string) string {
	return strings.Join(act, ",")
}

func (s *Service) officialName(ctx context.Context, ofs []*model.Official) {
	mids := make([]int64, 0, len(ofs))
	for _, o := range ofs {
		mids = append(mids, o.Mid)
	}
	ofds, err := s.dao.OfficialDocsByMids(ctx, mids)
	if err != nil {
		log.Error("Failed to s.dao.OfficialDocsByMids(%+v): %+v", mids, err)
		return
	}
	for _, o := range ofs {
		od, ok := ofds[o.Mid]
		if !ok {
			continue
		}
		o.Name = od.Name
	}
}

// Officials is.
func (s *Service) Officials(ctx context.Context, arg *model.ArgOfficial) ([]*model.Official, int, error) {
	if len(arg.Role) == 0 {
		arg.Role = i8Toi64(model.AllRoles)
	}
	if arg.ETime == 0 {
		arg.ETime = xtime.Time(time.Now().Unix())
	}
	ofs, total, err := s.dao.Officials(ctx, arg.Mid, i64Toi8(arg.Role), arg.STime.Time(), arg.ETime.Time(), arg.Pn, arg.Ps)
	if err != nil {
		return nil, 0, err
	}
	// 需要展示昵称
	s.officialName(ctx, ofs)
	return ofs, total, err
}

func (s *Service) blockResult(ctx context.Context, mid int64) (*model.BlockResult, error) {
	info, err := s.memberRPC.BlockInfo(ctx, &block.RPCArgInfo{MID: mid})
	if err != nil {
		err = errors.Wrapf(err, "%v", mid)
		return nil, err
	}
	block := &model.BlockResult{
		MID:         info.MID,
		BlockStatus: info.BlockStatus,
		StartTime:   info.StartTime,
		EndTime:     info.EndTime,
	}
	return block, nil
}

// OfficialDoc is.
func (s *Service) OfficialDoc(ctx context.Context, mid int64) (od *model.OfficialDoc, logs *model.SearchLogResult, block *model.BlockResult, spys []*spymodel.Statistics, realname *model.Realname, sameCreditCodeMids []int64, err error) {
	if od, err = s.dao.OfficialDoc(ctx, mid); err != nil {
		return
	}
	if od == nil {
		od = &model.OfficialDoc{
			Mid:           mid,
			OfficialExtra: &model.OfficialExtra{},
		}
	}
	logs, err = s.dao.SearchLog(ctx, 0, mid, "", actOr(_logActionAudit, _logActionEdit))
	if err != nil {
		log.Error("Failed to s.dao.SearchLog(%+v): %+v", mid, err)
		return
	}
	block, err = s.blockResult(ctx, mid)
	if err != nil {
		log.Error("Failed to s.blockResult(%+v): %+v", mid, err)
		return
	}
	arg := &spymodel.ArgStat{Mid: mid}
	spys, err = s.spyRPC.StatByID(ctx, arg)
	if err != nil {
		log.Error("Failed to s.spyRPC.StatByID: mid(%d): %+v", od.Mid, err)
		return
	}

	realname, err = s.officialRealname(ctx, mid)
	if err != nil {
		log.Error("Failed to get official realname with mid: %d: %+v", od.Mid, err)
		return
	}

	// 查询使用相同社会信用代码的mid
	sameCreditCodeMids = make([]int64, 0)
	if od.OfficialExtra.CreditCode != "" {
		func() {
			addits, err := s.OfficialDocAddits(ctx, "credit_code", od.OfficialExtra.CreditCode)
			if err != nil {
				log.Error("Failed to get official addit with mid: %d: %+v", od.Mid, err)
				return
			}
			for _, addit := range addits {
				if addit.Mid != od.Mid {
					sameCreditCodeMids = append(sameCreditCodeMids, addit.Mid)
				}
			}
		}()
	}

	return
}

func (s *Service) officialRealname(ctx context.Context, mid int64) (*model.Realname, error) {
	realname := &model.Realname{
		State: model.RealnameApplyStateNone,
	}
	dr, err := s.dao.RealnameInfo(ctx, mid)
	if err != nil {
		log.Error("Failed to get realname info with mid: %d: %+v", mid, err)
		return realname, nil
	}
	if dr != nil {
		realname.ParseInfo(dr)
	}

	imagesByMain := func() {
		apply, err := s.dao.LastPassedRealnameMainApply(ctx, mid)
		if err != nil {
			log.Error("Failed to get last passed realname main apply with mid: %d: %+v", mid, err)
			return
		}

		images, err := s.dao.RealnameApplyIMG(ctx, []int64{apply.HandIMG, apply.FrontIMG, apply.BackIMG})
		if err != nil {
			log.Error("Failed to get realname apply image by apply: %+v: %+v", apply, err)
			return
		}

		for _, image := range images {
			realname.ParseDBApplyIMG(image.IMGData)
		}
	}

	imagesByAlipay := func() {
		apply, err := s.dao.LastPassedRealnameAlipayApply(ctx, mid)
		if err != nil {
			log.Error("Failed to get last passed realname alipay apply with mid: %d: %+v", mid, err)
			return
		}
		realname.ParseDBApplyIMG(apply.IMG)
	}

	switch dr.Channel {
	case model.ChannelMain.DBChannel():
		imagesByMain()
	case model.ChannelAlipay.DBChannel():
		imagesByAlipay()
	default:
		log.Error("Failed to get realname apply images by realname info: %+v", dr)
	}

	return realname, nil
}

// OfficialDocs is.
func (s *Service) OfficialDocs(ctx context.Context, arg *model.ArgOfficialDoc) ([]*model.OfficialDoc, int, error) {
	if len(arg.Role) == 0 {
		arg.Role = i8Toi64(model.AllRoles)
	}
	if len(arg.State) == 0 {
		arg.State = i8Toi64(model.AllStates)
	}
	if arg.ETime == 0 {
		arg.ETime = xtime.Time(time.Now().Unix())
	}
	return s.dao.OfficialDocs(ctx, arg.Mid, i64Toi8(arg.Role), i64Toi8(arg.State), arg.Uname, arg.STime.Time(), arg.ETime.Time(), arg.Pn, arg.Ps)
}

// OfficialDocAudit is.
func (s *Service) OfficialDocAudit(ctx context.Context, arg *model.ArgOfficialAudit) (err error) {
	od, err := s.dao.OfficialDoc(ctx, arg.Mid)
	if err != nil {
		return
	}
	if arg.State == model.OfficialStatePass {
		if err = s.updateUname(ctx, od.Mid, od.Name, arg.UID, arg.Uname); err != nil {
			log.Error("Failed to update uname: mid(%d), name(%s): %+v", od.Mid, od.Name, err)
			return
		}
	}

	if err = s.dao.OfficialDocAudit(ctx, arg.Mid, arg.State, arg.Uname, arg.IsInternal, arg.Reason); err != nil {
		return
	}
	od, err = s.dao.OfficialDoc(ctx, arg.Mid)
	if err != nil {
		return
	}
	role := int8(model.OfficialRoleUnauth)
	cnt := `对不起，您的官方认证申请未通过，未通过原因："` + arg.Reason + `，重新申请点#{这里}{"https://account.bilibili.com/account/official/home"}`
	if arg.State == model.OfficialStatePass {
		role = od.Role
		cnt = "恭喜，您的官方认证申请已经通过啦！o(*￣▽￣*)o"
	}
	if _, err = s.dao.OfficialEdit(ctx, arg.Mid, role, od.Title, od.Desc); err != nil {
		return
	}

	if err = s.dao.Message(ctx, "官方认证审核通知", cnt, []int64{arg.Mid}); err != nil {
		log.Error("Failed to send message: %+v", err)
		err = nil
	}
	report.Manager(&report.ManagerInfo{
		Uname:    arg.Uname,
		UID:      arg.UID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionAudit,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{arg.State, int64(od.CTime), od.Role, od.Name, od.Title, od.Desc},
		Content: map[string]interface{}{
			"reason":    arg.Reason,
			"name":      od.Name,
			"extra":     od.Extra,
			"role":      od.Role,
			"title":     od.Title,
			"desc":      od.Desc,
			"state":     arg.State,
			"doc_ctime": int64(od.CTime),
		},
	})
	return
}

// OfficialDocEdit is.
func (s *Service) OfficialDocEdit(ctx context.Context, arg *model.ArgOfficialEdit) (err error) {
	od, _ := s.dao.OfficialDoc(ctx, arg.Mid)
	if od == nil {
		od = &model.OfficialDoc{
			Mid:           arg.Mid,
			Role:          arg.Role,
			OfficialExtra: &model.OfficialExtra{},
		}
	}
	od.State = int8(model.OfficialStatePass)
	if arg.Role == model.OfficialRoleUnauth {
		od.State = int8(model.OfficialStateNoPass)
	}
	od.Name = arg.Name
	od.Uname = arg.Uname
	od.Telephone = arg.Telephone
	od.Email = arg.Email
	od.Address = arg.Address
	od.Supplement = arg.Supplement
	od.Company = arg.Company
	od.Operator = arg.Operator
	od.CreditCode = arg.CreditCode
	od.Organization = arg.Organization
	od.OrganizationType = arg.OrganizationType
	od.BusinessLicense = arg.BusinessLicense
	od.BusinessLevel = arg.BusinessLevel
	od.BusinessScale = arg.BusinessScale
	od.BusinessAuth = arg.BusinessAuth
	od.OfficalSite = arg.OfficalSite
	od.RegisteredCapital = arg.RegisteredCapital
	extra, err := json.Marshal(od.OfficialExtra)
	if err != nil {
		err = errors.Wrap(err, "official doc edit")
		return
	}
	if err = s.updateUname(ctx, arg.Mid, arg.Name, arg.UID, arg.Uname); err != nil {
		log.Error("Failed to update uname: mid(%d), name(%s): %+v", arg.Mid, arg.Name, err)
		err = ecode.MemberNameFormatErr
		return
	}
	if err = s.dao.OfficialDocEdit(ctx, arg.Mid, arg.Name, arg.Role, od.State, arg.Title, arg.Desc, string(extra), arg.Uname, arg.IsInternal); err != nil {
		log.Error("Failed to update official doc: %+v", err)
		err = ecode.RequestErr
		return
	}
	if _, err = s.dao.OfficialEdit(ctx, arg.Mid, arg.Role, arg.Title, arg.Desc); err != nil {
		return
	}
	report.Manager(&report.ManagerInfo{
		Uname:    arg.Uname,
		UID:      arg.UID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionEdit,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{od.State, int64(od.CTime), arg.Role, arg.Name, arg.Title, arg.Desc},
		Content: map[string]interface{}{
			"extra":     string(extra),
			"name":      arg.Name,
			"role":      arg.Role,
			"title":     arg.Title,
			"desc":      arg.Desc,
			"state":     od.State,
			"doc_ctime": int64(od.CTime),
		},
	})
	if arg.SendMessage {
		if merr := s.dao.Message(ctx, arg.MessageTitle, arg.MessageContent, []int64{arg.Mid}); merr != nil {
			log.Error("Failed to send message: %+v", merr)
		}
	}
	return
}

// OfficialDocSubmit is.
func (s *Service) OfficialDocSubmit(ctx context.Context, arg *model.ArgOfficialSubmit) (err error) {
	od := &model.OfficialDoc{
		Mid:   arg.Mid,
		Name:  arg.Name,
		State: int8(model.OfficialStateWait),
		Role:  arg.Role,
		Title: arg.Title,
		Desc:  arg.Desc,

		Uname:        arg.Uname,
		IsInternal:   arg.IsInternal,
		SubmitSource: arg.SubmitSource,

		OfficialExtra: &model.OfficialExtra{
			Realname:          arg.Realname,
			Operator:          arg.Operator,
			Telephone:         arg.Telephone,
			Email:             arg.Email,
			Address:           arg.Address,
			Company:           arg.Company,
			CreditCode:        arg.CreditCode,
			Organization:      arg.Organization,
			OrganizationType:  arg.OrganizationType,
			BusinessLicense:   arg.BusinessLicense,
			BusinessScale:     arg.BusinessScale,
			BusinessLevel:     arg.BusinessLevel,
			BusinessAuth:      arg.BusinessAuth,
			Supplement:        arg.Supplement,
			Professional:      arg.Professional,
			Identification:    arg.Identification,
			OfficalSite:       arg.OfficalSite,
			RegisteredCapital: arg.RegisteredCapital,
		},
	}
	if !od.Validate() {
		log.Error("Failed to validate official doc: %+v", od)
		err = ecode.RequestErr
		return
	}
	return s.dao.OfficialDocSubmit(ctx, od.Mid, od.Name, od.Role, int8(model.OfficialStateWait), od.Title, od.Desc, od.OfficialExtra.String(), od.Uname, od.IsInternal, od.SubmitSource)
}

func (s *Service) updateUname(ctx context.Context, mid int64, name string, adminID int64, adminName string) error {
	b, err := s.dao.Base(ctx, mid)
	if err != nil {
		return err
	}
	if b.Name == name {
		return nil
	}
	if err := s.dao.UpdateUname(ctx, mid, name); err != nil {
		log.Error("Failed to update uname to aso: mid(%d), name(%s): %+v", mid, name, err)
		return err
	}
	if err := s.dao.UpName(ctx, mid, name); err != nil {
		log.Error("Failed to update uname to member: mid(%d), name(%s): %+v", mid, name, err)
		return err
	}
	report.Manager(&report.ManagerInfo{
		Uname:    adminName,
		UID:      adminID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      mid,
		Action:   _logActionEditName,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{0, 0, 0, "", "", ""},
		Content: map[string]interface{}{
			"old_name": b.Name,
			"new_name": name,
		},
	})
	return nil
}

// OfficialDocAddits  find mids by property and value
func (s *Service) OfficialDocAddits(ctx context.Context, property string, vstring string) ([]*model.OfficialDocAddit, error) {
	if property == "" {
		return nil, ecode.RequestErr
	}
	addits, err := s.dao.OfficialDocAddits(ctx, property, vstring)
	return addits, err
}
