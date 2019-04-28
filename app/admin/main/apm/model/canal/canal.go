package canal

import (
	xtime "go-common/library/time"
)

// TableName case tablename
func (*Canal) TableName() string {
	return "master_info"
}

// Canal canal
type Canal struct {
	ID       int64      `gorm:"column:id" json:"id"`
	Addr     string     `gorm:"column:addr" json:"addr" form:"addr" validate:"required"`
	BinName  string     `gorm:"column:bin_name" json:"bin_name" form:"bin_name"`
	BinPos   int32      `gorm:"column:bin_pos" json:"bin_pos" form:"bin_pos"`
	Remark   string     `gorm:"column:remark" json:"remark" form:"remark"`
	Leader   string     `gorm:"column:leader" json:"leader" form:"leader"`
	Cluster  string     `gorm:"column:cluster" json:"project" form:"project"`
	CTime    xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime    xtime.Time `gorm:"column:mtime" json:"mtime"`
	IsDelete int        `gorm:"column:is_delete" json:"is_delete"`
}

//ScanReq canal scan req
type ScanReq struct {
	Addr string `form:"addr" validate:"required"`
}

//Results canalscan resp
type Results struct {
	ID       int64     `json:"id"`
	Addr     string    `json:"addr"`
	Cluster  string    `json:"project"`
	Leader   string    `json:"leader"`
	Document *Document `json:"document"`
}

//EditReq canal edit req
type EditReq struct {
	ID      int64  `form:"id" validate:"required"`
	BinName string `form:"bin_name"`
	BinPos  int32  `form:"bin_pos"`
	Remark  string `form:"remark"`
	Leader  string `form:"leader"`
	Project string `form:"project"`
}

//ListReq canallist req
type ListReq struct {
	Addr    string `form:"addr"`
	Project string `form:"project"`
	Status  int8   `form:"status"`
	Pn      int    `form:"pn" default:"1"`
	Ps      int    `form:"ps" default:"20"`
}

//Paper canallist resp
type Paper struct {
	Total int         `json:"total"`
	Pn    int         `json:"pn"`
	Ps    int         `json:"ps"`
	Items interface{} `json:"items"`
}

//Conf is
type Conf struct {
	ID      int64  `json:"id"`
	Comment string `json:"comment"`
}

//Document document
type Document struct {
	Instance struct {
		User          string `json:"user" toml:"user"`
		Password      string `json:"password" toml:"password"`
		MonitorPeriod string `json:"monitor_period" toml:"monitor_period"`
		ServerID      int64  `json:"server_id" toml:"server_id"`
		Db            []*struct {
			Schema string `json:"schema" toml:"schema"`
			Table  []*struct {
				Name       string   `json:"name" toml:"name"`
				Primarykey []string `json:"primarykey,omitempty" toml:"primarykey"`
				Omitfield  []string `json:"omitfield,omitempty" toml:"omitfield"`
			} `json:"table" toml:"table"`
			Databus *struct {
				Group        string `json:"group" toml:"group"`
				Topic        string `json:"topic" toml:"topic"`
				Action       string `json:"action" toml:"action"`
				Name         string `json:"name" toml:"name"`
				Proto        string `json:"proto" toml:"proto"`
				Addr         string `json:"addr" toml:"addr"`
				Idle         int    `json:"idle" toml:"idle"`
				Active       int    `json:"active" toml:"active"`
				DialTimeout  string `json:"dialTimeout" toml:"dialTimeout"`
				ReadTimeout  string `json:"readTimeout" toml:"readTimeout"`
				WriteTimeout string `json:"writeTimeout" toml:"writeTimeout"`
				IdleTimeout  string `json:"idleTimeout" toml:"idleTimeout"`
			} `json:"databus" toml:"databus"`
			Infoc *struct {
				TaskID       string `json:"taskID" toml:"taskID"`
				Proto        string `json:"proto" toml:"proto"`
				Addr         string `json:"addr" toml:"addr"`
				ReporterAddr string `json:"reporterAddr" toml:"reporterAddr"`
			} `json:"infoc" toml:"infoc"`
		} `json:"db"`
	} `json:"instance"`
}
