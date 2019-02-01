package metadata

// metadata common key
const (

	// Network
	RemoteIP   = "remote_ip"
	RemotePort = "remote_port"
	ServerAddr = "server_addr"
	ClientAddr = "client_addr"

	// Router
	Cluster = "cluster"
	Color   = "color"

	// Trace
	Trace  = "trace"
	Caller = "caller"

	// Timeout
	Timeout = "timeout"

	// Dispatch
	CPUUsage = "cpu_usage"
	Errors   = "errors"
	Requests = "requests"

	// Mirror
	Mirror = "mirror"

	// Mid 外网账户用户id
	Mid = "mid" // NOTE: ！！！业务可重新修改key名！！！

	// Username LDAP平台的username
	Username = "username"

	// Device 客户端信息
	Device = "device"
)
