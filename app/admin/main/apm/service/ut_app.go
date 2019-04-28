package service

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"
)

var (
	_upsertUtApp = "INSERT INTO ut_app (path,owner) VALUES %s"
)

// AppsCache flush cache for apps.
func (s *Service) AppsCache(c context.Context) (err error) {
	var (
		appSlice []*ut.App
		appMap   = make(map[string]*ut.App)
		ownerMap = make(map[string][]*ut.App)
		deptMap  = make(map[string]*ut.Department)
	)
	if err = s.DB.Table("ut_app").Find(&appSlice).Error; err != nil {
		log.Error("s.AppsCache.Find() error(%v)", err)
		return
	}
	for _, app := range appSlice {
		appMap[app.Path] = app
		owners := strings.Split(app.Owner, ",")
		for _, owner := range owners {
			ownerMap[owner] = append(ownerMap[owner], app)
		}
		pathSlice := strings.Split(app.Path, "/")
		if len(pathSlice) < 2 {
			continue
		}
		deptName := pathSlice[len(pathSlice)-2]
		if deptMap[deptName] == nil {
			deptMap[deptName] = &ut.Department{}
		}
		deptMap[deptName].Name = deptName
		deptMap[deptName].Total++
		if app.HasUt == 1 {
			deptMap[deptName].Access++
			deptMap[deptName].Coverage += app.Coverage
		}
	}
	for _, v := range deptMap {
		if v.Access > 0 {
			v.Coverage = v.Coverage / float64(v.Access)
		}
	}
	s.appsCache.Lock()
	s.appsCache.Map = appMap
	s.appsCache.Slice = appSlice
	s.appsCache.Owner = ownerMap
	s.appsCache.Dept = deptMap
	s.appsCache.Unlock()
	return
}

// AddUTApp  add path to ut_app
func (s *Service) AddUTApp(c context.Context, apps []*ut.App) (err error) {
	var (
		valueStrings []string
		valueArgs    []interface{}
	)
	s.appsCache.Lock()
	for _, app := range apps {
		var cache, ok = s.appsCache.Map[app.Path]
		if !ok || cache == nil {
			valueStrings = append(valueStrings, "(?,?)")
			valueArgs = append(valueArgs, app.Path)
			valueArgs = append(valueArgs, app.Owner)
			continue
		}
		if cache.Owner == app.Owner {
			continue
		}
		if err = s.DB.Table("ut_app").Where("ID=?", cache.ID).
			Update("owner", app.Owner).Error; err != nil {
			log.Error("AddUTApp err (%v)", err)
			return
		}
		cache.Owner = app.Owner
	}
	s.appsCache.Unlock()
	if len(valueStrings) == 0 {
		return
	}
	stmt := fmt.Sprintf(_upsertUtApp, strings.Join(valueStrings, ","))
	if err = s.DB.Exec(stmt, valueArgs...).Error; err != nil {
		return
	}
	// update AppsCache
	if err = s.AppsCache(c); err != nil {
		return
	}
	return
}

// UpdateUTApp update has_ut=1
func (s *Service) UpdateUTApp(c context.Context, pkg *ut.PkgAnls) (err error) {
	s.appsCache.Lock()
	defer s.appsCache.Unlock()
	path := paserPkg(pkg.PKG)
	app, ok := s.appsCache.Map[path]
	if !ok || (app.HasUt != 0 && app.Coverage == pkg.Coverage) {
		log.Info("s.UpdateUTApp(%s) skiped.", pkg.PKG)
		return
	}
	app.HasUt = 1
	app.Coverage = pkg.Coverage
	if err = s.DB.Table("ut_app").Where("ID=?", app.ID).Updates(app).Error; err != nil {
		log.Error("UpdateUTApp err (%v)", err)
		return
	}
	return
}

func paserPkg(pkg string) (path string) {
	temp := strings.Split(pkg, "/")
	if len(temp) < int(5) {
		path = pkg
		return
	}
	path = strings.Join(temp[0:5], "/")
	return
}

// UTApps .
func (s *Service) UTApps(c context.Context, v *ut.AppReq) (result []*ut.App, count int, err error) {
	if v.Path != "" {
		if err = s.DB.Table("ut_app").Where("path LIKE ?", "%"+v.Path+"%").
			Count(&count).Find(&result).Error; err != nil {
			log.Error("UtProject err (%v)", err)
			return
		}
	} else {
		if err = s.DB.Table("ut_app").Where("has_ut=?", v.HasUt).Count(&count).Offset((v.Pn - 1) * v.Ps).
			Limit(v.Ps).Find(&result).Error; err != nil {
			log.Error("UtProject err (%v)", err)
			return
		}
	}
	for _, v := range result {
		v.Link = parsePath(v.Path)
	}
	return
}

// "go-common/app/service/main/share" to  "go-common/tree/master/app/service/main/share"
func parsePath(path string) (link string) {
	temp := strings.SplitN(path, "/", 2)
	link = fmt.Sprintf("%s%s%s", temp[0], "/tree/master/", temp[1])
	return
}
