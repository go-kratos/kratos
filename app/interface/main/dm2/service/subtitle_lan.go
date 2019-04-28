package service

import (
	"context"
	"sort"

	"encoding/json"
	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

const (
	_subtitleLanFileName = "subtitle_lan.json"
)

// SubtitleLanOp .
func (s *Service) SubtitleLanOp(c context.Context, code uint8, lan, docZh, docEn string, isDelete bool) (err error) {
	var (
		subtitleLans []*model.SubtitleLan
		subtitleLan  *model.SubtitleLan
		bs           []byte
	)
	subtitleLan = &model.SubtitleLan{
		Code:     int64(code),
		Lan:      lan,
		DocZh:    docZh,
		DocEn:    docEn,
		IsDelete: isDelete,
	}
	if err = s.dao.SubtitleLanAdd(c, subtitleLan); err != nil {
		log.Error("params(subtitleLan:%+v).error(%v)", subtitleLan, err)
		return
	}
	if subtitleLans, err = s.dao.SubtitleLans(c); err != nil {
		log.Error("SubtitleLans.error(%v)", err)
		return
	}
	if bs, err = json.Marshal(subtitleLans); err != nil {
		log.Error("json.Marshal.params(subtitleLan:%+v).error(%v)", subtitleLan, err)
		return
	}
	// reload bfs
	if _, err = s.dao.UploadBfs(c, _subtitleLanFileName, bs); err != nil {
		log.Error("UploadBfs.params.error(%v)", err)
		return
	}
	return
}

func (s *Service) isSubtitleLanLock(c context.Context, oid int64, tp int32, lan string) (subtitle *model.Subtitle, err error) {
	var (
		vss        []*model.VideoSubtitle
		subtitleID int64
	)
	if vss, err = s.getVideoSubtitles(c, oid, tp); err != nil {
		log.Error("params(oid:%v, tp:%v) error(%v)", oid, tp, err)
		return
	}
	for _, vs := range vss {
		if vs.Lan == lan {
			subtitleID = vs.ID
			break
		}
	}
	if subtitleID > 0 {
		if subtitle, err = s.getSubtitle(c, oid, subtitleID); err != nil {
			log.Error("params(oid:%v, subtitleID:%v) error(%v)", oid, subtitleID, err)
			return
		}
	}
	return
}

// SubtitleLans .
func (s *Service) SubtitleLans(c context.Context, oid int64, tp int32, mid int64) (lans []*model.Language, err error) {
	var (
		vss           []*model.VideoSubtitle
		res           *model.SearchSubtitleResult
		_searchStatus = []int64{int64(model.SubtitleStatusDraft), int64(model.SubtitleStatusToAudit), int64(model.SubtitleStatusAuditBack), int64(model.SubtitleStatusCheckToAudit)}
		subtitleIds   []int64
		subtitles     []*model.Subtitle
		subtitlesM    map[int64]*model.Subtitle
		mapLans       map[string]*model.Language
		_maxPageSize  = 200
		subtitleMap   map[string][]*model.Subtitle
		lan           *model.Language
		ok            bool
	)
	mapLans = make(map[string]*model.Language)
	if vss, err = s.getVideoSubtitles(c, oid, tp); err != nil {
		log.Error("params(oid:%v,tp:%v).error(%v)", oid, tp, err)
		return
	}
	if res, err = s.dao.SearchSubtitles(c, 1, int32(_maxPageSize), mid, nil, 0, oid, tp, _searchStatus); err == nil {
		if res != nil {
			for _, result := range res.Results {
				subtitleIds = append(subtitleIds, result.ID)
			}
		}
	} else {
		log.Error("SearchSubtitles.params(mid:%v,oid:%v),error(%v)", mid, oid, err)
		err = nil
	}
	for _, vs := range vss {
		mapLans[vs.Lan] = &model.Language{
			Lan:    vs.Lan,
			LanDoc: vs.LanDoc,
			Pub: &model.LanguagePub{
				SubtitleID: vs.ID,
				IsLock:     vs.IsLock,
				IsPub:      true,
			},
		}
	}
	if len(subtitleIds) > 0 {
		if subtitlesM, err = s.getSubtitles(c, oid, subtitleIds); err != nil {
			log.Error("params(oid:%v,subtitleDraftIds:%v).error(%v)", oid, subtitleIds, err)
			return
		}
		subtitleMap = make(map[string][]*model.Subtitle)
		for _, subtitle := range subtitlesM {
			tlan, tlanDoc := s.subtitleLans.GetByID(int64(subtitle.Lan))
			if lan, ok = mapLans[tlan]; !ok {
				lan = &model.Language{
					Lan:    tlan,
					LanDoc: tlanDoc,
				}
				mapLans[tlan] = lan
			}
			switch subtitle.Status {
			case model.SubtitleStatusDraft:
				lan.Draft = &model.LanguageID{
					SubtitleID: subtitle.ID,
				}
			case model.SubtitleStatusToAudit:
				lan.Audit = &model.LanguageID{
					SubtitleID: subtitle.ID,
				}
			case model.SubtitleStatusAuditBack:
				subtitleMap[tlan] = append(subtitleMap[tlan], subtitle)
			}
		}
		for _, subtitles = range subtitleMap {
			if len(subtitles) > 0 {
				sort.Slice(subtitles, func(i, j int) bool {
					return subtitles[i].PubTime > subtitles[j].PubTime
				})
				tlan, tlanDoc := s.subtitleLans.GetByID(int64(subtitles[0].Lan))
				if lan, ok = mapLans[tlan]; !ok {
					lan = &model.Language{
						Lan:    tlan,
						LanDoc: tlanDoc,
					}
					mapLans[tlan] = lan
				}
				lan.AuditBack = &model.LanguageID{
					SubtitleID: subtitles[0].ID,
				}
			}
		}
	}
	lans = make([]*model.Language, 0, len(mapLans))
	for _, mapLan := range mapLans {
		lans = append(lans, mapLan)
	}
	return
}
