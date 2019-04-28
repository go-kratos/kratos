package service

import (
	"database/sql"
	"fmt"

	"go-common/app/admin/main/appstatic/model"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	// file type
	_fullPackage = 0
	_diffPackge  = 1
	// diff file name format
	_diffFormat = "Mod_%d/V_%d-V_%d.bspatch"
	// limit column
	deviceCol  = "device"
	mobiAppCol = "mobi_app"
	platCol    = "plat"
	buildCol   = "build"
	sysverCol  = "sysver"
	scaleCol   = "scale"
	levelCol   = "level"
	archCol    = "arch"
	// condition column
	_bk        = "bk"
	_wt        = "wt"
	buildLtCdt = "lt"
	buildGtCdt = "gt"
	buildLeCdt = "le"
	buildGeCdt = "ge"
	_valid     = 1
)

// GenerateVer generates a new version ( resource ) and cover the diff logic
func (s *Service) GenerateVer(resName string, limitData *model.Limit, fInfo *model.FileInfo, pool *model.ResourcePool, defPkg int) (resID int, version int, err error) {
	// create a new version
	var tx = s.DB.Begin()
	var resource = tx.Create(transResource(resName, pool.ID))
	if err = resource.Error; err != nil {
		log.Error("GenerateVer DBCreate Resource Error(%v)", err)
		tx.Rollback()
		return
	}
	resID = int(resource.Value.(*model.Resource).ID)
	log.Info("Resource Generated: ID = %d", resID)
	// create the full package in File Table
	if err = tx.Create(transFile(fInfo, resID)).Error; err != nil {
		log.Error("GenerateVer DBCreate ResourceFile Error(%v)", err)
		tx.Rollback()
		return
	}
	// create the resource config
	var config = tx.Create(transConfig(int64(resID), limitData))
	if err = config.Error; err != nil {
		log.Error("GenerateVer DBCreate ResoureConfig Error(%v)", err)
		tx.Rollback()
		return
	}
	configID := int64(config.Value.(*model.ResourceConfig).ID)
	log.Info("Resource Config Generated: ID = %d", configID)
	// create the resource limits
	limits := createLimit(configID, limitData)
	if len(limits) != 0 {
		for _, v := range limits {
			if err = tx.Create(v).Error; err != nil {
				log.Error("GenerateVer DCreate ResourceLimit (%v) Error(%v)", v, err)
				tx.Rollback()
				return
			}
		}
	} else {
		log.Error("[GenerateVer]-[createLimit]-No limit to create")
	}
	// commit the transaction
	tx.Commit()
	log.Info("Transaction Committed, ResID: %d", resID)
	// treat the default package setting
	if defPkg == 1 {
		if err = s.DefaultPkg(resID, pool.ID, configID); err != nil {
			log.Error("defaultPkg Error (%v)", err)
			return
		}
	}
	// defines the version of this resource
	if version, err = s.defineVer(resID, pool.ID); err != nil {
		log.Error("defineVer Error (%v)", err)
		return
	}
	// create diff records
	if err = s.createDiff(resID); err != nil {
		log.Error("[GenerateVer]-[createDiff]-Error(%v)", err)
		return
	}
	return
}

// DefaultPkg sets the resID's config as default package, and resets the other resources' config
func (s *Service) DefaultPkg(resID int, poolID int64, confID int64) (err error) {
	var (
		tx   = s.DB.Begin()
		rows *sql.Rows
		rid  int // resource id
	)
	// find out all the resources under the same pool, and put them as non-default pkg
	if rows, err = s.DB.Model(&model.Resource{}).Where("pool_id = ?", poolID).Select("id").Rows(); err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&rid); err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Model(&model.ResourceConfig{}).Where("resource_id = ?", rid).Update("default_package", 0).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	// defines the new package as the default pkg
	if err = tx.Model(&model.ResourceConfig{}).Where("id = ?", confID).Update("default_package", 1).Error; err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// defines the version of the resource after the transaction commited
func (s *Service) defineVer(resID int, poolID int64) (version int, err error) {
	var maxVer = model.Resource{}
	if err = s.DB.Where("id < ?", resID).Where("pool_id = ?", poolID).Order("version desc").First(&maxVer).Error; err == gorm.ErrRecordNotFound {
		err = nil
	}
	if err != nil {
		log.Error("GenerateVer DBFind ResourceVer (%d) Error(%v)", resID, err)
		return
	}
	version = int(maxVer.Version) + 1
	if err = s.DB.Model(&model.Resource{}).Where("id = ?", resID).Update("version", version).Error; err != nil {
		log.Error("GenerateVer DBUpdate ResourceVer (%d)-(%d) Error(%v)", resID, maxVer.Version+1, err)
		return
	}
	return
}

// create diff packages for the latest version with the history versions
func (s *Service) createDiff(resID int) (err error) {
	var (
		prodVers, testVers []int64
		currRes            *model.Resource
	)
	// pick history versions to calculate diff
	if prodVers, testVers, currRes, err = s.pickDiff(resID); err != nil {
		return
	}
	// put diff packages in our DB
	if err = s.putDiff(resID, mergeSlice(prodVers, testVers), currRes); err != nil {
		return
	}
	return
}

// pick history versions to calculate diff
func (s *Service) pickDiff(resID int) (prodVers []int64, testVers []int64, currRes *model.Resource, err error) {
	var (
		VersProd = []*model.Resource{} // prod
		VersTest = []*model.Resource{} // test
		res      = model.Resource{}
	)
	if err = s.DB.Where("id = ?", resID).First(&res).Error; err != nil {
		log.Error("[createDiff]-[FindCurrentRes]-Error(%v)", err)
		return
	}
	currRes = &res
	poolID := currRes.PoolID
	// calculate prod diffs
	if err = s.DB.Joins("LEFT JOIN resource_config ON resource.id = resource_config.resource_id").
		Where("resource.pool_id = ?", poolID).
		Where("resource.id < ?", resID).
		Where("resource_config.valid = ?", _valid).
		Order("resource.version desc").Limit(s.c.Cfg.HistoryVer).
		Select("resource.*").
		Find(&VersProd).Error; err != nil {
		log.Error("[createDiff]-[FindHistoryVers]-Error(%v)", err)
		return
	}
	log.Info("Get Prod History Versions: %d", len(VersProd))
	// calculate test diffs
	if err = s.DB.Joins("LEFT JOIN resource_config ON resource.id = resource_config.resource_id").
		Where("resource.pool_id = ?", poolID).
		Where("resource.id < ?", resID).
		Where("resource_config.valid != ?", _valid).
		Where("resource_config.valid_test = ?", _valid).
		Order("resource.version desc").Limit(s.c.Cfg.HistoryVer).
		Select("resource.*").
		Find(&VersTest).Error; err != nil {
		log.Error("[createDiff]-[FindHistoryVers]-Error(%v)", err)
		return
	}
	log.Info("Get Test History Versions: %d", len(VersTest))
	// merge slices
	prodVers = pickVersion(VersProd)
	testVers = pickVersion(VersTest)
	return
}

// put diff package in our DB
func (s *Service) putDiff(resID int, historyVers []int64, currRes *model.Resource) (err error) {
	for _, v := range historyVers {
		var diffPkg = &model.ResourceFile{
			Name:       fmt.Sprintf(_diffFormat, currRes.PoolID, v, currRes.Version),
			FromVer:    v,
			ResourceID: resID,
			FileType:   _diffPackge,
		}
		if err = s.DB.Create(diffPkg).Error; err != nil {
			log.Error("[createDiff]-[createDiffPkg]-Error(%v)", err)
			return
		}
	}
	log.Info("[createDiff]-Create (%d) Diff Pkg for ResID:(%d)", len(historyVers), resID)
	return
}

// pick resource version
func pickVersion(s1 []*model.Resource) (res []int64) {
	if len(s1) > 0 {
		for _, v := range s1 {
			res = append(res, v.Version)
		}
	}
	return
}

// merge int64 slices
func mergeSlice(s1 []int64, s2 []int64) []int64 {
	slice := make([]int64, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

// transform to Resource struct
func transResource(resMame string, id int64) *model.Resource {
	return &model.Resource{
		Name:    resMame,
		Version: 0, // will be updated after the transaction commits
		PoolID:  id,
	}
}

// transform to File struct
func transFile(fInfo *model.FileInfo, resID int) *model.ResourceFile {
	return &model.ResourceFile{
		Name:       fInfo.Name,
		Type:       fInfo.Type,
		Md5:        fInfo.Md5,
		Size:       int(fInfo.Size),
		URL:        fInfo.URL,
		ResourceID: resID,
		FileType:   _fullPackage,
	}
}

// transform to Config struct
func transConfig(resID int64, limitData *model.Limit) *model.ResourceConfig {
	var cfg = &model.ResourceConfig{
		ResourceID: resID,
		Valid:      0,
		IsDeleted:  0,
		IsWifi:     limitData.IsWifi,
	}
	if limitData.TimeRange != nil {
		cfg.Etime = limitData.TimeRange.Etime
		cfg.Stime = limitData.TimeRange.Stime
	}
	return cfg
}

// createLimit
func createLimit(configID int64, limitData *model.Limit) (res []*model.ResourceLimit) {
	// create device limit
	if len(limitData.Device) != 0 {
		generateDevice(limitData.Device, &res, configID, deviceCol, _bk)
	}
	// create plat limit
	if len(limitData.Plat) != 0 {
		generateDevice(limitData.Plat, &res, configID, platCol, _wt)
	}
	// create mobi_app limit
	if len(limitData.MobiApp) != 0 {
		generateDevice(limitData.MobiApp, &res, configID, mobiAppCol, _wt)
	}
	// scale, level, arch
	if len(limitData.Level) != 0 {
		generateDevice(limitData.Level, &res, configID, levelCol, _wt)
	}
	if len(limitData.Scale) != 0 {
		generateDevice(limitData.Scale, &res, configID, scaleCol, _wt)
	}
	if len(limitData.Arch) != 0 {
		generateDevice(limitData.Arch, &res, configID, archCol, _wt)
	}
	// create build & sysver limit
	if build := limitData.Build; build != nil {
		generateBuild(build, configID, &res, buildCol)
	}
	if sysver := limitData.Sysver; sysver != nil {
		generateBuild(sysver, configID, &res, sysverCol)
	}
	log.Info("createLimit creates %d limits", len(res))
	return
}

// generate build-like limits data (json, range), insert them into the slice 'res'
func generateBuild(build *model.Build, configID int64, res *[]*model.ResourceLimit, column string) {
	if build.GT != 0 {
		*res = append(*res, transBuild(buildGtCdt, configID, build.GT, column))
	}
	if build.LT != 0 {
		*res = append(*res, transBuild(buildLtCdt, configID, build.LT, column))
	}
	if build.GE != 0 {
		*res = append(*res, transBuild(buildGeCdt, configID, build.GE, column))
	}
	if build.LE != 0 {
		*res = append(*res, transBuild(buildLeCdt, configID, build.LE, column))
	}
}

// generate device-like limits data([]string), insert them into the slice 'res'
func generateDevice(device []string, res *[]*model.ResourceLimit, configID int64, col string, cdt string) {
	for _, v := range device {
		*res = append(*res, &model.ResourceLimit{
			ConfigID:  configID,
			Column:    col,
			Condition: cdt,
			Value:     v,
			IsDeleted: 0,
		})
	}
}

func transBuild(condition string, configID int64, value int, column string) *model.ResourceLimit {
	return &model.ResourceLimit{
		ConfigID:  configID,
		Column:    column,
		Condition: condition,
		Value:     fmt.Sprintf("%d", value),
		IsDeleted: 0,
	}
}
