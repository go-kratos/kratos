package trace

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"go-common/library/conf/dsn"
	xtime "go-common/library/time"
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
	// DEPRECATED
	Proto string `dsn:"network"`
	// DEPRECATED
	Chan int `dsn:"query.chan,"`
	// Report timeout
	Timeout xtime.Duration `dsn:"query.timeout,200ms"`
	// DisableSample
	DisableSample bool `dsn:"query.disable_sample"`
	// probabilitySampling
	Probability float32 `dsn:"-"`
	// ProtocolVersion
	ProtocolVersion int32 `dsn:"query.protocol_version,2"`
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
	serviceName := serviceNameFromEnv()
	return newTracer(serviceName, report, cfg), nil
}

// Init 兼容以前的 Init 写法
func Init(cfg *Config) {
	serviceName := serviceNameFromEnv()
	if cfg != nil {
		// NOTE compatible proto field
		cfg.Network = cfg.Proto
		fmt.Fprintf(os.Stderr, "[deprecated] trace.Init() with conf is Deprecated, argument will be ignored. please use flag -trace or env TRACE to configure trace.\n")
		report := newReport(cfg.Network, cfg.Addr, time.Duration(cfg.Timeout), cfg.ProtocolVersion)
		SetGlobalTracer(newTracer(serviceName, report, cfg))
		return
	}
	// paser config from env
	cfg, err := parseDSN(_traceDSN)
	if err != nil {
		panic(fmt.Errorf("parse trace dsn error: %s", err))
	}
	report := newReport(cfg.Network, cfg.Addr, time.Duration(cfg.Timeout), cfg.ProtocolVersion)
	// disable sample if uat env
	cfg.DisableSample = isUATEnv()
	SetGlobalTracer(newTracer(serviceName, report, cfg))
}
