package model

import (
	"time"
)

//Script script
type Script struct {
	ID                  int                 `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	TreeID              int                 `json:"tree_id"`
	ProjectID           int                 `json:"project_id" form:"project_id"`
	Type                int                 `json:"type" form:"type"`
	ProjectName         string              `json:"project_name" form:"project_name"`
	TestName            string              `json:"test_name" form:"test_name"`
	ThreadsSum          int                 `json:"threads_sum" form:"threads_sum"`
	LoadTime            int                 `json:"load_time"`
	ReadyTime           int                 `json:"ready_time"`
	ProcType            string              `json:"proc_type"`
	URL                 string              `json:"url" form:"url" gorm:"url"`
	Domain              string              `json:"domain" form:"domain"`
	Port                string              `json:"port"`
	Login               bool                `json:"login"`
	Path                string              `json:"path"`
	Method              string              `json:"method" form:"method"`
	Cookie              string              `json:"cookie" form:"cookie"`
	ContentType         string              `json:"content_type"`
	Data                string              `json:"data" form:"data"`
	Assertion           string              `json:"assertion"`
	AssertionString     interface{}         `json:"assertion_string" gorm:"-"`
	UseAssertion        bool                `json:"use_assertion" gorm:"-"`
	UseBuiltinParam     bool                `json:"use_builtin_param" gorm:"-"`
	SavePath            string              `json:"save_path" form:"save_path"`
	ResJtl              string              `json:"res_jtl" form:"res_jtl"`
	JmeterLog           string              `json:"jmeter_log"`
	UpdateBy            string              `json:"update_by" form:"update_by"`
	Ctime               time.Time           `json:"ctime" form:"ctime"`
	Mtime               time.Time           `json:"mtime" form:"mtime"`
	Active              int                 `json:"active"`
	Upload              bool                `json:"upload" form:"upload"`
	Headers             []map[string]string `json:"headers" form:"headers" gorm:"-"` // true
	APIHeader           string              `json:"api_header"`
	ArgumentsMap        []map[string]string `json:"arguments_map" gorm:"-"` // true
	ArgumentString      string              `gorm:"column:argument_map"`
	RowQuery            string              `json:"row_query" form:"row_query" gorm:"-"`
	UseSign             bool                `json:"use_sign" form:"use_sign"`
	LabelIds            []int               `json:"label_ids" form:"label_ids" gorm:"-"`
	IsCopy              bool                `json:"is_copy" form:"is_copy" gorm:"-"`
	ConnTimeOut         int                 `json:"conn_time_out"`
	RespTimeOut         int                 `json:"resp_time_out"`
	IsSave              bool                `json:"is_save" gorm:"-"`
	TestType            int                 `json:"test_type" form:"test_type"`
	SceneID             int                 `json:"scene_id" form:"scene_id"`
	OutputParamsMap     []map[string]string `json:"output_params_map" form:"output_params_map" gorm:"-"`
	OutputParams        string              `json:"output_params" form:"output_params"`
	JSONPath            string              `json:"json_path"`
	GroupID             int                 `json:"group_id" form:"group_id"`
	RunOrder            int                 `json:"run_order" form:"run_order"`
	ScriptPath          string              `json:"script_path" form:"script_path"`
	JmeterSample        interface{}         `json:"jmeter_sample" gorm:"-"`
	JSONExtractor       interface{}         `json:"json_extractor" gorm:"-"`
	IsAsync             bool                `json:"is_async" form:"is_async"`
	AsyncInfo           interface{}         `json:"async_info" gorm:"-"`
	MultiPartInfo       interface{}         `json:"multi_part_info" gorm:"-"`
	UseMultipart        bool                `json:"use_multipart" gorm:"-"`
	MultipartPath       string              `json:"multipart_path"`
	MultipartFile       string              `json:"multipart_file"`
	MultipartParam      string              `json:"multipart_param"`
	MimeType            string              `json:"mime_type"`
	Fusing              int                 `json:"fusing"`
	UseBusinessStop     bool                `json:"use_business_stop" form:"use_business_stop"`
	BusinessStopPercent int                 `json:"business_stop_percent" form:"business_stop_percent"`
	KeepAlive           bool                `json:"keep_alive" form:"keep_alive"`
	ExecuDockerSum      int                 `json:"execu_docker_sum" gorm:"-"`
	ConstTimer          int                 `json:"const_timer"`
	ConstTimerInfo      interface{}         `json:"const_timer_info" gorm:"-"`
	RandomTimer         int                 `json:"random_timer"`
	RandomTimerInfo     interface{}         `json:"random_timer_info" gorm:"-"`
	DataFile
	TreePath
}

//APIH api headers
type APIH struct {
	APIHeader []map[string]string `json:"api_header"`
}

//ScriptScene script scene
type ScriptScene struct {
	Scripts []Script `json:"scripts" form:"scripts"`
}

// TreePath service tree
type TreePath struct {
	Department string `json:"department" form:"department"`
	Project    string `json:"project" form:"project"`
	App        string `json:"app" form:"app"`
}

//QueryScriptResponse query script response
type QueryScriptResponse struct {
	Scripts []*ScriptLabels `json:"scripts"`
	Pagination
}

//ScriptLabels script labels
type ScriptLabels struct {
	Script
	Labels []*LabelRelation `json:"labels"`
}

//DataFile ignore db
type DataFile struct {
	UseDataFile   bool        `json:"use_data_file" gorm:"use_data_file"` // true
	FileName      string      `json:"file_name" gorm:"file_name"`         // true
	ParamsName    string      `json:"params_name"  gorm:"params_name"`    // true
	Delimiter     string      `json:"delimiter"  gorm:"delimiter"`        // true
	Loops         int         `json:"loops"     gorm:"loops"`             // true
	ResLog        string      `json:"res_log" gorm:"-"`
	BeginTestName string      `json:"begin_test_name" gorm:"-"`
	IsDebug       bool        `json:"is_debug" gorm:"-"`
	HeaderString  interface{} `json:"header_string" gorm:"-"`
	Arguments     interface{} `json:"arguments" gorm:"-"`
	FileSplit     bool        `json:"file_split" form:"file_split"`
	SplitNum      int         `json:"split_num" form:"split_num"`
}

//QueryScriptRequest query script request
type QueryScriptRequest struct {
	Script
	Pagination
	Executor string `json:"executor" form:"executor"`
}

//ScrThreadGroup script thread group
type ScrThreadGroup struct {
	Scripts []*Script `json:"scripts"`
}

//URLEncode URL Encode
type URLEncode struct {
	ParamsType string `json:"params_type"`
	NewUrl     string `json:"new_url"`
}

//TableName db table name of script
func (st Script) TableName() string {
	return "script"
}

// FusingInfo Fusing List
type FusingInfo struct {
	Fusing int `json:"fusing"`
}

// FusingInfoList Fusing Info List
type FusingInfoList struct {
	FusingList []FusingInfo `json:"fusing_list"`
	SetNull    bool         `json:"set_null"`
}
