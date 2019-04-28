package service

import (
	"context"
	"fmt"

	"go-common/app/tool/saga/model"
	"go-common/library/log"
)

// MergeRequest ...
func (s *Service) MergeRequest(c context.Context, event *model.HookMR) (err error) {
	var (
		ok           bool
		repo         *model.Repo
		url          = event.Project.GitSSHURL
		projID       = event.Project.ID
		mrIID        = int(event.ObjectAttributes.IID)
		sourceBranch = event.ObjectAttributes.SourceBranch
		targetBranch = event.ObjectAttributes.TargetBranch
		wip          = event.ObjectAttributes.WorkInProgress
	)

	log.Info("Hook MergeRequest state: %s, projID: %d, mrid: %d", event.ObjectAttributes.State, projID, mrIID)
	if ok, repo = s.findRepo(url); !ok {
		return
	}

	if event.ObjectAttributes.State == model.MRStateOpened {
		if wip {
			log.Info("MergeRequest mr is wip, project ID: [%d], branch: [%s]", event.Project.ID, sourceBranch)
			return
		}
		if !s.validTargetBranch(targetBranch, repo) {
			s.gitlab.CreateMRNote(projID, mrIID, fmt.Sprintf("<pre>警告：目标分支 %s 不在SAGA白名单中，此MR不会触发SAGA行为！</pre>", targetBranch))
			return
		}

		authBranch := s.cmd.GetAuthBranch(targetBranch, repo.Config.AuthBranches)
		if err = s.ShowMRRoleInfo(c, authBranch, repo, event); err != nil {
			return
		}
	} else {
		if event.ObjectAttributes.State == model.MRStateMerged {
			log.Info("MergeRequest.UpdateContributor:%d,%s,%s,%+v", projID, sourceBranch, targetBranch, repo.Config.AuthBranches)
			if err = s.cmd.UpdateContributor(projID, mrIID, sourceBranch, targetBranch, repo.Config.AuthBranches); err != nil {
				return
			}
		}
		if err = s.d.DeleteReportStatus(c, projID, mrIID); err != nil {
			return
		}
	}
	return
}

// ShowMRRoleInfo ...
func (s *Service) ShowMRRoleInfo(c context.Context, authBranch string, repo *model.Repo, event *model.HookMR) (err error) {
	var (
		projID = event.Project.ID
		mrIID  = int(event.ObjectAttributes.IID)
		result bool
	)

	if result, err = s.d.ReportStatus(c, projID, mrIID); err != nil {
		return
	}

	log.Info("MergeRequest whether create note auth info: %v", result)
	if !result {
		if len(repo.Config.SuperAuthUsers) > 0 {
			if err = s.ReportSuperRoleInfo(c, repo.Config.SuperAuthUsers, event); err != nil {
				return
			}
		} else {
			if err = s.ReportMRRoleInfo(c, authBranch, event); err != nil {
				return
			}
		}

		if err = s.d.SetReportStatus(c, projID, mrIID, true); err != nil {
			return
		}
	}
	return
}

// ReportMRRoleInfo ...
func (s *Service) ReportMRRoleInfo(c context.Context, authBranch string, event *model.HookMR) (err error) {
	var (
		projID     = event.Project.ID
		mrIID      = int(event.ObjectAttributes.IID)
		pathOwners []model.RequireReviewFolder
	)

	if pathOwners, err = s.cmd.GetAllPathAuth(projID, mrIID, authBranch); err != nil {
		return
	}

	authInfo := ""
	authInfo = fmt.Sprintf("<pre>SAGA权限信息提示，请review: %s</pre>", event.ObjectAttributes.Title)
	authInfo += "\n"
	for _, os := range pathOwners {
		if len(os.Owners) > 0 {
			//authInfo += "----- PATH: " + os.Folder + "；OWNER: "
			authInfo += fmt.Sprintf("+ PATH: %s；OWNER: ", os.Folder)
			for _, o := range os.Owners {
				if o == "all" {
					authInfo += "所有人" + " 或 "
				} else {
					authInfo += "@" + o + " 或 "
				}
			}
			authInfo = authInfo[:len(authInfo)-len(" 或 ")]
		}
		if len(os.Reviewers) > 0 {
			authInfo += "；REVIEWER: "
			for _, o := range os.Reviewers {
				if o == "all" {
					authInfo += "所有人" + " 与 "
				} else {
					authInfo += "@" + o + " 与 "
				}
			}
			authInfo = authInfo[:len(authInfo)-len(" 与 ")]
		}
		authInfo += "\n\n"
	}

	if _, err = s.gitlab.CreateMRNote(projID, mrIID, authInfo); err != nil {
		return
	}
	return
}

// ReportSuperRoleInfo ...
func (s *Service) ReportSuperRoleInfo(c context.Context, superUsers []string, event *model.HookMR) (err error) {
	var (
		projID = event.Project.ID
		mrIID  = int(event.ObjectAttributes.IID)
	)

	authInfo := fmt.Sprintf("<pre>SAGA权限信息提示，已配置超级权限用户，请review: %s</pre>", event.ObjectAttributes.Title)
	authInfo += "\n"
	authInfo += fmt.Sprintf("+ SUPERMAN: ")
	for _, user := range superUsers {
		authInfo += "@" + user + " 或 "
	}
	authInfo = authInfo[:len(authInfo)-len(" 或 ")]

	if _, err = s.gitlab.CreateMRNote(projID, mrIID, authInfo); err != nil {
		return
	}
	return
}
