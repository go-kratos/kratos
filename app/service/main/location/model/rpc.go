package model

// ArgIP for Zone、Info、InfoComplete、Info2
type ArgIP struct {
	IP string
}

// Archive for Archive
type Archive struct {
	Aid int64
	Mid int64
	IP  string
	CIP string
}

// Group for group
type Group struct {
	Gid int64
	Mid int64
	IP  string
	CIP string
}

// ArgPids for PIDs
type ArgPids struct {
	Pids, IP, CIP string
}
