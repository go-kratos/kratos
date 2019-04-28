package http

import (
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/time"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

func createConfig(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CreateConfigReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	if _, err = svr.AuthApp(c, user, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	conf := &model.Config{}
	conf.Operator = user
	conf.Name = v.Name
	conf.Mark = v.Mark
	conf.Comment = v.Comment
	conf.State = v.State
	conf.From = v.From
	c.JSON(nil, svr.CreateConf(conf, v.TreeID, v.Env, v.Zone, v.SkipLint))
}

func lintConfig(c *bm.Context) {
	var req struct {
		Name    string `form:"name" validate:"required"`
		Comment string `form:"comment" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		// ignore error
		return
	}
	c.JSON(svr.LintConfig(req.Name, req.Comment))
}

func updateConfValue(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.UpdateConfValueReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	user := user(c)
	if _, err = svr.AuthApps(c, user, c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	conf := &model.Config{}
	conf.Name = v.Name
	conf.ID = v.ID
	conf.Operator = user
	conf.Mark = v.Mark
	conf.Comment = v.Comment
	conf.State = v.State
	conf.Mtime = time.Time(v.Mtime)
	var configs *model.Config
	configs, err = svr.Value(v.ID)
	if err != nil {
		res["message"] = "未找到源文件"
		c.JSONMap(res, err)
		return
	}
	if v.NewCommon > 0 {
		common := &model.CommonConf{}
		common2 := &model.CommonConf{}
		if err = svr.DB.Where("id = ?", configs.From).First(common).Error; err != nil {
			res["message"] = "未找到公共源文件"
			c.JSONMap(res, err)
			return
		}
		if err = svr.DB.Where("team_id = ? and name = ? and state = 2 and id = ?", common.TeamID, common.Name, v.NewCommon).Order("id desc").First(common2).Error; err != nil {
			res["message"] = "未找到最新的公共文件"
			c.JSONMap(res, err)
			return
		}
		conf.From = v.NewCommon
	}
	//验证是否最新源文件
	newConfig := &model.Config{}
	if err = svr.DB.Where("app_id = ? and name = ?", configs.AppID, configs.Name).Order("id desc").First(newConfig).Error; err != nil {
		res["message"] = "未找到最新文件"
		c.JSONMap(res, err)
		return
	}
	//默认验证ignore 0
	if newConfig.ID != v.ID && v.Ignore == 0 && user != newConfig.Operator {
		err = ecode.ConfigNotNow
		res["message"] = fmt.Sprintf("当前源文件:（%d）有最新源文件版本（%d）操作人:%s是否继续提交？", v.ID, newConfig.ID, newConfig.Operator)
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.UpdateConfValue(conf, v.SkipLint))
}

func value(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ValueReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	var TreeID int64
	TreeID, err = svr.ConfigGetTreeID(v.ConfigID)
	if err != nil {
		res["message"] = "未找到tree_id"
		c.JSONMap(res, err)
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Value(v.ConfigID))
}

func configsByBuildID(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigsByBuildIDReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigsByBuildID(v.BuildID))
}

func configsByTagID(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigsByTagIDReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigsByTagID(v.TagID))
}

func configsByAppName(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigsByAppNameReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigsByAppName(v.AppName, v.Env, v.Zone, v.TreeID, 0))
}

func configSearchAll(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigSearchAllReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	nodes, err := svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie"))
	if err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigSearchAll(c, v.Env, v.Zone, v.Like, nodes))
}

func configSearchApp(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigSearchAppReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigSearchApp(c, v.AppName, v.Env, v.Zone, v.Like, v.BuildID, v.TreeID))
}

func configsByName(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigsByNameReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigsByTree(v.TreeID, v.Env, v.Zone, v.Name))
}

func configs(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Configs(v.AppName, v.Env, v.Zone, v.BuildID, v.TreeID))
}

func configRefs(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigRefsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigRefs(v.AppName, v.Env, v.Zone, v.BuildID, v.TreeID))
}

func namesByAppName(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.NamesByAppNameReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.NamesByAppTree(v.AppName, v.Env, v.Zone, v.TreeID))
}

func diff(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.DiffReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.Diff(v.ConfigID, v.BuildID))
}

func configDel(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigDelReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApps(c, user(c), c.Request.Header.Get("Cookie")); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.ConfigDel(v.ConfigID))
}

func configBuildInfos(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigBuildInfosReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user(c), c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.BuildConfigInfos(v.AppName, v.Env, v.Zone, v.BuildID, v.TreeID))
}

func configUpdate(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.ConfigUpdateReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	app := &model.App{}
	if err = svr.DB.Where("name = ? and env = ? and zone = ? and tree_id = ? and token = ?", v.AppName, v.Env, v.Zone, v.TreeID, v.Token).First(app).Error; err != nil {
		res["message"] = "参数不正确，未找到该服务"
		c.JSONMap(res, err)
		return
	}
	var obj []map[string]string
	err = json.Unmarshal([]byte(v.Data), &obj)
	tx := svr.DB.Begin()
	for _, val := range obj {
		if len(val["name"]) > 0 {
			config := &model.Config{}
			if err = tx.Where("app_id = ? and name = ? and state = 1", app.ID, val["name"]).First(config).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					c.JSON(nil, err)
					tx.Rollback()
					return
				}
			} else {
				//把老的更新了再加新的
				ups := map[string]interface{}{
					"state": 2,
				}
				if err = tx.Model(&model.App{}).Where("id = ? ", config.ID).Updates(ups).Error; err != nil {
					c.JSON(nil, err)
					tx.Rollback()
					return
				}
			}
			//加新的
			m := &model.Config{
				AppID:    app.ID,
				Name:     val["name"],
				Comment:  val["comment"],
				State:    2,
				Mark:     val["mark"],
				Operator: v.User,
			}
			db := tx.Create(m)
			if err = db.Error; err != nil {
				res["message"] = "创建失败"
				c.JSONMap(res, err)
				tx.Rollback()
				return
			}
		} else {
			c.JSON(nil, ecode.RequestErr)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	c.JSON(nil, err)
}

func tagUpdate(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.TagUpdateReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if len(strings.TrimSpace(v.Build)) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	app := &model.App{}
	if err = svr.DB.Where("name = ? and env = ? and zone = ? and tree_id = ? and token = ?", v.AppName, v.Env, v.Zone, v.TreeID, v.Token).First(app).Error; err != nil {
		res["message"] = "参数不正确，未找到该服务"
		c.JSONMap(res, err)
		return
	}
	confs := []*model.Config{}
	tags := &model.Tag{}
	tag := &model.Tag{}
	build := &model.Build{}
	tagConfigs := []*model.Config{}
	var in []string
	var in2 []string
	var nameString string
	tmp := make(map[string]struct{})
	if v.ConfigIDs == "" && v.Names == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	} else if v.Names != "" {
		in = strings.Split(v.Names, ",")
		if err = svr.DB.Select("max(id) as id,name").Where("app_id = ? and state = 2 and is_delete = 0 and name in (?)", app.ID, in).Group("name").Find(&confs).Error; err != nil {
			res["message"] = "未找到发版文件"
			c.JSONMap(res, err)
			return
		}
		for _, vv := range confs {
			if len(nameString) > 0 {
				nameString = nameString + ","
			}
			nameString = nameString + strconv.FormatInt(vv.ID, 10)
		}
		tag.ConfigIDs = nameString
	} else if v.ConfigIDs != "" {
		in = strings.Split(v.ConfigIDs, ",")
		if err = svr.DB.Where("app_id = ? and state = 2 and is_delete = 0 and id in (?)", app.ID, in).Find(&confs).Error; err != nil {
			res["message"] = "未找到发版文件"
			c.JSONMap(res, err)
			return
		}
		tag.ConfigIDs = v.ConfigIDs
	}
	if v.Names != "" && v.Increment == 1 {
		if err = svr.DB.Where("app_id = ? and name = ?", app.ID, v.Build).Order("id desc").First(build).Error; err != nil {
			res["message"] = "未找到对应的build"
			c.JSONMap(res, err)
			return
		}
		if err = svr.DB.Where("app_id = ? and build_id = ?", app.ID, build.ID).Order("id desc").First(tags).Error; err != nil {
			res["message"] = "未找到对应的tag"
			c.JSONMap(res, err)
			return
		}
		in2 = strings.Split(tags.ConfigIDs, ",")
		if err = svr.DB.Where("app_id = ? and state = 2 and id in (?)", app.ID, in2).Find(&tagConfigs).Error; err != nil {
			res["message"] = "未找到tag中的文件"
			c.JSONMap(res, err)
			return
		}
		for _, vv := range tagConfigs {
			tss := 0
			for _, vvv := range confs {
				if vv.Name == vvv.Name {
					tss = 1
				}
			}
			if tss != 1 {
				if len(nameString) > 0 {
					nameString = nameString + ","
				}
				nameString = nameString + strconv.FormatInt(vv.ID, 10)
			}
		}
		tag.ConfigIDs = nameString
	} else {
		if len(confs) != len(in) {
			res["message"] = "发版数据不符"
			c.JSONMap(res, ecode.RequestErr)
			return
		}
		for _, vv := range confs {
			if _, ok := tmp[vv.Name]; !ok {
				tmp[vv.Name] = struct{}{}
			}
		}
		if len(tmp) != len(confs) {
			res["message"] = "有重复的文件名"
			c.JSONMap(res, ecode.RequestErr)
			return
		}
	}
	tag.Operator = v.User
	tag.Mark = v.Mark
	if v.Force == 1 {
		tag.Force = 1
	}
	c.JSON(nil, svr.UpdateTag(c, v.TreeID, v.Env, v.Zone, v.Build, tag))
}

func canalTagUpdate(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CanalTagUpdateReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.Force != 1 {
		v.Force = 0
	}
	if len(strings.TrimSpace(v.Build)) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	app := &model.App{}
	if err = svr.DB.Where("name = ? and env = ? and zone = ? and tree_id = ? and token = ?", v.AppName, v.Env, v.Zone, v.TreeID, v.Token).First(app).Error; err != nil {
		res["message"] = "参数不正确，未找到该服务"
		c.JSONMap(res, err)
		return
	}
	confs := []*model.Config{}
	configs := []*model.Config{}
	tags := &model.Tag{}
	tag := &model.Tag{}
	build := &model.Build{}
	tagConfigs := []*model.Config{}
	var in []string
	var in2 []string
	var nameString string
	tmp := make(map[string]struct{})
	in = strings.Split(v.ConfigIDs, ",")
	if err = svr.DB.Where("app_id = ? and state = 2 and id in(?)", app.ID, in).Find(&configs).Error; err != nil {
		res["message"] = "未找到文件"
		c.JSONMap(res, err)
		return
	}
	if len(configs) != len(in) {
		res["message"] = fmt.Sprintf("数据不匹配,传的数据为(%v)条,查到的数据为(%v)条,app_id(%v),config_ids(%v),in(%v)", len(in), len(configs), app.ID, v.ConfigIDs, in)
		err = ecode.RequestErr
		c.JSONMap(res, err)
		return
	}
	if err = svr.DB.Where("app_id = ? and name = ?", app.ID, v.Build).Order("id desc").First(build).Error; err != nil {
		res["message"] = "未找到对应的build"
		c.JSONMap(res, err)
		return
	}
	if err = svr.DB.Where("app_id = ? and build_id = ?", app.ID, build.ID).Order("id desc").First(tags).Error; err != nil {
		res["message"] = "未找到对应的tag"
		c.JSONMap(res, err)
		return
	}
	in2 = strings.Split(tags.ConfigIDs, ",")
	if err = svr.DB.Where("app_id = ? and state = 2 and id in (?)", app.ID, in2).Find(&tagConfigs).Error; err != nil {
		res["message"] = "未找到tag中的文件"
		c.JSONMap(res, err)
		return
	}
	if err = svr.DB.Select("id,app_id,name,`from`,state,mark,operator,ctime,mtime").Where(in2).Find(&confs).Error; err != nil {
		log.Error("ConfigsByIDs(%v) error(%v)", in2, err)
		res["message"] = "config文件未找到"
		c.JSONMap(res, err)
		return
	}
	for _, val := range confs {
		for _, vv := range configs {
			if val.Name == vv.Name {
				if len(nameString) > 0 {
					nameString = nameString + ","
				}
				nameString = nameString + strconv.FormatInt(vv.ID, 10)
				tmp[vv.Name] = struct{}{}
			}
		}
		if _, ok := tmp[val.Name]; !ok {
			if len(nameString) > 0 {
				nameString = nameString + ","
			}
			nameString = nameString + strconv.FormatInt(val.ID, 10)
		}
	}
	for _, val := range configs {
		if _, ok := tmp[val.Name]; !ok {
			if len(nameString) > 0 {
				nameString = nameString + ","
			}
			nameString = nameString + strconv.FormatInt(val.ID, 10)
		}
	}
	tag.ConfigIDs = nameString
	tag.Operator = v.User
	tag.Mark = v.Mark
	tag.Force = v.Force
	c.JSON(nil, svr.UpdateTag(c, v.TreeID, v.Env, v.Zone, v.Build, tag))
}

func canalConfigCreate(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CanalConfigCreateReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if err = svr.CanalCheckToken(v.AppName, v.Env, v.Zone, v.Token); err != nil {
		res["message"] = "未找到数据"
		c.JSONMap(res, err)
		return
	}
	conf := &model.Config{}
	conf.Operator = v.User
	conf.Name = v.Name
	conf.Mark = v.Mark
	conf.Comment = v.Comment
	conf.State = v.State
	conf.From = v.From
	c.JSON(nil, svr.CreateConf(conf, v.TreeID, v.Env, v.Zone, true))
}

func canalNameConfigs(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(model.CanalNameConfigsReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if err = svr.CanalCheckToken(v.AppName, v.Env, v.Zone, v.Token); err != nil {
		res["message"] = "未找到数据"
		c.JSONMap(res, err)
		return
	}
	c.JSON(svr.ConfigsByTree(v.TreeID, v.Env, v.Zone, v.Name))
}
