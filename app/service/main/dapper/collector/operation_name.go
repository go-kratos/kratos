package collector

import (
	"net/url"
	"strings"
)

import (
	"go-common/app/service/main/dapper/model"
)

// OperationNameProcess fix operation name so sad!
type OperationNameProcess struct{}

// Process implement operation name
func (o *OperationNameProcess) Process(span *model.Span) error {
	switch {
	case !span.IsServer() && strings.HasPrefix(span.OperationName, "http://"):
		o.fixHTTP(span)
	}
	return nil
}

func (o *OperationNameProcess) fixHTTP(span *model.Span) {
	oldOperationName := span.OperationName
	method := "UNKONWN"
	if methodTag := span.GetTagString("http.method"); methodTag != "" {
		method = methodTag
	}
	operationName := "HTTP:" + method
	span.SetOperationName(operationName)

	peerSign := oldOperationName
	if strings.HasPrefix(oldOperationName, "http://") {
		if reqURL, err := url.Parse(oldOperationName); err == nil {
			peerSign = reqURL.Path
			span.SetTag("http.url", oldOperationName)
		}
	}
	span.SetTag("_peer.sign", peerSign)
}

// NewOperationNameProcess .
func NewOperationNameProcess() Processer {
	return &OperationNameProcess{}
}
