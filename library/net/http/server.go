package http

import (
	"net"
	"net/http"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

var (
	_defaultAddr string
)

// ServerConfig is the alias to bm ServerConfig
//
// Deprecated: using bm.ServerConfig instead
type ServerConfig = bm.ServerConfig

func init() {
	if env.HTTPPort != "" {
		_defaultAddr = net.JoinHostPort("0.0.0.0", env.HTTPPort)
	} else {
		_defaultAddr = "0.0.0.0:8000"
	}
}

// Serve listen and serve bm engine by given config.
//
// Deprecated: using Engine.Start instead
func Serve(engine *bm.Engine, conf *ServerConfig) error {
	if conf == nil {
		conf = &ServerConfig{
			Addr:    _defaultAddr,
			Timeout: xtime.Duration(time.Second),
		}
	}
	l, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		errors.Wrapf(err, "listen tcp: %d", conf.Addr)
		return err
	}
	if err := engine.SetConfig(conf); err != nil {
		return err
	}

	log.Info("blademaster: start http listen addr: %s", conf.Addr)
	server := &http.Server{
		ReadTimeout:  time.Duration(conf.ReadTimeout),
		WriteTimeout: time.Duration(conf.WriteTimeout),
	}
	go func() {
		if err := engine.RunServer(server, l); err != nil {
			log.Error("blademaster: engine.ListenServer(%+v, %+v) error(%v)", server, l, err)
		}
	}()

	return nil
}
