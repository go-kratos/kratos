package service

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/melloi/model"
	"go-common/app/admin/ep/melloi/service/proto"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_upCase   = 1
	_downCase = 0
)

// capitalize 首字母大小写处理
func (s *Service) capitalize(str string, tp int) string {
	b := []rune(str)
	if tp == _upCase {
		if b[0] >= 97 && b[1] <= 122 {
			b[0] -= 32
		}
	} else {
		if b[0] >= 64 && b[0] <= 90 {
			b[0] += 32
		}
	}
	return string(b)
}

// ProtoParsing analyze proto file
func (s *Service) ProtoParsing(protoPath, protoName string) (res map[string]interface{}, err error) {
	var pkgRetMap []string
	res = make(map[string]interface{})
	reader, err := os.Open(path.Join(protoPath, protoName))
	if err != nil {
		log.Error("open proto file error(%v)", err)
		return nil, err
	}
	defer reader.Close()
	parser := proto.NewParser(reader)

	var proto *proto.Proto
	proto, err = parser.Parse()
	if err != nil {
		log.Error("parse proto file error(%v)", err)
		return nil, err
	}

	res["fileName"] = protoName
	proto.Filename = strings.TrimSuffix(protoName, ".proto")
	res["protoClass"] = s.capitalize(proto.Filename, _upCase)

	if len(proto.Imports) > 0 {
		var retImportMap []string
		for _, impt := range proto.Imports {
			retImportMap = append(retImportMap, impt.Filename)
		}
		res["import"] = retImportMap
	}

	if len(proto.Package) > 0 {
		for _, pkg := range proto.Package {
			pkgRetMap = append(pkgRetMap, pkg.Name)
		}
		res["package"] = pkgRetMap
	}

	if len(proto.Options) > 0 {
		var retOptionsList []string
		for _, opt := range proto.Options {
			retOptionsList = append(retOptionsList, opt.Name)
			// java_package proto
			if opt.Name == "java_package" && len(proto.Package) <= 0 {
				pkgRetMap = append(pkgRetMap, opt.Constant.Source)
			}
		}
		res["package"] = pkgRetMap
		res["options"] = retOptionsList
	}

	if len(proto.Services) > 0 {
		var retImportMap []map[string]interface{}
		for _, srv := range proto.Services {
			firstTmp := make(map[string]interface{})
			firstTmp["name"] = srv.Name

			var secondMapTmp []map[string]interface{}
			for _, rpc := range srv.RPCElements {
				secondTmp := make(map[string]interface{})
				method := s.capitalize(rpc.Name, _downCase)
				// transfer get_by_uid to getByUid
				if strings.ContainsAny(method, "_") {
					var nmethod string
					methodStrs := strings.Split(method, "_")
					for i, item := range methodStrs {
						if i != 0 {
							item = s.capitalize(item, _upCase)
						}
						nmethod += item
					}
					method = nmethod
				}
				secondTmp["method"] = method
				secondTmp["requestType"] = rpc.RequestType
				secondTmp["returnType"] = rpc.ReturnsType

				secondMapTmp = append(secondMapTmp, secondTmp)
			}
			firstTmp["rpc"] = secondMapTmp
			retImportMap = append(retImportMap, firstTmp)
		}
		res["service"] = retImportMap
	}
	return
}

// CreateProtoImportDir create import dir
func (s *Service) CreateProtoImportDir(pathModel *model.ProtoPathModel) (err error) {
	cMakeDir := fmt.Sprintf("mkdir -p %s ", path.Join(pathModel.RootPath, pathModel.ExtraPath))
	if err = exec.Command("/bin/bash", "-c", cMakeDir).Run(); err != nil {
		log.Error("create proto import dir error(%v)", err)

	}
	return
}

//GRPCQuickStart grpc quickstart
func (s *Service) GRPCQuickStart(c context.Context, request *model.GRPCQuickStartRequest, runUser string, cookies string) (ret map[string]string, err error) {
	var (
		g        *model.GRPC
		reportID int
	)
	addScriptReq := request.GRPCAddScriptRequest
	ret = make(map[string]string)
	if g, err = s.GRPCAddScript(c, &addScriptReq); err != nil {
		log.Error("Save grpc script failed, (%v)", err)
		return
	}
	if reportID, err = s.GRPCRunByModel(c, g, runUser, cookies); err != nil {
		log.Error("performance test execution failed, (%v)", err)
		return
	}
	ret["report_id"] = strconv.Itoa(reportID)
	return
}

// SaveGRPCQuickStart save gRPC script of quick start
func (s *Service) SaveGRPCQuickStart(c context.Context, request *model.GRPCQuickStartRequest) (err error) {
	addScriptReq := request.GRPCAddScriptRequest
	_, err = s.GRPCAddScript(c, &addScriptReq)
	return err
}

//GRPCAddScript create grpc script
func (s *Service) GRPCAddScript(c context.Context, request *model.GRPCAddScriptRequest) (g *model.GRPC, err error) {
	// 若是debug ，则执行固定次数
	if request.IsDebug == 1 {
		request.Loops = 4
		request.ThreadsSum = 1
		request.TaskName += "_perf_debug" // debug tag
	}
	if g, err = s.CreateJmx(c, request); err != nil {
		log.Error("create jmeter file error (%v)", err)
		return
	}
	g.IsDebug = request.IsDebug
	return s.dao.CreateGRPC(g)
}

// CreateJmx create jmx
func (s *Service) CreateJmx(c context.Context, request *model.GRPCAddScriptRequest) (g *model.GRPC, err error) {
	if g, err = s.createJmeterFile(request); err != nil {
		log.Error("create jmeter file error (%v)", err)
		return
	}
	return
}

//GRPCRunByScriptID grpc execution by id
func (s *Service) GRPCRunByScriptID(c context.Context, request *model.GRPCExecuteScriptRequest, runUser string, cookie string) (reportId int, err error) {
	var grpc *model.GRPC
	if grpc, err = s.dao.QueryGRPCByID(request.ScriptID); err != nil {
		return
	}
	return s.GRPCRunByModel(c, grpc, runUser, cookie)
}

//GRPCRunByModel execute grpc by model data
func (s *Service) GRPCRunByModel(c context.Context, grpc *model.GRPC, runUser string, cookie string) (reportId int, err error) {
	var resp model.DoPtestResp
	tim := strconv.FormatInt(time.Now().Unix(), 10)
	testNameNick := grpc.TaskName + tim
	log.Info("开始调用压测grpc job-------\n")

	ptestParam := model.DoPtestParam{
		UserName:   runUser,                      // 用户名
		LoadTime:   grpc.LoadTime,                //运行时间
		TestNames:  StringToSlice(grpc.TaskName), //接口名转数组
		FileName:   grpc.JmxPath,                 // jmx文件
		ResJtl:     grpc.JtlLog,                  // jtl时间戳
		JmeterLog:  grpc.JmxLog,                  // jmeterlog时间戳
		Department: grpc.Department,
		Project:    grpc.Project,
		IsDebug:    false,
		APP:        grpc.APP,
		ScriptID:   grpc.ID,
		Cookie:     cookie,           // 用不到
		URL:        grpc.ServiceName, // 用于微信通知
		LabelIDs:   nil,
		Domain:     grpc.HostName, // 微信通知Domain
		FileSplit:  false,         // 文件切割
		SplitNum:   0,             // 切割数量
		JarPath:    grpc.JarPath,
		Type:       model.PROTOCOL_GRPC, //grpc
	}
	if grpc.IsDebug == 1 {
		ptestParam.IsDebug = true
	}
	if resp, err = s.DoPtestByJmeter(c, ptestParam, StringToSlice(testNameNick)); err != nil {
		log.Error("s.DoPtestByJmeter err :(%v)", err)
		return
	}
	reportId = resp.ReportSuID
	return
}

//QueryGrpc query grpc list
func (s *Service) QueryGrpc(c context.Context, sessionID string, qgr *model.QueryGRPCRequest) (res *model.QueryGRPCResponse, err error) {
	// 获取服务树节点
	var (
		treeNodes  []string
		treeNodesd []string
	)
	if treeNodesd, err = s.QueryUserRoleNode(c, sessionID); err != nil {
		log.Error("QueryUserRoleNode err (%v):", err)
		return
	}
	treeNodes = append(treeNodesd, "")
	if ExistsInSlice(qgr.Executor, s.c.Melloi.Executor) {
		return s.dao.QueryGRPCByWhiteName(&qgr.GRPC, qgr.PageNum, qgr.PageSize)
	}
	return s.dao.QueryGRPC(&qgr.GRPC, qgr.PageNum, qgr.PageSize, treeNodes)
}

//QueryGrpcById query grpc by id
func (s *Service) QueryGrpcById(id int) (*model.GRPC, error) {
	return s.dao.QueryGRPCByID(id)
}

// UpdateGrpc  update grpc
func (s *Service) UpdateGrpc(grpc *model.GRPC) error {
	return s.dao.UpdateGRPC(grpc)
}

// DeleteGrpc  delete grpc by id
func (s *Service) DeleteGrpc(id int) error {
	return s.dao.DeleteGRPC(id)
}

//createJmx create jmx file, jtl file, jmx_log file
func (s *Service) createJmeterFile(grpcReq *model.GRPCAddScriptRequest) (g *model.GRPC, err error) {
	var (
		buff *template.Template
		file *os.File
	)

	if buff, err = template.ParseFiles(s.c.Jmeter.GRPCTemplatePath); err != nil {
		log.Error("open grpc.jmx failed! (%v)", err)
		return nil, ecode.MelloiJmeterGenerateErr
	}

	g = model.GRPCReqToGRPC(grpcReq)
	// 脚本执行限制
	if g.LoadTime > s.c.Jmeter.TestTimeLimit {
		g.LoadTime = s.c.Jmeter.TestTimeLimit
	}
	if g.Loops <= 0 {
		g.Loops = -1
	}

	if g.IsAsync {
		g.AsyncInfo = unescaped(model.AsyncInfo)
	}

	// 生成jmx文件
	if grpcReq.ScriptPath == "" {
		log.Error("proto file is not uploaded.")
		return nil, ecode.MelloiProtoFileNotUploaded
	}

	g.JmxPath = path.Join(grpcReq.ScriptPath, grpcReq.TaskName+".jmx")
	g.JmxLog = path.Join(grpcReq.ScriptPath, grpcReq.TaskName+".log") // 定义好即可，运行时生成
	g.JtlLog = path.Join(grpcReq.ScriptPath, grpcReq.TaskName+".jtl") // 定义好即可，运行时生成

	if g.ParamFilePath != "" {
		g.ParamEnable = "true"
	}

	if file, err = os.Create(g.JmxPath); err != nil {
		log.Error("create jmx file error :(%v)", err)
		return nil, ecode.MelloiJmeterGenerateErr
	}
	defer file.Close()
	// 写入模板数据
	buff = buff.Funcs(template.FuncMap{"unescaped": unescaped})
	if err = buff.Execute(io.Writer(file), g); err != nil {
		log.Error("write jmeter file failed, (%v)", err)
		return nil, ecode.MelloiJmeterGenerateErr
	}
	return
}
