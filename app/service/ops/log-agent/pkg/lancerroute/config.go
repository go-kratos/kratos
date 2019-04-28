package lancerroute

import (
	"time"
	"fmt"

	"go-common/app/service/ops/log-agent/conf/configcenter"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
)

const (
	lancerRouteFile = "lancerRoute.toml"
)

type Config struct {
	LancerRoute map[string]string `toml:"lancerroute"`
}

func (l *Lancerroute) InitConfig() (err error) {
	// read config from config center
	if err = l.readConfig(); err != nil {
		return err
	}
	// watch update and reload config
	go func() {
		currentVersion := configcenter.Version
		for {
			if currentVersion != configcenter.Version {
				log.Info("lancer route config reload")
				if err := l.readConfig(); err != nil {
					log.Error("lancer route config reload error (%v", err)
				}
				currentVersion = configcenter.Version
			}
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (l *Lancerroute) readConfig() (err error) {
	var (
		ok             bool
		value          string
		tmpLancerRoute Config
	)

	// sample config
	if value, ok = configcenter.Client.Value(lancerRouteFile); !ok {
		return fmt.Errorf("failed to get %s", lancerRouteFile)
	}
	if _, err = toml.Decode(value, &tmpLancerRoute); err != nil {
		return err
	}
	l.c = &tmpLancerRoute
	return nil
}
