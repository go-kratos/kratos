package service

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/tool/saga/model"
	"go-common/library/log"
)

// HandleGitlabComment handle comment webhook
func (s *Service) HandleGitlabComment(c context.Context, event *model.HookComment) (err error) {
	var (
		comment      = strings.TrimSpace(event.ObjectAttributes.Note)
		projName     = event.MergeRequest.Source.Name
		targetBranch = event.MergeRequest.TargetBranch
		gitRepo      *model.Repo
		ok           bool
	)
	if event.ObjectAttributes.NoteableType != model.HookCommentTypeMR {
		log.Info("Comment hook noteableType [%s] ignore", event.ObjectAttributes.NoteableType)
		return
	}
	if event.MergeRequest.State == model.MRStateMerged {
		log.Info("Gitlab MR [%d] has been merged", event.MergeRequest.ID)
		return
	}
	if gitRepo, ok = s.gitRepoMap[projName]; !ok {
		log.Info("Gitlab MR (%d) unknown projName (%s)", event.MergeRequest.ID, projName)
		return
	}
	if !s.validTargetBranch(targetBranch, gitRepo) {
		log.Info("Target branch (%s) is not in white list, won't serve comment!", targetBranch)
		return
	}
	fmt.Println("HandleGitlabComment:", event)
	if err = s.cmd.Exec(c, comment, event, gitRepo); err != nil {
		log.Error("Command Exec cmd: %s err: (%+v)", comment, err)
		return
	}
	return
}
