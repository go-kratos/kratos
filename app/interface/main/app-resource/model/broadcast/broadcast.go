package broadcast

import (
	wardensvr "go-common/app/service/main/broadcast/api/grpc/v1"
)

type ServerListReply struct {
	Domain       string   `json:"domain,omitempty"`
	TCPPort      int32    `json:"tcp_port,omitempty"`
	WsPort       int32    `json:"ws_port,omitempty"`
	WssPort      int32    `json:"wss_port,omitempty"`
	Heartbeat    int32    `json:"heartbeat,omitempty"`
	HeartbeatMax int32    `json:"heartbeat_max,omitempty"`
	Nodes        []string `json:"nodes,omitempty"`
	Backoff      *Backoff `json:"backoff,omitempty"`
}

type Backoff struct {
	MaxDelay  int32   `json:"max_delay,omitempty"`
	BaseDelay int32   `json:"base_delay,omitempty"`
	Factor    float32 `json:"factor,omitempty"`
	Jitter    float32 `json:"jitter,omitempty"`
}

func (l *ServerListReply) ServerListChange(w *wardensvr.ServerListReply) {
	l.Domain = w.Domain
	l.TCPPort = w.TcpPort
	l.WsPort = w.WsPort
	l.WssPort = w.WssPort
	l.Heartbeat = w.Heartbeat
	l.HeartbeatMax = w.HeartbeatMax
	if len(w.Nodes) > 0 {
		l.Nodes = w.Nodes
	}
	if w.Backoff != nil {
		l.Backoff = &Backoff{
			MaxDelay:  w.Backoff.MaxDelay,
			BaseDelay: w.Backoff.BaseDelay,
			Factor:    w.Backoff.Factor,
			Jitter:    w.Backoff.Jitter,
		}
	}
}
