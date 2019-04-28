package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AuthHub Auth Hub.
func (s *Service) AuthHub(c context.Context, sessionID string) (err error) {
	return s.dao.AuthHub(c, sessionID)
}

// AccessAuthHub Access Auth Hub.
func (s *Service) AccessAuthHub(c context.Context, username string) (accessHub bool, err error) {
	var (
		personProjectID int
	)

	if personProjectID, err = s.ProjectID(c, username); err != nil {
		return
	}

	if personProjectID > 0 {
		accessHub = true
	}

	return
}

// ProjectID Get Project ID.
func (s *Service) ProjectID(c context.Context, projectName string) (projectID int, err error) {
	var hubProjects []*model.HubProject

	if hubProjects, err = s.dao.HubProjects(c, projectName); err != nil {
		return
	}

	for _, hubProject := range hubProjects {
		if projectName == hubProject.Name {
			projectID = hubProject.ProjectID
			return
		}
	}
	return
}

// ProjectRepositories Get Project Repositories.
func (s *Service) ProjectRepositories(c context.Context, projectID, pn, ps int, keyWord string) (p *model.PaginateProjectRepoRecord, err error) {
	var (
		total        int
		projectRepos []*model.ProjectRepository
	)

	if total, err = s.dao.ProjectRepoCount(c, projectID); err != nil {
		return
	}

	if projectRepos, err = s.dao.ProjectRepos(c, projectID, pn, ps, keyWord); err != nil {
		return
	}

	p = &model.PaginateProjectRepoRecord{
		PageNum:           pn,
		PageSize:          ps,
		Total:             total,
		ProjectRepository: projectRepos,
	}
	return
}

// AccessPullProjects Access Pull Projects.
func (s *Service) AccessPullProjects(c context.Context, username string) (projectsName []string, err error) {
	var accessToHub bool

	if accessToHub, err = s.AccessAuthHub(c, username); err != nil {
		return
	}

	if accessToHub {
		projectsName = append(projectsName, s.c.BiliHub.HostName+"/"+username)
	}

	projectsName = append(projectsName, s.c.BiliHub.HostName+"/"+s.c.BiliHub.SharePub)

	return
}

// Projects Get Projects.
func (s *Service) Projects(c context.Context, username string) (projects map[string]int, err error) {
	var (
		personProjectID int
		merlinProjectID int
		publicProjectID int
	)

	projects = make(map[string]int)
	if personProjectID, err = s.ProjectID(c, username); err != nil {
		return
	}
	projects[username] = personProjectID

	if merlinProjectID, err = s.ProjectID(c, s.c.BiliHub.MerlinPub); err != nil {
		return
	}
	projects[s.c.BiliHub.MerlinPub] = merlinProjectID

	if publicProjectID, err = s.ProjectID(c, s.c.BiliHub.SharePub); err != nil {
		return
	}
	projects[s.c.BiliHub.SharePub] = publicProjectID

	return
}

// RepositoryTags Get Repository Tags.
func (s *Service) RepositoryTags(c context.Context, repoName string) (repoTags []*model.RepositoryTag, err error) {
	var (
		repoTagsRet         []*model.RepositoryTagResponse
		imageFullNamePrefix = s.c.BiliHub.HostName + "/" + repoName + ":"
	)
	if repoTagsRet, err = s.dao.RepoTags(c, repoName); err != nil {
		return
	}

	for _, repoTagRet := range repoTagsRet {
		// ignore image that size is 0
		if repoTagRet.Size == 0 {
			continue
		}

		repoTag := &model.RepositoryTag{
			RepositoryTagResponse: *repoTagRet,
			ImageFullName:         imageFullNamePrefix + repoTagRet.Name,
		}
		repoTags = append(repoTags, repoTag)
	}
	return
}

// DeleteRepositoryTag Delete Repository Tag.
func (s *Service) DeleteRepositoryTag(c context.Context, username, repoName, tagName string) (status int, err error) {
	var hasRight bool
	if hasRight = s.hasOperateHubRight(username, repoName); !hasRight {
		status = -1
		err = ecode.MerlinHubNoRight
		return
	}

	if err = s.dao.DeleteRepoTag(c, repoName, tagName); err != nil {
		status = -1
	}

	return
}

// DeleteRepository Delete Repository.
func (s *Service) DeleteRepository(c context.Context, username, repoName string) (status int, err error) {
	var hasRight bool
	if hasRight = s.hasOperateHubRight(username, repoName); !hasRight {
		status = -1
		err = ecode.MerlinHubNoRight
		return
	}

	if err = s.dao.DeleteRepo(c, repoName); err != nil {
		status = -1
	}

	return
}

// AddRepositoryTag Add Repository Tag.
func (s *Service) AddRepositoryTag(c context.Context, username, repoName, tagName, newRepoName, newTagName string) (status int, err error) {
	//操作src image 权限认证
	var hasRight bool
	if hasRight = s.hasOperateHubRight(username, repoName); !hasRight {
		status = -1
		err = ecode.MerlinHubNoRight
		return
	}

	//上传的src image权限认证。如果上传的路径 只能为公共路径和个人路径
	if strings.Split(newRepoName, "/")[0] != s.c.BiliHub.SharePub {
		if hasRight = s.hasOperateHubRight(username, newRepoName); !hasRight {
			status = -1
			err = ecode.MerlinHubNoRight
			return
		}
	}

	imageSrcName := s.getFullRepoName(repoName, tagName)
	imageTagName := s.getFullRepoName(newRepoName, newTagName)

	// pull and push image
	s.dao.ImageTask(func() {
		s.PullAndPush(username, imageSrcName, imageTagName, 0)
	})
	return
}

// AddTagToMachine Add Tag To Machine.
func (s *Service) AddTagToMachine(c context.Context, username, imageSrcName string, machineIDs []int64) (machineImageMap map[int64]string, err error) {
	machineImageMap = map[int64]string{}
	for _, machineID := range machineIDs {
		repoName := strings.Split(imageSrcName, ":")[0]
		absRepoName := strings.Replace(repoName, "/", "-", -1)
		newRepoName := s.c.BiliHub.MachineTagPri + "/" + absRepoName
		newTagName := strconv.FormatInt(machineID, 10) + "-" + time.Now().Format("20060102150405")
		imageTagName := s.getFullRepoName(newRepoName, newTagName)

		if _, err = s.PullAndPush(username, imageSrcName, imageTagName, machineID); err != nil {
			return
		}
		machineImageMap[machineID] = imageTagName
	}

	return
}

// GetAllImagesInDocker Get All Images In Docker.
func (s *Service) GetAllImagesInDocker() (imageNames []string, err error) {
	return s.dao.ImageGetAll()
}

// hasOperateHubRight has Operate HubRight.
func (s *Service) hasOperateHubRight(username, repoName string) (hasRight bool) {
	//判断是否又权限执行操作
	if username == strings.Split(repoName, "/")[0] {
		hasRight = true
		return
	}

	for _, super := range s.c.BiliHub.SuperOwner {
		if username == super {
			hasRight = true
			return
		}
	}

	return
}

// get Full RepoName getFullRepoName.
func (s *Service) getFullRepoName(repoName, tagName string) string {
	return s.c.BiliHub.HostName + "/" + repoName + ":" + tagName
}

//Push Push.
func (s *Service) Push(c context.Context, username, repoName, tagName string, machineID int64) (status int, err error) {
	imageSrcName := s.getFullRepoName(repoName, tagName)

	log.Info("start Push target %s", imageSrcName)

	if err = s.dao.ImagePush(imageSrcName); err != nil {
		status = model.ImagePushErr
		log.Error("ImagePush target %s,err (%+v)", imageSrcName, err)
	}

	hubImageLog := &model.HubImageLog{
		UserName:    username,
		MachineID:   machineID,
		ImageSrc:    "",
		ImageTag:    imageSrcName,
		Status:      status,
		OperateType: model.ImagePush,
	}
	s.dao.InsertHubImageLog(hubImageLog)

	log.Info("end Push target %s", imageSrcName)

	return
}

// ReTag ReTag.
func (s *Service) ReTag(c context.Context, username, repoName, tagName, newRepoName, newTagName string, machineID int64) (status int, err error) {
	imageSrcName := s.getFullRepoName(repoName, tagName)
	imageTagName := s.getFullRepoName(newRepoName, newTagName)

	log.Info("start ReTag source %s tag %s", imageSrcName, imageTagName)

	if err = s.dao.ImageTag(imageSrcName, imageTagName); err != nil {
		status = model.ImageReTagErr
		log.Error("ImageTag source %s, target %s,err (%+v)", imageSrcName, imageTagName, err)
	}
	hubImageLog := &model.HubImageLog{
		UserName:    username,
		MachineID:   machineID,
		ImageSrc:    imageSrcName,
		ImageTag:    imageTagName,
		Status:      status,
		OperateType: model.ImageTag,
	}
	s.dao.InsertHubImageLog(hubImageLog)

	log.Info("end ReTag source %s, tag %s", imageSrcName, imageTagName)

	return
}

// Pull Pull.
func (s *Service) Pull(c context.Context, username, repoName, tagName string, machineID int64) (status int, err error) {
	imageSrcName := s.getFullRepoName(repoName, tagName)

	log.Info("start Pull source %s", imageSrcName)

	if err = s.dao.ImagePull(imageSrcName); err != nil {
		status = model.ImagePullErr
		log.Error("ImagePull source %s,err (%+v)", imageSrcName, err)
	}

	hubImageLog := &model.HubImageLog{
		UserName:    username,
		MachineID:   machineID,
		ImageSrc:    imageSrcName,
		ImageTag:    "",
		Status:      status,
		OperateType: model.ImagePull,
	}
	s.dao.InsertHubImageLog(hubImageLog)

	log.Info("end Pull source %s", imageSrcName)

	return
}

// CreateSnapShot CreateSnapShot.
func (s *Service) CreateSnapShot(c context.Context, username string, machineID int64) (status int, err error) {
	status = -1
	//获取镜像名称
	var (
		machine           *model.Machine
		passMachineDetail *model.PaasMachineDetail
		pqadmr            *model.PaasQueryAndDelMachineRequest
		tmpSnapshotRecord *model.SnapshotRecord
	)
	if machine, err = s.dao.QueryMachine(machineID); err != nil {
		return
	}

	pqadmr = machine.ToPaasQueryAndDelMachineRequest()

	if passMachineDetail, err = s.dao.QueryPaasMachine(c, pqadmr); err != nil {
		return
	}
	snapshotRecord := &model.SnapshotRecord{
		MachineID: machineID,
		ImageName: passMachineDetail.Image,
		Username:  username,
		Status:    model.SnapshotInit,
	}

	if tmpSnapshotRecord, err = s.dao.FindSnapshotRecord(machineID); err != nil {
		return
	}

	if tmpSnapshotRecord.ID == 0 {
		//首次快照 插入
		if err = s.dao.InsertSnapshotRecord(snapshotRecord); err != nil {
			return
		}
	} else if tmpSnapshotRecord.Status == model.SnapshotDoing {
		//有快照记录，查看是否正在进行中
		err = ecode.MerlinSnapshotInDoingErr
		return
	} else if err = s.dao.UpdateSnapshotRecord(snapshotRecord); err != nil {
		return
	}

	//创建快照
	resultStatus := model.SnapshotDoing
	if _, err = s.dao.SnapshotPaasMachineStatus(c, pqadmr); err != nil {
		resultStatus = model.SnapShotFailed
	}

	if err = s.dao.UpdateSnapshotRecordStatus(machineID, resultStatus); err != nil {
		return
	}

	status = 0
	return
}

// QuerySnapShot Query SnapShot.
func (s *Service) QuerySnapShot(c context.Context, machineID int64) (snapshotRecord *model.SnapshotRecord, err error) {
	return s.dao.FindSnapshotRecord(machineID)
}

// QueryMachine2ImageLog Query Machine to ImageLog.
func (s *Service) QueryMachine2ImageLog(c context.Context, queryRequest *model.QueryMachine2ImageLogRequest) (p *model.PaginateHubImageLog, err error) {
	var (
		hubImageLogs []*model.HubImageLog
		total        int64
	)

	if total, hubImageLogs, err = s.dao.FindHubMachine2ImageLog(queryRequest); err != nil {
		return
	}
	p = &model.PaginateHubImageLog{
		PageNum:      queryRequest.PageNum,
		PageSize:     queryRequest.PageSize,
		Total:        total,
		HubImageLogs: hubImageLogs,
	}

	return
}

// CallBackSnapShot Call Back SnapShot.
func (s *Service) CallBackSnapShot(c context.Context, machineName, imageName, msg string, resultStatus bool) (err error) {
	var (
		machine              *model.Machine
		snapshotResultStatus string
	)
	if machine, err = s.dao.QueryOnlineMachineByName(machineName); err != nil {
		return
	}

	if resultStatus {
		snapshotResultStatus = model.SnapshotSuccess
	} else {
		snapshotResultStatus = model.SnapShotFailed
	}

	err = s.dao.UpdateSnapshotRecordStatus(machine.ID, snapshotResultStatus)
	return
}

// Machine2Image Machine to Image.
func (s *Service) Machine2Image(c context.Context, username, imageName, newImageName string, machineID int64) (err error) {
	var (
		accessToHub       bool
		hubImageLog       *model.HubImageLog
		hubImageLogs      []*model.HubImageLog
		machine           *model.Machine
		passMachineDetail *model.PaasMachineDetail
	)

	//判断镜像和机器是否一致
	if machine, err = s.dao.QueryMachine(machineID); err != nil {
		return
	}
	if passMachineDetail, err = s.dao.QueryPaasMachine(c, machine.ToPaasQueryAndDelMachineRequest()); err != nil {
		return
	}

	if passMachineDetail.Image != imageName {
		err = ecode.MerlinMachineImageNotSameErr
		return
	}

	//判断有无授权hub
	if accessToHub, err = s.AccessAuthHub(c, username); err != nil {
		return
	}

	if !accessToHub {
		err = ecode.MerlinHubNoRight
		return
	}

	//判断new image name是否重名
	if hubImageLog, err = s.dao.FindHubImageLogByImageTag(newImageName); err != nil {
		return
	}

	if hubImageLog.ID > 0 {
		err = ecode.MerlinDuplicateImageNameErr
		return
	}

	//判断该机器是否有正在进行的机器转镜像任务
	if hubImageLogs, err = s.dao.FindHubImageLogByMachineID(machineID); err != nil {
		return
	}

	for _, hil := range hubImageLogs {
		if hil.OperateType == model.ImageMachine2Image && hil.Status == model.ImageInit {
			err = ecode.MerlinMachine2ImageInDoingErr
			return
		}
	}

	status := model.ImageInit
	newHubImageLog := &model.HubImageLog{
		UserName:    username,
		MachineID:   machineID,
		ImageSrc:    imageName,
		ImageTag:    newImageName,
		Status:      status,
		OperateType: model.ImageMachine2Image,
	}

	if err = s.dao.InsertHubImageLog(newHubImageLog); err != nil {
		return
	}

	s.dao.ImageTask(func() {
		s.PullAndPushWithMachine2Image(username, imageName, newImageName, machineID, newHubImageLog.ID)
	})
	return
}

// PullAndPush pull And Push.
func (s *Service) PullAndPush(username, imageSrcName, imageTagName string, machineID int64) (status int, err error) {
	log.Info("start pullAndPush source %s, target %s", imageSrcName, imageTagName)

	//pull image
	if err = s.dao.ImagePull(imageSrcName); err != nil {
		status = model.ImagePullErr
		log.Error("ImagePull source %s,err (%+v)", imageSrcName, err)

	} else if err = s.dao.ImageTag(imageSrcName, imageTagName); err != nil {
		status = model.ImageReTagErr
		log.Error("ImageTag source %s, target %s,err (%+v)", imageSrcName, imageTagName, err)

	} else if err = s.dao.ImagePush(imageTagName); err != nil {
		status = model.ImagePushErr
		log.Error("ImagePush target %s,err (%+v)", imageTagName, err)
	}

	hubImageLog := &model.HubImageLog{
		UserName:    username,
		MachineID:   machineID,
		ImageSrc:    imageSrcName,
		ImageTag:    imageTagName,
		Status:      status,
		OperateType: model.ImagePullAndPush,
	}
	s.dao.InsertHubImageLog(hubImageLog)

	log.Info("end pullAndPush source %s, target %s", imageSrcName, imageTagName)

	return
}

// PullAndPushWithMachine2Image Pull And Push With Machine to Image.
func (s *Service) PullAndPushWithMachine2Image(username, imageSrcName, imageTagName string, machineID, hubImageLogID int64) (status int, err error) {
	log.Info("start PullAndPushWithMachine2Image source %s, target %s", imageSrcName, imageTagName)
	status = model.ImageSuccess
	//pull image
	if err = s.dao.ImagePull(imageSrcName); err != nil {
		status = model.ImagePullErr
		log.Error("ImagePull source %s,err (%+v)", imageSrcName, err)

	} else if err = s.dao.ImageTag(imageSrcName, imageTagName); err != nil {
		status = model.ImageReTagErr
		log.Error("ImageTag source %s, target %s,err (%+v)", imageSrcName, imageTagName, err)

	} else if err = s.dao.ImagePush(imageTagName); err != nil {
		status = model.ImagePushErr
		log.Error("ImagePush target %s,err (%+v)", imageTagName, err)
	}

	err = s.dao.UpdateHubImageLogStatus(hubImageLogID, status)
	log.Info("end PullAndPushWithMachine2Image source %s, target %s", imageSrcName, imageTagName)
	return
}

// Machine2ImageForceFailed Machine to Image Force  Failed.
func (s *Service) Machine2ImageForceFailed(c context.Context, machineID int64) (status int, err error) {
	if err = s.dao.UpdateHubImageLogStatusInDoingStatus(machineID, model.ImagePullErr); err != nil {
		status = -1
	}
	return
}

// UpdateImageConf Update Image Conf.
func (s *Service) UpdateImageConf(c context.Context, username string, ic *model.ImageConfiguration) (status int, err error) {
	var (
		hubImageConf *model.HubImageConf
		hasRight     bool
		envsJson     []byte
		hostsJson    []byte
	)
	status = -1

	// 超级用户才能改MerlinPub 配置模板
	ret := strings.Split(ic.ImageFullName, "/")
	if len(ret) > 1 && ret[1] == s.c.BiliHub.MerlinPub {
		for _, super := range s.c.BiliHub.SuperOwner {
			if username == super {
				hasRight = true
				continue
			}
		}

		if !hasRight {
			err = ecode.MerlinHubNoRight
			return
		}
	}

	if envsJson, err = json.Marshal(ic.Envs); err != nil {
		return
	}

	if hostsJson, err = json.Marshal(ic.HostAlias); err != nil {
		return
	}

	if hubImageConf, err = s.dao.FindHubImageConfByImageName(ic.ImageFullName); err != nil {
		return
	}

	newHubImageConf := &model.HubImageConf{
		ImageName: ic.ImageFullName,
		UpdateBy:  username,
		Command:   strings.TrimSpace(ic.Command),
		Envs:      string(envsJson),
		Hosts:     string(hostsJson),
	}

	if hubImageConf.ID == 0 {
		if err = s.dao.InsertHubImageConf(newHubImageConf); err != nil {
			return
		}
	} else {
		if err = s.dao.UpdateHubImageConf(newHubImageConf); err != nil {
			return
		}
	}

	status = 0
	return
}

// QueryImageConf  Query Image Conf.
func (s *Service) QueryImageConf(c context.Context, imageName string) (ic *model.ImageConfiguration, err error) {
	var (
		hubImageConf *model.HubImageConf
		envs         []*model.EnvVariable
		hostAlias    []*model.Host
	)

	if hubImageConf, err = s.dao.FindHubImageConfByImageName(imageName); err != nil || hubImageConf.ID == 0 {
		return
	}

	if err = json.Unmarshal([]byte(hubImageConf.Envs), &envs); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(hubImageConf.Hosts), &hostAlias); err != nil {
		return
	}

	ic = &model.ImageConfiguration{
		ImageFullName: imageName,
		PaasMachineSystem: model.PaasMachineSystem{
			Command:   hubImageConf.Command,
			Envs:      envs,
			HostAlias: hostAlias,
		},
	}

	if len(ic.Envs) == 0 && len(ic.HostAlias) == 0 && strings.TrimSpace(ic.Command) == "" {
		ic = nil
	}
	return
}
