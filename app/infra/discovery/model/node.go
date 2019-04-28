package model

// NodeStatus Status of instance
type NodeStatus int

const (
	// AppID is discvoery id
	AppID = "infra.discovery"
)

const (
	// NodeStatusUP Ready to receive register
	NodeStatusUP NodeStatus = iota
	// NodeStatusLost lost with each other
	NodeStatusLost
)

// Node node
type Node struct {
	Addr   string     `json:"addr"`
	Status NodeStatus `json:"status"`
	Zone   string     `json:"zone"`
}
