package model

import (
	"time"
)

//ScriptSnap script snap
type ScriptSnap struct {
	ID             int                 `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	ScriptID       int                 `json:"script_id"`
	TreeID         int                 `json:"tree_id"`
	ProjectID      int                 `json:"project_id" form:"project_id"`
	ExecuteID      string              `json:"execute_id" form:"execute_id"`
	Type           int                 `json:"type" form:"type"`
	ProjectName    string              `json:"project_name" form:"project_name"`
	TestName       string              `json:"test_name" form:"test_name"`
	ThreadsSum     int                 `json:"threads_sum" form:"threads_sum"`
	LoadTime       int                 `json:"load_time"`
	ReadyTime      int                 `json:"ready_time"`
	ProcType       string              `json:"proc_type"`
	URL            string              `json:"url"`
	Domain         string              `json:"domain" form:"domain"`
	Port           string              `json:"port"`
	Login          bool                `json:"login"`
	Path           string              `json:"path"`
	Method         string              `json:"method" form:"method"`
	Cookie         string              `json:"cookie" form:"cookie"`
	ContentType    string              `json:"content_type"`
	Data           string              `json:"data" form:"data"`
	Assertion      string              `json:"assertion"`
	SavePath       string              `json:"save_path" form:"save_path"`
	ResJtl         string              `json:"res_jtl" form:"res_jtl"`
	JmeterLog      string              `json:"jmeter_log"`
	UpdateBy       string              `json:"update_by" form:"update_by"`
	Ctime          time.Time           `json:"ctime" form:"ctime"`
	Mtime          time.Time           `json:"mtime" form:"mtime"`
	Active         int                 `json:"active"`
	Upload         bool                `json:"upload" form:"upload"`
	Headers        []map[string]string `json:"headers" form:"headers" gorm:"-"` // true
	APIHeader      string              `json:"api_header" gorm:"column:api_header"`
	ArgumentsMap   []map[string]string `json:"arguments_map" gorm:"-"` // true
	ArgumentString string              `gorm:"column:argument_map"`
	ConnTimeOut    int                 `json:"conn_time_out"`
	RespTimeOut    int                 `json:"resp_time_out"`
	UseSign        bool                `json:"use_sign" form:"use_sign"`
	SceneID        int                 `json:"scene_id" form:"scene_id"`
	GroupID        int                 `json:"group_id" form:"group_id"`
	IsAsync        bool                `json:"is_async" form:"is_async"`
	MultipartPath  string              `json:"multipart_path"`
	MultipartFile  string              `json:"multipart_file"`
	MultipartParam string              `json:"multipart_param"`
	MimeType       string              `json:"mime_type"`
	Fusing         int                 `json:"fusing"`
	KeepAlive      bool                `json:"keep_alive" form:"keep_alive"`
	DataFile
	TreePath
}

//TableName script
func (st ScriptSnap) TableName() string {
	return "script_snap"
}
