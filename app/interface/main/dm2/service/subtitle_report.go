package service

import (
	"context"

	"go-common/app/interface/main/dm2/model"
	"go-common/app/service/main/archive/api"
	archiveMdl "go-common/app/service/main/archive/model/archive"
	figureMdl "go-common/app/service/main/figure/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_workFlowSubtitleBid = 14
	_workFlowSubtitleRid = 1
)

// SubtitleReportList .
func (s *Service) SubtitleReportList(c context.Context) (data []*model.WorkFlowTag, err error) {
	var (
		cacheErr bool
	)
	if data, err = s.dao.SubtitleWorlFlowTagCache(c, _workFlowSubtitleBid, _workFlowSubtitleRid); err != nil {
		cacheErr = true
		err = nil
	}
	if len(data) > 0 {
		return
	}
	if data, err = s.dao.WorkFlowTagList(c, _workFlowSubtitleBid, _workFlowSubtitleRid); err != nil {
		return
	}
	if !cacheErr {
		temp := data
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetSubtitleWorlFlowTagCache(ctx, _workFlowSubtitleBid, _workFlowSubtitleRid, temp)
		})
	}
	return
}

// SubtitleReportAdd .
func (s *Service) SubtitleReportAdd(c context.Context, mid int64, param *model.SubtitleReportAddParam) (err error) {
	var (
		figureWithRank *figureMdl.FigureWithRank
		subtitle       *model.Subtitle
		archiveInfo    *api.Arc
		score          int32
	)
	if subtitle, err = s.getSubtitle(c, param.Oid, param.SubtitleID); err != nil {
		return
	}
	if subtitle == nil {
		err = ecode.NothingFound
		return
	}
	if figureWithRank, err = s.figureRPC.UserFigure(c, &figureMdl.ArgUserFigure{
		Mid: mid,
	}); err == nil {
		score = figureWithRank.Score
	} else {
		log.Error("UserFigure(mid:%v),error(%v)", mid, err)
	}
	if archiveInfo, err = s.arcRPC.Archive3(c, &archiveMdl.ArgAid2{
		Aid: subtitle.Aid,
	}); err != nil {
		log.Error("s.arcRPC.Archive3(aid:%v),error(%v)", subtitle.Aid, err)
		return
	}
	req := &model.WorkFlowAppealAddReq{
		Business:       _workFlowSubtitleBid,
		Oid:            param.Oid,
		Aid:            subtitle.Aid,
		Rid:            _workFlowSubtitleRid,
		LanCode:        int64(subtitle.Lan),
		SubtitleID:     param.SubtitleID,
		Score:          score,
		Tid:            param.Tid,
		Mid:            mid,
		Description:    param.MetaData,
		BusinessTypeID: archiveInfo.TypeID,
		BusinessTitle:  param.Content,
		BusinessMid:    subtitle.Mid,
		Extra: &model.WorkFlowAppealAddExtra{
			SubtitleStatus: int64(subtitle.Status),
			SubtitleURL:    subtitle.SubtitleURL,
			ArchiveName:    archiveInfo.Title,
		},
	}
	if err = s.dao.WorkFlowAppealAdd(c, req); err != nil {
		log.Error("SubtitleReportAdd(req:%+v),error(%v)", req, err)
		return
	}
	return
}

func (s *Service) subtitleReportDelete(c context.Context, oid, subtitleID int64) (err error) {
	if err = s.dao.WorkFlowAppealDelete(c, _workFlowSubtitleBid, oid, subtitleID); err != nil {
		return
	}
	return
}
