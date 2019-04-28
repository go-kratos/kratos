package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/app/interface/openplatform/monitor-end/model/monitor"
	"go-common/library/ecode"
)

var (
	_duplicateNameErr    = errors.New("组名已经存在！")
	_duplicateTargetErr  = errors.New("已经存在该告警sub_event!")
	_invalidTargetIDErr  = errors.New("无效的id！")
	_duplicateProductErr = errors.New("重复的product name")
)

func targetKey(t *model.Target) string {
	if t.Source == "" || t.Product == "" || t.Event == "" || t.SubEvent == "" {
		return ""
	}
	return fmt.Sprintf("%s_%s_%s_%s", t.Source, t.Product, t.Event, t.SubEvent)

}

func productKey(p string) string {
	if p == "" {
		return ""
	}
	return fmt.Sprintf("%s", p)
}

func (s *Service) loadalertsettings() {
	var (
		c   = context.Background()
		gm  = make(map[int64]*model.Group)
		tm  = make(map[int64]*model.Target)
		tmk = make(map[string]*model.Target)
		tmn = make(map[string]*model.Target)
		pm  = make(map[int64]*model.Product)
		pmk = make(map[string]*model.Product)
	)
	if gs, err := s.dao.AllGroups(c); err == nil {
		for _, g := range gs {
			gm[g.ID] = g
		}
	}

	if ts, err := s.dao.AllTargets(c, 1); err == nil {
		for _, t := range ts {
			tm[t.ID] = t
			if key := targetKey(t); key != "" {
				tmk[key] = t
			}
		}
	}
	if ts, err := s.dao.AllTargets(c, 0); err == nil {
		for _, t := range ts {
			if key := targetKey(t); key != "" {
				tmn[key] = t
			}
		}
	}
	if ps, err := s.dao.AllProducts(c); err == nil {
		for _, p := range ps {
			pm[p.ID] = p
			if key := productKey(p.Name); err == nil {
				pmk[key] = p
			}
		}
	}
	s.mapMutex.Lock()
	s.groups = gm
	s.targets = tm
	s.products = pm
	s.targetKeys = tmk
	s.productKeys = pmk
	s.newTargets = tmn
	s.mapMutex.Unlock()
}

// AddGroup add a new group.
func (s *Service) AddGroup(c context.Context, group *model.Group) (id int64, err error) {
	var g *model.Group
	if group.Name == "" || group.Receivers == "" || group.Interval < 0 {
		err = ecode.RequestErr
		return
	}
	if g, err = s.GroupByName(c, group.Name); err != nil {
		return
	}
	if g.ID > 0 {
		err = _duplicateNameErr
		return
	}
	if group.Interval == 0 {
		group.Interval = 30
	}
	return s.dao.AddGroup(c, group)
}

// UpdateGroup update group.
func (s *Service) UpdateGroup(c context.Context, group *model.Group) (err error) {
	var g *model.Group
	if group.ID == 0 || group.Name == "" || group.Receivers == "" || group.Interval < 0 {
		err = ecode.RequestErr
		return
	}
	if g, err = s.GroupByName(c, group.Name); err != nil {
		return
	}
	if g.ID != 0 && g.ID != group.ID {
		err = _duplicateNameErr
		return
	}
	if group.Interval == 0 {
		group.Interval = 30
	}
	_, err = s.dao.UpdateGroup(c, group)
	return
}

// DeleteGroup delete group.
func (s *Service) DeleteGroup(c context.Context, id int64) (err error) {
	if id == 0 {
		err = ecode.RequestErr
		return
	}
	_, err = s.dao.DeleteGroup(c, id)
	return
}

// GroupList return all groups.
func (s *Service) GroupList(c context.Context, params *model.GroupListParams) (res *model.Groups, err error) {
	res = &model.Groups{}
	if res.Groups, err = s.dao.AllGroups(c); err != nil {
		return
	}
	res.Total = len(res.Groups)
	return
}

// GroupByName get group by name.
func (s *Service) GroupByName(c context.Context, name string) (res *model.Group, err error) {
	return s.dao.GroupByName(c, name)
}

// Target get target by id.
func (s *Service) Target(c context.Context, id int64) (res *model.Target, err error) {
	return s.dao.Target(c, id)
}

// AddTarget add a new target.
func (s *Service) AddTarget(c context.Context, t *model.Target) (id int64, err error) {
	if t.SubEvent == "" || t.Event == "" || t.Product == "" || t.Source == "" {
		err = ecode.RequestErr
		return
	}
	if err = s.checkTarget(c, t, true); err != nil {
		return
	}
	return s.dao.AddTarget(c, t)
}

// UpdateTarget update target.
func (s *Service) UpdateTarget(c context.Context, t *model.Target) (err error) {
	var (
		oldTarget *model.Target
	)
	if oldTarget, err = s.Target(c, t.ID); err != nil {
		return
	}
	if oldTarget == nil {
		err = _invalidTargetIDErr
		return
	}
	mergeTarget(t, oldTarget)
	if err = s.checkTarget(c, t, false); err != nil {
		return
	}
	_, err = s.dao.UpdateTarget(c, t)
	return
}

// TargetList .
func (s *Service) TargetList(c context.Context, t *model.Target, pn int, ps int, sort string) (res *model.Targets, err error) {
	// query := "SELECT id, sub_event, event, product, source, group_id, threshold, duration, state FROM target"
	var (
		where      = ""
		order      = ""
		empty      struct{}
		sortFields = map[string]struct{}{
			"sub_event": empty,
			"mtime":     empty,
			"ctime":     empty,
			"state":     empty,
		}
		sortOrder = map[string]string{
			"0": "DESC",
			"1": "ASC",
		}
	)
	if t.SubEvent != "" {
		where += " sub_event LIKE '%" + t.SubEvent + "%'"
	}
	if t.Event != "" {
		where += " event = '" + t.Event + "'"
	}
	if t.Product != "" {
		where += " product = '" + t.Product + "'"
	}
	if t.Source != "" {
		where += " source = '" + t.Source + "'"
	}
	if t.States != "" {
		where += " state in (" + t.States + ")"
	}
	if where == "" {
		where = " WHERE" + where + " deleted_time = 0"
	} else {
		where = " WHERE" + where + " AND deleted_time = 0"
	}
	countWhere := where
	pn = (pn - 1) * ps
	sorts := strings.Split(sort, ",")
	if len(sorts) == 2 {
		_, ok1 := sortFields[sorts[0]]
		d, ok2 := sortOrder[sorts[1]]
		if ok1 && ok2 {
			order = " ORDER BY " + sorts[0] + " " + d
		}
	}
	where += order
	where += fmt.Sprintf(" LIMIT %d, %d", pn, ps)
	res = &model.Targets{}
	if res.Targets, err = s.dao.TargetsByQuery(c, where); err != nil {
		return
	}
	if res.Total, err = s.dao.CountTargets(c, countWhere); err != nil {
		return
	}
	res.Page = pn
	res.PageSize = ps
	return
}

// TargetSync sync target state.
func (s *Service) TargetSync(c context.Context, id int64, state int) (err error) {
	return s.dao.TargetSync(c, id, state)
}

// DeleteTarget delete target by id.
func (s *Service) DeleteTarget(c context.Context, id int64) (err error) {
	_, err = s.dao.DeleteTarget(c, id)
	return
}

func (s *Service) checkTarget(c context.Context, t *model.Target, isNew bool) (err error) {
	var id int64
	t.SubEvent = induceSubEvent(t.SubEvent)
	if id, err = s.dao.IsExisted(c, t); err != nil {
		return
	}
	if isNew && id != 0 {
		fmt.Println("id", id)
		err = _duplicateTargetErr
	}
	if !isNew && id != 0 && t.ID != id {
		err = _duplicateTargetErr
	}
	return
}

func mergeTarget(t *model.Target, o *model.Target) {
	te := reflect.ValueOf(t).Elem()
	oe := reflect.ValueOf(o).Elem()
	for i := 0; i < te.NumField()-2; i++ {
		switch v := te.Field(i).Interface().(type) {
		case int, int64:
			if v == 0 {
				te.Field(i).Set(oe.Field(i))
			}
		case string:
			if v == "" {
				te.Field(i).Set(oe.Field(i))
			}
		}
	}
}

// AddProduct add a new group.
func (s *Service) AddProduct(c context.Context, p *model.Product) (id int64, err error) {
	if p.Name == "" || p.GroupIDs == "" {
		err = ecode.RequestErr
		return
	}
	var a *model.Product
	if a, err = s.dao.ProductByName(c, p.Name); err != nil {
		return
	}
	if a != nil {
		err = _duplicateProductErr
		return
	}
	id, err = s.dao.AddProduct(c, p)
	return
}

// UpdateProduct update product.
func (s *Service) UpdateProduct(c context.Context, p *model.Product) (err error) {
	if p.ID == 0 || p.Name == "" || p.GroupIDs == "" {
		err = ecode.RequestErr
		return
	}
	var a *model.Product
	if a, err = s.dao.ProductByName(c, p.Name); err != nil {
		return
	}
	if a != nil && a.ID != p.ID {
		err = _duplicateProductErr
		return
	}
	_, err = s.dao.UpdateProduct(c, p)
	return
}

// DeleteProduct delete product.
func (s *Service) DeleteProduct(c context.Context, id int64) (err error) {
	_, err = s.dao.DeleteProduct(c, id)
	return
}

// AllProducts return all products.
func (s *Service) AllProducts(c context.Context) (res *model.Products, err error) {
	res = &model.Products{}
	if res.Products, err = s.dao.AllProducts(c); err != nil {
		return
	}
	res.Total = len(res.Products)
	return
}

// Collect collect.
func (s *Service) Collect(c context.Context, p *monitor.Log) {
	var (
		curr int
		key  string
		t    *model.Target
	)
	target := &model.Target{
		Source:   sourceFromLog(p),
		Product:  p.Product,
		Event:    p.Event,
		SubEvent: induceSubEvent(p.SubEvent),
		State:    0,
	}
	s.infoCh <- p
	//TODO 获取buvid, ip等信息，过滤重复请求
	if key = targetKey(target); key == "" {
		return
	}

	if t = s.targetKeys[key]; t == nil {
		if s.newTargets[key] == nil {
			// 添加新的target
			s.AddTarget(c, target)
		}
		s.mapMutex.Lock()
		s.newTargets[key] = t
		s.mapMutex.Unlock()
		return
	}
	if t.Threshold == 0 {
		return
	}
	p.CalCode()
	code := codeFromLog(p)
	curr = s.dao.TargetIncr(c, t, code)
	if curr > t.Threshold {
		go s.mail(c, p, t, curr, code)
	}
}

func sourceFromLog(l *monitor.Log) string {
	if l.Type == "web/h5" {
		return l.Type
	}
	if l.RequestURI == "" {
		return "app"
	}
	if res := strings.Split(l.RequestURI, "?"); len(res) > 1 {
		return res[1]
	}
	return l.RequestURI
}

func codeFromLog(l *monitor.Log) string {
	if l.Codes == "" {
		return "999"
	}
	if l.HTTPCode != "" && l.HTTPCode != "200" {
		return l.HTTPCode
	}
	if l.BusinessCode != "" && l.BusinessCode != "0" {
		return l.BusinessCode
	}
	return "-999"
}
