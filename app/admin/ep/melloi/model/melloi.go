package model

//PaginateScripts page script
type PaginateScripts struct {
	Total   int       `json:"total"`
	Pn      int       `json:"page_num"`
	Ps      int       `json:"page_size"`
	Scripts []*Script `json:"scripts"`
}

//PaginateReports page report
type PaginateReports struct {
	Total          int              `json:"total"`
	Pn             int              `json:"page_num"`
	Ps             int              `json:"page_size"`
	ReportSummarys []*ReportSummary `json:"reportInfos"`
}

//ReducePtest model for test stress reduce
type ReducePtest struct {
	ID      int    `json:"id" form:"id"`
	JobName string `json:"job_name" form:"job_name"`
}

//PtestBatch ptest batch
type PtestBatch struct {
	UserName string `json:"user_name"`
	IDArr    []int  `json:"id_arr"`
}

//JobBatch   batch job
type JobBatch struct {
	JobNames    []string `json:"job_names"`
	ReportSuIDs []int    `json:"report_su_ids" form:"report_su_ids"`
}

//DockerStats model for container status
type DockerStats struct {
	Container string      `json:"container" form:"container"`
	Memory    interface{} `json:"memory" form:"memory"`
	CPU       string      `json:"cpu" form:"cpu"`
}

//DoPtestParam   ptest param
type DoPtestParam struct {
	UserName            string    `json:"user_name"`
	LoadTime            int       `json:"load_time"`
	TestNames           []string  `json:"test_names"` // 人工上传的脚本，可能会有很多接口名
	SceneName           string    `json:"scene_name"`
	TestNameNick        string    `json:"test_name_nick"`
	TestNameNicks       []string  `json:"test_name_nicks"`
	FileName            string    `json:"file_name"`
	Upload              bool      `json:"upload"`
	ProjectName         string    `json:"project_name"`
	ResLog              string    `json:"res_log"`
	ResJtl              string    `json:"res_jtl"`
	JmeterLog           string    `json:"jmeter_log"`
	Department          string    `json:"department"`
	Project             string    `json:"project"`
	APP                 string    `json:"app"`
	ScriptID            int       `json:"script_id"`
	AddPtest            bool      `json:"add_ptest"`
	IsDebug             bool      `json:"is_debug"`
	Cookie              string    `json:"cookie"`
	URL                 string    `json:"url"`
	Domain              string    `json:"domain"`
	LabelIDs            []int     `json:"label_ids"`
	FileSplit           bool      `json:"file_split"`
	SplitNum            int       `json:"split_num"`
	DockerSum           int       `json:"docker_sum"`
	JarPath             string    `json:"jar_path"`
	EnvInfo             string    `json:"env_info"`
	IsScene             bool      `json:"is_scene"` //场景压测
	Type                int       `json:"type"`     // 0.http单接口  1.场景报告  2.grpc报告  3.全链路
	Scripts             []*Script `json:"scripts"`
	SceneID             int       `json:"scene_id"`
	Fusing              int       `json:"fusing"`
	APIHeader           string    `json:"api_header"`
	ExecuDockerSum      int       `json:"execu_docker_sum"`
	UseBusinessStop     bool      `json:"use_business_stop"`
	BusinessStopPercent int       `json:"business_stop_percent"`
}

//QueryReGraphParam query ReGraphParam
type QueryReGraphParam struct {
	TestNameNicks []string `json:"test_name_nicks" form:"test_name_nicks"`
}

//UploadParam   uplaod param
type UploadParam struct {
	Path                string `json:"path" form:"path" params:"path"`
	IsPtest             bool   `json:"is_ptest" form:"is_ptest" params:"is_ptest"`
	UserName            string `json:"user_name" form:"user_name" params:"user_name"`
	TestName            string `json:"test_name" form:"test_name" params:"test_name"`
	Department          string `json:"department" form:"department" params:"department"`
	Project             string `json:"project" form:"project" params:"project"`
	APP                 string `json:"app" form:"app" params:"app"`
	ScriptPath          string `json:"script_path" form:"script_path" params:"script_path"`
	Domains             string `json:"domains" form:"domains" params:"domains"`
	Fusing              int    `json:"fusing" form:"fusing"`
	UseBusinessStop     bool   `json:"use_business_stop" form:"use_business_stop"`
	BusinessStopPercent int    `json:"business_stop_percent" form:"business_stop_percent"`
}

//QueryReportsRequest query report request
type QueryReportsRequest struct {
	ID           string `params:"id" form:"id" json:"id"`
	TestNameNick string `params:"test_name_nick" form:"test_name_nick" json:"test_name_nick"`
	TestName     string `params:"test_name" form:"test_name" json:"test_name"`
	Ps           int    `params:"page_size" form:"page_size" json:"page_size"`
	Pn           int    `params:"page_num" form:"page_num" json:"page_num"`
}

//BfsUploadParam bfs upload param
type BfsUploadParam struct {
	BfsIP        string `json:"bfs_ip" form:"bfs_ip" params:"bfs_ip"`
	BfsPort      int    `json:"bfs_port" form:"bfs_port" params:"bfs_port"`
	BucketName   string `json:"bucket_name" form:"bucket_name" params:"bucket_name"`
	FileName     string `json:"file_name" form:"file_name" params:"file_name"`
	AccessKey    string `json:"access_key" form:"access_key" params:"access_key"`
	AccessSecret string `json:"access_secret" form:"access_secret" params:"access_secret"`
	Method       string `json:"method" form:"method" params:"method"`
}

//JSONExtractor JSON Extractor
type JSONExtractor struct {
	JSONName string `json:"json_name"`
	JSONPath string `json:"json_path"`
}

//ReportGraphAdd Report Graph Add
type ReportGraphAdd struct {
	ReportSuID          int      `json:"report_su_id"`
	JobName             string   `json:"job_name"`
	TestName            string   `json:"test_name"`
	BeginTime           string   `json:"begin_time"`
	AfterTime           string   `json:"after_time"`
	TestNameNick        string   `json:"test_name_nick"`
	PodNames            []string `json:"pod_names"`
	ElapsedTime         int      `json:"elapsed_time"`
	Fusing              int      `json:"fusing"`
	UseBusinessStop     bool     `json:"use_business_stop"`
	BusinessStopPercent int      `json:"business_stop_percent"`
}

//AllPtestStop all ptest stop
type AllPtestStop struct {
	ReportSuID int `json:"report_su_id" form:"report_su_id"`
}
