package service

import (
	"context"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
)

// FavoriteProjects list user's favorite projects
func (s *Service) FavoriteProjects(c context.Context, req *model.Pagination, userName string) (resp *model.FavoriteProjectsResp, err error) {
	var (
		favorite    *model.ProjectFavorite
		favorites   []*model.ProjectFavorite
		projectInfo *model.ProjectInfo
	)
	resp = &model.FavoriteProjectsResp{}

	if favorites, err = s.dao.FavoriteProjects(userName); err != nil {
		return
	}

	for _, favorite = range favorites {
		if projectInfo, err = s.dao.ProjectInfoByID(favorite.ProjID); err != nil {
			return
		}

		myProject := &model.MyProjectInfo{
			ProjectInfo: projectInfo,
		}
		myProject.Star = true
		resp.Projects = append(resp.Projects, myProject)
	}
	resp.PageSize = req.PageSize
	resp.PageNum = req.PageNum
	resp.Total = len(favorites)
	return
}

// EditFavorite edit user's favorites, star/unstar
func (s *Service) EditFavorite(c context.Context, req *model.EditFavoriteReq, userName string) (resp *model.EmptyResp, err error) {
	var (
		projID    = req.ProjID
		favorites []*model.ProjectFavorite
		exist     bool
	)
	resp = &model.EmptyResp{}

	if favorites, err = s.dao.FavoriteProjects(userName); err != nil {
		return
	}
	if req.Star {
		if !inFavorites(favorites, projID) {
			if exist, err = s.dao.ProjectExist(projID); err != nil {
				return
			}
			if exist {
				if err = s.dao.AddFavorite(userName, projID); err != nil {
					return
				}
			}
		}
	} else {
		if inFavorites(favorites, projID) {
			if err = s.dao.DelFavorite(userName, projID); err != nil {
				return
			}
		}
	}

	return
}

// inFavorites ...
func inFavorites(favorites []*model.ProjectFavorite, projID int) bool {
	var (
		f *model.ProjectFavorite
	)
	for _, f = range favorites {
		if projID == f.ProjID {
			return true
		}
	}

	return false
}

// QueryCommonProjects ...
func (s *Service) QueryCommonProjects(c context.Context) (result []string, err error) {
	var (
		projectInfo *model.ProjectInfo
	)
	ids := conf.Conf.Property.DefaultProject.ProjectIDs
	for _, id := range ids {
		if projectInfo, err = s.dao.ProjectInfoByID(id); err != nil {
			return
		}
		result = append(result, projectInfo.Name)
	}
	return
}

// QueryProjectInfo query project info.
func (s *Service) QueryProjectInfo(c context.Context, req *model.ProjectInfoRequest) (resp *model.ProjectInfoResp, err error) {
	var (
		projectInfo     []*model.ProjectInfo
		project         *model.ProjectInfo
		total           int
		saga            int
		runner          int
		sagaScale       int
		runnerScale     int
		mproject        *model.MyProjectInfo
		favorites       []*model.ProjectFavorite
		projectInfoResp []*model.MyProjectInfo
	)
	userName := req.Username
	if total, projectInfo, err = s.dao.QueryProjectInfo(true, req); err != nil {
		return
	}
	if favorites, err = s.dao.FavoriteProjects(userName); err != nil {
		return
	}
	for _, project = range projectInfo {
		mproject = &model.MyProjectInfo{
			ProjectInfo: project,
		}
		if !inFavorites(favorites, project.ProjectID) {
			mproject.Star = false
		} else {
			mproject.Star = true
		}
		projectInfoResp = append(projectInfoResp, mproject)
	}

	if saga, err = s.dao.QueryConfigInfo(req.Name, req.Department, req.Business, "saga"); err != nil {
		return
	}
	if runner, err = s.dao.QueryConfigInfo(req.Name, req.Department, req.Business, "runner"); err != nil {
		return
	}

	if total != 0 {
		sagaScale = saga * 100 / total
		runnerScale = runner * 100 / total
	}

	resp = &model.ProjectInfoResp{
		PageNum:     req.PageNum,
		PageSize:    req.PageSize,
		Total:       total,
		Saga:        saga,
		Runner:      runner,
		SagaScale:   sagaScale,
		RunnerScale: runnerScale,
		ProjectInfo: projectInfoResp,
	}
	return
}
