package service

import (
	"context"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"
)

// UpdateBuglyVersion update Bugly Version.
func (s *Service) UpdateBuglyVersion(c context.Context, username string, req *model.AddVersionRequest) (rep map[string]interface{}, err error) {
	var (
		bugVersionID int64
		buglyProject *model.BuglyProject
	)

	// check project
	if buglyProject, err = s.dao.QueryBuglyProject(req.BuglyProjectID); err != nil {
		return
	}

	if buglyProject.ID == 0 {
		err = ecode.NothingFound
		return
	}

	if req.ID > 0 {
		// update
		var tmpVersion *model.BuglyVersion

		if tmpVersion, err = s.dao.QueryBuglyVersion(req.ID); err != nil {
			return
		}

		if tmpVersion.ID == 0 {
			err = ecode.NothingFound
			return
		}

		tmpVersion.BuglyProjectID = req.BuglyProjectID
		tmpVersion.Version = req.Version
		tmpVersion.Action = req.Action

		if err = s.dao.UpdateBuglyVersion(tmpVersion); err != nil {
			return
		}
		bugVersionID = tmpVersion.ID

	} else {
		// add
		var tmpVersion *model.BuglyVersion

		//check name
		if tmpVersion, err = s.dao.QueryBuglyVersionByVersion(req.Version); err != nil {
			return
		}

		if tmpVersion.ID > 0 {
			err = ecode.MartheDuplicateErr
			return
		}

		buglyVersion := &model.BuglyVersion{
			ID:             req.ID,
			Version:        req.Version,
			BuglyProjectID: req.BuglyProjectID,
			Action:         req.Action,
			TaskStatus:     model.BuglyVersionTaskStatusReady,
			UpdateBy:       username,
		}

		if err = s.dao.InsertBuglyVersion(buglyVersion); err != nil {
			return
		}
		bugVersionID = buglyVersion.ID
	}

	rep = make(map[string]interface{})
	rep["bug_version_id"] = bugVersionID
	return
}

// QueryBuglyVersions Query Bugly Versions.
func (s *Service) QueryBuglyVersions(c context.Context, req *model.QueryBuglyVersionRequest) (rep *model.PaginateBuglyProjectVersions, err error) {
	var (
		total                int64
		buglyProjectVersions []*model.BuglyProjectVersion
	)
	if total, buglyProjectVersions, err = s.dao.FindBuglyProjectVersions(req); err != nil {
		return
	}
	rep = &model.PaginateBuglyProjectVersions{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		BuglyProjectVersions: buglyProjectVersions,
	}
	return
}

// BuglyVersionAndProjectList Bugly Version and project List.
func (s *Service) BuglyVersionAndProjectList(c context.Context) (rep map[string]interface{}, err error) {
	var (
		versionList []string
		projectList []string
	)

	rep = make(map[string]interface{})

	if versionList, err = s.dao.QueryBuglyVersionList(); err != nil {
		return
	}

	if projectList, err = s.dao.QueryBuglyProjectList(); err != nil {
		return
	}

	rep["versions"] = versionList

	rep["projects"] = projectList

	return
}

// UpdateCookie update Cookie.
func (s *Service) UpdateCookie(c context.Context, username string, req *model.AddCookieRequest) (rep map[string]interface{}, err error) {
	var cookieID int64

	if req.ID > 0 {
		buglyCookie := &model.BuglyCookie{
			ID:        req.ID,
			QQAccount: req.QQAccount,
			Cookie:    req.Cookie,
			Token:     req.Token,
			Status:    req.Status,
			UpdateBy:  username,
		}
		if err = s.dao.UpdateCookie(buglyCookie); err != nil {
			return
		}
		cookieID = req.ID
	} else {
		buglyCookie := &model.BuglyCookie{
			QQAccount: req.QQAccount,
			Cookie:    req.Cookie,
			Token:     req.Token,
			Status:    req.Status,
			UpdateBy:  username,
		}
		if err = s.dao.InsertCookie(buglyCookie); err != nil {
			return
		}
		cookieID = buglyCookie.ID
	}

	rep = make(map[string]interface{})
	rep["cookie_id"] = cookieID
	return
}

// QueryCookies Add Cookie.
func (s *Service) QueryCookies(c context.Context, req *model.QueryBuglyCookiesRequest) (rep *model.PaginateBuglyCookies, err error) {
	var (
		total        int64
		buglyCookies []*model.BuglyCookie
	)
	if total, buglyCookies, err = s.dao.FindCookies(req); err != nil {
		return
	}
	rep = &model.PaginateBuglyCookies{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		BuglyCookies: buglyCookies,
	}
	return
}

// UpdateCookieStatus Update Cookie Status.
func (s *Service) UpdateCookieStatus(c context.Context, cookieID int64, status int) (err error) {
	err = s.dao.UpdateCookieStatus(cookieID, status)
	return
}

// GetEnableCookie Get Enable Cookie.
func (s *Service) GetEnableCookie() (buglyCookie *model.BuglyCookie, err error) {
	var buglyCookies []*model.BuglyCookie

	if buglyCookies, err = s.dao.QueryCookieByStatus(model.BuglyCookieStatusEnable); err != nil {
		return
	}

	/*for _, ele := range buglyCookies {
		if ele.UsageCount < s.c.Bugly.CookieUsageUpper {
			buglyCookie = ele
			err = s.dao.UpdateCookieUsageCount(buglyCookie.ID, buglyCookie.UsageCount+1)
			return
		}
	}*/

	if len(buglyCookies) > 0 {
		buglyCookie = buglyCookies[0]
		err = s.dao.UpdateCookieUsageCount(buglyCookie.ID, buglyCookie.UsageCount+1)
		return
	}

	//not found enable cookie
	s.DoWhenNoEnableCookie()
	err = ecode.MartheNoCookie
	return
}

// DisableCookie Disable Cookie.
func (s *Service) DisableCookie(c context.Context, cookieID int64) (err error) {
	err = s.dao.UpdateCookieStatus(cookieID, model.BuglyCookieStatusDisable)
	return
}

// DoWhenNoEnableCookie Do When No Enable Cookie.
func (s *Service) DoWhenNoEnableCookie() {
	// todo notice
	s.SendMail(s.c.Mail.NoticeOwner, "marthe has no enable cookie", "")
}

// QueryBuglyIssueRecords Query Bugly Issue Records
func (s *Service) QueryBuglyIssueRecords(c context.Context, req *model.QueryBuglyIssueRequest) (rep *model.PaginateBuglyIssues, err error) {
	var (
		total       int64
		buglyIssues []*model.BuglyIssue
	)
	if total, buglyIssues, err = s.dao.FindBuglyIssues(req); err != nil {
		return
	}
	rep = &model.PaginateBuglyIssues{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		BuglyIssues: buglyIssues,
	}
	return

}

// QueryBatchRuns Query Batch Runs
func (s *Service) QueryBatchRuns(c context.Context, req *model.QueryBuglyBatchRunsRequest) (rep *model.PaginateBuglyBatchRuns, err error) {
	var (
		total          int64
		buglyBatchRuns []*model.BuglyBatchRun
	)
	if total, buglyBatchRuns, err = s.dao.FindBuglyBatchRuns(req); err != nil {
		return
	}
	rep = &model.PaginateBuglyBatchRuns{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		BuglyBatchRuns: buglyBatchRuns,
	}
	return
}

// UpdateBuglyProject update bugly project.
func (s *Service) UpdateBuglyProject(c context.Context, username string, req *model.AddProjectRequest) (rep map[string]interface{}, err error) {
	var buglyProjectID int64

	if req.ID > 0 {
		buglyProject := &model.BuglyProject{
			ID:        req.ID,
			ProjectID: req.ProjectID,
			//ProjectName:   req.ProjectName,
			PlatformID:    req.PlatformID,
			ExceptionType: req.ExceptionType,
			UpdateBy:      username,
		}
		if err = s.dao.UpdateBuglyProject(buglyProject); err != nil {
			return
		}
		buglyProjectID = req.ID
	} else {
		var buglyProjectInDB *model.BuglyProject
		if buglyProjectInDB, err = s.dao.QueryBuglyProjectByName(req.ProjectName); err != nil {
			return
		}

		if buglyProjectInDB.ID > 0 {
			err = ecode.MartheDuplicateErr
			return
		}

		buglyProject := &model.BuglyProject{
			ProjectID:     req.ProjectID,
			ProjectName:   req.ProjectName,
			PlatformID:    req.PlatformID,
			ExceptionType: req.ExceptionType,
			UpdateBy:      username,
		}
		if err = s.dao.InsertBuglyProject(buglyProject); err != nil {
			return
		}
		buglyProjectID = buglyProject.ID
	}

	rep = make(map[string]interface{})
	rep["bugly_project_id"] = buglyProjectID
	return
}

// QueryBuglyProjects Query Bugly Project.
func (s *Service) QueryBuglyProjects(c context.Context, req *model.QueryBuglyProjectRequest) (rep *model.PaginateBuglyProjects, err error) {
	var (
		total         int64
		buglyProjects []*model.BuglyProject
	)
	if total, buglyProjects, err = s.dao.FindBuglyProjects(req); err != nil {
		return
	}
	rep = &model.PaginateBuglyProjects{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		BuglyProjects: buglyProjects,
	}
	return
}

// QueryBuglyProject Query Bugly Project.
func (s *Service) QueryBuglyProject(c context.Context, id int64) (buglyProject *model.BuglyProject, err error) {
	return s.dao.QueryBuglyProject(id)
}

// QueryAllBuglyProjects Query All Bugly Projects.
func (s *Service) QueryAllBuglyProjects(c context.Context) (buglyProjects []*model.BuglyProject, err error) {
	return s.dao.QueryAllBuglyProjects()
}

// QueryBuglyProjectVersions Query Bugly Project Versions.
func (s *Service) QueryBuglyProjectVersions(c context.Context, buglyProjectID int64) (rep map[string]interface{}, err error) {
	var (
		buglyProject  *model.BuglyProject
		buglyCookie   *model.BuglyCookie
		buglyVersions []*model.BugVersion
		versions      []string
	)

	defer func() {
		if err != nil && err == ecode.MartheCookieExpired {
			s.DisableCookie(c, buglyCookie.ID)
		}
	}()

	if buglyProject, err = s.dao.QueryBuglyProject(buglyProjectID); err != nil {
		return
	}

	if buglyProject.ID == 0 {
		err = ecode.NothingFound
		return
	}

	//get enable cookie
	if buglyCookie, err = s.GetEnableCookie(); err != nil {
		return
	}

	if buglyVersions, err = s.dao.BuglyVersion(c, buglyCookie, buglyProject.ProjectID, buglyProject.PlatformID); err != nil || buglyVersions == nil {
		return
	}

	for _, buglyVersion := range buglyVersions {
		versions = append(versions, buglyVersion.Name)
	}

	rep = make(map[string]interface{})
	rep["versions"] = versions
	return
}
