package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"encoding/json"
	"go-common/app/admin/main/appstatic/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const nameFmt = `^[a-zA-Z0-9._-]+$`
const fileFmt = "Mod_%d-%s/%s"

func httpCode(c *bm.Context, message string, err error) {
	c.JSON(map[string]interface{}{
		"message": message,
	}, err)
}

// validate required data
func validateRequired(reqInfo *model.RequestVer) (err error) {
	reg := regexp.MustCompile(nameFmt)
	if res := reg.MatchString(reqInfo.ModName); !res {
		err = fmt.Errorf("mod_name %s contains illegal character", reqInfo.ModName)
		return
	}
	if res := reg.MatchString(reqInfo.Department); !res {
		err = fmt.Errorf("department %s contains illegal character", reqInfo.Department)
		return
	}
	return
}

// check whether the build range is valid
func checkRange(build *model.Build) (res bool) {
	if (build.GE != 0 && build.GT != 0) || (build.LE != 0 && build.LT != 0) { // two values by one side
		return false
	}
	var (
		gt = build.GT
		lt = build.LT
	)
	// transform E to T
	if build.GE != 0 {
		gt = build.GE - 1
	}
	if build.LE != 0 {
		lt = build.LE + 1
	}
	// range check
	if lt != 0 && gt != 0 && lt-gt <= 1 {
		return false
	}
	return true
}

// transform []int to []string
func sliceString(is []int) (ss []string) {
	for _, v := range is {
		ss = append(ss, fmt.Sprintf("%d", v))
	}
	return
}

// check limit data and build the Limit Struct, error is json error here
func checkLimit(reqInfo *model.RequestVer) (res *model.Limit, err error) {
	getFormat := "GetLimit Param (%s), Value = (%s)"
	res = &model.Limit{}
	// mobi_app
	if len(reqInfo.MobiAPP) != 0 {
		res.MobiApp = reqInfo.MobiAPP
	}
	// device
	if len(reqInfo.Device) != 0 {
		res.Device = reqInfo.Device
	}
	// plat
	if len(reqInfo.Plat) != 0 {
		res.Plat = reqInfo.Plat
	}
	if reqInfo.IsWifi != 0 {
		res.IsWifi = reqInfo.IsWifi
	}
	// Scale & Arch & Level
	if len(reqInfo.Scale) != 0 {
		res.Scale = sliceString(reqInfo.Scale)
	}
	if len(reqInfo.Arch) != 0 {
		res.Arch = sliceString(reqInfo.Arch)
	}
	if reqInfo.Level != 0 {
		res.Level = sliceString([]int{reqInfo.Level}) // treat level as others ( []int )
	}
	// build_range
	if buildStr := reqInfo.BuildRange; buildStr != "" {
		log.Info(getFormat, "build_range", buildStr)
		var build = model.Build{}
		if err = json.Unmarshal([]byte(buildStr), &build); err != nil { // json err
			log.Error("buildStr (%s) json.Unmarshal error(%v)", buildStr, err)
			return
		}
		if isValid := checkRange(&build); !isValid { // range not valid
			err = fmt.Errorf("build range (%s) not valid", buildStr)
			log.Error("buildStr CheckRange Error (%v)", err)
			return
		}
		res.Build = &build
	}
	// sysver
	if sysverStr := reqInfo.Sysver; sysverStr != "" {
		var build = model.Build{}
		if err = json.Unmarshal([]byte(sysverStr), &build); err != nil { // json err
			log.Error("buildStr (%s) json.Unmarshal error(%v)", sysverStr, err)
			return
		}
		if isValid := checkRange(&build); !isValid { // range not valid
			err = fmt.Errorf("build range (%s) not valid", sysverStr)
			log.Error("sysverStr CheckRange Error (%v)", err)
			return
		}
		res.Sysver = &build
	}
	// time_range
	if timeStr := reqInfo.TimeRange; timeStr != "" {
		log.Info(getFormat, "time_range", timeStr)
		var tr = model.TimeRange{}
		if err = json.Unmarshal([]byte(timeStr), &tr); err != nil {
			log.Error("timeStr (%s) json.Unmarshal error(%v)", timeStr, err)
			return
		}
		if tr.Stime != 0 && tr.Etime != 0 && tr.Stime > tr.Etime {
			err = fmt.Errorf("Stime(%d) is bigger than Etime(%d)", tr.Stime, tr.Etime)
			log.Error("Time Range Error(%v)", err)
			return
		}
		res.TimeRange = &tr
	}
	return
}

// validate the file type, content and upload it to the BFS storage
func validateFile(ctx *bm.Context, req *http.Request, pool *model.ResourcePool) (fInfo *model.FileInfo, err error) {
	// get the file
	file, header, err := req.FormFile("file")
	if err != nil {
		return
	}
	defer file.Close()
	// read the file
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("resource uploadFile.ReadAll error(%v)", err)
		return
	}
	// parse file, get type, size, md5
	fInfo, err = apsSvc.ParseFile(content)
	if err != nil {
		log.Error("[validateFile]-[ParseFile] Error-[%v]", err)
		return
	}
	if !apsSvc.TypeCheck(fInfo.Type) {
		log.Error("[validateFile]-[FileType] Error-[%v]", fInfo.Type)
		err = fmt.Errorf("请上传指定类型文件")
		return
	}
	// regex checking
	reg := regexp.MustCompile(nameFmt)
	if res := reg.MatchString(header.Filename); !res {
		err = fmt.Errorf("fileName %s contains illegal character", header.Filename)
		return
	}
	// upload file to BFS
	fInfo.Name = fmt.Sprintf(fileFmt, pool.ID, fInfo.Md5, header.Filename) // rename with the MD5 and poolID
	location, err := apsSvc.Upload(ctx, fInfo.Name, fInfo.Type, time.Now().Unix(), content)
	if err != nil {
		log.Error("[validateFile]-[UploadBFS] Error-[%v]", err)
		return
	}
	fInfo.URL = location
	return
}

// for other systems
func addVer(c *bm.Context) {
	var (
		pool       = model.ResourcePool{}
		department = model.Department{}
		req        = c.Request
		limitData  *model.Limit
		fInfo      *model.FileInfo
		err        error
		reqInfo    = model.RequestVer{}
		respData   = &model.RespAdd{}
	)
	req.ParseMultipartForm(apsSvc.MaxSize)
	if err = c.Bind(&reqInfo); err != nil {
		return
	}
	// validate required data
	if err = validateRequired(&reqInfo); err != nil {
		log.Error("addVer ModName, ResName Error (%v)", err)
		c.JSON(nil, err)
		return
	}
	// validate department
	if err = apsSvc.DB.Where("`name` = ?", reqInfo.Department).First(&department).Error; err != nil {
		log.Error("addVer First department Error (%v)", err)
		httpCode(c, fmt.Sprintf("department %s doesn't exist", reqInfo.Department), ecode.RequestErr)
		return
	}
	// validate mod Name
	if err = apsSvc.DB.Where("`name` = ? AND `department_id` = ? AND `deleted` = 0 AND `action` = 1", reqInfo.ModName, department.ID).First(&pool).Error; err != nil {
		log.Error("addVer First Pool Error (%v)", err)
		httpCode(c, fmt.Sprintf("Mod_name %s doesn't exist", reqInfo.ModName), ecode.RequestErr)
		return
	}
	// check limit & config data
	if limitData, err = checkLimit(&reqInfo); err != nil {
		log.Error("addVer CheckLimit Error (%v)", err)
		httpCode(c, fmt.Sprintf("Limit Params JSON Error:(%v)", err), ecode.RequestErr)
		return
	}
	// validate file data
	if fInfo, err = validateFile(c, req, &pool); err != nil {
		log.Error("addVer ValidateFile Error (%v)", err)
		httpCode(c, fmt.Sprintf("File Error:(%v)", err), ecode.RequestErr)
		return
	}
	// DB & storage operation
	if respData.ResID, respData.Version, err = apsSvc.GenerateVer(reqInfo.ResName, limitData, fInfo, &pool, reqInfo.DefaultPackage); err != nil {
		log.Error("addVer GenerateVer Error (%v)", err)
		httpCode(c, fmt.Sprintf("Generate Version Error:(%v)", err), ecode.ServerErr)
		return
	}
	c.JSON(respData, nil)
}
