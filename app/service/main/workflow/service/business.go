package service

import (
	"context"
	"encoding/json"

	"go-common/app/service/main/workflow/model"
	"go-common/library/log"
)

// UpBusinessExtra update business extra by cid && mid && business
// Deprecated
func (s *Service) UpBusinessExtra(c context.Context, cid int32, mid int64, business int8, key, value string) (err error) {
	var (
		bs    []byte
		bsns  = &model.Business{}
		extra = make(map[string]string)
	)
	if err = s.dao.DB.Where("cid=? and mid=? and business=?", cid, mid, business).Find(bsns).Error; err != nil {
		log.Error("d.DB.WhereBusiness(%d,%d,%d) error(%v)", cid, mid, business, err)
		return
	}
	if bsns.Extra != "" {
		json.Unmarshal([]byte(bsns.Extra), &extra)
	}
	extra[key] = value
	if bs, err = json.Marshal(extra); err != nil {
		log.Error("json.Marshal(%+v)", extra)
		return
	}
	if err = s.dao.DB.Model(&model.Business{}).Where("cid=? and mid=? and business=?", cid, mid, business).Update("extra", string(bs)).Error; err != nil {
		log.Error("s.dao.UpBusinessExtra(%d,%d,%d,%s) error(%v)", cid, mid, business, string(bs), err)
		return
	}
	return
}

// UpBusinessExtraV2 .
func (s *Service) UpBusinessExtraV2(c context.Context, cid int32, mid int64, business int8, key, value string) (err error) {
	var (
		bs    []byte
		bsns  = &model.Business{}
		chall = &model.Challenge{}
		extra = make(map[string]string)
	)
	if err = s.dao.DB.Where("id=?", cid).Find(chall).Error; err != nil {
		log.Error("failed to find challenge cid(%d) error(%v)", cid, err)
		return
	}

	if err = s.dao.DB.Table("workflow_business").Where("business=? and oid=?", chall.Business, chall.Oid).Last(bsns).Error; err != nil {
		log.Error("failed to find last business object business(%d) oid(%d) error(%v)", chall.Business, chall.Oid, err)
		return
	}

	if bsns.Extra != "" {
		json.Unmarshal([]byte(bsns.Extra), &extra)
	}
	extra[key] = value
	if bs, err = json.Marshal(extra); err != nil {
		log.Error("json.Marshal(%+v)", extra)
		return
	}
	if err = s.dao.DB.Model(&model.Business{}).Where("id=?", bsns.ID).Update("extra", string(bs)).Error; err != nil {
		log.Error("s.dao.UpBusinessExtra(%d,%s) error(%v)", bsns.ID, string(bs), err)
		return
	}
	return
}

// BusinessExtra get business extra field by cid
// Deprecated
func (s *Service) BusinessExtra(c context.Context, cid int32, mid int64, business int8) (extra json.RawMessage, err error) {
	bsns := &model.Business{}
	if err = s.dao.DB.Where("cid=? and mid=? and business=?", cid, mid, business).Find(bsns).Error; err != nil {
		log.Error("d.DB.WhereBusiness(%d,%d,%d) error(%v)", cid, mid, business, err)
		return
	}
	if bsns.Extra == "" {
		bsns.Extra = "{}"
	}
	extra = json.RawMessage(bsns.Extra)
	return
}

// BusinessExtraV2 get business extra field by cid
func (s *Service) BusinessExtraV2(c context.Context, cid int32, mid int64, business int8) (extra json.RawMessage, err error) {
	var (
		bsns  = &model.Business{}
		chall = &model.Challenge{}
	)
	if err = s.dao.DB.Where("id=?", cid).Find(chall).Error; err != nil {
		log.Error("failed find challenge cid(%d) error(%v)", cid, err)
		return
	}
	if err = s.dao.DB.Table("workflow_business").Where("business=? and oid=?", chall.Business, chall.Oid).Last(bsns).Error; err != nil {
		log.Error("failed to find last business object business(%d) oid(%d) error(%v)", chall.Business, chall.Oid, err)
		return
	}

	if bsns.Extra == "" {
		bsns.Extra = "{}"
	}
	extra = json.RawMessage(bsns.Extra)
	return
}
