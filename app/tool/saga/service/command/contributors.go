package command

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/notification"
	"go-common/library/log"
)

type contributor struct {
	Owner    []string
	Author   []string
	Reviewer []string
}

func readContributor(content []byte) (c *contributor) {
	var (
		lines      []string
		lineStr    string
		curSection string
	)
	c = &contributor{}
	lines = strings.Split(string(content), "\n")
	for _, lineStr = range lines {
		if lineStr == "" {
			continue
		}
		if strings.Contains(strings.ToLower(lineStr), "owner") {
			curSection = "owner"
			continue
		}
		if strings.Contains(strings.ToLower(lineStr), "author") {
			curSection = "author"
			continue
		}
		if strings.Contains(strings.ToLower(lineStr), "reviewer") {
			curSection = "reviewer"
			continue
		}
		switch curSection {
		case "owner":
			c.Owner = append(c.Owner, strings.TrimSpace(lineStr))
		case "author":
			c.Author = append(c.Author, strings.TrimSpace(lineStr))
		case "reviewer":
			c.Reviewer = append(c.Reviewer, strings.TrimSpace(lineStr))
		}
	}
	return
}

// BuildContributor ...
func (c *Command) BuildContributor(repo *model.RepoInfo) (err error) {
	var (
		host   string
		token  string
		files  []string
		projID int
		branch = repo.Branch
	)
	if host, token, err = c.gitlab.HostToken(); err != nil {
		return
	}
	if files, err = c.dao.RepoFiles(context.TODO(), host, token, repo); err != nil {
		return
	}
	if projID, err = c.gitlab.ProjectID(fmt.Sprintf("git@%s:%s/%s.git", host, repo.Group, repo.Name)); err != nil {
		return
	}
	if err = c.SaveContributor(projID, branch, files, false); err != nil {
		return
	}
	return
}

func hasbranch(branch string, branchs []string) bool {
	for _, b := range branchs {
		if strings.EqualFold(b, branch) {
			return true
		}
	}
	return false
}

// UpdateContributor ...
func (c *Command) UpdateContributor(projID int, mrIID int, sourceBranch string, targetBranch string, authBranches []string) (err error) {
	var changeFiles []string
	var deleteFiles []string
	if !hasbranch(targetBranch, authBranches) {
		log.Info("UpdateContributor not authBranches: %d, %s, %+v", projID, targetBranch, authBranches)
		return
	}
	if changeFiles, deleteFiles, err = c.gitlab.MRChanges(projID, mrIID); err != nil {
		return
	}

	//log.Info("UpdateContributor: projID: %d, changeFiles: %+v, deleteFiles: %+v", projID, changeFiles, deleteFiles)
	if err = c.SaveContributor(projID, targetBranch, changeFiles, false); err != nil {
		return
	}
	if err = c.SaveContributor(projID, targetBranch, deleteFiles, true); err != nil {
		return
	}
	return
}

// SaveContributor ...
func (c *Command) SaveContributor(projID int, branch string, files []string, clean bool) (err error) {
	var (
		ctx = context.TODO()
		raw []byte
		cb  *contributor
	)
	for _, file := range files {
		if strings.EqualFold(filepath.Base(file), model.SagaContributorsName) {
			path := filepath.Dir(file)
			if clean {
				//delete redis key
				if err = c.dao.DeletePathAuthR(ctx, projID, branch, path); err != nil {
					log.Error("delete Auth Redis key error(%+v) project ID:%d, branch:%s, path:%s", err, projID, branch, path)
					err = nil
				}
				//delete hbase key
				if err = c.dao.DeletePathAuthH(ctx, projID, branch, path); err != nil {
					log.Error("delete Auth hbase key error(%+v) project ID:%d, branch:%s, path:%s", err, projID, branch, path)
					err = nil
				}
			} else {
				if raw, err = c.gitlab.RepoRawFile(projID, branch, file); err != nil {
					return
				}
				cb = readContributor(raw)
				log.Info("SaveContributor projectID:%d, branch:%s, path:%s, owner:%+v, reviewer:%+v", projID, branch, path, cb.Owner, cb.Reviewer)

				//add to redis
				authUser := &model.AuthUsers{
					Owners:    cb.Owner,
					Reviewers: cb.Reviewer,
				}
				if err = c.dao.SetPathAuthR(ctx, projID, branch, path, authUser); err != nil {
					log.Error("Set Auth to Redis error(%+v) project ID:%d, path:%s, owner:%+v, reviewer:%+v", err, projID, path, cb.Owner, cb.Reviewer)
					err = nil
				}

				//add to hbase
				if err = c.dao.SetPathAuthH(ctx, projID, branch, path, cb.Owner, cb.Reviewer); err != nil {
					log.Error("Set Auth to Hbase error(%+v) projectID:%d, branch:%s, path:%s, owner:%+v, reviewer:%+v", err, projID, branch, path, cb.Owner, cb.Reviewer)
					err = nil
				}
			}
		}
	}
	return
}

// checkSuperAuth ...
func checkSuperAuth(username string, reviewedUsers []string, superUsers []string) (ok bool) {
	for _, super := range superUsers {
		if strings.EqualFold(super, username) {
			return true
		}
	}

	for _, r := range reviewedUsers {
		for _, s := range superUsers {
			if strings.EqualFold(r, s) {
				return true
			}
		}
	}
	return false
}

// checkAllPathAuth ...
func (c *Command) checkAllPathAuth(taskInfo *model.TaskInfo) (ok bool, err error) {
	var (
		comment              string
		files                []string
		folders              map[string]string
		owners               []string
		reviewers            []string
		requireOwners        []string
		requireReviewers     []string
		reviewedUsers        []string
		requireReviewFolders []*model.RequireReviewFolder
		username             = taskInfo.Event.User.UserName
		projID               = int(taskInfo.Event.MergeRequest.SourceProjectID)
		mrIID                = int(taskInfo.Event.MergeRequest.IID)
		sourceBranch         = taskInfo.Event.MergeRequest.SourceBranch
		targetBranch         = taskInfo.Event.MergeRequest.TargetBranch
		authBranch           = c.GetAuthBranch(targetBranch, taskInfo.Repo.Config.AuthBranches)
		url                  = taskInfo.Event.ObjectAttributes.URL
		minReviewer          = taskInfo.Repo.Config.MinReviewer
		limitAuth            = taskInfo.Repo.Config.LimitAuth
		minReviewTip         bool
		isOwner              bool
	)

	log.Info("checkAllPathAuth start ... MRIID: %d", mrIID)
	if reviewedUsers, err = c.reviewedUsers(projID, mrIID); err != nil {
		return
	}
	if len(taskInfo.Repo.Config.SuperAuthUsers) > 0 {
		log.Info("checkAllPathAuth MRIID: %d, reviewedUsers: %v, SuperAuthUsers: %v", mrIID, reviewedUsers, taskInfo.Repo.Config.SuperAuthUsers)
		ok = checkSuperAuth(username, reviewedUsers, taskInfo.Repo.Config.SuperAuthUsers)
		if !ok {
			comment, err = c.showRequireSuperAuthComment(taskInfo)
			go notification.WechatAuthor(c.dao, username, url, sourceBranch, targetBranch, comment)
		}
		return
	}

	if files, err = c.gitlab.CompareDiff(projID, targetBranch, sourceBranch); err != nil {
		return
	}

	log.Info("checkAllPathAuth projID:%d, MRIID:%d, reviewedUsers:%+v, files:%+v, targetBranch:%s, sourceBranch:%s", projID, mrIID, reviewedUsers, files, targetBranch, sourceBranch)
	// 去重目录，校验目录权限
	folders = make(map[string]string)
	for _, file := range files {
		folder := filepath.Dir(file)
		if _, has := folders[folder]; !has {
			folders[folder] = folder
			authEnough := false
			ownerPathReviewed := false
			ownerReviewedStat := false
			reviewedCount := 0
			requireOwners = []string{}
			requireReviewers = []string{}
			dir := folder
			for {
				if owners, reviewers, err = c.GetPathAuth(projID, authBranch, dir); err != nil {
					return
				}

				isOwner, ownerPathReviewed = reviewedOwner(owners, reviewedUsers, username)
				num := reviewedNum(reviewers, reviewedUsers)
				reviewedCount = reviewedCount + num

				if isOwner {
					log.Info("checkAllPathAuth user is owner, MRIID: %d, user: %s, owners %v, path: %s", mrIID, username, owners, dir)
					authEnough = true
					break
				}
				if ownerPathReviewed {
					ownerReviewedStat = true
					if reviewedCount >= minReviewer {
						log.Info("checkAllPathAuth owner reviewed and reach minReviewer, user: %s, owners %+v, path: %s, reviewedUsers: %+v, reviewedCount: %d, minReviewer: %d", username, owners, dir, reviewedUsers, reviewedCount, minReviewer)
						authEnough = true
						break
					}
				}

				if (len(requireOwners) == 0) && (len(owners) > 0) {
					log.Info("checkAllPathAuth required owners %v, path: %s, MRIID: %d", owners, dir, mrIID)
					requireOwners = owners
				}
				if (len(requireReviewers) == 0) && (len(reviewers) > 0) {
					log.Info("checkAllPathAuth required reviewers %v, path: %s, MRIID: %d", reviewers, dir, mrIID)
					requireReviewers = reviewers
				}
				if dir == "." {
					break
				}
				if (len(owners) > 0) && limitAuth {
					break
				}
				dir = filepath.Dir(dir)
			}
			if !authEnough {
				requireAuth := &model.RequireReviewFolder{}
				requireAuth.Folder = folder
				if !ownerReviewedStat {
					requireAuth.Owners = requireOwners
				}
				if reviewedCount < minReviewer {
					requireAuth.Reviewers = requireReviewers
					minReviewTip = true
				}
				requireReviewFolders = append(requireReviewFolders, requireAuth)
			}
		}
	}
	if len(requireReviewFolders) > 0 {
		comment, err = c.showRequireAuthComment(taskInfo, requireReviewFolders, minReviewTip)
		go notification.WechatAuthor(c.dao, username, url, sourceBranch, targetBranch, comment)
		return
	}
	ok = true
	return
}

// GetAllPathAuth ...
func (c *Command) GetAllPathAuth(projID int, mrIID int, authBranch string) (pathOwners []model.RequireReviewFolder, err error) {
	var (
		files        []string
		deleteFiles  []string
		folders      map[string]string
		pathOwner    []string
		pathReviewer []string
	)
	if files, deleteFiles, err = c.gitlab.MRChanges(projID, mrIID); err != nil {
		return
	}
	files = append(files, deleteFiles...)
	// 去重目录，校验目录权限
	folders = make(map[string]string)
	for _, file := range files {
		folder := filepath.Dir(file)
		if _, has := folders[folder]; !has {
			folders[folder] = folder
			dir := folder

			for {
				if pathOwner, pathReviewer, err = c.GetPathAuth(projID, authBranch, dir); err != nil {
					return
				}

				if len(pathOwner) > 0 {
					exist := false
					for _, os := range pathOwners {
						if os.Folder == dir {
							exist = true
							break
						}
					}
					if exist {
						break
					}

					requireOwner := model.RequireReviewFolder{}
					requireOwner.Folder = dir
					requireOwner.Owners = pathOwner
					if len(pathReviewer) > 0 {
						requireOwner.Reviewers = pathReviewer
					}

					pathOwners = append(pathOwners, requireOwner)
					break
				}
				if dir == "." {
					break
				}
				dir = filepath.Dir(dir)
			}
		}
	}
	return
}

// GetPathAuth ...
func (c *Command) GetPathAuth(projID int, branch string, path string) (owners []string, reviewers []string, err error) {
	var (
		ctx      = context.TODO()
		authUser *model.AuthUsers
	)

	if authUser, err = c.dao.PathAuthR(ctx, projID, branch, path); err != nil || authUser == nil {
		if err != nil {
			log.Error("GetPathAuthInfo error project ID:%d, branch:%s, path: %s (err: %+v) ", projID, branch, path, err)
		}

		if owners, reviewers, err = c.dao.PathAuthH(ctx, projID, branch, path); err != nil {
			return
		}
		if len(owners) <= 0 && len(reviewers) <= 0 {
			return
		}

		if authUser == nil {
			authUser = new(model.AuthUsers)
		}
		authUser.Owners = owners
		authUser.Reviewers = reviewers
		if err1 := c.dao.SetPathAuthR(ctx, projID, branch, path, authUser); err1 != nil {
			log.Error("SetPathAuthR error project ID:%d, branch:%s, path: %s (err1: %+v) ", projID, branch, path, err1)
		}
		return
	}
	if authUser != nil {
		owners = authUser.Owners
		reviewers = authUser.Reviewers
	}
	return
}

// GetAuthBranch ...
func (c *Command) GetAuthBranch(targetBranch string, authBranches []string) (authBranch string) {
	for _, r := range authBranches {
		if r == targetBranch {
			return r
		}
	}
	return authBranches[0]
}

// showRequireAuthComment ...
func (c *Command) showRequireAuthComment(taskInfo *model.TaskInfo, requireReviewFolders []*model.RequireReviewFolder, minReviewTip bool) (comment string, err error) {
	var (
		username    = taskInfo.Event.User.UserName
		mrIID       = int(taskInfo.Event.MergeRequest.IID)
		projID      = int(taskInfo.Event.MergeRequest.SourceProjectID)
		noteID      = taskInfo.NoteID
		minReviewer = taskInfo.Repo.Config.MinReviewer
	)

	if minReviewTip {
		comment = fmt.Sprintf("<pre>[%s]尝试合并失败，无权限的文件目录如下，并且最少需要review人数 [%d] 个，请寻找对应的大佬 +1 后再 +merge：</pre>", username, minReviewer)
	} else {
		comment = fmt.Sprintf("<pre>[%s]尝试合并失败，无权限的文件目录如下，请寻找对应的大佬 +1 后再 +merge：</pre>", username)
	}
	comment += "\n"

	for _, reviewFolder := range requireReviewFolders {
		comment += fmt.Sprintf("+ %s ", reviewFolder.Folder)
		if len(reviewFolder.Owners) > 0 {
			comment += "；OWNER: "
			for _, o := range reviewFolder.Owners {
				if o == "all" {
					comment += "所有人" + " 或 "
				} else {
					comment += o + " 或 "
				}
			}
			comment = comment[:len(comment)-len(" 或 ")]
		}
		if len(reviewFolder.Reviewers) > 0 {
			comment += "；REVIEWER: "
			for _, o := range reviewFolder.Reviewers {
				if o == "all" {
					comment += "所有人" + " 与 "
				} else {
					comment += o + " 与 "
				}
			}
			comment = comment[:len(comment)-len(" 与 ")]
			//comment += fmt.Sprintf("<pre>; 已review人数 [%d] 个</pre>", minReviewer)
		}
		comment += "\n\n"
	}

	err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
	return
}

// showRequireSuperAuthComment ...
func (c *Command) showRequireSuperAuthComment(taskInfo *model.TaskInfo) (comment string, err error) {
	var (
		username = taskInfo.Event.User.UserName
		mrIID    = int(taskInfo.Event.MergeRequest.IID)
		projID   = int(taskInfo.Event.MergeRequest.SourceProjectID)
		noteID   = taskInfo.NoteID
	)

	comment = fmt.Sprintf("<pre>[%s]尝试合并失败，已配置超级权限用户，请寻找其中至少一位大佬 +1 后再 +merge：</pre>", username)
	comment += "\n"

	comment += fmt.Sprintf("+ SUPERMAN: ")
	for _, user := range taskInfo.Repo.Config.SuperAuthUsers {
		comment += fmt.Sprintf("%s 或 ", user)
	}
	comment = comment[:len(comment)-len(" 或 ")]

	err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
	return
}
