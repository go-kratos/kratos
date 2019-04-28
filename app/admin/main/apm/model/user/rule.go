package user

import (
	"go-common/library/time"
)

// rules
const (
	UserView  = "USER_VIEW"
	UserEdit  = "USER_EDIT"
	UserAudit = "USER_AUDIT"

	EcodeView = "ECODE_VIEW"
	EcodeEdit = "ECODE_EDIT"

	DatabusKeyView    = "DATABUS_KEY_VIEW"
	DatabusKeyEdit    = "DATABUS_KEY_EDIT"
	DatabusGroupView  = "DATABUS_GROUP_VIEW"
	DatabusGroupEdit  = "DATABUS_GROUP_EDIT"
	DatabusTopicView  = "DATABUS_TOPIC_VIEW"
	DatabusTopicEdit  = "DATABUS_TOPIC_EDIT"
	DatabusNotifyView = "DATABUS_NOTIFY_VIEW"
	DatabusNotifyEdit = "DATABUS_NOTIFY_EDIT"
	DatabusGroupApply = "DATABUS_GROUP_APPLY"

	DapperView = "DAPPER_VIEW"

	CanalView = "CANAL_VIEW"
	CanalEdit = "CANAL_EDIT"

	ConfigView       = "CONFIG_VIEW"
	ConfigSearchView = "CONFIG_SEARCH_VIEW"
	ConfigPublicView = "CONFIG_PUBLIC_VIEW"

	DiscoveryView      = "DISCOVERY_VIEW"
	PerformanceManager = "PERFORMANCE_MANAGER"

	AppView         = "APP_VIEW"
	AppEdit         = "APP_EDIT"
	AppAuthView     = "APP_AUTH_VIEW"
	AppCallerSearch = "APP_CALLER_SEARCH"

	NeedVerify = "NEED_VERIFY"

	PlatformSearchView = "PLATFORM_SEARCH_VIEW"
	PlatformReplyView  = "PLATFORM_REPLY_VIEW"
	PlatformTagView    = "PLATFORM_TAG_VIEW"

	CacheOpsView = "CACHE_OPS_VIEW"

	OpenView = "OPEN_VIEW"

	BFSView = "BFS_VIEW"
	BFSEdit = "BFS_EDIT"
)

//PermitType value
const (
	PermitDefault = iota
	PermitAuth
	PermitSuper
)

// PermitType permit type
type PermitType int

// Permission descript modules and rules
type Permission struct {
	Name   string
	Permit PermitType
	Des    string
}

// rules
var (
	Rules = map[string]*Permission{
		UserView:           {Name: "UserView", Permit: PermitSuper, Des: "用户查看"},
		UserEdit:           {Name: "UserEdit", Permit: PermitSuper, Des: "用户管理"},
		UserAudit:          {Name: "UserAudit", Permit: PermitSuper, Des: "权限审核"},
		EcodeView:          {Name: "EcodeView", Permit: PermitDefault, Des: "错误码查看"},
		EcodeEdit:          {Name: "EcodeEdit", Permit: PermitDefault, Des: "错误码编辑"},
		DatabusKeyView:     {Name: "DatabusKeyView", Permit: PermitDefault, Des: "Key查看"},
		DatabusKeyEdit:     {Name: "DatabusKeyEdit", Permit: PermitAuth, Des: "Key编辑"},
		DatabusGroupView:   {Name: "DatabusGroupView", Permit: PermitDefault, Des: "Group查看"},
		DatabusGroupEdit:   {Name: "DatabusGroupEdit", Permit: PermitAuth, Des: "Group修改"},
		DatabusTopicView:   {Name: "DatabusTopicView", Permit: PermitDefault, Des: "Topic查看"},
		DatabusTopicEdit:   {Name: "DatabusTopicEdit", Permit: PermitAuth, Des: "Topic编辑"},
		DatabusNotifyView:  {Name: "DatabusNotifyView", Permit: PermitDefault, Des: "Notify查看"},
		DatabusNotifyEdit:  {Name: "DatabusNotifyEdit", Permit: PermitAuth, Des: "Notify编辑"},
		DatabusGroupApply:  {Name: "DatabusGroupApply", Permit: PermitAuth, Des: "Group审核"},
		DapperView:         {Name: "DapperView", Permit: PermitDefault, Des: "Dapper查询"},
		CanalView:          {Name: "CanalView", Permit: PermitDefault, Des: "Canal查看"},
		CanalEdit:          {Name: "CanalEdit", Permit: PermitAuth, Des: "Canal编辑"},
		ConfigView:         {Name: "ConfigView", Permit: PermitDefault, Des: "配置列表查看"},
		ConfigSearchView:   {Name: "ConfigSearchView", Permit: PermitAuth, Des: "搜索列表查看"},
		ConfigPublicView:   {Name: "ConfigPublicView", Permit: PermitDefault, Des: "公共配置查看"},
		DiscoveryView:      {Name: "DiscoveryView", Permit: PermitDefault, Des: "Discovery查看"},
		AppView:            {Name: "AppView", Permit: PermitDefault, Des: "APP查看"},
		AppEdit:            {Name: "AppEdit", Permit: PermitAuth, Des: "APP编辑"},
		AppAuthView:        {Name: "AppAuthView", Permit: PermitAuth, Des: "APP鉴权查看"},
		AppCallerSearch:    {Name: "AppCallerSearch", Permit: PermitAuth, Des: "APP调用方查询"},
		PlatformSearchView: {Name: "PlatformSearchView", Permit: PermitAuth, Des: "平台搜索"},
		PlatformReplyView:  {Name: "PlatformReplyView", Permit: PermitAuth, Des: "平台评论"},
		PlatformTagView:    {Name: "PlatformTagView", Permit: PermitAuth, Des: "Tag"},
		CacheOpsView:       {Name: "CacheOpsView", Permit: PermitAuth, Des: "overlord缓存集群管理"},
		NeedVerify:         {Name: "NeedVerify", Permit: PermitSuper, Des: "需求建议审核"},
		OpenView:           {Name: "OpenView", Permit: PermitAuth, Des: "open鉴权查看"},
		PerformanceManager: {Name: "PerformanceManager", Permit: PermitDefault, Des: "性能管理"},
		BFSView:            {Name: "BFSView", Permit: PermitAuth, Des: "BFS查看"},
		BFSEdit:            {Name: "BFSEdit", Permit: PermitAuth, Des: "BFS编辑"},
	}
)

// TableName case tablename
func (*Rule) TableName() string {
	return "user_rule"
}

// TableName case tablename
func (*Apply) TableName() string {
	return "user_apply"
}

// Rule rule model
type Rule struct {
	ID     int64     `gorm:"column:id" json:"id"`
	UserID int64     `gorm:"column:user_id" json:"user_id"`
	Rule   string    `gorm:"column:rule" json:"rule"`
	Ctime  time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime  time.Time `gorm:"column:mtime" json:"-"`
}

// Apply user apply
type Apply struct {
	ID     int64     `gorm:"column:id" json:"id"`
	UserID int64     `gorm:"column:user_id" json:"user_id" params:"user_id"`
	Rules  string    `gorm:"column:rules" json:"rules" params:"rules"`
	Admin  string    `gorm:"column:admin" json:"admin" params:"admin"`
	Status int8      `gorm:"column:status" json:"status" `
	Ctime  time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime  time.Time `gorm:"column:mtime" json:"-"`
}

// Applies ...
type Applies struct {
	ID       int64  `gorm:"column:id" json:"id"`
	UserID   int64  `gorm:"column:user_id" json:"user_id"`
	UserName string `gorm:"column:username" json:"username"`
	Rules    string `gorm:"column:rules" json:"rules"`
	Status   string `gorm:"column:status" json:"status"`
}
