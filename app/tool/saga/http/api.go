package http

import (
	"go-common/app/tool/saga/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func buildContributors(c *bm.Context) {
	var (
		err  error
		repo = &model.RepoInfo{}
	)

	if err = c.BindWith(repo, binding.JSON); err != nil {
		log.Error("BindWith error(%v)", err)
		return
	}
	c.JSON(nil, svc.HandleBuildContributors(c, repo))
}
