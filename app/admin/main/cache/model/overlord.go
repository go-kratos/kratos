package model

// OverlordReq .
type OverlordReq struct {
	Name  string `json:"name" form:"name"`
	Zone  string `json:"zone" form:"zone"`
	Type  string `json:"type" form:"type"`
	Alias string `json:"alias" form:"alias"`
	Addr  string `json:"addr" form:"addr"`
	AppID string `json:"appid" form:"appid"`

	PN int `form:"pn"  default:"1"`
	PS int `form:"ps"  default:"20"`

	Cookie string `json:"-"`
}

// OverlordResp .
type OverlordResp struct {
	Names    []string           `json:"names,omitempty"`
	Addrs    []string           `json:"addrs,omitempty"`
	Cluster  *OverlordCluster   `json:"cluster,omitempty"`
	Clusters []*OverlordCluster `json:"clusters,omitempty"`
	Total    int64              `json:"total"`
	Nodes    []*OverlordNode    `json:"nodes,omitempty"`
	Apps     []*OverlordApp     `json:"apps,omitempty"`
	AppIDs   []string           `json:"appids,omitempty"`
}

// TableName gorm table name.
func (*OverlordCluster) TableName() string {
	return "overlord_cluster"
}

// OverlordCluster .
type OverlordCluster struct {
	ID               int64  `json:"id" gorm:"column:id"`
	Name             string `json:"name" gorm:"column:name"`                           // 集群名字
	Type             string `json:"type" gorm:"column:type"`                           // 缓存类型. (memcache,redis,redis-cluster)
	Zone             string `json:"zone" gorm:"column:zone"`                           // 机房
	HashMethod       string `json:"hash_method" gorm:"column:hash_method"`             // 哈希方法  默认sha1
	HashDistribution string `json:"hash_distribution" gorm:"column:hash_distribution"` // key分布策略 默认为ketama一致性hash
	HashTag          string `json:"hash_tag" gorm:"column:hashtag"`                    // key hash 标识
	ListenProto      string `json:"listen_proto" gorm:"column:listen_proto"`
	ListenAddr       string `json:"listen_addr" gorm:"column:listen_addr"`
	DailTimeout      int32  `json:"dail_timeout" gorm:"column:dial"`               // dial 超时
	ReadTimeout      int32  `json:"read_timeout" gorm:"column:read"`               // read 超时
	WriteTimeout     int32  `json:"write_timeout" gorm:"column:write"`             // write 超时
	NodeConn         int8   `json:"node_conn" gorm:"column:nodeconn"`              // 集群内节点连接数
	PingFailLimit    int32  `json:"ping_fail_limit" gorm:"column:ping_fail_limit"` // 节点失败检测次数
	PingAutoEject    bool   `json:"ping_auto_eject" gorm:"column:auto_eject"`      // 是否自动剔除节点

	Nodes []*OverlordNode `json:"nodes" gorm:"-"`
}

// TableName gorm table name.
func (*OverlordNode) TableName() string {
	return "overlord_node"
}

// OverlordNode .
type OverlordNode struct {
	Cid    int64  `json:"cid" gorm:"column:cid"`
	Alias  string `json:"alias" gorm:"column:alias"`
	Addr   string `json:"addr" gorm:"column:addr"`
	Weight int8   `json:"weight" gorm:"column:weight"`
}

// TableName gorm table name.
func (*OverlordApp) TableName() string {
	return "overlord_appid"
}

// OverlordApp .
type OverlordApp struct {
	ID     int64  `json:"-" gorm:"column:id"`
	TreeID int64  `json:"treeid" gorm:"column:tree_id"`
	AppID  string `json:"appid" gorm:"column:app_id"`
	Cid    int64  `json:"cid" gorm:"column:cid"`

	Cluster *OverlordCluster `json:"cluster,omitempty"`
}

// OverlordApiserver resp result of clusters.
type OverlordApiserver struct {
	Group    string `json:"group"`
	Clusters []struct {
		Name string `json:"name"`
		Type string `json:"cache_type"`
		// HashMethod       string   `json:"hash_method"`
		// HashDistribution string   `json:"hash_distribution"`
		// HashTag          string   `json:"hash_tag"`
		// DailTimeout      int32    `json:"dail_timeout"`
		// ReadTimeout      int32    `json:"read_timeout"`
		// WriteTimeout     int32    `json:"write_timeout"`
		// NodeConn         int8     `json:"node_connections"`
		// PingFailLimit    int32    `json:"ping_fail_limit"`
		// PingAutoEject    bool     `json:"ping_auto_eject"`
		FrontEndPort int `json:"front_end_port"`

		Instances []struct {
			IP     string `json:"ip"`
			Port   int    `json:"port"`
			Weight int8   `json:"weight"`
			Alias  string `json:"alias"`
			State  string `json:"state"`
			Role   string `json:"role"`
		} `json:"instances"`
	} `json:"clusters"`
}

// OverlordToml resp result of clusters.
type OverlordToml struct {
	Name             string   `toml:"name"`
	Type             string   `toml:"cache_type"`
	HashMethod       string   `toml:"hash_method"`
	HashDistribution string   `toml:"hash_distribution"`
	HashTag          string   `toml:"hash_tag"`
	DailTimeout      int32    `toml:"dail_timeout"`
	ReadTimeout      int32    `toml:"read_timeout"`
	WriteTimeout     int32    `toml:"write_timeout"`
	NodeConn         int8     `toml:"node_connections"`
	PingFailLimit    int32    `toml:"ping_fail_limit"`
	PingAutoEject    bool     `toml:"ping_auto_eject"`
	ListenProto      string   `toml:"listen_proto"`
	ListenAddr       string   `toml:"listen_addr"`
	Servers          []string `toml:"servers"`
}
