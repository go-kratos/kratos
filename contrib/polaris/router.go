package polaris

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"
	"github.com/polarismesh/polaris-go/pkg/model/local"
	"github.com/polarismesh/polaris-go/pkg/model/pb"
	v1 "github.com/polarismesh/polaris-go/pkg/model/pb/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
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
		req := &polaris.ProcessRoutersRequest{
			ProcessRoutersRequest: model.ProcessRoutersRequest{
				SourceService: model.ServiceInfo{Namespace: p.namespace, Service: o.service},
				DstInstances:  buildPolarisInstance(p.namespace, nodes),
			},
		}
		if appInfo, ok := kratos.FromContext(ctx); ok {
			req.SourceService.Service = appInfo.Name()
		}

		req.AddArguments(model.BuildCallerServiceArgument(p.namespace, req.ProcessRoutersRequest.SourceService.Service))

		// process transport
		if tr, ok := transport.FromClientContext(ctx); ok {
			req.AddArguments(model.BuildMethodArgument(tr.Operation()))
			req.AddArguments(model.BuildPathArgument(tr.Operation()))

			for _, key := range tr.RequestHeader().Keys() {
				req.AddArguments(model.BuildHeaderArgument(strings.ToLower(key), tr.RequestHeader().Get(key)))
			}

			// http
			if ht, ok := tr.(http.Transporter); ok {
				req.AddArguments(model.BuildPathArgument(ht.Request().URL.Path))
				req.AddArguments(model.BuildCallerIPArgument(ht.Request().RemoteAddr))

				// cookie
				for _, cookie := range ht.Request().Cookies() {
					req.AddArguments(model.BuildCookieArgument(cookie.Name, cookie.Value))
				}

				// url query
				for key, values := range ht.Request().URL.Query() {
					req.AddArguments(model.BuildQueryArgument(key, strings.Join(values, ",")))
				}
			}
		}

		n := make(map[string]selector.Node, len(nodes))
		for _, node := range nodes {
			n[node.Address()] = node
		}

		m, err := p.router.ProcessRouters(req)
		if err != nil {
			log.Errorf("polaris process routers failed, err=%v", err)
			return nodes
		}

		newNode := make([]selector.Node, 0, len(m.Instances))
		for _, ins := range m.GetInstances() {
			if v, ok := n[fmt.Sprintf("%s:%d", ins.GetHost(), ins.GetPort())]; ok {
				newNode = append(newNode, v)
			}
		}
		if len(newNode) == 0 {
			return nodes
		}
		return newNode
	}
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
		Service:   &v1.Service{Name: wrapperspb.String(nodes[0].ServiceName()), Namespace: wrapperspb.String("default")},
		Instances: ins,
	}
	return pb.NewServiceInstancesInProto(d, func(string) local.InstanceLocalValue {
		return local.NewInstanceLocalValue()
	}, &pb.SvcPluginValues{Routers: nil, Loadbalancer: nil}, nil)
}
