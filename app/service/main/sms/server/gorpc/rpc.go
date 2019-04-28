package gorpc

import (
	pb "go-common/app/service/main/sms/api"
	"go-common/app/service/main/sms/conf"
	"go-common/app/service/main/sms/model"
	"go-common/app/service/main/sms/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
	"go-common/library/net/rpc/interceptor"
)

// RPC rpc server
type RPC struct {
	s *service.Service
}

// New new rpc server.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	in := interceptor.NewInterceptor("")
	svr.Interceptor = in
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Send rpc send.
func (r *RPC) Send(c context.Context, a *model.ArgSend, res *struct{}) (err error) {
	req := &pb.SendReq{Mid: a.Mid, Mobile: a.Mobile, Country: a.Country, Tcode: a.Tcode, Tparam: a.Tparam}
	_, err = r.s.Send(c, req)
	return
}

// SendBatch rpc sendbatch.
func (r *RPC) SendBatch(c context.Context, a *model.ArgSendBatch, res *struct{}) (err error) {
	req := &pb.SendBatchReq{Mids: a.Mids, Mobiles: a.Mobiles, Tcode: a.Tcode, Tparam: a.Tparam}
	_, err = r.s.SendBatch(c, req)
	return
}
