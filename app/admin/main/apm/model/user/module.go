package user

import (
	"go-common/library/time"
)

// modules
const (
	USER        = "USER"
	ECODE       = "ECODE"
	DATABUS     = "DATABUS"
	DAPPER      = "DAPPER"
	CONFIG      = "CONFIG"
	CANAL       = "CANAL"
	DISCOVERY   = "DISCOVERY"
	APP         = "APP"
	PLATFORM    = "PLATFORM"
	CACHE       = "CACHE"
	OPEN        = "OPEN"
	NEED        = "NEED"
	PERFORMANCE = "PERFORMANCE"
	BFS         = "BFS"
)

// Modules modules
var (
	Modules = map[string]*Permission{
		USER:        {Name: "USER", Permit: PermitSuper, Des: "用户管理"},
		ECODE:       {Name: "ECODE", Permit: PermitDefault, Des: "错误码管理"},
		DATABUS:     {Name: "DATABUS", Permit: PermitDefault, Des: "DATABUS管理"},
		DAPPER:      {Name: "DAPPER", Permit: PermitDefault, Des: "DAPPER查询"},
		CONFIG:      {Name: "CONFIG", Permit: PermitDefault, Des: "配置中心"},
		CANAL:       {Name: "CANAL", Permit: PermitDefault, Des: "CANAL管理"},
		DISCOVERY:   {Name: "DISCOVERY", Permit: PermitDefault, Des: "DISCOVERY管理"},
		APP:         {Name: "APP", Permit: PermitDefault, Des: "APP管理"},
		PLATFORM:    {Name: "PLATFORM", Permit: PermitAuth, Des: "平台管理"},
		CACHE:       {Name: "CACHE", Permit: PermitDefault, Des: "缓存集群"},
		OPEN:        {Name: "OPEN", Permit: PermitDefault, Des: "open鉴权管理"},
		NEED:        {Name: "NEED", Permit: PermitDefault, Des: "需求管理"},
		PERFORMANCE: {Name: "PERFORMANCE", Permit: PermitDefault, Des: "性能管理"},
		BFS:         {Name: "BFS", Permit: PermitAuth, Des: "BFS管理"},
	}
)

// var (
// 	Modules = map[string]string{
// 		USER:      "用户管理",
// 		ECODE:     "错误码管理",
// 		DATABUS:   "DATABUS管理",
// 		DAPPER:    "DAPPER查询",
// 		CONFIG:    "配置中心",
// 		CANAL:     "CANAL管理",
// 		DISCOVERY: "DISCOVERY管理",
// 		APP:       "APP管理",
// 		PLATFORM:  "平台管理",
// 		CACHE:     "缓存集群",
// 		OPEN:      "open鉴权管理",
// 		NEED:      "需求管理",
// 	}
// 	DefaultModules = []string{ECODE, DATABUS, DAPPER, CONFIG, CANAL, DISCOVERY, APP, CACHE, OPEN}
// )

// TableName case tablename
func (*Module) TableName() string {
	return "user_module"
}

// Module module model
type Module struct {
	ID     int64     `gorm:"column:id" json:"id"`
	UserID int64     `gorm:"column:user_id" json:"user_id"`
	Module string    `gorm:"column:module" json:"module"`
	Ctime  time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime  time.Time `gorm:"column:mtime" json:"-"`
}
