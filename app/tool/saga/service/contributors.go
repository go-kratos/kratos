package service

import (
	"context"
	"runtime/debug"
	"strings"

	"go-common/app/tool/saga/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// HandleBuildContributors ...
func (s *Service) HandleBuildContributors(c context.Context, repo *model.RepoInfo) (err error) {
	if strings.TrimSpace(repo.Group) == "" {
		err = errors.Errorf("repo Group is not valid")
		return
	}
	if strings.TrimSpace(repo.Name) == "" {
		err = errors.Errorf("repo Name is not valid")
		return
	}
	if strings.TrimSpace(repo.Branch) == "" {
		err = errors.Errorf("repo Branch is not valid")
		return
	}
	go func() {
		defer func() {
			if x := recover(); x != nil {
				log.Error("BuildContributor: %+v %s", x, debug.Stack())
			}
		}()
		if err = s.cmd.BuildContributor(repo); err != nil {
			log.Error("BuildContributor %+v", err)
		}
	}()
	return
}

// BuildContributors ...
func (s *Service) BuildContributors(repos []*model.Repo) (err error) {
	var (
		repo   *model.Repo
		branch string
	)

	log.Info("BuildContributors start ...")
	for _, repo = range repos {
		for _, branch = range repo.Config.AuthBranches {
			repoInfo := &model.RepoInfo{
				Group:  repo.Config.Group,
				Name:   repo.Config.Name,
				Branch: branch,
			}
			log.Info("BuildContributors project [%s], group [%s], Name [%s], branch [%s]", repo.Config.URL, repo.Config.Group, repo.Config.Name, branch)
			if err = s.HandleBuildContributors(context.TODO(), repoInfo); err != nil {
				log.Error("BuildContributors err (%+v)", err)
			}
		}
	}
	return
}
