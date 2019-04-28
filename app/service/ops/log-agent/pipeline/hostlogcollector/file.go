package hostlogcollector

import (
	"os"
	"io/ioutil"
	"path"
	"strings"
	"fmt"
	"time"
	"context"

	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pipeline"
)

type HostLogCollector struct {
	c      *Config
	ctx    context.Context
	cancel context.CancelFunc
}

func InitHostLogCollector(ctx context.Context, c *Config) (err error) {
	if err = c.ConfigValidate(); err != nil {
		return err
	}
	collector := new(HostLogCollector)
	collector.c = c

	collector.ctx, collector.cancel = context.WithCancel(ctx)

	go collector.scan()

	return nil
}

//
func (collector *HostLogCollector) scan() {
	ticker := time.Tick(time.Duration(collector.c.ScanInterval))
	for {
		select {
		case <-ticker:
			configPaths, err := collector.getConfigs()
			if err != nil {
				log.Error("failed to scan hostlogcollector config file list: %s", err)
				continue
			}
			for _, configPath := range configPaths {
				config, err := ioutil.ReadFile(configPath)
				if err != nil {
					log.Error("filed to read hostlogcollector config file %s: %s", configPath, err)
					continue
				}
				if !pipeline.PipelineManagement.PipelineExisted(configPath) {
					go pipeline.PipelineManagement.StartPipeline(collector.ctx, configPath, string(config))
				}
			}
		case <-collector.ctx.Done():
			return
		}
	}
}

// HostLogCollector get file collect configs under path
func (collector *HostLogCollector) getConfigs() ([]string, error) {
	var (
		err         error
		cinfos      []os.FileInfo
		configFiles = make([]string, 0)
	)

	dinfo, err := os.Lstat(collector.c.HostConfigPath)

	if err != nil {
		return nil, fmt.Errorf("lstat(%s) failed: %s", collector.c.HostConfigPath, err)
	}

	if !dinfo.IsDir() {
		return nil, fmt.Errorf("file collect config path must be dir")
	}

	if cinfos, err = ioutil.ReadDir(collector.c.HostConfigPath); err != nil {
		return nil, fmt.Errorf("ioutil.ReadDir(%s) error(%v)", collector.c.HostConfigPath, err)
	}

	for _, cinfo := range cinfos {
		name := path.Join(collector.c.HostConfigPath, cinfo.Name())
		if !cinfo.IsDir() && strings.HasSuffix(name, collector.c.ConfigSuffix) {
			configFiles = append(configFiles, name)
		}
	}
	return configFiles, nil
}
