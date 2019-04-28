package service

import (
	"context"
	"encoding/json"
	"sync"

	"go-common/app/admin/main/esports/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_arcsSize     = 50
	_ReplyTypeAct = "25"
)

var (
	_emptyModules = make([]*model.Module, 0)
	_emptMatchMod = make([]*model.MatchModule, 0)
	_emptyActive  = make([]*model.MatchModule, 0)
)

// AddAct .
func (s *Service) AddAct(c context.Context, param *model.ParamMA) (arcs map[string][]int64, err error) {
	var ms []*model.Module
	if param.Modules != "" {
		if err = json.Unmarshal([]byte(param.Modules), &ms); err != nil {
			err = ecode.EsportsActModErr
			return
		}
		if arcs, err = s.checkArc(c, ms); err != nil {
			return
		}
	}
	tx := s.dao.DB.Begin()
	if err = tx.Model(&model.MatchActive{}).Create(&param.MatchActive).Error; err != nil {
		log.Error("AddAct MatchActive tx.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	maID := param.ID
	if len(ms) > 0 {
		if err = tx.Model(&model.Module{}).Exec(model.BatchAddModuleSQL(maID, ms)).Error; err != nil {
			log.Error("AddAct Module tx.Model Create(%+v) error(%v)", param, err)
			err = tx.Rollback().Error
			return
		}
	}
	if err = tx.Commit().Error; err != nil {
		return
	}
	// register reply
	if err = s.dao.RegReply(c, maID, param.Adid, _ReplyTypeAct); err != nil {
		err = nil
	}
	return
}

// EditAct .
func (s *Service) EditAct(c context.Context, param *model.ParamMA) (arcs map[string][]int64, err error) {
	var (
		tmpMs, upM, addM, ms []*model.Module
		mapMID               map[int64]int64
		pMID                 map[int64]*model.Module
		delM                 []int64
	)
	if param.Modules != "" {
		if err = json.Unmarshal([]byte(param.Modules), &ms); err != nil {
			err = ecode.EsportsActModErr
			return
		}
		if arcs, err = s.checkArc(c, ms); err != nil {
			return
		}
	}
	// check module
	if err = s.dao.DB.Model(&model.Module{}).Where("ma_id=?", param.ID).Where("status=?", _notDeleted).Find(&tmpMs).Error; err != nil {
		log.Error("EditAct s.dao.DB.Model Find (%+v) error(%v)", param.ID, err)
		return
	}
	mapMID = make(map[int64]int64, len(tmpMs))
	for _, m := range tmpMs {
		mapMID[m.ID] = m.MaID
	}
	pMID = make(map[int64]*model.Module, len(ms))
	for _, m := range ms {
		if _, ok := mapMID[m.ID]; m.ID > 0 && !ok {
			err = ecode.EsportsActModNot
			return
		}
		pMID[m.ID] = m
		if m.ID == 0 {
			addM = append(addM, m)
		}
	}
	for _, m := range tmpMs {
		if mod, ok := pMID[m.ID]; ok {
			upM = append(upM, mod)
		} else {
			delM = append(delM, m.ID)
		}
	}
	// save
	tx := s.dao.DB.Begin()
	upFields := map[string]interface{}{"sid": param.Sid, "mid": param.Mid, "background": param.Background,
		"back_color": param.BackColor, "color_step": param.ColorStep, "live_id": param.LiveID, "intr": param.Intr,
		"focus": param.Focus, "url": param.URL, "status": param.Status,
		"h5_background": param.H5Background, "h5_back_color": param.H5BackColor,
		"h5_focus": param.H5Focus, "h5_url": param.H5URL,
		"intr_logo": param.IntrLogo, "intr_title": param.IntrTitle, "intr_text": param.IntrText}
	if err = tx.Model(&model.MatchActive{}).Where("id = ?", param.ID).Update(upFields).Error; err != nil {
		log.Error("EditAct MatchActive tx.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	if len(upM) > 0 {
		if err = tx.Model(&model.Module{}).Exec(model.BatchEditModuleSQL(upM)).Error; err != nil {
			log.Error("EditAct Module tx.Model Exec(%+v) error(%v)", upM, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(delM) > 0 {
		if err = tx.Model(&model.Module{}).Where("id IN (?)", delM).Updates(map[string]interface{}{"status": _deleted}).Error; err != nil {
			log.Error("EditAct Module tx.Model Updates(%+v) error(%v)", delM, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(addM) > 0 {
		if err = tx.Model(&model.Module{}).Exec(model.BatchAddModuleSQL(param.ID, addM)).Error; err != nil {
			log.Error("EditAct Module tx.Model Create(%+v) error(%v)", addM, err)
			err = tx.Rollback().Error
			return
		}
	}
	err = tx.Commit().Error
	return
}

func (s *Service) checkArc(c context.Context, ms []*model.Module) (rsAids map[string][]int64, err error) {
	var (
		name       string
		aids       []int64
		allAids    []int64
		tmpMap     map[int64]struct{}
		repeatAids []int64
		wrongAids  []int64
		isWrong    bool
	)
	rsAids = make(map[string][]int64, 2)
	for _, m := range ms {
		// check name only .
		if m.Name != "" && name == m.Name {
			err = ecode.EsportsModNameErr
			return
		}
		name = m.Name
		if aids, err = xstr.SplitInts(m.Oids); err != nil {
			err = ecode.RequestErr
			return
		}
		tmpMap = make(map[int64]struct{})
		for _, aid := range aids {
			if _, ok := tmpMap[aid]; ok {
				repeatAids = append(repeatAids, aid)
				continue
			}
			tmpMap[aid] = struct{}{}
		}
		allAids = append(allAids, aids...)
	}
	// check aids .
	if wrongAids, err = s.wrongArc(c, allAids); err != nil {
		err = ecode.EsportsArcServerErr
		return
	}
	if len(repeatAids) > 0 {
		rsAids["repeat"] = repeatAids
		isWrong = true
	}
	if len(wrongAids) > 0 {
		rsAids["wrong"] = wrongAids
		isWrong = true
	}
	if isWrong {
		err = ecode.EsportsModArcErr
	}
	return
}

func (s *Service) wrongArc(c context.Context, aids []int64) (list []int64, err error) {
	var (
		arcErr    error
		arcNormal map[int64]struct{}
		mutex     = sync.Mutex{}
	)
	group, errCtx := errgroup.WithContext(c)
	aidsLen := len(aids)
	arcNormal = make(map[int64]struct{}, aidsLen)
	for i := 0; i < aidsLen; i += _arcsSize {
		var partAids []int64
		if i+_arcsSize > aidsLen {
			partAids = aids[i:]
		} else {
			partAids = aids[i : i+_arcsSize]
		}
		group.Go(func() (err error) {
			var tmpRes *arcmdl.ArcsReply
			if tmpRes, arcErr = s.arcClient.Arcs(errCtx, &arcmdl.ArcsRequest{Aids: partAids}); arcErr != nil {
				log.Error("wrongArc s.arcClient.Arcs(%v) error %v", partAids, err)
				return arcErr
			}
			if tmpRes != nil {
				for _, arc := range tmpRes.Arcs {
					if arc != nil && arc.IsNormal() {
						mutex.Lock()
						arcNormal[arc.Aid] = struct{}{}
						mutex.Unlock()
					}
				}
			}
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		return
	}
	for _, aid := range aids {
		if _, ok := arcNormal[aid]; !ok {
			list = append(list, aid)
		}
	}
	return
}

// ForbidAct .
func (s *Service) ForbidAct(c context.Context, id int64, state int) (err error) {
	if err = s.dao.DB.Model(&model.MatchActive{}).Where("id=?", id).Updates(map[string]interface{}{"status": state}).Error; err != nil {
		log.Error("ForbidAct MatchActive s.dao.DB.Model Updates(%d) error(%v)", id, err)
	}
	return
}

// ListAct .
func (s *Service) ListAct(c context.Context, mid, pn, ps int64) (rs []*model.MatchModule, count int64, err error) {
	var (
		mas                                          []*model.MatchActive
		maIDs, matchIDs, seasonIDs                   []int64
		mapMs                                        map[int64][]*model.Module
		mapMatch                                     map[int64]*model.Match
		mapSeaon                                     map[int64]*model.Season
		matchTitle, matchSub, seasonTitle, seasonSub string
	)
	maDB := s.dao.DB.Model(&model.MatchActive{})
	if mid > 0 {
		maDB = maDB.Where("mid=?", mid)
	}
	maDB.Count(&count)
	if count == 0 {
		rs = _emptyActive
	}
	if err = maDB.Offset((pn - 1) * ps).Order("id ASC").Limit(ps).Find(&mas).Error; err != nil {
		log.Error("ListAct MatchActive s.dao.DB.Model Find error(%v)", err)
		return
	}
	if len(mas) == 0 {
		rs = _emptMatchMod
		return
	}
	for _, ma := range mas {
		maIDs = append(maIDs, ma.ID)
		matchIDs = append(matchIDs, ma.Mid)
		seasonIDs = append(seasonIDs, ma.Sid)
	}
	if ids := unique(matchIDs); len(ids) > 0 {
		var matchs []*model.Match
		if err = s.dao.DB.Model(&model.Match{}).Where("id IN (?)", ids).Find(&matchs).Error; err != nil {
			log.Error("ListAct match Error (%v)", err)
			return
		}
		mapMatch = make(map[int64]*model.Match, len(matchs))
		for _, v := range matchs {
			mapMatch[v.ID] = v
		}
	}
	if ids := unique(seasonIDs); len(ids) > 0 {
		var seasons []*model.Season
		if err = s.dao.DB.Model(&model.Match{}).Where("id IN (?)", ids).Find(&seasons).Error; err != nil {
			log.Error("ListAct season Error (%v)", err)
			return
		}
		mapSeaon = make(map[int64]*model.Season, len(seasonIDs))
		for _, v := range seasons {
			mapSeaon[v.ID] = v
		}
	}
	if mapMs, err = s.modules(maIDs, count); err != nil {
		log.Error("ListAct s.modules maIDs(%+v) faild(%+v)", maIDs, err)
		return
	}
	for _, ma := range mas {
		if match, ok := mapMatch[ma.Mid]; ok {
			matchTitle = match.Title
			matchSub = match.SubTitle
		} else {
			matchTitle = ""
			matchSub = ""
		}
		if season, ok := mapSeaon[ma.Sid]; ok {
			seasonTitle = season.Title
			seasonSub = season.SubTitle
		} else {
			seasonTitle = ""
			seasonSub = ""
		}
		if rsMs, ok := mapMs[ma.ID]; ok {
			tmpMs := rsMs
			rs = append(rs, &model.MatchModule{MatchActive: ma, Modules: tmpMs, MatchTitle: matchTitle, MatchSubTitle: matchSub, SeasonTitle: seasonTitle, SeasonSubTitle: seasonSub})
		} else {
			rs = append(rs, &model.MatchModule{MatchActive: ma, Modules: _emptyModules, MatchTitle: matchTitle, MatchSubTitle: matchSub, SeasonTitle: seasonTitle, SeasonSubTitle: seasonSub})
		}
	}
	return
}

func (s *Service) modules(maIDs []int64, count int64) (rs map[int64][]*model.Module, err error) {
	var ms []*model.Module
	if err = s.dao.DB.Model(&model.Module{}).Where("ma_id in(?)", maIDs).Where("status=?", _notDeleted).Find(&ms).Order("ma_id ASC").Error; err != nil {
		err = errors.Wrap(err, "modules map Model Find")
		return
	}
	rs = make(map[int64][]*model.Module, count)
	for _, m := range ms {
		tmpM := m
		rs[tmpM.MaID] = append(rs[tmpM.MaID], tmpM)
	}
	return
}
