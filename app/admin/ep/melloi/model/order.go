package model

import (
	"time"

	"go-common/library/ecode"
)

// Order perf order model
type Order struct {
	ID                 int64     `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	Name               string    `json:"name" form:"name"`
	Broker             string    `json:"broker" form:"broker"`
	TestBackGround     string    `json:"test_background" form:"test_background" gorm:"column:test_background"`
	Type               int32     `json:"type" form:"type" gorm:"type"`
	TestType           int32     `json:"test_type" form:"test_type" gorm:"test_type"`
	TestTarget         string    `json:"test_target" gorm:"test_target"`
	APIList            string    `json:"api_list" gorm:"api_list"`
	APIDoc             string    `json:"api_doc" gorm:"api_doc"`
	LimitUser          string    `json:"limit_user" gorm:"limit_user"`
	LimitIP            string    `json:"limit_ip" gorm:"limit_ip"`
	LimitVisit         string    `json:"limit_visit" gorm:"limit_visit"`
	ServerConf         string    `json:"server_conf" gorm:"server_conf"`
	DependentComponent string    `json:"dependent_component" gorm:"dependent_component"`
	DependentBusiness  string    `json:"dependent_business" gorm:"dependent_business"`
	TestDataFrom       string    `json:"test_data_from" gorm:"test_data_from"`
	TestHost           string    `json:"test_host" gorm:"test_host"`
	MoniRedis          string    `json:"moni_redis" gorm:"moni_redis"`
	MoniMemcache       string    `json:"moni_memcache" gorm:"moni_memcache"`
	MoniDocker         string    `json:"moni_docker" gorm:"moni_docker"`
	MoniAPI            string    `json:"moni_api" gorm:"moni_api"`
	MoniMysql          string    `json:"moni_mysql" gorm:"moni_mysql"`
	MoniElasticsearch  string    `json:"moni_elasticsearch" gorm:"moni_elasticsearch"`
	MoniOther          string    `json:"moni_other" gorm:"moni_other"`
	TestCycles         string    `json:"test_cycles" gorm:"moni_cycles"`
	ScriptID           string    `json:"script_id" gorm:"script_id"`
	MachineID          string    `json:"machine_id" gorm:"machine_id"`
	Department         string    `json:"department" form:"department" gorm:"department"`
	Project            string    `json:"project" form:"project" gorm:"project"`
	App                string    `json:"app" form:"app" gorm:"app"`
	Status             int32     `json:"status" form:"status" gorm:"status"`
	UpdateBy           string    `json:"update_by" form:"update_by" gorm:"update_by"`
	Handler            string    `json:"handler" form:"handler" gorm:"handler"`
	ApplyDate          time.Time `json:"apply_date" gorm:"apply_date"`
	Active             int32     `json:"active" gorm:"active"`
}

// QueryOrderRequest queryOrderRequest
type QueryOrderRequest struct {
	Order
	Pagination
}

// QueryOrderResponse queryOrderResponse
type QueryOrderResponse struct {
	Orders []*Order `json:"orders"`
	Pagination
}

// Verify verify the value of pageNum and pageSize.
func (p *Pagination) Verify() error {
	if p.PageNum < 0 {
		return ecode.MeilloiIllegalPageNumErr
	} else if p.PageNum == 0 {
		p.PageNum = 1
	}
	if p.PageSize < 0 {
		return ecode.MeilloillegalPageSizeErr
	} else if p.PageSize == 0 {
		p.PageSize = 10
	}
	return nil
}

// Pagination page num
type Pagination struct {
	PageNum   int32 `form:"page_num" json:"page_num"`
	PageSize  int32 `form:"page_size" json:"page_size"`
	TotalSize int32 `form:"total_size" json:"total_size"`
}

// TableName get table name model
func (w Order) TableName() string {
	return "order"
}
