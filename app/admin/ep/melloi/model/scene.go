package model

import (
	"time"
)

// Scene GRPCReqToGRPC
type Scene struct {
	ID             int         `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	SceneName      string      `json:"scene_name" form:"scene_name"`
	SceneType      int         `json:"scene_type" form:"scene_type"`
	UserName       string      `json:"user_name" form:"user_name"`
	IsDraft        int         `json:"is_draft" form:"is_draft"`
	IsDebug        bool        `json:"is_debug" form:"is_debug"`
	IsBatch        bool        `json:"is_batch" gorm:"-"`
	Scripts        []*Script   `json:"scripts" gorm:"-"`
	ThreadGroup    interface{} `json:"thread_group" gorm:"-"`
	ScriptPath     string      `json:"script_path" gorm:"-"`
	IsExecute      bool        `json:"is_execute" gorm:"-"`
	JmeterFilePath string      `json:"jmeter_file_path"`
	Department     string      `json:"department" form:"department"`
	Project        string      `json:"project" form:"project"`
	APP            string      `json:"app" form:"app"`
	Fusing         int         `json:"fusing" form:"fusing"`
	IsUpdate       bool        `json:"is_update" form:"is_update" gorm:"-"`
	JmeterLog      string      `json:"jmeter_log"`
	ResJtl         string      `json:"res_jtl"`
	IsActive       bool        `json:"is_active" form:"is_active"`
	Ctime          time.Time   `json:"ctime" form:"ctime"`
	Mtime          time.Time   `json:"mtime" form:"mtime"`
}

// Draft Draft
type Draft struct {
	SceneID   int    `json:"scene_id" form:"scene_id"`
	SceneName string `json:"scene_name" form:"scene_name"`
}

// QueryDraft QueryDraft
type QueryDraft struct {
	Total  int      `json:"total"`
	Drafts []*Draft `json:"draft_list"`
}

// Relation Relation
type Relation struct {
	GroupID int `json:"group_id"`
	Count   int `json:"count"`
}

// QueryRelation Query Relation
type QueryRelation struct {
	RelationList []*Relation `json:"relation_list"`
}

// QueryAPIs Query APIs
type QueryAPIs struct {
	Total      int        `json:"total"`
	SceneID    int        `json:"scene_id"`
	SceneName  string     `json:"scene_name"`
	SceneType  int        `json:"scene_type"`
	Department string     `json:"department"`
	Project    string     `json:"project"`
	App        string     `json:"app"`
	APIs       []*TestAPI `json:"api_list"`
}

// TestAPI Test API
type TestAPI struct {
	GroupID      int    `json:"group_id" form:"group_id"`
	RunOrder     int    `json:"run_order" form:"run_order"`
	ID           int    `json:"id" form:"id"`
	TestName     string `json:"test_name" form:"test_name"`
	URL          string `json:"url" form:"url"`
	OutputParams string `json:"output_params" form:"output_params"`
	ThreadsSum   string `json:"threads_sum" form:"threads_sum"`
	LoadTime     string `json:"load_time" form:"load_time"`
}

// ShowTree Show Tree
type ShowTree struct {
	IsShow int     `json:"is_show" form:"is_show"`
	Tree   []*Tree `json:"tree" form:"tree"`
}

// RunOrder Run Order
type RunOrder struct {
	SceneID   int `json:"scene_id" form:"scene_id"`
	SceneType int `json:"scene_type" form:"scene_type"`
}

// RunOrderList Run Order List
type RunOrderList struct {
	Total     int         `json:"total"`
	RunOrders []*RunOrder `json:"run_order_list"`
}

// Params Params
type Params struct {
	ID           int    `json:"id" form:"id"`
	GroupID      int    `json:"group_id" form:"group_id"`
	RunOrder     int    `json:"run_order" form:"run_order"`
	OutputParams string `json:"output_params" form:"output_params"`
}

// ParamList ParamList
type ParamList struct {
	ParamList []*Params `json:"param_list" form:"param_list"`
}

// SaveOrderReq Save Order Req
type SaveOrderReq struct {
	GroupOrderList []*GroupOrder `json:"group_order_list"`
}

// GroupOrder Group Order
type GroupOrder struct {
	ID       int    `json:"id"`
	TestName string `json:"test_name"`
	GroupID  int    `json:"group_id"`
	RunOrder int    `json:"run_order"`
}

// TableName Table Name
func (w Scene) TableName() string {
	return "scene"
}

//QuerySceneResponse query scene response
type QuerySceneResponse struct {
	Scenes []*Scene `json:"scenes"`
	Pagination
}

//QuerySceneRequest query script request
type QuerySceneRequest struct {
	Scene
	Pagination
	Executor string `json:"executor" form:"executor"`
}

// DoPtestSceneParam Do Ptest Scene Param
type DoPtestSceneParam struct {
	SceneID  int    `json:"scene_id" form:"scene_id"`
	UserName string `json:"user_name" form:"user_name"`
}

// DoPtestSceneParams Do Ptests Scene Param
type DoPtestSceneParams struct {
	SceneIDs []int  `json:"scene_ids" form:"scene_ids"`
	UserName string `json:"user_name" form:"user_name"`
}

// SceneInfo Scene Info
type SceneInfo struct {
	MaxLoadTime   int       `json:"max_load_time" form:"max_load_time"`
	TestNames     []string  `json:"test_names"`
	TestNameNicks []string  `json:"test_name_nicks"`
	JmeterLog     string    `json:"jmeter_log"`
	ResJtl        string    `json:"res_jtl"`
	LoadTimes     []int     `json:"load_times"`
	SceneName     string    `json:"scene_name"`
	Scripts       []*Script `json:"scripts"`
}

// APIInfo API Info
type APIInfo struct {
	ID         int    `json:"id" form:"id"`
	TestName   string `json:"test_name" form:"test_name"`
	URL        string `json:"url" form:"url"`
	ThreadsSum int    `json:"threads_sum" form:"threads_sum"`
	LoadTime   int    `json:"load_time" form:"load_time"`
}

// APIInfoList API Info List
type APIInfoList struct {
	SceneID int `json:"scene_id" form:"scene_id"`
	Pagination
	ScriptList []*Script `json:"script_list"`
}

// APIInfoRequest API Info Request
type APIInfoRequest struct {
	Script
	Pagination
	//DeliverySceneID int `json:"delivery_scene_id" form:"delivery_scene_id"`
}

// PreviewInfoList Preview Info List
type PreviewInfoList struct {
	PreviewInfoList []*PreviewInfo `json:"preview_info_list"`
}

// PreviewInfo Preview Info
type PreviewInfo struct {
	GroupInfo
	InfoList []*Preview `json:"info_list"`
}

// Preview Preview
type Preview struct {
	ID          int    `json:"id"`
	TestName    string `json:"test_name"`
	RunOrder    int    `json:"run_order"`
	GroupID     int    `json:"group_id"`
	ConstTimer  int    `json:"const_timer"`
	RandomTimer int    `json:"random_timer"`
}

// PreviewList Preview List
type PreviewList struct {
	PreList []*Preview `json:"pre_list"`
}

// GroupList Group List
type GroupList struct {
	GroupList []*GroupInfo `json:"group_list"`
}

// GroupInfo Group Info
type GroupInfo struct {
	GroupID    int `json:"group_id"`
	ThreadsSum int `json:"threads_sum"`
	LoadTime   int `json:"load_time"`
	ReadyTime  int `json:"ready_time"`
}

// UsefulParams Useful Params
type UsefulParams struct {
	OutputParams string `json:"output_params"`
}

// UsefulParamsList Useful Params
type UsefulParamsList struct {
	ParamsList []*UsefulParams `json:"params_list"`
}

// Test test
type Test struct {
	Count int `json:"count"`
}

// BindScene Bind Scene
type BindScene struct {
	SceneID int    `json:"scene_id"`
	ID      string `json:"id"`
}

// DrawRelationList Draw Relation List
type DrawRelationList struct {
	//Nodes []*string `json:"nodes"`
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}

// Node Node
type Node struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Edge Edge
type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}
