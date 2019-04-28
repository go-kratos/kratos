package model

// LogParams .
type LogParams struct {
	Bsp       *BasicSearchParams
	Business  int    `form:"business" params:"business"`
	CTimeFrom string `form:"ctime_from" params:"ctime_from"`
	CTimeTo   string `form:"ctime_to" params:"ctime_to"`
}

// Business .
type Business struct {
	ID                int
	AppID             string
	Name              string
	AdditionalMapping string
	Mapping           map[string]string
	IndexFormat       string
	IndexCluster      string
	PermissionPoint   string
}

// UDepTsData .
type UDepTsData struct {
	Code int `json:"code"`
	Data map[string]string
}

// IPData .
type IPData struct {
	Code int `json:"code"`
	Data map[string]struct {
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
		Isp      string `json:"isp"`
	}
}

// LogAuditDefaultMapping .
var LogAuditDefaultMapping = map[string]string{
	"uname":      "string",
	"uid":        "string",
	"type":       "string",
	"oid":        "string",
	"action":     "string",
	"ctime":      "time",
	"int_0":      "int",
	"int_1":      "int",
	"int_2":      "int",
	"str_0":      "string",
	"str_1":      "string",
	"str_2":      "string",
	"extra_data": "string",
}

// LogUserActionDefaultMapping .
var LogUserActionDefaultMapping = map[string]string{
	"mid":        "string",
	"platform":   "string",
	"build":      "string",
	"buvid":      "string",
	"type":       "string",
	"oid":        "string",
	"action":     "string",
	"ip":         "string",
	"ctime":      "time",
	"int_0":      "int",
	"int_1":      "int",
	"int_2":      "int",
	"str_0":      "string",
	"str_1":      "string",
	"str_2":      "string",
	"extra_data": "string",
}
