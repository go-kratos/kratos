package metadata

// metadata common key
const (

	// Network

	RemoteIP   = "remote_ip"
	RemotePort = "remote_port"
	ServerAddr = "server_addr"
	ClientAddr = "client_addr"

	// Router

	Color = "color"

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

	// Mid
	// 外网账户用户id

	Mid = "mid"

	// Uid
	// 内网manager平台的用户id user_id

	Uid = "uid"

	// Username
	// LDAP平台的username

	Username = "username"

	// Device
	Device = "device"

	// Cluster cluster info key
	Cluster = "cluster"
)
