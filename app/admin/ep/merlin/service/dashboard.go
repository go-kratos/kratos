package service

import (
	"context"
	"math"
	"time"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/log"
)

// QueryMachineLifeCycle Query Machine Life Cycle.
func (s *Service) QueryMachineLifeCycle(c context.Context) (machineResponse map[string]interface{}, err error) {
	var machineLifeCycles []*model.MachineLifeCycle

	machineResponse = make(map[string]interface{})

	if machineLifeCycles, err = s.dao.MachineLifeCycle(); err != nil {
		return
	}
	machineResponse["machine_life_cycle"] = machineLifeCycles

	return
}

// QueryMachineCount Query Machine Count.
func (s *Service) QueryMachineCount(c context.Context) (machineResponse *model.MachineCountGroupResponse, err error) {
	var (
		machinesCount          []*model.MachineCountGroupByBusiness
		machinesCountInRunning []*model.MachineCountGroupByBusiness
		businessUnits          []string
		machineCntData         []int
		machineCntInRunData    []int
	)

	if machinesCount, err = s.dao.MachineCountGroupByBusiness(); err != nil {
		return
	}

	if machinesCountInRunning, err = s.dao.MachineCountGroupByBusinessInRunning(); err != nil {
		return
	}

	for _, machineCount := range machinesCount {
		var runningCount int
		businessUnits = append(businessUnits, machineCount.BusinessUnit)
		machineCntData = append(machineCntData, machineCount.Count)
		for _, machineCountInRunning := range machinesCountInRunning {
			if machineCountInRunning.BusinessUnit == machineCount.BusinessUnit {
				runningCount = machineCountInRunning.Count
				break
			}
		}
		machineCntInRunData = append(machineCntInRunData, runningCount)
	}

	machineCountItem := &model.MachineCountGroupItem{
		Type: "已创建机器数量",
		Data: machineCntData,
	}

	machineCountInRunningItem := &model.MachineCountGroupItem{
		Type: "正在运行机器数量",
		Data: machineCntInRunData,
	}

	machineResponse = &model.MachineCountGroupResponse{
		BusinessUnits: businessUnits,
		Items:         []*model.MachineCountGroupItem{machineCountItem, machineCountInRunningItem},
	}

	return
}

// QueryMachineCreatedAndEndTime Query Machine Created And End Time.
func (s *Service) QueryMachineCreatedAndEndTime(c context.Context) (machineResponse map[string]interface{}, err error) {
	var (
		machineCreatedTime []*model.MachineCreatedAndEndTime
		machineExpiredTime []*model.MachineCreatedAndEndTime
	)

	machineResponse = make(map[string]interface{})

	if machineCreatedTime, err = s.dao.MachineLatestCreated(); err != nil {
		return
	}

	if machineExpiredTime, err = s.dao.MachineWillBeExpired(); err != nil {
		return
	}

	machineResponse["machine_created_top"] = machineCreatedTime
	machineResponse["machine_expired_top"] = machineExpiredTime

	return
}

// QueryMachineUsage Query Machine Usage.
func (s *Service) QueryMachineUsage(c context.Context) (machineResponse map[string]interface{}, err error) {
	var (
		machines              []*model.Machine
		pmds                  []*model.PaasMachineDetail
		pqadmrs               []*model.PaasQueryAndDelMachineRequest
		totalCPU              float32
		totalMemory           float32
		totalMachine          int
		totalMachineInRunnint int
	)
	machineResponse = make(map[string]interface{})

	if totalMachine, err = s.dao.QueryMachineCount(); err != nil {
		return
	}

	if machines, err = s.dao.QueryMachineInRunning(); err != nil {
		log.Error("query MachineInRunning err(%v)", err)
		return
	}

	totalMachineInRunnint = len(machines)

	for _, machine := range machines {
		pqadmrs = append(pqadmrs, machine.ToPaasQueryAndDelMachineRequest())
	}

	if pmds, err = s.dao.QueryMachineUsageSummaryFromCache(c, pqadmrs); err != nil {
		return
	}

	for _, pmd := range pmds {
		pmd.ConvertUnits()
		totalCPU = totalCPU + pmd.CPULimit/model.CPURatio
		totalMemory = totalMemory + pmd.MemoryLimit/model.MemoryRatio
	}

	machineResponse["total_machine"] = totalMachine
	machineResponse["total_machine_in_running"] = totalMachineInRunnint
	machineResponse["total_cpu_usage"] = totalCPU
	machineResponse["total_memory_usage"] = totalMemory

	return
}

// QueryMobileMachineUsageCount Query Mobile Machines Count.
func (s *Service) QueryMobileMachineUsageCount(c context.Context) (res map[string]interface{}, err error) {
	var (
		mobileMachinesUserUsageCount []*model.MobileMachineUserUsageCount
		mobileMachinesUserLendCount  []*model.MobileMachineUserUsageCount
		mobileMachinesUsageCount     []*model.MobileMachineUsageCount
		mobileMachinesLendCount      []*model.MobileMachineUsageCount
	)

	if mobileMachinesUserUsageCount, err = s.dao.MobileMachineUserUsageCount(); err != nil {
		return
	}

	if mobileMachinesUserLendCount, err = s.dao.MobileMachineUserLendCount(); err != nil {
		return
	}

	if mobileMachinesUsageCount, err = s.dao.MobileMachineUsageCount(); err != nil {
		return
	}

	if mobileMachinesLendCount, err = s.dao.MobileMachineLendCount(); err != nil {
		return
	}

	res = make(map[string]interface{})
	res["user_usage_count"] = mobileMachinesUserUsageCount
	res["user_lend_count"] = mobileMachinesUserLendCount
	res["mobile_machine_usage_count"] = mobileMachinesUsageCount
	res["mobile_machine_lend_count"] = mobileMachinesLendCount

	return
}

// QueryMobileMachineModeCount Query Mobile Machine Mode Count.
func (s *Service) QueryMobileMachineModeCount(c context.Context) (res map[string]interface{}, err error) {
	var mobileMachinesTypeCount []*model.MobileMachineTypeCount

	if mobileMachinesTypeCount, err = s.dao.MobileMachineModeCount(); err != nil {
		return
	}

	res = make(map[string]interface{})

	res["mobile_machine_mode_count"] = mobileMachinesTypeCount

	return
}

// QueryMobileMachineUsageTime Query Mobile Machine Usage Time.
func (s *Service) QueryMobileMachineUsageTime(c context.Context) (ret []*model.MobileMachineUsageTimeResponse, err error) {
	var (
		mobileMachineLogs    []*model.MobileMachineLog
		mobileMachineLogsMap = make(map[int64][]*model.MobileMachineLog)
	)

	if mobileMachineLogs, err = s.dao.MobileMachineUseRecord(); err != nil {
		return
	}

	//按机器id 分组
	for _, mobileMachineLog := range mobileMachineLogs {
		mobileMachineLogsMap[mobileMachineLog.MachineID] = append(mobileMachineLogsMap[mobileMachineLog.MachineID], mobileMachineLog)
	}

	//按机器计算 使用时长

	for machineID := range mobileMachineLogsMap {
		var (
			mobileMachine              *model.MobileMachine
			isStartTimeFound           bool
			isEndTimeFound             bool
			startTime                  time.Time
			endTime                    time.Time
			username                   string
			preMobileMachinesUsageTime []*model.MobileMachineUsageTime
			totalDuration              float64
		)

		if mobileMachine, err = s.dao.FindMobileMachineByID(machineID); err != nil {
			continue
		}

		for index, preMobileMachineLog := range mobileMachineLogsMap[machineID] {
			if preMobileMachineLog.OperateType == model.MBBindLog {
				startTime = preMobileMachineLog.OperateTime
				username = preMobileMachineLog.Username
				isStartTimeFound = true
			}

			if preMobileMachineLog.OperateType == model.MBReleaseLog && isStartTimeFound {
				endTime = preMobileMachineLog.OperateTime
				isEndTimeFound = true
			}

			//处理中间绑定和解绑 计算使用时间
			if isStartTimeFound && isEndTimeFound {
				duration := math.Trunc(endTime.Sub(startTime).Minutes()*1e2+0.5) * 1e-2

				MobileMachineUsageTime := &model.MobileMachineUsageTime{
					Username:  username,
					StartTime: startTime,
					EndTime:   endTime,
					Duration:  duration,
				}
				preMobileMachinesUsageTime = append(preMobileMachinesUsageTime, MobileMachineUsageTime)

				isStartTimeFound = false
				isEndTimeFound = false

				totalDuration = totalDuration + duration
			}

			// 最后次记录和绑定记录且机器在相应人名下 计算实时机器使用时长
			if index == (len(mobileMachineLogsMap[machineID])-1) && isStartTimeFound && !isEndTimeFound && mobileMachine.Username != "" {
				timeNow := time.Now()
				duration := math.Trunc(timeNow.Sub(startTime).Minutes()*1e2+0.5) * 1e-2

				MobileMachineUsageTime := &model.MobileMachineUsageTime{
					Username:  username,
					StartTime: timeNow,
					EndTime:   endTime,
					Duration:  duration,
				}
				preMobileMachinesUsageTime = append(preMobileMachinesUsageTime, MobileMachineUsageTime)

				isStartTimeFound = false
				isEndTimeFound = false

				totalDuration = totalDuration + duration
			}
		}

		mobileMachineUsageTimeResponse := &model.MobileMachineUsageTimeResponse{
			MobileMachineID:         machineID,
			MobileMachineName:       mobileMachine.Name,
			ModeName:                mobileMachine.Mode,
			TotalDuration:           totalDuration,
			MobileMachinesUsageTime: preMobileMachinesUsageTime,
		}

		ret = append(ret, mobileMachineUsageTimeResponse)
	}
	return
}
