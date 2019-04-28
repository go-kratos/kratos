package service

import (
	"context"

	pb "go-common/app/service/main/sms/api"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_tableTemplate = "sms_template_new"

	_searchTypeCode    = "code"
	_searchTypeContent = "content"
)

func (s *Service) templateByID(ctx context.Context, id int64) (res *smsmdl.ModelTemplate, err error) {
	res = new(smsmdl.ModelTemplate)
	if err = s.db.Table(_tableTemplate).Where("id=?", id).First(&res).Error; err != nil {
		res = nil
		if err == ecode.NothingFound {
			err = nil
			return
		}
		log.Error("templateByID(%d) error(%v)", id, err)
		return
	}
	return
}

func (s *Service) templateByCode(ctx context.Context, code string) (res *smsmdl.ModelTemplate, err error) {
	res = new(smsmdl.ModelTemplate)
	if err = s.db.Table(_tableTemplate).Where("code=?", code).First(&res).Error; err != nil {
		res = nil
		if err == ecode.NothingFound {
			err = nil
			return
		}
		log.Error("templateByCode(%s) error(%v)", code, err)
		return
	}
	return
}

// AddTemplate add template
func (s *Service) AddTemplate(ctx context.Context, req *pb.AddTemplateReq) (res *pb.AddTemplateReply, err error) {
	tpl, err := s.templateByCode(ctx, req.Tcode)
	if err != nil {
		return
	}
	if tpl != nil {
		err = ecode.SmsTemplateCodeExist
		return
	}
	t := &smsmdl.ModelTemplate{
		Code:      req.Tcode,
		Template:  req.Template,
		Stype:     req.Stype,
		Status:    smsmdl.TemplateStatusApprovel,
		Submitter: req.Submitter,
	}
	if err = s.db.Table(_tableTemplate).Create(t).Error; err != nil {
		log.Error("s.AddTemplate(%+v) error(%v)", req, err)
	}
	return
}

// UpdateTemplate update template
func (s *Service) UpdateTemplate(ctx context.Context, req *pb.UpdateTemplateReq) (res *pb.UpdateTemplateReply, err error) {
	tpl, err := s.templateByID(ctx, req.ID)
	if err != nil {
		return
	}
	if tpl == nil {
		err = ecode.SmsTemplateNotExist
		return
	}
	if tpl.Code != req.Tcode {
		var t *smsmdl.ModelTemplate
		if t, err = s.templateByCode(ctx, req.Tcode); err != nil {
			return
		}
		if t != nil {
			err = ecode.SmsTemplateCodeExist
			return
		}
	}
	m := map[string]interface{}{
		"code":      req.Tcode,
		"template":  req.Template,
		"stype":     req.Stype,
		"status":    req.Status,
		"submitter": req.Submitter,
	}
	if err = s.db.Table(_tableTemplate).Where("id=?", req.ID).Update(m).Error; err != nil {
		log.Error("s.UpdateTemplate(%+v) error(%v)", req, err)
	}
	return
}

// TemplateList template list
func (s *Service) TemplateList(ctx context.Context, req *pb.TemplateListReq) (res *pb.TemplateListReply, err error) {
	res = new(pb.TemplateListReply)
	cond := "1=1"
	if req.St == _searchTypeCode {
		cond = "code like '%" + req.Sw + "%'"
	} else if req.St == _searchTypeContent {
		cond = "template like '%" + req.Sw + "%'"
	}
	start := (req.Pn - 1) * req.Ps
	if err = s.db.Table(_tableTemplate).Where(cond).Order("id desc").Offset(start).Limit(req.Ps).Find(&res.List).Error; err != nil {
		log.Error("s.TemplateList(%d,%d,%s,%s) error(%v)", req.Pn, req.Ps, req.St, req.Sw, err)
		return
	}
	err = s.db.Table(_tableTemplate).Where(cond).Count(&res.Total).Error
	return
}
