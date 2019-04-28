package service

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

// QueryDraft Query Draft
func (s *Service) QueryDraft(scene *model.Scene) (*model.QueryDraft, error) {
	return s.dao.QueryDraft(scene)
}

// UpdateScene Update Scene
func (s *Service) UpdateScene(scene *model.Scene) (fusing int, err error) {
	scriptIDList := make([]int, len(scene.Scripts))
	for index, script := range scene.Scripts {
		scriptIDList[index] = script.ID
	}
	fusing, err = s.dao.UpdateScene(scene, scriptIDList)
	return
}

// AddScene Add Scene
func (s *Service) AddScene(scene *model.Scene) (id int, err error) {
	id, err = s.dao.AddScene(scene)
	return
}

// QueryAPI Query API
func (s *Service) QueryAPI(scene *model.Scene) (*model.QueryAPIs, error) {
	return s.dao.QueryAPI(scene)
}

// AddConfig Add Config
func (s *Service) AddConfig(script *model.Script) error {
	return s.dao.AddConfig(script)
}

// QueryTree Query Tree
func (s *Service) QueryTree(script *model.Script) (*model.ShowTree, error) {
	return s.dao.QueryTree(script)
}

// QueryScenesByPage Query Scene By Page
func (s *Service) QueryScenesByPage(c context.Context, sessionID string, qsrq *model.QuerySceneRequest) (rsp *model.QuerySceneResponse, err error) {
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
		if rsp, err = s.dao.QueryScenesByPageWhiteName(&qsrq.Scene, qsrq.PageNum, qsrq.PageSize); err != nil {
			return
		}
	} else {
		if rsp, err = s.dao.QueryScenesByPage(&qsrq.Scene, qsrq.PageNum, qsrq.PageSize, treeNodes); err != nil {
			return
		}
	}
	return
}

//AddAndExecuScene AddSceneAuto add scene auto
func (s *Service) AddAndExecuScene(c context.Context, scene model.Scene, cookie string) (resp model.DoPtestResp, err error) {
	var (
		buff           *template.Template
		file           *os.File
		scrThreadGroup model.ScrThreadGroup
		fileName       string
		testNames      []string
		testNameNicks  []string
		loadTimes      []int
		scripts        []*model.Script
		jmeterLog      string
		resJtl         string
		sceneID        int
		threadGroup    string
	)
	// id 不是0 ，说明是输入接口创建的场景，根据sceneID 查询到接口列表，再根据接口列表生成jmx
	// 不是前端 quick-start 页的批量选择，走如下逻辑
	if scene.ID != 0 && !scene.IsBatch {
		script := model.Script{SceneID: scene.ID}
		if scripts, err = s.dao.QueryScripts(&script, 1, 200); err != nil {
			log.Error("s.dao.QueryScripts err :(%v)", err)
			return
		}
		scene.Scripts = scripts
	}
	if scene.IsDebug {
		for _, script := range scene.Scripts {
			script.ThreadsSum = 1
			script.Loops = 5
			script.TestName = script.TestName + "_perf_debug"
		}
		scene.SceneName = scene.SceneName + "_perf_debug"
	}
	// 非 quick-start 页面的批量选择，走如下逻辑，即 quick-start 页面的批量选择，不生成 scene.jmx 文件
	sceneSuffix := strconv.FormatInt(time.Now().Unix(), 10)
	if !scene.IsBatch {
		scrThreadGroup.Scripts = scene.Scripts
		if threadGroup, err = s.GetThreadGroup(scrThreadGroup); err != nil {
			log.Error("s.dao.GetThreadGroup (%v)", err)
			return
		}
		scene.ThreadGroup = unescaped(threadGroup)
		if buff, err = template.ParseFiles(s.c.Jmeter.JmeterSceneTmp); err != nil {
			log.Error("file is not exists (%v)", err)
			return
		}
		scene.ScriptPath = "/data/jmeter-log/" + scene.Department + "/" + scene.Project + "/" + scene.APP + "/" + "scene" + "/" + sceneSuffix + "/"
		if err = os.MkdirAll(scene.ScriptPath, 0755); err != nil {
			log.Error("Create SavePath Err (%v)", err)
			return
		}
		// 创建脚本文件 脚本路径+场景名+后缀.jmx
		if file, err = os.Create(scene.ScriptPath + sceneSuffix + ".jmx"); err != nil {
			log.Error("os.Create file err :(%v)", err)
			return
		}
		defer file.Close()
		buff = buff.Funcs(template.FuncMap{"unescaped": unescaped})
		buff.Execute(io.Writer(file), scene)
		fileName = file.Name()
		jmeterLog = scene.ScriptPath + "jmg" + sceneSuffix
		resJtl = scene.ScriptPath + "jtl" + sceneSuffix
	}
	scene.JmeterLog = jmeterLog
	scene.JmeterFilePath = fileName
	scene.ResJtl = resJtl
	sceneID = scene.ID

	// 批量选择脚本执行场景压测，则走如下逻辑，写 script 表
	if scene.ID == 0 || scene.IsBatch {
		//前端传入的scene 有 scripts , sceneId 没有传，走如下逻辑,写入scene 表
		if !scene.IsBatch {
			scene.SceneType = 1
			if sceneID, err = s.AddScene(&scene); err != nil {
				log.Error("s.AddScene err :(%v)", err)
				return
			}
			scene.ID = sceneID
		}
		for _, script := range scene.Scripts {
			script.IsSave = true
			script.SceneID = sceneID
			script.ID = 0
			script.TestType = model.SCENE_SCRIPT_TYPE
			//脚本落库
			if resp, err = s.AddAndExcuScript(c, script, "", &scene, false, true); err != nil {
				log.Error("s.AddAndExcuScript err :(%v)", err)
				return
			}
			script.ID = resp.ScriptID
			scripts = append(scripts, script)
		}
		scene.Scripts = scripts
	}

	// 有id，需要更新数据库
	if scene.ID != 0 && !scene.IsDebug {
		if _, err = s.UpdateScene(&scene); err != nil {
			log.Error("s.UpdateScene error :(%v)", err)
			return
		}
	}
	//获取 testNames  testNameNicks loadTimes
	for _, script := range scene.Scripts {
		testNames = append(testNames, script.TestName)
		testNameNicks = append(testNameNicks, script.TestName+sceneSuffix)
		loadTimes = append(loadTimes, script.LoadTime)
	}
	sort.Ints(loadTimes)
	//执行压测
	if scene.IsExecute && len(scene.Scripts) > 0 {
		ptestParam := model.DoPtestParam{
			TestNames:     testNames,
			SceneName:     scene.SceneName,
			UserName:      scene.UserName,
			LoadTime:      loadTimes[len(loadTimes)-1],
			FileName:      fileName,
			Upload:        false,
			ProjectName:   scene.SceneName,
			JmeterLog:     scene.JmeterLog,
			ResJtl:        scene.ResJtl,
			Department:    scene.Department,
			Project:       scene.Project,
			APP:           scene.APP,
			ScriptID:      0,
			DockerSum:     1,
			TestNameNick:  SliceToString(testNameNicks, ","),
			TestNameNicks: testNameNicks,
			Scripts:       scene.Scripts,
			SceneID:       sceneID,
			Type:          model.PROTOCOL_SCENE, // 场景
			IsDebug:       scene.IsDebug,
			Cookie:        cookie,
			Fusing:        scene.Fusing,
		}
		if resp, err = s.DoPtest(c, ptestParam); err != nil {
			log.Error("s.DoPtest err :(%v)", err)
			return
		}
		//临时写法
		resp.ScriptID = sceneID
	}

	return
}

// SaveScene Save Scene
func (s *Service) SaveScene(scene *model.Scene) error {
	return s.dao.SaveScene(scene)
}

// SaveOrder Save Order
func (s *Service) SaveOrder(req model.SaveOrderReq, scene *model.Scene) error {

	var (
		flag   = 1
		length = len(req.GroupOrderList)
		index  = length - 1
		bak    int
	)
	if scene.SceneType == 1 {
		for i := 0; i < length; i++ {
			bak = req.GroupOrderList[i].GroupID
			req.GroupOrderList[i].GroupID = flag
			//防止越界
			if i == index {
				break
			}
			if bak != req.GroupOrderList[i+1].GroupID {
				flag++
			}
		}
	} else {
		for i := 0; i < length; i++ {
			bak = req.GroupOrderList[i].RunOrder
			req.GroupOrderList[i].RunOrder = flag
			//防止越界
			if i == index {
				break
			}
			if bak != req.GroupOrderList[i+1].RunOrder {
				flag++
			}
		}
	}

	return s.dao.SaveOrder(req.GroupOrderList, scene)
}

// QueryRelation Query Relation
func (s *Service) QueryRelation(script *model.Script) (*model.QueryRelation, error) {
	var (
		groupId int
		err     error
	)
	if _, groupId, err = s.dao.QueryGroupId(script); err != nil {
		log.Error("s.dao.QueryGroupId err :(%v)", err)
		return nil, err
	}
	return s.dao.QueryRelation(groupId, script)
}

// DeleteAPI Delete API
func (s *Service) DeleteAPI(script *model.Script) error {
	return s.dao.DeleteAPI(script)
}

// DoScenePtest Do Scene Ptest
func (s *Service) DoScenePtest(c context.Context, ptestScene model.DoPtestSceneParam, addPtest bool, cookie string) (resp model.DoPtestResp, err error) {
	var (
		scenes  []*model.Scene
		scripts []*model.Script
	)

	scene := model.Scene{ID: ptestScene.SceneID}
	if scenes, err = s.dao.QueryScenes(&scene, 1, 1); err != nil {
		log.Error("s.dao.QueryScenes err :(%v)", err)
		return
	}
	script := model.Script{SceneID: ptestScene.SceneID}
	if scripts, err = s.dao.QueryScripts(&script, 1, 300); err != nil {
		log.Error("s.dao.QueryScripts err :(%v)", err)
		return
	}
	scenes[0].Scripts = scripts
	if len(scenes) > 0 {
		sceneInfo := GetSceneInfo(scenes[0])
		ptestParam := model.DoPtestParam{
			TestNames:     sceneInfo.TestNames,
			SceneName:     sceneInfo.SceneName,
			UserName:      ptestScene.UserName,
			LoadTime:      sceneInfo.MaxLoadTime,
			FileName:      scenes[0].JmeterFilePath,
			Upload:        false,
			ProjectName:   sceneInfo.SceneName,
			ResJtl:        sceneInfo.ResJtl,
			JmeterLog:     sceneInfo.JmeterLog,
			Department:    scenes[0].Department,
			Project:       scenes[0].Project,
			APP:           scenes[0].APP,
			ScriptID:      0,
			DockerSum:     1,
			TestNameNick:  SliceToString(sceneInfo.TestNameNicks, ","),
			TestNameNicks: sceneInfo.TestNameNicks,
			Scripts:       sceneInfo.Scripts,
			SceneID:       ptestScene.SceneID,
			Type:          model.PROTOCOL_SCENE, // 场景
			AddPtest:      addPtest,
			Cookie:        cookie,
		}
		return s.DoPtest(c, ptestParam)
	}
	return
}

//DoScenePtestBatch dosceneptest batch
func (s *Service) DoScenePtestBatch(c context.Context, ptestScenes model.DoPtestSceneParams, cookie string) (err error) {
	for _, SceneID := range ptestScenes.SceneIDs {
		ptestScene := model.DoPtestSceneParam{SceneID: SceneID, UserName: ptestScenes.UserName}
		go s.DoScenePtest(context.TODO(), ptestScene, false, cookie)
	}
	return
}

//GetSceneInfo get sceneInfo
func GetSceneInfo(scene *model.Scene) (sceneInfo model.SceneInfo) {
	var (
		loadTimes     []int
		testNames     []string
		testNameNicks []string
		scripts       []*model.Script
	)
	if scene == nil {
		return
	}
	sceneSuffix := strconv.FormatInt(time.Now().Unix(), 10)
	for _, script := range scene.Scripts {
		loadTimes = append(loadTimes, script.LoadTime)
		testNames = append(testNames, script.TestName)
		testNameNicks = append(testNameNicks, script.TestName+sceneSuffix)
		scripts = append(scripts, script)
	}
	sort.Ints(loadTimes)
	sceneInfo = model.SceneInfo{
		MaxLoadTime:   loadTimes[len(loadTimes)-1],
		JmeterLog:     scene.JmeterLog + sceneSuffix,
		ResJtl:        scene.ResJtl + sceneSuffix,
		LoadTimes:     loadTimes,
		TestNames:     testNames,
		TestNameNicks: testNameNicks,
		SceneName:     scene.SceneName,
		Scripts:       scripts,
	}
	return
}

// QueryExistAPI Query Exist API
func (s *Service) QueryExistAPI(c context.Context, sessionID string, req *model.APIInfoRequest) (res *model.APIInfoList, err error) {
	// 获取服务树节点
	var (
		treeNodes  []string
		treeNodesd []string
	)
	if treeNodesd, err = s.QueryUserRoleNode(c, sessionID); err != nil {
		log.Error("QueryUserRoleNode err (%v):", err)
	}
	treeNodes = append(treeNodesd, "")
	if res, err = s.dao.QueryExistAPI(&req.Script, req.PageNum, req.PageSize, req.SceneID, treeNodes); err != nil {
		return
	}
	for _, script := range res.ScriptList {
		if script.APIHeader != "" {
			if err = json.Unmarshal([]byte(script.APIHeader), &script.Headers); err != nil {
				log.Error("get script header err : (%v),scriptId:(%d)", err, script.ID)
			}
		}
		if script.ArgumentString != "" {
			if err = json.Unmarshal([]byte(script.ArgumentString), &script.ArgumentsMap); err != nil {
				log.Error("get script argument err: (%v), scriptId:(%d)", err, script.ID)
			}
		}
		if script.OutputParams != "" {
			if err = json.Unmarshal([]byte(script.OutputParams), &script.OutputParamsMap); err != nil {
				log.Error("get script OutputParams err: (%v),scriptId:(%d)", err, script.ID)
			}
		}
	}
	return
}

// QueryPreview Query Preview
func (s *Service) QueryPreview(req *model.Script) (preRes *model.PreviewInfoList, err error) {
	var (
		list    *model.GroupList
		preList *model.PreviewList
		preResd model.PreviewInfoList
	)

	if list, err = s.dao.QueryGroup(req.SceneID); err != nil {
		return
	}
	for i := 0; i < len(list.GroupList); i++ {
		// 或者使用var preInfo = &model.PreviewInfo{}
		preInfo := new(model.PreviewInfo)
		groupId := list.GroupList[i].GroupID
		threadsSum := list.GroupList[i].ThreadsSum
		loadTime := list.GroupList[i].LoadTime
		readyTime := list.GroupList[i].ReadyTime
		if preList, err = s.dao.QueryPreview(req.SceneID, groupId); err != nil {
			return
		}
		preInfo.GroupID = groupId
		preInfo.ThreadsSum = threadsSum
		preInfo.LoadTime = loadTime
		preInfo.ReadyTime = readyTime
		preInfo.InfoList = preList.PreList

		preResd.PreviewInfoList = append(preResd.PreviewInfoList, preInfo)
		preRes = &preResd
	}
	return
}

// QueryParams Query Params
func (s *Service) QueryParams(req *model.Script) (res *model.UsefulParamsList, tempRes *model.UsefulParamsList, err error) {
	//var paramList []string
	var (
		uParam *model.UsefulParams
	)
	res = &model.UsefulParamsList{}
	if tempRes, err = s.dao.QueryUsefulParams(req.SceneID); err != nil {
		return
	}
	if len(tempRes.ParamsList) > 0 {
		for _, tempParam := range tempRes.ParamsList {

			if strings.Contains(tempParam.OutputParams, ",") {
				tempParamList := strings.Split(tempParam.OutputParams, ",")
				for _, param := range tempParamList {
					uParam = &model.UsefulParams{}
					uParam.OutputParams = strings.Split(strings.Split(param, "\":\"")[0], "\"")[1]
					res.ParamsList = append(res.ParamsList, uParam)
				}
			} else {
				uParam = &model.UsefulParams{}
				uParam.OutputParams = strings.Split(strings.Split(tempParam.OutputParams, "\":\"")[0], "\"")[1]
				res.ParamsList = append(res.ParamsList, uParam)
			}
		}
	} else {
		// 赋值空数组
		res.ParamsList = []*model.UsefulParams{}
	}
	return
}

// UpdateBindScene Update Bind Scene
func (s *Service) UpdateBindScene(bindScene *model.BindScene) (err error) {
	return s.dao.UpdateBindScene(bindScene)
}

// QueryDrawRelation Query Draw Relation
func (s *Service) QueryDrawRelation(scene *model.Scene) (res model.DrawRelationList, tempRes *model.SaveOrderReq, err error) {
	var (
		//relationList model.DrawRelationList
		edges    []*model.Edge
		edge     *model.Edge
		tempEdge *model.Edge
		//nodes []*model.Node
		node *model.Node
	)

	if tempRes, err = s.dao.QueryDrawRelation(scene); err != nil {
		return
	}
	for k := 0; k < len(tempRes.GroupOrderList); k++ {
		node = &model.Node{}
		node.ID = tempRes.GroupOrderList[k].ID
		node.Name = tempRes.GroupOrderList[k].TestName
		res.Nodes = append(res.Nodes, node)
	}
	for i := 0; i < len(tempRes.GroupOrderList)-1; i++ {
		edge = &model.Edge{}
		tempEdge = &model.Edge{}
		if tempRes.GroupOrderList[i].GroupID == tempRes.GroupOrderList[i+1].GroupID && tempRes.GroupOrderList[i].RunOrder == tempRes.GroupOrderList[i+1].RunOrder {
			edge.Source = tempRes.GroupOrderList[i].TestName
			edges = append(edges, edge)
			if len(edges) > 0 {
				//edges[i+1].Source =
				tempEdge.Source = edges[i-1].Source
				tempEdge.Target = tempRes.GroupOrderList[i+1].TestName
				edges = append(edges, tempEdge)
			}
		} else if tempRes.GroupOrderList[i].GroupID == tempRes.GroupOrderList[i+1].GroupID && tempRes.GroupOrderList[i].RunOrder != tempRes.GroupOrderList[i+1].RunOrder {
			edge.Source = tempRes.GroupOrderList[i].TestName
			edge.Target = tempRes.GroupOrderList[i+1].TestName
			edges = append(edges, edge)
		} else if tempRes.GroupOrderList[i].GroupID != tempRes.GroupOrderList[i+1].GroupID {
			edge.Source = tempRes.GroupOrderList[i].TestName
			edges = append(edges, edge)
		}
	}
	for j := 0; j < len(edges); j++ {
		if edges[j].Target != "" {
			//relationList.Edges = append(relationList.Edges, edges[j])
			res.Edges = append(res.Edges, edges[j])
		}
	}
	return
}

// DeleteDraft Delete Draft
func (s *Service) DeleteDraft(scene *model.Scene) error {
	return s.dao.DeleteDraft(scene)
}

// QueryConfig Query Config
func (s *Service) QueryConfig(script *model.Script) (*model.GroupInfo, error) {
	return s.dao.QueryConfig(script)
}

// DeleteScene Delete Scene
func (s *Service) DeleteScene(scene *model.Scene) error {
	return s.dao.DeleteScene(scene)
}

//CopyScene copy scene
func (s *Service) CopyScene(c context.Context, scene *model.Scene, cookie string) (addScene model.AddScene, err error) {
	var (
		scripts []*model.Script
		scenes  []*model.Scene
		sceneID int
		resp    model.DoPtestResp
	)
	script := model.Script{SceneID: scene.ID}
	scened := model.Scene{ID: scene.ID}
	//先根据 scene.ID 查询 scenes 和 scripts
	if scenes, err = s.dao.QueryScenes(&scened, 1, 10); err != nil {
		log.Error("s.dao.QueryScenes err :(%v)", err)
		return
	}
	if scripts, err = s.dao.QueryScripts(&script, 1, 200); err != nil {
		log.Error("s.dao.QueryScripts err :(%v)", err)
		return
	}

	if len(scenes) != 0 {
		scenes[0].ID = 0
		scenes[0].UserName = scene.UserName
		scenes[0].SceneName = scene.SceneName
		//将新的scene 写入 数据库，并返回 sceneID
		if sceneID, err = s.dao.AddScene(scenes[0]); err != nil {
			log.Error("s.dao.AddScene err :(%v)", err)
			return
		}
		// 将新的 script 写入 script 表
		for _, sctd := range scripts {
			sctd.SceneID = sceneID
			sctd.ID = 0
			if _, _, _, err = s.dao.AddScript(sctd); err != nil {
				return
			}
		}
		sce := model.Scene{
			APP:        scenes[0].APP,
			Department: scenes[0].Department,
			Project:    scenes[0].Project,
			ID:         sceneID,
			IsDebug:    false,
			IsExecute:  false,
			UserName:   scene.UserName,
			SceneName:  scene.SceneName,
		}
		if resp, err = s.AddAndExecuScene(c, sce, cookie); err != nil {
			log.Error("s.AddAndExecuScene err :(%v), (%v)", err, resp)
			return
		}
		addScene.SceneID = sceneID
		addScene.UserName = scene.UserName
	}
	return
}

// QueryFusing Query Fusing
func (s *Service) QueryFusing(script *model.Script) (res *model.FusingInfoList, err error) {
	res, err = s.dao.QueryFusing(script)
	for i := 0; i < len(res.FusingList)-1; i++ {
		if res.FusingList[i] != res.FusingList[i+1] {
			res.SetNull = true
			return
		}
	}
	return
}
