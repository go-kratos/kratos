package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/satori/go.uuid"
)

//SyncMobileDeviceList Sync Device Farm List .
func (s *Service) SyncMobileDeviceList(c context.Context) (uid string, err error) {
	uid = uuid.NewV4().String()
	s.deviceChan.Do(c, func(c context.Context) {
		var (
			startCnt  int
			status    int
			totalCnt  int
			addCnt    int
			updateCnt int
			delCnt    int
			err       error
		)

		if startCnt, err = s.dao.FindMobileSyncLogStartStatus(); err != nil {
			return
		}

		if startCnt > 0 {
			return
		}

		mobileSyncLog := &model.MobileSyncLog{
			UUID: uid,
		}
		s.dao.InsertMobileSyncLog(mobileSyncLog)

		if totalCnt, addCnt, updateCnt, delCnt, err = s.SyncMobileDeviceListInChan(uid); err != nil {
			status = -1
		} else {
			status = 1
		}

		mobileSyncLog = &model.MobileSyncLog{
			UUID:      uid,
			TotalCnt:  totalCnt,
			AddCnt:    addCnt,
			UpdateCnt: updateCnt,
			DeleteCnt: delCnt,
			Status:    status,
		}
		s.dao.UpdateMobileSyncLog(mobileSyncLog)

	})
	return
}

// SyncMobileDeviceListInChan  Sync Mobile Device List In Chan.
func (s *Service) SyncMobileDeviceListInChan(uuid string) (totalCnt, addCnt, updateCnt, delCnt int, err error) {
	var (
		resTotal map[string][]*model.Device
	)

	if resTotal, err = s.dao.MobileDeviceList(context.Background()); err != nil {
		log.Error("SyncMobileDeviceListInChan httpClient err(%+v)", err)
		err = nil
	}

	if len(resTotal) == 0 {
		return
	}

	for host, devices := range resTotal {
		totalCnt = totalCnt + len(devices)
		for _, device := range devices {
			var deviceType int
			var tmpMobileMachine *model.MobileMachine

			if device.IsSimulator {
				deviceType = model.MBSimulator
			}

			if tmpMobileMachine, err = s.dao.FindMobileMachineBySerial(device.Serial); err != nil {
				log.Error("d.SyncMobileDeviceList  FindMobileMachineBySerial err(%+v)", err)
				continue
			}
			mobileMachine := &model.MobileMachine{
				Serial:  device.Serial,
				Name:    device.Name,
				CPU:     device.CPU,
				Version: device.Version,
				Mode:    device.Mode,
				State:   device.State,
				Host:    host,
				Type:    deviceType,
				Action:  model.MBOnline,
				UUID:    uuid,
			}

			if tmpMobileMachine.ID == 0 {
				if err = s.dao.InsertMobileMachine(mobileMachine); err != nil {
					log.Error("d.SyncMobileDeviceList InsertMobileMachine err(%+v)", err)
				} else {
					addCnt = addCnt + 1
				}

			} else {
				if err = s.dao.UpdateMobileMachine(mobileMachine); err != nil {
					log.Error("d.SyncMobileDeviceList UpdateMobileMachine err(%+v)", err)
				} else {
					updateCnt = updateCnt + 1
				}
			}
		}
	}

	s.dao.DeleteMobileMachineNotInHost(s.c.DeviceFarm.HostList)

	delCnt, err = s.dao.DeleteMobileMachineByUUID(uuid)

	return
}

// QueryMobileDeviceList Query Mobile Device List.
func (s *Service) QueryMobileDeviceList(c context.Context, queryRequest *model.QueryMobileDeviceRequest) (p *model.PaginateMobileMachines, err error) {
	var (
		total                  int64
		mobileMachines         []*model.MobileMachine
		mobileMachinesResponse []*model.MobileMachineResponse
	)
	if total, mobileMachines, err = s.dao.FindMobileMachines(queryRequest); err != nil {
		return
	}

	for _, mobileMachine := range mobileMachines {
		var mobileMachineResponse model.MobileMachineResponse
		mobileMachineResponse.MobileMachine = mobileMachine
		mobileMachinesResponse = append(mobileMachinesResponse, &mobileMachineResponse)
	}

	if err = s.getMobileImages(mobileMachinesResponse); err != nil {
		return
	}

	p = &model.PaginateMobileMachines{
		PageNum:        queryRequest.PageNum,
		PageSize:       queryRequest.PageSize,
		Total:          total,
		MobileMachines: mobileMachinesResponse,
	}
	return
}

// RefreshMobileDeviceDetail Refresh Mobile Device Detail.
func (s *Service) RefreshMobileDeviceDetail(c context.Context, deviceID int64) (mobileMachineResponse *model.MobileMachineResponse, err error) {
	var (
		device        *model.Device
		mobileMachine *model.MobileMachine
	)

	if mobileMachine, err = s.dao.FindMobileMachineByID(deviceID); err != nil {
		return
	}

	if device, err = s.dao.MobileDeviceDetail(c, mobileMachine.Host, mobileMachine.Serial); err != nil {
		return
	}

	//update db
	mobileMachine.State = device.State
	mobileMachine.Action = 1
	if err = s.dao.UpdateMobileMachine(mobileMachine); err != nil {
		log.Error("d.RefreshMobileDeviceDetail err(%+v)", err)
		return
	}

	mobileMachineResponse = &model.MobileMachineResponse{}
	mobileMachineResponse.MobileMachine = mobileMachine
	if err = s.getMobileImage(mobileMachineResponse); err != nil {
		return
	}

	return
}

// MobileMachineCategory Mobile Machine Category.
func (s *Service) MobileMachineCategory(c context.Context, isShowOffline bool) (categoryList []*model.MobileDeviceCategoryResponse, err error) {
	var mobileCategory *model.MobileCategory
	if mobileCategory, err = s.dao.FindMobileMachineCategory(isShowOffline); err != nil {
		return
	}

	interfaceCPU := make([]interface{}, len(mobileCategory.CPUs))
	for i, v := range mobileCategory.CPUs {
		interfaceCPU[i] = v
	}
	categoryCPU := &model.MobileDeviceCategoryResponse{
		Name:   "cpu",
		Label:  "CPU",
		Values: interfaceCPU,
	}

	interfaceVersion := make([]interface{}, len(mobileCategory.Versions))
	for i, v := range mobileCategory.Versions {
		interfaceVersion[i] = v
	}
	categoryVersion := &model.MobileDeviceCategoryResponse{
		Name:   "version",
		Label:  "版本",
		Values: interfaceVersion,
	}

	interfaceMode := make([]interface{}, len(mobileCategory.Modes))
	for i, v := range mobileCategory.Modes {
		interfaceMode[i] = v
	}
	categoryMode := &model.MobileDeviceCategoryResponse{
		Name:   "mode",
		Label:  "型号",
		Values: interfaceMode,
	}

	interfaceType := make([]interface{}, len(mobileCategory.Types))
	for i, v := range mobileCategory.Types {
		interfaceType[i] = v
	}
	categoryType := &model.MobileDeviceCategoryResponse{
		Name:   "type",
		Label:  "是否真机",
		Values: interfaceType,
	}

	interfaceState := make([]interface{}, len(mobileCategory.States))
	for i, v := range mobileCategory.States {
		interfaceState[i] = v
	}
	categoryState := &model.MobileDeviceCategoryResponse{
		Name:   "state",
		Label:  "状态",
		Values: interfaceState,
	}

	interfaceUsage := make([]interface{}, len(mobileCategory.Usages))
	for i, v := range mobileCategory.Usages {
		interfaceUsage[i] = v
	}
	categoryUsage := &model.MobileDeviceCategoryResponse{
		Name:   "usage",
		Label:  "使用情况",
		Values: interfaceUsage,
	}

	categoryList = append(categoryList, categoryCPU)
	categoryList = append(categoryList, categoryVersion)
	categoryList = append(categoryList, categoryMode)
	categoryList = append(categoryList, categoryType)
	categoryList = append(categoryList, categoryState)
	categoryList = append(categoryList, categoryUsage)
	return
}

// BindMobileDevice Bind Mobile Device.
func (s *Service) BindMobileDevice(c context.Context, deviceID int64, username string) (status int, err error) {
	var mobileMachine *model.MobileMachine

	status = -1
	if mobileMachine, err = s.dao.FindMobileMachineByID(deviceID); err != nil {
		return
	}

	if mobileMachine.Action < model.MBOnline {
		err = ecode.MerlinDeviceFarmMachineStatusErr
		return
	}

	//真机借出状态， 无法操作
	if mobileMachine.Type == model.MBReal && mobileMachine.IsLendOut == model.MBLendOut {
		err = ecode.MerlinDeviceIsLendOut
		return
	}

	if mobileMachine.Username != "" {
		return
	}

	mobileMachine.Username = username
	mobileMachine.EndTime = time.Now().AddDate(0, s.c.DeviceFarm.ApplyMonthTime, 0)
	mobileMachine.LastBindTime = time.Now()
	if err = s.dao.UpdateMobileMachine(mobileMachine); err != nil {
		return
	}

	if err = s.addMobileDeviceLog(deviceID, username, model.MBBindLog); err != nil {
		return
	}

	status = 0
	return

}

// ReleaseMobileDevice Release Mobile Device.
func (s *Service) ReleaseMobileDevice(c context.Context, deviceID int64, username string) (status int, err error) {
	var mobileMachine *model.MobileMachine

	status = -1
	if mobileMachine, err = s.isAccessToMobileDevice(deviceID, username); err != nil {
		return
	}

	mobileMachine.Username = ""
	if mobileMachine.EndTime, err = time.ParseInLocation(model.TimeFormat, "1990-01-01 01:01:01", time.Local); err != nil {
		return
	}

	if err = s.dao.UpdateMobileMachineByRelease(mobileMachine); err != nil {
		return
	}

	if err = s.addMobileDeviceLog(deviceID, username, model.MBReleaseLog); err != nil {
		return
	}

	status = 0
	return

}

// LendOutMobileDevice Pull Out Mobile Device.
func (s *Service) LendOutMobileDevice(c context.Context, deviceID int64, username string) (status int, err error) {
	var (
		mobileMachine *model.MobileMachine
	)

	status = -1

	if mobileMachine, err = s.isAccessToMobileDevice(deviceID, username); err != nil {
		return
	}

	//只有真机支持借出 真机type=0
	if mobileMachine.Type != model.MBReal {
		err = ecode.MerlinDeviceIsNotRealMachineErr
		return
	}

	if err = s.dao.LendOutMobileMachine(deviceID, mobileMachine.Serial, username); err != nil {
		return
	}
	status = 0
	return
}

// ReturnMobileDevice Return Mobile Device.
func (s *Service) ReturnMobileDevice(c context.Context, deviceID int64, username string) (status int, err error) {
	var (
		mobileMachine *model.MobileMachine
		isSuperOwner  = false
	)

	status = -1

	for _, superOwner := range s.c.DeviceFarm.SuperOwner {
		if username == superOwner {
			isSuperOwner = true
			break
		}
	}

	if !isSuperOwner {
		err = ecode.MerlinDeviceNoRight
		return
	}

	if mobileMachine, err = s.dao.FindMobileMachineByID(deviceID); err != nil {
		return
	}

	if mobileMachine.IsLendOut == model.MBOnSite {
		err = ecode.MerlinDeviceFarmErr
		return
	}

	if err = s.dao.ReturnMobileMachine(deviceID, mobileMachine.Serial, username); err != nil {
		return
	}

	status = 0
	return
}

// StartMobileDevice Start Mobile Device.
func (s *Service) StartMobileDevice(c context.Context, deviceID int64, username string) (deviceBootData *model.DeviceBootData, err error) {
	var (
		mobileMachine *model.MobileMachine
		device        *model.Device
	)

	if mobileMachine, err = s.isAccessToMobileDevice(deviceID, username); err != nil {
		return
	}

	if deviceBootData, err = s.dao.BootMobileDevice(c, mobileMachine.Host, mobileMachine.Serial); err != nil {
		return
	}

	if device, err = s.dao.MobileDeviceDetail(c, mobileMachine.Host, mobileMachine.Serial); err != nil {
		return
	}

	mobileMachine.State = device.State
	mobileMachine.WsURL = deviceBootData.WSRUL
	mobileMachine.UploadURL = deviceBootData.UploadURL
	if err = s.dao.UpdateMobileMachineWsurlAndState(mobileMachine); err != nil {
		return
	}

	err = s.addMobileDeviceLog(deviceID, username, model.MBStartLog)
	return

}

// ShutdownMobileDevice Shutdown Mobile Device.
func (s *Service) ShutdownMobileDevice(c context.Context, deviceID int64, username string) (status int, err error) {
	var (
		mobileMachine *model.MobileMachine
		device        *model.Device
	)
	status = -1
	if mobileMachine, err = s.isAccessToMobileDevice(deviceID, username); err != nil {
		return
	}

	if err = s.dao.ShutdownMobileDevice(c, mobileMachine.Host, mobileMachine.Serial); err != nil {
		return
	}

	if device, err = s.dao.MobileDeviceDetail(c, mobileMachine.Host, mobileMachine.Serial); err != nil {
		return
	}

	mobileMachine.State = device.State
	mobileMachine.WsURL = ""
	mobileMachine.UploadURL = ""
	if err = s.dao.UpdateMobileMachineWsurlAndState(mobileMachine); err != nil {
		return
	}

	if err = s.addMobileDeviceLog(deviceID, username, model.MBShutDownLog); err != nil {
		return
	}
	status = 0
	return
}

func (s *Service) isAccessToMobileDevice(deviceID int64, username string) (mobileMachine *model.MobileMachine, err error) {
	if mobileMachine, err = s.dao.FindMobileMachineByID(deviceID); err != nil {
		return
	}

	if mobileMachine.Action < model.MBOnline {
		err = ecode.MerlinDeviceFarmMachineStatusErr
		return
	}

	//真机借出状态， 无法操作
	if mobileMachine.Type == model.MBReal && mobileMachine.IsLendOut == model.MBLendOut {
		err = ecode.MerlinDeviceIsLendOut
		return
	}

	for _, superOwner := range s.c.DeviceFarm.SuperOwner {
		if username == superOwner {
			return
		}
	}

	if mobileMachine.Username != username {
		err = ecode.MerlinDeviceNotBind
	}
	return
}

func (s *Service) addMobileDeviceLog(deviceID int64, username, operateType string) (err error) {
	machineLog := &model.MobileMachineLog{
		OperateType:   operateType,
		Username:      username,
		MachineID:     deviceID,
		OperateResult: model.OperationSuccessForMachineLog,
	}
	return s.dao.InsertMobileMachineLog(machineLog)
}

func (s *Service) getMobileImages(mobileMachinesResponse []*model.MobileMachineResponse) (err error) {
	var mobileImages []*model.MobileImage
	if mobileImages, err = s.dao.FindAllMobileImages(); err != nil {
		return
	}
	for _, mobileMachineResponse := range mobileMachinesResponse {
		for _, mobileImage := range mobileImages {
			if mobileImage.Mode == mobileMachineResponse.Mode {
				mobileMachineResponse.ImageSrc = mobileImage.ImageSrc
				continue
			}
		}
	}
	return
}

func (s *Service) getMobileImage(mobileMachinesResponse *model.MobileMachineResponse) (err error) {
	var mobileImage *model.MobileImage
	if mobileImage, err = s.dao.FindMobileImageByMode(mobileMachinesResponse.Mode); err != nil {
		return
	}
	mobileMachinesResponse.ImageSrc = mobileImage.ImageSrc
	return
}

// QueryMobileMachineLogs Query Mobile Machine Logs.
func (s *Service) QueryMobileMachineLogs(c context.Context, queryRequest *model.QueryMobileMachineLogRequest) (p *model.PaginateMobileMachineLog, err error) {
	var (
		total       int64
		machineLogs []*model.AboundMobileMachineLog
	)
	if total, machineLogs, err = s.dao.FindMobileMachineLogs(queryRequest); err != nil {
		return
	}
	p = &model.PaginateMobileMachineLog{
		PageNum:     queryRequest.PageNum,
		PageSize:    queryRequest.PageSize,
		Total:       total,
		MachineLogs: machineLogs,
	}
	return
}

// QueryMobileMachineLendOutLogs Query Mobile Machine Lend Out Logs.
func (s *Service) QueryMobileMachineLendOutLogs(c context.Context, username string, queryRequest *model.QueryMobileMachineLogRequest) (p *model.PaginateMobileMachineLendOutLog, err error) {
	var (
		total                 int64
		machineLogs           []*model.AboundMobileMachineLog
		machineLendOutRecords []*model.MachineLendOutRecord
		mobileMachine         *model.MobileMachine
		existMachineIDs       []int64
		isSuperUser           bool
	)
	isSuperUser = s.isSuperUser(username)

	queryRequest.OperateType = model.MBLendOutLog

	if total, machineLogs, err = s.dao.FindMobileMachineLogs(queryRequest); err != nil {
		return
	}

	for _, machineLog := range machineLogs {

		var (
			enableReturn      bool
			isExistLendRecord bool
		)

		if mobileMachine, err = s.dao.FindMobileMachineByID(machineLog.MachineID); err != nil {
			return
		}

		for _, existMachineID := range existMachineIDs {
			if existMachineID == machineLog.MachineID {
				isExistLendRecord = true
				break
			}
		}

		//移动设备借出状态 且 用户是超级用户 才有资格归还手机 只有最新一条记录 才能点击返还
		if isExistLendRecord {

			machineLendOutRecord := &model.MachineLendOutRecord{
				MachineID:    machineLog.MachineID,
				Lender:       machineLog.Username,
				LendTime:     machineLog.OperateTime,
				EnableReturn: false,
				Status:       model.MBOnSite,
			}
			machineLendOutRecords = append(machineLendOutRecords, machineLendOutRecord)

		} else {

			if isSuperUser && (mobileMachine.IsLendOut == model.MBLendOut) {
				enableReturn = true
			}

			machineLendOutRecord := &model.MachineLendOutRecord{
				MachineID:    machineLog.MachineID,
				Lender:       machineLog.Username,
				LendTime:     machineLog.OperateTime,
				EnableReturn: enableReturn,
				Status:       mobileMachine.IsLendOut,
			}
			machineLendOutRecords = append(machineLendOutRecords, machineLendOutRecord)
			existMachineIDs = append(existMachineIDs, machineLog.MachineID)
		}

	}

	p = &model.PaginateMobileMachineLendOutLog{
		PageNum:  queryRequest.PageNum,
		PageSize: queryRequest.PageSize,
		Total:    total,
		MachineLendOutRecords: machineLendOutRecords,
	}
	return
}

// QueryMobileMachineErrorLogs Query Mobile Machine Error Logs.
func (s *Service) QueryMobileMachineErrorLogs(c context.Context, queryRequest *model.QueryMobileMachineErrorLogRequest) (p *model.PaginateMobileMachineErrorLog, err error) {
	var (
		total       int64
		machineLogs []*model.MobileMachineErrorLog
	)
	if total, machineLogs, err = s.dao.FindMobileMachineErrorLog(queryRequest); err != nil {
		return
	}
	p = &model.PaginateMobileMachineErrorLog{
		PageNum:     queryRequest.PageNum,
		PageSize:    queryRequest.PageSize,
		Total:       total,
		MachineLogs: machineLogs,
	}
	return
}

// CallBackMobileDeviceError Call Back Mobile Device Error.
func (s *Service) CallBackMobileDeviceError(c context.Context, serialName, errMsg string, errCode int) (err error) {
	var mobileMachine *model.MobileMachine

	if mobileMachine, err = s.dao.FindMobileMachineBySerial(serialName); err != nil {
		return
	}

	mobileMachineErrorLog := &model.MobileMachineErrorLog{
		MachineID:    mobileMachine.ID,
		SerialName:   serialName,
		ErrorMessage: errMsg,
		ErrorCode:    errCode,
	}

	err = s.dao.InsertMobileMachineErrorLog(mobileMachineErrorLog)
	return
}

// ReportMobileDeviceError Report Mobile Device Error.
func (s *Service) ReportMobileDeviceError(c context.Context, serialName, errMsg string, errCode int) (err error) {
	var mobileMachine *model.MobileMachine

	if mobileMachine, err = s.dao.FindMobileMachineBySerial(serialName); err != nil {
		return
	}

	if mobileMachine.ID == 0 {
		err = ecode.NothingFound
		return
	}

	content := fmt.Sprintf("通知\n设备ID: %d\n系列号: %s\n型号: %s\n系统版本: %s\n使用者: %s\n错误信息: %s\n错误码: %d", mobileMachine.ID, serialName, mobileMachine.Mode, mobileMachine.Version, mobileMachine.Username, errMsg, errCode)
	msgSendReq := &model.MsgSendReq{
		Chatid:  s.c.WeChat.WeChatDeviceFarm.ChatID,
		Msgtype: s.c.WeChat.WeChatDeviceFarm.MsgType,
		Text:    model.MsgSendReqText{Content: content},
		Safe:    s.c.WeChat.WeChatDeviceFarm.Safe,
	}
	if _, err = s.dao.WeChatSendMessage(c, msgSendReq); err != nil {
		return
	}

	mobileMachineErrorLog := &model.MobileMachineErrorLog{
		MachineID:    mobileMachine.ID,
		SerialName:   serialName,
		ErrorMessage: errMsg,
		ErrorCode:    errCode,
	}

	err = s.dao.InsertMobileMachineErrorLog(mobileMachineErrorLog)
	return
}

// DeviceFarmSuperUser Device Farm Super User.
func (s *Service) DeviceFarmSuperUser() (responseMap map[string]interface{}) {
	responseMap = make(map[string]interface{})
	responseMap["superuser_list"] = s.c.DeviceFarm.SuperOwner
	return
}

// IsBindMobileDeviceByTheUser Is Bind Mobile Device By The User.
func (s *Service) IsBindMobileDeviceByTheUser(c context.Context, deviceID int64, username string) (isBind bool, err error) {
	var mobileMachine *model.MobileMachine

	if mobileMachine, err = s.dao.FindMobileMachineByID(deviceID); err != nil {
		return
	}

	if mobileMachine.Username == username {
		isBind = true
	}
	return
}

func (s *Service) isSuperUser(username string) bool {
	for _, superOwner := range s.c.DeviceFarm.SuperOwner {
		if username == superOwner {
			return true
		}
	}
	return false
}
