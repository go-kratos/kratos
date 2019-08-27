package metadata

// metadata common key
const (
	// Network
	RemoteIP   = "remote_ip"
	RemotePort = "remote_port"

	// Router
	Color  = "color"
	Caller = "caller"

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

	// Criticality
	Criticality = "criticality"

	// Locale locale language.
	Locale = "locale"
)

var outgoingKey = map[string]struct{}{
	Color:       struct{}{},
	RemoteIP:    struct{}{},
	RemotePort:  struct{}{},
	Mirror:      struct{}{},
	Criticality: struct{}{},
}

var incomingKey = map[string]struct{}{
	Caller: struct{}{},
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
