package service

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"path/filepath"
	"strconv"
	"strings"

	"go-common/app/admin/main/config/model"
	"go-common/app/admin/main/config/pkg/lint"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

func lintConfig(filename, content string) error {
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	err := lint.Lint(ext, bytes.NewBufferString(content))
	if err != nil && err != lint.ErrLintNotExists {
		return ecode.Error(ecode.RequestErr, err.Error())
	}
	return nil
}

// CreateConf create config.
func (s *Service) CreateConf(conf *model.Config, treeID int64, env, zone string, skiplint bool) error {
	// lint config
	if !skiplint {
		if err := lintConfig(conf.Name, conf.Comment); err != nil {
			return err
		}
	}
	app, err := s.AppByTree(treeID, env, zone)
	if err != nil {
		return err
	}
	conf.AppID = app.ID
	if _, err := s.configIng(conf.Name, app.ID); err == nil { // judge config is configIng
		return ecode.TargetBlocked
	}
	return s.dao.DB.Create(conf).Error
}

// LintConfig lint config file
func (s *Service) LintConfig(filename, content string) ([]lint.LineErr, error) {
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	err := lint.Lint(ext, bytes.NewBufferString(content))
	if err == nil || err == lint.ErrLintNotExists {
		return make([]lint.LineErr, 0), nil
	}
	lintErr, ok := err.(lint.Error)
	if !ok {
		return nil, lintErr
	}
	return []lint.LineErr(lintErr), nil
}

// UpdateConfValue update config state.
func (s *Service) UpdateConfValue(conf *model.Config, skiplint bool) (err error) {
	// lint config
	if !skiplint {
		if err := lintConfig(conf.Name, conf.Comment); err != nil {
			return err
		}
	}
	var confDB *model.Config
	if confDB, err = s.Config(conf.ID); err != nil {
		return
	}
	if confDB.State == model.ConfigIng { //judge config is configIng.
		if conf.Mtime != confDB.Mtime {
			err = ecode.TargetBlocked
			return
		}
		conf.Mtime = 0
		err = s.dao.DB.Model(&model.Config{State: model.ConfigIng}).Updates(conf).Error
		return
	}
	if _, err = s.configIng(confDB.Name, confDB.AppID); err == nil {
		err = ecode.TargetBlocked
		return
	}
	if err == sql.ErrNoRows || err == ecode.NothingFound {
		conf.ID = 0
		conf.AppID = confDB.AppID
		conf.Name = confDB.Name
		if conf.From == 0 {
			conf.From = confDB.From
		}
		conf.Mtime = 0

		return s.dao.DB.Create(conf).Error
	}
	return
}

// UpdateConfState update config state.
func (s *Service) UpdateConfState(ID int64) (err error) {
	err = s.dao.DB.Model(&model.Config{ID: ID}).Update("state", model.ConfigEnd).Error
	return
}

// ConfigsByIDs get Config by IDs.
func (s *Service) ConfigsByIDs(ids []int64) (confs []*model.Config, err error) {
	if err = s.dao.DB.Select("id,app_id,name,`from`,state,mark,operator,ctime,mtime,is_delete").Where(ids).Find(&confs).Error; err != nil {
		log.Error("ConfigsByIDs(%v) error(%v)", ids, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// ConfigsByAppName get configs by app name.
func (s *Service) ConfigsByAppName(appName, env, zone string, treeID int64, state int8) (confs []*model.Config, err error) {
	var app *model.App
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		if err == ecode.NothingFound {
			err = s.CreateApp(appName, env, zone, treeID)
			return
		}
	}
	if state != 0 {
		err = s.dao.DB.Select("id,app_id,name,`from`,state,mark,operator,is_delete,ctime,mtime").Where("app_id = ? and state =?", app.ID, state).Order("id desc").Find(&confs).Error
	} else {
		err = s.dao.DB.Select("id,app_id,name,`from`,state,mark,operator,is_delete,ctime,mtime").Where("app_id = ? ", app.ID).Order("id desc").Find(&confs).Error
	}
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// ConfigsByAppID configs by app ID.
func (s *Service) ConfigsByAppID(appID int64, state int8) (confs []*model.Config, err error) {
	if state != 0 {
		err = s.dao.DB.Select("id,app_id,name,`from`,state,mark,operator,is_delete,ctime,mtime").Where("app_id = ? and state =?", appID, state).Order("id desc").Find(&confs).Error
	} else {
		err = s.dao.DB.Select("id,app_id,name,`from`,state,mark,operator,is_delete,ctime,mtime").Where("app_id = ? ", appID).Order("id desc").Find(&confs).Error
	}
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

//ConfigSearchApp configSearchApp.
func (s *Service) ConfigSearchApp(ctx context.Context, appName, env, zone, like string, buildID, treeID int64) (searchs []*model.ConfigSearch, err error) {
	var (
		app     *model.App
		builds  []*model.Build
		tags    []*model.Tag
		confs   []*model.Config
		tagIDs  []int64
		confIDs []int64
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	if builds, err = s.BuildsByApp(app.ID); err != nil {
		return
	}
	if len(builds) == 0 {
		return
	}
	for _, build := range builds {
		tagIDs = append(tagIDs, build.TagID)
	}
	if err = s.dao.DB.Where("id in(?)", tagIDs).Find(&tags).Error; err != nil {
		log.Error("tagsByIDs(%v) error(%v)", tagIDs, err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	tmp := make(map[int64]struct{})
	for _, tag := range tags {
		tmpIDs := strings.Split(tag.ConfigIDs, ",")
		for _, tmpID := range tmpIDs {
			var id int64
			if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
				log.Error("strconv.ParseInt() error(%v)", err)
				return
			}
			if _, ok := tmp[id]; !ok {
				tmp[id] = struct{}{}
				confIDs = append(confIDs, id)
			}
		}
	}
	if err = s.dao.DB.Where("id in (?) AND comment like(?) ", confIDs, "%"+like+"%").Find(&confs).Error; err != nil {
		log.Error("confsByIDs(%v) error(%v)", confIDs, err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	for _, conf := range confs {
		search := new(model.ConfigSearch)
		search.App = appName
		for _, tag := range tags {
			tmpIDs := strings.Split(tag.ConfigIDs, ",")
			for _, tmpID := range tmpIDs {
				var id int64
				if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
					log.Error("strconv.ParseInt() error(%v)", err)
					return
				}
				if id != conf.ID { //judge config is in build.
					continue
				}
				for _, build := range builds {
					if build.ID == tag.BuildID {
						search.Builds = append(search.Builds, build.Name)
					}
				}
			}
		}
		//generate comment.
		search.ConfValues = genComments(conf.Comment, like)
		search.ConfID = conf.ID
		search.Mark = conf.Mark
		search.ConfName = conf.Name
		searchs = append(searchs, search)
	}
	return
}

//ConfigSearchAll configSearchAll.
func (s *Service) ConfigSearchAll(ctx context.Context, env, zone, like string, nodes *model.CacheData) (searchs []*model.ConfigSearch, err error) {
	var (
		apps      []*model.App
		builds    []*model.Build
		tags      []*model.Tag
		confs     []*model.Config
		names     []string
		appIDs    []int64
		tagIDs    []int64
		configIDs []int64
	)
	searchs = make([]*model.ConfigSearch, 0)
	if len(nodes.Data) == 0 {
		return
	}
	for _, node := range nodes.Data {
		names = append(names, node.Path)
	}
	if err = s.dao.DB.Where("env =? and zone =?  and name in(?)", env, zone, names).Find(&apps).Error; err != nil {
		log.Error("AppList() find apps() error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	for _, app := range apps {
		appIDs = append(appIDs, app.ID)
	}
	if err = s.dao.DB.Where("app_id in(?) ", appIDs).Find(&builds).Error; err != nil {
		log.Error("BuildsByAppIDs(%v) error(%v)", appIDs, err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	for _, build := range builds {
		tagIDs = append(tagIDs, build.TagID)
	}
	if err = s.dao.DB.Where("id in(?)", tagIDs).Find(&tags).Error; err != nil {
		log.Error("tagsByIDs(%v) error(%v)", tagIDs, err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	tmp := make(map[int64]struct{})
	for _, tag := range tags {
		tmpIDs := strings.Split(tag.ConfigIDs, ",")
		for _, tmpID := range tmpIDs {
			var id int64
			if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
				log.Error("strconv.ParseInt() error(%v)", err)
				return
			}
			if _, ok := tmp[id]; !ok {
				tmp[id] = struct{}{}
				configIDs = append(configIDs, id)
			}
		}
	}
	if err = s.dao.DB.Where("id in (?) and comment like(?) ", configIDs, "%"+like+"%").Find(&confs).Error; err != nil {
		log.Error("confsByIDs(%v) error(%v)", configIDs, err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	if len(confs) == 0 {
		return
	}
	// convert confs to confSearch.
	for _, conf := range confs {
		search := new(model.ConfigSearch)
		for _, app := range apps {
			if app.ID == conf.AppID {
				search.App = app.Name
				search.TreeID = app.TreeID
			}
		}
		for _, tag := range tags {
			tmpIDs := strings.Split(tag.ConfigIDs, ",")
			for _, tmpID := range tmpIDs {
				var id int64
				if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
					log.Error("strconv.ParseInt() error(%v)", err)
					return
				}
				if id != conf.ID { //judge config is in build.
					continue
				}
				for _, build := range builds {
					if build.ID == tag.BuildID {
						search.Builds = append(search.Builds, build.Name)
					}
				}
			}
		}
		//generate comment.
		search.ConfValues = genComments(conf.Comment, like)
		search.ConfID = conf.ID
		search.Mark = conf.Mark
		search.ConfName = conf.Name
		searchs = append(searchs, search)
	}
	return
}

func genComments(comment, like string) (comments []string) {
	var (
		line   []byte
		l, cur []byte
		err    error
	)
	wbuf := new(bytes.Buffer)
	rbuf := bufio.NewReader(strings.NewReader(comment))
	for {
		l = line
		if line, _, err = rbuf.ReadLine(); err != nil {
			break
		}
		if bytes.Contains(line, []byte(like)) {
			cur = line
			wbuf.Write(l)
			wbuf.WriteString("\n")
			wbuf.Write(cur)
			wbuf.WriteString("\n")
			line, _, _ = rbuf.ReadLine()
			wbuf.Write(line)
			wbuf.WriteString("\n")
			comments = append(comments, wbuf.String())
			wbuf.Reset()
		}
	}
	return
}

//Configs configs.
func (s *Service) Configs(appName, env, zone string, buildID, treeID int64) (res *model.ConfigRes, err error) {
	var (
		allConfs   []*model.Config
		buildConfs []*model.Config
		lastConfs  []*model.Config
		app        *model.App
		build      *model.Build
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		if err == ecode.NothingFound {
			err = s.CreateApp(appName, env, zone, treeID)
			return
		}
	}
	if allConfs, err = s.ConfigsByAppID(app.ID, 0); err != nil {
		return
	}
	if buildID != 0 {
		if build, err = s.Build(buildID); err != nil {
			return
		}
		if build.AppID != app.ID {
			err = ecode.NothingFound
			return
		}
		tagID := build.TagID
		if tagID == 0 {
			return
		}
		if buildConfs, err = s.ConfigsByTagID(tagID); err != nil {
			return
		}
		if lastConfs, err = s.LastBuildConfigConfigs(build.AppID, buildID); err != nil {
			if err != ecode.NothingFound {
				return
			}
			err = nil
		}
	}
	tmpMap := make(map[string]struct{})
	res = new(model.ConfigRes)
	for _, conf := range allConfs {
		if _, ok := tmpMap[conf.Name]; ok {
			continue
		}
		//new common
		if conf.From > 0 {
			conf.NewCommon, _ = s.NewCommon(conf.From)
		}
		tmpMap[conf.Name] = struct{}{}
		var bool bool
		for _, bconf := range buildConfs {
			//new common
			if bconf.From > 0 {
				bconf.NewCommon, _ = s.NewCommon(bconf.From)
			}
			if bconf.Name == conf.Name {
				if bconf.ID != conf.ID {
					if conf.IsDelete != 1 {
						res.BuildNewFile = append(res.BuildNewFile, conf)
					}
				}
				bf := &model.BuildFile{Config: bconf}
				if bf.IsDelete == 0 {
					res.BuildFiles = append(res.BuildFiles, bf)
				}
				for _, lconf := range lastConfs {
					if lconf.Name == bconf.Name {
						if lconf.ID != bconf.ID {
							bf.LastConf = lconf
						}
					}
				}
				bool = true
				break
			}
		}
		if !bool {
			if conf.IsDelete != 1 {
				res.Files = append(res.Files, conf)
			}
			continue
		}
	}
	return
}

// NewCommon get new common id
func (s *Service) NewCommon(from int64) (new int64, err error) {
	common := &model.CommonConf{}
	newCommon := &model.CommonConf{}
	if err = s.dao.DB.First(&common, from).Error; err != nil {
		log.Error("NewCommon.First.from(%d) error(%v)", from, err)
		return
	}
	if err = s.dao.DB.Where("team_id = ? and name = ? and state = 2", common.TeamID, common.Name).Order("id desc").First(&newCommon).Error; err != nil {
		log.Error("NewCommon.Order.First.common(%v) error(%v)", common, err)
		return
	}
	new = newCommon.ID
	return
}

//ConfigRefs configRefs.
func (s *Service) ConfigRefs(appName, env, zone string, buildID, treeID int64) (res []*model.ConfigRefs, err error) {
	var (
		allConfs   []*model.Config
		buildConfs []*model.Config
		num        int
		ok         bool
		ref        *model.ConfigRefs
	)
	if allConfs, err = s.ConfigsByAppName(appName, env, zone, treeID, model.ConfigEnd); err != nil {
		return
	}
	if buildID != 0 {
		if buildConfs, err = s.ConfigsByBuildID(buildID); err != nil {
			return
		}
	}
	tmpMap := make(map[string]int)
	refs := make(map[string]*model.ConfigRefs)
	for _, conf := range allConfs {
		if num, ok = tmpMap[conf.Name]; !ok {
			ref = new(model.ConfigRefs)
			refs[conf.Name] = ref
			tmpMap[conf.Name] = num
		} else {
			ref = refs[conf.Name]
		}
		if num <= 5 {
			ref.Configs = append(ref.Configs, &model.ConfigRef{ID: conf.ID, Mark: conf.Mark})
			tmpMap[conf.Name] = num + 1
		}
		if ref.Ref != nil {
			continue
		}
		for _, bconf := range buildConfs {
			if bconf.Name == conf.Name {
				ref.Ref = &model.ConfigRef{ID: bconf.ID, Mark: bconf.Mark}
				break
			}
		}
	}
	for k, v := range refs {
		v.Name = k
		res = append(res, v)
	}
	return
}

// ConfigsByBuildID get configs by build ID.
func (s *Service) ConfigsByBuildID(buildID int64) (confs []*model.Config, err error) {
	var (
		build *model.Build
	)
	if build, err = s.Build(buildID); err != nil {
		return
	}
	tagID := build.TagID
	if tagID == 0 {
		return
	}
	return s.ConfigsByTagID(tagID)
}

// LastBuildConfigs get configs by build ID.
func (s *Service) LastBuildConfigs(appID, buildID int64) (confs []*model.Config, err error) {
	var (
		tag *model.Tag
		ids []int64
		id  int64
	)
	if tag, err = s.LastTagByAppBuild(appID, buildID); err != nil {
		return
	}
	confIDs := tag.ConfigIDs
	if len(confIDs) == 0 {
		return
	}
	tmpIDs := strings.Split(confIDs, ",")
	for _, tmpID := range tmpIDs {
		if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", tmpID, err)
			return
		}
		ids = append(ids, id)
	}
	return s.ConfigsByIDs(ids)
}

// LastBuildConfigConfigs get configs by build ID.
func (s *Service) LastBuildConfigConfigs(appID, buildID int64) (confs []*model.Config, err error) {
	var (
		tag     *model.Tag
		ids     []int64
		id      int64
		tmps    []*model.Config
		cids    []int64
		lastIDS []int64
	)
	if tag, err = s.TagByAppBuild(appID, buildID); err != nil {
		return
	}
	confIDs := tag.ConfigIDs
	if len(confIDs) == 0 {
		return
	}
	tmpIDs := strings.Split(confIDs, ",")
	for _, tmpID := range tmpIDs {
		if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", tmpID, err)
			return
		}
		ids = append(ids, id)
	}
	tmps, err = s.ConfigsByIDs(ids)
	if err != nil {
		log.Error("LastBuildConfigConfigs ids(%v) error(%v)", ids, err)
		return
	}
	for _, val := range tmps {
		cs, err := s.ConfigList(val.AppID, val.Name)
		if err == nil {
			cids = nil
		csloop:
			for _, vv := range cs {
				for _, vvv := range ids {
					if vvv == vv.ID {
						continue csloop
					}
				}
				cids = append(cids, vv.ID)
			}
			if configID, err := s.TagByAppBuildLastConfig(appID, buildID, tag.ID, cids); err == nil {
				lastIDS = append(lastIDS, configID)
			}
		}
	}
	return s.ConfigsByIDs(lastIDS)
}

// ConfigList ...
func (s *Service) ConfigList(appID int64, name string) (confs []*model.Config, err error) {
	if err = s.dao.DB.Where("app_id = ? and name = ?", appID, name).Order("id desc").Find(&confs).Error; err != nil {
		log.Error("ConfigList appid(%v) name(%v) error(%v)", appID, name, err)
	}
	return
}

// ConfigsByTagID get configs by tag id.
func (s *Service) ConfigsByTagID(tagID int64) (confs []*model.Config, err error) {
	var (
		tag *model.Tag
		ids []int64
		id  int64
	)
	if tag, err = s.Tag(tagID); err != nil {
		return
	}
	confIDs := tag.ConfigIDs
	if len(confIDs) == 0 {
		err = ecode.NothingFound
		return
	}
	tmpIDs := strings.Split(confIDs, ",")
	for _, tmpID := range tmpIDs {
		if id, err = strconv.ParseInt(tmpID, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", tmpID, err)
			return
		}
		ids = append(ids, id)
	}
	return s.ConfigsByIDs(ids)
}

//Config get Config by Config ID.
func (s *Service) Config(ID int64) (conf *model.Config, err error) {
	conf = new(model.Config)
	err = s.dao.DB.First(&conf, ID).Error
	return
}

func (s *Service) configIng(name string, appID int64) (conf *model.Config, err error) {
	conf = new(model.Config)
	if err = s.dao.DB.Select("id").Where("name = ? and app_id = ? and state=?", name, appID, model.ConfigIng).First(&conf).Error; err != nil {
		log.Error("configIng(%v) error(%v)", name, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

//Value get value by Config ID.
func (s *Service) Value(ID int64) (conf *model.Config, err error) {
	conf = new(model.Config)
	if err = s.dao.DB.First(&conf, ID).Error; err != nil {
		log.Error("Value() error(%v)", err)
	}
	return
}

//ConfigsByTree get Config by Config name.
func (s *Service) ConfigsByTree(treeID int64, env, zone, name string) (confs []*model.Config, err error) {
	var app *model.App
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		return
	}
	if err = s.dao.DB.Order("id desc").Limit(10).Find(&confs, "name = ? and app_id = ?", name, app.ID).Error; err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return
}

// NamesByAppTree get configs by app name.
func (s *Service) NamesByAppTree(appName, env, zone string, treeID int64) (names []string, err error) {
	var (
		app   *model.App
		confs []*model.Config
	)
	if app, err = s.AppByTree(treeID, env, zone); err != nil {
		if err == ecode.NothingFound {
			err = s.CreateApp(appName, env, zone, treeID)
			return
		}
	}
	if err = s.dao.DB.Select("name").Where("app_id = ?", app.ID).Order("id desc").Group("name").Find(&confs).Error; err != nil {
		log.Error("NamesByAppName(%v) error(%v)", app.ID, err)
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	for _, conf := range confs {
		names = append(names, conf.Name)
	}
	return
}

//Diff get value by Config ID.
func (s *Service) Diff(ID, BuildID int64) (data *model.Config, err error) {
	var tmpID int64
	var idArr []string
	conf := new(model.Config)
	if err = s.dao.DB.First(&conf, ID).Error; err != nil {
		log.Error("Diff() config_id(%v) error(%v)", ID, err)
		return
	}
	config := []*model.Config{}
	if err = s.dao.DB.Where("`app_id` = ? and `name` = ? and `state` = 2", conf.AppID, conf.Name).Order("id desc").Find(&config).Error; err != nil {
		log.Error("Diff() app_id(%v) name(%v) error(%v)", conf.AppID, conf.Name, err)
		return
	}
	build := &model.Build{}
	if err = s.dao.DB.First(build, BuildID).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Diff() build_id(%v) error(%v)", BuildID, err)
		return
	}
	err = nil
	if build.ID > 0 {
		tag := &model.Tag{}
		if err = s.dao.DB.First(tag, build.TagID).Error; err != nil {
			log.Error("Diff() tag_id(%v) error(%v)", build.TagID, err)
			return
		}
		idArr = strings.Split(tag.ConfigIDs, ",")
	}
	if len(idArr) > 0 {
		for _, val := range config {
			for _, vv := range idArr {
				tmpID, _ = strconv.ParseInt(vv, 10, 64)
				if tmpID == val.ID {
					data = val
					return
				}
			}
		}
	}
	for _, val2 := range config {
		if val2.ID < ID {
			data = val2
			return
		}
	}
	data = conf
	return
}

//ConfigDel config is delete.
func (s *Service) ConfigDel(ID int64) (err error) {
	conf := &model.Config{}
	if err = s.dao.DB.First(conf, ID).Error; err != nil {
		log.Error("ConfigDel first id(%v) error(%v)", ID, err)
		return
	}
	confs := []*model.Config{}
	if err = s.dao.DB.Where("app_id = ? and name = ?", conf.AppID, conf.Name).Find(&confs).Error; err != nil {
		log.Error("ConfigDel find error(%v)", err)
		return
	}
	build := []*model.Build{}
	if err = s.dao.DB.Where("app_id = ?", conf.AppID).Find(&build).Error; err != nil {
		log.Error("ConfigDel find app_id(%v) error(%v)", conf.AppID, err)
		return
	}
	for _, val := range build {
		tag := &model.Tag{}
		if err = s.dao.DB.Where("id = ?", val.TagID).First(tag).Error; err != nil {
			log.Error("ConfigDel first tag_id(%v) error(%v)", val.TagID, err)
			return
		}
		arr := strings.Split(tag.ConfigIDs, ",")
		for _, vv := range arr {
			for _, vvv := range confs {
				if vv == strconv.FormatInt(vvv.ID, 10) {
					err = ecode.NothingFound
					return
				}
			}
		}
	}
	ups := map[string]interface{}{
		"is_delete": 1,
		"state":     2,
	}
	if err = s.dao.DB.Model(conf).Where("id = ?", ID).Updates(ups).Error; err != nil {
		log.Error("ConfigDel updates error(%v)", err)
	}
	return
}

//BuildConfigInfos configRefs.
func (s *Service) BuildConfigInfos(appName, env, zone string, buildID, treeID int64) (res map[string][]*model.ConfigRefs, err error) {
	var (
		allConfs   []*model.Config
		buildConfs []*model.Config
		num        int
		ok         bool
		ref        *model.ConfigRefs
	)
	if allConfs, err = s.ConfigsByAppName(appName, env, zone, treeID, model.ConfigEnd); err != nil {
		return
	}
	if buildID != 0 {
		if buildConfs, err = s.ConfigsByBuildID(buildID); err != nil {
			return
		}
	}
	tmpMap := make(map[string]int)
	tmpBuildConfs := make(map[string]map[int64]struct{})
	refs := make(map[string]*model.ConfigRefs)
	for _, conf := range allConfs {
		if num, ok = tmpMap[conf.Name]; !ok {
			ref = new(model.ConfigRefs)
			refs[conf.Name] = ref
			tmpMap[conf.Name] = num
		} else {
			ref = refs[conf.Name]
		}
		if num <= 20 {
			ref.Configs = append(ref.Configs, &model.ConfigRef{ID: conf.ID, Mark: conf.Mark, IsDelete: conf.IsDelete})
			tmpMap[conf.Name] = num + 1
			if tmpBuildConfs[conf.Name] == nil {
				tmpBuildConfs[conf.Name] = make(map[int64]struct{})
			}
			tmpBuildConfs[conf.Name][conf.ID] = struct{}{}
		} else {
			for _, bconf := range buildConfs {
				if bconf.Name == conf.Name {
					if _, ok = tmpBuildConfs[conf.Name][bconf.ID]; !ok {
						tmpBuildConfs[conf.Name][bconf.ID] = struct{}{}
						ref.Configs = append(ref.Configs, &model.ConfigRef{ID: bconf.ID, Mark: bconf.Mark, IsDelete: bconf.IsDelete})
					}
					break
				}
			}
		}
		if ref.Ref != nil {
			continue
		}
		for _, bconf := range buildConfs {
			if bconf.Name == conf.Name {
				ref.Ref = &model.ConfigRef{ID: bconf.ID, Mark: bconf.Mark, IsDelete: bconf.IsDelete}
				break
			}
		}
	}
	res = make(map[string][]*model.ConfigRefs)
	var tp int64
	var IsDelete int64
	capacity := len(refs)
	res["new"] = make([]*model.ConfigRefs, 0, capacity)
	res["nothing"] = make([]*model.ConfigRefs, 0, capacity)
	res["notused"] = make([]*model.ConfigRefs, 0, capacity)
	for k, v := range refs {
		v.Name = k
		IsDelete = 0
		for i, tv := range v.Configs {
			if tv.IsDelete == 1 && tv.ID > IsDelete {
				IsDelete = tv.ID
				v.Configs = v.Configs[:i]
			}
		}
		v.DeleteMAX = IsDelete
		if len(v.Configs) == 0 {
			continue
		}
		if v.Ref != nil {
			tp = 0
			for _, vv := range v.Configs {
				if vv.ID > v.Ref.ID {
					tp = vv.ID
				}
			}
			if tp > 0 {
				res["new"] = append(res["new"], v)
			} else {
				res["nothing"] = append(res["nothing"], v)
			}
		} else {
			res["notused"] = append(res["notused"], v)
		}
	}
	return
}

// GetConfigs ...
func (s *Service) GetConfigs(ids []int64, name string) (configs []*model.Config, err error) {
	if err = s.dao.DB.Where("name = ? AND id in (?)", name, ids).Find(&configs).Error; err != nil {
		log.Error("GetConfigs error(%v)", err)
	}
	return
}

// GetConfig ...
func (s *Service) GetConfig(ids []int64, name string) (config *model.Config, err error) {
	config = new(model.Config)
	if err = s.dao.DB.Where("name = ? AND id in (?)", name, ids).First(config).Error; err != nil {
		log.Error("GetConfigs error(%v)", err)
	}
	return
}
