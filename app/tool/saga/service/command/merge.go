package command

import (
	"context"
	"fmt"
	"runtime/debug"

	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/notification"
	"go-common/library/log"

	ggitlab "github.com/xanzy/go-gitlab"
)

func (c *Command) runTryMerge(ctx context.Context, event *model.HookComment, repo *model.Repo) (err error) {
	var (
		ok       bool
		canMerge bool
		projID   = int(event.MergeRequest.SourceProjectID)
		mrIID    = int(event.MergeRequest.IID)
		wip      = event.MergeRequest.WorkInProgress
		noteID   int
		taskInfo = &model.TaskInfo{
			Event: event,
			Repo:  repo,
		}
	)
	log.Info("runTryMerge start ... MRIID: %d, Repo Config: %+v", mrIID, repo.Config)

	if ok, err = c.dao.ExistMRIID(ctx, mrIID); err != nil || ok {
		return
	}
	if noteID, err = c.gitlab.CreateMRNote(projID, mrIID, "<pre>SAGA 开始执行，请大佬稍后......</pre>"); err != nil {
		return
	}
	taskInfo.NoteID = noteID

	// 1, check wip
	if wip {
		c.gitlab.UpdateMRNote(projID, mrIID, noteID, "<pre>警告：当前MR处于WIP状态，请待开发结束后再merge！</pre>")
		return
	}
	// 2, check labels
	if ok, err = c.checkLabels(projID, mrIID, noteID, repo); err != nil || !ok {
		return
	}

	// 3, check merge status
	if canMerge, err = c.checkMergeStatus(projID, mrIID, noteID); err != nil || !canMerge {
		return
	}

	// 4, check pipeline status
	if repo.Config.RelatePipeline {
		if repo.Config.DelayMerge {
			if ok, _, err = c.checkPipeline(projID, mrIID, noteID, 0, model.QueryProcessing); err != nil || !ok {
				return
			}
		} else {
			if ok, _, err = c.checkPipeline(projID, mrIID, noteID, 0, model.QuerySuccess); err != nil || !ok {
				return
			}
		}
	}

	// 5, check path auth
	if ok, err = c.checkAllPathAuth(taskInfo); err != nil || !ok {
		return
	}

	// 6, show current mr queue info
	c.showMRQueueInfo(ctx, taskInfo)

	if err = c.dao.PushMergeTask(ctx, model.TaskStatusWaiting, taskInfo); err != nil {
		return
	}
	if err = c.dao.AddMRIID(ctx, mrIID, int(repo.Config.LockTimeout)); err != nil {
		return
	}
	log.Info("runTryMerge merge task 已加入 waiting 任务列队中... MRIID: %d", mrIID)
	return
}

func (c *Command) execMergeTask(taskInfo *model.TaskInfo) (err error) {
	var (
		ctx          = context.TODO()
		projID       = int(taskInfo.Event.MergeRequest.SourceProjectID)
		mrIID        = int(taskInfo.Event.MergeRequest.IID)
		sourceBranch = taskInfo.Event.MergeRequest.SourceBranch
		pipeline     = &ggitlab.Pipeline{}
		noteID       = taskInfo.NoteID
		mergeInfo    = &model.MergeInfo{
			ProjID:       projID,
			MRIID:        mrIID,
			URL:          taskInfo.Event.ObjectAttributes.URL,
			AuthBranches: taskInfo.Repo.Config.AuthBranches,
			SourceBranch: taskInfo.Event.MergeRequest.SourceBranch,
			TargetBranch: taskInfo.Event.MergeRequest.TargetBranch,
			AuthorID:     int(taskInfo.Event.MergeRequest.AuthorID),
			UserName:     taskInfo.Event.User.UserName,
			MinReviewer:  taskInfo.Repo.Config.MinReviewer,
			LockTimeout:  taskInfo.Repo.Config.LockTimeout,
			Title:        taskInfo.Event.MergeRequest.Title,
			Description:  taskInfo.Event.MergeRequest.Description,
		}
	)
	mergeInfo.NoteID = noteID
	// 从等待任务列队移除
	if err = c.dao.DeleteMergeTask(ctx, model.TaskStatusWaiting, taskInfo); err != nil {
		return
	}
	// 加入到正在执行任务列队
	if err = c.dao.PushMergeTask(ctx, model.TaskStatusRunning, taskInfo); err != nil {
		return
	}

	if taskInfo.Repo.Config.RelatePipeline {
		if taskInfo.Repo.Config.DelayMerge {
			if err = c.HookDelayMerge(projID, sourceBranch, mergeInfo); err != nil {
				return
			}
			return
		}

		if err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, "<pre>SAGA 提示：为了保证合进主干后能正常编译，正在重跑pipeline，等待时间取决于pipeline运行时间！请耐心等待！</pre>"); err != nil {
			return
		}
		if pipeline, err = c.retryPipeline(taskInfo.Event); err != nil {
			return
		}

		mergeInfo.PipelineID = pipeline.ID
		if err = c.dao.SetMergeInfo(ctx, projID, sourceBranch, mergeInfo); err != nil {
			return
		}
	} else {
		if err = c.HookMerge(projID, sourceBranch, mergeInfo); err != nil {
			return
		}
	}
	return
}

func (c *Command) retryPipeline(event *model.HookComment) (pipeline *ggitlab.Pipeline, err error) {
	var (
		trigger      *ggitlab.PipelineTrigger
		triggers     []*ggitlab.PipelineTrigger
		projID       = int(event.MergeRequest.SourceProjectID)
		sourceBranch = event.MergeRequest.SourceBranch
	)
	if triggers, err = c.gitlab.Triggers(projID); err != nil {
		return
	}
	if len(triggers) == 0 {
		log.Info("No triggers were found for project %d, try to create it now.", projID)
		if trigger, err = c.gitlab.CreateTrigger(projID); err != nil {
			return
		}
		triggers = []*ggitlab.PipelineTrigger{trigger}
	}
	trigger = triggers[0]
	if trigger.Owner == nil || trigger.Owner.ID == 0 {
		log.Info("Legacy trigger (without owner), take ownership now.")
		if trigger, err = c.gitlab.TakeOwnership(projID, trigger.ID); err != nil {
			return
		}
	}
	if pipeline, err = c.gitlab.TriggerPipeline(projID, sourceBranch, trigger.Token); err != nil {
		return
	}
	return
}

// HookPipeline ...
func (c *Command) HookPipeline(projID int, branch string, pipelineID int) (err error) {
	var (
		ok        bool
		canMerge  bool
		mergeInfo *model.MergeInfo
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("HookPipeline: %+v %s", x, debug.Stack())
		}
	}()

	if ok, mergeInfo, err = c.dao.MergeInfo(context.TODO(), projID, branch); err != nil || !ok {
		return
	}
	log.Info("HookPipeline projID: %d, MRIID: %d, branch: %s, pipelineId: %d", projID, mergeInfo.MRIID, branch, mergeInfo.PipelineID)
	if pipelineID < mergeInfo.PipelineID {
		return
	}

	defer func() {
		if err = c.resetMergeStatus(projID, mergeInfo.MRIID, branch, true); err != nil {
			log.Error("resetMergeStatus MRIID: %d, error: %+v", mergeInfo.MRIID, err)
		}
	}()

	// 1, check pipeline id
	if ok, _, err = c.checkPipeline(projID, mergeInfo.MRIID, mergeInfo.NoteID, mergeInfo.PipelineID, model.QueryID); err != nil || !ok {
		return
	}
	// 2, check pipeline status
	if ok, _, err = c.checkPipeline(projID, mergeInfo.MRIID, mergeInfo.NoteID, 0, model.QuerySuccess); err != nil || !ok {
		return
	}
	// 3, check merge status
	if canMerge, err = c.checkMergeStatus(projID, mergeInfo.MRIID, mergeInfo.NoteID); err != nil || !canMerge {
		return
	}

	log.Info("HookPipeline acceptMerge ... MRIID: %d", mergeInfo.MRIID)
	if ok, err = c.acceptMerge(mergeInfo); err != nil || !ok {
		return
	}
	return
}

// HookMerge ...
func (c *Command) HookMerge(projID int, branch string, mergeInfo *model.MergeInfo) (err error) {
	var (
		ok       bool
		canMerge bool
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("HookMerge: %+v %s", x, debug.Stack())
		}
	}()
	defer func() {
		if err = c.resetMergeStatus(projID, mergeInfo.MRIID, branch, true); err != nil {
			log.Error("resetMergeStatus MRIID: %d, error: %+v", mergeInfo.MRIID, err)
		}
	}()

	log.Info("HookMerge projID: %d, MRIID: %d, branch: %s", projID, mergeInfo.MRIID, branch)
	if canMerge, err = c.checkMergeStatus(projID, mergeInfo.MRIID, mergeInfo.NoteID); err != nil || !canMerge {
		return
	}

	log.Info("HookMerge acceptMerge ... MRIID: %d", mergeInfo.MRIID)
	if ok, err = c.acceptMerge(mergeInfo); err != nil || !ok {
		return
	}
	return
}

// HookDelayMerge ...
func (c *Command) HookDelayMerge(projID int, branch string, mergeInfo *model.MergeInfo) (err error) {
	var (
		ctx        = context.TODO()
		ok         bool
		noteID     = mergeInfo.NoteID
		mrIID      = mergeInfo.MRIID
		pipelineID int
		status     string
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("HookDelayMerge: %+v %s", x, debug.Stack())
		}
	}()

	//if ok, pipelineID, err = c.checkPipeline(projID, mrIID, noteID, 0, model.QuerySuccessRmNote); err != nil {
	//return
	//}
	if pipelineID, status, err = c.gitlab.MRPipelineStatus(projID, mrIID); err != nil {
		return
	}
	if status == model.PipelineSuccess || status == model.PipelineSkipped {
		ok = true
	} else if status != model.PipelineRunning && status != model.PipelinePending {
		comment := fmt.Sprintf("<pre>警告：pipeline状态异常，请确保pipeline状态正常后再执行merge操作！</pre>")
		err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
		return
	}

	log.Info("HookDelayMerge projID: %d, MRIID: %d, branch: %s, pipeline status: %t", projID, mergeInfo.MRIID, branch, ok)
	if ok {
		if err = c.HookMerge(projID, branch, mergeInfo); err != nil {
			return
		}
	} else {
		mergeInfo.PipelineID = pipelineID
		if err = c.dao.SetMergeInfo(ctx, projID, branch, mergeInfo); err != nil {
			return
		}
	}
	return
}

func (c *Command) acceptMerge(mergeInfo *model.MergeInfo) (ok bool, err error) {
	var (
		comment      string
		author       string
		canMerge     bool
		state        string
		authorID     = mergeInfo.AuthorID
		username     = mergeInfo.UserName
		projID       = mergeInfo.ProjID
		mrIID        = mergeInfo.MRIID
		url          = mergeInfo.URL
		sourceBranch = mergeInfo.SourceBranch
		targetBranch = mergeInfo.TargetBranch
		noteID       = mergeInfo.NoteID
		content      = mergeInfo.Title
	)
	if author, err = c.gitlab.UserName(authorID); err != nil {
		return
	}
	if canMerge, err = c.checkMergeStatus(projID, mrIID, noteID); err != nil {
		return
	}
	if !canMerge {
		go notification.WechatAuthor(c.dao, author, url, sourceBranch, targetBranch, comment)
		return
	}

	if len(mergeInfo.Description) > 0 {
		content = content + "\n\n" + mergeInfo.Description
	}
	mergeMSG := fmt.Sprintf("Merge branch [%s] into [%s] by [%s]\n%s", sourceBranch, targetBranch, username, content)
	if state, err = c.gitlab.AcceptMR(projID, mrIID, mergeMSG); err != nil || state != model.MRStateMerged {
		if err != nil {
			comment = fmt.Sprintf("<pre>[%s]尝试合并失败，当前状态不允许合并，请查看上方merge按钮旁的提示！</pre>", username)
		} else {
			comment = fmt.Sprintf("<pre>[%s]尝试合并失败，请检查当前状态或同步目标分支代码后再试！</pre>", username)
		}
		go notification.WechatAuthor(c.dao, author, url, sourceBranch, targetBranch, comment)
		c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
		return
	}
	ok = true
	comment = fmt.Sprintf("<pre>[%s]尝试合并成功！</pre>", username)
	go notification.WechatAuthor(c.dao, author, url, sourceBranch, targetBranch, comment)
	c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
	return
}

func (c *Command) resetMergeStatus(projID int, MRIID int, branch string, taskRunning bool) (err error) {
	var (
		ctx = context.TODO()
	)

	log.Info("resetMergeStatus projID: %d, MRIID: %d start", projID, MRIID)
	if err = c.dao.UnLock(ctx, fmt.Sprintf(model.SagaRepoLockKey, projID)); err != nil {
		log.Error("UnLock error: %+v", err)
	}
	if err = c.dao.DeleteMergeInfo(ctx, projID, branch); err != nil {
		log.Error("DeleteMergeInfo error: %+v", err)
	}
	if err = c.dao.DeleteMRIID(ctx, MRIID); err != nil {
		log.Error("Delete MRIID :%d, error: %+v", MRIID, err)
	}
	if taskRunning {
		if err = c.DeleteRunningTask(projID, MRIID); err != nil {
			log.Error("DeleteRunningTask: %+v", err)
		}
	}
	log.Info("resetMergeStatus projID: %d, MRIID: %d end!", projID, MRIID)
	return
}

// DeleteRunningTask ...
func (c *Command) DeleteRunningTask(projID int, mrID int) (err error) {
	var (
		ctx       = context.TODO()
		taskInfos []*model.TaskInfo
	)

	if _, taskInfos, err = c.dao.MergeTasks(ctx, model.TaskStatusRunning); err != nil {
		return
	}
	for _, taskInfo := range taskInfos {
		pID := int(taskInfo.Event.MergeRequest.SourceProjectID)
		mID := int(taskInfo.Event.MergeRequest.IID)
		if pID == projID && mID == mrID {
			// 从正在运行的任务列队中移除
			err = c.dao.DeleteMergeTask(ctx, model.TaskStatusRunning, taskInfo)
			return
		}
	}
	return
}

func (c *Command) checkMergeStatus(projID int, mrIID int, noteID int) (canMerge bool, err error) {
	var (
		wip     bool
		state   string
		status  string
		comment string
	)
	if wip, state, status, err = c.gitlab.MergeStatus(projID, mrIID); err != nil {
		return
	}
	if wip {
		comment = "<pre>SAGA 尝试合并失败，当前MR是一项正在进行的工作！若已完成请先点击“Resolve WIP status”按钮处理后再+merge！</pre>"
	} else if state != model.MergeStateOpened {
		comment = "<pre>SAGA 尝试合并失败，当前MR已经关闭或者已经合并！</pre>"
	} else if status != model.MergeStatusOk {
		comment = "<pre>SAGA 尝试合并失败，请先解决合并冲突！</pre>"
	} else {
		canMerge = true
	}

	if len(comment) > 0 {
		c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
	}

	return
}

// checkLabels ios or android need checkout label when release app stage
func (c *Command) checkLabels(projID int, mrIID int, noteID int, repo *model.Repo) (ok bool, err error) {
	var (
		labels  []string
		comment = fmt.Sprintf("<pre>警告：SAGA 无法执行+merge，发版阶段只允许合入指定label的MR！</pre>")
	)
	if len(repo.Config.AllowLabel) <= 0 {
		ok = true
		return
	}

	if labels, err = c.gitlab.MergeLabels(projID, mrIID); err != nil {
		return
	}

	log.Info("checkMrLabels MRIID: %d, labels: %+v", mrIID, labels)
	for _, label := range labels {
		if label == repo.Config.AllowLabel {
			ok = true
			return
		}
	}
	if err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment); err != nil {
		return
	}
	return
}

func (c *Command) checkPipeline(projID int, mrIID int, noteID int, lastPipelineID int, queryStatus model.QueryStatus) (ok bool, pipelineID int, err error) {
	var status string
	if pipelineID, status, err = c.gitlab.MRPipelineStatus(projID, mrIID); err != nil {
		return
	}
	log.Info("checkPipeline MRIID: %d, queryStatus: %d, pipeline status: %s", mrIID, queryStatus, status)

	// query pipeline id index
	if queryStatus == model.QueryID {
		if pipelineID > lastPipelineID {
			comment := fmt.Sprintf("<pre>警告：SAGA 检测到重新提交代码了，+merge中断！请重新review代码！</pre>")
			err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
			return
		}
		ok = true
		return
	}

	// query process status
	if queryStatus == model.QueryProcessing {
		if status == model.PipelineRunning || status == model.PipelinePending {
			comment := fmt.Sprintf("<pre>警告：pipeline正在运行中，暂不能立即merge，待pipeline运行通过后会自动执行merge操作！</pre>")
			err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
			ok = true
			return
		} else if status == model.PipelineSuccess || status == model.PipelineSkipped {
			ok = true
			return
		}
		comment := fmt.Sprintf("<pre>警告：pipeline状态异常，请确保pipeline状态正常后再执行merge操作！</pre>")
		err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
		return
	}

	// query success status
	if queryStatus == model.QuerySuccess {
		if status != model.PipelineSuccess {
			comment := fmt.Sprintf("<pre>警告：SAGA 无法执行+merge，pipeline还未成功，请大佬先让pipeline执行通过！</pre>")
			err = c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
			ok = false
			return
		}
	}

	ok = true
	return
}

// showMRQueueInfo ...
func (c *Command) showMRQueueInfo(ctx context.Context, taskInfo *model.TaskInfo) (err error) {
	var (
		mrIID      = int(taskInfo.Event.MergeRequest.IID)
		projID     = int(taskInfo.Event.MergeRequest.SourceProjectID)
		noteID     = taskInfo.NoteID
		taskInfos  []*model.TaskInfo
		comment    string
		waitNum    int
		runningNum int
	)

	if _, taskInfos, err = c.dao.MergeTasks(ctx, model.TaskStatusWaiting); err != nil {
		return
	}
	for _, waitTaskInfo := range taskInfos {
		if waitTaskInfo.Event.ProjectID == taskInfo.Event.ProjectID {
			waitNum++
		}
	}

	if _, taskInfos, err = c.dao.MergeTasks(ctx, model.TaskStatusRunning); err != nil {
		return
	}
	for _, runningTaskInfo := range taskInfos {
		if runningTaskInfo.Event.ProjectID == taskInfo.Event.ProjectID {
			runningNum++
		}
	}

	if waitNum > 0 {
		comment = fmt.Sprintf("<pre>SAGA 提示：当前还有 [%d] 个 MR 等待合并，请大佬耐心等待！</pre>", waitNum)
		c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
	} else if runningNum > 0 {
		comment = fmt.Sprintf("<pre>SAGA 提示：当前还有 [%d] 个 MR 正在执行，请大佬耐心等待！</pre>", runningNum)
		c.gitlab.UpdateMRNote(projID, mrIID, noteID, comment)
	}
	return
}
