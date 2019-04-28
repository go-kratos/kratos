package pipeline

import (
	"sync"
	"errors"
	"context"
	"sort"
	"time"
	"os"
	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/input"
	"go-common/app/service/ops/log-agent/processor"
	"go-common/app/service/ops/log-agent/output"
	"go-common/app/service/ops/log-agent/pkg/common"
	"go-common/library/log"
	"github.com/BurntSushi/toml"
)

type PipelineMng struct {
	Pipelines     map[string]*Pipeline
	PipelinesLock sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	scanInterval  time.Duration
}

var PipelineManagement *PipelineMng

const defaultPipeline = "defaultPipeline"

func InitPipelineMng(ctx context.Context) (err error) {
	m := new(PipelineMng)
	m.Pipelines = make(map[string]*Pipeline)
	m.ctx, m.cancel = context.WithCancel(ctx)
	if err = m.StartDefaultOutput(); err != nil {
		return err
	}

	m.scanInterval = time.Second * 10

	go m.scan()

	m.StartDefaultPipeline()
	// Todo check defaultPipeline
	//if !m.PipelineExisted(defaultPipeline) {
	//	return errors.New("failed to start defaultPipeline, see log for more details")
	//}

	PipelineManagement = m

	return nil
}

func (m *PipelineMng) RegisterHostFileCollector(configPath string, p *Pipeline) {
	m.PipelinesLock.Lock()
	defer m.PipelinesLock.Unlock()
	m.Pipelines[configPath] = p
}

func (m *PipelineMng) UnRegisterHostFileCollector(configPath string) {
	m.PipelinesLock.Lock()
	defer m.PipelinesLock.Unlock()
	delete(m.Pipelines, configPath)
}

func (m *PipelineMng) PipelineExisted(configPath string) bool {
	m.PipelinesLock.RLock()
	defer m.PipelinesLock.RUnlock()
	_, ok := m.Pipelines[configPath]
	return ok
}

func (m *PipelineMng) GetPipeline(configPath string) *Pipeline {
	m.PipelinesLock.RLock()
	defer m.PipelinesLock.RUnlock()
	if pipe, ok := m.Pipelines[configPath]; ok {
		return pipe
	}
	return nil
}

// pipelines get configPath list of registered pipeline
func (m *PipelineMng) configPaths() []string {
	m.PipelinesLock.RLock()
	defer m.PipelinesLock.RUnlock()
	result := make([]string, 0, len(m.Pipelines))
	for p, _ := range m.Pipelines {
		result = append(result, p)
	}
	return result
}

func (m *PipelineMng) scan() {
	ticker := time.Tick(m.scanInterval)
	whiteList := make(map[string]struct{})
	whiteList[defaultPipeline] = struct{}{}
	for {
		select {
		case <-ticker:
			for _, configPath := range m.configPaths() {
				if _, ok := whiteList[configPath]; ok {
					continue
				}
				pipe := m.GetPipeline(configPath)
				// config removed
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					if pipe != nil {
						log.Info("config file not exist any more, stop pipeline: %s", configPath)
						pipe.Stop()
						continue
					}
				}
				// config updated
				oldMd5 := pipe.configMd5
				newMd5 := common.FileMd5(configPath)
				if oldMd5 != newMd5 {
					log.Info("config file updated, stop old pipeline: %s", configPath)
					pipe.Stop()
					continue
				}
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *PipelineMng) StartPipeline(ctx context.Context, configPath string, config string) () {
	var err error
	p := new(Pipeline)
	p.configPath = configPath
	p.configMd5 = common.FileMd5(configPath)

	ctx = context.WithValue(ctx, "configPath", configPath)

	p.ctx, p.cancel = context.WithCancel(ctx)

	defer p.Stop()

	var sortedOrder []string

	p.c = new(Config)

	md, err := toml.Decode(config, p.c)
	if err != nil {
		p.logError(err)
		return
	}
	inputToProcessor := make(chan *event.ProcessorEvent)
	// start input
	inputName := p.c.Input.Name
	if inputName == "" {
		p.logError(errors.New("type of Config can't be nil"))
		return
	}

	c, err := DecodeInputConfig(inputName, md, p.c.Input.Config)
	if err != nil {
		p.logError(err)
		return
	}

	InputFactory, err := input.GetFactory(inputName)

	if err != nil {
		p.logError(err)
		return
	}

	i, err := InputFactory(p.ctx, c, inputToProcessor)
	if err != nil {
		p.logError(err)
		return
	}

	if err = i.Run(); err != nil {
		p.logError(err)
		return
	}

	// start processor
	var ProcessorConnector chan *event.ProcessorEvent
	ProcessorConnector = inputToProcessor

	sortedOrder = make([]string, 0)
	for order, _ := range p.c.Processor {
		sortedOrder = append(sortedOrder, order)
	}

	sort.Strings(sortedOrder)
	for _, order := range sortedOrder {
		name := p.c.Processor[order].Name
		if name == "" {
			p.logError(errors.New("type of Processor can't be nil"))
			return
		}
		c, err := DecodeProcessorConfig(name, md, p.c.Processor[order].Config)
		if err != nil {
			p.logError(err)
			return
		}

		proc, err := processor.GetFactory(name)
		if err != nil {
			p.logError(err)
			return
		}

		ProcessorConnector, err = proc(p.ctx, c, ProcessorConnector)
		if err != nil {
			p.logError(err)
			return
		}
	}

	// add classify and fileLog processor by default if inputName == "file"
	if inputName == "file" {
		config := `
	[processor]
	[processor.1]
	type = "classify"
	[processor.2]
	type = "fileLog"
	`
		fProcessor := new(Config)
		md, _ := toml.Decode(config, fProcessor)
		fsortedOrder := make([]string, 0)
		for order, _ := range fProcessor.Processor {
			fsortedOrder = append(fsortedOrder, order)
		}

		sort.Strings(fsortedOrder)
		for _, order := range fsortedOrder {
			name := fProcessor.Processor[order].Name
			if name == "" {
				p.logError(errors.New("type of Processor can't be nil"))
				return
			}
			fc, err := DecodeProcessorConfig(name, md, fProcessor.Processor[order].Config)
			if err != nil {
				p.logError(err)
				return
			}

			proc, err := processor.GetFactory(name)
			if err != nil {
				p.logError(err)
				return
			}

			ProcessorConnector, err = proc(p.ctx, fc, ProcessorConnector)
			if err != nil {
				p.logError(err)
				return
			}
		}
	}

	// start output
	if p.c.Output != nil {
		if len(p.c.Output) > 1 {
			p.logError(errors.New("only One Output is allowed in One pipeline"))
			return
		}
		var first string
		for key, _ := range p.c.Output {
			first = key
			break
		}

		o, err := StartOutput(p.ctx, md, p.c.Output[first])
		if err != nil {
			p.logError(err)
			return
		}
		// connect processor and output
		output.ChanConnect(m.ctx, ProcessorConnector, o.InputChan())

	} else {
		// write to default output
		if err := processor.WriteToOutput(p.ctx, "", ProcessorConnector); err != nil {
			p.logError(err)
			return
		}
	}
	m.RegisterHostFileCollector(configPath, p)

	defer m.UnRegisterHostFileCollector(configPath)

	<-p.ctx.Done()
}

func (m *PipelineMng) StartDefaultPipeline() {
	//	config := `
	//[input]
	//type = "file"
	//[input.config]
	//paths = ["/data/log-agent/log/info.log.2018-11-07.001"]
	//appId = "ops.billions.test"
	//[processor]
	//[output]
	//[output.1]
	//type = "stdout"
	//`
	//	config := `
	//[input]
	//type = "file"
	//[input.config]
	//paths = ["/data/log-agent/log/info.log.2018-*"]
	//appId = "ops.billions.test"
	//logId = "000069"
	//[processor]
	//[processor.1]
	//type = "fileLog"
	//`

	config := `
	[input]
	type = "sock"
	[input.config]

	[processor]

	[processor.1]
	type = "jsonLog"

	[processor.2]
	type = "lengthCheck"

	[processor.3]
	type = "httpStream"

	[processor.4]
	type = "sample"

	[processor.5]
	type = "classify"
	`
	go m.StartPipeline(context.Background(), defaultPipeline, config)
}

func (m *PipelineMng) StartDefaultOutput() (err error) {
	var value string
	if value, err = output.ReadConfig(); err != nil {
		return err
	}
	p := new(Pipeline)
	p.c = new(Config)
	md, err := toml.Decode(value, p.c)
	if err != nil {
		return err
	}
	return StartOutputs(m.ctx, md, p.c.Output)
}

func StartOutputs(ctx context.Context, md toml.MetaData, config map[string]ConfigItem) (err error) {
	for _, item := range config {
		name := item.Name
		if name == "" {
			return errors.New("type of Output can't be nil")
		}

		if _, err = StartOutput(ctx, md, item); err != nil {
			return err
		}
	}
	return nil
}

func StartOutput(ctx context.Context, md toml.MetaData, config ConfigItem) (o output.Output, err error) {
	name := config.Name
	if name == "" {
		return nil, errors.New("type of Output can't be nil")
	}
	c, err := DecodeOutputConfig(name, md, config.Config)
	if err != nil {
		return nil, err
	}
	OutputFactory, err := output.GetFactory(name)
	if err != nil {
		return nil, err
	}

	o, err = OutputFactory(ctx, c)
	if err != nil {
		return nil, err
	}

	if err = o.Run(); err != nil {
		return nil, err
	}

	return o, nil
}
