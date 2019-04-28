package conf

import (
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/naming/discovery"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log       *log.Config
	Ecode     *ecode.Config
	HTTP      *bm.ServerConfig
	RPC       *warden.ClientConfig
	Databus   *databus.Config
	Discovery *discovery.Config
	Routine   *Routine
	Room      *Room
}

// Routine routine.
type Routine struct {
	Size uint64
	Chan uint64
}

// Room room.
type Room struct {
	Refresh  xtime.Duration
	Idle     xtime.Duration
	Batch    int
	Signal   xtime.Duration
	Compress bool
}

// Set set config and decode.
func (c *Config) Set(text string) error {
	var tmp Config
	if _, err := toml.Decode(text, &tmp); err != nil {
		return err
	}
	*c = tmp
	return nil
}
