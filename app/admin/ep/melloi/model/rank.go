package model

//Rank rank
type Rank struct {
	Duration   int    `json:"duration" form:"duration"`
	TimeDegree string `json:"time_degree" form:"time_degree"`
	StartTime  string `json:"start_time" form:"start_time"`
	EndTime    string `json:"end_time" form:"end_time"`
	SearchAll  bool   `json:"search_all" form:"search_all"`
}

//Tree node in service tree
type Tree struct {
	Department string `json:"department" form:"department"`
	Project    string `json:"project" form:"project"`
	App        string `json:"app" form:"app"`
}

//TreeNum model for department, project and app performance test count statistic
type TreeNum struct {
	DeptNum int `json:"dept_num" form:"dept_num"`
	ProNum  int `json:"pro_num" form:"pro_num"`
	AppNum  int `json:"app_num" form:"app_num"`
}

//TreeList tree list
type TreeList struct {
	TreeList []*Tree `json:"tree_list"`
}

//NumList num list
type NumList struct {
	NumList TreeNum `json:"num"`
}

//API api
type API struct {
	URL   string `json:"url" form:"url"`
	Count int    `json:"count" form:"count"`
}

// GrpcInfo Grpc Info
type GrpcInfo struct {
	ServiceName   string `json:"service_name" form:"service_name"`
	RequestMethod string `json:"request_method" form:"request_method"`
	Count         int    `json:"count" form:"count"`
}

// GrpcRes Grpc Res
type GrpcRes struct {
	GrpcList []*GrpcInfo `json:"grpc_list" form:"grpc_list"`
}

// SceneCount Scene Count
type SceneCount struct {
	Department string `json:"department" form:"department"`
	SceneName  string `json:"scene_name" form:"scene_name"`
	Count      int    `json:"count" form:"count"`
}

// SceneRes SceneRes
type SceneRes struct {
	SceneList []*SceneCount `json:"scene_list" form:"scene_list"`
}

//Department department
type Department struct {
	Department string `json:"department" form:"department"`
	Count      int    `json:"count" form:"count"`
}

//Build performance test count by date
type Build struct {
	Date  string `json:"date" form:"date"`
	Count int    `json:"count" form:"count"`
}

//State model for test status
type State struct {
	TestStatus int `json:"test_status" form:"test_status"`
	Count      int `json:"count" form:"count"`
}

//TopAPIRes performance test top apis
type TopAPIRes struct {
	//Total    int      `json:"total"`
	APIList []*API `json:"api_list"`
}

//TopDeptRes performance test top departments
type TopDeptRes struct {
	DeptList []*Department `json:"dept_list"`
}

//BuildLineRes test line
type BuildLineRes struct {
	BuildList []*Build `json:"build_list"`
}

//StateLineRes test state chart
type StateLineRes struct {
	StateList []*State `json:"state_list"`
}
