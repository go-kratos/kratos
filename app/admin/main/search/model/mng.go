package model

const (
	// MngAssetTypeDB .
	MngAssetTypeDB = 1
	// MngAssetTypeES .
	MngAssetTypeES = 2
	// MngAssetTypeDatabus .
	MngAssetTypeDatabus = 3
	// MngAssetTypeTable .
	MngAssetTypeTable = 4
)

// MngBusiness .
type MngBusiness struct {
	ID       int64             `json:"id"`
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Apps     []*MngBusinessApp `json:"apps"`
	AppsJSON string            `json:"-"`
}

// MngBusinessApp .
type MngBusinessApp struct {
	AppID    string `json:"appid"`
	IncrWay  string `json:"incr_way"`
	IncrOpen bool   `json:"incr_open"`
}

// MngAsset .
type MngAsset struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Config string `json:"config"`
	Desc   string `json:"desc"`
}

// MngAssetTable .
type MngAssetTable struct {
	TablePrefix string `json:"prefix"`
	TableFormat string `json:"format"`
}

// MngAssetDatabus .
type MngAssetDatabus struct {
	DatabusInfo    string `json:"info"`
	DatabusIndexID string `json:"index_id"`
}

// MngApp .
type MngApp struct {
	ID              int64   `json:"id" form:"id"`
	Business        string  `json:"business" form:"business"`
	AppID           string  `json:"appid" form:"appid"`
	Desc            string  `json:"desc" form:"desc"`
	DBName          string  `json:"db_name" form:"db_name"`
	ESName          string  `json:"es_name" form:"es_name"`
	TableName       string  `json:"table_name" form:"table_name"`
	TablePrefix     string  `json:"-"`
	TableFormat     string  `json:"-"`
	DatabusName     string  `json:"databus_name" form:"databus_name"`
	DatabusInfo     string  `json:"-"`
	DatabusIndexID  string  `json:"-"`
	IndexPrefix     string  `json:"index_prefix" form:"index_prefix"`
	IndexVersion    string  `json:"index_version" form:"index_version"`
	IndexFormat     string  `json:"index_format" form:"index_format"`
	IndexType       string  `json:"index_type" form:"index_type"`
	IndexID         string  `json:"index_id" form:"index_id"`
	DataIndexSuffix string  `json:"data_index_suffix" form:"data_index_suffix"`
	IndexMapping    string  `json:"index_mapping" form:"index_mapping"`
	DataFields      string  `json:"data_fields" form:"data_fields"`
	DataExtra       string  `json:"data_extra" form:"data_extra"`
	ReviewNum       int     `json:"review_num" form:"review_num"`
	ReviewTime      int     `json:"review_time" form:"review_time"`
	Sleep           float64 `json:"sleep" form:"sleep"`
	Size            int     `json:"size" form:"size"`
	SQLByID         string  `json:"sql_by_id" form:"sql_by_id"`
	SQLByMtime      string  `json:"sql_by_mtime" form:"sql_by_mtime"`
	SQLByIDMtime    string  `json:"sql_by_idmtime" form:"sql_by_idmtime"`
	QueryMaxIndexes int     `json:"query_max_indexes" form:"query_max_indexes"`
}

// MngCount .
type MngCount struct {
	Business string `json:"business" form:"business"`
	Type     string `json:"type" form:"type"`
	Name     string `json:"name"`
	Chart    string `json:"chart"`
	Param    string `json:"param"`
}

// MngCountRes .
type MngCountRes struct {
	Time  string `json:"time"`
	Count string `json:"count"`
}

// MngPercentRes .
type MngPercentRes struct {
	Name  string `json:"name"`
	Count string `json:"count"`
}

// UnamesData .
type UnamesData struct {
	Code int `json:"code"`
	Data map[string]string
}
