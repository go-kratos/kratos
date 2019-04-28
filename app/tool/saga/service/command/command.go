package command

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"go-common/app/tool/saga/conf"
	"go-common/app/tool/saga/dao"
	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/gitlab"
	"go-common/library/log"
)

// Command Command def.
type Command struct {
	dao    *dao.Dao
	gitlab *gitlab.Gitlab
	cmds   map[string]cmdFunc
}

type cmdFunc func(ctx context.Context, event *model.HookComment, repo *model.Repo) (err error)

// New ...
func New(dao *dao.Dao, gitlab *gitlab.Gitlab) (c *Command) {
	c = &Command{
		dao:    dao,
		gitlab: gitlab,
		cmds:   make(map[string]cmdFunc),
	}
	return
}

// Exec ...
func (c *Command) Exec(ctx context.Context, cmd string, event *model.HookComment, repo *model.Repo) (err error) {
	var (
		f      cmdFunc
		ok     bool
		projID = int(event.MergeRequest.SourceProjectID)
		mrIID  = int(event.MergeRequest.IID)
	)
	if f, ok = c.cmds[cmd]; !ok {
		return
	}
	if err = f(ctx, event, repo); err != nil {
		c.gitlab.CreateMRNote(projID, mrIID, fmt.Sprintf("<pre>SAGA 异常：%+v %s</pre>", err, debug.Stack()))
		return
	}
	return
}

func (c *Command) register(cmd string, f cmdFunc) {
	c.cmds[cmd] = f
}

// Registers ...
func (c *Command) Registers() {
	c.register(model.SagaCommandPlusOne, c.runPlusOne)
	c.register(model.SagaCommandMerge, c.runTryMerge)
	c.register(model.SagaCommandMerge1, c.runTryMerge)
	c.register(model.SagaCommandPlusOne1, c.runPlusOne)
}

// ListenTask ...
func (c *Command) ListenTask() {
	var (
		ctx       = context.TODO()
		err       error
		t         *time.Timer
		ok        bool
		taskInfos []*model.TaskInfo
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("ListenTask: %+v %s", x, debug.Stack())
		}
	}()

	t = time.NewTimer(time.Duration(conf.Conf.Property.TaskInterval))
	for range t.C {
		if _, taskInfos, err = c.dao.MergeTasks(ctx, model.TaskStatusWaiting); err != nil {
			log.Error("request MergeTasks: %+v", err)
			continue
		}
		for _, taskInfo := range taskInfos {
			if ok, err = c.dao.TryLock(ctx, fmt.Sprintf(model.SagaRepoLockKey, int(taskInfo.Event.ProjectID)), model.SagaLockValue, int(taskInfo.Repo.Config.LockTimeout)); err != nil {
				log.Error("TryLock ProjectID: %d, MRIID: %d, err: %+v", taskInfo.Event.ProjectID, taskInfo.Event.MergeRequest.IID, err)
				continue
			}
			if !ok {
				log.Info("TryLock ok: %t, ProjectID: %d, MRIID: %d", ok, taskInfo.Event.ProjectID, taskInfo.Event.MergeRequest.IID)
				continue
			}
			log.Info("request MRIID: %d, MergeTasks:%+v", taskInfo.Event.MergeRequest.IID, taskInfo)
			go func(taskInfo *model.TaskInfo) {
				var (
					projID = int(taskInfo.Event.MergeRequest.SourceProjectID)
					mrIID  = int(taskInfo.Event.MergeRequest.IID)
					branch = taskInfo.Event.MergeRequest.SourceBranch
				)
				if err = c.execMergeTask(taskInfo); err != nil {
					c.gitlab.UpdateMRNote(projID, mrIID, taskInfo.NoteID, fmt.Sprintf("<pre>SAGA 异常：%+v</pre>", err))
					if err = c.resetMergeStatus(projID, mrIID, branch, false); err != nil {
						log.Error("resetMergeStatus error: %+v", err)
					}
					log.Error("execMergeTask: %+v", err)
				}
			}(taskInfo)
		}
		t.Reset(time.Duration(conf.Conf.Property.TaskInterval))
	}
}
