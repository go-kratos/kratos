package http

import (
	"encoding/json"
	"io/ioutil"

	"go-common/app/tool/saga/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func gitlabComment(c *bm.Context) {
	var (
		bytes       []byte
		err         error
		hookComment = &model.HookComment{}
	)
	if bytes, err = ioutil.ReadAll(c.Request.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal(bytes, hookComment); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if hookComment == nil || hookComment.User == nil || hookComment.ObjectAttributes == nil || hookComment.MergeRequest == nil {
		log.Error("hookComment event not standard")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("Got new Gitlab Comment Event kind(%s) attr(%+v) user(%+v)", hookComment.ObjectKind, hookComment.ObjectAttributes, hookComment.User)
	c.JSON(nil, svc.HandleGitlabComment(c, hookComment))
}

func gitlabPipeline(c *bm.Context) {
	var (
		err          error
		hookPipeline = &model.HookPipeline{}
	)

	if err = c.BindWith(hookPipeline, binding.JSON); err != nil {
		return
	}
	if hookPipeline == nil || hookPipeline.User == nil || hookPipeline.ObjectAttributes == nil || hookPipeline.Project == nil {
		log.Error("hookPipeline event not standard")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.PipelineChanged(c, hookPipeline))
}

func gitlabMR(c *bm.Context) {
	var (
		err    error
		hookMR = &model.HookMR{}
	)
	if err = c.BindWith(hookMR, binding.JSON); err != nil {
		return
	}
	if hookMR == nil || hookMR.ObjectAttributes == nil || hookMR.Project == nil {
		log.Error("hookMR event not standard")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.MergeRequest(c, hookMR))
}
