package server

import (
	"errors"
	"fmt"
	"github.com/ipipdotnet/ipdb-go"
	"go-common/app/service/live/broadcast-proxy/conf"
	"go-common/app/service/live/broadcast-proxy/dispatch"
	"go-common/app/service/live/broadcast-proxy/grocery"
	"go-common/library/log"
	"sync"
)

type CometDispatcher struct {
	sven     *grocery.SvenClient
	stopper  chan struct{}
	wg       sync.WaitGroup
	locker   sync.RWMutex
	ipDataV4 *ipdb.City
	ipDataV6 *ipdb.City
	matcher  *dispatch.Matcher
	config   *conf.DispatchConfig
}

func NewCometDispatcher(ipipConfig *conf.IpipConfig,
	dispatchConfig *conf.DispatchConfig, svenConfig *conf.SvenConfig) (*CometDispatcher, error) {
	sven, err := grocery.NewSvenClient(svenConfig.TreeID, svenConfig.Zone, svenConfig.Env, svenConfig.Build,
		svenConfig.Token)
	if err != nil {
		return nil, err
	}
	ipDataV4, err := ipdb.NewCity(ipipConfig.V4)
	if err != nil {
		return nil, err
	}
	ipDataV6, err := ipdb.NewCity(ipipConfig.V6)
	if err != nil {
		return nil, err
	}
	dispatcher := &CometDispatcher{
		sven:     sven,
		stopper:  make(chan struct{}),
		ipDataV4: ipDataV4,
		ipDataV6: ipDataV6,
		config:   dispatchConfig,
	}
	config := sven.Config()
	if data, ok := config.Config[dispatchConfig.FileName]; ok {
		dispatcher.updateDispatchConfig(data)
	} else {
		return nil, errors.New(fmt.Sprintf("cannot find %s in sven config", dispatchConfig.FileName))
	}
	dispatcher.wg.Add(1)
	go func() {
		defer dispatcher.wg.Done()
		dispatcher.configWatcherProcess(dispatchConfig.FileName)
	}()
	return dispatcher, nil
}

func (dispatcher *CometDispatcher) Close() {
	close(dispatcher.stopper)
	dispatcher.wg.Wait()
	dispatcher.sven.Close()
}

func (dispatcher *CometDispatcher) updateDispatchConfig(config string) error {
	matcher, err := dispatch.NewMatcher([]byte(config), dispatcher.ipDataV4, dispatcher.ipDataV6, dispatcher.config)
	if err != nil {
		log.Error("parse rule config error:%v, data:%s", err, config)
		return err
	}
	dispatcher.locker.Lock()
	dispatcher.matcher = matcher
	dispatcher.locker.Unlock()
	log.Info("parse rule config ok, data:%s", config)
	return nil
}

func (dispatcher *CometDispatcher) configWatcherProcess(filename string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for config := range dispatcher.sven.ConfigNotify() {
			log.Info("[sven]New version:%d", config.Version)
			log.Info("[sven]New config: %v", config.Config)
			if data, ok := config.Config[filename]; ok {
				dispatcher.updateDispatchConfig(data)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range dispatcher.sven.LogNotify() {
			log.Info("[sven]log level:%v, message:%v", e.Level, e.Message)
		}
	}()
	wg.Wait()
}

func (dispatcher *CometDispatcher) Dispatch(ip string, uid int64) ([]string, []string) {
	var matcher *dispatch.Matcher
	dispatcher.locker.RLock()
	matcher = dispatcher.matcher
	dispatcher.locker.RUnlock()
	if matcher == nil {
		return []string{dispatcher.config.DefaultDomain}, []string{dispatcher.config.DefaultDomain}
	}
	return matcher.Dispatch(ip, uid)
}
