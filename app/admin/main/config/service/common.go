package service

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CreateComConf create config.
func (s *Service) CreateComConf(conf *model.CommonConf, name, env, zone string, skiplint bool) (err error) {
	if !skiplint {
		if err = lintConfig(conf.Name, conf.Comment); err != nil {
			return
		}
	}
	var team *model.Team
	if team, err = s.TeamByName(name, env, zone); err != nil {
		return
	}
	conf.TeamID = team.ID
	return s.dao.DB.Create(conf).Error
}

// ComConfig get common config by id.
func (s *Service) ComConfig(id int64) (conf *model.CommonConf, err error) {
	conf = new(model.CommonConf)
	if err = s.dao.DB.First(&conf, "id = ?", id).Error; err != nil {
		log.Error("ComConfig() error(%v)", err)
	}
	return
}

// ComConfigsByTeam common config by team.
func (s *Service) ComConfigsByTeam(name, env, zone string, ps, pn int64) (pager *model.CommonConfPager, err error) {
	var (
		team   *model.Team
		confs  []*model.CommonConf
		temp   []*model.CommonTemp
		counts model.CommonCounts
		array  []int64
	)
	if team, err = s.TeamByName(name, env, zone); err != nil {
		return
	}
	if err = s.dao.DB.Raw("select max(id) as id,count(distinct name) as counts from common_config where team_id =? group by name order by id desc", team.ID).Scan(&temp).Error; err != nil {
		log.Error("NamesByTeam(%v) error(%v)", team.ID, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	if err = s.dao.DB.Raw("select count(distinct name) as counts from common_config where team_id = ?", team.ID).Scan(&counts).Error; err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	for _, v := range temp {
		array = append(array, v.ID)
	}
	if err = s.dao.DB.Raw("select id,team_id,name,state,mark,operator,ctime,mtime from common_config where id in (?) limit ?,?", array, (pn-1)*ps, ps).Scan(&confs).Error; err != nil {
		log.Error("NamesByTeam(%v) temp(%v) error(%v)", team.ID, temp, err)
		return
	}
	return &model.CommonConfPager{Total: counts.Counts, Pn: pn, Ps: ps, Items: confs}, nil
}

//ComConfigsByName get Config by Config name.
func (s *Service) ComConfigsByName(teamName, env, zone, name string) (confs []*model.CommonConf, err error) {
	var team *model.Team
	if team, err = s.TeamByName(teamName, env, zone); err != nil {
		return
	}
	if err = s.dao.DB.Select("id,team_id,name,state,mark,operator,ctime,mtime").Where("name = ? and team_id = ?",
		name, team.ID).Order("id desc").Limit(10).Find(&confs).Error; err != nil {
		return
	}
	return
}

// UpdateComConfValue update config value.
func (s *Service) UpdateComConfValue(conf *model.CommonConf, skiplint bool) (err error) {
	if !skiplint {
		if err = lintConfig(conf.Name, conf.Comment); err != nil {
			return
		}
	}
	var confDB *model.CommonConf
	if confDB, err = s.ComConfig(conf.ID); err != nil {
		return
	}
	if confDB.State == model.ConfigIng { //judge config is configIng.
		if conf.Mtime != confDB.Mtime {
			err = ecode.TargetBlocked
			return
		}
		conf.Mtime = 0
		err = s.dao.DB.Model(&model.CommonConf{ID: confDB.ID}).Updates(conf).Error
		return
	}
	if _, err = s.comConfigIng(confDB.Name, confDB.TeamID); err == nil { //judge have configing.
		err = ecode.TargetBlocked
		return
	}
	if err == sql.ErrNoRows || err == ecode.NothingFound {
		conf.ID = 0
		conf.TeamID = confDB.TeamID
		conf.Name = confDB.Name
		conf.Mtime = 0
		return s.dao.DB.Create(conf).Error
	}
	return
}

func (s *Service) comConfigIng(name string, teamID int64) (conf *model.CommonConf, err error) {
	conf = new(model.CommonConf)
	if err = s.dao.DB.Select("id").Where("name = ? and team_id = ? and state=?", name, teamID, model.ConfigIng).First(&conf).Error; err != nil {
		log.Error("configIng(%v) error(%v)", name, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

// NamesByTeam get configs by team name.
func (s *Service) NamesByTeam(teamName, env, zone string) (names []*model.CommonName, err error) {
	var (
		team  *model.Team
		confs []*model.CommonConf
	)
	if team, err = s.TeamByName(teamName, env, zone); err != nil {
		if err == ecode.NothingFound {
			err = s.CreateTeam(teamName, env, zone)
			return
		}
	}
	if err = s.dao.DB.Where("team_id = ? and state = 2", team.ID).Order("id desc").Find(&confs).Error; err != nil {
		log.Error("NamesByTeam(%v) error(%v)", team.ID, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	tmp := make(map[string]struct{})
	for _, conf := range confs {
		if _, ok := tmp[conf.Name]; !ok {
			names = append(names, &model.CommonName{Name: conf.Name, ID: conf.ID})
			tmp[conf.Name] = struct{}{}
		}
	}
	return
}

//AppByTeam get tagMap
func (s *Service) AppByTeam(commonConfigID int64) (tagMap map[int64]*model.TagMap, err error) {
	var commonConfig *model.CommonConf
	if commonConfig, err = s.ComConfig(commonConfigID); err != nil {
		return
	}
	team := &model.Team{}
	if err = s.dao.DB.Where("id = ?", commonConfig.TeamID).First(team).Error; err != nil {
		return
	}
	commonConf := []*model.CommonConf{}
	if err = s.dao.DB.Select("id").Where("name = ? and team_id = ? and state = 2", commonConfig.Name, commonConfig.TeamID).Find(&commonConf).Error; err != nil {
		log.Error("AppByTeam() common_config error(%v)", err)
	}
	var commonConfTmp []int64
	for _, val := range commonConf {
		commonConfTmp = append(commonConfTmp, val.ID)
	}
	app := []*model.App{}
	if err = s.dao.DB.Where("name like ? and env = ? and zone = ?", team.Name+".%", team.Env, team.Zone).Find(&app).Error; err != nil {
		log.Error("AppByTeam() app error(%v)", err)
	}
	var appTmp []int64
	appMap := make(map[int64]*model.App)
	for _, val := range app {
		appMap[val.ID] = val
		appTmp = append(appTmp, val.ID)
	}
	conf := []*model.Config{}
	if err = s.dao.DB.Where("`from` in (?) and app_id in (?) and state = 2 and is_delete = 0", commonConfTmp, appTmp).Find(&conf).Error; err != nil {
		log.Error("AppByTeam() config error(%v)", err)
	}
	confMap := make(map[int64]struct{})
	for _, val := range conf {
		confMap[val.ID] = struct{}{}
	}
	build := []*model.Build{}
	if err = s.dao.DB.Where("app_id in (?)", appTmp).Find(&build).Error; err != nil {
		log.Error("AppByTeam() build error(%v)", err)
	}
	var buildTmp []int64
	buildMap := make(map[int64]string)
	for _, val := range build {
		buildMap[val.ID] = val.Name
		buildTmp = append(buildTmp, val.TagID)
	}
	tagMap = make(map[int64]*model.TagMap)
	tag := []*model.Tag{}
	if err = s.dao.DB.Where("id in (?)", buildTmp).Find(&tag).Error; err != nil {
		log.Error("AppByTeam() tag error(%v)", err)
	}
	for _, val := range tag {
		tmp := strings.Split(val.ConfigIDs, ",")
		for _, vv := range tmp {
			vv, _ := strconv.ParseInt(vv, 10, 64)
			if _, ok := confMap[vv]; !ok {
				continue
			}
			tagMap[val.ID] = &model.TagMap{Tag: val}
			if _, ok := appMap[val.AppID]; ok {
				tagMap[val.ID].AppName = appMap[val.AppID].Name
				tagMap[val.ID].TreeID = appMap[val.AppID].TreeID
			}
			if _, ok := buildMap[val.BuildID]; ok {
				tagMap[val.ID].BuildName = buildMap[val.BuildID]
			}
		}
	}
	return
}

// CommonPush ...
func (s *Service) CommonPush(ctx context.Context, tagID, commonConfigID int64, user string) (err error) {
	var tag *model.Tag
	tag, err = s.Tag(tagID)
	if err != nil {
		log.Error("CommonPush() tagid(%v) error(%v)", tagID, err)
		return
	}
	configIDS := strings.Split(tag.ConfigIDs, ",")
	app := &model.App{}
	if err = s.dao.DB.Where("id = ?", tag.AppID).First(app).Error; err != nil {
		log.Error("CommonPush() app error(%v)", err)
		return
	}
	build := &model.Build{}
	if err = s.dao.DB.Where("id = ?", tag.BuildID).First(build).Error; err != nil {
		log.Error("CommonPush() build error(%v)", err)
		return
	}
	var commonConfig *model.CommonConf
	if commonConfig, err = s.ComConfig(commonConfigID); err != nil {
		return
	}
	team := &model.Team{}
	if err = s.dao.DB.Where("id = ?", commonConfig.TeamID).First(team).Error; err != nil {
		log.Error("CommonPush() team error(%v)", err)
		return
	}
	commonConf := []*model.CommonConf{}
	if err = s.dao.DB.Select("id").Where("name = ? and team_id = ? and state = 2", commonConfig.Name, commonConfig.TeamID).Find(&commonConf).Error; err != nil {
		log.Error("CommonPush() common_config error(%v)", err)
		return
	}
	var commonConfTmp []int64
	for _, val := range commonConf {
		commonConfTmp = append(commonConfTmp, val.ID)
	}
	conf := []*model.Config{}
	if err = s.dao.DB.Where("id in (?) and `from` in (?)", configIDS, commonConfTmp).Find(&conf).Error; err != nil {
		log.Error("CommonPush() config error(%v)", err)
		return
	}
	if len(conf) != 1 {
		log.Error("CommonPush() count config(%v) error(数据有误更新数据非1条)", conf)
		return
	}
	var newConfigIDS string
	for _, val := range conf {
		newConf := &model.Config{}
		newConf.AppID = val.AppID
		newConf.Comment = commonConfig.Comment
		newConf.Mark = commonConfig.Mark
		newConf.Name = val.Name
		newConf.State = 2
		newConf.From = commonConfigID
		newConf.Operator = user
		if _, err = s.configIng(newConf.Name, app.ID); err == nil { // judge config is configIng
			err = ecode.TargetBlocked
			return
		}
		if err = s.dao.DB.Create(newConf).Error; err != nil {
			log.Error("CommonPush() create newConf error(%v)", err)
			return
		}
		newConfigIDS = strconv.FormatInt(newConf.ID, 10)
		for _, vv := range configIDS {
			if strconv.FormatInt(val.ID, 10) != vv {
				if len(newConfigIDS) > 0 {
					newConfigIDS += ","
				}
				newConfigIDS += vv
			}
		}
		//tag发版
		newTag := &model.Tag{}
		newTag.Operator = user
		newTag.Mark = tag.Mark
		newTag.ConfigIDs = newConfigIDS
		newTag.Force = 1
		err = s.UpdateTag(ctx, app.TreeID, app.Env, app.Zone, build.Name, newTag)
	}
	return
}
