package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	showModel "go-common/app/admin/main/feed/model/show"
	"go-common/app/admin/main/feed/util"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	//_ActAddSearchWebCard log action
	_ActAddSearchWebCard = "ActAddSearchWebCard"
	//_ActUpSearchWebCard log action
	_ActUpSearchWebCard = "ActUpSearchWebCard"
	//_ActDelSearchWebCard log action
	_ActDelSearchWebCard = "ActDelSearchWebCard"
	//_ActAddSearchWeb log action
	_ActAddSearchWeb = "ActAddSearchWeb"
	//_ActUpSearchWeb log action
	_ActUpSearchWeb = "ActUpSearchWeb"
	//_ActDelSearchWeb log action
	_ActDelSearchWeb = "ActDelSearchWeb"
	//_ActOptSearchWeb log action
	_ActOptSearchWeb = "ActOptSearchWeb"
)

var (
	_emptyWebQuery = make([]*show.SearchWebQuery, 0)
)

//SearchWebCardList channel SearchWebCard list
func (s *Service) SearchWebCardList(lp *show.SearchWebCardLP) (pager *show.SearchWebCardPager, err error) {
	pager = &show.SearchWebCardPager{
		Page: common.Page{
			Num:  lp.Pn,
			Size: lp.Ps,
		},
	}
	w := map[string]interface{}{
		"deleted": common.NotDeleted,
	}
	query := s.showDao.DB.Model(&show.SearchWebCard{})
	if lp.ID > 0 {
		w["id"] = lp.ID
	}
	if lp.Person != "" {
		query = query.Where("person like ?", "%"+lp.Person+"%")
	}
	if lp.Title != "" {
		query = query.Where("title like ?", "%"+lp.Title+"%")
	}
	if lp.STime != "" {
		query = query.Where("ctime >= ?", lp.STime)
	}
	if lp.ETime != "" {
		query = query.Where("ctime <= ?", lp.ETime)
	}
	if err = query.Where(w).Count(&pager.Page.Total).Error; err != nil {
		log.Error("searchWebSvc.SearchWebCardList count error(%v)", err)
		return
	}
	SearchWebCards := make([]*show.SearchWebCard, 0)
	if err = query.Where(w).Order("`id` DESC").Offset((lp.Pn - 1) * lp.Ps).Limit(lp.Ps).Find(&SearchWebCards).Error; err != nil {
		log.Error("searchWebSvc.SearchWebCardList Find error(%v)", err)
		return
	}
	pager.Item = SearchWebCards
	return
}

//AddSearchWebCard add channel SearchWebCard
func (s *Service) AddSearchWebCard(c context.Context, param *show.SearchWebCardAP, name string, uid int64) (err error) {
	if err = s.showDao.SearchWebCardAdd(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogSWEBCard, name, uid, 0, _ActAddSearchWebCard, param); err != nil {
		log.Error("searchWebSvc.AddSearchWebCard AddLog error(%v)", err)
		return
	}
	return
}

//UpdateSearchWebCard update channel SearchWebCard
func (s *Service) UpdateSearchWebCard(c context.Context, param *show.SearchWebCardUP, name string, uid int64) (err error) {
	if err = s.showDao.SearchWebCardUpdate(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogSWEBCard, name, uid, 0, _ActUpSearchWebCard, param); err != nil {
		log.Error("searchWebSvc.UpdateSearchWebCard AddLog error(%v)", err)
		return
	}
	return
}

//DeleteSearchWebCard delete channel SearchWebCard
func (s *Service) DeleteSearchWebCard(id int64, name string, uid int64) (err error) {
	if err = s.showDao.SearchWebCardDelete(id); err != nil {
		return
	}
	if err = util.AddLogs(common.LogSWEBCard, name, uid, id, _ActDelSearchWebCard, id); err != nil {
		log.Error("searchWebSvc.DeleteSearchWebCard AddLog error(%v)", err)
		return
	}
	return
}

//SearchWebList SearchWeb list
func (s *Service) SearchWebList(lp *show.SearchWebLP) (pager *show.SearchWebPager, err error) {
	pager = &show.SearchWebPager{
		Page: common.Page{
			Num:  lp.Pn,
			Size: lp.Ps,
		},
	}
	w := map[string]interface{}{
		"deleted": common.NotDeleted,
	}
	query := s.showDao.DB.Model(&show.SearchWeb{})
	if lp.ID > 0 {
		w["id"] = lp.ID
	}
	if lp.Person != "" {
		query = query.Where("person like ?", "%"+lp.Person+"%")
	}
	if lp.STime != "" {
		query = query.Where("stime >= ?", lp.STime)
	}
	if lp.ETime != "" {
		query = query.Where("etime <= ?", lp.ETime)
	}
	cTimeStr := util.CTimeStr()
	if lp.Check != 0 {
		if lp.Check == common.Pass {
			//已通过 未生效
			query = query.Where("`check` = ?", common.Pass)
			query = query.Where("stime > ?", cTimeStr)
		} else if lp.Check == common.Valid {
			//已通过 已生效
			query = query.Where("`check` = ?", common.Pass)
			query = query.Where("stime <= ?", cTimeStr).Where("etime >= ?", cTimeStr)
		} else if lp.Check == common.InValid {
			//已通过 已失效
			query = query.Where("(`check` = ? AND etime <= ?) OR (`check` = ?)", common.Pass, cTimeStr, common.InValid)
		} else {
			query = query.Where("`check` = ? ", lp.Check)
		}
	}
	if err = query.Where(w).Count(&pager.Page.Total).Error; err != nil {
		log.Error("searchSvc.SearchWebList count error(%v)", err)
		return
	}
	SearchWebs := make([]*show.SearchWeb, 0)
	if err = query.Where(w).Order("`id` DESC").Offset((lp.Pn - 1) * lp.Ps).Limit(lp.Ps).Find(&SearchWebs).Error; err != nil {
		log.Error("searchSvc.SearchWebList Find error(%v)", err)
		return
	}
	if len(SearchWebs) > 0 {
		var (
			ids      []int64
			queryMap map[int64][]*show.SearchWebQuery
		)
		for _, v := range SearchWebs {
			if v.Check == common.Pass {
				c := time.Now().Unix()
				if (c >= v.Stime.Time().Unix()) && (c <= v.Etime.Time().Unix()) {
					v.Check = common.Valid
				} else if c > v.Etime.Time().Unix() && v.Check != common.InValid {
					v.Check = common.InValid
					v.Status = common.StatusDownline
				}
			}
			webCard := &show.SearchWebCard{}
			cardWhere := map[string]interface{}{
				"deleted": common.NotDeleted,
				"id":      v.CardValue,
			}
			if err = s.showDao.DB.Model(&show.SearchWebCard{}).Where(cardWhere).First(webCard).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = nil
				} else {
					log.Error("searchSvc.SearchWebCard Find error(%v)", err)
					return
				}
			}
			v.Card = webCard
			ids = append(ids, v.ID)
		}
		where := map[string]interface{}{
			"deleted": common.NotDeleted,
		}
		SearchWebQuery := make([]*show.SearchWebQuery, 0)
		if err = s.showDao.DB.Model(&show.SearchWebQuery{}).Where(where).Where("sid in (?)", ids).Find(&SearchWebQuery).Error; err != nil {
			log.Error("searchSvc.SearchWebList Find error(%v)", err)
			return
		}
		queryMap = make(map[int64][]*show.SearchWebQuery, len(SearchWebQuery))
		for _, v := range SearchWebQuery {
			queryMap[v.SID] = append(queryMap[v.SID], v)
		}
		for _, v := range SearchWebs {
			if value, ok := queryMap[v.ID]; ok {
				v.Query = value
			} else {
				v.Query = _emptyWebQuery
			}
		}
	}
	pager.Item = SearchWebs
	return
}

//OpenSearchWebList SearchWeb list
func (s *Service) OpenSearchWebList() (SearchWebs []*show.SearchWeb, err error) {
	cTimeStr := util.CTimeStr()
	SearchWebs = make([]*show.SearchWeb, 0)
	w := map[string]interface{}{
		"deleted": common.NotDeleted,
		"check":   common.Pass,
	}
	query := s.showDao.DB.Model(&show.SearchWeb{})
	//已通过 已生效
	query = query.Where("stime <= ?", cTimeStr).Where("etime >= ?", cTimeStr)
	if err = query.Where(w).Order("`id` DESC").Find(&SearchWebs).Error; err != nil {
		log.Error("searchSvc.OpenSearchWebList Find error(%v)", err)
		return
	}
	if len(SearchWebs) > 0 {
		var (
			ids      []int64
			queryMap map[int64][]*show.SearchWebQuery
		)
		for _, v := range SearchWebs {
			webCard := &show.SearchWebCard{}
			cardWhere := map[string]interface{}{
				"deleted": common.NotDeleted,
				"id":      v.CardValue,
			}
			if err = s.showDao.DB.Model(&show.SearchWebCard{}).Where(cardWhere).First(webCard).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = nil
					webCard = nil
				} else {
					log.Error("searchSvc.OpenSearchWebList Find error(%v)", err)
				}
			}
			if webCard != nil {
				v.Card = webCard
			} else {
				v.Card = struct{}{}
			}
			ids = append(ids, v.ID)
		}
		where := map[string]interface{}{
			"deleted": common.NotDeleted,
		}
		SearchWebQuery := make([]*show.SearchWebQuery, 0)
		if err = s.showDao.DB.Model(&show.SearchWebQuery{}).Where(where).Where("sid in (?)", ids).Find(&SearchWebQuery).Error; err != nil {
			log.Error("searchSvc.OpenSearchWebList Find error(%v)", err)
			return
		}
		queryMap = make(map[int64][]*show.SearchWebQuery, len(SearchWebQuery))
		for _, v := range SearchWebQuery {
			queryMap[v.SID] = append(queryMap[v.SID], v)
		}
		for _, v := range SearchWebs {
			if value, ok := queryMap[v.ID]; ok {
				v.Query = value
			} else {
				v.Query = _emptyWebQuery
			}
		}
	}
	return
}

//Validate validate search web card
func (s *Service) Validate(p *show.SWTimeValid) (err error) {
	var (
		querys  []*show.SearchWebQuery
		webCard *showModel.SearchWebCard
		id      int64
	)
	if id, err = strconv.ParseInt(p.CardValue, 10, 64); err != nil {
		return
	}
	if webCard, err = s.showDao.SWBFindByID(id); err != nil {
		return err
	}
	if webCard == nil {
		return fmt.Errorf("无效web卡片ID(%d)", id)
	}
	if err = json.Unmarshal([]byte(p.Query), &querys); err != nil {
		log.Error("searchSvc.Validate json.Unmarshal(%v) error(%v)", p, err)
		return
	}
	if len(querys) == 0 {
		err = fmt.Errorf("query不能为空")
		return
	}
	for _, v := range querys {
		count := 0
		p.Query = v.Value
		if count, err = s.showDao.SWTimeValid(p); err != nil {
			return
		}
		if count > 0 {
			err = fmt.Errorf("相同query(%s)该位置已有运营卡片", v.Value)
		}
	}
	return
}

//AddSearchWeb add SearchWeb
func (s *Service) AddSearchWeb(c context.Context, param *show.SearchWebAP, name string, uid int64) (err error) {
	p := &show.SWTimeValid{
		Priority:  param.Priority,
		STime:     param.Stime,
		ETime:     param.Etime,
		Query:     param.Query,
		CardValue: param.CardValue,
	}
	if err = s.Validate(p); err != nil {
		return
	}
	if err = s.showDao.SearchWebAdd(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogSWEB, name, uid, 0, _ActAddSearchWeb, param); err != nil {
		log.Error("searchSvc.AddSearchWeb AddLog error(%v)", err)
		return
	}
	return
}

//UpdateSearchWeb update SearchWeb
func (s *Service) UpdateSearchWeb(c context.Context, param *show.SearchWebUP, name string, uid int64) (err error) {
	var (
		swValue *show.SearchWeb
	)
	p := &show.SWTimeValid{
		ID:        param.ID,
		Priority:  param.Priority,
		STime:     param.Stime,
		ETime:     param.Etime,
		Query:     param.Query,
		CardValue: param.CardValue,
	}
	if err = s.Validate(p); err != nil {
		return
	}
	if swValue, err = s.showDao.SWFindByID(param.ID); err != nil {
		log.Error("searchSvc.UpdateSearchWeb AddLog error(%v)", err)
		return
	}
	//待审核&已通过&已生效-》编辑-》状态不变；其它-》编辑-》审待核
	cTime := time.Now().Unix()
	if (swValue.Check == common.Verify) ||
		(swValue.Check == common.Pass && swValue.Stime.Time().Unix() > cTime ||
			(swValue.Check == common.Pass && (cTime > swValue.Stime.Time().Unix() && cTime <= swValue.Stime.Time().Unix()))) {
		param.Check = swValue.Check
		param.Status = swValue.Status
	} else {
		param.Check = common.Verify
		param.Status = common.StatusDownline
	}
	if err = s.showDao.SearchWebUpdate(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogSWEB, name, uid, 0, _ActUpSearchWeb, param); err != nil {
		log.Error("searchSvc.UpdateSearchWeb AddLog error(%v)", err)
		return
	}
	return
}

//DeleteSearchWeb delete SearchWeb
func (s *Service) DeleteSearchWeb(id int64, name string, uid int64) (err error) {
	if err = s.showDao.SearchWebDelete(id); err != nil {
		return
	}
	if err = util.AddLogs(common.LogSWEB, name, uid, id, _ActDelSearchWeb, id); err != nil {
		log.Error("searchSvc.DeleteSearchWeb AddLog error(%v)", err)
		return
	}
	return
}

//OptionSearchWeb option SearchWeb
func (s *Service) OptionSearchWeb(id int64, opt string, name string, uid int64) (err error) {
	up := &show.SearchWebOption{}
	if opt == common.OptionOnline {
		up.Status = common.StatusOnline
		up.Check = common.Pass
	} else if opt == common.OptionHidden {
		up.Status = common.StatusDownline
		up.Check = common.InValid
	} else if opt == common.OptionPass {
		up.Status = common.StatusOnline
		up.Check = common.Pass
	} else if opt == common.OptionReject {
		up.Status = common.StatusDownline
		up.Check = common.Rejecte
	} else {
		err = fmt.Errorf("参数不合法")
		return
	}
	up.ID = id
	if err = s.showDao.SearchWebOption(up); err != nil {
		return
	}
	logParam := map[string]interface{}{
		"id":  id,
		"opt": opt,
		"up":  up,
	}
	if err = util.AddLogs(common.LogSWEB, name, uid, id, _ActOptSearchWeb, logParam); err != nil {
		log.Error("searchSvc.OptionSearchWeb AddLog error(%v)", err)
		return
	}
	return
}
