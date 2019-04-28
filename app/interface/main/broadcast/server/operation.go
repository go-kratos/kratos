package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	iModel "go-common/app/interface/main/broadcast/model"
	pb "go-common/app/service/main/broadcast/api/grpc/v1"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster/render"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	_apiConnect     = "/x/broadcast/conn/connect"
	_apiDisconnect  = "/x/broadcast/conn/disconnect"
	_apiHeartbeat   = "/x/broadcast/conn/heartbeat"
	_apiRenewOnline = "/x/broadcast/online/renew"
)

func (s *Server) failedForword(ctx context.Context, url string, req, reply proto.Message) (err error) {
	var (
		b       []byte
		httpReq *http.Request
		res     = new(render.PB)
		api     = fmt.Sprintf("%s%s?token=%s", s.c.Broadcast.APIHost, url, s.c.Broadcast.APIToken)
	)
	if b, err = proto.Marshal(req); err != nil {
		return
	}
	if httpReq, err = http.NewRequest("POST", api, bytes.NewBuffer(b)); err != nil {
		return
	}
	if err = s.httpCli.PB(ctx, httpReq, res); err != nil {
		return
	}
	if int(res.Code) != ecode.OK.Code() {
		err = ecode.Int(int(res.Code))
		return
	}
	err = types.UnmarshalAny(res.Data, reply)
	return
}

// Connect .
func (s *Server) Connect(ctx context.Context, p *model.Proto, cookie string) (mid int64, key, rid, platform string, accepts []int32, err error) {
	var (
		req = &pb.ConnectReq{
			Server:    s.serverID,
			ServerKey: s.NextKey(),
			Cookie:    cookie,
			Token:     p.Body,
		}
		reply *pb.ConnectReply
	)
	if !s.c.Broadcast.Failover {
		reply, err = s.rpcClient.Connect(ctx, req)
	}
	if s.c.Broadcast.Failover || err != nil {
		reply = new(pb.ConnectReply)
		if err = s.failedForword(ctx, _apiConnect, req, reply); err != nil {
			return
		}
	}
	return reply.Mid, reply.Key, reply.RoomID, reply.Platform, reply.Accepts, nil
}

// Disconnect .
func (s *Server) Disconnect(ctx context.Context, mid int64, key string) (err error) {
	var (
		req = &pb.DisconnectReq{
			Mid:    mid,
			Server: s.serverID,
			Key:    key,
		}
		reply *pb.DisconnectReply
	)
	if !s.c.Broadcast.Failover {
		reply, err = s.rpcClient.Disconnect(ctx, req)
	}
	if s.c.Broadcast.Failover || err != nil {
		reply = new(pb.DisconnectReply)
		if err = s.failedForword(ctx, _apiDisconnect, req, reply); err != nil {
			return
		}
	}
	return
}

// Heartbeat .
func (s *Server) Heartbeat(ctx context.Context, mid int64, key string) (err error) {
	var (
		req = &pb.HeartbeatReq{
			Mid:    mid,
			Server: s.serverID,
			Key:    key,
		}
		reply *pb.HeartbeatReply
	)
	if !s.c.Broadcast.Failover {
		reply, err = s.rpcClient.Heartbeat(ctx, req)
	}
	if s.c.Broadcast.Failover || err != nil {
		reply = new(pb.HeartbeatReply)
		if err = s.failedForword(ctx, _apiHeartbeat, req, reply); err != nil {
			return
		}
	}
	return
}

// RenewOnline .
func (s *Server) RenewOnline(ctx context.Context, serverID string, shard int32, rommCount map[string]int32) (allRoom map[string]int32, err error) {
	var (
		req = &pb.OnlineReq{
			Server:    s.serverID,
			RoomCount: rommCount,
			Sharding:  shard,
		}
		reply *pb.OnlineReply
	)
	if !s.c.Broadcast.Failover {
		for r := 0; r < s.c.Broadcast.OnlineRetries; r++ {
			if reply, err = s.rpcClient.RenewOnline(ctx, req, grpc.UseCompressor(gzip.Name)); err != nil {
				time.Sleep(s.backoff.Backoff(r))
				continue
			}
			break
		}
	}
	if s.c.Broadcast.Failover || err != nil {
		reply = new(pb.OnlineReply)
		if err = s.failedForword(ctx, _apiRenewOnline, req, reply); err != nil {
			return
		}
	}
	return reply.RoomCount, nil
}

// Report .
func (s *Server) Report(mid int64, proto *model.Proto) (rp *model.Proto, err error) {
	var (
		reply *pb.ReceiveReply
	)
	if reply, err = s.rpcClient.Receive(context.Background(), &pb.ReceiveReq{
		Mid:   mid,
		Proto: proto,
	}); err != nil {
		return
	}
	return reply.Proto, nil
}

// Operate .
func (s *Server) Operate(p *model.Proto, ch *Channel, b *Bucket) error {
	var err error
	switch {
	case p.Operation >= model.MinBusinessOp && p.Operation <= model.MaxBusinessOp:
		_, err = s.Report(ch.Mid, p)
		if err != nil {
			log.Error("s.Reprot(%d,%v) error(%v)", ch.Mid, p, err)
			return nil
		}
		p.Body = nil
		// ignore down message
	case p.Operation == model.OpChangeRoom:
		p.Operation = model.OpChangeRoomReply
		var req iModel.ChangeRoomReq
		if err = json.Unmarshal(p.Body, &req); err == nil {
			if err = b.ChangeRoom(req.RoomID, ch); err == nil {
				p.Body = iModel.Message(map[string]interface{}{"room_id": string(p.Body)}, nil)
			}
		}
	case p.Operation == model.OpRegister:
		p.Operation = model.OpRegisterReply
		var req iModel.RegisterOpReq
		if err = json.Unmarshal(p.Body, &req); err == nil {
			if len(req.Operations) > 0 {
				ch.Watch(req.Operations...)
				p.Body = iModel.Message(map[string]interface{}{"operations": req.Operations}, nil)
			} else {
				ch.Watch(req.Operation)
				p.Body = iModel.Message(map[string]interface{}{"operation": req.Operation}, nil)
			}
		}
	case p.Operation == model.OpUnregister:
		p.Operation = model.OpUnregisterReply
		var req iModel.UnregisterOpReq
		if err = json.Unmarshal(p.Body, &req); err == nil {
			if len(req.Operations) > 0 {
				ch.UnWatch(req.Operations...)
				p.Body = iModel.Message(map[string]interface{}{"operations": req.Operations}, nil)
			} else {
				ch.UnWatch(req.Operation)
				p.Body = iModel.Message(map[string]interface{}{"operation": req.Operation}, nil)
			}
		}
	default:
		err = ErrOperation
	}
	if err != nil {
		log.Error("Operate (%+v) failed!err:=%v", p, err)
		p.Body = iModel.Message(nil, err)
	}
	return nil
}
