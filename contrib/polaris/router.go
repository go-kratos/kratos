package polaris

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/polarismesh/polaris-go"
	polarisconfig "github.com/polarismesh/polaris-go/pkg/config"
	"github.com/polarismesh/polaris-go/pkg/model"
	"github.com/polarismesh/polaris-go/pkg/model/local"
	"github.com/polarismesh/polaris-go/pkg/model/pb"
	"github.com/polarismesh/polaris-go/pkg/plugin/common"
	"github.com/polarismesh/polaris-go/pkg/plugin/servicerouter"
	v1 "github.com/polarismesh/specification/source/go/api/v1/service_manage"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/selector"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/go-kratos/kratos/v3/transport/http"
)

type router struct {
	service string
}

type RouterOption func(o *router)

// WithRouterService set the caller service name used by the route
func WithRouterService(service string) RouterOption {
	return func(o *router) {
		o.service = service
	}
}

// NodeFilter polaris dynamic router selector
func (p *Polaris) NodeFilter(opts ...RouterOption) selector.NodeFilter {
	o := router{service: p.service}
	for _, opt := range opts {
		opt(&o)
	}
	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		if len(nodes) == 0 {
			return nodes
		}
		sourceService := &model.ServiceInfo{
			Namespace: p.namespace,
			Service:   o.service,
			Metadata:  map[string]string{},
		}
		if appInfo, ok := kratos.FromContext(ctx); ok {
			sourceService.Service = appInfo.Name()
		}

		sourceService.AddArgument(model.BuildCallerServiceArgument(p.namespace, sourceService.Service))

		// process transport
		if tr, ok := transport.FromClientContext(ctx); ok {
			sourceService.AddArgument(model.BuildMethodArgument(tr.Operation()))
			sourceService.AddArgument(model.BuildPathArgument(tr.Operation()))

			for _, key := range tr.RequestHeader().Keys() {
				sourceService.AddArgument(model.BuildHeaderArgument(strings.ToLower(key), tr.RequestHeader().Get(key)))
			}

			// http
			if ht, ok := tr.(http.Transporter); ok {
				sourceService.AddArgument(model.BuildPathArgument(ht.Request().URL.Path))
				sourceService.AddArgument(model.BuildCallerIPArgument(ht.Request().RemoteAddr))

				// cookie
				for _, cookie := range ht.Request().Cookies() {
					sourceService.AddArgument(model.BuildCookieArgument(cookie.Name, cookie.Value))
				}

				// url query
				for key, values := range ht.Request().URL.Query() {
					sourceService.AddArgument(model.BuildQueryArgument(key, strings.Join(values, ",")))
				}
			}
		}

		n := make(map[string]selector.Node, len(nodes))
		for _, node := range nodes {
			n[node.Address()] = node
		}

		instances, err := p.processRouters(sourceService, buildPolarisInstance(p.namespace, nodes))
		if err != nil {
			log.Errorf("polaris process routers failed, err=%v", err)
			return nodes
		}

		newNode := make([]selector.Node, 0, len(instances))
		for _, ins := range instances {
			if v, ok := n[net.JoinHostPort(ins.GetHost(), strconv.FormatUint(uint64(ins.GetPort()), 10))]; ok {
				newNode = append(newNode, v)
			}
		}
		if len(newNode) == 0 {
			return nodes
		}
		return newNode
	}
}

func (p *Polaris) processRouters(sourceService *model.ServiceInfo, dstInstances *pb.ServiceInstancesInProto) ([]model.Instance, error) {
	if dstInstances == nil {
		return nil, nil
	}
	plugins := p.discovery.SDKContext().GetPlugins()
	ruleBasedRouter, err := plugins.GetPlugin(common.TypeServiceRouter, polarisconfig.DefaultServiceRouterRuleBased)
	if err != nil {
		return nil, err
	}
	filterOnlyRouter, err := plugins.GetPlugin(common.TypeServiceRouter, polarisconfig.DefaultServiceRouterFilterOnly)
	if err != nil {
		return nil, err
	}
	dstRouteRule, err := p.getRouteRule(dstInstances)
	if err != nil {
		return nil, err
	}
	routeInfo := &servicerouter.RouteInfo{
		SourceService:    sourceService,
		DestService:      dstInstances,
		DestRouteRule:    dstRouteRule,
		FilterOnlyRouter: filterOnlyRouter.(servicerouter.ServiceRouter),
	}
	if sourceService.HasService() {
		sourceRouteRule, routeErr := p.getRouteRule(sourceService)
		if routeErr != nil {
			return nil, routeErr
		}
		routeInfo.SourceRouteRule = sourceRouteRule
	}
	routeInfo.Init(plugins)
	result, err := servicerouter.GetFilterCluster(
		p.discovery.SDKContext().GetValueContext(),
		[]servicerouter.ServiceRouter{ruleBasedRouter.(servicerouter.ServiceRouter)},
		routeInfo,
		dstInstances.GetServiceClusters(),
	)
	if err != nil {
		return nil, err
	}
	if result == nil || result.OutputCluster == nil {
		return nil, nil
	}
	instances, _ := result.OutputCluster.GetInstances()
	return instances, nil
}

func (p *Polaris) getRouteRule(service model.ServiceMetadata) (model.ServiceRule, error) {
	resp, err := p.discovery.GetRouteRule(&polaris.GetServiceRuleRequest{
		GetServiceRuleRequest: model.GetServiceRuleRequest{
			Namespace: service.GetNamespace(),
			Service:   service.GetService(),
		},
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return &serviceRuleResponse{resp: resp}, nil
}

type serviceRuleResponse struct {
	resp *model.ServiceRuleResponse
}

func (r serviceRuleResponse) GetType() model.EventType {
	return r.resp.Type
}

func (r serviceRuleResponse) IsInitialized() bool {
	return true
}

func (r serviceRuleResponse) GetRevision() string {
	return r.resp.Revision
}

func (r serviceRuleResponse) GetHashValue() uint64 {
	return r.resp.HashValue
}

func (r serviceRuleResponse) IsNotExists() bool {
	return r.resp.NotExists
}

func (r serviceRuleResponse) GetNamespace() string {
	return r.resp.Service.Namespace
}

func (r serviceRuleResponse) GetService() string {
	return r.resp.Service.Service
}

func (r serviceRuleResponse) GetValue() interface{} {
	return r.resp.Value
}

func (r serviceRuleResponse) GetRuleCache() model.RuleCache {
	return r.resp.RuleCache
}

func (r serviceRuleResponse) GetValidateError() error {
	return r.resp.ValidateError
}

func (r serviceRuleResponse) IsCacheLoaded() bool {
	return false
}

func buildPolarisInstance(namespace string, nodes []selector.Node) *pb.ServiceInstancesInProto {
	ins := make([]*v1.Instance, 0, len(nodes))
	for _, node := range nodes {
		host, port, err := net.SplitHostPort(node.Address())
		if err != nil {
			log.Errorf("split host port failed error: %v", err)
			return nil
		}
		portUint64, err := strconv.ParseUint(port, 10, 32)
		if err != nil {
			log.Errorf("parse port failed error: %v", err)
			return nil
		}
		ins = append(ins, &v1.Instance{
			Id:        wrapperspb.String(node.Metadata()["merge"]),
			Service:   wrapperspb.String(node.ServiceName()),
			Namespace: wrapperspb.String(namespace),
			Host:      wrapperspb.String(host),
			Port:      wrapperspb.UInt32(uint32(portUint64)),
			Protocol:  wrapperspb.String(node.Scheme()),
			Version:   wrapperspb.String(node.Version()),
			Weight:    wrapperspb.UInt32(uint32(*node.InitialWeight())),
			Metadata:  node.Metadata(),
		})
	}

	d := &v1.DiscoverResponse{
		Code:      wrapperspb.UInt32(1),
		Info:      wrapperspb.String("ok"),
		Type:      v1.DiscoverResponse_INSTANCE,
		Service:   &v1.Service{Name: wrapperspb.String(nodes[0].ServiceName()), Namespace: wrapperspb.String(namespace)},
		Instances: ins,
	}
	return pb.NewServiceInstancesInProto(d, func(string) local.InstanceLocalValue {
		return local.NewInstanceLocalValue()
	}, &pb.SvcPluginValues{Routers: nil, Loadbalancer: nil}, nil)
}
