package service

import (
	"context"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

//QueryProjectRunners query project runners info according to project id
func (s *Service) QueryProjectRunners(c context.Context, req *model.ProjectDataReq) (resp []*gitlab.Runner, err error) {
	var (
		runners  []*gitlab.Runner
		response *gitlab.Response
	)

	for page := 1; ; page++ {
		if runners, response, err = s.gitlab.ListProjectRunners(req.ProjectID, page); err != nil {
			return
		}
		resp = append(resp, runners...)
		if response.NextPage == 0 {
			break
		}
	}
	return
}

/*-------------------------------------- sync runner ----------------------------------------*/

// SyncAllRunners ...
func (s *Service) SyncAllRunners(projectID int) (totalPage, totalNum int, err error) {
	var (
		runners     []*gitlab.Runner
		resp        *gitlab.Response
		projectInfo *model.ProjectInfo
	)
	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}

	for page := 1; ; page++ {
		totalPage++
		if runners, resp, err = s.gitlab.ListProjectRunners(projectID, page); err != nil {
			return
		}
		for _, runner := range runners {
			var (
				ipAddress string
			)
			//ipAddress = runner.IPAddress.String()
			runnerDB := &model.StatisticsRunners{
				ProjectID:   projectID,
				ProjectName: projectInfo.Name,
				RunnerID:    runner.ID,
				Description: runner.Description,
				Active:      runner.Active,
				IsShared:    runner.IsShared,
				IPAddress:   ipAddress,
				Name:        runner.Name,
				Online:      runner.Online,
				Status:      runner.Status,
				Token:       runner.Token,
			}

			if err = s.SaveDatabaseRunner(runnerDB); err != nil {
				log.Error("runner Save Database err: projectID(%d), RunnerID(%d)", projectID, runner.ID)
				err = nil
				continue
			}
			totalNum++
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return
}

// SaveDatabaseRunner ...
func (s *Service) SaveDatabaseRunner(runnerDB *model.StatisticsRunners) (err error) {
	var total int

	if total, err = s.dao.HasRunner(runnerDB.ProjectID, runnerDB.RunnerID); err != nil {
		log.Error("SaveDatabaseRunner HasRunner(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		if err = s.dao.UpdateRunner(runnerDB.ProjectID, runnerDB.RunnerID, runnerDB); err != nil {
			log.Error("SaveDatabaseRunner UpdateRunner(%+v)", err)
			return
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseRunner commit has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateRunner(runnerDB); err != nil {
		log.Error("SaveDatabaseRunner CreateRunner(%+v)", err)
		return
	}

	return
}
