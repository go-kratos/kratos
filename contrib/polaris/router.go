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

// NodeFilter polaris dynamic router selector
func (p *Polaris) NodeFilter() selector.NodeFilter {
	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		if len(nodes) == 0 {
			return nodes
		}
		if appInfo, ok := kratos.FromContext(ctx); ok {
			req := &polaris.ProcessRoutersRequest{
				ProcessRoutersRequest: model.ProcessRoutersRequest{
					SourceService: model.ServiceInfo{
						Service:   appInfo.Name(),
						Namespace: p.namespace,
					},
					DstInstances: buildPolarisInstance(p.namespace, nodes),
				},
			}
			req.AddArguments(model.BuildCallerServiceArgument(p.namespace, appInfo.Name()))

			// process transport
			if tr, ok := transport.FromServerContext(ctx); ok {
				req.AddArguments(model.BuildMethodArgument(tr.Operation()))
				req.AddArguments(model.BuildPathArgument(tr.Operation()))

				for _, key := range tr.RequestHeader().Keys() {
					req.AddArguments(model.BuildHeaderArgument(key, tr.RequestHeader().Get(key)))
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
			return newNode
		}
		return nodes
	}
}

func buildPolarisInstance(namespace string, nodes []selector.Node) *pb.ServiceInstancesInProto {
	ins := make([]*v1.Instance, 0, len(nodes))
	for _, node := range nodes {
		host, port, err := net.SplitHostPort(node.Address())
		if err != nil {
			return nil
		}
		portInt, err := strconv.Atoi(port)
		if err != nil {
			return nil
		}
		ins = append(ins, &v1.Instance{
			Id:        wrapperspb.String(node.Metadata()["merge"]),
			Service:   wrapperspb.String(node.ServiceName()),
			Namespace: wrapperspb.String(namespace),
			Host:      wrapperspb.String(host),
			Port:      wrapperspb.UInt32(uint32(portInt)),
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
	return pb.NewServiceInstancesInProto(d, func(s string) local.InstanceLocalValue {
		return local.NewInstanceLocalValue()
	}, &pb.SvcPluginValues{Routers: nil, Loadbalancer: nil}, nil)
}
