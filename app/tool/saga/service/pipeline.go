package service

import (
	"context"

	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/notification"
	"go-common/library/log"
)

// PipelineChanged handle pipeline changed webhook
func (s *Service) PipelineChanged(c context.Context, event *model.HookPipeline) (err error) {
	var (
		ok                bool
		repo              *model.Repo
		state             bool
		lastPipeLineState bool
		wip               bool
	)
	//get associated pipeline mr and check it's wip status, if wip is true and return
	if wip, err = s.checkMrStatus(c, event.Project.ID, event.ObjectAttributes.Ref); wip {
		log.Info("Pipeline associated mr is wip, project ID: [%d], branch: [%s]", event.Project.ID, event.ObjectAttributes.Ref)
		return
	}

	if ok, repo = s.findRepo(event.Project.GitSSHURL); !ok || !repo.Config.RelatePipeline {
		log.Info("PipelineChanged return repo: %s, ok: %t", event.Project.GitSSHURL, ok)
		return
	}

	if event.ObjectKind != model.HookPipelineType {
		log.Info("Pipeline hook object kind [%s] ignore", event.ObjectKind)
		return
	}

	if (event.ObjectAttributes.Status != model.PipelineFailed) && (event.ObjectAttributes.Status != model.PipelineSuccess) && (event.ObjectAttributes.Status != model.PipelineCanceled) {
		log.Info("Pipeline status [%s] ignore for project [%s]", event.ObjectAttributes.Status, event.Project.Name)
		return
	}

	go func() {
		if err = s.cmd.HookPipeline(event.Project.ID, event.ObjectAttributes.Ref, int(event.ObjectAttributes.ID)); err != nil {
			log.Error("CheckPipeline: %d %s %+v", event.Project.ID, event.ObjectAttributes.Ref, err)
		}
	}()

	// 查询上次pipeline状态，状态变化才发通知
	if lastPipeLineState, err = s.gitlab.LastPipeLineState(event.Project.ID, event.ObjectAttributes.Ref); err != nil {
		return
	}
	state = event.ObjectAttributes.Status == model.PipelineSuccess
	log.Info("status:%t, lastPipeLineState: %t", state, lastPipeLineState)
	if lastPipeLineState != state {
		go notification.MailPipeline(event)
		go notification.WechatPipeline(s.d, event)
	}
	return
}

func (s *Service) checkMrStatus(c context.Context, projectID int, branch string) (wip bool, err error) {
	var (
		ok        bool
		mergeInfo *model.MergeInfo
	)

	if ok, mergeInfo, err = s.d.MergeInfo(c, projectID, branch); err != nil || !ok {
		return
	}
	if wip, _, _, err = s.gitlab.MergeStatus(projectID, mergeInfo.MRIID); err != nil {
		return
	}
	return
}
