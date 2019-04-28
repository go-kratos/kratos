package command

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/notification"
	"go-common/library/log"
)

func (c *Command) runPlusOne(ctx context.Context, event *model.HookComment, repo *model.Repo) (err error) {
	var (
		author       string
		url          = event.ObjectAttributes.URL
		commit       = event.MergeRequest.LastCommit.ID
		reviewer     = event.User.UserName
		authorID     = int(event.MergeRequest.AuthorID)
		sourceBranch = event.MergeRequest.SourceBranch
		targetBranch = event.MergeRequest.TargetBranch
		wip          = event.MergeRequest.WorkInProgress
	)
	log.Info("runPlusOne start ...")
	if wip {
		c.gitlab.CreateMRNote(event.Project.ID, int(event.MergeRequest.IID), fmt.Sprintf("<pre>警告：当前MR处于WIP状态，请待开发结束后再review！</pre>"))
		return
	}
	if author, err = c.gitlab.UserName(authorID); err != nil {
		log.Error("%+v", err)
		return
	}

	log.Info("runPlusOne notification author: %s", author)
	if author != "" {
		go func() {
			notification.MailPlusOne(author, reviewer, url, commit, sourceBranch, targetBranch)
		}()
		go func() {
			notification.WechatPlusOne(c.dao, author, reviewer, url, commit, sourceBranch, targetBranch)
		}()
	}
	return
}

func reviewedOwner(owners []string, reviewedUsers []string, username string) (isowner bool, reviewed bool) {
	for _, owner := range owners {
		if strings.EqualFold(owner, username) || strings.EqualFold(owner, "all") {
			return true, true
		}
	}
	for _, owner := range owners {
		for _, user := range reviewedUsers {
			if strings.EqualFold(user, owner) {
				return false, true
			}
		}
	}
	return false, false
}

func (c *Command) reviewedUsers(projID int, mrIID int) (reviewedUsers []string, err error) {
	var awardUsers []string
	if reviewedUsers, err = c.gitlab.PlusUsernames(projID, mrIID); err != nil {
		return
	}
	if awardUsers, err = c.gitlab.AwardEmojiUsernames(projID, mrIID); err != nil {
		return
	}
OUTER:
	for _, au := range awardUsers {
		for _, mu := range reviewedUsers {
			if au == mu {
				continue OUTER
			}
		}
		reviewedUsers = append(reviewedUsers, au)
	}
	return
}

func reviewedNum(reviewers []string, reviewedUsers []string) (num int) {
	for _, reviewer := range reviewers {
		for _, user := range reviewedUsers {
			if strings.EqualFold(user, reviewer) || strings.EqualFold(reviewer, "all") {
				num++
			}
		}
	}
	return
}
