package http

import (
	"encoding/json"
	"strconv"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryScripts(c *bm.Context) {
	qsr := model.QueryScriptRequest{}
	if err := c.BindWith(&qsr, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	qsr.Active = 1
	if err := qsr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	sessionID, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}

	res, err := srv.QueryScriptsByPage(c, sessionID.Value, &qsr)
	if err != nil {
		log.Error("queryScripts Error", err)
		return
	}
	c.JSON(res, err)
}

func queryScripysFree(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	script.Active = 1
	res, err := srv.QueryScripts(&script, 1, 200)
	if err != nil {
		log.Error("queryScripts Error", err)
		return
	}
	var resMap = make(map[string]interface{})
	resMap["scripts"] = res
	c.JSON(resMap, err)
}

func queryScriptSnap(c *bm.Context) {
	scriptSnap := model.ScriptSnap{}
	if err := c.BindWith(&scriptSnap, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	scriptSnaps, err := srv.QueryScriptSnap(&scriptSnap)
	if err != nil {
		log.Error("QueryScriptSnap Error", err)
		return
	}
	var resMap = make(map[string]interface{})
	resMap["scriptSnaps"] = scriptSnaps
	c.JSON(resMap, nil)
}

func addJmeterSample(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.AddJmeterSample(&script))
}

func addThreadGroup(c *bm.Context) {
	script := model.Script{}
	sceneTyped := c.Request.Form.Get("scene_type")
	sceneType, err := strconv.Atoi(sceneTyped)
	if err != nil {
		log.Error("your string cannot strconv to int ")
		return
	}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.AddThreadGroup(&script, sceneType))
}

func getThreadGroup(c *bm.Context) {
	scrThreadGroup := model.ScrThreadGroup{}
	if err := c.BindWith(&scrThreadGroup, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.GetThreadGroup(scrThreadGroup))
}

func addAndExecuteScript(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}

	// 压测时间段检查
	userName, _ := c.Request.Cookie("username")
	// 非debug逻辑 && 需要检查 && 非保存 && 非复制
	if !script.IsDebug && conf.Conf.Melloi.CheckTime && !script.IsSave && !script.IsCopy {
		//if !srv.CheckRunTime() {
		//	if !srv.CheckRunPermission(userName.Value) {
		//		c.JSON("Non-working time", ecode.MelloiRunNotInTime)
		//		return
		//	}
		//}
		if !srv.CheckRunPermission(userName.Value) {
			c.JSON("Non-working time", ecode.MelloiRunNotInTime)
			return
		}
	}
	JSON, _ := json.Marshal(script)
	log.Info("script:----------", string(JSON))
	cookie := c.Request.Header.Get("Cookie")
	resp, err := srv.AddAndExcuScript(c, &script, cookie, &scene, true, false)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(resp, err)
}

func delScript(c *bm.Context) {
	id := c.Request.Form.Get("id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.DelScript(ID))
}

func updateScript(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	updateResult, err := srv.UpdateScript(&script)
	if err != nil {
		log.Error("UpdateScript err (%v)", err)
		return
	}
	var resMap = make(map[string]string)
	resMap["update_result"] = updateResult
	c.JSON(resMap, err)
}

func updateScriptAll(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	updateResult, err := srv.UpdateScriptAll(&script, &scene)
	if err != nil {
		log.Error("UpdateScriptAll err (%v)", err)
		return
	}
	var resMap = make(map[string]string)
	resMap["update_result"] = updateResult
	c.JSON(resMap, err)
}

func runTimeCheck(c *bm.Context) {
	// 压测时间段检查
	userName, _ := c.Request.Cookie("username")
	// 暂时去掉压测时间
	//if !srv.CheckRunTime() {
	//	if !srv.CheckRunPermission(userName.Value) {
	//		c.JSON("Non-working time", ecode.MelloiRunNotInTime)
	//		return
	//	}
	//}
	if !srv.CheckRunPermission(userName.Value) {
		c.JSON("Non-working time", ecode.MelloiRunNotInTime)
		return
	}
	c.JSON(nil, nil)
}

func urlCheck(c *bm.Context) {
	// 检查url是否包含json串
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.URLCheck(&script))
}

// 增加定时器配置
func addTimer(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.AddTimer(&script))
}
