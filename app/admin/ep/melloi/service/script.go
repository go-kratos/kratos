package service

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

//QueryScripts query scripts
func (s *Service) QueryScripts(script *model.Script, pn, ps int) (scripts []*model.Script, err error) {
	if ps <= 0 || pn <= 0 {
		ps = 10
		pn = 1
	}
	if scripts, err = s.dao.QueryScripts(script, pn, ps); err != nil {
		log.Error("s.dao.QueryScripts error :(%v)", err)
		return
	}
	for _, script := range scripts {
		if script.APIHeader != "" {
			if err = json.Unmarshal([]byte(script.APIHeader), &script.Headers); err != nil {
				log.Error("get script header err : (%v), scriptId:(%d)", err, script.ID)
			}
		}
		if script.ArgumentString != "" {
			if err = json.Unmarshal([]byte(script.ArgumentString), &script.ArgumentsMap); err != nil {
				log.Error("get script argument err: (%v), scriptId:(%d)", err, script.ID)
			}
		}
		if script.OutputParams != "" {
			if err = json.Unmarshal([]byte(script.OutputParams), &script.OutputParamsMap); err != nil {
				log.Error("get script OutputParams err: (%v), scriptId:(%d)", err, script.ID)
			}
		}
	}
	return
}

//CountQueryScripts count query scripts
func (s *Service) CountQueryScripts(script *model.Script) (total int) {
	return s.dao.CountQueryScripts(script)
}

//QueryScriptSnap query scriptSnap
func (s *Service) QueryScriptSnap(scriptSnap *model.ScriptSnap) (snaps []*model.ScriptSnap, err error) {
	return s.dao.QueryScriptSnap(scriptSnap)
}

//QueryScriptsByPage query script by page
func (s *Service) QueryScriptsByPage(c context.Context, sessionID string, qsrq *model.QueryScriptRequest) (rsp *model.QueryScriptResponse, err error) {
	// 获取服务树节点
	var (
		treeNodes  []string
		treeNodesd []string
	)
	if treeNodesd, err = s.QueryUserRoleNode(c, sessionID); err != nil {
		log.Error("QueryUserRoleNode err (%v):", err)
	}
	treeNodes = append(treeNodesd, "")
	if ExistsInSlice(qsrq.Executor, conf.Conf.Melloi.Executor) {
		//rsp, err = s.dao.QueryScriptsByPageWhiteName(&qsrq.Script, qsrq.PageNum, qsrq.PageSize)
		if rsp, err = s.dao.QueryScriptsByPageWhiteName(&qsrq.Script, qsrq.PageNum, qsrq.PageSize); err != nil {
			return
		}
	} else {
		if rsp, err = s.dao.QueryScriptsByPage(&qsrq.Script, qsrq.PageNum, qsrq.PageSize, treeNodes); err != nil {
			return
		}
	}

	for _, script := range rsp.Scripts {
		if script.APIHeader != "" {
			if err = json.Unmarshal([]byte(script.APIHeader), &script.Headers); err != nil {
				log.Error("get script header err : (%v), scriptId:(%d)", err, script.ID)
			}
		}
		if script.ArgumentString != "" {
			if err = json.Unmarshal([]byte(script.ArgumentString), &script.ArgumentsMap); err != nil {
				log.Error("get script argument err: (%v), scriptId:(%d)", err, script.ID)
			}
		}
		if script.OutputParams != "" {
			if err = json.Unmarshal([]byte(script.OutputParams), &script.OutputParamsMap); err != nil {
				log.Error("get script OutputParams err: (%v), scriptId:(%d)", err, script.ID)
			}
		}
	}

	// 获取label
	for _, script := range rsp.Scripts {
		lr := model.LabelRelation{Type: model.ScriptType, TargetID: int64(script.ID)}
		if script.Labels, err = s.dao.QueryLabelRelation(&lr); err != nil {
			return
		}
	}
	return
}

//AddAndExcuScript add script and excute it
func (s *Service) AddAndExcuScript(c context.Context, script *model.Script, cookie string, scene *model.Scene, fileWrite, convjson bool) (resp model.DoPtestResp, err error) {
	var (
		scriptID     int
		groupID      int
		runOrder     int
		testNameNick string
		message      string
		scriptd      model.Script
		paramList    *model.ParamList
		isEdit       = false
	)

	if message, scriptd, err = s.AddScript(script, fileWrite); err != nil {
		log.Error("sorry, addScript fail : (%v)", err)
		return
	}
	headerString := MapToString(scriptd.Headers)
	argumentString := MapToString(scriptd.ArgumentsMap)
	scriptd.ArgumentString = argumentString
	scriptd.APIHeader = headerString
	scriptd.OutputParams = MapToString(scriptd.OutputParamsMap)
	//场景压测中选择已有接口时，去修改默认的熔断成功率的值（新增接口时不修改）
	if scene.IsBatch {
		scriptd.Fusing = scene.Fusing
	}

	if convjson {
		scriptd.ArgumentString = script.ArgumentString
		scriptd.APIHeader = script.APIHeader
		scriptd.OutputParams = script.OutputParams
	}
	//场景压测独有的逻辑
	if scriptd.TestType == model.SCENE_SCRIPT_TYPE {
		//根据接口参数依赖，计算出接口分组与执行顺序
		if paramList, err = s.dao.QueryParams(&scriptd, scene); err != nil {
			log.Error("s.dao.QueryParams err :(%v)", err)
			return
		}
		scriptd.GroupID, scriptd.RunOrder = GroupOrder(isEdit, &scriptd, scene, paramList)
	}

	//非debug 执行的 循环次数为-1 即永久
	if scriptd.Loops == -1 {
		if len(script.OutputParamsMap) == 0 || script.OutputParams == "[]" {
			scriptd.OutputParams = "[{\"\":\"\"}]"
		}
		if scriptID, groupID, runOrder, err = s.dao.AddScript(&scriptd); err != nil {
			log.Error("s.dao.AddScript err : (%v)", err)
			message = "sql excute err"
			return
		}
	}

	//如果不是复制脚本操作 或者 不是保存，则执行jmeter压测
	if !script.IsCopy && !script.IsSave {
		tim := strconv.FormatInt(time.Now().Unix(), 10)
		testNameNick = scriptd.TestName + tim
		log.Info("开始调用压测job-------\n")
		ptestParam := model.DoPtestParam{
			UserName:            scriptd.UpdateBy,                // 用户名
			LoadTime:            scriptd.LoadTime,                //运行时间
			TestNames:           StringToSlice(scriptd.TestName), //接口名转数组
			FileName:            scriptd.SavePath,                // jmx文件
			ResJtl:              scriptd.ResJtl,                  // jtl 文件
			JmeterLog:           scriptd.JmeterLog,               // jmeterlog时间戳
			Department:          script.Department,
			Project:             script.Project,
			APP:                 script.App,
			ScriptID:            scriptID,
			IsDebug:             script.IsDebug, //false
			Cookie:              cookie,         // 用不到
			URL:                 scriptd.URL,    // 微信通知URL
			LabelIDs:            script.LabelIds,
			Domain:              scriptd.Domain,
			FileSplit:           script.FileSplit, // 文件切割
			SplitNum:            script.SplitNum,  // 切割数量
			Fusing:              script.Fusing,
			APIHeader:           scriptd.APIHeader,
			ExecuDockerSum:      script.ExecuDockerSum,
			UseBusinessStop:     script.UseBusinessStop,
			BusinessStopPercent: script.BusinessStopPercent,
		}
		if resp, err = s.DoPtestByJmeter(c, ptestParam, StringToSlice(testNameNick)); err != nil {
			log.Error("DoPtestByJmeter err :(%v)", err)
			message = "DoPtestByJmeter err "
			return
		}
		resp.Message = message
	}
	resp.ScriptID = scriptID
	resp.GroupID = groupID
	resp.RunOrder = runOrder
	return
}

//AddScript add script
func (s *Service) AddScript(script *model.Script, fileWrite bool) (message string, scriptd model.Script, err error) {
	var (
		buff           *template.Template
		file           *os.File
		scriptPath     string
		u              *url.URL
		sinfoJSON      []byte
		JSON           []byte
		JSONExtractors string
		dataFileName   string
		Assertionsd    []string
	)
	log.Info("threadsSum :(%d)", script.ThreadsSum)
	if script.LoadTime > s.c.Jmeter.TestTimeLimit {
		script.LoadTime = s.c.Jmeter.TestTimeLimit
	}
	//判断url
	if !strings.Contains(script.URL, "http://") && !strings.Contains(script.URL, "https://") {
		script.URL = "http://" + script.URL
	}
	if script.Delimiter == "" {
		script.Delimiter = ","
	}
	if u, err = url.Parse(script.URL); err != nil || u == nil {
		log.Error("sorry,url is not standard,error(%v)", err)
		err = ecode.MelloiUrlParseError
		return
	}
	if script.Data == "" {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterScUcodedTmp)
	} else {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterScTmp)
	}
	if err != nil {
		log.Error("open file failed! error %v", err)
		return
	}
	if script.Assertion != "" {
		script.UseAssertion = true
		Assertions := strings.Split(script.Assertion, "|")
		for _, assert := range Assertions {
			Assertionsd = append(Assertionsd, model.AssertionStart+assert+model.AssertionEnd)
		}
		script.AssertionString = unescaped(SliceToString(Assertionsd, " "))
	}
	if len(script.OutputParamsMap) > 0 && script.OutputParams != "[{\"\":\"\"}]" && script.OutputParams != "[]" && script.OutputParams != "" {
		if JSONExtractors, err = s.SetJSONExtractors(script); err != nil {
			return
		}
		if JSONExtractors != "" {
			script.JSONExtractor = unescaped(JSONExtractors)
		}
	}

	// POST 请求的 binary 逻辑
	if script.MultipartFile != "" && script.MultipartPath != "" {
		script.UseMultipart = true
		multiPartInfo := SetMultiPartInfo(script)
		script.MultiPartInfo = unescaped(multiPartInfo)
	}
	stringHeader := SetHeaders(script.Headers, script.ID)
	ArgmentStr := SetArguments(script.ArgumentsMap, script.ID)
	ArgmentString := unescaped(ArgmentStr)
	script.Arguments = ArgmentString
	stringHeaderd := unescaped(stringHeader)
	script.HeaderString = stringHeaderd
	script.ProcType = u.Scheme
	script.Domain = u.Hostname()
	script.Port = u.Port()
	if u.RawQuery == "" {
		script.Path = u.Path
		script.RowQuery = "appkey=" + conf.Conf.Melloi.AppkeyProd
	} else {
		script.Path = u.Path + "?" + u.RawQuery
		script.RowQuery = u.RawQuery + "&appkey=" + conf.Conf.Melloi.AppkeyProd
	}

	if script.UseSign {
		script.Path = u.Path + "?${params}"
	}
	dataFileName = script.FileName
	script.FileName = script.ScriptPath + script.FileName
	if script.ScriptPath != "" && script.FileName == "" && !script.Upload {
		message = "dataFile is not exists"
		log.Error(message)
		return
	}

	if script.Loops == 0 {
		script.Loops = -1
	}
	//极速模式（线程NIO模式）
	if script.IsAsync {
		script.AsyncInfo = unescaped(model.AsyncInfo)
	}
	//生成压测脚本
	if sinfoJSON, err = json.Marshal(script); err != nil {
		log.Error("script err (%v):", err)
		return
	}
	sc := string(sinfoJSON)
	log.Info("script :(%s)", sc)
	//判断是否使用内置参数
	if strings.Contains(sc, "${mid}") || strings.Contains(sc, "${access_key}") || script.UseSign {
		script.UseBuiltinParam = true
	}
	//fileWrite = true ,需要创建文件，为 false ，不需要创建文件
	if fileWrite {
		//创建脚本保存路径
		if script.ScriptPath == "" {
			if scriptPath, err = s.uniqueFolderPath(script.SavePath); err != nil {
				return
			}
		} else {
			scriptPath = script.ScriptPath
		}
		SavePath := scriptPath + script.TestName + "/"
		log.Info("SavePath****:  " + SavePath)
		if err = os.MkdirAll(SavePath, 0755); err != nil {
			log.Error("Create SavePath Err :(%v)", err)
			return
		}
		///创建脚本保存路径 脚本存放路径：部门/项目/应用/项目名/,须判断路径是否包含.jmx
		if !strings.Contains(script.SavePath, ".jmx") {
			// 创建脚本文件 部门/项目/应用/项目名/**.jmx
			if file, err = os.Create(SavePath + script.TestName + ".jmx"); err != nil {
				log.Error("create file error :(%v)", err)
				message = "create .jmx error"
				return
			}
			defer file.Close()
			//jtl 路径
			script.ResJtl = SavePath + script.TestName + ".jtl"
			//jmeter-log 路径
			script.JmeterLog = SavePath + script.TestName + ".log"

		} else {
			// 如果路径钟包含 .jmx ,则直接根据路径创建
			if file, err = os.Create(script.SavePath); err != nil {
				log.Error("create file error :(%v)", err)
				message = "create .jmx error"
				return
			}
			defer file.Close()
		}
		buff = buff.Funcs(template.FuncMap{"unescaped": unescaped})
		buff.Execute(io.Writer(file), script)
		script.SavePath = file.Name()
	}
	script.Ctime = time.Now()
	script.ProjectName = script.TestName
	if JSON, err = json.Marshal(script); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(JSON), &scriptd); err != nil {
		log.Error("sorry,scriptd err(%v)", err)
		return
	}
	scriptd.HeaderString = stringHeaderd
	scriptd.Arguments = ArgmentString
	scriptd.JSONExtractor = script.JSONExtractor
	scriptd.Data = script.Data
	scriptd.FileName = dataFileName
	if script.FileName == "" || script.ParamsName == "" {
		scriptd.UseDataFile = false
	} else {
		scriptd.UseDataFile = true
	}

	if len(Assertionsd) > 0 {
		scriptd.AssertionString = unescaped(SliceToString(Assertionsd, " "))
	}
	return
}

//DelScript sel script
func (s *Service) DelScript(id int) error {
	return s.dao.DelScript(id)
}

//UpdateScript update script info
func (s *Service) UpdateScript(script *model.Script) (updateResult string, err error) {
	var (
		buff           *template.Template
		stringHeader   string
		Headers        []map[string]string
		stringArgument string
		Arguments      []map[string]string
		scriptd        *model.Script
		u              *url.URL
		file           *os.File
		sinfoJSON      []byte
		sinfoJSONd     []byte
		JSONExtractors string
		Assertionsd    []string
	)
	if scriptd, err = s.dao.QueryScriptByID(script.ID); err != nil {
		log.Error("query script fail , error %v", err)
		return
	}
	if script.Data == "" {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterScUcodedTmp)
	} else {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterScTmp)
	}
	if err != nil {
		log.Error("open file failed! error %v", err)
		return
	}
	if u, err = url.Parse(script.URL); err != nil || u == nil {
		log.Error("sorry,url is not standard,error(%v)", err)
		err = ecode.MelloiUrlParseError
		return
	}
	// 修改laodTime 和 threadsSum
	// ThreadGroup.ramp_time">{{.ReadyTime}}</stringProp>
	if file, err = os.Create(scriptd.SavePath); err != nil {
		log.Error("create file failed! error %v", err)
		return
	}
	defer file.Close()
	if err = json.Unmarshal([]byte(scriptd.APIHeader), &Headers); err != nil {
		updateResult = " scriptd.APIHeade json 不合法"
		log.Error("scriptd.APIHeader switch map failed! error %v", err)
		return
	}

	if err = json.Unmarshal([]byte(scriptd.ArgumentString), &Arguments); err != nil {
		updateResult = "scriptd.ArgumentString  json 不合法"
		log.Error("scriptd.ArgumentString switch map failed! error %v", err)
		return
	}
	if len(scriptd.OutputParamsMap) > 0 && scriptd.OutputParams != "" && scriptd.OutputParams != "[{\"\":\"\"}]" && scriptd.OutputParams != "[]" {
		if err = json.Unmarshal([]byte(scriptd.OutputParams), &scriptd.OutputParamsMap); err != nil {
			log.Error("scriptd.OutputParams switch map failed! error %v", err)
			return
		}
	}
	stringHeader = SetHeaders(Headers, scriptd.ID)
	stringArgument = SetArguments(Arguments, scriptd.ID)
	if len(scriptd.OutputParamsMap) > 0 && scriptd.OutputParams != "[{\"\":\"\"}]" && scriptd.OutputParams != "[]" && script.OutputParams != "" {
		if JSONExtractors, err = s.SetJSONExtractors(scriptd); err != nil {
			return
		}
		if JSONExtractors != "" {
			script.JSONExtractor = unescaped(JSONExtractors)
		}
	}
	scriptd.HeaderString = unescaped(stringHeader)
	scriptd.Arguments = unescaped(stringArgument)
	scriptd.Data = script.Data
	scriptd.Assertion = script.Assertion
	scriptd.ReadyTime = script.ReadyTime
	scriptd.ThreadsSum = script.ThreadsSum
	scriptd.LoadTime = script.LoadTime
	scriptd.ProcType = u.Scheme
	scriptd.Domain = u.Hostname()
	scriptd.Port = u.Port()
	scriptd.TestName = script.TestName
	scriptd.FileName = scriptd.ScriptPath + scriptd.FileName
	if scriptd.Assertion != "" {
		scriptd.UseAssertion = true
		Assertions := strings.Split(scriptd.Assertion, "|")
		for _, assert := range Assertions {
			Assertionsd = append(Assertionsd, model.AssertionStart+assert+model.AssertionEnd)
		}
		scriptd.AssertionString = unescaped(SliceToString(Assertionsd, " "))
	}

	//判断是否使用内置参数
	if sinfoJSON, err = json.Marshal(script); err != nil {
		log.Error("script err :(%s)", err)
		return
	}
	if sinfoJSONd, err = json.Marshal(scriptd); err != nil {
		log.Error("script err :(%s)", err)
		return
	}
	sc := string(sinfoJSON) + string(sinfoJSONd)
	if strings.Contains(sc, "${mid}") || strings.Contains(sc, "${access_key}") || script.UseSign {
		scriptd.UseBuiltinParam = true
	}
	if u.RawQuery == "" {
		scriptd.Path = u.Path
		scriptd.RowQuery = "appkey=" + conf.Conf.Melloi.AppkeyProd
	} else {
		scriptd.Path = u.Path + "?" + u.RawQuery
		scriptd.RowQuery = u.RawQuery + "&appkey=" + conf.Conf.Melloi.AppkeyProd
	}
	if scriptd.UseSign {
		scriptd.Path = u.Path + "?${params}"
	}
	// POST 请求的 binary 逻辑
	if scriptd.MultipartFile != "" && scriptd.MultipartPath != "" {
		scriptd.UseMultipart = true
		multiPartInfo := SetMultiPartInfo(scriptd)
		scriptd.MultiPartInfo = unescaped(multiPartInfo)
	}
	if scriptd.IsAsync {
		scriptd.AsyncInfo = unescaped(model.AsyncInfo)
	}
	buff = buff.Funcs(template.FuncMap{"unescaped": unescaped})
	buff.Execute(io.Writer(file), scriptd)
	script.ProjectName = script.TestName
	script.KeepAlive = scriptd.KeepAlive
	if err = s.dao.UpdateScriptPart(script); err != nil {
		log.Error("s.dao.UpdateScript err :(%v)", err)
		return
	}
	updateResult = "success"
	return
}

//UpdateScriptAll update script all
func (s *Service) UpdateScriptAll(script *model.Script, scene *model.Scene) (updateResult string, err error) {
	var (
		scriptd   model.Script
		paramList *model.ParamList
		isEdit    = true
		groupId   int
		runOrder  int
		fileWrite bool
	)
	id := script.ID
	script.ID = 0
	// script test_type 为1 表示 http 脚本
	if script.TestType == model.HTTP_SCRIPT_TYPE {
		fileWrite = true
	}
	if _, scriptd, err = s.AddScript(script, fileWrite); err != nil {
		log.Error("sorry, addScript fail : (%v)", err)
		return
	}

	if script.TestType == model.SCENE_SCRIPT_TYPE {
		//根据接口参数依赖，计算出接口分组与执行顺序
		if paramList, err = s.dao.QueryParams(&scriptd, scene); err != nil {
			log.Error("s.dao.QueryParams err :(%v)", err)
			return
		}
		groupId, runOrder = GroupOrderByEdit(isEdit, id, &scriptd, scene, paramList)
		if groupId != 0 && runOrder != 0 {
			scriptd.GroupID = groupId
			scriptd.RunOrder = runOrder
		}
	}

	scriptd.ID = id
	scriptd.SavePath = ""

	if scriptd.Loops == -1 {
		headerString := MapToString(scriptd.Headers)
		argumentString := MapToString(scriptd.ArgumentsMap)
		scriptd.OutputParams = MapToString(scriptd.OutputParamsMap)
		scriptd.ArgumentString = argumentString
		scriptd.APIHeader = headerString
		if err = s.dao.UpdateScript(&scriptd); err != nil {
			log.Error("s.dao.UpdateScript err : (%v)", err)
			return
		}
	}
	updateResult = "success"
	return
}

//AddJmeterSample add jmeter sample
func (s *Service) AddJmeterSample(script *model.Script) (result string, err error) {
	var (
		buff           *template.Template
		u              *url.URL
		b              = bytes.NewBuffer(nil)
		JSONExtractors string
		Assertionsd    []string
	)

	if script.Data == "" {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterSampleTmp)
	} else {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterSamplePostTmp)
	}
	if err != nil {
		log.Info("open template failed! (%v)", err)
		return
	}
	//判断url
	if !strings.Contains(script.URL, "http://") && !strings.Contains(script.URL, "https://") {
		script.URL = "http://" + script.URL
	}
	if script.Delimiter == "" {
		script.Delimiter = ","
	}
	if u, err = url.Parse(script.URL); err != nil || u == nil {
		log.Error("sorry,url is not standard,error(%v)", err)
		err = ecode.MelloiUrlParseError
		return
	}
	stringHeader := SetHeaders(script.Headers, script.ID)
	ArgmentStr := SetArguments(script.ArgumentsMap, script.ID)
	ArgmentString := unescaped(ArgmentStr)
	if len(script.OutputParamsMap) > 0 && script.OutputParams != "[{\"\":\"\"}]" && script.OutputParams != "[]" && script.OutputParams != "" {
		if JSONExtractors, err = s.SetJSONExtractors(script); err != nil {
			return
		}
		if JSONExtractors != "" {
			script.JSONExtractor = unescaped(JSONExtractors)
		}
	}
	script.Arguments = ArgmentString
	stringHeaderd := unescaped(stringHeader)
	script.HeaderString = stringHeaderd
	script.ProcType = u.Scheme
	script.Domain = u.Hostname()
	script.Port = u.Port()
	script.FileName = script.ScriptPath + script.FileName
	if u.RawQuery == "" {
		script.Path = u.Path
		script.RowQuery = "appkey=" + conf.Conf.Melloi.AppkeyProd
	} else {
		script.Path = u.Path + "?" + u.RawQuery
		script.RowQuery = u.RawQuery + "&appkey=" + conf.Conf.Melloi.AppkeyProd
	}
	if script.UseSign {
		script.Path = u.Path + "?${params}"
	}
	if script.Assertion != "" {
		script.UseAssertion = true
		Assertions := strings.Split(script.Assertion, "|")
		for _, assert := range Assertions {
			Assertionsd = append(Assertionsd, model.AssertionStart+assert+model.AssertionEnd)
		}
		script.AssertionString = unescaped(SliceToString(Assertionsd, " "))
	}
	if script.ConstTimer > 0 && script.RandomTimer == 0 {
		constTimer := strings.Replace(model.ConstTimer, "1000", strconv.Itoa(script.ConstTimer), -1)
		script.ConstTimerInfo = unescaped(constTimer)
	}
	if script.ConstTimer > 0 && script.RandomTimer > 0 {
		constTimerd := strings.Replace(model.RandomTimer, "1000", strconv.Itoa(script.ConstTimer), -1)
		randomTimerd := strings.Replace(constTimerd, "500", strconv.Itoa(script.RandomTimer), -1)
		script.RandomTimerInfo = unescaped(randomTimerd)
	}
	buff = buff.Funcs(template.FuncMap{"unescaped": unescaped})
	buff.Execute(b, script)
	result = b.String()
	return
}

//AddThreadGroup add thread group
func (s *Service) AddThreadGroup(script *model.Script, sceneType int) (result string, err error) {
	var (
		buff    *template.Template
		scriptd model.Script
		b       = bytes.NewBuffer(nil)
	)

	if _, scriptd, err = s.AddScript(script, false); err != nil {
		log.Error("s.AddScript err :(%v)", err)
		return
	}
	//并行get
	if sceneType == 1 && script.Data == "" {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterThGroupTmp)
	}

	//并行post
	if sceneType == 1 && script.Data != "" {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterThGroupPostTmp)
	}

	//串行get
	if sceneType == 2 && script.Data == "" {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterThGroupDuliTmp)
	}

	//串行 post
	if sceneType == 2 && script.Data != "" {
		buff, err = template.ParseFiles(s.c.Jmeter.JmeterThGroupPostDuliTmp)
	}

	if err != nil {
		log.Info("open template failed! (%v)", err)
		return
	}
	buff = buff.Funcs(template.FuncMap{"unescaped": unescaped})
	scriptd.FileName = scriptd.ScriptPath + scriptd.FileName
	buff.Execute(b, scriptd)
	result = b.String()
	return
}

//GetThreadGroup get thread group
func (s *Service) GetThreadGroup(scrThreadGroup model.ScrThreadGroup) (threadGroup string, err error) {
	var (
		threadGroups   []string
		IsSerial       bool
		groups         []int
		groupMap       = make(map[int]int)
		RunOrders      []int
		scriptsd       []*model.Script
		scriptsds      [][]*model.Script
		threadGroupd   string
		threadGroupds  []string
		threadSamples  []string
		threadSample   string
		threadSamplesd string
		threadGrou     string
	)
	for _, script := range scrThreadGroup.Scripts {
		RunOrders = append(RunOrders, script.RunOrder)
		groupMap[script.GroupID]++
		if script.APIHeader != "" && script.APIHeader != "[{\"\":\"\"}]" {
			if err = json.Unmarshal([]byte(script.APIHeader), &script.Headers); err != nil {
				log.Error("get script header err : (%v), scriptId:(%d)", err, script.ID)
			}
		}
		if script.ArgumentString != "" && script.ArgumentString != "[{\"\":\"\"}]" {
			if err = json.Unmarshal([]byte(script.ArgumentString), &script.ArgumentsMap); err != nil {
				log.Error("get script argument err: (%v), scriptId:(%d)", err, script.ID)
			}
		}
		if script.OutputParams != "" && script.OutputParams != "[{\"\":\"\"}]" && script.OutputParams != "[]" {
			if err = json.Unmarshal([]byte(script.OutputParams), &script.OutputParamsMap); err != nil {
				log.Error("get script OutputParams err: (%v), scriptId:(%d)", err, script.ID)
			}
		}
	}
	for _, order := range RunOrders {
		if order > 1 {
			IsSerial = true
			break
		}
	}
	for k := range groupMap {
		groups = append(groups, k)
	}
	//并行脚本逻辑
	if !IsSerial {
		for _, script := range scrThreadGroup.Scripts {
			if threadGrou, err = s.AddThreadGroup(script, 1); err != nil {
				log.Error("d.GetThroupGroupJmeter:  error(%v)", err)
				return
			}
			threadGroups = append(threadGroups, threadGrou)
		}
		for _, thgroup := range threadGroups {
			threadGroupd = threadGroupd + thgroup
		}
		threadGroup = threadGroupd
		return
	}

	//串行脚本逻辑
	for _, group := range groups {
		for i := 0; i < len(scrThreadGroup.Scripts); i++ {
			if scrThreadGroup.Scripts[i].GroupID == group {
				scriptsd = append(scriptsd, scrThreadGroup.Scripts[i])
			}
		}
		scriptsds = append(scriptsds, scriptsd)
		scriptsd = []*model.Script{}
	}

	// 每个线程组内的接口，按照接口的 RunOrder 属性进行排序
	for _, scriptsd := range scriptsds {
		for i := 0; i < len(scriptsd)-1; i++ {
			for j := i + 1; j < len(scriptsd); j++ {
				if scriptsd[i].RunOrder > scriptsd[j].RunOrder {
					tmp := scriptsd[i]
					scriptsd[i] = scriptsd[j]
					scriptsd[j] = tmp
				}
			}
		}
	}

	for k := 0; k < len(scriptsds); k++ {
		if len(scriptsds[k]) == 1 {
			if threadGroupd, err = s.AddThreadGroup(scriptsds[k][0], 1); err != nil {
				log.Error("d.GetThroupGroupJmeter:  error(%v)", err)
				return
			}
		} else {
			//生成一个线程组，再往线程组插入 sample
			if threadGroupd, err = s.AddThreadGroup(scriptsds[k][0], 2); err != nil {
				log.Error(" s.AddScript err :(%v)", err)
				return
			}
			for i := 0; i < len(scriptsds[k]); i++ {
				//从第二个接口开始生成 sample
				if i > 0 {
					if threadSample, err = s.AddJmeterSample(scriptsds[k][i]); err != nil {
						log.Error("s.dao.GetThroupSeriesSample err :(%v)", err)
						return
					}
					threadSamples = append(threadSamples, threadSample)
				}
			}
		}
		for j := 0; j < len(threadSamples); j++ {
			threadSamplesd = threadSamplesd + threadSamples[j]
		}
		//生成一个 线程组
		threadGroupM := strings.Replace(threadGroupd, "+-+-+*", threadSamplesd, -1)

		//多个线程组
		threadGroupds = append(threadGroupds, threadGroupM)
		threadSamples = []string{}
		threadSamplesd = ""
	}

	for _, threadGrou := range threadGroupds {
		threadGroup = threadGroup + threadGrou
	}
	return
}

//URLCheck url check
func (s *Service) URLCheck(script *model.Script) (urlEncode *model.URLEncode, err error) {
	var (
		tempUrl *url.URL
	)
	urlEncode = &model.URLEncode{}
	if strings.Contains(script.URL, "?") {
		if tempUrl, err = url.Parse(script.URL); err != nil {
			return
		}
		params := tempUrl.RawQuery
		paramList := strings.Split(params, "&")
		for _, param := range paramList {
			value := strings.Split(param, "=")[1]
			var myMap map[string]interface{}
			if err = json.Unmarshal([]byte(value), &myMap); err != nil {
				urlEncode.ParamsType = "nojson"
				urlEncode.NewUrl = script.URL
				continue
			}
			urlEncode.ParamsType = "json"
			newParams := tempUrl.Query().Encode()
			urlEncode.NewUrl = strings.Split(script.URL, "?")[0] + "?" + newParams
			return
		}
	}
	return
}

//AddTimer add Timer
func (s *Service) AddTimer(script *model.Script) error {
	return s.dao.AddScriptTimer(script)
}

//SetHeaders set headers
func SetHeaders(Headers []map[string]string, scriptId int) (stringHeader string) {
	var (
		stringHeaders []string
		HeaderString  string
	)

	// 表示直接生成的脚本，非修改
	if scriptId == 0 {
		for _, header := range Headers {
			k := header["key"]
			v := header["value"]
			HeaderString = model.HeaderStart + k + model.HeaderMid + v + model.HeaderEnd
			stringHeaders = append(stringHeaders, HeaderString)
		}
	} else {
		//修改脚本，header 从数据库取出的
		for _, header := range Headers {
			for k, v := range header {
				HeaderString = model.HeaderStart + k + model.HeaderMid + v + model.HeaderEnd
				stringHeaders = append(stringHeaders, HeaderString)
			}
		}
	}

	for _, str := range stringHeaders {
		stringHeader = stringHeader + str + "\r\n"
	}
	return
}

//SetArguments  set argument
func SetArguments(Arguments []map[string]string, scriptId int) (stringArgument string) {
	var (
		stringArguments []string
		ArgumentString  string
	)

	// 表示直接生成的脚本，非修改
	if scriptId == 0 {
		for _, argument := range Arguments {
			k := argument["key"]
			v := argument["value"]
			ArgumentString = model.ElementPropName + k + model.HTTPArgument + model.HTTPArgumentEncode + model.ArgumentStart + v + model.ArgumentMid + k + model.ArgumentEnd
			stringArguments = append(stringArguments, ArgumentString)
		}
	} else {
		//修改脚本，argument 从数据库取出的
		for _, argument := range Arguments {
			for k, v := range argument {
				ArgumentString = model.ElementPropName + k + model.HTTPArgument + model.HTTPArgumentEncode + model.ArgumentStart + v + model.ArgumentMid + k + model.ArgumentEnd
				stringArguments = append(stringArguments, ArgumentString)
			}
		}
	}

	for _, str := range stringArguments {
		stringArgument = stringArgument + str + "\r\n"
	}
	return
}

//SetJSONExtractor set JSONExtractorgo
func (s *Service) SetJSONExtractor(jsonExt model.JSONExtractor) (JSONExtractor string, err error) {
	var (
		buff *template.Template
		b    = bytes.NewBuffer(nil)
	)
	buff, err = template.ParseFiles(s.c.Jmeter.JSONExtractorTmp)
	if err != nil {
		log.Info("open template failed! (%v)", err)
		return
	}
	buff = buff.Funcs(template.FuncMap{"unescaped": unescaped})
	buff.Execute(b, jsonExt)
	JSONExtractor = b.String()
	return
}

//SetJSONExtractors Set JSONExtractors
func (s *Service) SetJSONExtractors(script *model.Script) (JSONExtractors string, err error) {
	var (
		JSONExtractords []string
		jsonExt         model.JSONExtractor
		JSONExtractord  string
	)
	// 表示直接生成的脚本，非修改
	if script.ID == 0 {
		for _, OutputParam := range script.OutputParamsMap {
			jsonExt.JSONName = OutputParam["key"]
			jsonExt.JSONPath = OutputParam["value"]
			if jsonExt.JSONName != "" && jsonExt.JSONPath != "" {
				if JSONExtractord, err = s.SetJSONExtractor(jsonExt); err != nil {
					log.Error("SetJSONExtractor err :(%v) ", err)
					return
				}
				JSONExtractords = append(JSONExtractords, JSONExtractord)
			}
		}
	} else {
		//修改脚本，header 从数据库取出的
		for _, OutputParam := range script.OutputParamsMap {
			for k, v := range OutputParam {
				jsonExt.JSONName = k
				jsonExt.JSONPath = v
			}
			if jsonExt.JSONName != "" && jsonExt.JSONPath != "" {
				if JSONExtractord, err = s.SetJSONExtractor(jsonExt); err != nil {
					log.Error("SetJSONExtractor err :(%v) ", err)
					return
				}
				JSONExtractords = append(JSONExtractords, JSONExtractord)
			}
		}
	}
	JSONExtractors = SliceToString(JSONExtractords, " ")
	return
}

//SetMultiPartInfo Set MultiPartInfo
func SetMultiPartInfo(script *model.Script) (multiPartInfo string) {
	multiPartInfo = model.MultipartName + script.MultipartPath + script.MultipartFile + model.MultipartFilePath +
		script.MultipartPath + script.MultipartFile + model.MultipartFilePathd + script.MultipartParam +
		model.MultipartMimetype + script.MimeType + model.MultipartEnd
	return
}

//MapToString map to string
func MapToString(strMapArr []map[string]string) (headerString string) {
	for _, strMap := range strMapArr {
		k, keyOk := strMap["key"]
		v, valueOK := strMap["value"]
		if keyOk != valueOK {
			return
		}
		if keyOk {
			delete(strMap, "key")
			delete(strMap, "value")
			strMap[k] = v
		}
	}
	headers, _ := json.Marshal(strMapArr)
	headerString = string(headers)
	return
}

//StringToMap string to map
func StringToMap(headerString string) (mapHeaders []map[string]string) {
	if headerString != "[{\"\":\"\"}]" {
		if err := json.Unmarshal([]byte(headerString), &mapHeaders); err != nil {
			return
		}
	}
	return
}

func unescaped(x string) interface{} {
	return template.HTML(x)
}

//HostInfo add hostinfo
func HostInfo(domain string, apiHeader string) (hostInfo string) {
	mapHeaders := StringToMap(apiHeader)
	reg := regexp.MustCompile(`.com|.co`)
	strs := reg.FindAllString(domain, -1)
	if len(strs) != 0 {
		for _, str := range strs {
			if (strings.Contains(str, ".com") || strings.Contains(str, ".co")) && strings.Contains(apiHeader, domain) { //绑定至用户指定的ip
				for _, header := range mapHeaders {
					for host := range header {
						conf.Conf.Paas.HostInfo = domain + ":" + host
						hostInfo = conf.Conf.Paas.HostInfo
					}
				}
			} else if strings.Contains(str, ".com") || strings.Contains(str, ".co") { //默认绑定至172.22.22.222
				conf.Conf.Paas.HostInfo = domain + ":" + conf.Conf.Melloi.DefaultHost
				hostInfo = conf.Conf.Paas.HostInfo
			}
		}
	}
	return
}

//HostInfoList add hostinfo list
func HostInfoList(scripts []*model.Script) (hostInfo string) {
	for _, script := range scripts {
		mapHeaders := StringToMap(script.APIHeader)
		reg := regexp.MustCompile(`.com|.co`)
		strs := reg.FindAllString(script.Domain, -1)
		if len(strs) != 0 {
			for _, str := range strs {
				if (strings.Contains(str, ".com") || strings.Contains(str, ".co")) && strings.Contains(script.APIHeader, script.Domain) { //绑定至用户指定的ip
					for _, header := range mapHeaders {
						for host := range header {
							conf.Conf.Paas.HostInfo = script.Domain + ":" + host
							hostInfo = conf.Conf.Paas.HostInfo
						}
					}
				} else if strings.Contains(str, ".com") || strings.Contains(str, ".co") { //默认绑定至172.22.22.222
					hostInfo = hostInfo + script.Domain + ":" + conf.Conf.Melloi.DefaultHost + ","
				}
			}
		}
	}
	conf.Conf.Paas.HostInfo = hostInfo
	return
}

//HostInfoByUploadSc host info By UploadSc
func HostInfoByUploadSc(domain string) (hostInfo string) {
	domains := strings.Split(domain, ",")
	for _, dom := range domains {
		hostInfo = hostInfo + dom + ":" + conf.Conf.Melloi.DefaultHost + ","
	}
	conf.Conf.Paas.HostInfo = hostInfo
	return
}

// GroupOrder Group Order
func GroupOrder(isEdit bool, script *model.Script, scene *model.Scene, paramList *model.ParamList) (groupId int, runOrder int) {
	var (
		flag = 0
		//tmpIndex = 0
		u   *url.URL
		err error
	)
	if u, err = url.Parse(script.URL); err != nil || u == nil {
		log.Error("sorry, url is wrong, please check. error(%v)", err)
		err = ecode.MelloiUrlParseError
		return
	}
	if len(paramList.ParamList) > 0 {
		for _, paramList := range paramList.ParamList {
			tempList := strings.Split(u.RawQuery, "&")
			tempHeaderList := strings.Split(script.APIHeader, ",")
			tempArgumentList := strings.Split(script.ArgumentString, ",")
			tempDataList := strings.Split(script.Data, ",")
			if u.RawQuery == "" && (script.APIHeader == "[{\"\":\"\"}]" || script.APIHeader == "[]") && script.Data == "" && (script.ArgumentString == "[{\"\":\"\"}]" || script.ArgumentString == "[]") {
				if scene.SceneType == 1 || scene.SceneType == 0 {
					groupId = paramList.GroupID + 1
					runOrder = 1
					continue
				} else if scene.SceneType == 2 {
					groupId = paramList.GroupID
					runOrder = paramList.RunOrder + 1
					continue
				}
			}
			//url params check
			if u.RawQuery != "" {
				if groupId, runOrder, flag, err = OrderEach(isEdit, tempList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
					return
					//} else if flag == 2 {
					//	tmpIndex = index
					//	break
				}
			}
			//APIHeader params check
			if script.APIHeader != "[{\"\":\"\"}]" && script.APIHeader != "[]" {
				if groupId, runOrder, flag, err = OrderEach(isEdit, tempHeaderList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
					return
				}
			}
			//ArgumentString params check
			if script.ArgumentString != "[{\"\":\"\"}]" && script.ArgumentString != "[]" {
				if groupId, runOrder, flag, err = OrderEach(isEdit, tempArgumentList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
					return
				}
			}
			//body params check
			if script.Data != "" {
				if groupId, runOrder, flag, err = OrderEach(isEdit, tempDataList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
					return
				}
			}
		}
		//for i := tmpIndex; i < len(paramList.ParamList); i++ {
		//	if err = s.dao.UpdateRunOrder(paramList.ParamList[i].); err != nil {
		//		log.Error("s.dao.UpdateScript err :(%v)", err)
		//		return
		//	}
		//}
	} else if script.GroupID != 0 {
		groupId = script.GroupID
		runOrder = 1
	} else {
		groupId = 1
		runOrder = 1
	}
	return
}

// GroupOrderByEdit Group Order By Edit
func GroupOrderByEdit(isEdit bool, id int, script *model.Script, scene *model.Scene, paramList *model.ParamList) (groupId int, runOrder int) {
	var (
		flag = 0
		//tmpIndex = 0
		u   *url.URL
		err error
	)
	if u, err = url.Parse(script.URL); err != nil || u == nil {
		log.Error("sorry, url is wrong, please check. error(%v)", err)
		err = ecode.MelloiUrlParseError
		return
	}
	if len(paramList.ParamList) > 0 {
		for _, paramList := range paramList.ParamList {
			if paramList.ID != id {
				tempList := strings.Split(u.RawQuery, "&")
				tempHeaderList := strings.Split(script.APIHeader, ",")
				tempArgumentList := strings.Split(script.ArgumentString, ",")
				tempDataList := strings.Split(script.Data, ",")
				//url params check
				if u.RawQuery != "" {
					if groupId, runOrder, flag, err = OrderEach(isEdit, tempList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
						return
					}
				}
				//APIHeader params check
				if script.APIHeader != "[{\"\":\"\"}]" && script.APIHeader != "[]" {
					if groupId, runOrder, flag, err = OrderEach(isEdit, tempHeaderList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
						return
					}
				}
				//ArgumentString params check
				if script.ArgumentString != "[{\"\":\"\"}]" && script.ArgumentString != "[]" {
					if groupId, runOrder, flag, err = OrderEach(isEdit, tempArgumentList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
						return
					}
				}
				//body params check
				if script.Data != "" {
					if groupId, runOrder, flag, err = OrderEach(isEdit, tempDataList, paramList, scene.SceneType); err != nil || (runOrder != 1 && flag == 0) {
						return
					}
				}
			}
		}
	}
	return
}

// OrderEach Order Each
func OrderEach(isEdit bool, strList []string, paramList *model.Params, sceneType int) (groupId int, runOrder int, flag int, err error) {
	for _, paramPage := range strList {
		if strings.Contains(paramPage, "$") {
			tempParam := strings.Split(paramPage, "$")[1]
			param := RegularTrans(tempParam)
			//several params
			if strings.Contains(paramList.OutputParams, ",") {
				dbTempList := strings.Split(paramList.OutputParams, ",")
				for _, dbParam := range dbTempList {
					flag = 0
					finalParam := strings.Split(strings.Split(dbParam, "\":\"")[0], "\"")[1]
					if param == finalParam {
						groupId = paramList.GroupID
						runOrder = paramList.RunOrder + 1
						return
					}
					if !isEdit {
						if sceneType == 1 || sceneType == 0 {
							groupId = paramList.GroupID + 1
							runOrder = 1
						} else if sceneType == 2 {
							groupId = paramList.GroupID
							runOrder = paramList.RunOrder + 1
							flag = 1
						}
					}
				}
			} else { //single param
				dbParam := paramList.OutputParams
				finalParam := strings.Split(strings.Split(dbParam, "\":\"")[0], "\"")[1]
				if param == finalParam {
					groupId = paramList.GroupID
					runOrder = paramList.RunOrder + 1
					return
				}
				if !isEdit {
					if sceneType == 1 || sceneType == 0 {
						groupId = paramList.GroupID + 1
						runOrder = 1
					} else if sceneType == 2 {
						groupId = paramList.GroupID
						runOrder = paramList.RunOrder + 1
						flag = 1
					}
				}
			}
		} else {
			if !isEdit {
				if sceneType == 1 || sceneType == 0 {
					groupId = paramList.GroupID + 1
					runOrder = 1
				} else if sceneType == 2 {
					groupId = paramList.GroupID
					runOrder = paramList.RunOrder + 1
					flag = 1
				}
			}
		}
	}
	return
}

// RegularTrans Regular Trans
func RegularTrans(tempStr string) (targetStr string) {
	reg := regexp.MustCompile(`[\w]+`)
	str := reg.FindAllString(tempStr, -1)
	if len(str) != 0 {
		targetStr = str[0]
	}
	return targetStr
}
