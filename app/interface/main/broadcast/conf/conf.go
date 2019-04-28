package conf

import (
	"time"

	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/naming/discovery"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Broadcast *Broadcast
	Log       *log.Config
	HTTP      *bm.ServerConfig
	Tracer    *trace.Config
	Ecode     *ecode.Config

	WardenServer *warden.ServerConfig
	WardenClient *warden.ClientConfig
	Discovery    *discovery.Config
	HTTPClient   *bm.ClientConfig

	TCP          *TCP
	WebSocket    *WebSocket
	Timer        *Timer
	ProtoSection *ProtoSection
	Whitelist    *Whitelist
	Bucket       *Bucket
}

// Broadcast config.
type Broadcast struct {
	Debug         bool
	MaxProc       int
	ServerTick    xtime.Duration
	OnlineTick    xtime.Duration
	Failover      bool
	APIHost       string
	APIToken      string
	OnlineRetries int
	OpenPortV1    bool
}

// TCP config
type TCP struct {
	Bind         []string
	BindV1       []string
	Sndbuf       int
	Rcvbuf       int
	Keepalive    bool
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// WebSocket  config
type WebSocket struct {
	Bind        []string
	BindV1      []string
	TLSOpen     bool
	TLSBind     []string
	TLSBindV1   []string
	CertFile    string
	PrivateFile string
}

// Timer config
type Timer struct {
	Timer     int
	TimerSize int
}

// ProtoSection config
type ProtoSection struct {
	HandshakeTimeout xtime.Duration
	WriteTimeout     xtime.Duration
	SvrProto         int
	CliProto         int
}

// Whitelist .
type Whitelist struct {
	Whitelist []int64
	WhiteLog  string
}

// Bucket .
type Bucket struct {
	Size          int
	Channel       int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}

// Fix fix config to default.
func (c *Config) Fix() {
	if c.Broadcast == nil {
		c.Broadcast = new(Broadcast)
	}
	if c.Broadcast.MaxProc <= 0 {
		c.Broadcast.MaxProc = 32
	}
	if c.Broadcast.ServerTick <= 0 {
		c.Broadcast.ServerTick = xtime.Duration(5 * time.Second)
	}
	if c.Broadcast.OnlineTick <= 0 {
		c.Broadcast.OnlineTick = xtime.Duration(10 * time.Second)
	}
	if c.Broadcast.APIHost == "" {
		c.Broadcast.APIHost = "http://api.bilibili.com"
	}
}

// Set set config and decode.
func (c *Config) Set(text string) error {
	var tmp Config
	if _, err := toml.Decode(text, &tmp); err != nil {
		return err
	}
	tmp.Fix()
	*c = tmp
	return nil
}
