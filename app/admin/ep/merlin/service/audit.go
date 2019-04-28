package service

import (
	"context"
	"time"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
)

// DelayMachineEndTime Delay Machine End Time.
func (s *Service) DelayMachineEndTime(c context.Context, machineID int64, username string) (status int, err error) {
	var (
		machine    *model.Machine
		machineLog = &model.MachineLog{}
	)
	status = -1
	//查询机器延期状态
	if machine, err = s.dao.QueryMachine(machineID); err != nil {
		return
	}

	//机器必须存在，机器状态必须是有效的，延期状态必须为可自动延期，结束时间必须是前后7天之内
	if machine == nil || machine.DelayStatus != model.DelayStatusAuto || machine.Status < 0 || machine.EndTime.Before(time.Now().AddDate(0, 0, -8)) || machine.EndTime.After(time.Now().AddDate(0, 0, 8)) {
		err = ecode.MerlinDelayMachineErr
		return
	}

	//执行延期操作
	if err = s.dao.UpdateMachineEndTime(machineID, model.DelayStatusApply, machine.EndTime.AddDate(0, 0, 30)); err != nil {
		machineLog.OperateResult = model.OperationFailedForMachineLog
	} else {
		machineLog.OperateResult = model.OperationSuccessForMachineLog
		status = 0
	}

	//记录操作日志
	machineLog.OperateType = model.DelayMachineEndTime
	machineLog.Username = username
	machineLog.MachineID = machineID
	err = s.dao.InsertMachineLog(machineLog)
	return

}

// ApplyDelayMachineEndTime Apply Delay Machine End Time.
func (s *Service) ApplyDelayMachineEndTime(c context.Context, username string, applyEndTime *model.ApplyEndTimeRequest) (status int, err error) {
	var (
		machine    *model.Machine
		machineLog = &model.MachineLog{}
	)

	status = -1
	//查询机器延期状态
	if machine, err = s.dao.QueryMachine(applyEndTime.MachineID); err != nil {
		return
	}

	//申请延期时，机器必须手动申请延期过，且机器状态为正常
	if machine == nil || machine.DelayStatus != model.DelayStatusApply || machine.Status < 0 {
		err = ecode.MerlinApplyMachineErr
		return
	}
	machineLog.OperateResult = model.OperationFailedForMachineLog

	var applyET time.Time
	if applyET, err = time.ParseInLocation(model.TimeFormat, applyEndTime.ApplyEndTime, time.Local); err != nil {
		return
	}

	//延长日期超过过期时间三个月或小于当前过期时间，驳回
	if machine.EndTime.AddDate(0, 3, 1).Before(applyET) || machine.EndTime.After(applyET) {
		err = ecode.MerlinApplyMachineByApplyEndTimeMore3MErr
		return
	}

	ar := &model.ApplicationRecord{
		Applicant:    username,
		MachineID:    machine.ID,
		ApplyEndTime: applyET,
		Status:       model.ApplyDelayInit,
		Auditor:      applyEndTime.Auditor,
	}

	//添加延期申请记录
	if err = s.dao.InsertApplicationRecordAndUpdateMachineDelayStatus(ar, machine.ID, model.DelayStatusDisable); err == nil {
		machineLog.OperateResult = model.OperationSuccessForMachineLog
		status = 0
		//发送邮件通知
		s.SendMailApplyDelayMachineEndTime(c, applyEndTime.Auditor, username, machine.ID, machine.EndTime, applyET)

	}
	//记录操作日志
	machineLog.OperateType = model.DelayMachineEndTime
	machineLog.Username = username
	machineLog.MachineID = machine.ID
	if err = s.dao.InsertMachineLog(machineLog); err != nil {
		return
	}

	return
}

// CancelMachineEndTime Cancel Machine End Time.
func (s *Service) CancelMachineEndTime(c context.Context, auditID int64, username string) (status int, err error) {
	var (
		applicationRecord *model.ApplicationRecord
		machineLog        = &model.MachineLog{}
		machine           *model.Machine
	)
	status = -1

	if applicationRecord, err = s.dao.FindApplicationRecordsByID(auditID); err != nil {
		return
	}

	//查询机器延期状态
	if machine, err = s.dao.QueryMachine(applicationRecord.MachineID); err != nil {
		return
	}

	//取消延期时，审批状态应为 ApplyDelayInit，且用户名是当时申请延期名,机器状态为正常
	if applicationRecord == nil || applicationRecord.Status != model.ApplyDelayInit || applicationRecord.Applicant != username || machine.Status < 0 {
		err = ecode.MerlinCancelMachineErr
		return
	}
	machineLog.OperateResult = model.OperationFailedForMachineLog

	if err = s.dao.UpdateAuditStatusAndUpdateMachineDelayStatus(applicationRecord.MachineID, auditID, model.DelayStatusApply, model.ApplyDelayCancel); err == nil {
		machineLog.OperateResult = model.OperationSuccessForMachineLog
		status = 0
	}
	//记录操作日志
	machineLog.OperateType = model.CancelDelayMachineEndTime
	machineLog.Username = username
	machineLog.MachineID = applicationRecord.MachineID
	if err = s.dao.InsertMachineLog(machineLog); err != nil {
		return
	}

	return

}

// AuditMachineEndTime Audit Machine End Time.
func (s *Service) AuditMachineEndTime(c context.Context, auditID int64, username string, auditResult bool, comment string) (status int, err error) {
	var (
		applicationRecord *model.ApplicationRecord
		machineLog        = &model.MachineLog{}
		machine           *model.Machine
	)
	status = -1
	if applicationRecord, err = s.dao.FindApplicationRecordsByID(auditID); err != nil {
		return
	}

	//查询机器延期状态
	if machine, err = s.dao.QueryMachine(applicationRecord.MachineID); err != nil {
		return
	}

	//审核时，延期状态必须为ApplyDelayInit，同时用户名为审核者名一致,机器状态为正常
	if applicationRecord == nil || applicationRecord.Status != model.ApplyDelayInit || applicationRecord.Auditor != username || machine.Status < 0 {
		err = ecode.MerlinAuditMachineErr
		return
	}
	machineLog.OperateResult = model.OperationFailedForMachineLog

	if auditResult {
		//批准  //修改机器状态和延期时间，审批状态
		err = s.dao.UpdateAuditStatusAndUpdateMachineEndTime(applicationRecord.MachineID, auditID, model.DelayStatusApply, model.ApplyDelayApprove, applicationRecord.ApplyEndTime, comment)
	} else {
		//驳回  //修改机器状态,审批状态
		err = s.dao.UpdateAuditStatusAndUpdateMachineDelayStatusComment(applicationRecord.MachineID, auditID, model.DelayStatusApply, model.ApplyDelayDecline, comment)
	}

	//发送邮件通知
	if err == nil {
		machineLog.OperateResult = model.OperationSuccessForMachineLog
		status = 0
		s.SendMailAuditResult(c, applicationRecord.Auditor, applicationRecord.Applicant, applicationRecord.MachineID, auditResult)
	}
	//记录操作日志
	machineLog.OperateType = model.AuditDelayMachineEndTime
	machineLog.Username = username
	machineLog.MachineID = applicationRecord.MachineID
	if err = s.dao.InsertMachineLog(machineLog); err != nil {
		return
	}

	return

}

// GetApplicationRecordsByApplicant Get Application Records By Applicant.
func (s *Service) GetApplicationRecordsByApplicant(c context.Context, username string, pn, ps int) (p *model.PaginateApplicationRecord, err error) {
	var (
		total              int64
		applicationRecords []*model.ApplicationRecord
	)
	if total, applicationRecords, err = s.dao.FindApplicationRecordsByApplicant(username, pn, ps); err != nil {
		return
	}
	p = &model.PaginateApplicationRecord{
		PageNum:            pn,
		PageSize:           ps,
		Total:              total,
		ApplicationRecords: applicationRecords,
	}
	return
}

// GetApplicationRecordsByAuditor Get Application Records By Auditor.
func (s *Service) GetApplicationRecordsByAuditor(c context.Context, username string, pn, ps int) (p *model.PaginateApplicationRecord, err error) {
	var (
		total              int64
		applicationRecords []*model.ApplicationRecord
	)
	if total, applicationRecords, err = s.dao.FindApplicationRecordsByAuditor(username, pn, ps); err != nil {
		return
	}
	p = &model.PaginateApplicationRecord{
		PageNum:            pn,
		PageSize:           ps,
		Total:              total,
		ApplicationRecords: applicationRecords,
	}
	return
}

// GetApplicationRecordsByMachineID Get Application Records By MachineID.
func (s *Service) GetApplicationRecordsByMachineID(c context.Context, machineID int64, pn, ps int) (p *model.PaginateApplicationRecord, err error) {
	var (
		total              int64
		applicationRecords []*model.ApplicationRecord
	)
	if total, applicationRecords, err = s.dao.FindApplicationRecordsByMachineID(machineID, pn, ps); err != nil {
		return
	}
	p = &model.PaginateApplicationRecord{
		PageNum:            pn,
		PageSize:           ps,
		Total:              total,
		ApplicationRecords: applicationRecords,
	}
	return
}
