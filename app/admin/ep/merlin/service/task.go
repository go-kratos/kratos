package service

import (
	"context"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/log"
)

func (s *Service) taskGetExpiredMachinesIntoTask() {
	machines, err := s.dao.FindExpiredMachine()
	if err != nil {
		log.Error("Task get expired machines into task (%v)", err)
		return
	}

	if machines != nil {
		log.Info("machines will be expired on tomorrow and add into task")
		s.dao.InsertDeleteMachinesTasks(machines)
	}
}

// 定时发邮件通知将要过期机器
func (s *Service) taskSendTaskMailMachinesWillExpired() {
	var (
		machines []*model.Machine
		err      error
	)

	if machines, err = s.dao.FindExpiredMachineByDay(s.c.Scheduler.ExpiredDate); err != nil {
		log.Error("Task send task mail machines will expired (%v)", err)
		return
	}
	for _, machine := range machines {
		log.Info("Machine named [%s] will be expired on next week and send a mail", machine.Name)
		if machine.DelayStatus == model.DelayStatusInit {
			s.dao.UpdateMachineDelayStatus(machine.ID, model.DelayStatusAuto)
		}
		if err = s.SendMail(model.MailTypeMachineWillExpired, machine); err != nil {
			log.Error("Send mail failed (%v)", err)
		}
	}
}

// 定时删除过期机器
func (s *Service) taskDeleteExpiredMachines() {
	var (
		tasks    []*model.Task
		instance *model.ReleaseInstance
		err      error
		machine  *model.Machine
	)

	if tasks, err = s.dao.FindDeleteMachineTasks(); err != nil {
		log.Error("Task delete expired machines error (%v)", err)
		return
	}
	for _, taskEle := range tasks {
		if instance, err = s.DelMachine(context.TODO(), taskEle.MachineID, "机器删除"); err != nil {
			log.Error("Task delete expired machines error (%v)", err)
			continue
		}
		if instance.InstanceReleaseStatus != model.SuccessDeletePaasMachines {
			if machine, err = s.dao.QueryMachine(taskEle.MachineID); err != nil {
				log.Error("Task delete expired machines error (%v)", err)
				continue
			}
			if err = s.SendMail(model.MailTypeTaskDeleteMachineFailed, machine); err != nil {
				log.Error("Send mail failed (%v)", err)
			}
			s.dao.UpdateTaskStatusByTaskID(taskEle.ID, model.TaskFailed)
		}
	}
}

func (s *Service) taskMachineStatus() {
	var (
		pathAndPodNames  map[string][]string
		machineStatuses  map[string]bool
		createdPodNames  []string
		creatingPodNames []string
		creatingMachines []*model.Machine
		err              error
		c                = context.TODO()
	)
	if pathAndPodNames, err = s.dao.QueryPathAndPodNamesMapping(); err != nil {
		log.Error("Query creating machines in db err(%v)", err)
		return
	}
	if len(pathAndPodNames) == 0 {
		return
	}
	log.Info("Get pathAndPodNames(%v) from Service Tree", pathAndPodNames)
	if machineStatuses, err = s.TreeMachinesIsExist(c, pathAndPodNames); err != nil {
		log.Error("Query service tree machine status err(%v)", err)
		return
	}
	for k, v := range machineStatuses {
		if v {
			createdPodNames = append(createdPodNames, k)
		} else {
			creatingPodNames = append(creatingPodNames, k)
		}
	}
	if len(createdPodNames) > 0 {
		if s.dao.UpdateMachineStatusByPodNames(createdPodNames, model.BootMachineInMerlin); err != nil {
			log.Error("update creating machines to boot in db err(%v)", err)
			return
		}
	}
	if len(creatingPodNames) <= 0 {
		return
	}
	if creatingMachines, err = s.dao.QueryMachinesByPodNames(creatingPodNames); err != nil {
		log.Error("Query creating machines in db err(%v)", err)
		return
	}
	for _, m := range creatingMachines {
		log.Info("Create machine(%v) deadline exceeded", m)
		if _, err = s.verifyPassStatus(c, m, false); err != nil {
			log.Error("Del verify machine(%v) in db err(%v)", m, err)
			continue
		}
		if _, err = s.dao.DelPaasMachine(c, m.ToPaasQueryAndDelMachineRequest()); err != nil {
			log.Error("Del creating machine(%v) in db err(%v)", m, err)
			continue
		}
	}
}

// 同步ios移动设备状态
func (s *Service) taskSyncMobileDeviceList() {
	s.SyncMobileDeviceList(context.Background())
}

// 定时清理由于回调失败，而在进行中的快照
func (s *Service) taskUpdateSnapshotStatusInDoing() {
	var (
		err             error
		snapshotRecords []*model.SnapshotRecord
	)
	if snapshotRecords, err = s.dao.FindSnapshotStatusInDoingOver2Hours(); err != nil {
		return
	}

	for _, snapshotRecord := range snapshotRecords {
		if err = s.dao.UpdateSnapshotRecordStatus(snapshotRecord.MachineID, model.SnapShotFailed); err != nil {
			continue
		}
	}
}
