package dao

import (
	"context"
	"time"

	"go-common/app/admin/ep/saga/model"

	pkgerr "github.com/pkg/errors"
)

/*-------------------------------------- commit ----------------------------------------*/

// HasCommit ...
func (d *Dao) HasCommit(projID int, commitID string) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsCommits{}).Where("project_id = ? AND commit_id = ?", projID, commitID).Count(&total).Error)
	return
}

// DelCommit ...
func (d *Dao) DelCommit(projID int, commitID string) (err error) {
	return pkgerr.WithStack(d.db.Where("project_id = ? AND commit_id = ?", projID, commitID).Delete(&model.StatisticsCommits{}).Error)
}

// UpdateCommit ...
func (d *Dao) UpdateCommit(projID int, commitID string, commit *model.StatisticsCommits) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsCommits{}).Where("project_id = ? AND commit_id = ?", projID, commitID).Update(commit).Error)
}

// CreateCommit ...
func (d *Dao) CreateCommit(req *model.StatisticsCommits) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

// QueryCommitByID query commit info by ID
func (d *Dao) QueryCommitByID(projID int, commitID string) (commit *model.StatisticsCommits, err error) {
	commit = &model.StatisticsCommits{}
	err = pkgerr.WithStack(d.db.Where("project_id = ? AND commit_id = ?", projID, commitID).First(commit).Error)
	return
}

// CountCommitByTime ...
func (d *Dao) CountCommitByTime(projID int, since, until string) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsCommits{}).Where("project_id = ? AND created_at > ? AND created_at < ?", projID, since, until).Count(&total).Error)
	return
}

// QueryProjectCommits ...
func (d *Dao) QueryProjectCommits(c context.Context, projectID int) (commits []*model.StatisticsCommits, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsCommits{}).Where("project_id = ?", projectID).Find(&commits).Error)
	return
}

/*-------------------------------------- issue ----------------------------------------*/

// HasIssue ...
func (d *Dao) HasIssue(projID, issueID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsIssues{}).Where("project_id = ? AND issue_id = ?", projID, issueID).Count(&total).Error)
	return
}

// UpdateIssue ...
func (d *Dao) UpdateIssue(projID, issueID int, issue *model.StatisticsIssues) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsIssues{}).Where("project_id = ? AND issue_id = ?", projID, issueID).Update(issue).Error)
}

// CreateIssue ...
func (d *Dao) CreateIssue(req *model.StatisticsIssues) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

/*-------------------------------------- runner ----------------------------------------*/

// HasRunner ...
func (d *Dao) HasRunner(projID, runnerID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsRunners{}).Where("project_id = ? AND runner_id = ?", projID, runnerID).Count(&total).Error)
	return
}

// UpdateRunner ...
func (d *Dao) UpdateRunner(projID, runnerID int, runner *model.StatisticsRunners) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsRunners{}).Where("project_id = ? AND runner_id = ?", projID, runnerID).Update(runner).Error)
}

// CreateRunner ...
func (d *Dao) CreateRunner(req *model.StatisticsRunners) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

/*-------------------------------------- job ----------------------------------------*/

// HasJob ...
func (d *Dao) HasJob(projID, jobID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsJobs{}).Where("project_id = ? AND job_id=?", projID, jobID).Count(&total).Error)
	return
}

// UpdateJob ...
func (d *Dao) UpdateJob(projID, jobID int, runner *model.StatisticsJobs) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsJobs{}).Where("project_id = ? AND job_id = ?", projID, jobID).Update(runner).Error)

}

// CreateJob ...
func (d *Dao) CreateJob(req *model.StatisticsJobs) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

// QueryJobsByTime ...
func (d *Dao) QueryJobsByTime(projID int, req *model.ProjectJobRequest, since, until string) (total int, jobs []*model.StatisticsJobs, err error) {

	gDB := d.db.Model(&model.StatisticsJobs{}).Where("project_id = ? AND created_at > ? AND created_at < ?", projID, since, until)
	if req.Branch != "" {
		gDB = gDB.Where("ref = ?", req.Branch)
	}
	if req.User != "" {
		gDB = gDB.Where("user_name = ?", req.User)
	}
	if req.Machine != "" {
		gDB = gDB.Where("runner_description = ?", req.Machine)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}
	if err = pkgerr.WithStack(gDB.Find(&jobs).Error); err != nil {
		return
	}

	return
}

/*-------------------------------------- note ----------------------------------------*/

// HasNote ...
func (d *Dao) HasNote(c context.Context, projectID, noteID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsNotes{}).Where("project_id = ? AND note_id = ?", projectID, noteID).Count(&total).Error)
	return
}

// UpdateNote ...
func (d *Dao) UpdateNote(c context.Context, projectID, noteID int, note *model.StatisticsNotes) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsNotes{}).Where("project_id = ? AND note_id = ?", projectID, noteID).Update(note).Error)
}

// CreateNote ...
func (d *Dao) CreateNote(c context.Context, note *model.StatisticsNotes) (err error) {
	return pkgerr.WithStack(d.db.Create(note).Error)
}

// NoteByMRIID ...
func (d *Dao) NoteByMRIID(c context.Context, projectID, mrIID int) (notes []*model.StatisticsNotes, err error) {
	err = pkgerr.WithStack(d.db.Where("project_id = ? AND mr_iid = ? ", projectID, mrIID).Order("note_id desc").Find(&notes).Error)
	return
}

// NoteByID ...
func (d *Dao) NoteByID(c context.Context, projectID, mrIID, noteID int) (note *model.StatisticsNotes, err error) {
	note = &model.StatisticsNotes{}
	err = pkgerr.WithStack(d.db.Where("project_id = ? AND mr_iid = ? AND note_id = ?", projectID, mrIID, noteID).First(note).Error)
	return
}

/*-------------------------------------- mr ----------------------------------------*/

// CountMRByTime ...
func (d *Dao) CountMRByTime(projID int, since, until string) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsMrs{}).Where("project_id = ? AND created_at > ? AND created_at < ?", projID, since, until).Count(&total).Error)
	return
}

// MRByProjectID query mr by projectID
func (d *Dao) MRByProjectID(c context.Context, projectID int, since *time.Time, until *time.Time) (mrs []*model.StatisticsMrs, err error) {
	gDB := d.db.Model(&model.StatisticsMrs{}).Where("project_id = ? ", projectID)
	if since != nil || until != nil {
		gDB = gDB.Where("updated_at >= ? AND updated_at <=  ?", since, until)
	}
	err = pkgerr.WithStack(gDB.Find(&mrs).Error)
	return
}

// HasMR ...
func (d *Dao) HasMR(c context.Context, projectID, mrIID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsMrs{}).Where("project_id = ? AND mr_iid = ?", projectID, mrIID).Count(&total).Error)
	return
}

// UpdateMR update
func (d *Dao) UpdateMR(c context.Context, projectID, mrIID int, mr *model.StatisticsMrs) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsMrs{}).Where("project_id = ? AND mr_iid = ?", projectID, mrIID).Update(mr).Error)
}

// CreateMR ...
func (d *Dao) CreateMR(c context.Context, req *model.StatisticsMrs) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

/*-------------------------------------- aggMR ----------------------------------------*/

// HasAggregateReviewer ...
func (d *Dao) HasAggregateReviewer(c context.Context, projectID, mrIID, reviewerID, reviewID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.AggregateMrReviewer{}).Where("project_id = ? and mr_iid = ? and reviewer_id=? and review_id = ? ", projectID, mrIID, reviewerID, reviewID).Count(&total).Error)
	return
}

// UpdateAggregateReviewer ...
func (d *Dao) UpdateAggregateReviewer(c context.Context, projectID, mrIID, reviewerID, reviewID int, aggregateReviewer *model.AggregateMrReviewer) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.AggregateMrReviewer{}).Where("project_id = ? and mr_iid = ? and reviewer_id=? and review_id = ? ", projectID, mrIID, reviewerID, reviewID).Update(aggregateReviewer).Error)
}

// CreateAggregateReviewer ...
func (d *Dao) CreateAggregateReviewer(c context.Context, req *model.AggregateMrReviewer) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

/*-------------------------------------- pipeline ----------------------------------------*/

// HasPipeline ...
func (d *Dao) HasPipeline(projID int, pipelineID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsPipeline{}).Where("project_id = ? AND pipeline_id = ?", projID, pipelineID).Count(&total).Error)
	return
}

// UpdatePipeline ...
func (d *Dao) UpdatePipeline(projID int, pipelineID int, pipeline *model.StatisticsPipeline) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsPipeline{}).Where("project_id = ? AND pipeline_id = ?", projID, pipelineID).Update(pipeline).Error)
}

// CreatePipeline ...
func (d *Dao) CreatePipeline(req *model.StatisticsPipeline) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

// QueryPipelinesByTime ...
func (d *Dao) QueryPipelinesByTime(projID int, req *model.PipelineDataReq, since, until string) (total, statNum int, pipelines []*model.StatisticsPipeline, err error) {

	gDB := d.db.Model(&model.StatisticsPipeline{}).Where("project_id = ? AND created_at > ? AND created_at < ?", projID, since, until)
	if req.Branch != "" {
		gDB = gDB.Where("ref = ?", req.Branch)
	}
	if req.User != "" {
		gDB = gDB.Where("user = ?", req.User)
	}
	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	if req.State != "" {
		gDB = gDB.Where("status = ?", req.State)
	}
	if err = pkgerr.WithStack(gDB.Count(&statNum).Error); err != nil {
		return
	}

	if err = pkgerr.WithStack(gDB.Find(&pipelines).Error); err != nil {
		return
	}

	return
}

/*-------------------------------------- member ----------------------------------------*/

// HasMember ...
func (d *Dao) HasMember(c context.Context, projectID, memberID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsMembers{}).Where("project_id = ? AND member_id = ?", projectID, memberID).Count(&total).Error)
	return
}

// UpdateMember ...
func (d *Dao) UpdateMember(c context.Context, projectID, memberID int, member *model.StatisticsMembers) (err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsMembers{}).Where("project_id = ? AND member_id = ?", projectID, memberID).Update(member).Error)
	return
}

// QueryMemberByID ...
func (d *Dao) QueryMemberByID(c context.Context, projectID, memberID int) (member *model.StatisticsMembers, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsMembers{}).Where("project_id = ? AND member_id = ?", projectID, memberID).First(&member).Error)
	return
}

// CreateMember ...
func (d *Dao) CreateMember(c context.Context, member *model.StatisticsMembers) (err error) {
	return pkgerr.WithStack(d.db.Create(member).Error)
}

/*-------------------------------------- emoji ----------------------------------------*/

// HasMRAwardEmoji ...
func (d *Dao) HasMRAwardEmoji(c context.Context, projectID, mrIID, awardEmojiID int) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsMRAwardEmojis{}).Where("project_id = ? and mr_iid = ? and award_emoji_id=? ", projectID, mrIID, awardEmojiID).Count(&total).Error)
	return
}

// UpdateMRAwardEmoji update
func (d *Dao) UpdateMRAwardEmoji(c context.Context, projectID, mrIID, awardEmojiID int, awardEmoji *model.StatisticsMRAwardEmojis) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsMRAwardEmojis{}).Where("project_id = ? and mr_iid = ? and award_emoji_id=? ", projectID, mrIID, awardEmojiID).Update(awardEmoji).Error)
}

// CreateMRAwardEmoji ...
func (d *Dao) CreateMRAwardEmoji(c context.Context, awardEmoji *model.StatisticsMRAwardEmojis) (err error) {
	return pkgerr.WithStack(d.db.Create(awardEmoji).Error)
}

// AwardEmojiByMRIID ...
func (d *Dao) AwardEmojiByMRIID(c context.Context, projectID, mrIID int) (AwardEmojis []*model.StatisticsMRAwardEmojis, err error) {
	err = pkgerr.WithStack(d.db.Where("project_id = ? AND mr_iid = ? ", projectID, mrIID).Find(&AwardEmojis).Error)
	return
}

/*-------------------------------------- discussion ----------------------------------------*/

// HasDiscussion ...
func (d *Dao) HasDiscussion(c context.Context, projectID, mrIID int, discussionID string) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsDiscussions{}).Where("project_id = ? and mr_iid = ? and discussion_id=? ", projectID, mrIID, discussionID).Count(&total).Error)
	return
}

// UpdateDiscussion ...
func (d *Dao) UpdateDiscussion(c context.Context, projectID, mrIID int, discussionID string, discussion *model.StatisticsDiscussions) (err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsDiscussions{}).Where("project_id = ? and mr_iid = ? and discussion_id=? ", projectID, mrIID, discussionID).Update(discussion).Error)
	return
}

// CreateDiscussion ...
func (d *Dao) CreateDiscussion(c context.Context, discussion *model.StatisticsDiscussions) (err error) {
	return pkgerr.WithStack(d.db.Create(discussion).Error)
}

// DiscussionsByMRIID ...
func (d *Dao) DiscussionsByMRIID(c context.Context, projectID, mrIID int) (discussions []*model.StatisticsDiscussions, err error) {
	err = pkgerr.WithStack(d.db.Where("project_id = ? AND mr_iid = ? ", projectID, mrIID).Find(&discussions).Error)
	return
}

/*-------------------------------------- Branch ----------------------------------------*/

// DeleteProjectBranch ...
func (d *Dao) DeleteProjectBranch(c context.Context, projectID int) error {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsBranches{}).Where("project_id = ?", projectID).Update("is_deleted", model.BranchDeleted).Error)
}

// HasBranch is exist project branch info in database.
func (d *Dao) HasBranch(c context.Context, projectID int, branchName string) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsBranches{}).Where("project_id = ? AND branch_name = ? ", projectID, branchName).Count(&total).Error)
	return
}

// UpdateBranch update
func (d *Dao) UpdateBranch(c context.Context, projectID int, branchName string, branch *model.StatisticsBranches) error {
	return pkgerr.WithStack(d.db.Model(&model.StatisticsBranches{}).Where("project_id = ? AND branch_name = ? ", projectID, branchName).Save(branch).Error)
}

// CreateBranch query all the records in contact_infos
func (d *Dao) CreateBranch(c context.Context, branch *model.StatisticsBranches) error {
	return pkgerr.WithStack(d.db.Create(branch).Error)
}

// QueryProjectBranch ...
func (d *Dao) QueryProjectBranch(c context.Context, projectID int) (b []*model.StatisticsBranches, err error) {
	err = pkgerr.WithStack(d.db.Where("project_id = ? and is_deleted != ?", projectID, model.BranchDeleted).Find(&b).Error)
	return
}

// QueryFirstBranch ...
func (d *Dao) QueryFirstBranch(c context.Context, projectID int, branchName string) (branch *model.StatisticsBranches, err error) {
	branch = &model.StatisticsBranches{}
	err = pkgerr.WithStack(d.db.Model(&model.StatisticsBranches{}).Where("project_id = ? AND branch_name = ? ", projectID, branchName).First(branch).Error)
	return
}

// DeleteAggregateBranch ...
func (d *Dao) DeleteAggregateBranch(c context.Context, projectID int) error {
	return pkgerr.WithStack(d.db.Model(&model.AggregateBranches{}).Where("project_id = ?", projectID).Update("is_deleted", model.BranchDeleted).Error)
}

// HasAggregateBranch ...
func (d *Dao) HasAggregateBranch(c context.Context, projectID int, branchName string) (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.AggregateBranches{}).Where("project_id = ? AND branch_name = ? ", projectID, branchName).Count(&total).Error)
	return
}

// UpdateAggregateBranch ...
func (d *Dao) UpdateAggregateBranch(c context.Context, projectID int, branchName string, aggregateBranch *model.AggregateBranches) error {
	return pkgerr.WithStack(d.db.Model(&model.AggregateBranches{}).Where("project_id = ? AND branch_name = ?", projectID, branchName).Save(aggregateBranch).Error)
}

// CreateAggregateBranch ...
func (d *Dao) CreateAggregateBranch(c context.Context, aggregateBranch *model.AggregateBranches) error {
	return pkgerr.WithStack(d.db.Create(aggregateBranch).Error)
}

// QueryFirstAggregateBranch ...
func (d *Dao) QueryFirstAggregateBranch(c context.Context, projectID int, branchName string) (branch *model.AggregateBranches, err error) {
	branch = &model.AggregateBranches{}
	err = pkgerr.WithStack(d.db.Model(&model.AggregateBranches{}).Where("project_id = ? AND branch_name = ? ", projectID, branchName).First(branch).Error)
	return
}
