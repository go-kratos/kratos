package model

import "time"

// CommonResp ...
type CommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
}

// BuildNewFile ...
type BuildNewFile struct {
	ID        int    `json:"id"`
	AppID     int    `json:"app_id"`
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	From      int    `json:"from"`
	State     int    `json:"state"`
	Mark      string `json:"mark"`
	Operator  string `json:"operator"`
	IsDelete  int    `json:"is_delete"`
	NewCommon int    `json:"new_common"`
	Ctime     int    `json:"ctime"`
	Mtime     int    `json:"mtime"`
}

// BuildFile ...
type BuildFile struct {
	*BuildNewFile
	LastConf *BuildNewFile `json:"last_conf"`
}

// ConfigData ...
type ConfigData struct {
	Files        []*BuildFile    `json:"files"`
	BuildFiles   []*BuildFile    `json:"build_files"`
	BuildNewFile []*BuildNewFile `json:"build_new_file"`
}

// ConfigsParam ...
type ConfigsParam struct {
	AppName              string
	TreeID               int
	Env                  string
	Zone                 string
	BuildId              int
	Build                string
	Token                string
	FilenameGo           string
	FilenameRunnerJava   string
	FilenameTokenJava    string
	FilenameRunnerCommon string
	Increment            int
	Force                int
	AutoRequiredParams   []string
	RequiredParams       []string
	Comment              *ConfigComment
}

// SagaConfigsParam ...
type SagaConfigsParam struct {
	FileName  string
	AppName   string
	TreeID    int
	Env       string
	Zone      string
	BuildId   int
	Build     string
	Token     string
	Increment int
	Force     int
	UserList  []string
}

// ConfigComment ...
type ConfigComment struct {
	CommentURL string
}

// TagUpdate ...
type TagUpdate struct {
	Mark  string `form:"mark"`
	Names string `form:"names"`
}

// SvenResp ...
type SvenResp struct {
	CommonResp
	Data *ConfigData `json:"data"`
}

// ConfigValueResp ...
type ConfigValueResp struct {
	CommonResp
	Data *BuildNewFile `json:"data"`
}

// Config ...
type Config struct {
	Property *Property
}

// Property ...
type Property struct {
	Repos []*RepoConfig
}

// ConfigList ...
type ConfigList struct {
	ProjectID int              `json:"project_id" validate:"required"`
	Configs   []ConfigSagaItem `json:"configs"`
}

// ConfigSagaItem ...
type ConfigSagaItem struct {
	Name  string      `json:"name" validate:"required"`
	Value interface{} `json:"value"`
}

// SagaConfigLogResp ...
type SagaConfigLogResp struct {
	Id         int       `form:"id" gorm:"column:id"`
	Username   string    `form:"username" json:"username" gorm:"column:username"`
	ProjectId  int       `form:"project_id" json:"project_id" gorm:"column:project_id"`
	Content    string    `form:"content" json:"content" gorm:"column:content"`
	Ctime      time.Time `form:"ctime" json:"ctime" gorm:"column:ctime"`
	Mtime      time.Time `form:"mtime" json:"mtime" gorm:"column:mtime"`
	UpdateUser string    `form:"update_user" json:"update_user" gorm:"column:update_user"`
	Status     int       `form:"status" json:"status" gorm:"column:status"` //1创建 2修改 3同步中 4同步完成 5同步失败
}

// UpdateConfigReq ...
type UpdateConfigReq struct {
	Ids        []int  `json:"ids" validate:"required"`
	ConfigID   string `json:"config_id" validate:"required"`
	ConfigName string `json:"config_name" validate:"required"`
	Mark       string `json:"mark" validate:"required"`
}

// OptionSagaItem ...
type OptionSagaItem struct {
	ConfigSagaItem
	CNName  string `json:"cn_name"`
	Remark  string `json:"remark"`
	Type    string `json:"type"`
	Require bool   `json:"require"`
}

// RepoConfig ...
type RepoConfig struct {
	URL            string
	Group          string
	Name           string
	GName          string // gitlab仓库别名
	Language       string
	AuthBranches   []string // 鉴权分支
	TargetBranches []string // 分支白名单
	LockTimeout    int32
	MinReviewer    int
	RelatePipeline bool
	DelayMerge     bool
	LimitAuth      bool
	AllowLabel     string
	SuperAuthUsers []string
}
