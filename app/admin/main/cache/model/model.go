package model

// ClustersReq params of cluster list.
type ClustersReq struct {
	PN int `form:"pn"  default:"1"`
	PS int `form:"ps"  default:"20"`
}

// ClusterResp resp result of clusters.
type ClusterResp struct {
	Clusters []*Cluster `json:"clusters"`
	Total    int64      `json:"total"` // 总数量
}

// TableName gorm table name.
func (*Cluster) TableName() string {
	return "cluster"
}

// Cluster resp result of clusters.
type Cluster struct {
	ID               int64    `gorm:"column:id" json:"id" toml:"id"`
	Name             string   `json:"name" gorm:"column:name" toml:"name"`       // 集群名字
	Type             string   `json:"type" gorm:"column:type" toml:"cache_type"` // 缓存类型. (memcache,redis,redis-cluster)
	AppID            string   `json:"app_id" gorm:"column:appids" toml:"app_id"`
	Zone             string   `json:"zone" gorm:"column:zone" toml:"zone"`                                        // 机房
	HashMethod       string   `json:"hash_method" gorm:"column:hash_method" toml:"hash_method"`                   // 哈希方法  默认sha1
	HashDistribution string   `json:"hash_distribution" gorm:"column:hash_distribution" toml:"hash_distribution"` // key分布策略 默认为ketama一致性hash
	HashTag          string   `json:"hash_tag" gorm:"column:hashtag" toml:"hash_tag"`                             // key hash 标识
	DailTimeout      int32    `json:"dail_timeout" gorm:"column:dial" toml:"dail_timeout"`                        // dial 超时
	ReadTimeout      int32    `json:"read_timeout" gorm:"column:read" toml:"read_timeout"`
	WriteTimeout     int32    `json:"write_timeout" gorm:"column:write" toml:"write_timeout"`               // read 超时
	NodeConn         int8     `json:"node_conn" gorm:"column:nodeconn" toml:"node_connections"`             // 集群内节点连接数
	PingFailLimit    int32    `json:"ping_fail_limit" gorm:"column:ping_fail_limit" toml:"ping_fail_limit"` // 节点失败检测次数
	PingAutoEject    bool     `json:"ping_auto_eject" gorm:"column:auto_eject" toml:"ping_auto_eject"`      // 是否自动剔除节点
	ListenProto      string   `json:"listen_proto" toml:"listen_proto"`
	ListenAddr       string   `json:"listen_addr" toml:"listen_addr"`
	Servers          []string `json:"-" gorm:"-" toml:"servers"`

	Hit      int `json:"hit" gorm:"-" toml:"-"`       // 集群命中率
	QPS      int `json:"qps" gorm:"-" toml:"-"`       // 集群qps
	State    int `json:"state" gorm:"-" toml:"-"`     // 集群状态 （0-online ;1-offline）
	MemRatio int `json:"mem_ratio" gorm:"-" toml:"-"` // 内存使用率

	Nodes []NodeDtl `json:"nodes" gorm:"-" toml:"-"`
}

// AddClusterReq params of add a cluster.
type AddClusterReq struct {
	ID               int64  `json:"id" form:"id"`                                                // 主键id 更新的时候使用
	Type             string `json:"type" form:"type" validate:"required"`                        // 缓存类型(memcache,redis,redis_cluster)
	AppID            string `json:"app_id" form:"app_id"`                                        // 集群关联的appid
	Zone             string `json:"zone" form:"zone"`                                            // 机房
	HashMethod       string `json:"hash_method" form:"hash_method" default:"fnv1a_64"`           // 哈希方法  默认fvn1a_64
	HashDistribution string `json:"hash_distribution" form:"hash_distribution" default:"ketama"` // key分布策略 默认为ketama一致性hash
	HashTag          string `json:"hash_tag" form:"hash_tag"`                                    // key hash 标识
	Name             string `json:"name" form:"name" validate:"required"`                        // 集群名字
	DailTimeout      int32  `json:"dail_timeout" form:"dail_timeout" default:"100"`              // dial 超时
	ReadTimeout      int32  `json:"read_timeout" form:"read_timeout" default:"100"`              // read 超时
	WriteTimeout     int32  `json:"write_timeout" form:"write_timeout" default:"100"`            // write 超时
	NodeConn         int8   `json:"node_conn" form:"node_conn" default:"10"`                     // 集群内及诶单连接数
	PingFailLimit    int32  `json:"ping_fail_limit" form:"ping_fail_limit"`                      // 节点失败检测次数
	PingAutoEject    bool   `json:"ping_auto_eject" form:"ping_auto_eject"`                      // 是否自动剔除节点
	ListenProto      string `json:"listen_proto" form:"listen_proto" default:"tcp"`              // 协议
	ListenAddr       string `json:"listen_addr" form:"listen_addr"`                              // 监听地址
}

// DelClusterReq params of del cluster.
type DelClusterReq struct {
	ID int64 `form:"id"` // 集群主键id
}

// ClusterReq get cluster by appid or cluster name.
type ClusterReq struct {
	AppID string `json:"app_id" form:"app_id"` // 关联的appid
	Zone  string `json:"zone" form:"zone"`     // 机房信息
	Type  string `json:"type" form:"type"`     // 缓存类型
	PN    int    `form:"pn"  default:"1"`
	PS    int    `form:"ps"  default:"20"`
	//	Cluster string `json:"cluster" form:"cluster"` // 集群名字
	Cookie string `form:"-"`
}

// ModifyClusterReq params of modify cluster detail.
type ModifyClusterReq struct {
	Name   string `json:"name" form:"name"`
	ID     int64  `json:"id" form:"id"`         // 集群id
	Action int8   `json:"action" form:"action"` // 操作(1 添加节点，2 删除节点；删除节点时只需要传alias)
	Nodes  string `json:"nodes" form:"nodes"`   // 节点信息 json数组 [{"id":11,"addr":"11","weight":1,"alias":"alias"}]
	// Addrs  []string `json:"addrs" form:"addrs,split"`   // 节点地址
	// Weight []int8   `json:"weight" form:"weight,split"` // 节点权重，必须与地址一一对应
	// Alias  []string `json:"alias" form:"alias,split"`   // 节点别名，必须与地址一一对应
}

// ClusterDtlReq params of get cluster detail.
type ClusterDtlReq struct {
	ID int64 `json:"id" form:"id"` // 集群id
}

// ClusterDtlResp resp result of cluster detail.
type ClusterDtlResp struct {
	Nodes []NodeDtl `json:"nodes"`
}

// ClusterFromYml get cluster from tw yml.
type ClusterFromYml struct {
	AppID string `json:"app_id" form:"app_id"` // 关联的appid
	Zone  string `json:"zone" form:"zone"`     // 机房信息
	TwYml string `json:"tw_yml" form:"tw_yml"`
}

// TableName gorm table name.
func (*NodeDtl) TableName() string {
	return "nodes"
}

// NodeDtl cluster node detaiwl
type NodeDtl struct {
	ID      int64   `json:"id" gorm:"column:id"`
	Cid     int64   `json:"cid" gorm:"column:cid"`
	Addr    string  `json:"addr" gorm:"column:addr"`
	Weight  int8    `json:"weight" gorm:"column:weight"`
	Alias   string  `json:"alias" gorm:"column:alias"`
	State   int8    `json:"state" gorm:"column:state"`
	QPS     int64   `json:"qps" gorm:"-"`
	MemUse  float32 `json:"mem_use" gorm:"-"`
	MemToal float32 `json:"mem_toal" gorm:"-"`
}

// EmpResp is empty resp.
type EmpResp struct {
}
