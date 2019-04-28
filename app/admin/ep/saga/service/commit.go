package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/utils"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

// QueryProjectCommit query project commit info according to project id.
func (s *Service) QueryProjectCommit(c context.Context, req *model.ProjectDataReq) (resp *model.ProjectDataResp, err error) {

	if resp, err = s.QueryProject(c, "commit", req); err != nil {
		return
	}
	return
}

// QueryTeamCommit query team commit info according to department and business
func (s *Service) QueryTeamCommit(c context.Context, req *model.TeamDataRequest) (resp *model.TeamDataResp, err error) {

	if resp, err = s.QueryTeam(c, "commit", req); err != nil {
		return
	}
	return
}

// QueryCommit query commit info according to department、 business and time.
func (s *Service) QueryCommit(c context.Context, req *model.CommitRequest) (resp *model.CommitResp, err error) {
	var (
		layout        = "2006-01-02"
		projectInfo   []*model.ProjectInfo
		reqProject    = &model.ProjectInfoRequest{}
		ProjectCommit []*model.ProjectCommit
		respCommit    *gitlab.Response
		since         time.Time
		until         time.Time
		commitNum     int
	)

	if len(req.Department) <= 0 && len(req.Business) <= 0 {
		log.Warn("query department and business are empty!")
		return
	}

	reqProject.Department = req.Department
	reqProject.Business = req.Business
	reqProject.Username = req.Username
	if _, projectInfo, err = s.dao.QueryProjectInfo(false, reqProject); err != nil {
		return
	}

	if len(projectInfo) <= 0 {
		log.Warn("Found no project!")
		return
	}

	//since, err = time.Parse("2006-01-02 15:04:05", "2018-08-13 00:00:00")
	if since, err = time.ParseInLocation(layout, req.Since, time.Local); err != nil {
		return
	}
	if until, err = time.ParseInLocation(layout, req.Until, time.Local); err != nil {
		return
	}

	log.Info("query commit start!")
	for _, project := range projectInfo {

		if _, respCommit, err = s.gitlab.ListProjectCommit(project.ProjectID, 1, &since, &until); err != nil {
			return
		}
		//log.Info("query: %s, result: %+v", project.Name, respCommit)

		CommitPer := &model.ProjectCommit{
			ProjectID: project.ProjectID,
			Name:      project.Name,
			CommitNum: respCommit.TotalItems,
		}
		ProjectCommit = append(ProjectCommit, CommitPer)
		commitNum = commitNum + respCommit.TotalItems
	}
	log.Info("query commit end!")

	resp = &model.CommitResp{
		Total:         commitNum,
		ProjectCommit: ProjectCommit,
	}
	return
}

/*-------------------------------------- sync commit ----------------------------------------*/

// SyncProjectCommit ...
func (s *Service) SyncProjectCommit(projectID int) (result *model.SyncResult, err error) {
	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		commits     []*gitlab.Commit
		resp        *gitlab.Response
		since       *time.Time
		until       *time.Time
		projectInfo *model.ProjectInfo
	)
	result = &model.SyncResult{}

	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}

	if !syncAllTime {
		since, until = utils.CalSyncTime()
	}
	log.Info("sync project(%d) commit time since: %v, until: %v", projectID, since, until)

	for page := 1; ; page++ {
		result.TotalPage++
		if commits, resp, err = s.gitlab.ListProjectCommit(projectID, page, since, until); err != nil {
			return
		}

		for _, commit := range commits {
			var (
				statsAdditions int
				statsDeletions int
				parentIDs      string
				commitStatus   string
			)

			if commit.Stats != nil {
				statsAdditions = commit.Stats.Additions
				statsDeletions = commit.Stats.Deletions
			}
			if commit.Status != nil {
				commitStatusByte, _ := json.Marshal(commit.Status)
				commitStatus = string(commitStatusByte)
			}
			parentIDsByte, _ := json.Marshal(commit.ParentIDs)
			parentIDs = string(parentIDsByte)

			commitDB := &model.StatisticsCommits{
				CommitID:       commit.ID,
				ProjectID:      projectID,
				ProjectName:    projectInfo.Name,
				ShortID:        commit.ShortID,
				Title:          commit.Title,
				AuthorName:     commit.AuthorName,
				AuthoredDate:   commit.AuthoredDate,
				CommitterName:  commit.CommitterName,
				CommittedDate:  commit.CommittedDate,
				CreatedAt:      commit.CreatedAt,
				Message:        commit.Message,
				ParentIDs:      parentIDs,
				StatsAdditions: statsAdditions,
				StatsDeletions: statsDeletions,
				Status:         commitStatus,
			}
			if len(commitDB.Message) > model.MessageMaxLen {
				commitDB.Message = commitDB.Message[0 : model.MessageMaxLen-1]
			}

			if err = s.SaveDatabaseCommit(commitDB); err != nil {
				log.Error("Commit Save Database err: projectID(%d), commitID(%s)", projectID, commit.ID)
				err = nil

				errData := &model.FailData{
					ChildIDStr: commit.ID,
				}
				result.FailData = append(result.FailData, errData)

				continue
			}
			result.TotalNum++
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return
}

// SaveDatabaseCommit ...
func (s *Service) SaveDatabaseCommit(commitDB *model.StatisticsCommits) (err error) {
	var total int

	if total, err = s.dao.HasCommit(commitDB.ProjectID, commitDB.CommitID); err != nil {
		log.Error("SaveDatabaseCommit HasCommit(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		if err = s.dao.UpdateCommit(commitDB.ProjectID, commitDB.CommitID, commitDB); err != nil {
			if strings.Contains(err.Error(), model.DatabaseErrorText) {
				commitDB.Title = strconv.QuoteToASCII(commitDB.Title)
				commitDB.Message = strconv.QuoteToASCII(commitDB.Message)
				commitDB.Title = utils.Unicode2Chinese(commitDB.Title)
				commitDB.Message = utils.Unicode2Chinese(commitDB.Message)
			}
			if err = s.dao.UpdateCommit(commitDB.ProjectID, commitDB.CommitID, commitDB); err != nil {
				log.Error("SaveDatabaseCommit UpdateCommit err(%+v)", err)
				return
			}
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseCommit commit has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateCommit(commitDB); err != nil {
		if strings.Contains(err.Error(), model.DatabaseErrorText) {
			commitDB.Title = strconv.QuoteToASCII(commitDB.Title)
			commitDB.Message = strconv.QuoteToASCII(commitDB.Message)
			commitDB.Title = utils.Unicode2Chinese(commitDB.Title)
			commitDB.Message = utils.Unicode2Chinese(commitDB.Message)
		}
		if err = s.dao.CreateCommit(commitDB); err != nil {
			log.Error("SaveDatabaseCommit CreateCommit err(%+v)", err)
			return
		}
	}

	return
}
