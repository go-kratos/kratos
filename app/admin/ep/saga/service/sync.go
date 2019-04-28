package service

import (
	"context"
	"fmt"
	"sync"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
	"go-common/library/log"
)

const _baseBranch = "master"

// syncalldataproc ...
func (s *Service) syncweekdataproc() {
	//the blow object will be sync or update all time data
	go s.SyncMember()
	go s.SyncRunners()
	go s.SyncContacts(context.TODO())
}

// syncalldataproc ...
func (s *Service) syncalldataproc() {
	if err := s.SyncBranch(); err != nil {
		return
	}
	// AggregateBranch base on branch, if SyncBranch err, so return
	s.AggregateBranch()
}

// syncdataproc ...
func (s *Service) syncdataproc() {
	var (
		wg   sync.WaitGroup
		done = func(f func() error) {
			defer wg.Done()
			f()
		}
	)

	wg.Add(2)
	go done(s.SyncIssues)
	go done(s.SyncPipelines)
	wg.Wait()

	wg.Add(1)
	go done(s.SyncCommit)
	wg.Wait()

	wg.Add(2)
	go done(s.SyncJobs)
	go done(s.SyncMR)
	wg.Wait()

	wg.Add(3)
	go done(s.SyncMRNote)
	go done(s.SyncMRAwardEmoji)
	go done(s.SyncMRDiscussion)
	wg.Wait()

	s.AggregateMR()
}

// SyncCommit ...
func (s *Service) SyncCommit() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		allPage    int
		allNum     int
		result     *model.SyncResult
	)
	log.Info("===================== SyncCommit start ========================")
	for _, projectID := range projectIDs {
		if result, err = s.SyncProjectCommit(projectID); err != nil {
			log.Error("SyncCommit projectID(%d), err(%+v)", projectID, err)
			go s.WechatFailData(model.DataTypeCommit, projectID, result, err)
			return
		}
		log.Info(">>>>>>>>> SyncCommit projectID(%d) complete Page(%d), Num(%d)", projectID, result.TotalPage, result.TotalNum)

		if result != nil && len(result.FailData) > 0 {
			go s.WechatFailData(model.DataTypeCommit, projectID, result, nil)
		}

		allPage = allPage + result.TotalPage
		allNum = allNum + result.TotalNum
	}
	log.Info("===================== SyncCommit finished totalPage(%d), totalNum(%d)========================", allPage, allNum)
	return
}

// SyncMR ...
func (s *Service) SyncMR() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		allPage    int
		allNum     int
		result     *model.SyncResult
	)
	log.Info("===================== SyncMR start ========================")
	for _, projectID := range projectIDs {
		if result, err = s.SyncProjectMR(context.TODO(), projectID); err != nil {
			log.Error("SyncMR projectID(%d), err(%+v)", projectID, err)
			go s.WechatFailData(model.DataTypeMR, projectID, result, err)
			return
		}
		log.Info(">>>>>>>>> SyncMR projectID(%d) complete Page(%d), Num(%d)", projectID, result.TotalPage, result.TotalNum)
		if result != nil && len(result.FailData) > 0 {
			go s.WechatFailData(model.DataTypeMR, projectID, result, nil)
		}
		allPage = allPage + result.TotalPage
		allNum = allNum + result.TotalNum
	}
	log.Info("===================== SyncMR finished totalPage(%d), totalNum(%d)========================", allPage, allNum)
	return
}

// AggregateMR ...
func (s *Service) AggregateMR() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
	)
	log.Info("===================== AggMR start ========================")
	for _, projectID := range projectIDs {
		if err = s.AggregateProjectMR(context.TODO(), projectID); err != nil {
			log.Error("AggMR projectID(%d), err(%+v)", projectID, err)
			return
		}
		log.Info(">>>>>>>>> AggMR projectID(%d) complete", projectID)
	}

	log.Info("===================== AggMR finished ========================")
	return
}

// SyncJobs ...
func (s *Service) SyncJobs() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		allPage    int
		allNum     int
		result     *model.SyncResult
	)
	log.Info("===================== SyncJobs start ========================")
	for _, projectID := range projectIDs {

		if result, err = s.SyncProjectJobs(projectID); err != nil {
			log.Error("SyncJobs projectID(%d), err(%+v)", projectID, err)
			go s.WechatFailData(model.DataTypeJob, projectID, result, err)
			return
		}
		log.Info(">>>>>>>>> SyncJobs projectID(%d) complete Page(%d), Num(%d)", projectID, result.TotalPage, result.TotalNum)

		if result != nil && len(result.FailData) > 0 {
			go s.WechatFailData(model.DataTypeJob, projectID, result, nil)
		}

		allPage = allPage + result.TotalPage
		allNum = allNum + result.TotalNum
	}
	log.Info("===================== SyncJobs finished totalPage(%d), totalNum(%d)========================", allPage, allNum)
	return
}

// SyncPipelines ...
func (s *Service) SyncPipelines() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		allPage    int
		allNum     int
		result     *model.SyncResult
	)
	log.Info("===================== SyncPipelines start ========================")
	for _, projectID := range projectIDs {

		if result, err = s.SyncProjectPipelines(projectID); err != nil {
			log.Error("SyncPipelines projectID(%d), err(%+v)", projectID, err)
			go s.WechatFailData(model.DataTypePipeline, projectID, result, err)
			return
		}
		log.Info(">>>>>>>>> SyncPipelines projectID(%d) complete Page(%d), Num(%d)", projectID, result.TotalPage, result.TotalNum)

		if result != nil && len(result.FailData) > 0 {
			go s.WechatFailData(model.DataTypePipeline, projectID, result, nil)
		}

		allPage = allPage + result.TotalPage
		allNum = allNum + result.TotalNum
	}
	log.Info("===================== SyncPipelines finished totalPage(%d), totalNum(%d)========================", allPage, allNum)
	return
}

// SyncMRNote ...
func (s *Service) SyncMRNote() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		totalPage  int
		totalNum   int
		page       int
		num        int
	)
	log.Info("===================== SyncNote start ========================")
	for _, projectID := range projectIDs {
		if page, num, err = s.SyncProjectNotes(context.TODO(), projectID); err != nil {
			log.Error("SyncNotes projectID(%d), err(%+v)", projectID, err)
			return
		}
		log.Info(">>>>>>>>> SyncNote projectID(%d) complete Page(%d), Num(%d)", projectID, page, num)
		totalPage = totalPage + page
		totalNum = totalNum + num
	}
	log.Info("===================== SyncNote finished totalPage(%d), totalNum(%d)========================", totalPage, totalNum)
	return
}

// SyncMember ...
func (s *Service) SyncMember() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		totalPage  int
		totalNum   int
		page       int
		num        int
	)
	log.Info("===================== SyncMember start ========================")
	for _, projectID := range projectIDs {
		if page, num, err = s.SyncProjectMember(context.TODO(), projectID); err != nil {
			log.Error("SyncMember projectID(%d), err(%+v)", projectID, err)
			return
		}
		log.Info(">>>>>>>>> SyncMember projectID(%d) complete Page(%d), Num(%d)", projectID, page, num)
		totalPage = totalPage + page
		totalNum = totalNum + num
	}
	log.Info("===================== SyncMember finished totalPage(%d), totalNum(%d)========================", totalPage, totalNum)
	return
}

// SyncMRAwardEmoji ...
func (s *Service) SyncMRAwardEmoji() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		totalPage  int
		totalNum   int
		page       int
		num        int
	)
	log.Info("===================== SyncMRAwardEmoji start ========================")
	for _, projectID := range projectIDs {
		if page, num, err = s.SyncProjectAwardEmoji(context.TODO(), projectID); err != nil {
			log.Error("SyncMRAwardEmoji projectID(%d), err(%+v)", projectID, err)
			return
		}
		log.Info(">>>>>>>>> SyncMRAwardEmoji projectID(%d) complete Page(%d), Num(%d)", projectID, page, num)
		totalPage = totalPage + page
		totalNum = totalNum + num
	}
	log.Info("===================== SyncMRAwardEmoji finished totalPage(%d), totalNum(%d)========================", totalPage, totalNum)
	return
}

// SyncMRDiscussion ...
func (s *Service) SyncMRDiscussion() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		totalPage  int
		totalNum   int
		page       int
		num        int
	)
	log.Info("===================== SyncMRDiscussion start ========================")
	for _, projectID := range projectIDs {
		if page, num, err = s.SyncProjectDiscussion(context.TODO(), projectID); err != nil {
			log.Error("SyncMRDiscussion projectID(%d), err(%+v)", projectID, err)
			return
		}
		log.Info(">>>>>>>>> SyncMRDiscussion projectID(%d) complete Page(%d), Num(%d)", projectID, page, num)
		totalPage = totalPage + page
		totalNum = totalNum + num
	}
	log.Info("===================== SyncMRDiscussion finished totalPage(%d), totalNum(%d)========================", totalPage, totalNum)
	return
}

// SyncRunners ...
func (s *Service) SyncRunners() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		totalPage  int
		totalNum   int
		page       int
		num        int
	)
	log.Info("===================== SyncRunners start ========================")
	for _, projectID := range projectIDs {
		if page, num, err = s.SyncAllRunners(projectID); err != nil {
			log.Error("SyncRunners projectID(%d), err(%+v)", projectID, err)
			return
		}
		log.Info(">>>>>>>>> SyncRunners projectID(%d) complete Page(%d), Num(%d)", projectID, page, num)
		totalPage = totalPage + page
		totalNum = totalNum + num
	}
	log.Info("===================== SyncRunners finished totalPage(%d), totalNum(%d)========================", totalPage, totalNum)
	return
}

// SyncIssues ...
func (s *Service) SyncIssues() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		totalPage  int
		totalNum   int
		page       int
		num        int
	)
	log.Info("===================== SyncIssues start ========================")
	for _, projectID := range projectIDs {

		if page, num, err = s.SyncAllIssues(projectID); err != nil {
			log.Error("SyncIssues projectID(%d), err(%+v)", projectID, err)
			return
		}
		log.Info(">>>>>>>>> SyncIssues projectID(%d) complete Page(%d), Num(%d)", projectID, page, num)
		totalPage = totalPage + page
		totalNum = totalNum + num
	}
	log.Info("===================== SyncIssues finished totalPage(%d), totalNum(%d)========================", totalPage, totalNum)
	return
}

// SyncBranch ...
func (s *Service) SyncBranch() (err error) {
	var (
		projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs
		allPage    int
		allNum     int
		result     *model.SyncResult
	)
	log.Info("===================== SyncBranch start ========================")
	for _, projectID := range projectIDs {
		if result, err = s.SyncProjectBranch(context.TODO(), projectID); err != nil {
			log.Error("SyncBranch projectID(%d), err(%+v)", projectID, err)
			go s.WechatFailData(model.DataTypeBranch, projectID, result, err)
			return
		}
		log.Info(">>>>>>>>> SyncBranch projectID(%d) complete Page(%d), Num(%d)", projectID, result.TotalPage, result.TotalNum)
		if result != nil && len(result.FailData) > 0 {
			go s.WechatFailData(model.DataTypeBranch, projectID, result, nil)
		}
		allPage = allPage + result.TotalPage
		allNum = allNum + result.TotalNum
	}
	log.Info("===================== SyncBranch finished totalPage(%d), totalNum(%d)========================", allPage, allNum)
	return
}

// AggregateBranch ...
func (s *Service) AggregateBranch() (err error) {
	var projectIDs = conf.Conf.Property.DefaultProject.ProjectIDs

	log.Info("===================== SyncAggregateBranch start ========================")
	for _, projectID := range projectIDs {
		if err = s.AggregateProjectBranch(context.TODO(), projectID, _baseBranch); err != nil {
			log.Error("SyncAggregateBranch projectID(%d), err(%+v)", projectID, err)
			return
		}
	}
	log.Info("===================== SyncAggregateBranch finished ========================")
	return
}

// WechatFailData ...
func (s *Service) WechatFailData(dataType string, projectID int, result *model.SyncResult, error error) (err error) {
	var (
		ctx     = context.Background()
		users   = conf.Conf.Property.SyncData.WechatUser
		content string
	)

	title := "[SAGA-ADMIN] SYNC ERROR !!!"
	subject := fmt.Sprintf("%s\n\nDataType ( %s )\nProjectID ( %d )", title, dataType, projectID)

	if error != nil {
		return s.wechat.PushMsg(ctx, users, fmt.Sprintf("%s\n\n%s", subject, error.Error()))
	}

	if result == nil || len(result.FailData) <= 0 {
		return
	}

	content = "LIST : \n"
	for _, data := range result.FailData {

		if dataType == model.DataTypeJob || dataType == model.DataTypePipeline || dataType == model.DataTypeMR {
			content += fmt.Sprintf("%d,  ", data.ChildID)
		} else if dataType == model.DataTypeCommit || dataType == model.DataTypeBranch {
			content += fmt.Sprintf("%s\n", data.ChildIDStr)
		}
	}

	return s.wechat.PushMsg(ctx, users, fmt.Sprintf("%s\n\n%s", subject, content))
}
