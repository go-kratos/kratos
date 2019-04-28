package grpc

import (
	"context"

	pb "go-common/app/service/main/thumbup/api"
	"go-common/app/service/main/thumbup/model"
	"go-common/app/service/main/thumbup/service"
	"go-common/library/net/rpc/warden"

	"github.com/golang/protobuf/ptypes/empty"
)

// New Coin warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterThumbupServer(ws.Server(), &server{s: svr})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	s *service.Service
}

func (r server) Like(c context.Context, req *pb.LikeReq) (reply *pb.LikeReply, err error) {
	res, err := r.s.Like(c, req.Business, req.Mid, req.OriginID, req.MessageID, int8(req.Action), req.UpMid)
	reply = &pb.LikeReply{
		OriginID:      res.OriginID,
		MessageID:     res.ID,
		LikeNumber:    res.Likes,
		DislikeNumber: res.Dislikes,
	}
	return
}

func (r server) Stats(c context.Context, req *pb.StatsReq) (reply *pb.StatsReply, err error) {
	res, err := r.s.StatsWithLike(c, req.Business, req.Mid, req.OriginID, req.MessageIds)
	reply = &pb.StatsReply{Stats: map[int64]*pb.StatState{}}
	for name, item := range res {
		reply.Stats[name] = &pb.StatState{
			OriginID:      item.OriginID,
			MessageID:     item.ID,
			LikeNumber:    item.Likes,
			DislikeNumber: item.Dislikes,
			LikeState:     pb.State(item.LikeState),
		}
	}
	return
}

func (r server) MultiStats(c context.Context, req *pb.MultiStatsReq) (reply *pb.MultiStatsReply, err error) {
	arg := &model.MultiBusiness{
		Mid:        req.Mid,
		Businesses: make(map[string][]*model.MultiBusinessItem),
	}
	for name, b := range req.Business {
		for _, i := range b.Records {
			arg.Businesses[name] = append(arg.Businesses[name], &model.MultiBusinessItem{
				OriginID:  i.OriginID,
				MessageID: i.MessageID,
			})
		}
	}
	res, err := r.s.MultiStatsWithLike(c, arg)
	reply = &pb.MultiStatsReply{}
	if res != nil {
		reply.Business = make(map[string]*pb.MultiStatsReply_Records)
		for k, v := range res {
			items := &pb.MultiStatsReply_Records{
				Records: make(map[int64]*pb.StatState),
			}
			for id, state := range v {
				items.Records[id] = &pb.StatState{
					OriginID:      state.OriginID,
					MessageID:     state.ID,
					LikeNumber:    state.Likes,
					DislikeNumber: state.Dislikes,
					LikeState:     pb.State(state.LikeState),
				}
			}
			reply.Business[k] = items
		}
	}
	return
}

func (r server) HasLike(c context.Context, req *pb.HasLikeReq) (reply *pb.HasLikeReply, err error) {
	_, res, err := r.s.HasLike(c, req.Business, req.Mid, req.MessageIds)
	reply = &pb.HasLikeReply{States: res}
	return
}

func (r server) UserLikes(c context.Context, req *pb.UserLikesReq) (reply *pb.UserLikesReply, err error) {
	res, err := r.s.UserTotalLike(c, req.Business, req.Mid, int(req.Pn), int(req.Ps))
	reply = &pb.UserLikesReply{}
	if res != nil {
		reply.Total = int64(res.Total)
		for _, item := range res.List {
			reply.Items = append(reply.Items, &pb.ItemRecord{
				MessageID: item.MessageID,
				Time:      item.Time,
			})
		}
	}
	return
}

func (r server) ItemLikes(c context.Context, req *pb.ItemLikesReq) (reply *pb.ItemLikesReply, err error) {
	res, err := r.s.ItemLikes(c, req.Business, req.OriginID, req.MessageID, int(req.Pn), int(req.Ps), req.LastMid)
	reply = &pb.ItemLikesReply{}
	for _, item := range res {
		reply.Users = append(reply.Users, &pb.UserRecord{
			Mid:  item.Mid,
			Time: item.Time,
		})
	}
	return
}

func (r server) UpdateCount(c context.Context, req *pb.UpdateCountReq) (reply *empty.Empty, err error) {
	reply = &empty.Empty{}
	err = r.s.UpdateCount(c, req.Business, req.OriginID, req.MessageID, req.LikeChange, req.DislikeChange, req.IP, req.Operator)
	return
}

func (r server) RawStat(c context.Context, req *pb.RawStatReq) (reply *pb.RawStatReply, err error) {
	res, err := r.s.RawStats(c, req.Business, req.OriginID, req.MessageID)
	reply = &pb.RawStatReply{
		OriginID:      res.OriginID,
		MessageID:     res.ID,
		LikeNumber:    res.Likes,
		DislikeNumber: res.Dislikes,
		LikeChange:    res.LikesChange,
		DislikeChange: res.DislikesChange,
	}
	return
}
