package http

import (
	"go-common/app/admin/ep/melloi/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// AddAndExecuScene Add And Execu Scene
func AddAndExecuScene(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	c.JSON(srv.AddAndExecuScene(c, scene, cookie))
}

func queryDraft(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryDraft(&scene))
}

func updateScene(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.UpdateScene(&scene))
}

func addScene(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	id, _ := srv.AddScene(&scene)
	c.JSON(id, nil)
}

func saveScene(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.SaveScene(&scene))
}

func saveOrder(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	var req model.SaveOrderReq
	if err := c.BindWith(&req, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.SaveOrder(req, &scene))
}

func queryRelation(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryRelation(&script))
}

func queryAPI(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryAPI(&scene))
}

func deleteAPI(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.DeleteAPI(&script))
}

func addConfig(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.AddConfig(&script))
}

func queryTree(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryTree(&script))
}

func queryScenes(c *bm.Context) {
	qsr := model.QuerySceneRequest{}
	if err := c.BindWith(&qsr, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	if err := qsr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	sessionID, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}

	res, err := srv.QueryScenesByPage(c, sessionID.Value, &qsr)
	if err != nil {
		return
	}
	c.JSON(res, err)
}

func doScenePtest(c *bm.Context) {
	ptestScene := model.DoPtestSceneParam{}
	if err := c.BindWith(&ptestScene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	c.JSON(srv.DoScenePtest(c, ptestScene, false, cookie))
}

func doScenePtestBatch(c *bm.Context) {
	ptestScenes := model.DoPtestSceneParams{}
	if err := c.BindWith(&ptestScenes, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	c.JSON(nil, srv.DoScenePtestBatch(c, ptestScenes, cookie))
}
func queryExistAPI(c *bm.Context) {
	apiInfoReq := model.APIInfoRequest{}
	if err := c.BindWith(&apiInfoReq, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	sessionID, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryExistAPI(c, sessionID.Value, &apiInfoReq))
}

func queryPreview(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
	}
	c.JSON(srv.QueryPreview(&script))
}

func queryParams(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
	}
	res, _, err := srv.QueryParams(&script)
	c.JSON(res, err)
}

func updateBindScene(c *bm.Context) {
	bindScene := model.BindScene{}
	if err := c.BindWith(&bindScene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.UpdateBindScene(&bindScene))
}

func queryDrawRelation(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.Form); err != nil {
		c.JSON(nil, err)
	}
	res, _, err := srv.QueryDrawRelation(&scene)
	c.JSON(res, err)
}

func addSceneScript(c *bm.Context) {
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
	cookie := c.Request.Header.Get("Cookie")
	c.JSON(srv.AddAndExcuScript(c, &script, cookie, &scene, false, false))
}

func deleteDraft(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.DeleteDraft(&scene))
}

func queryConfig(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryConfig(&script))
}

func deleteScene(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.DeleteScene(&scene))
}

func copyScene(c *bm.Context) {
	scene := model.Scene{}
	if err := c.BindWith(&scene, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	c.JSON(srv.CopyScene(c, &scene, cookie))
}

func queryFusing(c *bm.Context) {
	script := model.Script{}
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryFusing(&script))
}
