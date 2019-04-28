package model

import (
	"encoding/json"

	"gopkg.in/olivere/elastic.v5"
)

// QueryParams .
type QueryParams struct {
	Business     string `form:"business" params:"business;Required" validate:"required"`
	QueryBodyStr string `form:"query" params:"query;Required" validate:"required"`
	DebugLevel   int    `form:"debug_level" params:"debug_level" default:"0"` // 2 默认全局debug，（包含dsl执行后的分析），1 dsl执行前的分析，防止504啥分析看不到。通过 /x/admin/search/query/debug （※包含es query体是否正确 + dsl体 + explain 返回信息） 方式，非签名请求
	QueryBody    *QueryBody
	AppIDConf    *QueryConfDetail
}

// QueryBody .
type QueryBody struct {
	Fields          []string            `json:"fields"` // default:"*" _source，default = *
	From            string              `json:"from"`   //索引名，多个用逗号隔开
	Where           *QueryBodyWhere     `json:"where"`
	Order           []map[string]string `json:"order"`
	OrderScoreFirst bool                `json:"order_score_first"`
	OrderRandomSeed string              `json:"order_random_seed"` // 随机排序种子
	Scroll          bool                `json:"scroll"`
	Highlight       bool                `json:"highlight"` //default:"false"
	Pn              int                 `json:"pn"`        //Range(1,5000) default:"1"
	Ps              int                 `json:"ps"`        //Range(1,1000) default:"10"
}

// QueryBodyWhere .
type QueryBodyWhere struct {
	EQ       map[string]interface{}     `json:"eq"`    //可能是数据或字符,[12,333,67] ["asd", "你好"]
	Or       map[string]interface{}     `json:"or"`    //暂时不支持minimum should
	In       map[string][]interface{}   `json:"in"`    //TODO改造为slice
	Range    map[string]string          `json:"range"` //[10,20)  (2018-05-10 00:00:00,2018-05-31 00:00:00]  (,30]
	Like     []QueryBodyWhereLike       `json:"like"`
	Enhanced []QueryBodyWhereEnhanced   `json:"enhanced"` //包含GourpBy Collapse
	Combo    []QueryBodyWhereCombo      `json:"combo"`    //混合与或
	Not      map[string]map[string]bool `json:"not"`      //对eq、in、range条件取反
}

// QueryBodyWhereLike .
type QueryBodyWhereLike struct {
	KWFields []string `json:"kw_fields"`
	KW       []string `json:"kw"`    //将kw的值使用空白间隔给query
	Or       bool     `json:"or"`    //default:"false"
	Level    string   `json:"level"` //默认default
}

// QueryBodyWhereEnhanced .
type QueryBodyWhereEnhanced struct {
	Mode  string              `json:"mode"`
	Field string              `json:"field"`
	Order []map[string]string `json:"order"`
	Size  int                 `json:"size"` //todo：sdk增加子集返回数
	// more conditions...
}

// QueryBodyWhereCombo .
type QueryBodyWhereCombo struct {
	EQ       []map[string]interface{}   `json:"eq"`
	In       []map[string][]interface{} `json:"in"`
	Range    []map[string]string        `json:"range"`
	NotEQ    []map[string]interface{}   `json:"not_eq"`
	NotIn    []map[string][]interface{} `json:"not_in"`
	NotRange []map[string]string        `json:"not_range"`
	Min      struct {
		EQ       int `json:"eq"`
		In       int `json:"in"`
		Range    int `json:"range"`
		NotEQ    int `json:"not_eq"`
		NotIn    int `json:"not_in"`
		NotRange int `json:"not_range"`
		Min      int `json:"min"`
	} `json:"min"`
}

// QueryConfDetail .
type QueryConfDetail struct {
	ESCluster     string
	IndexPrefix   string
	IndexType     string
	IndexID       string
	IndexMapping  string
	MaxIndicesNum int
	QueryMode     int //1:默认完全走查询体 2:基于查询体的定制 3:nested查询
	MaxPageSize   int //最大page size
}

// QueryResult query result.
type QueryResult struct {
	Order  string            `json:"order"`
	Sort   string            `json:"sort"`
	Result json.RawMessage   `json:"result"`
	Debug  *QueryDebugResult `json:"debug"`
	Page   *Page             `json:"page"`
}

// QueryDebugResult query result.
type QueryDebugResult struct {
	ErrMsg    []string               `json:"err_msg"`
	QueryBody string                 `json:"query_body"`
	DSL       string                 `json:"dsl"`
	Mapping   map[string]interface{} `json:"mapping"`
	Profile   *elastic.SearchProfile `json:"profile"` //性能分析
}

// AddErrMsg .
func (qdr *QueryDebugResult) AddErrMsg(msg ...string) {
	qdr.ErrMsg = append(qdr.ErrMsg, msg...)
}

// UpsertResult upsert result.
type UpsertResult struct {
}

var (
	// QueryModeBasic completely using basic query & nested .
	QueryModeBasic = 1
	// QueryModeExtra write some extra conditions under basic query .
	QueryModeExtra = 2

	// EnhancedModeGroupBy group by .
	EnhancedModeGroupBy = "group_by"
	// EnhancedModeSum sum from a filed .
	EnhancedModeSum = "sum"
	// EnhancedModeCollapse collapse .
	EnhancedModeCollapse = "collapse"
	// EnhancedModeDistinct distinct .
	EnhancedModeDistinct = "distinct"
	// EnhancedModeDistinctCount distinct .
	EnhancedModeDistinctCount = "distinct_count"
	// EnhancedModeGroupBySum group by sum .
	EnhancedModeGroupBySum = "group_by_sum"
	// EnhancedModeGroupByTop top hits .
	EnhancedModeGroupByTop = "group_by_tophits"

	// LikeLevelHigh high level .
	LikeLevelHigh = "high"
	// LikeLevelMiddel middle level .
	LikeLevelMiddel = "middle"
	// LikeLevelLow low level .
	LikeLevelLow = "low"

	// QueryConf 自定义部分
	QueryConf = map[string]*QueryConfDetail{
		"archive_video_score":     {ESCluster: "ssd_archive", IndexPrefix: "archive_video", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		"archive_score":           {ESCluster: "ssd_archive", IndexPrefix: "archive", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		"task_qa_random":          {ESCluster: "internalPublic", IndexPrefix: "task_qa", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		"esports_contests_date":   {ESCluster: "pcie_pub_out01", IndexPrefix: "esports_contests_map", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		"creative_archive_search": {ESCluster: "pcie_pub_out01", IndexPrefix: "creative_archive", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		"creative_archive_staff":  {ESCluster: "pcie_pub_out02", IndexPrefix: "creative_archive", IndexID: "%d,id", IndexType: "base", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		"creative_archive_apply":  {ESCluster: "pcie_pub_out02", IndexPrefix: "creative_archive", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		"dm_history":              {ESCluster: "dmout", IndexPrefix: "dm_search", MaxIndicesNum: 1, QueryMode: QueryModeExtra},
		// "pgc_contract_info":     {ESCluster: "pcie_pub_out01", IndexPrefix: "pgc_contract_info", MaxIndicesNum: 1, QueryMode: QueryModeNested},
		// "pgc_contract_video":    {ESCluster: "pcie_pub_out01", IndexPrefix: "pgc_contract_video", MaxIndicesNum: 1, QueryMode: QueryModeNested},
	}

	// PermConf 权限业务
	PermConf = map[string]map[string]string{
		"star":     {"ops_log_billions": "true"},                                                     // 业务使用*批量获取索引
		"scroll":   {"dm_search": "true"},                                                            // 业务使用scroll
		"oht":      {"creative_reply": "true", "creative_reply_isreport": "true", "esports": "true"}, // 业务max_result_window 100k
		"es_cache": {"comics_firebird": "true", "pgc_media": "true", "pgc_season": "true"},           // request cache(失效时间和索引的refresh_interval一致)
		//"routing":  {"creative_reply": "o_mid"},
	}
)
