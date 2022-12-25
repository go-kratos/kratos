package polaris

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NodeFilter polaris dynamic router selector
func (p *Polaris) NodeFilter(service, namespace string) (selector.NodeFilter, error) {
	watcher, err := newWatcher(context.Background(), namespace, service, p.discovery)
	// FIXME: 是否需要返回 err, 还是后续不生效即可？
	if err != nil {
		return nil, err
	}
	go func() {
		// 缓存北极星数据
		for {
			_, err = watcher.Next()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Errorf("failed to watch polaris discovery endpoint: %v", err)
				time.Sleep(time.Second)
				continue
			}
		}
	}()

	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		if len(nodes) == 0 {
			return nodes
		}
		if appInfo, ok := kratos.FromContext(ctx); ok {
			req := &polaris.ProcessRoutersRequest{
				ProcessRoutersRequest: model.ProcessRoutersRequest{
					SourceService: model.ServiceInfo{
						Service:   appInfo.Name(),
						Namespace: namespace,
					},
					DstInstances: watcher.service,
				},
			}
			req.AddArguments(model.BuildCallerServiceArgument(appInfo.Name(), namespace))

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
	}, nil
}
