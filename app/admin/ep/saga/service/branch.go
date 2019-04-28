package service

import (
	"context"
	"encoding/json"
	"sort"

	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/utils"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

// QueryProjectBranchList query project commit info according to project id.
func (s *Service) QueryProjectBranchList(c context.Context, req *model.ProjectDataReq) (resp []*gitlab.Branch, err error) {
	var (
		branch     []*gitlab.Branch
		response   *gitlab.Response
		branchItem []*gitlab.Branch
	)

	if branch, response, err = s.gitlab.ListProjectBranch(req.ProjectID, 1); err != nil {
		return
	}

	page := 2
	for page <= response.TotalPages {
		if branchItem, _, err = s.gitlab.ListProjectBranch(req.ProjectID, page); err != nil {
			return
		}
		branch = append(branch, branchItem...)
		page++
	}
	resp = branch
	return
}

// QueryBranchDiffWith ...
func (s *Service) QueryBranchDiffWith(c context.Context, req *model.BranchDiffWithRequest) (resp []*model.BranchDiffWithResponse, err error) {
	var (
		branchDiff []*model.BranchDiffWithResponse
	)

	// 默认主分支为master
	if req.Master == "" {
		req.Master = "master"
	}

	if branchDiff, err = s.AllBranchDiffWithMaster(c, req.ProjectID, req.Master); err != nil {
		return
	}
	if req.Branch != "" {
		for _, branch := range branchDiff {
			if branch.Branch == req.Branch {
				resp = append(resp, branch)
				return
			}
		}
	}
	switch req.SortBy {
	case "update":
		sort.Slice(branchDiff, func(i, j int) bool {
			if branchDiff[i].LatestUpdateTime.After(*branchDiff[j].LatestUpdateTime) {
				return false
			}
			if branchDiff[i].LatestUpdateTime.Before(*branchDiff[j].LatestUpdateTime) {
				return true
			}
			return true
		})
	case "sync":
		sort.Slice(branchDiff, func(i, j int) bool {
			if branchDiff[i].LatestSyncTime.After(*branchDiff[j].LatestSyncTime) {
				return false
			}
			if branchDiff[i].LatestSyncTime.Before(*branchDiff[j].LatestSyncTime) {
				return true
			}
			return true
		})
	case "ahead":
		sort.Slice(branchDiff, func(i, j int) bool {
			if branchDiff[i].Ahead > branchDiff[j].Ahead {
				return true
			}
			if branchDiff[i].Ahead < branchDiff[j].Ahead {
				return false
			}
			return true
		})
	case "behind":
		sort.Slice(branchDiff, func(i, j int) bool {
			if branchDiff[i].Behind > branchDiff[j].Behind {
				return true
			}
			if branchDiff[i].Behind < branchDiff[j].Behind {
				return false
			}
			return true
		})
	}

	resp = branchDiff[:100]

	return
}

// AllBranchDiffWithMaster ...
func (s *Service) AllBranchDiffWithMaster(c context.Context, project int, master string) (branchDiff []*model.BranchDiffWithResponse, err error) {
	var (
		branches map[string]*model.CommitTreeNode
		tree     []*model.CommitTreeNode
		cacheKey string
	)

	log.Info("sync Branch diff start => %s", cacheKey)

	if branches, err = s.AllProjectBranchInfo(c, project); err != nil {
		return
	}
	if tree, err = s.BuildCommitTree(c, project); err != nil {
		return
	}

	for k, v := range branches {
		var (
			base   *gitlab.Commit
			ahead  int
			behind int
		)
		if base, _, err = s.gitlab.MergeBase(c, project, []string{k, master}); err != nil {
			return
		}

		// 计算领先commit
		ahead = s.ComputeCommitNum(v, base, &tree)

		// 计算落后commit
		behind = s.ComputeCommitNum(branches[master], base, &tree)

		log.Info("%s: ahead: %d behind:%d\n", k, ahead, behind)
		branchDiff = append(branchDiff, &model.BranchDiffWithResponse{Branch: k, Behind: behind, Ahead: ahead, LatestUpdateTime: v.CreatedAt, LatestSyncTime: base.CreatedAt})
	}

	log.Info("Redis: Save Branch Diff Info Successfully!")

	return
}

// ComputeCommitNum ...
func (s *Service) ComputeCommitNum(head *model.CommitTreeNode, base *gitlab.Commit, tree *[]*model.CommitTreeNode) (num int) {
	var (
		visitedQueue []string
		commitTree   = *tree
		i            int
	)

	if head.CommitID == base.ID {
		return
	}
	visitedQueue = append(visitedQueue, head.CommitID)
	for i = 0; ; i++ {
		pointer := visitedQueue[i]
		if pointer == base.ID {
			break
		}
		for _, node := range commitTree {
			if node.CommitID == pointer {
				if utils.InSlice(base.ID, node.Parents) {
					visitedQueue = append(visitedQueue, base.ID)
					break
				}
				for _, p := range node.Parents {
					if !utils.InSlice(p, visitedQueue) {
						visitedQueue = append(visitedQueue, p)
					}
				}
				break
			}
		}
		if i > 999 {
			break
		}
		if i == len(visitedQueue)-1 {
			break
		}
	}

	num = i

	return
}

// AllMergeBase ...
func (s *Service) AllMergeBase(c context.Context, project int, branches []string, master string) (mergeBase map[string]*gitlab.Commit, err error) {
	var base *gitlab.Commit

	mergeBase = make(map[string]*gitlab.Commit)
	for _, branch := range branches {
		if base, _, err = s.gitlab.MergeBase(c, project, []string{master, branch}); err != nil {
			return
		}
		mergeBase[branch] = base
	}

	return
}

//AllProjectBranchInfo 获取所有分支信息
func (s *Service) AllProjectBranchInfo(c context.Context, projectID int) (branches map[string]*model.CommitTreeNode, err error) {
	var projectBranches []*model.StatisticsBranches

	branches = make(map[string]*model.CommitTreeNode)
	if projectBranches, err = s.dao.QueryProjectBranch(c, projectID); err != nil {
		return
	}

	for _, b := range projectBranches {
		var (
			commitInfo    *model.StatisticsCommits
			commitParents []string
		)
		if commitInfo, err = s.dao.QueryCommitByID(projectID, b.CommitID); err != nil {
			log.Error("AllProjectBranchInfo QueryCommitByID err(%+v)", err)
			err = nil
			continue
		}
		if commitInfo.ParentIDs != "" {
			if err = json.Unmarshal([]byte(commitInfo.ParentIDs), &commitParents); err != nil {
				return
			}
		}
		branches[b.BranchName] = &model.CommitTreeNode{CommitID: b.CommitID, Parents: commitParents, CreatedAt: commitInfo.CreatedAt, Author: commitInfo.AuthorName}
	}
	return
}

// BuildCommitTree 获取到所有的Commits
func (s *Service) BuildCommitTree(c context.Context, project int) (tree []*model.CommitTreeNode, err error) {
	var projectCommits []*model.StatisticsCommits
	tree = []*model.CommitTreeNode{}
	if projectCommits, err = s.dao.QueryProjectCommits(c, project); err != nil {
		return
	}

	for _, c := range projectCommits {
		var commitParents []string
		commitParentsString := c.ParentIDs
		if commitParentsString != "" {
			if err = json.Unmarshal([]byte(commitParentsString), &commitParents); err != nil {
				log.Error("解析commit parents报错(%+v)", err)
			}
		}
		tree = append(tree, &model.CommitTreeNode{CommitID: c.CommitID, Parents: commitParents, CreatedAt: c.CreatedAt, Author: c.AuthorName})
	}
	return
}

/*-------------------------------------- sync branch ----------------------------------------*/

// SyncProjectBranch ...
func (s *Service) SyncProjectBranch(c context.Context, projectID int) (result *model.SyncResult, err error) {
	var (
		branches    []*gitlab.Branch
		resp        *gitlab.Response
		projectInfo *model.ProjectInfo
	)
	result = &model.SyncResult{}
	if err = s.dao.DeleteProjectBranch(c, projectID); err != nil {
		return
	}
	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}
	for page := 1; ; page++ {
		result.TotalPage++
		if branches, resp, err = s.gitlab.ListProjectBranch(projectID, page); err != nil {
			return
		}
		for _, branch := range branches {
			branchDB := &model.StatisticsBranches{
				ProjectID:          projectID,
				ProjectName:        projectInfo.Name,
				CommitID:           branch.Commit.ID,
				BranchName:         branch.Name,
				Protected:          branch.Protected,
				Merged:             branch.Merged,
				DevelopersCanPush:  branch.DevelopersCanPush,
				DevelopersCanMerge: branch.DevelopersCanMerge,
			}
			if err = s.SaveDatabaseBranch(c, branchDB); err != nil {
				log.Error("Branch 存数据库报错(%+V)", err)
				err = nil

				errData := &model.FailData{
					ChildIDStr: branch.Name,
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

// SaveDatabaseBranch ...
func (s *Service) SaveDatabaseBranch(c context.Context, branchDB *model.StatisticsBranches) (err error) {
	var (
		total  int
		branch *model.StatisticsBranches
	)

	if total, err = s.dao.HasBranch(c, branchDB.ProjectID, branchDB.BranchName); err != nil {
		log.Error("SaveDatabaseBranch HasBranch(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		if branch, err = s.dao.QueryFirstBranch(c, branchDB.ProjectID, branchDB.BranchName); err != nil {
			log.Error("SaveDatabaseBranch QueryFirstBranch(%+v)", err)
			return
		}
		branchDB.ID = branch.ID
		if err = s.dao.UpdateBranch(c, branchDB.ProjectID, branchDB.BranchName, branchDB); err != nil {
			log.Error("SaveDatabaseBranch UpdateBranch(%+v)", err)
			return
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseBranch Branch has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateBranch(c, branchDB); err != nil {
		log.Error("SaveDatabaseBranch CreateBranch(%+v)", err)
		return
	}

	return
}

/*-------------------------------------- agg branch ----------------------------------------*/

// AggregateProjectBranch ...
func (s *Service) AggregateProjectBranch(c context.Context, projectID int, master string) (err error) {
	var (
		branches    map[string]*model.CommitTreeNode
		tree        []*model.CommitTreeNode
		projectInfo *model.ProjectInfo
	)
	if err = s.dao.DeleteAggregateBranch(c, projectID); err != nil {
		return
	}
	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}
	if branches, err = s.AllProjectBranchInfo(c, projectID); err != nil {
		return
	}
	if tree, err = s.BuildCommitTree(c, projectID); err != nil {
		return
	}

	for k, v := range branches {
		var (
			base   *gitlab.Commit
			ahead  int
			behind int
		)
		if base, _, err = s.gitlab.MergeBase(c, projectID, []string{k, master}); err != nil {
			return
		}

		// 计算领先commit
		ahead = s.ComputeCommitNum(v, base, &tree)

		// 计算落后commit
		behind = s.ComputeCommitNum(branches[master], base, &tree)
		branch := &model.AggregateBranches{
			ProjectID:        projectID,
			ProjectName:      projectInfo.Name,
			BranchName:       k,
			BranchUserName:   v.Author,
			BranchMaster:     master,
			Behind:           behind,
			Ahead:            ahead,
			LatestUpdateTime: v.CreatedAt,
			LatestSyncTime:   base.CreatedAt,
		}

		if err = s.SaveAggregateBranchDatabase(c, branch); err != nil {
			log.Error("AggregateBranch 存数据库报错(%+V)", err)
			continue
		}
	}
	return
}

// SaveAggregateBranchDatabase ...
func (s *Service) SaveAggregateBranchDatabase(c context.Context, branchDB *model.AggregateBranches) (err error) {
	var (
		total           int
		aggregateBranch *model.AggregateBranches
	)

	if total, err = s.dao.HasAggregateBranch(c, branchDB.ProjectID, branchDB.BranchName); err != nil {
		log.Error("SaveAggregateBranchDatabase HasAggregateBranch(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		if aggregateBranch, err = s.dao.QueryFirstAggregateBranch(c, branchDB.ProjectID, branchDB.BranchName); err != nil {
			log.Error("SaveDatabaseBranch QueryFirstAggregateBranch(%+v)", err)
			return
		}
		branchDB.ID = aggregateBranch.ID
		if err = s.dao.UpdateAggregateBranch(c, branchDB.ProjectID, branchDB.BranchName, branchDB); err != nil {
			log.Error("SaveAggregateBranchDatabase UpdateAggregateBranch(%+v)", err)
			return
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveAggregateBranchDatabase AggregateBranch has more rows(%d)", total)
		return

	}

	// insert row now
	if err = s.dao.CreateAggregateBranch(c, branchDB); err != nil {
		log.Error("SaveAggregateBranchDatabase CreateAggregateBranch(%+v)", err)
		return
	}

	return
}
