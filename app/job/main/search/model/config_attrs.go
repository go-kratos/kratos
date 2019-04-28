package model

// SQLAttrs get attrs from db
type SQLAttrs struct {
	AppID            string
	DBName           string
	ESName           string
	DtbName          string
	TablePrefix      string
	TableFormat      string
	IndexAliasPrefix string //写和读数据时，指向的索引。是别名索引
	IndexVersion     string //创建索引时，指向的索引。这里是实体索引版本号
	IndexFormat      string
	IndexType        string
	IndexID          string
	IndexMapping     string
	DataIndexSuffix  string
	DataExtraInfo    string
	ReviewNum        int64
	ReviewTime       int64
	Sleep            float64
	Size             int
	Business         string
	DataFields       string
	SQLByID          string
	SQLByMTime       string
	SQLByIDMTime     string
	DatabusInfo      string
	DatabusIndexID   string
}

// Attrs parse AppAttrs
type Attrs struct {
	Business   string
	AppID      string
	DBName     string
	ESName     string
	DtbName    string
	Table      *AttrTable
	Index      *AttrIndex
	DataSQL    *AttrDataSQL
	DataExtras []AttrDataExtra //appID需要关联其他库的数据
	Databus    *AttrDatabus
	Other      *AttrOther
}

// AttrTable .
type AttrTable struct {
	TablePrefix string
	TableFormat string
	TableSplit  string
	TableFrom   int
	TableTo     int
	TableZero   string
	TableFixed  bool
}

// AttrIndex .
type AttrIndex struct {
	IndexAliasPrefix  string //写和读数据时，指向的索引。是别名索引
	IndexEntityPrefix string //创建索引时，指向的索引。是实体索引名
	IndexFormat       string
	IndexSplit        string
	IndexFrom         int
	IndexTo           int
	IndexType         string
	IndexID           string
	IndexMapping      string
	IndexZero         string
	IndexFixed        bool
}

// AttrDataSQL .
type AttrDataSQL struct {
	DataIndexSuffix       string //索引数据归属
	DataFields            string
	DataFieldsV2          map[string]AttrDataFields //存放json转换得到的data_fields字段信息, 替换老的DataFields
	DataIndexFields       []string                  //来自DataFields左数第一位
	DataIndexRemoveFields []string                  //ES需要移除的字段
	DataIndexFormatFields map[string]string         //ES每个字段的格式化，如int64,time,int等
	DataDtbFields         map[string][]string       //databus的字段对应的es字段, TODO 改成 map[string]map[string]bool 或 map[string][]string，应对一个数据库字段用在多个es字段
	DataExtraInfo         string
	SQLFields             string //来自DataFields左数第二位，含表名和字段的alias以及mysql函数等其他表达式
	SQLByID               string //因为有left join的缘故，顾提供完整sql（抛除字段部分，下同）
	SQLByMTime            string
	SQLByIDMTime          string
}

// AttrDataFields .
type AttrDataFields struct {
	ESField string `json:"es"`
	Field   string `json:"field"`
	SQL     string `json:"sql"`
	Expect  string `json:"expect"`
	Stored  string `json:"stored"`
	InDtb   string `json:"in_dtb"`
}

// AttrDataExtra .
type AttrDataExtra struct {
	Type         string            `json:"type"`
	Tag          string            `json:"tag"`
	Condition    map[string]string `json:"condition"`
	SliceField   string            `json:"slice_field"` // 逗号分隔，支持多个字段
	DBName       string            `json:"dbname"`
	Table        string            `json:"table"`
	TableFormat  string            `json:"table_format"`
	InField      string            `json:"in_field"`
	FieldsStr    string            `json:"fields_str"`
	Fields       []string          `json:"fields"`
	RemoveFields []string          `json:"remove_fields"`
	SQL          string            `json:"sql"`
}

// AttrDatabus .
type AttrDatabus struct {
	DatabusInfo string
	Ticker      int    // 定时时间(毫秒)
	AggCount    int    // 聚合数量
	Databus     string // databus Map key
	PrimaryID   string // 主表索引id
	RelatedID   string // 关联表索引id
}

// AttrOther .
type AttrOther struct {
	ReviewNum  int64
	ReviewTime int64
	Sleep      float64
	Size       int
}
