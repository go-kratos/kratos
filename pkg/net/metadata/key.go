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

	// Device 客户端信息
	Device = "device"

	// Criticality 重要性
	Criticality = "criticality"
)

var outgoingKey = map[string]struct{}{
	Color:       {},
	RemoteIP:    {},
	RemotePort:  {},
	Mirror:      {},
	Criticality: {},
}

var incomingKey = map[string]struct{}{
	Caller: {},
}

// IsOutgoingKey represent this key should propagate by rpc.
func IsOutgoingKey(key string) bool {
	_, ok := outgoingKey[key]
	return ok
}

// IsIncomingKey represent this key should extract from rpc metadata.
func IsIncomingKey(key string) (ok bool) {
	_, ok = outgoingKey[key]
	if ok {
		return
	}
	_, ok = incomingKey[key]
	return
}
