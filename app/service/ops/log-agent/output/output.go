package output

import (
	"fmt"
	"context"
	"go-common/app/service/ops/log-agent/event"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
)

type configDecodeFunc = func(md toml.MetaData, primValue toml.Primitive) (c interface{}, err error)

type Output interface {
	Run() (err error)
	Stop()
	InputChan() (chan *event.ProcessorEvent)
}

// Factory is used to register functions creating new Input instances.
type Factory = func(ctx context.Context, config interface{}) (Output, error)

var registry = make(map[string]Factory)

var runningOutput = make(map[string]Output)

func Register(name string, factory Factory) error {
	log.Info("Registering output factory")
	if name == "" {
		return fmt.Errorf("Error registering output: name cannot be empty")
	}
	if factory == nil {
		return fmt.Errorf("Error registering output '%v': factory cannot be empty", name)
	}
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering output '%v': already registered", name)
	}

	registry[name] = factory
	log.Info("Successfully registered output")

	return nil
}

func OutputExists(name string) bool {
	_, exists := registry[name]
	return exists
}

func GetFactory(name string) (Factory, error) {
	if _, exists := registry[name]; !exists {
		return nil, fmt.Errorf("Error creating output. No such output type exist: '%v'", name)
	}
	return registry[name], nil
}

func GetOutputChan(name string) (chan *event.ProcessorEvent, error) {
	if name == "" {
		name = "lancer-ops-log"
	}
	if _, exists := runningOutput[name]; !exists {
		return nil, fmt.Errorf("Error getting output chan. No such output chan exist: '%v'", name)
	}
	return runningOutput[name].InputChan(), nil
}

func OutputRunning(name string) bool {
	_, exists := runningOutput[name]
	return exists
}

func RegisterOutput(name string, o Output) (error) {
	if name == "" {
		return nil
	}
	if _, exists := runningOutput[name]; exists {
		return fmt.Errorf("output %s already running", name)
	}
	runningOutput[name] = o
	return nil
}

func ChanConnect(ctx context.Context, from <-chan *event.ProcessorEvent, to chan<- *event.ProcessorEvent) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-from:
				to <- e
			}
		}
	}()
	return
}
