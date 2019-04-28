package conf

import (
	"fmt"

	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// TiDBInsConf tidb instance config
type TiDBInsConf struct {
	Name          string
	ClusterID     string
	Addrs         []string
	Offset        int64
	CommitTS      int64
	MonitorPeriod xtime.Duration `toml:"monitor_period"`
	Databases     []*Database    `toml:"db"`
}

func newTiDBConf(fn, fc string) (c *TiDBInsConf, err error) {
	var ic struct {
		InsConf *TiDBInsConf `toml:"instance"`
	}
	if _, err = toml.Decode(fc, &ic); err != nil {
		return
	}
	if ic.InsConf == nil {
		err = fmt.Errorf("file(%s) cannot decode toml", fn)
		return
	}
	return ic.InsConf, nil
}

// TiDBEvent .
func TiDBEvent() chan *TiDBInsConf {
	return tidbEvent
}
