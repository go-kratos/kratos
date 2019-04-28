package sock

import (
	"errors"
	"time"

	xtime "go-common/library/time"
)

type Config struct {
	TcpAddr          string         `toml:"tcpAddr"`
	UdpAddr          string         `toml:"udpAddr"`
	ReadChanSize     int            `toml:"readChanSize"`
	TcpBatchMaxBytes int            `toml:"tcpBatchMaxBytes"`
	UdpPacketMaxSize int            `toml:"udpPacketMaxSize"`
	LogMaxBytes      int            `toml:"logMaxBytes"`
	UdpReadTimeout   xtime.Duration `toml:"udpReadTimeout"`
	TcpReadTimeout   xtime.Duration `toml:"tcpReadTimeout"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of Sock Input is nil")
	}

	if c.TcpAddr == "" {
		c.TcpAddr = "/var/run/lancer/collector_tcp.sock"
	}
	if c.UdpAddr == "" {
		c.UdpAddr = "/var/run/lancer/collector.sock"
	}

	if c.ReadChanSize == 0 {
		c.ReadChanSize = 5000
	}

	if c.TcpBatchMaxBytes == 0 {
		c.TcpBatchMaxBytes = 10240000 // 10MB
	}

	if c.UdpPacketMaxSize == 0 {
		c.UdpPacketMaxSize = 1024 * 64 //64KB
	}

	if c.UdpReadTimeout == 0 {
		c.UdpReadTimeout = xtime.Duration(time.Second * 10)
	}

	if c.TcpReadTimeout == 0 {
		c.TcpReadTimeout = xtime.Duration(time.Minute * 5)
	}
	return nil
}
