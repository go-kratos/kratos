package service

import (
	"context"
	"time"

	"go-common/app/admin/main/activity/model"
	articlemodel "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/jinzhu/gorm"
)

// GetArticleMetas from rpc .
func (s *Service) GetArticleMetas(c context.Context, aids []int64) (res map[int64]*articlemodel.Meta, err error) {
	if res, err = s.artRPC.ArticleMetas(c, &articlemodel.ArgAids{Aids: aids}); err != nil {
		log.Error("s.ArticleMetas(%v) error(%v)", aids, err)
	}
	return
}

// SubjectList get subject list .
func (s *Service) SubjectList(c context.Context, listParams *model.ListSub) (listRes *model.SubListRes, err error) {
	var (
		count int64
		list  []*model.ActSubject
	)
	db := s.DB
	if listParams.Keyword != "" {
		names := listParams.Keyword + "%"
		db = db.Where("`id` = ? or `name` like ? or `author` like ?", listParams.Keyword, names, names)
	}
	if listParams.Sctime != 0 {
		parseScime := time.Unix(listParams.Sctime, 0)
		db = db.Where("ctime >= ?", parseScime.Format("2006-01-02 15:04:05"))
	}
	if listParams.Ectime != 0 {
		parseEcime := time.Unix(listParams.Ectime, 0)
		db = db.Where("etime <= ?", parseEcime.Format("2006-01-02 15:04:05"))
	}
	if len(listParams.States) > 0 {
		db = db.Where("state in (?)", listParams.States)
	}
	if len(listParams.Types) > 0 {
		db = db.Where("type in (?)", listParams.Types)
	}
	if err = db.Offset((listParams.Page - 1) * listParams.PageSize).Limit(listParams.PageSize).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error(" db.Model(&model.ActSubject{}).Find() args(%v) error(%v)", listParams, err)
		return
	}
	if err = db.Model(&model.ActSubject{}).Count(&count).Error; err != nil {
		log.Error("db.Model(&model.ActSubject{}).Count() args(%v) error(%v)", listParams, err)
		return
	}
	listRes = &model.SubListRes{
		List: list,
		Page: &model.PageRes{
			Num:   listParams.Page,
			Size:  listParams.PageSize,
			Total: count,
		},
	}
	return
}

// VideoList .
func (s *Service) VideoList(c context.Context) (res []*model.ActSubjectResult, err error) {
	var (
		types    = []int{1, 4}
		list     []*model.ActSubject
		likeList []*model.Like
	)
	db := s.DB
	if err = db.Where("state = ?", 1).Where("type in (?)", types).Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db.Model(&model.ActSubject{}).Where(state = ?, 1).Where(type in (?), %v).Find() error(%v)", types, err)
		return
	}
	listCount := len(list)
	if listCount == 0 {
		return
	}
	sids := make([]int64, 0, listCount)
	for _, value := range list {
		sids = append(sids, value.ID)
	}
	if err = db.Where("sid in (?)", sids).Where("wid > ?", 0).Find(&likeList).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db.Model(&model.Like{}).Where(sid in (?), %v).Find() error(%v)", sids, err)
		return
	}
	hashList := make(map[int64][]int64)
	for _, value := range likeList {
		hashList[value.Sid] = append(hashList[value.Sid], value.Wid)
	}
	res = make([]*model.ActSubjectResult, 0, len(list))
	for _, value := range list {
		rs := &model.ActSubjectResult{
			ActSubject: value,
		}
		if v, ok := hashList[value.ID]; ok {
			rs.Aids = v
		}
		res = append(res, rs)
	}
	return
}

// AddActSubject .
func (s *Service) AddActSubject(c context.Context, params *model.AddList) (res int64, err error) {
	if params.ScreenSet != 2 {
		params.ScreenSet = 1
	}
	protect := &model.ActSubjectProtocol{
		Protocol:  params.Protocol,
		Types:     params.Types,
		Pubtime:   params.Pubtime,
		Deltime:   params.Deltime,
		Editime:   params.Editime,
		Tags:      params.Tags,
		Hot:       params.Hot,
		BgmID:     params.BgmID,
		Oids:      params.Oids,
		ScreenSet: params.ScreenSet,
		PasterID:  params.PasterID,
	}
	actTime := &model.ActTimeConfig{
		Interval: params.Interval,
		Tlimit:   params.Tlimit,
		Ltime:    params.Ltime,
	}
	if params.Tags != "" {
		if err = s.dao.AddTags(c, params.Tags, metadata.String(c, metadata.RemoteIP)); err != nil {
			log.Error("s.AddTags(%s,) error(%v)", params.Tags, err)
			return
		}
	}
	actSub := &model.ActSubject{
		Oid:        params.ActSubject.Oid,
		Type:       params.ActSubject.Type,
		State:      params.ActSubject.State,
		Level:      params.ActSubject.Level,
		Flag:       params.ActSubject.Flag,
		Rank:       params.ActSubject.Rank,
		Stime:      params.ActSubject.Stime,
		Etime:      params.ActSubject.Etime,
		Lstime:     params.ActSubject.Lstime,
		Letime:     params.ActSubject.Letime,
		Uetime:     params.ActSubject.Uetime,
		Ustime:     params.ActSubject.Ustime,
		Name:       params.ActSubject.Name,
		Author:     params.ActSubject.Author,
		ActURL:     params.ActSubject.ActURL,
		Cover:      params.ActSubject.Cover,
		Dic:        params.ActSubject.Dic,
		H5Cover:    params.ActSubject.H5Cover,
		LikeLimit:  params.ActSubject.LikeLimit,
		AndroidURL: params.ActSubject.AndroidURL,
		IosURL:     params.ActSubject.IosURL,
	}
	if err = s.DB.Create(actSub).Error; err != nil {
		log.Error("s.DB.Create(%v) error(%v)", actSub, err)
		return
	}
	protect.Sid = actSub.ID
	if err = s.DB.Create(protect).Error; err != nil {
		log.Error("s.DB.Create(%v) error(%v)", protect, err)
		return
	}
	if params.Type == model.ONLINEVOTE {
		actTime.Sid = actSub.ID
		if err = s.DB.Create(actTime).Error; err != nil {
			log.Error("s.DB.Create(%v) error(%v)", actTime, err)
			return
		}
	}
	res = actSub.ID
	return
}

// UpActSubject .
func (s *Service) UpActSubject(c context.Context, params *model.AddList, sid int64) (res int64, err error) {
	if params.ScreenSet != 2 {
		params.ScreenSet = 1
	}
	onlineData := &model.ActTimeConfig{
		Interval: params.Interval,
		Tlimit:   params.Tlimit,
		Ltime:    params.Ltime,
	}
	actSubject := new(model.ActSubject)
	if err = s.DB.Where("id = ?", sid).Last(actSubject).Error; err != nil {
		log.Error("s.DB.Where(id = ?, %d).Last(%v) error(%v)", sid, actSubject, err)
		return
	}
	data := map[string]interface{}{
		"Oid":        params.Oid,
		"Type":       params.Type,
		"State":      params.State,
		"Level":      params.Level,
		"Flag":       params.Flag,
		"Rank":       params.Rank,
		"Stime":      params.Stime,
		"Etime":      params.Etime,
		"Lstime":     params.Lstime,
		"Uetime":     params.Uetime,
		"Ustime":     params.Ustime,
		"Name":       params.Name,
		"Author":     params.Author,
		"ActURL":     params.ActURL,
		"Cover":      params.Cover,
		"Dic":        params.Dic,
		"H5Cover":    params.H5Cover,
		"LikeLimit":  params.LikeLimit,
		"AndroidURL": params.AndroidURL,
		"IosURL":     params.IosURL,
	}
	if err = s.DB.Model(&model.ActSubject{}).Where("id = ?", sid).Update(data).Error; err != nil {
		log.Error("s.DB.Model(&model.ActSubject{}).Where(id = ?, %d).Update(%v) error(%v)", sid, data, err)
		return
	}
	item := new(model.ActSubjectProtocol)
	if err = s.DB.Where("sid = ? ", sid).Last(item).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.DB.Where(sid = ? , %d).Last(%v) error(%v)", sid, item, err)
		return
	}
	//item有值
	if item.ID > 0 {
		if params.Tags != "" {
			if item.Tags != params.Tags {
				if err = s.dao.AddTags(c, params.Tags, metadata.String(c, metadata.RemoteIP)); err != nil {
					log.Error("s.AddTags(%s) error(%v)", params.Tags, err)
					return
				}
			}
		}
		upProtectData := map[string]interface{}{
			"Protocol":  params.Protocol,
			"Types":     params.Types,
			"Pubtime":   params.Pubtime,
			"Deltime":   params.Deltime,
			"Editime":   params.Editime,
			"Hot":       params.Hot,
			"BgmID":     params.BgmID,
			"Oids":      params.Oids,
			"ScreenSet": params.ScreenSet,
			"PasterID":  params.PasterID,
			"Tags":      params.Tags,
		}
		if err = s.DB.Model(&model.ActSubjectProtocol{}).Where("id = ?", item.ID).Update(upProtectData).Error; err != nil {
			log.Error("s.DB.Model(&model.ActSubjectProtocol{}).Where(id = ?, %d).Update(%v) error(%v)", item.ID, upProtectData, err)
			return
		}
	} else {
		protectDtata := &model.ActSubjectProtocol{
			Protocol:  params.Protocol,
			Types:     params.Types,
			Pubtime:   params.Pubtime,
			Deltime:   params.Deltime,
			Editime:   params.Editime,
			Hot:       params.Hot,
			BgmID:     params.BgmID,
			Oids:      params.Oids,
			ScreenSet: params.ScreenSet,
			PasterID:  params.PasterID,
			Sid:       sid,
		}
		if err = s.DB.Create(protectDtata).Error; err != nil {
			log.Error("s.DB.Create(%v) error(%v)", protectDtata, err)
			return
		}
	}
	if actSubject.Type == model.ONLINEVOTE {
		onlineData.Sid = sid
		output := new(model.ActTimeConfig)
		if err = s.DB.Where("sid = ?", sid).Last(output).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("s.DB.Where(sid = ?, %d).Last(%v) error(%v)", sid, output, err)
			return
		}
		if output.ID > 0 {
			if err = s.DB.Model(&model.ActTimeConfig{}).Where("id = ?", output.ID).Update(onlineData).Error; err != nil {
				log.Error("s.DB.Model(&model.ActTimeConfig{}).Where(id = ?, %d).Update(%v) error(%v)", output.ID, onlineData, err)
				return
			}
		}

	}
	res = sid
	return
}

// SubProtocol .
func (s *Service) SubProtocol(c context.Context, sid int64) (res *model.ActSubjectProtocol, err error) {
	res = &model.ActSubjectProtocol{}
	if err = s.DB.Where("sid = ?", sid).First(res).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.DB.Where(sid = %d ) error(%v)", sid, err)
	}
	return
}

// TimeConf .
func (s *Service) TimeConf(c context.Context, sid int64) (res *model.ActTimeConfig, err error) {
	res = new(model.ActTimeConfig)
	if err = s.DB.Where("sid = ?", sid).First(res).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("actSrv.DB.Where(sid = ?, %d) error(%v)", sid, err)
	}
	return
}
