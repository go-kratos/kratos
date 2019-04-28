package service

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// GenMachines create multiple machines.
func (s *Service) GenMachines(c context.Context, gmr *model.GenMachinesRequest, u string) (err error) {
	var (
		ins        []*model.CreateInstance
		cluster    *model.Cluster
		hasMachine bool
		ms         []*model.Machine
	)
	if hasMachine, err = s.dao.HasMachine(gmr.Name); err != nil {
		return
	}
	if hasMachine {
		err = ecode.MerlinDuplicateMachineNameErr
		return
	}
	if cluster, err = s.dao.QueryCluster(c, gmr.NetworkID); err != nil {
		return
	}
	if cluster == nil {
		err = ecode.MerlinInvalidClusterErr
		return
	}
	if err = cluster.Verify(); err != nil {
		return
	}
	pgmr := gmr.ToPaasGenMachineRequest(s.c.Paas.MachineLimitRatio)
	if ms, err = s.dao.InsertMachinesV2(u, gmr, pgmr); err != nil {
		return
	}

	//insert image log
	for _, m := range ms {
		hubImageLog := &model.HubImageLog{
			UserName:    u,
			MachineID:   m.ID,
			ImageSrc:    gmr.Image,
			ImageTag:    "",
			Status:      model.ImageSuccess,
			OperateType: model.ImageNoSnapshot,
		}
		s.dao.InsertHubImageLog(hubImageLog)
	}

	if ins, err = s.dao.GenPaasMachines(c, pgmr); err != nil {
		return
	}
	for _, in := range ins {
		if in.InstanceCreateStatus == model.CreateFailedMachineInPaas {
			s.dao.UpdateMachineStatusByName(model.ImmediatelyFailedMachineInMerlin, in.InstanceName)
		}
	}
	return
}

// DelMachineWhenCanBeDel Del Machine When Can Be Del.
func (s *Service) DelMachineWhenCanBeDel(c context.Context, id int64, username string) (instance *model.ReleaseInstance, err error) {
	var (
		machine *model.Machine
	)
	if machine, err = s.dao.QueryMachine(id); err != nil {
		return
	}

	if machine.IsCreating() {
		err = ecode.MerlinCanNotBeDel
		return
	}

	return s.DelMachine(c, id, username)

}

// DelMachine delete machine by giving id.
func (s *Service) DelMachine(c context.Context, id int64, username string, beforeDelMachineFuncs ...model.BeforeDelMachineFunc) (instance *model.ReleaseInstance, err error) {
	for _, beforeDelMachine := range beforeDelMachineFuncs {
		if err = beforeDelMachine(c, id, username); err != nil {
			return
		}
	}
	return s.delMachineWithoutComponent(c, id, username)
}

func (s *Service) delMachineWithoutComponent(c context.Context, mID int64, username string) (instance *model.ReleaseInstance, err error) {
	var (
		machine *model.Machine
	)
	if machine, err = s.dao.QueryMachine(mID); err != nil {
		return
	}
	machineLog := &model.MachineLog{
		OperateType: model.DeleteForMachineLog,
		Username:    username,
		MachineID:   mID,
	}
	if instance, err = s.delMachineWithoutLog(c, machine, username); err != nil {
		machineLog.OperateResult = model.OperationFailedForMachineLog
	} else {
		machineLog.OperateResult = model.OperationSuccessForMachineLog
	}
	machineLog.OperateType = model.DeleteForMachineLog
	machineLog.Username = username
	machineLog.MachineID = machine.ID
	err = s.dao.InsertMachineLog(machineLog)
	return
}

func (s *Service) delMachineWithoutLog(c context.Context, machine *model.Machine, username string) (instance *model.ReleaseInstance, err error) {
	if instance, err = s.dao.DelPaasMachine(c, machine.ToPaasQueryAndDelMachineRequest()); err != nil {
		return
	}
	if err = s.dao.UpdateTaskStatusByMachines([]int64{machine.ID}, model.TaskDone); err != nil {
		return
	}
	if err = s.dao.DelMachine(machine.ID, username); err != nil {
		return
	}
	if err = s.SendMailDeleteMachine(username, machine); err != nil {
		err = nil
		log.Error("Send mail failed (%v)", err)
	}
	return
}

// QueryMachineDetail query detail information of machine.
func (s *Service) QueryMachineDetail(c context.Context, mID int64) (detail model.MachineDetail, err error) {
	var (
		machine           *model.Machine
		passMachineDetail *model.PaasMachineDetail
		cluster           *model.Cluster
		nodes             []*model.MachineNode
		isSnapShot        bool
	)
	if machine, err = s.dao.QueryMachine(mID); err != nil {
		return
	}
	if passMachineDetail, err = s.dao.QueryPaasMachine(c, machine.ToPaasQueryAndDelMachineRequest()); err != nil {
		return
	}
	if cluster, err = s.dao.QueryCluster(c, machine.NetworkID); err != nil {
		return
	}
	if cluster == nil {
		err = ecode.MerlinInvalidClusterErr
		return
	}
	if err = cluster.Verify(); err != nil {
		return
	}
	if nodes, err = s.dao.QueryMachineNodes(mID); err != nil {
		return
	}

	isSnapShot = strings.Contains(passMachineDetail.Image, s.c.BiliHub.MachineTagPri)

	detail = model.MachineDetail{
		Machine:           *machine,
		Nodes:             nodes,
		PaasMachineDetail: passMachineDetail.ConvertUnits(),
		Name:              machine.Name,
		NetworkName:       cluster.Networks[0].Name,
		IsSnapShot:        isSnapShot,
	}
	return
}

// QueryMachinePackages query packages of machine.
func (s *Service) QueryMachinePackages(c context.Context) (mp []*model.MachinePackage, err error) {
	return s.dao.FindAllMachinePackages()
}

//QueryMachines query multiple machines and update machine status
func (s *Service) QueryMachines(c context.Context, session string, qmr *model.QueryMachineRequest) (p *model.PaginateMachine, err error) {
	var (
		machines     []*model.Machine
		total        int64
		cluster      *model.Cluster
		podNames     []string
		treeInstance *model.TreeInstance
		mapping      map[string]*model.TreeInstance
	)
	if mapping, err = s.QueryTreeInstanceForMerlin(c, session, &qmr.TreeNode); err != nil {
		return
	}
	for k := range mapping {
		podNames = append(podNames, k)
	}

	if total, machines, err = s.dao.QueryMachines(podNames, qmr); err != nil {
		return
	}

	genMachines := make([]model.GenMachine, len(machines))
	for i, m := range machines {
		if cluster, err = s.dao.QueryCluster(c, m.NetworkID); err != nil {
			return
		}
		genMachines[i] = model.GenMachine{
			Machine:     *m,
			ClusterName: cluster.Name,
			NetworkName: cluster.Networks[0].Name,
		}
		treeInstance = mapping[m.PodName]
		if treeInstance != nil {
			genMachines[i].IP = treeInstance.InternalIP
		}
	}
	p = &model.PaginateMachine{
		PageNum:  qmr.PageNum,
		PageSize: qmr.PageSize,
		Total:    total,
		Machines: genMachines,
	}
	return
}

// QueryMachineLogs query machine logs.
func (s *Service) QueryMachineLogs(c context.Context, queryRequest *model.QueryMachineLogRequest) (p *model.PaginateMachineLog, err error) {
	var (
		total       int64
		machineLogs []*model.AboundMachineLog
	)
	if total, machineLogs, err = s.dao.FindMachineLogs(queryRequest); err != nil {
		return
	}
	p = &model.PaginateMachineLog{
		PageNum:     queryRequest.PageNum,
		PageSize:    queryRequest.PageSize,
		Total:       total,
		MachineLogs: machineLogs,
	}
	return
}

// QueryMachineStatus query the status of machine.
func (s *Service) QueryMachineStatus(c context.Context, machineID int64) (msr *model.MachineStatusResponse, err error) {
	var m *model.Machine
	if m, err = s.dao.QueryMachine(machineID); err != nil {
		return
	}
	if m.IsFailed() || m.IsDeleted() {
		return
	}
	if msr, err = s.verifyPassStatus(c, m, true); err != nil {
		return
	}
	return
}

// TransferMachine Transfer Machine.
func (s *Service) TransferMachine(c context.Context, machineID int64, username, receiver string) (status int, err error) {
	var (
		machine                  *model.Machine
		treePath                 string
		treeRoles                []*model.TreeRole
		isReceiverAccessTreeNode bool
	)

	if _, err = s.dao.FindUserByUserName(receiver); err != nil {
		err = ecode.MerlinUserNotExist
		return
	}

	if machine, err = s.dao.QueryMachine(machineID); err != nil {
		return
	}

	// 查看接受者有无机器所在服务树节点权限
	treePath = machine.BusinessUnit + "." + machine.Project + "." + machine.App

	if treeRoles, err = s.dao.TreeRolesAsPlatform(c, treePath); err != nil {
		return
	}

	for _, treeRole := range treeRoles {
		if treeRole.UserName == receiver {
			isReceiverAccessTreeNode = true
			break
		}
	}

	if !isReceiverAccessTreeNode {
		err = ecode.MerlinUserNoAccessTreeNode
		return
	}

	machineLog := &model.MachineLog{
		OperateType: model.TransferForMachineLog,
		Username:    username,
		MachineID:   machineID,
	}

	if err = s.dao.UpdateMachineUser(machineID, receiver); err != nil {
		machineLog.OperateResult = model.OperationFailedForMachineLog

	} else {
		machineLog.OperateResult = model.OperationSuccessForMachineLog

		//send mail
		mailHeader := fmt.Sprintf("机器:[%s] 被用户[%s]从用户[%s]名下转移至用户[%s]", machine.Name, username, machine.Username, receiver)

		var sendMailUsers []string
		sendMailUsers = append(sendMailUsers, username)
		sendMailUsers = append(sendMailUsers, receiver)
		sendMailUsers = append(sendMailUsers, machine.Username)

		s.SendMailForMultiUsers(c, sendMailUsers, mailHeader)
	}

	err = s.dao.InsertMachineLog(machineLog)
	return
}

func (s *Service) verifyPassStatus(c context.Context, m *model.Machine, targetCreating bool) (msr *model.MachineStatusResponse, err error) {
	var (
		ms     *model.MachineStatus
		b      bool
		status int
		shadow *model.MachineStatusResponse
	)
	if msr = model.InstanceMachineStatusResponse(m.Status); msr != nil {
		return
	}
	if ms, err = s.dao.QueryPaasMachineStatus(c, m.ToPaasQueryAndDelMachineRequest()); err != nil {
		return
	}
	if b, err = s.TreeMachineIsExist(c, m.PodName, m.ToTreeNode()); err != nil {
		return
	}
	shadow = ms.ToMachineStatusResponse()
	if b {
		shadow.SynTree = model.True
	} else {
		shadow.SynTree = model.False
	}
	if targetCreating {
		status = shadow.CreatingMachineStatus()
	} else {
		status = shadow.FailedMachineStatus()
	}
	if err = s.dao.UpdateMachineStatus(m.ID, status); err != nil {
		return
	}
	msr = shadow
	return
}
