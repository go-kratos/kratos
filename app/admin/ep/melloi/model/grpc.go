package model

import "go-common/library/time"

// GRPCQuickStartRequest GRPCModel grpc models
type GRPCQuickStartRequest struct {
	GRPCAddScriptRequest
}

// GRPCAddScriptRequest GRPC Add Script Request
type GRPCAddScriptRequest struct {
	TaskName       string      `json:"task_name"`
	Department     string      `json:"department"`
	Project        string      `json:"project"`
	APP            string      `json:"app"`
	Active         int         `json:"active"`
	HostName       string      `json:"host_name,"`
	Port           int         `json:"port"`
	ServiceName    string      `json:"service_name"`
	ProtoClassName string      `json:"proto_class_name"`
	PkgPath        string      `json:"pkg_path"`
	AsynCall       int         `json:"asyn_call"`
	RequestType    string      `json:"request_type"`
	RequestMethod  string      `json:"request_method"`
	RequestContent string      `json:"request_content"`
	ResponseType   string      `json:"response_type"`
	ScriptPath     string      `json:"script_path"`
	JarPath        string      `json:"jar_path"`
	ThreadsSum     int         `json:"threads_sum"`
	RampUp         int         `json:"ramp_up"`
	Loops          int         `json:"loops"`
	LoadTime       int         `json:"load_time"`
	UpdateBy       string      `json:"update_by"`
	Ctime          time.Time   `json:"ctime"`
	Mtime          time.Time   `json:"mtime"`
	IsDebug        int         `json:"is_debug" gorm:"-"` // 判断是否调试，不落库
	IsAsync        bool        `json:"is_async" form:"is_async"`
	AsyncInfo      interface{} `json:"async_info" gorm:"-"`
	ParamEnable    string      `json:"param_enable" `
	ParamDelimiter string      `json:"param_delimiter"`
	ParamFilePath  string      `json:"param_file_path" gorm:"param_file_path"`
	ParamNames     string      `json:"param_names"`
}

// GRPCUpdateScriptRequest GRPC Update Script Request
type GRPCUpdateScriptRequest struct {
	ID int `json:"id"`
	GRPCAddScriptRequest
}

// GRPCExecuteScriptRequest GRPC Execute Script Request
type GRPCExecuteScriptRequest struct {
	ScriptID int `json:"script_id"`
}

// TableName Table Name
func (g GRPC) TableName() string {
	return "grpc"
}

// TableName Table Name
func (g GRPCSnap) TableName() string {
	return "grpc_snap"
}

// QueryGRPCRequest  grpc of query request
type QueryGRPCRequest struct {
	Executor string `json:"executor" form:"executor"`
	GRPC
	Pagination
}

// QueryGRPCResponse grpc of query response
type QueryGRPCResponse struct {
	GRPCS []*GRPC `json:"grpcs"`
	Pagination
}

// GRPCReqToGRPC GRPC Req To GRPC
func GRPCReqToGRPC(grpcReq *GRPCAddScriptRequest) (grpc *GRPC) {
	grpc = &GRPC{}
	grpc.TaskName = grpcReq.TaskName
	grpc.Department = grpcReq.Department
	grpc.Project = grpcReq.Project
	grpc.APP = grpcReq.APP
	grpc.Active = grpcReq.Active
	grpc.HostName = grpcReq.HostName
	grpc.Port = grpcReq.Port
	grpc.ServiceName = grpcReq.ServiceName
	grpc.ProtoClassName = grpcReq.ProtoClassName
	grpc.PkgPath = grpcReq.PkgPath
	grpc.AsynCall = grpcReq.AsynCall
	grpc.RequestType = grpcReq.RequestType
	grpc.RequestMethod = grpcReq.RequestMethod
	grpc.RequestContent = grpcReq.RequestContent
	grpc.ResponseType = grpcReq.ResponseType
	grpc.ScriptPath = grpcReq.ScriptPath
	grpc.JarPath = grpcReq.JarPath
	grpc.ThreadsSum = grpcReq.ThreadsSum
	grpc.RampUp = grpcReq.RampUp
	grpc.Loops = grpcReq.Loops
	grpc.LoadTime = grpcReq.LoadTime
	grpc.UpdateBy = grpcReq.UpdateBy
	grpc.IsAsync = grpcReq.IsAsync
	grpc.AsyncInfo = grpcReq.AsyncInfo
	grpc.ParamEnable = grpcReq.ParamEnable
	grpc.ParamDelimiter = grpcReq.ParamDelimiter
	grpc.ParamFilePath = grpcReq.ParamFilePath
	grpc.ParamNames = grpcReq.ParamNames
	return
}
