package conf

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"go-common/library/conf"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// ENV Key
const (
	BNSDNSAddr  = "BNS_DNS_HOST"
	BNSDNSPort  = "BNS_DNS_PORT"
	BNSHTTPAddr = "BNS_HTTP_ADDR"
	BNSHTTPPort = "BNS_HTTP_PORT"
)

// default value
const (
	defaultBNSDNSAddr  = "0.0.0.0"
	defaultBNSDNSPort  = 15353
	defaultBNSHTTPAddr = "0.0.0.0"
	defaultBNSHTTPPort = 15380
)

var defaultConfig Config

func init() {
	// default dns config
	defaultDNSConfig := DNSConfig{
		TTL:             0,
		AllowStale:      true,
		UDPAnswerLimit:  3,
		MaxStale:        xtime.Duration(time.Second * 87600),
		Domain:          "bili.",
		RecursorTimeout: xtime.Duration(time.Second),
	}
	defaultDNSServer := &DNSServer{
		Addr:   defaultBNSDNSAddr,
		Port:   defaultBNSDNSPort,
		Config: &defaultDNSConfig,
	}

	// default http config
	defaultHTTPServer := &HTTPServer{
		Addr: defaultBNSHTTPAddr,
		Port: defaultBNSHTTPPort,
	}

	defaultBackend := &Backend{
		Backend: "discovery",
		Config: map[string]interface{}{
			"url": "http://api.bilibili.co",
		},
	}

	defaultConfig = Config{
		Backend: defaultBackend,
		HTTP:    defaultHTTPServer,
		DNS:     defaultDNSServer,
	}
}

// LoadConfig from source
func LoadConfig(source string) (*Config, error) {
	cfg := defaultConfig
	var err error
	if strings.HasPrefix(source, "remote://") {
		var u *url.URL
		if u, err = url.Parse(source); err != nil {
			return nil, err
		}
		err = loadRemoteConfig(u.Path, &cfg)
	} else if source != "" {
		err = loadLocalConfig(source, &cfg)
	}
	if err != nil {
		return nil, err
	}
	overwriteByEnv(&cfg)
	return &cfg, nil
}

func loadRemoteConfig(key string, pcfg *Config) error {
	client, err := conf.New()
	if err != nil {
		return err
	}

	data, ok := client.Value(key)
	if !ok {
		return errors.New("load config center error")
	}

	if _, err = toml.Decode(data, pcfg); err != nil {
		return errors.New("could not decode config")
	}

	go func() {
		for range client.Event() {
			log.Warn("ignore config reload")
		}
	}()
	return nil
}

func loadLocalConfig(fpath string, pcfg *Config) error {
	_, err := toml.DecodeFile(fpath, pcfg)
	return err
}

// Config config struct
type Config struct {
	Log     *log.Config
	Backend *Backend
	HTTP    *HTTPServer
	DNS     *DNSServer
}

// overwrite config from env
func overwriteByEnv(pcfg *Config) {
	if addr := os.Getenv(BNSDNSAddr); addr != "" {
		pcfg.DNS.Addr = addr
	}
	if portStr := os.Getenv(BNSDNSPort); portStr != "" {
		if port, err := strconv.Atoi(portStr); err != nil {
			log.Warn("parse port from env error: %s", err)
		} else {
			pcfg.DNS.Port = port
		}
	}

	if addr := os.Getenv(BNSHTTPAddr); addr != "" {
		pcfg.HTTP.Addr = addr
	}
	if portStr := os.Getenv(BNSHTTPPort); portStr != "" {
		if port, err := strconv.Atoi(portStr); err != nil {
			log.Warn("parse port from env error: %s", err)
		} else {
			pcfg.HTTP.Port = port
		}
	}
}

// Backend Config
type Backend struct {
	Backend string
	Config  map[string]interface{}
}

// HTTPServer http server config
type HTTPServer struct {
	Addr string
	Port int
}

// DNSServer dns server config
type DNSServer struct {
	Addr   string
	Port   int
	Config *DNSConfig
}

// DNSConfig dns config
type DNSConfig struct {
	// TTL provides the TTL value for a easyns path query for given path.
	// The "*" wildcard can be used to set a default to a highlevel path, such as project level path.
	TTL xtime.Duration `toml:"ttl"`

	// AllowStale is used to enable lookups with stale
	// data. This gives horizontal read scalability since
	// any easyns server can service the query instead of
	// only the leader.
	AllowStale bool

	// EnableTruncate is used to enable setting the truncate
	// flag for UDP DNS queries.  This allows unmodified
	// clients to re-query the easyns server using TCP
	// when the total number of records exceeds the number
	// returned by default for UDP.
	EnableTruncate bool

	// UDPAnswerLimit is used to limit the maximum number of DNS Resource
	// Records returned in the ANSWER section of a DNS response. This is
	// not normally useful and will be limited based on the querying
	// protocol, however systems that implemented §6 Rule 9 in RFC3484
	// may want to set this to `1` in order to subvert §6 Rule 9 and
	// re-obtain the effect of randomized resource records (i.e. each
	// answer contains only one IP, but the IP changes every request).
	// RFC3484 sorts answers in a deterministic order, which defeats the
	// purpose of randomized DNS responses.  This RFC has been obsoleted
	// by RFC6724 and restores the desired behavior of randomized
	// responses, however a large number of Linux hosts using glibc(3)
	// implemented §6 Rule 9 and may need this option (e.g. CentOS 5-6,
	// Debian Squeeze, etc).
	UDPAnswerLimit int `toml:"udpAnswerLimit"`

	// MaxStale is used to bound how stale of a result is
	// accepted for a DNS lookup. This can be used with
	// AllowStale to limit how old of a value is served up.
	// If the stale result exceeds this, another non-stale
	// stale read is performed.
	MaxStale xtime.Duration

	// DisableCompression is used to control whether DNS responses are
	// compressed. This was turned on by default and this
	// config was added as an opt-out.
	DisableCompression bool

	// RecursorTimeout specifies the timeout in seconds
	// for Easyns agent's internal dns client used for recursion.
	// This value is used for the connection, read and write timeout.
	// Default: 2s
	RecursorTimeout xtime.Duration

	// Managed domain suffix
	Domain string

	// Upstream recursor dns servers
	Recursors []string
}
