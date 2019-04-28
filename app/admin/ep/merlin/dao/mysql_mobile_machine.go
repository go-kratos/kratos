package dao

import (
	"database/sql"
	"fmt"
	"strconv"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// FindMobileMachineBySerial Find Device FarmBy Serial.
func (d *Dao) FindMobileMachineBySerial(serial string) (mobileMachine *model.MobileMachine, err error) {
	mobileMachine = &model.MobileMachine{}
	if err = d.db.Where("serial = ?", serial).Find(mobileMachine).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// DeleteMobileMachineByUUID Delete Mobile Machine By UUID.
func (d *Dao) DeleteMobileMachineByUUID(UUID string) (delCnt int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.MobileMachine{}).Where("uuid <> ? and action <>?", UUID, model.MBHostDel).Count(&delCnt).Update("action", model.MBOffline).Error)
	return
}

// DeleteMobileMachineNotInHost Delete Mobile Machine Not In Host.
func (d *Dao) DeleteMobileMachineNotInHost(hostList []string) (delCnt int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.MobileMachine{}).Where("host not in (?)", hostList).Count(&delCnt).Update("action", model.MBHostDel).Error)
	return
}

// FindMobileMachineByID Find Mobile Machine By Id.
func (d *Dao) FindMobileMachineByID(deviceID int64) (mobileMachine *model.MobileMachine, err error) {
	mobileMachine = &model.MobileMachine{}
	err = pkgerr.WithStack(d.db.Where("id = ?", deviceID).Find(mobileMachine).Error)
	return
}

// InsertMobileMachine Insert Device Farm.
func (d *Dao) InsertMobileMachine(mobileMachine *model.MobileMachine) (err error) {
	return pkgerr.WithStack(d.db.Create(&mobileMachine).Error)
}

// UpdateMobileMachineWsurlAndState Update Mobile Machine Wsurl And State.
func (d *Dao) UpdateMobileMachineWsurlAndState(mobileMachine *model.MobileMachine) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.MobileMachine{}).Where("serial=?", mobileMachine.Serial).Updates(map[string]interface{}{"wsurl": mobileMachine.WsURL, "state": mobileMachine.State, "upload_url": mobileMachine.UploadURL}).Error)
}

// UpdateMobileMachine Update Mobile Machine.
func (d *Dao) UpdateMobileMachine(mobileMachine *model.MobileMachine) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.MobileMachine{}).Where("serial=?", mobileMachine.Serial).Updates(mobileMachine).Error)
}

// UpdateMobileMachineByRelease Update Mobile Machine By Release.
func (d *Dao) UpdateMobileMachineByRelease(mobileMachine *model.MobileMachine) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.MobileMachine{}).Where("serial=?", mobileMachine.Serial).Updates(map[string]interface{}{"username": mobileMachine.Username, "end_time": mobileMachine.EndTime}).Error)
}

//FindMobileMachines Find Mobile Machines.
func (d *Dao) FindMobileMachines(queryRequest *model.QueryMobileDeviceRequest) (total int64, mobileMachines []*model.MobileMachine, err error) {
	gDB := d.db.Model(&model.MobileMachine{}).Where("action> ? ", model.MBHostDel)

	if queryRequest.MobileID > 0 {
		gDB = gDB.Where("id=?", queryRequest.MobileID)
	}
	if queryRequest.Serial != "" {
		gDB = gDB.Where("serial=?", queryRequest.Serial)
	}
	if queryRequest.Name != "" {
		gDB = gDB.Where("name=?", queryRequest.Name)
	}
	if queryRequest.Username != "" {
		gDB = gDB.Where("username=?", queryRequest.Username)
	}
	if queryRequest.OwnerName != "" {
		gDB = gDB.Where("owner_name=?", queryRequest.OwnerName)
	}
	if queryRequest.CPU != "" {
		gDB = gDB.Where("cpu=?", queryRequest.CPU)
	}
	if queryRequest.Version != "" {
		gDB = gDB.Where("version=?", queryRequest.Version)
	}
	if queryRequest.Mode != "" {
		gDB = gDB.Where("mode=?", queryRequest.Mode)
	}
	if queryRequest.Type > -1 {
		gDB = gDB.Where("type=?", queryRequest.Type)
	}
	if queryRequest.State != "" {
		gDB = gDB.Where("state=?", queryRequest.State)
	}
	if queryRequest.Online {
		gDB = gDB.Where("action=?", model.MBOnline)
	}
	if queryRequest.Usage == model.MBFree {
		gDB = gDB.Where("username='' and action=?", model.MBOnline)
	}
	if queryRequest.Usage == model.MBInUse {
		gDB = gDB.Where("username<>'' and action=?", model.MBOnline)
	}
	if queryRequest.Usage == model.MBNoConnect {
		gDB = gDB.Where("action=?", model.MBOffline)
	}
	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	//先出action>0的，未连接的排在后面
	err = pkgerr.WithStack(gDB.Order("action desc,id desc").Offset((queryRequest.PageNum - 1) * queryRequest.PageSize).Limit(queryRequest.PageSize).Find(&mobileMachines).Error)

	return
}

// FindMobileMachineCategory Find Mobile Machine Category.
func (d *Dao) FindMobileMachineCategory(isShowOffline bool) (mobileCategory *model.MobileCategory, err error) {
	var (
		rows     *sql.Rows
		cpus     []string
		versions []string
		modes    []string
		states   []string
		types    []int
		usages   []int

		sqlGroupTpl string

		inUseCnt int
		freeCnt  int
		noConCnt int
	)

	if isShowOffline {
		sqlGroupTpl = "select %s from mobile_machines where action> " + strconv.Itoa(model.MBHostDel) + " group by %s"
	} else {
		sqlGroupTpl = "select %s from mobile_machines where action>-1 group by %s"
	}

	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}

	//cpu
	if rows, err = tx.Raw(fmt.Sprintf(sqlGroupTpl, "cpu", "cpu")).Rows(); err != nil {
		return
	}
	for rows.Next() {
		var cpu string
		if err = rows.Scan(&cpu); err != nil {
			return
		}
		cpus = append(cpus, cpu)
	}

	//version
	if rows, err = tx.Raw(fmt.Sprintf(sqlGroupTpl, "version", "version")).Rows(); err != nil {
		return
	}
	for rows.Next() {
		var version string
		if err = rows.Scan(&version); err != nil {
			return
		}
		versions = append(versions, version)
	}

	//mode
	if rows, err = tx.Raw(fmt.Sprintf(sqlGroupTpl, "mode", "mode")).Rows(); err != nil {
		return
	}
	for rows.Next() {
		var mode string
		if err = rows.Scan(&mode); err != nil {
			return
		}
		modes = append(modes, mode)
	}

	//state
	if rows, err = tx.Raw(fmt.Sprintf(sqlGroupTpl, "state", "state")).Rows(); err != nil {
		return
	}
	for rows.Next() {
		var state string
		if err = rows.Scan(&state); err != nil {
			return
		}
		states = append(states, state)
	}

	//type
	if rows, err = tx.Raw(fmt.Sprintf(sqlGroupTpl, "type", "type")).Rows(); err != nil {
		return
	}
	for rows.Next() {
		var stype int
		if err = rows.Scan(&stype); err != nil {
			return
		}
		types = append(types, stype)
	}

	//usage
	if err = tx.Model(&model.MobileMachine{}).Where("username<>'' and action>0").Count(&inUseCnt).Error; err != nil {
		return
	}
	if err = tx.Model(&model.MobileMachine{}).Where("username='' and action>0").Count(&freeCnt).Error; err != nil {
		return
	}

	if err = tx.Model(&model.MobileMachine{}).Where("action<0").Count(&noConCnt).Error; err != nil {
		return
	}

	if inUseCnt > 0 {
		usages = append(usages, model.MBInUse)
	}

	if freeCnt > 0 {
		usages = append(usages, model.MBFree)
	}

	if noConCnt > 0 {
		usages = append(usages, model.MBNoConnect)
	}

	mobileCategory = &model.MobileCategory{}
	mobileCategory.CPUs = cpus
	mobileCategory.Versions = versions
	mobileCategory.Modes = modes
	mobileCategory.States = states
	mobileCategory.Types = types
	mobileCategory.Usages = usages

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

// FindAllMobileImages Find All Mobile Images.
func (d *Dao) FindAllMobileImages() (mobileImages []*model.MobileImage, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.MobileImage{}).Find(&mobileImages).Error)
	return
}

// FindMobileImageByMode Find Mobile Image By Mode.
func (d *Dao) FindMobileImageByMode(mode string) (mobileImage *model.MobileImage, err error) {
	mobileImage = &model.MobileImage{}
	if err = d.db.Where("mode=?", mode).First(mobileImage).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// LendOutMobileMachine Lend Out Mobile.
func (d *Dao) LendOutMobileMachine(deviceID int64, serial, username string) (err error) {
	tx := d.db.Begin()

	if err = tx.Error; err != nil {
		return
	}

	mobileMachineLog := &model.MobileMachineLog{
		OperateType:   model.MBLendOutLog,
		Username:      username,
		MachineID:     deviceID,
		OperateResult: model.OperationSuccessForMachineLog,
	}
	if err = tx.Create(mobileMachineLog).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Model(&model.MobileMachine{}).Where("serial=?", serial).Update("is_lendout", model.MBLendOut).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

// ReturnMobileMachine return Mobile.
func (d *Dao) ReturnMobileMachine(deviceID int64, serial, username string) (err error) {
	tx := d.db.Begin()

	if err = tx.Error; err != nil {
		return
	}

	mobileMachineLog := &model.MobileMachineLog{
		OperateType:   model.MBReturnLog,
		Username:      username,
		MachineID:     deviceID,
		OperateResult: model.OperationSuccessForMachineLog,
	}
	if err = tx.Create(mobileMachineLog).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Model(&model.MobileMachine{}).Where("serial=?", serial).Update("is_lendout", model.MBOnSite).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}
