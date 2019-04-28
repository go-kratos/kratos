package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/admin/main/dm/model"
	accountApi "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	archiveMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_workFlowSubtitleBid = 14
)

// SubtitleLanList .
func (s *Service) SubtitleLanList(c context.Context) (res map[int64]string, err error) {
	var (
		sLans []*model.SubtitleLan
	)
	if sLans, err = s.dao.SubtitleLans(c); err != nil {
		return
	}
	res = make(map[int64]string)
	for _, sLan := range sLans {
		res[sLan.Code] = sLan.DocZh
	}
	return
}

// WorkFlowEditSubtitle .
func (s *Service) WorkFlowEditSubtitle(c context.Context, arg *model.WorkFlowSubtitleArg) (err error) {
	var (
		argEdit *model.EditSubtitleArg
		status  model.SubtitleStatus
	)
	if arg == nil || arg.Object == nil || len(arg.Object.Ids) == 0 || len(arg.Targets) == 0 {
		err = ecode.RequestErr
		return
	}
	switch arg.Object.DisposeMode {
	case model.WorkFlowSubtitleDisposeManagerBack:
		status = model.SubtitleStatusManagerBack
	case model.WorkFlowSubtitleDisposeManagerDelete:
		status = model.SubtitleStatusManagerRemove
	default:
		return
	}
	for _, target := range arg.Targets {
		if target == nil {
			continue
		}
		argEdit = &model.EditSubtitleArg{
			Oid:       target.Oid,
			SubtileID: target.Eid,
			Status:    uint8(status),
		}
		// 容错
		if argEdit.Oid == 0 {
			continue
		}
		if err = s.editSubtitle(c, argEdit, false); err != nil {
			log.Error("s.EditSubtitle(arg:%+v),error(%v)", argEdit, err)
			err = nil // ignore error
			return
		}
	}
	return
}

// EditSubtitle .
func (s *Service) EditSubtitle(c context.Context, arg *model.EditSubtitleArg) (err error) {
	return s.editSubtitle(c, arg, true)
}

func (s *Service) editSubtitle(c context.Context, arg *model.EditSubtitleArg, removeWorkFlow bool) (err error) {
	// 更新表
	var (
		subtitle     *model.Subtitle
		argStatus    = model.SubtitleStatus(arg.Status)
		subtitleLans model.SubtitleLans
		sLans        []*model.SubtitleLan
		err1         error
		lanDoc       string
		archiveInfo  *api.Arc
		archiveName  string
	)
	if subtitle, err = s.dao.GetSubtitle(c, arg.Oid, arg.SubtileID); err != nil {
		log.Error("params(oid:%v,subtitleID:%v),error(%v)", arg.Oid, arg.SubtileID, err)
		return
	}
	if subtitle == nil {
		err = ecode.NothingFound
		return
	}
	if argStatus == subtitle.Status {
		err = ecode.SubtitleStatusUnValid
		return
	}
	if subtitle.Status != model.SubtitleStatusPublish && argStatus != model.SubtitleStatusPublish {
		arg.NotifyUpper = false
	}
	switch argStatus {
	case model.SubtitleStatusDraft, model.SubtitleStatusToAudit,
		model.SubtitleStatusAuditBack, model.SubtitleStatusRemove,
		model.SubtitleStatusPublish, model.SubtitleStatusManagerBack,
		model.SubtitleStatusManagerRemove:
	default:
		err = ecode.SubtitleStatusUnValid
		return
	}
	if err = s.changeSubtitleStatus(c, subtitle, argStatus); err != nil {
		log.Error("params(subtitle:%+v,status:%v),error(%v)", subtitle, arg.Status, err)
		return
	}
	if arg.NotifyAuthor || arg.NotifyUpper {
		if sLans, err1 = s.dao.SubtitleLans(c); err1 == nil {
			subtitleLans = model.SubtitleLans(sLans)
		}
		_, lanDoc = subtitleLans.GetByID(int64(subtitle.Lan))
		if archiveInfo, err1 = s.arcRPC.Archive3(c, &archiveMdl.ArgAid2{
			Aid: subtitle.Aid,
		}); err1 != nil {
			log.Error("s.arcRPC.Archive3(aid:%v),error(%v)", subtitle.Aid, err1)
			err1 = nil
		} else {
			archiveName = archiveInfo.Title
		}
	}
	if arg.NotifyAuthor {
		argUser := &model.NotifySubtitleUser{
			Mid:         subtitle.Mid,
			Aid:         subtitle.Aid,
			Oid:         subtitle.Oid,
			SubtitleID:  subtitle.ID,
			ArchiveName: archiveName,
			LanDoc:      lanDoc,
			Status:      model.StatusContent[uint8(subtitle.Status)],
		}
		if err1 = s.dao.SendMsgToSubtitleUser(c, argUser); err1 != nil {
			log.Error("SendMsgToSubtitleUser(argUser:%+v),error(%v)", argUser, err1)
			err1 = nil
		}
	}
	if arg.NotifyUpper {
		var (
			accountInfo *accountApi.InfoReply
			authorName  string
		)
		if accountInfo, err1 = s.accountRPC.Info3(c, &accountApi.MidReq{
			Mid: subtitle.Mid,
		}); err1 != nil {
			log.Error("s.accRPC.Info3(mid:%v),error(%v)", subtitle.Mid, err1)
			err1 = nil
		} else {
			authorName = accountInfo.GetInfo().GetName()
		}
		argUp := &model.NotifySubtitleUp{
			Mid:         subtitle.UpMid,
			AuthorID:    subtitle.Mid,
			AuthorName:  authorName,
			Aid:         subtitle.Aid,
			Oid:         subtitle.Oid,
			SubtitleID:  subtitle.ID,
			ArchiveName: archiveName,
			LanDoc:      lanDoc,
			Status:      model.StatusContent[uint8(subtitle.Status)],
		}
		if err1 = s.dao.SendMsgToSubtitleUp(c, argUp); err1 != nil {
			log.Error("SendMsgToSubtitleUp(argUp:%+v),error(%v)", argUp, err1)
			err1 = nil
		}
	}
	if removeWorkFlow && (argStatus == model.SubtitleStatusRemove || argStatus == model.SubtitleStatusManagerRemove) {
		if err1 := s.dao.WorkFlowAppealDelete(c, _workFlowSubtitleBid, subtitle.Oid, subtitle.ID); err1 != nil {
			log.Error("s.dao.WorkFlowAppealDelete(oid:%v,subtitleID:%v),error(%v)", subtitle.Oid, subtitle.ID, err1)
			return
		}
	}
	return
}

// TODO 确认状态扭转
func (s *Service) changeSubtitleStatus(c context.Context, subtitle *model.Subtitle, status model.SubtitleStatus) (err error) {
	var (
		sc       *model.SubtitleContext
		hasDraft bool
	)
	sc = &model.SubtitleContext{}
	sc.Build(subtitle.Status, status)
	subtitle.PubTime = time.Now().Unix()
	if sc.CheckHasDraft {
		if hasDraft, err = s.CheckHasDraft(c, subtitle); err != nil {
			log.Error("params(subtitle:%+v),error(%v)", subtitle, err)
			return
		}
		if hasDraft {
			err = ecode.SubtitleAlreadyHasDraft
			return
		}
		subtitle.PubTime = 0
	}
	subtitle.Status = status
	if sc.RebuildPub {
		if err = s.RebuildSubtitle(c, subtitle); err != nil {
			log.Error("RebuildSubtitle.params(subtitle:%+v),error(%v)", subtitle, err)
			return
		}
	} else {
		if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
			log.Error("UpdateSubtitle.params(subtitle:%+v),error(%v)", subtitle, err)
			return
		}
	}
	if sc.DraftCache {
		s.dao.DelSubtitleDraftCache(c, subtitle.Oid, subtitle.Type, subtitle.Mid, subtitle.Lan)
	}
	if sc.SubtitleCache {
		s.dao.DelSubtitleCache(c, subtitle.Oid, subtitle.ID)
	}
	if sc.RebuildPub {
		s.dao.DelVideoSubtitleCache(c, subtitle.Oid, subtitle.Type)
	}
	return
}

// SubtitleList .
func (s *Service) SubtitleList(c context.Context, arg *model.SubtitleArg) (res *model.SubtitleList, err error) {
	var (
		searchResult    *model.SearchSubtitleResult
		oidSubtitleIds  map[int64][]int64
		eg              errgroup.Group
		lock            sync.Mutex
		subtitleMap     map[string]*model.Subtitle
		searchSubtitles []*model.SearchSubtitle
		searchSubtitle  *model.SearchSubtitle
		subtitle        *model.Subtitle
		ok              bool
		aids            []int64
		aidMap          map[int64]struct{}
		archives        map[int64]*api.Arc
		archive         *api.Arc
		archiveVideo    *api.Page
		subtitleLans    model.SubtitleLans
		searchArg       *model.SubtitleSearchArg
	)
	if arg.Ps > 100 {
		err = ecode.RequestErr
		return
	}
	key := func(oid, subtitleID int64) string {
		return fmt.Sprintf("%d_%d", oid, subtitleID)
	}
	if sLans, err1 := s.dao.SubtitleLans(c); err1 == nil {
		subtitleLans = model.SubtitleLans(sLans)
	}
	lanCode := subtitleLans.GetByLan(arg.Lan)
	searchArg = &model.SubtitleSearchArg{
		Aid:      arg.Aid,
		Oid:      arg.Oid,
		Mid:      arg.Mid,
		UpperMid: arg.UpperMid,
		Status:   arg.Status,
		Lan:      uint8(lanCode),
		Ps:       arg.Ps,
		Pn:       arg.Pn,
	}
	if searchResult, err = s.dao.SearchSubtitle(c, searchArg); err != nil {
		log.Error("params(arg:%+v).error(%v)", arg, err)
		return
	}
	if searchResult == nil || len(searchResult.Result) == 0 {
		err = ecode.NothingFound
		return
	}
	oidSubtitleIds = make(map[int64][]int64)
	subtitleMap = make(map[string]*model.Subtitle)
	aidMap = make(map[int64]struct{})
	for _, r := range searchResult.Result {
		oidSubtitleIds[r.Oid] = append(oidSubtitleIds[r.Oid], r.ID)
	}
	for oid, subtitleIds := range oidSubtitleIds {
		tempOid := oid
		tempSubtitleIds := subtitleIds
		eg.Go(func() (err error) {
			var subtitles []*model.Subtitle
			if subtitles, err = s.dao.GetSubtitles(context.Background(), tempOid, tempSubtitleIds); err != nil {
				log.Error("params(oid:%v,subtitleIds:%+v).error(%v)", tempOid, tempSubtitleIds, err)
				return
			}
			for _, subtitle := range subtitles {
				lock.Lock()
				aidMap[subtitle.Aid] = struct{}{}
				subtitleMap[key(subtitle.Oid, subtitle.ID)] = subtitle
				lock.Unlock()
			}
			return
		})
	}
	if err = eg.Wait(); err != nil {
		return
	}
	for aid := range aidMap {
		aids = append(aids, aid)
	}
	if archives, err = s.arcRPC.Archives3(c, &archiveMdl.ArgAids2{
		Aids: aids,
	}); err != nil {
		log.Error("prams(aid:%v),error(%v)", aids, err)
		archives = make(map[int64]*api.Arc)
		err = nil
	}
	searchSubtitles = make([]*model.SearchSubtitle, 0, len(searchResult.Result))
	for _, r := range searchResult.Result {
		if subtitle, ok = subtitleMap[key(r.Oid, r.ID)]; !ok {
			continue
		}
		lan, lanDoc := subtitleLans.GetByID(int64(subtitle.Lan))
		searchSubtitle = &model.SearchSubtitle{
			ID:          subtitle.ID,
			Oid:         subtitle.Oid,
			Aid:         subtitle.Aid,
			AuthorID:    subtitle.Mid,
			Status:      uint8(subtitle.Status),
			Lan:         lan,
			LanDoc:      lanDoc,
			IsSign:      subtitle.IsSign,
			IsLock:      subtitle.IsLock,
			Mtime:       subtitle.Mtime,
			SubtitleURL: subtitle.SubtitleURL,
		}
		if archive, ok = archives[subtitle.Aid]; ok {
			searchSubtitle.ArchiveName = archive.Title
		}
		if archiveVideo, err = s.arcRPC.Video3(c, &archiveMdl.ArgVideo2{
			Aid: subtitle.Aid,
			Cid: subtitle.Oid,
		}); err != nil {
			log.Error("params(aid:%v,oid:%v) error(%v)", subtitle.Aid, subtitle.Oid, err)
			err = nil
		} else {
			searchSubtitle.VideoName = archiveVideo.Part
		}
		searchSubtitles = append(searchSubtitles, searchSubtitle)
	}
	res = &model.SubtitleList{
		Page: searchResult.Page,
	}
	res.Subtitles = searchSubtitles
	return
}
