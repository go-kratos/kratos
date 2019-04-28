package service

import (
	"context"
	"fmt"
	"sync"

	"go-common/app/interface/main/dm2/model"
	account "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	archive "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_searchMaxSize = 100
)

// SearchAuthor .
func (s *Service) SearchAuthor(c context.Context, mid int64, status int32, page, size int32) (authorResult *model.SearchSubtitleAuthor, err error) {
	var (
		res           *model.SearchSubtitleResponse
		searchItems   []*model.SearchSubtitleAuthorItem
		formUpper     = false
		countSubtitle *model.CountSubtitleResult
	)
	if size > _searchMaxSize {
		err = ecode.RequestErr
		return
	}
	if res, err = s.searchSubtitle(c, 0, 0, 0, mid, nil, status, page, size, formUpper); err != nil {
		log.Error("psearchSubtitle.params(mid:%v,status:%v,page:%v,size:%v),error(%v)", mid, status, page, size, err)
		return
	}
	if res == nil || res.Page == nil {
		log.Error("psearchSubtitle.params(mid:%v,status:%v,page:%v,size:%v),error(%v)", mid, status, page, size, err)
		err = ecode.NothingFound
		return
	}
	if countSubtitle, err = s.dao.CountSubtitles(c, mid, nil, 0, 0, 0); err != nil {
		log.Error("CountSubtitles.params(mid:%v),error(%v)", mid, err)
		return
	}
	for _, rs := range res.Subtitles {
		searchItem := &model.SearchSubtitleAuthorItem{
			ID:          rs.ID,
			Oid:         rs.Oid,
			Aid:         rs.Aid,
			Type:        rs.Type,
			ArchiveName: rs.ArchiveName,
			ArchivePic:  rs.ArchivePic,
			VideoName:   rs.VideoName,
			Lan:         rs.Lan,
			LanDoc:      rs.LanDoc,
			Status:      rs.Status,
			IsSign:      rs.IsSign,
			IsLock:      rs.IsLock,
			Mtime:       rs.Mtime,
		}
		if searchItem.Status == int32(model.SubtitleStatusManagerBack) {
			searchItem.Status = int32(model.SubtitleStatusAuditBack)
		}
		if searchItem.Status == int32(model.SubtitleStatusAuditBack) {
			searchItem.RejectComment = rs.RejectComment
		}
		searchItems = append(searchItems, searchItem)
	}
	authorResult = &model.SearchSubtitleAuthor{
		Page:         res.Page,
		Subtitles:    searchItems,
		DraftCount:   countSubtitle.Draft,
		AuditCount:   countSubtitle.ToAudit,
		BackCount:    countSubtitle.AuditBack,
		PublishCount: countSubtitle.Publish,
	}
	authorResult.Total = authorResult.DraftCount + authorResult.AuditCount + authorResult.BackCount + authorResult.PublishCount
	return
}

// SearchAssist .
func (s *Service) SearchAssist(c context.Context, aid, oid int64, tp int32, mid int64, status int32, page, size int32) (assistResult *model.SearchSubtitleAssit, err error) {
	var (
		res           *model.SearchSubtitleResponse
		upMids        []int64
		fromUpper     = true
		countSubtitle *model.CountSubtitleResult
	)
	if size > _searchMaxSize {
		err = ecode.RequestErr
		return
	}
	upMids = append(upMids, mid)
	if res, err = s.searchSubtitle(c, aid, oid, tp, 0, upMids, status, page, size, fromUpper); err != nil {
		return
	}
	if res == nil || res.Page == nil {
		err = ecode.NothingFound
		return
	}
	if countSubtitle, err = s.dao.CountSubtitles(c, 0, upMids, aid, oid, tp); err != nil {
		log.Error("CountSubtitles.params(mid:%v),error(%v)", mid, err)
		return
	}
	assistResult = &model.SearchSubtitleAssit{
		Page:         res.Page,
		Subtitles:    res.Subtitles,
		AuditCount:   countSubtitle.ToAudit,
		PublishCount: countSubtitle.Publish,
	}
	assistResult.Total = assistResult.AuditCount + assistResult.PublishCount
	return
}

func (s *Service) buildSearchStatus(c context.Context, status int32, fromUpper bool) (searchStatus []int64, err error) {
	if fromUpper {
		switch status {
		case 0:
			searchStatus = []int64{
				int64(model.SubtitleStatusPublish),
				int64(model.SubtitleStatusToAudit),
				int64(model.SubtitleStatusCheckPublish),
			}
		case int32(model.SubtitleStatusPublish):
			searchStatus = []int64{
				int64(model.SubtitleStatusPublish),
				int64(model.SubtitleStatusCheckPublish),
			}
		case int32(model.SubtitleStatusToAudit):
			searchStatus = []int64{int64(status)}
		default:
			err = ecode.SubtitlePermissionDenied
			return
		}
	} else {
		switch status {
		case int32(model.SubtitleStatusPublish):
			searchStatus = []int64{
				int64(model.SubtitleStatusPublish),
				int64(model.SubtitleStatusCheckPublish),
			}
		case int32(model.SubtitleStatusAuditBack):
			searchStatus = []int64{
				int64(model.SubtitleStatusAuditBack),
				int64(model.SubtitleStatusManagerBack),
			}
		case int32(model.SubtitleStatusDraft):
			searchStatus = []int64{int64(status)}
		case int32(model.SubtitleStatusToAudit):
			searchStatus = []int64{
				int64(model.SubtitleStatusCheckToAudit),
				int64(model.SubtitleStatusToAudit),
			}
		case 0:
			searchStatus = []int64{
				int64(model.SubtitleStatusPublish),
				int64(model.SubtitleStatusToAudit),
				int64(model.SubtitleStatusDraft),
				int64(model.SubtitleStatusAuditBack),
				int64(model.SubtitleStatusCheckToAudit),
				int64(model.SubtitleStatusCheckPublish),
				int64(model.SubtitleStatusManagerBack),
			}
		default:
			err = ecode.SubtitlePermissionDenied
			return
		}
	}
	return
}

func (s *Service) searchSubtitle(c context.Context, aid, oid int64, tp int32, mid int64, upMids []int64, status int32, page, size int32, fromUpper bool) (result *model.SearchSubtitleResponse, err error) {
	var (
		res            *model.SearchSubtitleResult
		dmidsMap       map[int64][]int64
		eg             errgroup.Group
		subtitleMap    map[string]*model.Subtitle
		subtitle       *model.Subtitle
		results        []*model.SearchSubtitle
		searchSubtitle *model.SearchSubtitle
		mutex          sync.Mutex
		ok             bool
		profileReply   *account.ProfileReply
		archiveAids    []int64
		archiveMap     map[int64]*api.Arc
		searchStatus   []int64
	)
	key := func(oid, id int64) string {
		return fmt.Sprintf("%d:%d", oid, id)
	}
	if searchStatus, err = s.buildSearchStatus(c, status, fromUpper); err != nil {
		return
	}
	if res, err = s.dao.SearchSubtitles(c, page, size, mid, upMids, aid, oid, tp, searchStatus); err != nil {
		return
	}
	if res == nil || res.Page == nil {
		return
	}
	result = &model.SearchSubtitleResponse{
		Page: &model.SearchPage{
			Num:   res.Page.Num,
			Size:  res.Page.Size,
			Total: res.Page.Total,
		},
	}
	dmidsMap = make(map[int64][]int64)
	subtitleMap = make(map[string]*model.Subtitle)
	for _, rs := range res.Results {
		dmidsMap[rs.Oid] = append(dmidsMap[rs.Oid], rs.ID)
	}
	for oid, ids := range dmidsMap {
		tempOid := oid
		tempIds := ids
		eg.Go(func() (err error) {
			var (
				subtitles map[int64]*model.Subtitle
			)
			if subtitles, err = s.getSubtitles(c, tempOid, tempIds); err != nil {
				log.Error("params(tempOid:%v, tempIds:%v) error(%v)", tempOid, tempIds, err)
				return
			}
			mutex.Lock()
			for _, subtitle := range subtitles {
				subtitleMap[key(subtitle.Oid, subtitle.ID)] = subtitle
				archiveAids = append(archiveAids, subtitle.Aid)
			}
			mutex.Unlock()
			return
		})
	}
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait() error(%v)", err)
		return
	}
	if archiveMap, err = s.arcRPC.Archives3(c, &archive.ArgAids2{
		Aids: archiveAids,
	}); err != nil {
		log.Error("params(aids:%v) error(%v)", archiveAids, err)
		return
	}
	results = make([]*model.SearchSubtitle, 0, len(subtitleMap))
	for _, rs := range res.Results {
		var (
			archiveName  string
			archivePic   string
			archiveVideo *api.Page
			videoName    string
		)
		if subtitle, ok = subtitleMap[key(rs.Oid, rs.ID)]; !ok {
			continue
		}
		if a, ok := archiveMap[subtitle.Aid]; ok {
			archiveName = a.Title
			archivePic = a.Pic
		}
		if archiveVideo, err = s.arcRPC.Video3(c, &archive.ArgVideo2{
			Aid: subtitle.Aid,
			Cid: subtitle.Oid,
		}); err != nil {
			log.Error("params(aid:%v,oid:%v) error(%v)", subtitle.Aid, subtitle.Oid, err)
			err = nil
		} else {
			videoName = archiveVideo.Part
		}
		lan, lanDoc := s.subtitleLans.GetByID(int64(subtitle.Lan))
		searchSubtitle = &model.SearchSubtitle{
			ID:            rs.ID,
			Oid:           rs.Oid,
			Aid:           subtitle.Aid,
			Type:          subtitle.Type,
			ArchiveName:   archiveName,
			ArchivePic:    archivePic,
			VideoName:     videoName,
			Lan:           lan,
			LanDoc:        lanDoc,
			Status:        int32(subtitle.Status),
			IsSign:        subtitle.IsSign,
			IsLock:        subtitle.IsLock,
			RejectComment: subtitle.RejectComment,
			Mtime:         subtitle.Mtime,
		}
		switch subtitle.Status {
		case model.SubtitleStatusCheckToAudit:
			searchSubtitle.Status = int32(model.SubtitleStatusToAudit)
		case model.SubtitleStatusCheckPublish:
			searchSubtitle.Status = int32(model.SubtitleStatusPublish)
		}
		// up主搜索需要以下字段
		if fromUpper {
			var (
				profileName string
				authorPic   string
			)
			if profileReply, err = s.accountRPC.Profile3(c, &account.MidReq{
				Mid: subtitle.Mid,
			}); err != nil {
				log.Error("params(Mid:%v) error(%v)", subtitle.Mid, err)
				err = nil
			} else {
				profileName = profileReply.GetProfile().GetName()
				authorPic = profileReply.GetProfile().GetFace()
			}
			searchSubtitle.Author = profileName
			searchSubtitle.AuthorPic = authorPic
			searchSubtitle.AuthorID = subtitle.AuthorID

		}
		results = append(results, searchSubtitle)
	}
	result.Subtitles = results
	return
}
