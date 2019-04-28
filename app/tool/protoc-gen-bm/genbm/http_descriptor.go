package genbm

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// BMServerDescriptor descriptor for BM server
type BMServerDescriptor struct {
	Name         string
	ProtoService *descriptor.ServiceDescriptorProto
	Methods      []*BMMethodDescriptor
}

// BMMethodDescriptor descriptor for BM http method
type BMMethodDescriptor struct {
	Name        string
	Method      string
	PathPattern string
	RequestType string
	ReplyType   string
	ProtoMethod *descriptor.MethodDescriptorProto
	HTTPRule    *annotations.HttpRule
}

// ParseBMServer parse BMServerDescriptor form service descriptor proto
func ParseBMServer(service *descriptor.ServiceDescriptorProto) (*BMServerDescriptor, error) {
	glog.V(1).Infof("parse bmserver from service %s", service.GetName())
	serverDesc := &BMServerDescriptor{
		Name:         service.GetName(),
		ProtoService: service,
	}
	for _, method := range service.GetMethod() {
		if !HasHTTPRuleOptions(method) {
			glog.V(5).Infof("method %s not include http rule, skipped", method.GetName())
			continue
		}
		bmMethod, err := ParseBMMethod(method)
		if err != nil {
			return nil, err
		}
		serverDesc.Methods = append(serverDesc.Methods, bmMethod)
	}
	return serverDesc, nil
}

// ParseBMMethod parse BMMethodDescriptor form method descriptor proto
func ParseBMMethod(method *descriptor.MethodDescriptorProto) (*BMMethodDescriptor, error) {
	glog.V(1).Infof("parse bmmethod from method %s", method.GetName())
	ext, err := proto.GetExtension(method.GetOptions(), annotations.E_Http)
	if err != nil {
		return nil, fmt.Errorf("get extension error: %s", err)
	}
	rule := ext.(*annotations.HttpRule)
	var httpMethod string
	var pathPattern string
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		pathPattern = pattern.Get
		httpMethod = http.MethodGet
	case *annotations.HttpRule_Put:
		pathPattern = pattern.Put
		httpMethod = http.MethodPut
	case *annotations.HttpRule_Post:
		pathPattern = pattern.Post
		httpMethod = http.MethodPost
	case *annotations.HttpRule_Patch:
		pathPattern = pattern.Patch
		httpMethod = http.MethodPatch
	case *annotations.HttpRule_Delete:
		pathPattern = pattern.Delete
		httpMethod = http.MethodDelete
	default:
		return nil, fmt.Errorf("unsupport http pattern %s", rule.Pattern)
	}
	if len(rule.AdditionalBindings) != 0 {
		glog.Warningf("unsupport additional binding, additional binding will be ignored")
	}
	// TODO: support use type from other package
	requestType := splitLastElem(method.GetInputType(), ".")
	replyType := splitLastElem(method.GetOutputType(), ".")
	bmMethod := &BMMethodDescriptor{
		Name:        method.GetName(),
		Method:      httpMethod,
		PathPattern: pathPattern,
		RequestType: requestType,
		ReplyType:   replyType,
		ProtoMethod: method,
		HTTPRule:    rule,
	}
	glog.V(5).Infof("bmMethod %s: %s %s, Request:%s Reply: %s", bmMethod.Name, bmMethod.Method, bmMethod.PathPattern, bmMethod.RequestType, bmMethod.ReplyType)
	return bmMethod, nil
}

// HasHTTPRuleOptions check method has httprule extension
func HasHTTPRuleOptions(method *descriptor.MethodDescriptorProto) bool {
	options := method.GetOptions()
	if options == nil {
		return false
	}
	return proto.HasExtension(options, annotations.E_Http)
}
