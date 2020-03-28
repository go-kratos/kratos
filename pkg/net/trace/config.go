package trace

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/go-kratos/kratos/pkg/conf/dsn"
	"github.com/go-kratos/kratos/pkg/conf/env"
	xtime "github.com/go-kratos/kratos/pkg/time"
)

var _traceDSN = "unixgram:///var/run/dapper-collect/dapper-collect.sock"

func init() {
	if v := os.Getenv("TRACE"); v != "" {
		_traceDSN = v
	}
	flag.StringVar(&_traceDSN, "trace", _traceDSN, "trace report dsn, or use TRACE env.")
}

// Config config.
type Config struct {
	// Report network e.g. unixgram, tcp, udp
	Network string `dsn:"network"`
	// For TCP and UDP networks, the addr has the form "host:port".
	// For Unix networks, the address must be a file system path.
	Addr string `dsn:"address"`
	// Report timeout
	Timeout xtime.Duration `dsn:"query.timeout,200ms"`
	// DisableSample
	DisableSample bool `dsn:"query.disable_sample"`
	// ProtocolVersion
	ProtocolVersion int32 `dsn:"query.protocol_version,1"`
	// Probability probability sampling
	Probability float32 `dsn:"-"`
}

func parseDSN(rawdsn string) (*Config, error) {
	d, err := dsn.Parse(rawdsn)
	if err != nil {
		return nil, errors.Wrapf(err, "trace: invalid dsn: %s", rawdsn)
	}
	cfg := new(Config)
	if _, err = d.Bind(cfg); err != nil {
		return nil, errors.Wrapf(err, "trace: invalid dsn: %s", rawdsn)
	}
	return cfg, nil
}

// TracerFromEnvFlag new tracer from env and flag
func TracerFromEnvFlag() (Tracer, error) {
	cfg, err := parseDSN(_traceDSN)
	if err != nil {
		return nil, err
	}
	report := newReport(cfg.Network, cfg.Addr, time.Duration(cfg.Timeout), cfg.ProtocolVersion)
	return NewTracer(env.AppID, report, cfg.DisableSample), nil
}

// Init init trace report.
func Init(cfg *Config) {
	if cfg == nil {
		// paser config from env
		var err error
		if cfg, err = parseDSN(_traceDSN); err != nil {
			panic(fmt.Errorf("parse trace dsn error: %s", err))
		}
	}
	report := newReport(cfg.Network, cfg.Addr, time.Duration(cfg.Timeout), cfg.ProtocolVersion)
	SetGlobalTracer(NewTracer(env.AppID, report, cfg.DisableSample))
}
