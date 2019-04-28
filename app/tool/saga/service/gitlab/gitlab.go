package gitlab

import (
	"net/url"
	"strings"

	"go-common/app/tool/saga/model"
	"go-common/library/log"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
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

// HostToken ...
func (g *Gitlab) HostToken() (host string, token string, err error) {
	var u *url.URL
	if u, err = url.Parse(g.url); err != nil {
		return
	}
	host = u.Host
	token = g.token
	return
}

// LastGreenCommit get last green pipeline commit
func (g *Gitlab) LastGreenCommit(ciProjectID string, ciCommitRefName string) (commit string, err error) {
	var (
		status    = gitlab.Success
		pipelines gitlab.PipelineList
	)
	opt := &gitlab.ListProjectPipelinesOptions{
		Status: &status,
		Ref:    gitlab.String(ciCommitRefName),
	}
	if pipelines, _, err = g.client.Pipelines.ListProjectPipelines(ciProjectID, opt); err != nil {
		return
	}
	if len(pipelines) == 0 {
		return
	}
	commit = pipelines[0].Sha
	return
}

// AcceptMR accept Merge request
func (g *Gitlab) AcceptMR(projID int, mrIID int, msg string) (state string, err error) {
	var (
		opt = &gitlab.AcceptMergeRequestOptions{
			ShouldRemoveSourceBranch: gitlab.Bool(true),
			MergeCommitMessage:       gitlab.String(msg),
		}
		mergeRequest *gitlab.MergeRequest
	)
	if mergeRequest, _, err = g.client.MergeRequests.AcceptMergeRequest(projID, mrIID, opt); err != nil {
		err = errors.Wrapf(err, "AcceptMergeRequest(%d,%d)", projID, mrIID)
		return
	}
	state = mergeRequest.State
	return
}

// CloseMR close merge request
func (g *Gitlab) CloseMR(projID int, mrIID int) (err error) {
	var (
		opt = &gitlab.UpdateMergeRequestOptions{
			StateEvent: gitlab.String("close"),
		}
	)
	if _, _, err = g.client.MergeRequests.UpdateMergeRequest(projID, mrIID, opt); err != nil {
		err = errors.Wrapf(err, "CloseMR(%d,%d)", projID, mrIID)
	}
	return
}

// CreateMRNote create note to merge request
func (g *Gitlab) CreateMRNote(projID int, mrIID int, content string) (noteID int, err error) {
	var (
		opt = &gitlab.CreateMergeRequestNoteOptions{
			Body: gitlab.String(content),
		}
		note *gitlab.Note
	)
	if note, _, err = g.client.Notes.CreateMergeRequestNote(projID, mrIID, opt); err != nil {
		err = errors.Wrapf(err, "CreateMergeRequestNote(%d,%d)", projID, mrIID)
		return
	}
	noteID = note.ID
	return
}

// UpdateMRNote update the specified merge request note
func (g *Gitlab) UpdateMRNote(projID int, mrIID int, noteID int, content string) (err error) {
	var (
		opt = &gitlab.UpdateMergeRequestNoteOptions{
			Body: gitlab.String(content),
		}
	)
	if _, _, err = g.client.Notes.UpdateMergeRequestNote(projID, mrIID, noteID, opt); err != nil {
		err = errors.Wrapf(err, "UpdateMergeRequestNote(%d,%d,%d)", projID, mrIID, noteID)
		return
	}
	return
}

// DeleteMRNote delete the specified note for merge request
func (g *Gitlab) DeleteMRNote(projID int, mrIID int, noteID int) (err error) {
	if _, err = g.client.Notes.DeleteMergeRequestNote(projID, mrIID, noteID); err != nil {
		err = errors.Wrapf(err, "DeleteMergeRequestNote(%d,%d,%d)", projID, mrIID, noteID)
	}
	return
}

// UserName get user name for user id
func (g *Gitlab) UserName(userID int) (userName string, err error) {
	var (
		user *gitlab.User
	)

	if user, _, err = g.client.Users.GetUser(userID); err != nil {
		err = errors.WithStack(err)
		return
	}
	userName = user.Username
	return
}

// Triggers get pipeline triggers
func (g *Gitlab) Triggers(projectID int) (triggers []*gitlab.PipelineTrigger, err error) {
	var (
		opt = &gitlab.ListPipelineTriggersOptions{}
	)

	if triggers, _, err = g.client.PipelineTriggers.ListPipelineTriggers(projectID, opt); err != nil {
		err = errors.Wrapf(err, "ListPipelineTriggers (projectID: %d)", projectID)
		return nil, err
	}

	return triggers, nil
}

// CreateTrigger create the trigger to trigger pipeline
func (g *Gitlab) CreateTrigger(projectID int) (trigger *gitlab.PipelineTrigger, err error) {
	var (
		opt = &gitlab.AddPipelineTriggerOptions{
			Description: gitlab.String("gitlab-ci-build-on-merge-request"),
		}
	)

	if trigger, _, err = g.client.PipelineTriggers.AddPipelineTrigger(projectID, opt); err != nil {
		err = errors.Wrapf(err, "CreateTrigger (projectID: %d)", projectID)
		return nil, err
	}

	return trigger, nil
}

// TakeOwnership take ownership of the trigger
func (g *Gitlab) TakeOwnership(projectID int, triggerID int) (trigger *gitlab.PipelineTrigger, err error) {
	if trigger, _, err = g.client.PipelineTriggers.TakeOwnershipOfPipelineTrigger(projectID, triggerID); err != nil {
		err = errors.Wrapf(err, "TakeOwnershipOfPipelineTrigger (projectID: %d, triggerID: %d)", projectID, triggerID)
		return nil, err
	}

	return trigger, nil
}

// TriggerPipeline trigger the pipeline
func (g *Gitlab) TriggerPipeline(projectID int, ref string, token string) (pipeline *gitlab.Pipeline, err error) {
	var (
		opt = &gitlab.RunPipelineTriggerOptions{
			Ref:   gitlab.String(ref),
			Token: gitlab.String(token),
		}
	)
	if pipeline, _, err = g.client.PipelineTriggers.RunPipelineTrigger(projectID, opt); err != nil {
		err = errors.Wrapf(err, "RunPipelineTrigger (projectID: %d, ref: %s, token: %s)", projectID, ref, token)
		return nil, err
	}

	return pipeline, nil
}

// AwardEmojiUsernames get the username who gave the emoji
func (g *Gitlab) AwardEmojiUsernames(projID int, mrIID int) (usernames []string, err error) {
	var (
		opt = &gitlab.ListAwardEmojiOptions{
			Page:    1,
			PerPage: 100,
		}
		awards []*gitlab.AwardEmoji
	)
	if awards, _, err = g.client.AwardEmoji.ListMergeRequestAwardEmoji(projID, mrIID, opt); err != nil {
		return
	}
	for _, a := range awards {
		if a.Name == "thumbsup" {
			usernames = append(usernames, a.User.Username)
		}
	}
	return
}

// AuditProjects audit all the given projects
func (g *Gitlab) AuditProjects(repos []*model.Repo, hooks []*model.WebHook) (err error) {
	var (
		repo   *model.Repo
		projID int
	)

	log.Info("AuditProjects start")
	for _, repo = range repos {
		log.Info("AuditProjects project [%s]", repo.Config.URL)
		if projID, err = g.ProjectID(repo.Config.URL); err != nil {
			log.Error("Failed to get project ID [%s], err: %+v", repo.Config.URL, err)
			continue
		}
		if err = g.AuditProjectHooks(projID, hooks); err != nil {
			log.Error("Failed to get hooks [%s], err: %+v", repo.Config.URL, err)
			continue
		}
		log.Info("AuditProjects project [%s]", repo.Config.URL)
	}
	log.Info("AuditProjects end")
	return
}

// AuditProjectHooks check and add or edit the project webhooks
func (g *Gitlab) AuditProjectHooks(projID int, hooks []*model.WebHook) (err error) {
	var (
		webHooks []*gitlab.ProjectHook
		h        *model.WebHook
	)

	if webHooks, err = g.ListProjectHook(projID); err != nil {
		return
	}

	// 配置文件中的webhook是否配置在项目中
	inWebHooks := func(h *model.WebHook) (in bool, ph *gitlab.ProjectHook) {
		for _, ph = range webHooks {
			if h.URL == ph.URL {
				in = true
				return
			}
		}
		in = false
		return
	}

	// 找到的webhook是否和配置文件中webhook配置相等
	equalHook := func(h *model.WebHook, ph *gitlab.ProjectHook) (equal bool) {
		if h.PushEvents == ph.PushEvents && h.IssuesEvents == ph.IssuesEvents &&
			h.ConfidentialIssuesEvents == ph.ConfidentialIssuesEvents && h.MergeRequestsEvents == ph.MergeRequestsEvents &&
			h.TagPushEvents == ph.TagPushEvents && h.NoteEvents == ph.NoteEvents &&
			h.JobEvents == ph.JobEvents && h.PipelineEvents == ph.PipelineEvents &&
			h.WikiPageEvents == ph.WikiPageEvents && !ph.EnableSSLVerification {
			return true
		}
		return false
	}

	// 配置的webhook与查询出来的对比，如果存在，则比较属性，属性不同，则修改；不存在，则创建
	for _, h = range hooks {
		if in, ph := inWebHooks(h); in {
			if !equalHook(h, ph) {
				if err = g.EditProjectHook(projID, ph.ID, h); err != nil {
					return
				}
			}
		} else {
			if err = g.AddProjectHook(projID, h); err != nil {
				return
			}
		}
	}

	return
}

// AddProjectHook add project hook for the project
func (g *Gitlab) AddProjectHook(projID int, hook *model.WebHook) (err error) {
	var (
		opt = &gitlab.AddProjectHookOptions{
			URL:                      &hook.URL,
			PushEvents:               &hook.PushEvents,
			IssuesEvents:             &hook.IssuesEvents,
			ConfidentialIssuesEvents: &hook.ConfidentialIssuesEvents,
			MergeRequestsEvents:      &hook.MergeRequestsEvents,
			TagPushEvents:            &hook.TagPushEvents,
			NoteEvents:               &hook.NoteEvents,
			JobEvents:                &hook.JobEvents,
			PipelineEvents:           &hook.PipelineEvents,
			WikiPageEvents:           &hook.WikiPageEvents,
			EnableSSLVerification:    gitlab.Bool(false),
		}
	)

	if _, _, err = g.client.Projects.AddProjectHook(projID, opt); err != nil {
		return
	}

	return
}

// ListProjectHook list all the webhook for the project
func (g *Gitlab) ListProjectHook(projID int) (hooks []*gitlab.ProjectHook, err error) {
	var (
		opt = &gitlab.ListProjectHooksOptions{
			Page:    1,
			PerPage: 100,
		}
	)
	if hooks, _, err = g.client.Projects.ListProjectHooks(projID, opt); err != nil {
		return
	}
	return
}

// EditProjectHook edit the specified webhook for the project
func (g *Gitlab) EditProjectHook(projID int, hookID int, hook *model.WebHook) (err error) {
	var (
		opt = &gitlab.EditProjectHookOptions{
			URL:                      &hook.URL,
			PushEvents:               &hook.PushEvents,
			IssuesEvents:             &hook.IssuesEvents,
			ConfidentialIssuesEvents: &hook.ConfidentialIssuesEvents,
			MergeRequestsEvents:      &hook.MergeRequestsEvents,
			TagPushEvents:            &hook.TagPushEvents,
			NoteEvents:               &hook.NoteEvents,
			JobEvents:                &hook.JobEvents,
			PipelineEvents:           &hook.PipelineEvents,
			WikiPageEvents:           &hook.WikiPageEvents,
			EnableSSLVerification:    gitlab.Bool(false),
		}
	)

	if _, _, err = g.client.Projects.EditProjectHook(projID, hookID, opt); err != nil {
		return
	}
	return
}

// DeletePojectHook delete the specified hook for the project
func (g *Gitlab) DeletePojectHook(projID int, hookID int) (err error) {
	if _, err = g.client.Projects.DeleteProjectHook(projID, hookID); err != nil {
		return
	}
	return
}

// ProjectID get project id by project URL-encoded name
func (g *Gitlab) ProjectID(url string) (projID int, err error) {
	var (
		projectName = url[strings.LastIndex(url, ":")+1 : strings.LastIndex(url, ".git")]
		project     *gitlab.Project
	)

	if project, _, err = g.client.Projects.GetProject(projectName); err != nil {
		return
	}
	projID = project.ID
	return
}

// MRPipelineStatus query PipelineState for mr
func (g *Gitlab) MRPipelineStatus(projID int, mrIID int) (id int, status string, err error) {
	var pipelineList gitlab.PipelineList
	if pipelineList, _, err = g.client.MergeRequests.ListMergeRequestPipelines(projID, mrIID); err != nil {
		return
	}
	if len(pipelineList) == 0 {
		status = model.PipelineSuccess
		return
	}
	id = pipelineList[0].ID
	status = pipelineList[0].Status
	return
}

// LastPipeLineState query Last PipeLineState
func (g *Gitlab) LastPipeLineState(projID int, ciCommitRefName string) (state bool, err error) {
	var pipelines gitlab.PipelineList
	opt := &gitlab.ListProjectPipelinesOptions{
		Ref: gitlab.String(ciCommitRefName),
	}
	state = true
	if pipelines, _, err = g.client.Pipelines.ListProjectPipelines(projID, opt); err != nil {
		return
	}
	if len(pipelines) < 2 {
		return
	}
	state = pipelines[1].Status == model.PipelineSuccess
	return
}

// MergeStatus query MergeStatus
func (g *Gitlab) MergeStatus(projID int, mrIID int) (wip bool, state string, status string, err error) {
	var mergerequest *gitlab.MergeRequest
	if mergerequest, _, err = g.client.MergeRequests.GetMergeRequest(projID, mrIID); err != nil {
		return
	}
	state = mergerequest.State
	status = mergerequest.MergeStatus
	wip = mergerequest.WorkInProgress
	return
}

// MergeLabels get Merge request labels
func (g *Gitlab) MergeLabels(projID int, mrIID int) (labels []string, err error) {
	var mergerequest *gitlab.MergeRequest
	if mergerequest, _, err = g.client.MergeRequests.GetMergeRequest(projID, mrIID); err != nil {
		return
	}
	labels = mergerequest.Labels
	return
}

// PlusUsernames get the username who +1
func (g *Gitlab) PlusUsernames(projID int, mrIID int) (usernames []string, err error) {
	var (
		opt = &gitlab.ListMergeRequestNotesOptions{
			Page:    1,
			PerPage: 100,
		}
		notes []*gitlab.Note
		exist bool
	)
	if notes, _, err = g.client.Notes.ListMergeRequestNotes(projID, mrIID, opt); err != nil {
		return
	}
	for _, note := range notes {
		if (strings.TrimSpace(note.Body) == model.SagaCommandPlusOne) || (strings.TrimSpace(note.Body) == model.SagaCommandPlusOne1) {
			exist = false
			for _, user := range usernames {
				if user == note.Author.Username {
					exist = true
					break
				}
			}
			if !exist {
				usernames = append(usernames, note.Author.Username)
			}
		}
	}
	return
}

// CompareDiff ...
func (g *Gitlab) CompareDiff(projID int, from string, to string) (files []string, err error) {
	var (
		opt = &gitlab.CompareOptions{
			From: &from,
			To:   &to,
		}
		compare *gitlab.Compare
	)
	if compare, _, err = g.client.Repositories.Compare(projID, opt); err != nil {
		return
	}
	for _, diff := range compare.Diffs {
		if diff.NewFile {
			files = append(files, diff.NewPath)
		} else {
			files = append(files, diff.OldPath)
		}
	}
	return
}

// CommitDiff ...
func (g *Gitlab) CommitDiff(projID int, sha string) (files []string, err error) {
	var (
		diffs []*gitlab.Diff
		opt   = &gitlab.GetCommitDiffOptions{
			Page:    1,
			PerPage: 100,
		}
	)
	if diffs, _, err = g.client.Commits.GetCommitDiff(projID, sha, opt); err != nil {
		return
	}
	for _, diff := range diffs {
		if diff.NewFile {
			files = append(files, diff.NewPath)
		} else {
			files = append(files, diff.OldPath)
		}
	}
	return
}

// MRChanges ...
func (g *Gitlab) MRChanges(projID, mrIID int) (changeFiles []string, deleteFiles []string, err error) {
	var (
		mr *gitlab.MergeRequest
	)
	if mr, _, err = g.client.MergeRequests.GetMergeRequestChanges(projID, mrIID); err != nil {
		return
	}
	for _, c := range mr.Changes {
		if c.RenamedFile {
			deleteFiles = append(deleteFiles, c.OldPath)
			changeFiles = append(changeFiles, c.NewPath)
		} else if c.DeletedFile {
			deleteFiles = append(deleteFiles, c.OldPath)
		} else if c.NewFile {
			changeFiles = append(changeFiles, c.NewPath)
		} else {
			changeFiles = append(changeFiles, c.OldPath)
		}
	}
	return
}

// RepoRawFile ...
func (g *Gitlab) RepoRawFile(projID int, branch string, filename string) (content []byte, err error) {
	var (
		opt = &gitlab.GetRawFileOptions{
			Ref: &branch,
		}
	)
	if content, _, err = g.client.RepositoryFiles.GetRawFile(projID, filename, opt); err != nil {
		return
	}
	return
}
