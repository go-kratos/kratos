package http

import (
	"strings"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/user"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

func name(ctx *bm.Context) (name string) {
	usernameI, _ := ctx.Get("username")
	name, _ = usernameI.(string)
	return
}

func userAuth(c *bm.Context) {
	var (
		usr      = &user.User{}
		username = name(c)
		err      error
		mdls     []*user.Module
		rls      []*user.Rule
		super    bool
	)
	if usr, err = apmSvc.GetUser(c, username); err != nil {
		log.Error("apmSvc.userAuth error(%v)", err)
		c.JSON(nil, err)
		return
	}
	// err := apmSvc.DB.Where("username = ?", username).First(usr).Error
	// if err == gorm.ErrRecordNotFound {
	// 	usr.UserName = username
	// 	usr.NickName = username
	// 	err = apmSvc.DB.Create(usr).Error
	// }
	// if err != nil {
	// 	log.Error("apmSvc.userAuth error(%v)", err)
	// 	c.JSON(nil, err)
	// 	return
	// }
	for _, u := range conf.Conf.Superman {
		if u == username {
			super = true
			break
		}
	}
	var (
		ms []string
		rs []string
	)
	if super {
		for m := range user.Modules {
			ms = append(ms, m)
			for rl := range user.Rules {
				if strings.HasPrefix(rl+"_", m) {
					rs = append(rs, rl)
				}
			}
		}
	} else {
		ms, rs = apmSvc.GetDefaultPermission(c)
		if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&mdls).Error; err != nil {
			log.Error("apmSvc.userAuth modules error(%v)", err)
			c.JSON(nil, err)
			return
		}
		if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&rls).Error; err != nil {
			log.Error("apmSvc.userAuth rules error(%v)", err)
			c.JSON(nil, err)
			return
		}
		for _, m := range mdls {
			ms = append(ms, m.Module)
		}
		for _, r := range rls {
			rs = append(rs, r.Rule)
		}
	}
	data := user.Result{
		Super: super,
		User:  usr,
		Env:   env.DeployEnv,
		Rules: append(ms, rs...),
	}
	c.JSON(data, nil)
}

func userRuleStates(c *bm.Context) {
	username := name(c)
	usr := &user.User{}
	err := apmSvc.DB.Where("username = ?", username).First(usr).Error
	if err != nil {
		log.Error("apmSvc.userRuleStates error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, u := range conf.Conf.Superman {
		if u == username {
			c.JSONMap(map[string]interface{}{
				"message": "超级管理员拥有所有权限",
			}, nil)
			return
		}
	}
	var (
		//app *user.Apply
		rls []*user.Rule
	)
	app := &user.Apply{}
	if err = apmSvc.DB.Where("user_id=? AND status=?", usr.ID, 1).First(app).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("apm.Svc.userRuleStates error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&rls).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("apm.Svc.userRuleStates error(%v)", err)
		c.JSON(nil, err)
		return
	}
	trs := strings.Split(app.Rules, ",")
	type ruleRes struct {
		Rule  string `json:"rule"`
		Name  string `json:"name"`
		State int    `json:"state"`
	}
	data := map[string][]*ruleRes{}
	for module := range user.Modules {
		if env.DeployEnv != env.DeployEnvProd && module == "CONFIG" {
			continue
		}
		if user.Modules[module].Permit == user.PermitSuper {
			continue
		}
		// if module == "USER" {
		// 	continue
		// }
	NEXTRULE:
		for rule := range user.Rules {
			if !strings.HasPrefix(rule, module) {
				continue
			}
			rr := &ruleRes{Rule: rule, Name: user.Rules[rule].Des, State: 0}
			_, rdft := apmSvc.GetDefaultPermission(c)
			for _, rl := range rdft {
				if rule == rl {
					rr.State = 1
					data[module] = append(data[module], rr)
					continue NEXTRULE
				}
			}
			for _, rl := range rls {
				if rule == rl.Rule {
					rr.State = 1
					data[module] = append(data[module], rr)
					continue NEXTRULE
				}
			}
			for _, tr := range trs {
				if rule == tr {
					rr.State = 2
					data[module] = append(data[module], rr)
					continue NEXTRULE
				}
			}
			data[module] = append(data[module], rr)
		}
	}
	c.JSON(map[string]interface{}{
		"user":        usr,
		"rule_states": data,
	}, nil)
}

func userApply(c *bm.Context) {
	username := name(c)
	usr := &user.User{}
	if err := apmSvc.DB.Where("username = ?", username).First(usr).Error; err != nil {
		log.Error("apmSvc.userApply error(%v)", err)
		c.JSON(nil, err)
		return
	}
	v := new(struct {
		Rules []string `form:"rules,split" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, rule := range v.Rules {
		if _, ok := user.Rules[rule]; !ok {
			c.JSONMap(map[string]interface{}{
				"message": "申请的操作不存在",
			}, nil)
			return
		}
	}
	istr := strings.Join(v.Rules, ",")
	m := &user.Apply{
		UserID: usr.ID,
		Rules:  istr,
		Status: 1,
	}
	db := apmSvc.DB.Model(&user.Apply{}).Create(m)
	if err := db.Error; err != nil {
		log.Error("apmSvc.userApply error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"message": "申请成功",
	}, nil)
}

func userApplyEdit(c *bm.Context) {
	v := new(struct {
		ID    int64    `form:"id" validate:"required"`
		Rules []string `form:"rules,split" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	for _, r := range v.Rules {
		if _, ok := user.Rules[r]; !ok {
			c.JSONMap(map[string]interface{}{
				"message": "申请的操作不存在",
			}, nil)
			return
		}
	}
	username := name(c)
	if err = apmSvc.DB.Model(&user.Apply{}).Where("status = 1 AND id = ?", v.ID).Update(map[string]interface{}{
		"rules": strings.Join(v.Rules, ","), "admin": username}).Error; err != nil {
		log.Error("apmSvc.userApplyEdit error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"message": "修改成功",
	}, nil)
}

func userAudit(c *bm.Context) {
	username := name(c)
	super := false
	for _, u := range conf.Conf.Superman {
		if u == username {
			super = true
			break
		}
	}
	if !super {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	v := new(struct {
		ID     int64 `form:"id" validate:"required"`
		Status int8  `form:"status" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if !(v.Status == 2 || v.Status == 3) {
		log.Error("apmSvc.userAudit error(%v)", v.Status)
		c.JSONMap(map[string]interface{}{
			"message": "status值范围为2，3",
		}, ecode.RequestErr)
		return
	}
	if err := apmSvc.DB.Model(&user.Apply{}).Where("id = ? AND status = ?", v.ID, 1).Updates(map[string]interface{}{"status": v.Status, "admin": username}).Error; err != nil {
		log.Error("apmSvc.userAudit update user_apply error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.Status == 3 {
		c.JSONMap(map[string]interface{}{
			"message": "权限审核不通过",
		}, nil)
		return
	}
	apps := &user.Apply{}
	if err := apmSvc.DB.Where("id=?", v.ID).First(apps).Error; err != nil {
		log.Error("apmSvc.userAudit find user_apply error(%v)", err)
		c.JSON(nil, err)
		return
	}
	rules := strings.Split(apps.Rules, ",")
	for _, rule := range rules {
		r := &user.Rule{}
		apmSvc.DB.FirstOrCreate(r, &user.Rule{UserID: apps.UserID, Rule: rule})
		for module := range user.Modules {
			if strings.HasPrefix(rule, module) {
				m := &user.Module{}
				apmSvc.DB.FirstOrCreate(m, &user.Module{UserID: apps.UserID, Module: module})
			}
		}
	}
	c.JSONMap(map[string]interface{}{
		"message": "权限审核通过",
	}, nil)
}

func userApplies(c *bm.Context) {
	username := name(c)
	v := new(struct {
		Pn   int    `form:"pn" default:"1" validate:"min=1"`
		Ps   int    `form:"ps" default:"20" validate:"min=1"`
		Name string `form:"name"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		super   bool
		total   int
		applies []*user.Applies
	)
	for _, u := range conf.Conf.Superman {
		if u == username {
			super = true
			break
		}
	}
	if !super {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if v.Name != "" {
		err = apmSvc.DB.Raw(`SELECT user_apply.id, user_apply.user_id,user.username,user_apply.rules,user_apply.status 
			FROM user_apply LEFT JOIN user ON user_apply.user_id=user.id WHERE user_apply.status=? AND (user.username like ? OR user.nickname like ?)`,
			1, "%"+v.Name+"%", "%"+v.Name+"%").Order("user_apply.id desc").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&applies).Error
	} else {
		err = apmSvc.DB.Raw(`SELECT user_apply.id, user_apply.user_id,user.username,user_apply.rules,user_apply.status 
			FROM user_apply LEFT JOIN user ON user_apply.user_id=user.id WHERE user_apply.status=?`,
			1).Order("user_apply.id desc").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&applies).Error
	}
	// err := apmSvc.DB.Raw(`SELECT user_apply.id, user_apply.user_id,user.username,user_apply.rules,user_apply.status
	// FROM user_apply LEFT JOIN user ON user_apply.user_id=user.id WHERE user_apply.status=?`, 1).Scan(&applies).Error
	if err == gorm.ErrRecordNotFound {
		c.JSONMap(map[string]interface{}{
			"message": "当前没有任何申请",
		}, nil)
		return
	}
	if v.Name != "" {
		err = apmSvc.DB.Model(&user.Apply{}).Joins("LEFT JOIN user ON user_apply.user_id=user.id").Where(`user_apply.status=? 
			AND (user.username like ? OR user.nickname like ?)`, 1, "%"+v.Name+"%", "%"+v.Name+"%").Count(&total).Error
	} else {
		err = apmSvc.DB.Model(&user.Apply{}).Joins(`LEFT JOIN user ON user_apply.user_id=user.id`).Where(`user_apply.status=?`, 1).Count(&total).Error
	}
	if err != nil {
		log.Error("apmSvc.userApplies error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: applies,
		Total: total,
	}
	c.JSON(data, nil)
}

func userList(c *bm.Context) {
	v := new(struct {
		Pn   int    `form:"pn" default:"1" validate:"min=1"`
		Ps   int    `form:"ps" default:"20" validate:"min=1"`
		Name string `form:"name"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		pts   []*user.User
		total int
	)
	s := "%" + v.Name + "%"
	if v.Name != "" {
		err = apmSvc.DB.Where("username LIKE ? OR nickname LIKE ?", s, s).Order("id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&pts).Error
	} else {
		err = apmSvc.DB.Order("id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&pts).Error
	}
	if err != nil {
		log.Error("apmSvc.Users error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.Name != "" {
		err = apmSvc.DB.Where("username LIKE ? OR nickname LIKE ?", s, s).Model(&user.User{}).Count(&total).Error
	} else {
		err = apmSvc.DB.Model(&user.User{}).Count(&total).Error
	}
	if err != nil {
		log.Error("apmSvc.Users count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: pts,
		Total: total,
	}
	c.JSON(data, nil)
}

func userInfo(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	usr := &user.User{}
	if err = apmSvc.DB.First(usr, v.ID).Error; err != nil {
		log.Error("apmSvc.userInfo error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(usr, nil)
}

func userEdit(c *bm.Context) {
	v := new(struct {
		ID       int64  `form:"id" validate:"required"`
		Nickname string `form:"nickname"`
		Email    string `form:"email"`
		Phone    string `form:"phone"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	usr := &user.User{}
	if err = apmSvc.DB.First(usr, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DB.Model(&user.User{}).Where("id = ?", v.ID).Omit("id").UpdateColumns(v).Error; err != nil {
		log.Error("apmSvc.userEdit error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "id = ?",
		"Value1":  v.ID,
		"Update":  v,
		"Old":     usr,
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 2, int64(v.ID), "apmSvc.userEdit", sqlLog)
	c.JSON(nil, err)
}

func userModules(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	usr := &user.User{}
	if err = apmSvc.DB.First(usr, v.ID).Error; err != nil {
		log.Error("apmSvc.userInfo error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var mdls []*user.Module
	if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&mdls).Error; err != nil {
		log.Error("apmSvc.userAuth modules error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var ms []string
	for _, m := range mdls {
		ms = append(ms, m.Module)
	}
	allMds := make(map[string]string)
	for module := range user.Modules {
		allMds[module] = user.Modules[module].Des
	}
	data := map[string]interface{}{
		"owns":    ms,
		"modules": allMds,
	}
	c.JSON(data, nil)
}

func userRules(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	usr := &user.User{}
	if err = apmSvc.DB.First(usr, v.ID).Error; err != nil {
		log.Error("apmSvc.userInfo error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var (
		mdls []*user.Module
		rls  []*user.Rule
	)
	if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&mdls).Error; err != nil {
		log.Error("apmSvc.userAuth modules error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&rls).Error; err != nil {
		log.Error("apmSvc.userAuth rules error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var rs []string
	for _, r := range rls {
		rs = append(rs, r.Rule)
	}
	allRls := map[string]string{}
	for _, mdl := range mdls {
		for rl, rlM := range user.Rules {
			if strings.HasPrefix(rl+"_", mdl.Module) {
				allRls[rl] = rlM.Des
			}
		}
	}
	data := map[string]interface{}{
		"owns":  rs,
		"rules": allRls,
	}
	c.JSON(data, nil)
}

func userModulesEdit(c *bm.Context) {
	v := new(struct {
		ID      int64    `form:"id" validate:"required"`
		Modules []string `form:"modules,split"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	usr := &user.User{}
	if err = apmSvc.DB.First(usr, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	var mdls []*user.Module
	if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&mdls).Error; err != nil {
		log.Error("apmSvc.moduleEdit modules error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var (
		ns []*user.Module
		ds []string
	)
	for _, m := range v.Modules {
		if len(mdls) == 0 {
			ns = append(ns, &user.Module{UserID: usr.ID, Module: m})
			continue
		}
		for j, mdl := range mdls {
			if m == mdl.Module {
				break
			}
			if j+1 == len(mdls) {
				ns = append(ns, &user.Module{UserID: usr.ID, Module: m})
			}
		}
	}
	for _, mdl := range mdls {
		if len(v.Modules) == 0 {
			ds = append(ds, mdl.Module)
			continue
		}
		for j, m := range v.Modules {
			if m == mdl.Module {
				break
			}
			if j+1 == len(v.Modules) {
				ds = append(ds, mdl.Module)
			}
		}
	}
	if err = apmSvc.DB.Exec("DELETE FROM user_module WHERE user_id=? AND module IN (?)", usr.ID, ds).Error; err != nil {
		log.Error("apmSvc.moduleEdit delModule error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var sqlLogs []*map[string]interface{}
	sqlLog := &map[string]interface{}{
		"SQLType": "delete",
		"Where":   "DELETE FROM user_module WHERE user_id=? AND module IN (?)",
		"Value1":  usr.ID,
		"Value2":  ds,
		"Update":  "",
		"Old":     "",
	}
	sqlLogs = append(sqlLogs, sqlLog)
	username := name(c)
	// apmSvc.SendLog(c, username, 0, 2, int64(v.ID), "apmSvc.moduleEdit", sqlLog)
	for _, d := range ds {
		if err = apmSvc.DB.Exec("DELETE FROM user_rule WHERE user_id=? AND rule LIKE ?", usr.ID, d+"_%").Error; err != nil {
			log.Error("apmSvc.moduleEdit delModule error(%v)", err)
			c.JSON(nil, err)
			apmSvc.SendLog(*c, username, 0, 2, 0, "apmSvc.moduleEdit", sqlLogs)
			return
		}
		sqlLog := &map[string]interface{}{
			"SQLType": "delete",
			"Where":   "DELETE FROM user_rule WHERE user_id=? AND rule LIKE ?",
			"Value1":  usr.ID,
			"Value2":  d + "_%",
			"Update":  "",
			"Old":     "",
		}
		sqlLogs = append(sqlLogs, sqlLog)
	}
	for _, n := range ns {
		if err = apmSvc.DB.Create(n).Error; err != nil {
			log.Error("apmSvc.moduleEdit addModule error(%v)", err)
			c.JSON(nil, err)
			apmSvc.SendLog(*c, username, 0, 2, 0, "apmSvc.moduleEdit", sqlLogs)
			return
		}
		sqlLog := &map[string]interface{}{
			"SQLType": "add",
			"Content": n,
		}
		sqlLogs = append(sqlLogs, sqlLog)
	}
	apmSvc.SendLog(*c, username, 0, 2, 0, "apmSvc.moduleEdit", sqlLogs)
	c.JSON(nil, err)
}

func userRulesEdit(c *bm.Context) {
	v := new(struct {
		ID    int64    `form:"id" validate:"required"`
		Rules []string `form:"rules,split"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	usr := &user.User{}
	if err = apmSvc.DB.First(usr, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	var mdls []*user.Module
	if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&mdls).Error; err != nil {
		log.Error("apmSvc.moduleEdit modules error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if len(mdls) == 0 {
		log.Error("apmSvc.moduleEdit have not module error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, rl := range v.Rules {
		var has bool
		for _, mdl := range mdls {
			if has = strings.HasPrefix(rl, mdl.Module); has {
				break
			}
		}
		if !has {
			log.Error("apmSvc.moduleEdit have not module error(%v)", err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	var rls []*user.Rule
	if err = apmSvc.DB.Where("user_id=?", usr.ID).Find(&rls).Error; err != nil {
		log.Error("apmSvc.ruleEdit modules error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var (
		ns []*user.Rule
		ds []string
	)
	for _, m := range v.Rules {
		if len(rls) == 0 {
			ns = append(ns, &user.Rule{UserID: usr.ID, Rule: m})
			continue
		}
		for j, rl := range rls {
			if m == rl.Rule {
				break
			}
			if j+1 == len(rls) {
				ns = append(ns, &user.Rule{UserID: usr.ID, Rule: m})
			}
		}
	}
	for _, rl := range rls {
		if len(v.Rules) == 0 {
			ds = append(ds, rl.Rule)
			continue
		}
		for j, m := range v.Rules {
			if m == rl.Rule {
				break
			}
			if j+1 == len(v.Rules) {
				ds = append(ds, rl.Rule)
			}
		}
	}
	var sqlLogs []*map[string]interface{}
	if err = apmSvc.DB.Exec("DELETE FROM user_rule WHERE user_id=? AND rule IN (?)", usr.ID, ds).Error; err != nil {
		log.Error("apmSvc.ruleEdit delModule error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "delete",
		"Where":   "DELETE FROM user_rule WHERE user_id=? AND rule IN (?)",
		"Value1":  usr.ID,
		"Value2":  ds,
		"Update":  "",
		"Old":     "",
	}
	username := name(c)
	sqlLogs = append(sqlLogs, sqlLog)
	for _, n := range ns {
		if err = apmSvc.DB.Create(n).Error; err != nil {
			log.Error("apmSvc.ruleEdit addModule error(%v)", err)
			c.JSON(nil, err)
			apmSvc.SendLog(*c, username, 0, 2, 0, "apmSvc.ruleEdit", sqlLogs)
			return
		}
		sqlLog := &map[string]interface{}{
			"SQLType": "add",
			"Content": n,
		}
		sqlLogs = append(sqlLogs, sqlLog)
	}
	apmSvc.SendLog(*c, username, 0, 2, 0, "apmSvc.ruleEdit", sqlLogs)
	c.JSON(nil, err)
}

func userSyncTree(c *bm.Context) {
	username := name(c)
	apmSvc.TreeSync(c, username, c.Request.Header.Get("Cookie"))
	c.JSON(nil, nil)
}

func userTreeAppids(c *bm.Context) {
	username := name(c)
	appids, err := apmSvc.Appids(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		log.Error("%v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(appids, nil)
}

func userTreeDiscovery(c *bm.Context) {
	username := name(c)
	appids, err := apmSvc.DiscoveryID(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		log.Error("%v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(appids, nil)
}
