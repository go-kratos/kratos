package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/config/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

var (
	rdsEnvs = []*model.Env{
		{Name: "dev", NikeName: "开发环境"},
		{Name: "fat1", NikeName: "功能环境1"},
		{Name: "uat", NikeName: "集成环境"},
		{Name: "pre", NikeName: "预发环境"},
		{Name: "prod", NikeName: "线上环境"},
	}
	opsEnvs = []*model.Env{
		{Name: "dev", NikeName: "开发环境"},
		{Name: "fat1", NikeName: "功能环境1"},
		{Name: "uat", NikeName: "集成环境"},
	}
)

// CreateApp create App.
func (s *Service) CreateApp(name, env, zone string, treeID int64) error {
	bytes := [16]byte(uuid.NewV1())
	token := md5.Sum([]byte(hex.EncodeToString(bytes[:])))
	app := &model.App{Name: name, Env: env, Zone: zone, Token: hex.EncodeToString(token[:]), TreeID: treeID, Status: model.StatusShow}
	return s.dao.DB.Create(app).Error
}

// UpdateToken update token.
func (s *Service) UpdateToken(c context.Context, env, zone string, treeID int64) (err error) {
	bytes := [16]byte(uuid.NewV1())
	token := hex.EncodeToString(bytes[:])
	if err = s.dao.DB.Model(&model.App{}).Where("tree_id =? AND env=? AND zone=?", treeID, env, zone).Update("token", token).Error; err != nil {
		return
	}
	err = s.SetToken(c, treeID, env, zone, token)
	return
}

// AppByTree get token by Name.
func (s *Service) AppByTree(treeID int64, env, zone string) (app *model.App, err error) {
	app = &model.App{}
	row := s.dao.DB.Select("id,token").Where("tree_id=? AND env=? AND zone=?", treeID, env, zone).Model(&model.App{}).Row()
	if err = row.Scan(&app.ID, &app.Token); err != nil {
		log.Error("AppByTree(%v) err(%v)", treeID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

// AppsByTreeZone get token by Name and zone.
func (s *Service) AppsByTreeZone(treeID int64, zone string) (apps []*model.App, err error) {
	if err = s.dao.DB.Select("id,env,token").Where("tree_id=? AND zone=?", treeID, zone).Find(&apps).Error; err != nil {
		log.Error("AppsByTreezone(%d) error(%v)", treeID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

// AppList get token by Name.
func (s *Service) AppList(ctx context.Context, bu, team, name, env, zone string, ps, pn int64, nodes *model.CacheData, status int8) (pager *model.AppPager, err error) {
	var (
		like     string
		apps     []*model.App
		total    int64
		ids      []int64
		statusIn []int8
	)
	like = "%"
	if len(bu) != 0 {
		like = bu + ".%"
	}
	if len(team) != 0 {
		like = team + ".%"
	}
	if len(name) != 0 {
		like = "%.%.%" + name + "%"
	}
	if len(name) != 0 && len(team) != 0 {
		like = team + ".%" + name + "%"
	} else if len(name) != 0 && len(bu) != 0 {
		like = bu + ".%.%" + name + "%"
	}
	if status > 0 {
		statusIn = append(statusIn, status)
	} else {
		statusIn = []int8{model.StatusShow, model.StatusHidden}
	}
	for _, node := range nodes.Data {
		ids = append(ids, node.ID)
	}
	if err = s.dao.DB.Where("env=? AND zone=? AND name like ? AND tree_id in (?) And status in (?)", env, zone, like, ids, statusIn).
		Offset((pn - 1) * ps).Limit(ps).Find(&apps).Error; err != nil {
		log.Error("AppList() find page apps() error(%v)", err)
		return
	}
	if err = s.dao.DB.Model(&model.App{}).Where("env=? AND zone=? AND name like ? AND tree_id in (?) And status in (?)", env, zone, like, ids, statusIn).
		Count(&total).Error; err != nil {
		log.Error("AppList() count page apps() error(%v)", err)
		return
	}
	pager = &model.AppPager{Total: total, Pn: pn, Ps: ps, Items: apps}
	return
}

// Tree get service tree.
func (s *Service) Tree(ctx context.Context, user string) (data interface{}, err error) {
	var (
		parme []byte
		msg   map[string]interface{}
		tmp   interface{}
		token string
		ok    bool
	)
	if parme, err = json.Marshal(map[string]string{"user_name": user, "platform_id": s.c.Tree.Platform}); err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	if msg, err = s.dao.Token(ctx, string(parme)); err != nil {
		return
	}
	if tmp, ok = msg["token"]; !ok {
		err = ecode.NothingFound
		return
	}
	if token, ok = tmp.(string); !ok {
		err = ecode.NothingFound
		return
	}
	return s.dao.Tree(ctx, token)
}

// Node node.
func (s *Service) Node(ctx context.Context, user, node, team, cookie string, nodes *model.CacheData) (res []*model.Node, err error) {
	var nNodes *model.CacheData
	//bu list.
	if node == "" && team == "" {
		if nNodes, err = s.SyncTree(ctx, user, cookie); err == nil {
			nodes = nNodes
		}
		tmp := make(map[string]struct{})
		for _, app := range nodes.Data {
			idx := strings.Index(app.Path, ".")
			bu := string([]byte(app.Path)[:idx])
			if _, ok := tmp[bu]; ok {
				continue
			}
			n := new(model.Node)
			n.Name = bu
			n.Path = bu
			res = append(res, n)
			tmp[bu] = struct{}{}
		}
		return
	}
	//team list.
	if node != "" && team == "" {
		tmp := make(map[string]struct{})
		for _, app := range nodes.Data {
			s := []byte(app.Path)
			sep := []byte(".")
			fidx := bytes.Index(s, sep)
			lidx := bytes.LastIndex(s, sep)
			team = string(s[:lidx])
			if node == string(s[:fidx]) {
				if _, ok := tmp[team]; ok {
					continue
				}
				n := new(model.Node)
				n.Name = string([]byte(app.Path)[fidx+1 : lidx])
				n.Path = team
				tmp[team] = struct{}{}
				res = append(res, n)
			}
		}
		return
	}
	//app list.
	if team == "" {
		return
	}
	for _, app := range nodes.Data {
		s := []byte(app.Path)
		sep := []byte(".")
		lidx := bytes.LastIndex(s, sep)
		t := string(s[:lidx])
		if team != t {
			continue
		}
		n := new(model.Node)
		n.Name = string(s[lidx+1:])
		n.Path = app.Path
		n.TreeID = app.ID
		res = append(res, n)
	}
	return
}

//Envs envs.
func (s *Service) Envs(ctx context.Context, user, appName, zone string, treeID int64, nodes *model.CacheData) (envs []*model.Env, err error) {
	var (
		ok   bool
		node *model.RoleNode
		apps []*model.App
	)
	envs = rdsEnvs
	if node, ok = nodes.Data[treeID]; !ok {
		return
	}
	if node.Role == model.Ops {
		envs = opsEnvs
	}
	apps, err = s.AppsByTreeZone(treeID, zone)
	for _, env := range envs {
		env.Token = ""
		for _, app := range apps {
			if app.Env == env.Name {
				env.Token = app.Token
				break
			}
		}
	}
	return
}

//EnvsByTeam envs.
func (s *Service) EnvsByTeam(ctx context.Context, appName, zone string, nodes *model.CacheData) (envs []*model.Env, err error) {
	envs = rdsEnvs
	return
}

//SyncTree syncTree.
func (s *Service) SyncTree(ctx context.Context, user string, cookie string) (nodes *model.CacheData, err error) {
	var (
		msg   map[string]interface{}
		tmp   interface{}
		token string
		ok    bool
	)

	if msg, err = s.dao.Auth(ctx, cookie); err != nil {
		return
	}
	if tmp, ok = msg["token"]; !ok {
		err = ecode.NothingFound
		return
	}
	if token, ok = tmp.(string); !ok {
		err = ecode.NothingFound
		return
	}
	if nodes, err = s.dao.Role(ctx, user, token); err != nil {
		return
	}
	s.cLock.Lock()
	s.cache[user] = nodes
	s.cLock.Unlock()
	return
}

//AuthApps authApps.
func (s *Service) AuthApps(ctx context.Context, user, cookie string) (nodes *model.CacheData, err error) {
	if len(user) == 0 {
		err = ecode.NothingFound
		return
	}
	var ok bool
	s.cLock.RLock()
	nodes, ok = s.cache[user]
	s.cLock.RUnlock()
	if !ok || (time.Since(nodes.CTime) > 60*time.Second) {
		s.SyncTree(ctx, user, cookie)
		s.cLock.RLock()
		nodes, ok = s.cache[user]
		s.cLock.RUnlock()
		if !ok {
			err = ecode.NothingFound
		}
	}
	return
}

//AuthApp authApp.
func (s *Service) AuthApp(ctx context.Context, user, cookie string, treeID int64) (rule int8, err error) {
	var (
		ok    bool
		node  *model.RoleNode
		nodes *model.CacheData
	)
	if nodes, err = s.AuthApps(ctx, user, cookie); err != nil {
		return
	}
	if node, ok = nodes.Data[treeID]; !ok {
		err = ecode.AccessDenied
		return
	}
	return node.Role, nil
}

//ConfigGetTreeID ...
func (s *Service) ConfigGetTreeID(configID int64) (TreeID int64, err error) {
	conf := new(model.Config)
	if err = s.dao.DB.First(&conf, configID).Error; err != nil {
		log.Error("ConfigGetTreeID(%v) error(%v)", configID, err)
		return
	}
	TreeID, err = s.AppIDGetTreeID(conf.AppID)
	return
}

//AppIDGetTreeID ...
func (s *Service) AppIDGetTreeID(appID int64) (TreeID int64, err error) {
	app := new(model.App)
	if err = s.dao.DB.First(&app, appID).Error; err != nil {
		log.Error("AppIDGetTreeID(%v) error(%v)", appID, err)
		return
	}
	TreeID = app.TreeID
	return
}

//BuildGetTreeID ...
func (s *Service) BuildGetTreeID(buildID int64) (TreeID int64, err error) {
	build := new(model.Build)
	if err = s.dao.DB.First(&build, buildID).Error; err != nil {
		log.Error("BuildGetTreeID(%v) error(%v)", buildID, err)
		return
	}
	TreeID, err = s.AppIDGetTreeID(build.AppID)
	return
}

//TagGetTreeID ...
func (s *Service) TagGetTreeID(tagID int64) (TreeID int64, err error) {
	tag := new(model.Tag)
	if err = s.dao.DB.First(&tag, tagID).Error; err != nil {
		log.Error("TagGetTreeID(%v) error(%v)", tagID, err)
		return
	}
	TreeID, err = s.AppIDGetTreeID(tag.AppID)
	return
}

//ZoneCopy ...
func (s *Service) ZoneCopy(ctx context.Context, AppName, From, To string, TreeID int64) (err error) {
	apps := []*model.App{}
	if err = s.dao.DB.Where("name = ? and tree_id = ? and zone = ?", AppName, TreeID, From).Find(&apps).Error; err != nil {
		log.Error("ZoneCopy from apps error(%v)", err)
		return
	}
	tx := s.dao.DB.Begin()
	for _, v := range apps {
		app := &model.App{}
		if err = s.dao.DB.Where("name = ? and tree_id = ? and zone = ? and env = ?", AppName, TreeID, To, v.Env).First(app).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error("ZoneCopy to app error(%v)", err)
				return
			}
			//add app
			if err = s.CreateApp(AppName, v.Env, To, TreeID); err != nil {
				log.Error("ZoneCopy add app error(%v)", err)
				return
			}
			if err = s.dao.DB.Where("name = ? and tree_id = ? and zone = ? and env = ?", AppName, TreeID, To, v.Env).First(app).Error; err != nil {
				log.Error("ZoneCopy first app error(%v)", err)
				return
			}
		}
		//
		configs := []*model.Config{}
		if err = tx.Where("app_id = ?", app.ID).Find(&configs).Error; err != nil {
			log.Error("ZoneCopy find configs error(%v)", err)
			tx.Rollback()
			return
		}

		builds := []*model.Build{}
		if err = tx.Where("app_id = ?", v.ID).Find(&builds).Error; err != nil {
			log.Error("ZoneCopy find builds error(%v)", err)
			tx.Rollback()
			return
		}
		for _, val := range builds {
			tag := &model.Tag{}
			if err = tx.First(tag, val.TagID).Error; err != nil {
				log.Error("ZoneCopy find tag error(%v)", err)
				tx.Rollback()
				return
			}
			configs = []*model.Config{}
			in := strings.Split(tag.ConfigIDs, ",")
			if err = tx.Where("id in (?)", in).Find(&configs).Error; err != nil {
				log.Error("ZoneCopy find build configs error(%v)", err)
				tx.Rollback()
				return
			}
			config := &model.Config{}
			var configIDS string
			for _, vvv := range configs {
				config = &model.Config{}
				config.Operator = vvv.Operator
				config.Name = vvv.Name
				config.Mark = vvv.Mark
				config.Comment = vvv.Comment
				config.State = vvv.State
				config.From = 0 //公共文件变私人文件
				if err = s.CreateConf(config, TreeID, v.Env, To, true); err != nil {
					log.Error("ZoneCopy config create error(%v)", err)
					tx.Rollback()
					return
				}
				if len(configIDS) > 0 {
					configIDS += ","
				}
				configIDS += strconv.FormatInt(config.ID, 10)
			}
			newTag := &model.Tag{}
			newTag.Operator = tag.Operator
			newTag.Mark = tag.Mark
			newTag.ConfigIDs = configIDS
			s.UpdateTag(ctx, TreeID, v.Env, To, val.Name, newTag)
		}
	}
	tx.Commit()
	return
}

// CanalCheckToken ...
func (s *Service) CanalCheckToken(AppName, Env, Zone, Token string) (err error) {
	app := &model.App{}
	if err = s.dao.DB.Where("name = ? and env = ? and zone = ? and tree_id = ? and token = ?", AppName, Env, Zone, 3766, Token).First(app).Error; err != nil {
		log.Error("canalCheckToken error(%v)", err)
	}
	return
}

// CasterEnvs ...
func (s *Service) CasterEnvs(zone string, treeID int64) (envs []*model.Env, err error) {
	var (
		apps []*model.App
	)
	envs = rdsEnvs
	apps, err = s.AppsByTreeZone(treeID, zone)
	for _, env := range envs {
		env.Token = ""
		for _, app := range apps {
			if app.Env == env.Name {
				env.Token = app.Token
				break
			}
		}
	}
	return
}

// AppRename ...
func (s *Service) AppRename(treeID int64, user, cookie string) (err error) {
	var (
		ok    bool
		node  *model.RoleNode
		nodes *model.CacheData
	)
	s.cLock.RLock()
	nodes, ok = s.cache[user]
	s.cLock.RUnlock()
	if !ok {
		err = ecode.NothingFound
		return
	}
	if node, ok = nodes.Data[treeID]; !ok {
		err = ecode.AccessDenied
		return
	}
	if len(node.Path) == 0 {
		err = ecode.NothingFound
		return
	}
	if err = s.dao.DB.Model(&model.App{}).Where("tree_id =?", treeID).Update("name", node.Path).Error; err != nil {
		log.Error("AppRename update error(%v)", err)
		return
	}
	return
}

// GetApps ...
func (s *Service) GetApps(env string) (apps []*model.App, err error) {
	if err = s.dao.DB.Where("env = ?", env).Find(&apps).Error; err != nil {
		log.Error("GetApps error(%v)", err)
	}
	return
}

// IdsGetApps ...
func (s *Service) IdsGetApps(ids []int64) (apps []*model.App, err error) {
	if err = s.dao.DB.Where("id in (?)", ids).Find(&apps).Error; err != nil {
		log.Error("IdsGetApps error(%v)", err)
	}
	return
}

// UpAppStatus edit status.
func (s *Service) UpAppStatus(ctx context.Context, status int8, treeID int64) (err error) {
	var (
		apps []*model.App
	)
	ups := map[string]interface{}{
		"status": status,
	}
	if err = s.dao.DB.Model(apps).Where("tree_id = ?", treeID).Updates(ups).Error; err != nil {
		log.Error("AppStatus error(%v) status(%v)", err, status)
	}
	return
}
