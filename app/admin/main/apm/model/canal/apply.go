package canal

import (
	"go-common/library/time"
)

//type and explanation
const (
	TypeApply = int8(iota)
	TypeReview
	ReviewReject
	ReviewSuccess
	ReviewFailed
)

//TypeMap struct
var (
	TypeMap = map[int8]string{
		TypeApply:     "申请",
		TypeReview:    "审核",
		ReviewReject:  "驳回",
		ReviewSuccess: "通过",
		ReviewFailed:  "失败",
	}
)

// TableName case tablename
func (*Apply) TableName() string {
	return "canal_apply"
}

//Apply apply model
type Apply struct {
	ID       int       `gorm:"column:id" json:"id"`
	Addr     string    `gorm:"column:addr" json:"addr"`
	Remark   string    `gorm:"column:remark" json:"remark"`
	Cluster  string    `gorm:"column:cluster" json:"project"`
	Leader   string    `gorm:"column:leader" json:"leader"`
	Comment  string    `gorm:"column:comment" json:"comment"`
	State    int8      `gorm:"column:state" json:"status"`
	Operator string    `gorm:"column:operator" json:"operator"`
	Ctime    time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime    time.Time `gorm:"column:mtime" json:"mtime"`
	ConfID   int       `gorm:"column:conf_id" json:"conf_id"`
}

//Config struct
type Config struct {
	Instance *Instance `json:"instance" toml:"instance"`
}

//Instance struct
type Instance struct {
	Caddr           string        `json:"caddr" toml:"addr"`
	User            string        `json:"user" toml:"user"`
	Password        string        `json:"password" toml:"password"`
	MonitorPeriod   string        `json:"monitor_period" toml:"monitor_period,omitempty"`
	ServerID        int64         `json:"server_id" toml:"server_id"`
	Flavor          string        `json:"flavor" toml:"flavor"`
	HeartbeatPeriod time.Duration `json:"heartbeat_period" toml:"heartbeat_period"`
	ReadTimeout     time.Duration `json:"read_timeout" toml:"read_timeout"`
	DB              []*DB         `json:"db" toml:"db"`
}

//DB struct
type DB struct {
	Schema  string   `json:"schema" toml:"schema"`
	Table   []*Table `json:"table" toml:"table"`
	Databus *Databus `json:"databus" toml:"databus"`
	Infoc   *Infoc   `json:"infoc" toml:"infoc"`
}

//Table struct
type Table struct {
	Name       string   `json:"name" toml:"name"`
	Primarykey []string `json:"primarykey" toml:"primarykey"`
	Omitfield  []string `json:"omitfield" toml:"omitfield"`
}

//Databus struct
type Databus struct {
	Key          string `json:"key" toml:"key"`
	Secret       string `json:"secret" toml:"secret"`
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
}

//Infoc struct
type Infoc struct {
	TaskID       string `json:"taskID" toml:"taskID"`
	Proto        string `json:"proto" toml:"proto"`
	Addr         string `json:"addr" toml:"addr"`
	ReporterAddr string `json:"reporterAddr" toml:"reporterAddr"`
}

//ConfigReq struct is
type ConfigReq struct {
	Addr          string `form:"addr" validate:"required"`
	User          string `form:"user"`
	Password      string `form:"password"`
	MonitorPeriod string `form:"monitor_period"`
	Databases     string `form:"databases" validate:"required"`
	Project       string `form:"project"`
	Leader        string `form:"leader"`
	Mark          string `form:"mark" validate:"required"`
}
