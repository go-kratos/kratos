package server

import (
	"context"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/app/service/openplatform/ticket-item/service"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// New 新建warden实例
func New(c *conf.Config) *warden.Server {
	//server := warden.NewServer(c.RPCServer)
	server := warden.NewServer(nil)
	server.Use(middleware())
	item.RegisterItemServer(server.Server(), service.New(c))
	item.RegisterGuestServer(server.Server(), service.New(c))
	item.RegisterBulletinServer(server.Server(), service.New(c))
	item.RegisterVenueServer(server.Server(), service.New(c))
	item.RegisterPlaceServer(server.Server(), service.New(c))
	item.RegisterAreaServer(server.Server(), service.New(c))
	item.RegisterSeatServer(server.Server(), service.New(c))
	/**go func() {
		err := server.Run(c.RPCServer.Addr)
		if err != nil {
			panic("run server failed!" + err.Error())
		}
	}()
	log.Info("warden run@%s", c.RPCServer.Addr)
	**/
	_, err := server.Start()
	if err != nil {
		panic("run server failed!" + err.Error())
	}
	return server
}

// 拦截器，作用类似于http中间件
func middleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 记录调用方法
		log.Info("method:", info.FullMethod)
		// call chain
		resp, err = handler(ctx, req)
		return
	}
}

type server struct {
	is *service.ItemService
}

func (s *server) Info(ctx context.Context, in *item.InfoRequest) (*item.InfoReply, error) {
	return s.is.Info(ctx, in)
}

func (s *server) Cards(ctx context.Context, in *item.CardsRequest) (*item.CardsReply, error) {
	return s.is.Cards(ctx, in)
}

func (s *server) BillInfo(ctx context.Context, in *item.BillRequest) (*item.BillReply, error) {
	return s.is.BillInfo(ctx, in)
}

func (s *server) GuestInfo(ctx context.Context, in *item.GuestInfoRequest) (*item.GuestInfoReply, error) {
	return s.is.GuestInfo(ctx, in)
}

func (s *server) GuestStatus(ctx context.Context, in *item.GuestStatusRequest) (*item.GuestInfoReply, error) {
	return s.is.GuestStatus(ctx, in)
}

func (s *server) BulletinInfo(ctx context.Context, in *item.BulletinInfoRequest) (*item.BulletinReply, error) {
	return s.is.BulletinInfo(ctx, in)
}

func (s *server) BulletinCheck(ctx context.Context, in *item.BulletinCheckRequest) (*item.BulletinReply, error) {
	return s.is.BulletinCheck(ctx, in)
}

func (s *server) BulletinState(ctx context.Context, in *item.BulletinStateRequest) (*item.BulletinReply, error) {
	return s.is.BulletinState(ctx, in)
}

func (s *server) Wish(ctx context.Context, in *item.WishRequest) (*item.WishReply, error) {
	return s.is.Wish(ctx, in)
}

func (s *server) Fav(ctx context.Context, in *item.FavRequest) (*item.FavReply, error) {
	return s.is.Fav(ctx, in)
}

// VenueInfo 添加/修改场馆信息
func (s *server) VenueInfo(ctx context.Context, in *item.VenueInfoRequest) (*item.VenueInfoReply, error) {
	return s.is.VenueInfo(ctx, in)
}

// PlaceInfo 添加/修改场地信息
func (s *server) PlaceInfo(ctx context.Context, in *item.PlaceInfoRequest) (*item.PlaceInfoReply, error) {
	return s.is.PlaceInfo(ctx, in)
}

// AreaInfo 添加/修改区域信息
func (s *server) AreaInfo(ctx context.Context, in *item.AreaInfoRequest) (*item.AreaInfoReply, error) {
	return s.is.AreaInfo(ctx, in)
}

// DeleteArea 删除区域信息
func (s *server) DeleteArea(ctx context.Context, in *item.DeleteAreaRequest) (*item.DeleteAreaReply, error) {
	return s.is.DeleteArea(ctx, in)
}

// SeatInfo 添加/修改座位信息
func (s *server) SeatInfo(ctx context.Context, in *item.SeatInfoRequest) (*item.SeatInfoReply, error) {
	return s.is.SeatInfo(ctx, in)
}

// SeatStock 设置座位库存
func (s *server) SeatStock(ctx context.Context, in *item.SeatStockRequest) (*item.SeatStockReply, error) {
	return s.is.SeatStock(ctx, in)
}

// RemoveSeatOrders 删除坐票票价下所有座位
func (s *server) RemoveSeatOrders(ctx context.Context, in *item.RemoveSeatOrdersRequest) (*item.RemoveSeatOrdersReply, error) {
	return s.is.RemoveSeatOrders(ctx, in)
}

// Version 添加/修改项目信息
func (s *server) Version(ctx context.Context, in *item.VersionRequest) (*item.VersionReply, error) {
	return s.is.Version(ctx, in)
}
