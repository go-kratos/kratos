package rpc

// ArgFilter rpc args.
type ArgFilter struct {
	Area      string
	Message   string
	TypeID    int16    //分区id 	非必要，默认0
	ID        int64    //内容id 	非必要，默认0
	OID       int64    //内容作用域id 	非必要，默认0
	MID       int64    //内容产生者mid 	非必要，默认0
	Keys      []string //key维度过滤 	非必要，默认0
	ReplyType int8     //评论稿件类型 非必要 默认0
}

type ArgMfilter struct {
	Area    string
	Message map[string]string
	TypeID  int16 // 可以不填
}

type FilterRes struct {
	Result string `json:"result"`
	Level  int8   `json:"level"`
	Limit  int    `json:"limit"`
}
