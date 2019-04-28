package model

const (
	// UserSearchDB bili_search.
	UserSearchDB = "bili_search"
	// DBDsnFormat .
	DBDsnFormat = "%s:%s@tcp(%s:%s)/%s?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
)

// GFAsset .
type GFAsset struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	DSN         string `json:"dsn"`
	DB          string `json:"db"`
	Regex       string `json:"relex"`
	Fields      string `json:"fields"`
	Description string `json:"description"`
	State       int8   `json:"state"`
}

// GFBusiness .
type GFBusiness struct {
	ID           int64  `json:"id"`
	PID          int64  `json:"pid"`
	Name         string `json:"name"`
	DataConf     string `json:"data_conf"`
	IndexConf    string `json:"index_conf"`
	BusinessConf string `json:"business_conf"`
	Description  string `json:"description"`
	State        int8   `json:"state"`
	Mtime        string `json:"mtime"`
}

// TableField .
type TableField struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Count   int    `json:"count"`
	Primary bool   `json:"primary"`
}
