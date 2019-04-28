package collector

import (
	"sync"

	"go-common/app/service/main/dapper/model"
)

type peerServiceDetect struct {
	rmx  sync.RWMutex
	pair map[string]string
}

func (p *peerServiceDetect) detect(operationName string) (string, bool) {
	p.rmx.RLock()
	serviceName, ok := p.pair[operationName]
	p.rmx.RUnlock()
	return serviceName, ok
}

func (p *peerServiceDetect) add(serviceName, operationName string) {
	if operationName == "" || serviceName == "" {
		// ignored empty
		return
	}
	p.rmx.RLock()
	val, ok := p.pair[operationName]
	p.rmx.RUnlock()
	if !ok || serviceName != val {
		p.rmx.Lock()
		p.pair[operationName] = serviceName
		p.rmx.Unlock()
	}
}

func (p *peerServiceDetect) process(span *model.Span) {
	if span.IsServer() {
		p.add(span.ServiceName, span.OperationName)
		return
	}
	if span.GetTagString("peer.service") != "" {
		return
	}
	peerService, ok := p.detect(span.OperationName)
	if ok {
		span.SetTag("peer.service", peerService)
		span.SetTag("_auto.peer.service", true)
		return
	}
	if peerSign := span.StringTag("_peer.sign"); peerSign != "" {
		peerService, ok := p.detect(peerSign)
		if ok {
			span.SetTag("peer.service", peerService)
			span.SetTag("_auto.peer.service", true)
		}
	}
}

func (p *peerServiceDetect) Process(span *model.Span) error {
	p.process(span)
	return nil
}

// NewPeerServiceDetectProcesser .
func NewPeerServiceDetectProcesser(data map[string]map[string]struct{}) Processer {
	p := &peerServiceDetect{pair: make(map[string]string)}
	for serviceName, operationNames := range data {
		for operationName := range operationNames {
			p.add(serviceName, operationName)
		}
	}
	return p
}
