// Package server generate by warden_gen
package grpc

import (
	"context"
	"fmt"

	pb "go-common/app/service/main/location/api"
	"go-common/app/service/main/location/model"
	"go-common/app/service/main/location/service"
	"go-common/library/net/rpc/warden"
)

// New Location warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterLocationServer(ws.Server(), &server{svr})
	ws, err := ws.Start()
	if err != nil {
		panic(fmt.Sprintf("start warden server fail! %s", err))
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ pb.LocationServer = &server{}

// Info get ip info.
func (s *server) Info(c context.Context, arg *pb.InfoReq) (res *pb.InfoReply, err error) {
	var ipinfo *model.Info
	if ipinfo, err = s.svr.Info(c, arg.Addr); err != nil {
		return
	}
	res = &pb.InfoReply{
		Addr:      ipinfo.Addr,
		Country:   ipinfo.Country,
		Province:  ipinfo.Province,
		City:      ipinfo.City,
		Isp:       ipinfo.ISP,
		Latitude:  ipinfo.Latitude,
		Longitude: ipinfo.Longitude,
		ZoneId:    ipinfo.ZoneID}
	return
}
