package model

// ArgRegister define register param.
type ArgRegister struct {
	Region          string   `form:"region"`
	Zone            string   `form:"zone" validate:"required"`
	Env             string   `form:"env" validate:"required"`
	Appid           string   `form:"appid" validate:"required"`
	Treeid          int64    `form:"treeid"`
	Hostname        string   `form:"hostname" validate:"required"`
	Status          uint32   `form:"status" validate:"required"`
	HTTP            string   `form:"http"`
	RPC             string   `form:"rpc"`
	Version         string   `form:"version"`
	Metadata        string   `form:"metadata"`
	Replication     bool     `form:"replication"`
	Addrs           []string `form:"addrs,split"`
	LatestTimestamp int64    `form:"latest_timestamp"`
	DirtyTimestamp  int64    `form:"dirty_timestamp"`
}

// ArgRenew define renew params.
type ArgRenew struct {
	Region         string `form:"region"`
	Zone           string `form:"zone" validate:"required"`
	Env            string `form:"env" validate:"required"`
	Appid          string `form:"appid" validate:"required"`
	Treeid         int64  `form:"treeid"`
	Hostname       string `form:"hostname" validate:"required"`
	Replication    bool   `form:"replication"`
	DirtyTimestamp int64  `form:"dirty_timestamp"`
}

// ArgCancel define cancel params.
type ArgCancel struct {
	Region          string `form:"region"`
	Zone            string `form:"zone" validate:"required"`
	Env             string `form:"env" validate:"required"`
	Appid           string `form:"appid" validate:"required"`
	Treeid          int64  `form:"treeid"`
	Hostname        string `form:"hostname" validate:"required"`
	Replication     bool   `form:"replication"`
	LatestTimestamp int64  `form:"latest_timestamp"`
}

// ArgFetch define fetch param.
type ArgFetch struct {
	Region string `form:"region"`
	Zone   string `form:"zone"`
	Env    string `form:"env" validate:"required"`
	Appid  string `form:"appid"`
	Treeid int64  `form:"treeid"`
	Status uint32 `form:"status" validate:"required"`
}

// ArgFetchs define fetchs arg.
type ArgFetchs struct {
	Zone   string   `form:"zone"`
	Env    string   `form:"env" validate:"required"`
	Appid  []string `form:"appid,split"`
	Status uint32   `form:"status" validate:"required"`
}

// ArgPoll define poll param.
type ArgPoll struct {
	Region          string `form:"region"`
	Zone            string `form:"zone"`
	Env             string `form:"env" validate:"required"`
	Appid           string `form:"appid"`
	Treeid          int64  `form:"treeid"`
	Hostname        string `form:"hostname" validate:"required"`
	LatestTimestamp int64  `form:"latest_timestamp"`
}

// ArgPolling define polling arg.
type ArgPolling struct {
	Zone  string `form:"zone"`
	Env   string `form:"env" validate:"required"`
	Appid string `form:"appid"`
}

// ArgPolls define poll param.
type ArgPolls struct {
	Region          string   `form:"region"`
	Zone            string   `form:"zone"`
	Env             string   `form:"env" validate:"required"`
	Appid           []string `form:"appid,split"`
	Treeid          []int64  `form:"treeid,split"`
	Hostname        string   `form:"hostname,split" validate:"required"`
	LatestTimestamp []int64  `form:"latest_timestamp,split"`
}

// ArgSet define set param.
type ArgSet struct {
	Region       string   `form:"region"`
	Zone         string   `form:"zone" validate:"required"`
	Env          string   `form:"env" validate:"required"`
	Appid        string   `form:"appid" validate:"required"`
	Hostname     []string `form:"hostname,split"`
	Status       []int64  `form:"status,split"`
	Metadata     []string `form:"metadata"`
	Replication  bool     `form:"replication"`
	SetTimestamp int64    `form:"set_timestamp"`
}
