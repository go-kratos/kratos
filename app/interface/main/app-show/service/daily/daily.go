package daily

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	arcdao "go-common/app/interface/main/app-show/dao/archive"
	carddao "go-common/app/interface/main/app-show/dao/card"
	tagdao "go-common/app/interface/main/app-show/dao/tag"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/card"
	"go-common/app/interface/main/app-show/model/daily"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
)

const (
	_initDailyKey      = "daily_key_%d_%d"
	_initColumnKey     = "column_key_%d_%d"
	_initColumnListKey = "columnlist_key_%d_%d"
)

var (
	_emptyDaily = []*daily.Show{}
	_emptyList  = []*daily.Item{}
)

type Service struct {
	c    *conf.Config
	cdao *carddao.Dao
	arc  *arcdao.Dao
	tag  *tagdao.Dao
	// tick
	tick time.Duration
	// columnsCache
	columnsCache map[string]*card.Column
	// card
	cardCache       map[string][]*daily.Show
	columnCache     map[string]*daily.Show
	columnListCache map[string][]*daily.Item
	cardListCache   map[string][]*card.ColumnList
}

// New new a daily service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		cdao: carddao.New(c),
		arc:  arcdao.New(c),
		tag:  tagdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// columnsCache
		columnsCache: map[string]*card.Column{},
		// card
		cardCache:       map[string][]*daily.Show{},
		columnCache:     map[string]*daily.Show{},
		columnListCache: map[string][]*daily.Item{},
		cardListCache:   map[string][]*card.ColumnList{},
	}
	now := time.Now()
	s.loadColumnListCache(now)
	s.loadColumnsCache()
	s.loadNperCache(now)
	go s.cacheproc()
	return
}

// Daily
func (s *Service) Daily(c context.Context, plat int8, build, dailyID, pn, ps int) (res []*daily.Show) {
	if pn > 0 {
		pn = pn - 1
	}
	start := pn * ps
	end := start + ps
	key := fmt.Sprintf(_initColumnKey, plat, dailyID)
	if column, ok := s.columnsCache[key]; ok {
		if model.InvalidBuild(build, column.Build, column.Condition) {
			res = _emptyDaily
			return
		}
		cardKey := fmt.Sprintf(_initDailyKey, plat, dailyID)
		if cards, ok := s.cardCache[cardKey]; ok {
			for _, sw := range cards {
				if model.InvalidBuild(build, sw.Build, sw.Condition) {
					continue
				}
				res = append(res, sw)
			}
			resLen := len(res)
			if resLen > end {
				res = res[start:end]
			} else if resLen > start {
				res = res[start:]
			} else {
				res = _emptyDaily
			}
		}
	}
	if len(res) == 0 {
		res = _emptyDaily
	}
	return
}

// ColumnList
func (s *Service) ColumnList(plat int8, build, columnID int) (res *daily.ColumnList) {
	var (
		column []*daily.ColumnList
	)
	key := fmt.Sprintf(_initColumnListKey, plat, columnID)
	if columns, ok := s.cardListCache[key]; ok {
		for _, c := range columns {
			if model.InvalidBuild(build, c.Build, c.Condition) {
				continue
			}
			tmp := &daily.ColumnList{
				Cid:   c.Cid,
				Name:  c.Name,
				Ceid:  c.Ceid,
				Cname: c.Cname,
			}
			column = append(column, tmp)
		}
		if len(column) > 0 {
			res = &daily.ColumnList{
				Ceid:     column[0].Ceid,
				Name:     column[0].Cname,
				Children: column,
			}
		}
	}
	return
}

// Category
func (s *Service) Category(plat int8, build, categoryID, columnID, pn, ps int) (res *daily.Show) {
	var (
		key string
	)
	if pn > 0 {
		pn = pn - 1
	}
	start := pn * ps
	end := start + ps
	if columnID > 0 {
		key = fmt.Sprintf(_initDailyKey, plat, columnID)
	} else {
		listKey := fmt.Sprintf(_initColumnListKey, plat, categoryID)
		if columns, ok := s.cardListCache[listKey]; ok {
			for _, c := range columns {
				if model.InvalidBuild(build, c.Build, c.Condition) {
					continue
				}
				key = fmt.Sprintf(_initDailyKey, plat, c.Cid)
				break
			}
		}
	}
	if columns, ok := s.columnCache[key]; ok {
		res = columns
		if pn*ps > 400 {
			res.Body = _emptyList
			return
		}
		if res.Body, ok = s.columnListCache[key]; ok {
			resLen := len(res.Body)
			if resLen > end {
				res.Body = res.Body[start:end]
			} else if resLen > start {
				res.Body = res.Body[start:]
			} else {
				res.Body = _emptyList
			}
		}
		if len(res.Body) == 0 {
			res.Body = _emptyList
		}
	}
	return
}

// loadColumnsCache load all columns cache
func (s *Service) loadColumnsCache() {
	res, err := s.cdao.Columns(context.TODO())
	if err != nil {
		log.Error("s.cdao.Columns error(%v)", err)
		return
	}
	tmp := map[string]*card.Column{}
	for plat, columns := range res {
		for _, column := range columns {
			key := fmt.Sprintf(_initColumnKey, plat, column.ID)
			tmp[key] = column
		}
	}
	s.columnsCache = tmp
}

// loadColumnListCache
func (s *Service) loadColumnListCache(now time.Time) {
	var (
		tmp = map[string][]*card.ColumnList{}
	)
	platColumns, err := s.cdao.ColumnPlatList(context.TODO(), now)
	if err != nil {
		log.Error("s.cdao.ColumnPlatList error(%v)", err)
		return
	}
	for plat, columns := range platColumns {
		for _, column := range columns {
			key := fmt.Sprintf(_initColumnListKey, plat, column.Ceid)
			tmp[key] = append(tmp[key], column)
		}
	}
	s.cardListCache = tmp
}

// loadNperCache
func (s *Service) loadNperCache(now time.Time) {
	hdm, err := s.cdao.ColumnNpers(context.TODO(), now)
	if err != nil {
		log.Error("s.cdao.ColumnNpers error(%v)", err)
		return
	}
	itm, aids, err := s.cdao.NperContents(context.TODO(), now)
	if err != nil {
		log.Error("s.cdao.NperContents error(%v)", err)
		return
	}
	tmp, tmpColumns, tmpList := s.mergeCard(context.TODO(), hdm, itm, aids, now)
	s.cardCache = tmp
	s.columnCache = tmpColumns
	s.columnListCache = tmpList
}

// cacheproc load all cache.
func (s *Service) cacheproc() {
	for {
		time.Sleep(s.tick)
		now := time.Now()
		s.loadColumnListCache(now)
		s.loadColumnsCache()
		s.loadNperCache(now)
	}
}

// mergeCard
func (s *Service) mergeCard(c context.Context, hdm map[int8][]*card.ColumnNper, itm map[int][]*card.Content, itmaids map[int][]int64, now time.Time) (res map[string][]*daily.Show, columns map[string]*daily.Show, columnList map[string][]*daily.Item) {
	var (
		dailyMAX = 31
	)
	res = map[string][]*daily.Show{}
	columnList = map[string][]*daily.Item{}
	columns = map[string]*daily.Show{}
	for plat, hds := range hdm {
		for _, hd := range hds {
			var (
				ok     bool
				column *card.Column
			)
			columnskey := fmt.Sprintf(_initColumnKey, plat, hd.ColumnID)
			if column, ok = s.columnsCache[columnskey]; !ok {
				continue
			}
			switch column.Type {
			case model.GotoDaily:
				if dailykey := fmt.Sprintf(_initDailyKey, plat, hd.ColumnID); len(res[dailykey]) > dailyMAX {
					continue
				}
			}
			var (
				sis []*daily.Item
			)
			its, ok := itm[hd.ID]
			if !ok {
				its = []*card.Content{}
			}
			switch column.Tpl {
			case 1, 2:
				var tmpItem = map[int64]*daily.Item{}
				if aids, ok := itmaids[hd.ID]; ok {
					tmpItem = s.fromCardAids(context.TODO(), aids)
				}
				for _, ci := range its {
					si := s.fillCardItem(ci, tmpItem)
					if si.Title == "" {
						continue
					}
					if ci.TagID > 0 {
						si.TagName, si.TagID = s.fromTagIDByName(c, ci.TagID, now)
					}
					sis = append(sis, si)
				}
			}
			if len(sis) == 0 {
				continue
			}
			sw := &daily.Show{}
			sw.Head = &daily.Head{
				ColumnID:  hd.ID,
				Build:     hd.Build,
				Condition: hd.Condition,
				Plat:      hd.Plat,
				Desc:      hd.Desc,
				Type:      column.Type,
			}
			if hd.Cover != "" {
				sw.Cover = hd.Cover
			}
			var key string
			switch sw.Head.Type {
			case model.GotoDaily:
				key = fmt.Sprintf(_initDailyKey, plat, hd.ColumnID)
				sw.Head.Title = hd.Name
				if len(res[key]) == 0 {
					sw.Head.Date = now.Unix()
				} else {
					sw.Head.Date = int64(hd.NperTime)
				}
				sw.Body = sis
				res[key] = append(res[key], sw)
			case model.GotoColumn:
				key = fmt.Sprintf(_initDailyKey, plat, hd.ID)
				sw.Head.Title = hd.Name
				sw.Head.Goto = hd.Goto
				sw.Head.Param = hd.Param
				sw.Head.URI = hd.URI
				columnList[key] = sis
				columns[key] = sw
			}
		}
	}
	return
}

// fillCardItem
func (s *Service) fillCardItem(csi *card.Content, tsi map[int64]*daily.Item) (si *daily.Item) {
	si = &daily.Item{}
	switch csi.Type {
	case model.CardGotoAv:
		si.Goto = model.GotoAv
		si.Param = csi.Value
	}
	si.URI = model.FillURI(si.Goto, si.Param, nil)
	if si.Goto == model.GotoAv {
		aid, err := strconv.ParseInt(si.Param, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", si.Param, err)
		} else {
			if it, ok := tsi[aid]; ok {
				si = it
				if csi.Title != "" {
					si.Title = csi.Title
				}
			} else {
				si = &daily.Item{}
			}
		}
	}
	return
}

// fromCardAids get Aids.
func (s *Service) fromCardAids(ctx context.Context, aids []int64) (data map[int64]*daily.Item) {
	var (
		arc *api.Arc
		ok  bool
	)
	as, err := s.arc.ArchivesPB(ctx, aids)
	if err != nil {
		log.Error("s.arc.ArchivesPB(%v) error(%v)", aids, err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB(%v) length is 0", aids)
		return
	}
	data = map[int64]*daily.Item{}
	for _, aid := range aids {
		if arc, ok = as[aid]; ok {
			if !arc.IsNormal() {
				continue
			}
			i := &daily.Item{}
			i.FromArchivePB(arc)
			data[aid] = i
		}
	}
	return
}

// fromTagIDByName from tag_id by tag_name
func (s *Service) fromTagIDByName(ctx context.Context, tagID int, now time.Time) (tagName string, tagIDInt int64) {
	tag, err := s.tag.TagInfo(ctx, 0, tagID, now)
	if err != nil {
		log.Error("s.tag.TagInfo(%d) error(%v)", tagID, err)
		return
	}
	tagName = tag.Name
	tagIDInt = tag.Tid
	return
}
