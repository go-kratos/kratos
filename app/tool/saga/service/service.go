package service

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/tool/saga/conf"
	"go-common/app/tool/saga/dao"
	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/command"
	"go-common/app/tool/saga/service/gitlab"
	"go-common/library/log"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

// Service biz service def.
type Service struct {
	missch     chan func()
	d          *dao.Dao
	gitlab     *gitlab.Gitlab
	gitRepoMap map[string]*model.Repo // map[repoName]*repo
	cmd        *command.Command
	cron       *cron.Cron
}

// New a DirService and return.
func New() (s *Service) {
	s = &Service{
		d:      dao.New(),
		missch: make(chan func(), 10240),
	}
	// init gitlab client
	s.gitlab = gitlab.New(conf.Conf.Property.Gitlab.API, conf.Conf.Property.Gitlab.Token)
	s.cmd = command.New(s.d, s.gitlab)
	s.cmd.Registers()
	go s.cmd.ListenTask()
	s.loadRepos(false)
	go s.updateproc()
	// start cron
	s.cron = cron.New()
	if err := s.cron.AddFunc(conf.Conf.Property.SyncContact.CheckCron, s.synccontactsproc); err != nil {
		panic(err)
	}
	s.cron.Start()
	//
	return
}

func (s *Service) validTargetBranch(targetBranch string, gitRepo *model.Repo) bool {
	for _, r := range gitRepo.Config.TargetBranchRegexes {
		if r.MatchString(targetBranch) {
			return true
		}
	}
	/*for _, r := range gitRepo.Config.TargetBranches {
		if r == targetBranch {
			return true
		}
	}*/
	return false
}

func (s *Service) loadRepos(reload bool) {
	var (
		repo *model.Repo
		ok   bool
	)
	// init code repo
	if s.gitRepoMap == nil {
		s.gitRepoMap = make(map[string]*model.Repo)
	}

	webHookRepos := make([]*model.Repo, 0)
	authRepos := make([]*model.Repo, 0)
	for _, r := range conf.Conf.Property.Repos {
		if repo, ok = s.gitRepoMap[r.GName]; !ok {
			repo = &model.Repo{
				Config: r,
			}
			s.gitRepoMap[r.GName] = repo
			webHookRepos = append(webHookRepos, repo)
			authRepos = append(authRepos, repo)
		} else {
			if repo.AuthUpdate(r) {
				authRepos = append(authRepos, repo)
			}
			if repo.WebHookUpdate(r) {
				webHookRepos = append(webHookRepos, repo)
			}
			if repo.Update(r) {
				s.gitRepoMap[r.GName] = repo
			}
		}
	}

	if reload {
		s.BuildContributors(authRepos)
	}
	if err := s.gitlab.AuditProjects(webHookRepos, conf.Conf.Property.WebHooks); err != nil {
		log.Error("loadRepos err (%+v)", err)
	}
}

func (s *Service) findRepo(gitURL string) (ok bool, repo *model.Repo) {
	for _, r := range conf.Conf.Property.Repos {
		if strings.EqualFold(r.URL, gitURL) {
			repo = &model.Repo{
				Config: r,
			}
			ok = true
			return
		}
	}
	return
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.d.Ping(c)
}

// Wait wait all closed.
func (s *Service) Wait() {
}

// Close close all dao.
func (s *Service) Close() {
	s.d.Close()
}

func (s *Service) updateproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("updateproc panic(%v)", errors.WithStack(fmt.Errorf("%v", x)))
			go s.updateproc()
			log.Info("updateproc recover")
		}
	}()
	for range conf.ReloadEvents() {
		log.Info("DirService reload")
		s.loadRepos(true)
	}
}
