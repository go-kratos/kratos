package model

import "go-common/library/time"

// GRPC grpc model
type GRPC struct {
	ID             int         `json:"id" form:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	TaskName       string      `json:"task_name" form:"task_name" gorm:"task_name"`
	Department     string      `json:"department" form:"department" gorm:"department"`
	Project        string      `json:"project" form:"project" gorm:"project"`
	APP            string      `json:"app" form:"app" gorm:"app"`
	Active         int         `json:"active" form:"active" gorm:"active"`
	HostName       string      `json:"host_name," form:"host_name" gorm:"host_name"`
	Port           int         `json:"port" form:"port" gorm:"port"`
	ServiceName    string      `json:"service_name" form:"service_name" gorm:"service_name"`
	ProtoClassName string      `json:"proto_class_name" form:"proto_class_name" gorm:"proto_class_name"`
	PkgPath        string      `json:"pkg_path" form:"pkg_path" gorm:"pkg_path"`
	AsynCall       int         `json:"asyn_call" form:"asyn_call" gorm:"asyn_call"`
	RequestType    string      `json:"request_type" form:"request_type" gorm:"request_type"`
	RequestMethod  string      `json:"request_method" form:"request_method" gorm:"request_method"`
	RequestContent string      `json:"request_content" form:"request_content" gorm:"request_content"`
	ResponseType   string      `json:"response_type" form:"response_type" gorm:"response_type"`
	ScriptPath     string      `json:"script_path" form:"script_path" gorm:"script_path"`
	JarPath        string      `json:"jar_path" form:"jar_path" gorm:"jar_path"`
	JmxPath        string      `json:"jmx_path" form:"jmx_path" gorm:"jmx_path"`
	JmxLog         string      `json:"jmx_log" form:"jmx_log" gorm:"jmx_log"`
	JtlLog         string      `json:"jtl_log" form:"jtl_log" gorm:"jtl_log"`
	ThreadsSum     int         `json:"threads_sum" form:"threads_sum" gorm:"threads_sum"`
	RampUp         int         `json:"ramp_up" form:"ramp_up" gorm:"ramp_up"`
	Loops          int         `json:"loops" form:"loop" gorm:"loops"`
	LoadTime       int         `json:"load_time" form:"load_time" gorm:"load_time"`
	UpdateBy       string      `json:"update_by" form:"update_by" gorm:"update_by"`
	Ctime          time.Time   `json:"ctime"`
	Mtime          time.Time   `json:"mtime"`
	IsDebug        int         `json:"is_debug" gorm:"-"`
	IsAsync        bool        `json:"is_async" form:"is_async" gorm:"is_async"`
	AsyncInfo      interface{} `json:"async_info" gorm:"-"`
	ParamEnable    string      `json:"param_enable" `
	ParamDelimiter string      `json:"param_delimiter"`
	ParamFilePath  string      `json:"param_file_path" gorm:"param_file_path"`
	ParamNames     string      `json:"param_names"`
}

// GRPCSnap grpc snap model
type GRPCSnap struct {
	ID             int    `json:"id" form:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	GRPCID         int    `json:"grpc_id" form:"grpc_id" gorm:"column:grpc_id"`
	TaskName       string `json:"task_name" form:"task_name" gorm:"task_name"`
	Department     string `json:"department" form:"department" gorm:"department"`
	Project        string `json:"project" form:"project" gorm:"project"`
	APP            string `json:"app" form:"app" gorm:"app"`
	Active         int    `json:"active" form:"active" gorm:"active"`
	HostName       string `json:"host_name," form:"host_name" gorm:"host_name"`
	Port           int    `json:"port" form:"port" gorm:"port"`
	ServiceName    string `json:"service_name" form:"service_name" gorm:"service_name"`
	ProtoClassName string `json:"proto_class_name" form:"proto_class_name" gorm:"proto_class_name"`
	PkgPath        string `json:"pkg_path" form:"pkg_path" gorm:"pkg_path"`
	AsynCall       int    `json:"asyn_call" form:"asyn_call" gorm:"asyn_call"`
	RequestType    string `json:"request_type" form:"request_type" gorm:"request_type"`
	RequestMethod  string `json:"request_method" form:"request_method" gorm:"request_method"`
	RequestContent string `json:"request_content" form:"request_content" gorm:"request_content"`
	ResponseType   string `json:"response_type" form:"response_type" gorm:"response_type"`
	ScriptPath     string `json:"script_path" form:"script_path" gorm:"script_path"`
	JarPath        string `json:"jar_path" form:"jar_path" gorm:"jar_path"`
	JmxPath        string `json:"jmx_path" form:"jmx_path" gorm:"jmx_path"`
	JmxLog         string `json:"jmx_log" form:"jmx_log" gorm:"jmx_log"`
	JtlLog         string `json:"jtl_log" form:"jtl_log" gorm:"jtl_log"`
	ThreadsSum     int    `json:"threads_sum" form:"threads_sum" gorm:"threads_sum"`
	RampUp         int    `json:"ramp_up" form:"ramp_up" gorm:"ramp_up"`
	Loops          int    `json:"loops" form:"loop" gorm:"loops"`
	LoadTime       int    `json:"load_time" form:"load_time" gorm:"load_time"`
	UpdateBy       string `json:"update_by" form:"update_by" gorm:"update_by"`
	ExecuteID      string `json:"execute_id" gorm:"execute_id"`
	IsAsync        bool   `json:"is_async" form:"is_async" gorm:"is_async"`
	ParamEnable    string `json:"param_enable" `
	ParamDelimiter string `json:"param_delimiter"`
	ParamFilePath  string `json:"param_file_path" gorm:"param_file_path"`
	ParamNames     string `json:"param_names"`
}

// ProtoPathModel create proto dependency path
type ProtoPathModel struct {
	RootPath  string `json:"root_path"`
	ExtraPath string `json:"extra_path"`
}

// DependResponse depend reponse
type DependResponse struct {
	Items []Item `json:"items"`
}

// Item serivce name
type Item struct {
	ServiceName string `json:"service_name"`
}
