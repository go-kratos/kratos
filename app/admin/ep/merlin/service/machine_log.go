package service

import "go-common/app/admin/ep/merlin/model"

// AddMachineLog add machine log.
func (s *Service) AddMachineLog(username string, machineID int64, operateType string, operateResult string) (err error) {
	machineLog := &model.MachineLog{}
	machineLog.OperateType = operateType
	machineLog.Username = username
	machineLog.MachineID = machineID
	machineLog.OperateResult = operateResult

	err = s.dao.InsertMachineLog(machineLog)
	return
}
