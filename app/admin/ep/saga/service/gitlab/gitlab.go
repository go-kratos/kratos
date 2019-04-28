package gitlab

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

const (
	_defaultPerPage = 20
	_maxPerPage     = 100
)

// Gitlab def
type Gitlab struct {
	url    string
	token  string
	client *gitlab.Client
}

// New new gitlab structure
func New(url string, token string) (g *Gitlab) {
	g = &Gitlab{
		url:    url,
		token:  token,
		client: gitlab.NewClient(nil, token),
	}
	g.client.SetBaseURL(url)
	return
}

// ListProjects list all project
func (g *Gitlab) ListProjects(page int) (projects []*gitlab.Project, err error) {

	opt := &gitlab.ListProjectsOptions{}
	opt.ListOptions.Page = page
	opt.ListOptions.PerPage = _defaultPerPage

	if projects, _, err = g.client.Projects.ListProjects(opt); err != nil {
		err = errors.Wrapf(err, "ListProjects err(%+v)", err)
		return
	}
	return
}

// ListProjectCommit ...
func (g *Gitlab) ListProjectCommit(projectID, page int, since, until *time.Time) (commits []*gitlab.Commit, resp *gitlab.Response, err error) {
	var (
		opt = &gitlab.ListCommitsOptions{}
		all = true
	)
	opt.Page = page
	opt.PerPage = _maxPerPage
	opt.All = &all

	if since != nil {
		opt.Since = since
	}
	if until != nil {
		opt.Until = until
	}

	if commits, resp, err = g.client.Commits.ListCommits(projectID, opt); err != nil {
		err = errors.Wrapf(err, "ListCommits projectId(%d), err(%+v)", projectID, err)
		return
	}
	return
}

// ListProjectMergeRequests list all MR
func (g *Gitlab) ListProjectMergeRequests(projectID int, since, until *time.Time, page int) (mrs []*gitlab.MergeRequest, resp *gitlab.Response, err error) {

	opt := &gitlab.ListProjectMergeRequestsOptions{}
	opt.UpdatedAfter = since
	opt.UpdatedBefore = until
	if page != -1 {
		opt.PerPage = _defaultPerPage
		opt.Page = page
	}
	if mrs, resp, err = g.client.MergeRequests.ListProjectMergeRequests(projectID, opt); err != nil {
		err = errors.Wrapf(err, "ListProjectMergeRequests projectId(%d), err(%+v)", projectID, err)
		return
	}
	return
}

// ListProjectPipelines list all pipelines
func (g *Gitlab) ListProjectPipelines(page, projectID int, status gitlab.BuildStateValue) (pipelineList gitlab.PipelineList, resp *gitlab.Response, err error) {

	opt := &gitlab.ListProjectPipelinesOptions{}
	opt.ListOptions.Page = page
	if status != "" {
		opt.Status = gitlab.BuildState(status)
	}
	if pipelineList, resp, err = g.client.Pipelines.ListProjectPipelines(projectID, opt); err != nil {
		err = errors.Wrapf(err, "ListProjectPipelines projectId(%d), err(%+v)", projectID, err)
		return
	}
	return
}

// ListPipelineJobs list all pipeline jobs
func (g *Gitlab) ListPipelineJobs(opt *gitlab.ListJobsOptions, projectID, pipelineID int) (jobList []*gitlab.Job, resp *gitlab.Response, err error) {
	if jobList, resp, err = g.client.Jobs.ListPipelineJobs(projectID, pipelineID, opt); err != nil {
		err = errors.Wrapf(err, "ListPipelineJobs projectId(%d), pipelineId(%d), err(%+v)", projectID, pipelineID, err)
		return
	}
	return
}

// GetPipeline get a pipeline info
func (g *Gitlab) GetPipeline(projectID, pipelineID int) (pipeline *gitlab.Pipeline, resp *gitlab.Response, err error) {

	if pipeline, resp, err = g.client.Pipelines.GetPipeline(projectID, pipelineID); err != nil {
		err = errors.Wrapf(err, "GetPipeline projectId(%d), err(%+v)", projectID, err)
		return
	}
	return
}

// ListProjectBranch ...
func (g *Gitlab) ListProjectBranch(projectID, page int) (branches []*gitlab.Branch, resp *gitlab.Response, err error) {
	var opt = &gitlab.ListBranchesOptions{Page: page, PerPage: 20}
	if branches, resp, err = g.client.Branches.ListBranches(projectID, opt); err != nil {
		err = errors.Wrapf(err, "ListBranches projectId(%d), err(%+v)", projectID, err)
		return
	}
	return
}

// ListMRNotes get notes of the MR
func (g *Gitlab) ListMRNotes(c context.Context, projectID, mrID, page int) (notes []*gitlab.Note, resp *gitlab.Response, err error) {
	var opt = &gitlab.ListMergeRequestNotesOptions{Page: page, PerPage: 20}

	if notes, resp, err = g.client.Notes.ListMergeRequestNotes(projectID, mrID, opt); err != nil {
		err = errors.Wrapf(err, "ListMergeRequestNotes projectId(%d), mrId(%d), err(%+v)", projectID, mrID, err)
		return
	}
	return
}

// MergeBase 获取两个分支最近的合并commit
func (g *Gitlab) MergeBase(c context.Context, projectID int, refs []string) (commit *gitlab.Commit, resp *gitlab.Response, err error) {
	var opt = &gitlab.MergeBaseOptions{Ref: refs}
	if commit, resp, err = g.client.Repositories.MergeBase(projectID, opt); err != nil {
		err = errors.Wrapf(err, "MergeBase projectId(%d), refs(%v), err(%+v)", projectID, refs, err)
		return
	}
	return
}

// ListMRAwardEmoji ...
func (g *Gitlab) ListMRAwardEmoji(projectID, mrIID, page int) (emojis []*gitlab.AwardEmoji, resp *gitlab.Response, err error) {
	var opt = &gitlab.ListAwardEmojiOptions{Page: page, PerPage: _defaultPerPage}
	if emojis, resp, err = g.client.AwardEmoji.ListMergeRequestAwardEmoji(projectID, mrIID, opt); err != nil {
		err = errors.Wrapf(err, "ListMergeRequestAwardEmoji projectId(%d), mrIID(%d), err(%+v)", projectID, mrIID, err)
		return
	}
	return
}

//ListMRDiscussions ...
func (g *Gitlab) ListMRDiscussions(projectID, mrIID, page int) (result []*gitlab.Discussion, resp *gitlab.Response, err error) {
	var opt = &gitlab.ListMergeRequestDiscussionsOptions{Page: page, PerPage: _defaultPerPage}
	if result, resp, err = g.client.Discussions.ListMergeRequestDiscussions(projectID, mrIID, opt); err != nil {
		err = errors.Wrapf(err, "ListMergeRequestDiscussions projectId(%d), mrIID(%d), err(%+v)", projectID, mrIID, err)
		return
	}
	return
}

// GetMergeRequestDiff ...
func (g *Gitlab) GetMergeRequestDiff(projectID, mrID int) (mr *gitlab.MergeRequest, resp *gitlab.Response, err error) {
	if mr, resp, err = g.client.MergeRequests.GetMergeRequestChanges(projectID, mrID); err != nil {
		err = errors.Wrapf(err, "GetMergeRequestChanges projectId(%d), err(%+v)", projectID, err)
		return
	}
	return
}

// ListProjectJobs get jobs of the project
func (g *Gitlab) ListProjectJobs(projID, page int) (jobs []gitlab.Job, resp *gitlab.Response, err error) {
	var opt = &gitlab.ListJobsOptions{}
	opt.Page = page
	// perPage max is 100
	opt.PerPage = _maxPerPage

	if jobs, resp, err = g.client.Jobs.ListProjectJobs(projID, opt); err != nil {
		err = errors.Wrapf(err, "ListProjectJobs projectId(%d), err(%+v)", projID, err)
		return
	}
	return
}

// ListProjectMembers ...
func (g *Gitlab) ListProjectMembers(projectID, page int) (members []*gitlab.ProjectMember, resp *gitlab.Response, err error) {
	opt := &gitlab.ListProjectMembersOptions{}
	opt.Page = page
	opt.PerPage = _defaultPerPage
	if members, resp, err = g.client.ProjectMembers.ListProjectMembers(projectID, opt); err != nil {
		err = errors.Wrapf(err, "ListProjectMembers projectId(%d), err(%+v)", projectID, err)
		return
	}
	return
}

// ListProjectIssues ...
func (g *Gitlab) ListProjectIssues(projID, page int, since, until *time.Time) (issues []*gitlab.Issue, resp *gitlab.Response, err error) {
	var opt = &gitlab.ListProjectIssuesOptions{}
	opt.Page = page
	// perPage max is 100
	opt.PerPage = _maxPerPage
	if since != nil {
		opt.CreatedBefore = until
	}
	if until != nil {
		opt.CreatedAfter = since
	}

	if issues, resp, err = g.client.Issues.ListProjectIssues(projID, opt); err != nil {
		err = errors.Wrapf(err, "ListProjectIssues projectId(%d), err(%+v)", projID, err)
		return
	}
	return
}

// ListProjectRunners ...
func (g *Gitlab) ListProjectRunners(projID, page int) (result []*gitlab.Runner, resp *gitlab.Response, err error) {
	var opt = &gitlab.ListProjectRunnersOptions{}
	opt.ListOptions.Page = page
	opt.PerPage = _defaultPerPage
	if result, resp, err = g.client.Runners.ListProjectRunners(projID, opt); err != nil {
		err = errors.Wrapf(err, "ListProjectRunners projectId(%d), err(%+v)", projID, err)
		return
	}
	return
}
