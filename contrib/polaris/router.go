package polaris

import (
	"context"

	"github.com/polarismesh/polaris-go/pkg/model"
	"github.com/polarismesh/polaris-go/pkg/model/pb"
	"github.com/polarismesh/polaris-go/pkg/model/pb/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/selector"
)

// NodeFilter polaris dynamic router selector
func (p *Polaris) NodeFilter(namespace string) selector.NodeFilter {
	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		if appInfo, ok := kratos.FromContext(ctx); ok {
			req := model.ProcessRoutersRequest{
				SourceService: model.ServiceInfo{
					Service:   appInfo.Name(),
					Namespace: namespace,
				},
			}
		}
		return nodes
	}
}

func newPolarisServiceInstance(nodes []selector.Node, namespace string) {
	d := v1.DiscoverResponse{
		Code: wrapperspb.UInt32(0),
		Info: wrapperspb.String(""),
		Type: 1,
		Service: &v1.Service{
			Name:      wrapperspb.String(nodes[0].ServiceName()),
			Namespace: wrapperspb.String(namespace),
			Metadata:  nodes[0].Metadata(),
		},
		Instances: nil,
		Routing:   nil,
		RateLimit: nil,
		Services:  nil,
	}
	pb.NewServiceInstancesInProto()
}
