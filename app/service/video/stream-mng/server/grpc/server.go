// Package server generate by warden_gen
package grpc

import (
	"context"

	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-common/app/service/video/stream-mng/api/v1"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/service"
	"go-common/library/ecode"
	"go-common/library/log"
	nmd "go-common/library/net/metadata"
	"go-common/library/net/rpc/warden"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// New Stream warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	//ws := warden.NewServer(c, grpc.MaxRecvMsgSize(32*1024*1024), grpc.MaxSendMsgSize(32*1024*1024)) 这里需要考虑配置问题
	ws := warden.NewServer(c)
	ws.Use(middleware())
	v1.RegisterStreamServer(ws.Server(), &server{svr})

	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}

	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.StreamServer = &server{}

// GetSingleScreeShotByRoomID
func (s *server) GetSingleScreeShot(ctx context.Context, req *v1.GetSingleScreeShotReq) (*v1.GetSingleScreeShotReply, error) {
	roomID := req.RoomId

	start := req.StartTime
	end := req.EndTime
	channel := req.Channel

	resp := &v1.GetSingleScreeShotReply{}
	if roomID <= 0 || start == "" || end == "" {
		st, _ := ecode.Error(ecode.RequestErr, "some fields are empty").WithDetails(resp)
		return nil, st
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, "Start time format is incorrect").WithDetails(resp)
		return nil, st
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, "End time format is incorrect").WithDetails(resp)
		return nil, st
	}

	info, err := s.svr.GetSingleScreeShot(ctx, roomID, startTime.Unix(), endTime.Unix(), channel)
	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	return &v1.GetSingleScreeShotReply{
		List: info,
	}, nil

}

// GetMultiScreenShotByRommID
func (s *server) GetMultiScreenShot(ctx context.Context, req *v1.GetMultiScreenShotReq) (*v1.GetMultiScreenShotReply, error) {
	rooms := req.RoomIds
	ts := req.Ts
	channel := req.Channel

	resp := &v1.GetMultiScreenShotReply{}
	if rooms == "" || ts == 0 {
		st, _ := ecode.Error(ecode.RequestErr, "some fields are empty").WithDetails(resp)
		return nil, st
	}

	// 切割room_id
	roomIDs := strings.Split(rooms, ",")

	if len(roomIDs) <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_ids is not right").WithDetails(resp)
		return nil, st
	}

	res := v1.GetMultiScreenShotReply{
		List: map[int64]string{},
	}

	rids := []int64{}
	for _, v := range roomIDs {
		roomID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}

		rids = append(rids, roomID)
	}

	urls, err := s.svr.GetMultiScreenShot(ctx, rids, ts, channel)

	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	res.List = urls

	return &res, nil
}

// GetOriginScreenShotPic
func (s *server) GetOriginScreenShotPic(ctx context.Context, req *v1.GetOriginScreenShotPicReq) (*v1.GetOriginScreenShotPicReply, error) {
	rooms := req.RoomIds
	ts := req.Ts

	resp := &v1.GetOriginScreenShotPicReply{}
	if rooms == "" || ts == 0 {
		st, _ := ecode.Error(ecode.RequestErr, "some fields are empty").WithDetails(resp)
		return nil, st
	}

	// 切割room_id
	roomIDs := strings.Split(rooms, ",")

	if len(roomIDs) <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_ids is not right").WithDetails(resp)
		return nil, st
	}

	res := v1.GetOriginScreenShotPicReply{
		List: map[int64]string{},
	}

	rids := []int64{}
	for _, v := range roomIDs {
		roomID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}

		rids = append(rids, roomID)
	}

	urls, err := s.svr.GetOriginScreenShotPic(ctx, rids, ts)

	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	res.List = urls

	return &res, nil
}

// CreateOfficeStream 创建正式流
func (s *server) CreateOfficalStream(ctx context.Context, req *v1.CreateOfficalStreamReq) (*v1.CreateOfficalStreamReply, error) {
	key := req.Key
	streamName := req.StreamName
	if req.Uid != 0 {
		key = mockStreamKey(fmt.Sprintf("%d", req.Uid))
		streamName = mockStreamName(fmt.Sprintf("%d", req.Uid))
	}

	resp := &v1.CreateOfficalStreamReply{}
	// 检查参数
	if streamName == "" || key == "" || req.RoomId <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "some fields are empty").WithDetails(resp)
		return nil, st
	}

	flag := s.svr.CreateOfficalStream(ctx, streamName, key, req.RoomId)

	return &v1.CreateOfficalStreamReply{
		Success: flag,
	}, nil
}

// GetStreamInfo 获取单个流信息
func (s *server) GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (*v1.GetStreamInfoReply, error) {
	rid := req.RoomId
	sname := req.StreamName

	resp := &v1.GetStreamInfoReply{}
	if rid == 0 && sname == "" {
		resp.Code = -400
		resp.Message = "some fields are empty"
		return resp, nil
	}

	info, err := s.svr.GetStreamInfo(ctx, int64(rid), sname)

	if err != nil {
		resp.Code = -400
		resp.Message = err.Error()
		return resp, nil
	}

	baseList := []*v1.StreamBase{}
	for _, v := range info.List {
		forward := []uint32{}
		for _, f := range v.Forward {
			forward = append(forward, uint32(f))
		}

		baseList = append(baseList, &v1.StreamBase{
			StreamName:      v.StreamName,
			DefaultUpstream: uint32(v.DefaultUpStream),
			Origin:          uint32(v.Origin),
			Forward:         forward,
			Type:            uint32(v.Type),
			Options:         uint32(v.Options),
			//Key:             v.Key,
		})
	}

	resp.Code = 0
	resp.Data = &v1.StreamFullInfo{
		RoomId: uint32(info.RoomID),
		Hot:    uint32(info.Hot),
		List:   baseList,
	}
	return resp, nil
}

// GetMultiStreamInfo 批量获取流信息
func (s *server) GetMultiStreamInfo(ctx context.Context, req *v1.GetMultiStreamInfoReq) (*v1.GetMultiStreamInfoReply, error) {
	rids := req.RoomIds

	resp := &v1.GetMultiStreamInfoReply{}

	// 切割room_id

	if len(rids) <= 0 {
		resp.Code = 0
		resp.Message = "success"
		return resp, nil
	}

	if len(rids) > 30 {
		resp.Code = -400
		resp.Message = "The number of rooms must be less than 30"
		return resp, nil
	}

	roomIDs := []int64{}
	for _, v := range rids {
		roomID := int64(v)
		if roomID <= 0 {
			continue
		}

		roomIDs = append(roomIDs, roomID)
	}

	info, err := s.svr.GetMultiStreamInfo(ctx, roomIDs)

	if err != nil {
		log.Infov(ctx, log.KV("log", err.Error()))
		resp.Code = 0
		resp.Message = "success"
		return resp, nil
	}

	if info == nil || len(info) == 0 {
		log.Infov(ctx, log.KV("log", "can find any things"))
		resp.Code = 0
		resp.Message = "success"
		return resp, nil
	}

	res := map[uint32]*v1.StreamFullInfo{}
	for id, v := range info {
		item := &v1.StreamFullInfo{}
		item.Hot = uint32(v.Hot)
		item.RoomId = uint32(v.RoomID)

		baseList := []*v1.StreamBase{}
		for _, i := range v.List {
			forward := []uint32{}
			for _, f := range i.Forward {
				forward = append(forward, uint32(f))
			}

			baseList = append(baseList, &v1.StreamBase{
				StreamName:      i.StreamName,
				DefaultUpstream: uint32(i.DefaultUpStream),
				Origin:          uint32(i.Origin),
				Forward:         forward,
				Type:            uint32(i.Type),
				Options:         uint32(i.Options),
				//Key:             i.Key,
			})
		}

		item.List = baseList

		res[uint32(id)] = item
	}

	resp.Code = 0
	resp.Data = res
	return resp, nil
}

// ChangeSrc 切换cdn
func (s *server) ChangeSrc(ctx context.Context, req *v1.ChangeSrcReq) (*v1.EmptyStruct, error) {
	resp := &v1.EmptyStruct{}
	if req.RoomId <= 0 || req.Src == 0 || req.Source == "" || req.OperateName == "" {
		st, _ := ecode.Error(ecode.RequestErr, "some fields are empty").WithDetails(resp)
		return nil, st
	}

	// todo 后续改为新的src
	src := int8(req.Src)

	if _, ok := common.SrcMapBitwise[src]; !ok {
		st, _ := ecode.Error(ecode.RequestErr, "src is not right").WithDetails(resp)
		return nil, st
	}

	err := s.svr.ChangeSrc(ctx, req.RoomId, common.SrcMapBitwise[src], req.Source, req.OperateName, req.Reason)

	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	return resp, nil
}

// GetStreamLastTime 得到流到最后推流时间;主流的推流时间up_rank = 1
func (s *server) GetStreamLastTime(ctx context.Context, req *v1.GetStreamLastTimeReq) (*v1.GetStreamLastTimeReply, error) {
	rid := req.RoomId

	resp := &v1.GetStreamLastTimeReply{}
	if rid <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_id is not right").WithDetails(resp)
		return nil, st
	}

	t, err := s.svr.GetStreamLastTime(ctx, rid)
	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	return &v1.GetStreamLastTimeReply{
		LastTime: t,
	}, nil
}

// GetStreamNameByRoomID 需要考虑备用流 + 考虑短号
func (s *server) GetStreamNameByRoomID(ctx context.Context, req *v1.GetStreamNameByRoomIDReq) (*v1.GetStreamNameByRoomIDReply, error) {
	rid := req.RoomId

	resp := &v1.GetStreamNameByRoomIDReply{}

	if rid <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_id is not right").WithDetails(resp)
		return nil, st
	}

	res, err := s.svr.GetStreamNameByRoomID(ctx, rid, false)
	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	if len(res) == 0 {
		st, _ := ecode.Error(ecode.RequestErr, fmt.Sprintf("can not find info by room_id=%d", rid)).WithDetails(resp)
		return nil, st
	}

	return &v1.GetStreamNameByRoomIDReply{
		StreamName: res[0],
	}, nil
}

// GetRoomIDByStreamName 查询房间号
func (s *server) GetRoomIDByStreamName(ctx context.Context, req *v1.GetRoomIDByStreamNameReq) (*v1.GetRoomIDByStreamNameReply, error) {
	resp := &v1.GetRoomIDByStreamNameReply{}
	if req.StreamName == "" {
		st, _ := ecode.Error(ecode.RequestErr, "stream name is empty").WithDetails(resp)
		return nil, st
	}

	res, err := s.svr.GetRoomIDByStreamName(ctx, req.StreamName)
	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	if res <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, fmt.Sprintf("can not find any info by name = %s", req.StreamName)).WithDetails(resp)
		return nil, st
	}

	return &v1.GetRoomIDByStreamNameReply{
		RoomId: res,
	}, nil
}

// GetAdapterStreamByStreamName 适配结果输出, 此处也可以输入备用流， 该结果只输出直推上行
func (s *server) GetAdapterStreamByStreamName(ctx context.Context, req *v1.GetAdapterStreamByStreamNameReq) (*v1.GetAdapterStreamByStreamNameReply, error) {
	res := v1.GetAdapterStreamByStreamNameReply{
		List: map[string]*v1.AdapterStream{},
	}
	snames := req.StreamNames
	if snames == "" {
		st, _ := ecode.Error(ecode.RequestErr, "stream_names is empty").WithDetails(&res)
		return nil, st
	}

	nameSlice := strings.Split(snames, ",")
	if len(nameSlice) == 0 {
		st, _ := ecode.Error(ecode.RequestErr, "stream_names is empty").WithDetails(&res)
		return nil, st
	}

	if len(nameSlice) > 500 {
		st, _ := ecode.Error(ecode.RequestErr, "too many names").WithDetails(&res)
		return nil, st
	}

	info := s.svr.GetAdapterStreamByStreamName(ctx, nameSlice)

	if info != nil {
		for name, v := range info {
			res.List[name] = &v1.AdapterStream{
				Src:     v.Src,
				RoomId:  v.RoomID,
				UpRank:  v.UpRank,
				SrcName: v.SrcName,
			}
		}
	}

	return &res, nil
}

// GetSrcByRoomID
func (s *server) GetSrcByRoomID(ctx context.Context, req *v1.GetSrcByRoomIDReq) (*v1.GetSrcByRoomIDReply, error) {
	rid := req.RoomId

	resp := &v1.GetSrcByRoomIDReply{}
	if rid <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_id is not right").WithDetails(resp)
		return nil, st
	}

	info, err := s.svr.GetSrcByRoomID(ctx, rid)

	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	if info == nil || len(info) == 0 {
		st, _ := ecode.Error(ecode.RequestErr, "获取线路失败").WithDetails(resp)
		return nil, st
	}

	res := &v1.GetSrcByRoomIDReply{
		List: []*v1.RoomSrcCheck{},
	}

	for _, v := range info {
		res.List = append(res.List, &v1.RoomSrcCheck{
			Src:     v.Src,
			Checked: int32(v.Checked),
			Desc:    v.Desc,
		})
	}
	return res, nil
}

// GetLineListByRoomID
func (s *server) GetLineListByRoomID(ctx context.Context, req *v1.GetLineListByRoomIDReq) (*v1.GetLineListByRoomIDReply, error) {
	resp := &v1.GetLineListByRoomIDReply{}

	if req.RoomId <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_id is not right").WithDetails(resp)
		return nil, st
	}

	info, err := s.svr.GetLineListByRoomID(ctx, req.RoomId)

	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	if info == nil || len(info) == 0 {
		st, _ := ecode.Error(ecode.RequestErr, "获取线路失败").WithDetails(resp)
		return nil, st
	}

	res := &v1.GetLineListByRoomIDReply{
		List: []*v1.LineList{},
	}

	for _, v := range info {
		res.List = append(res.List, &v1.LineList{
			Src:  v.Src,
			Use:  v.Use,
			Desc: v.Desc,
		})
	}
	return res, nil
}

// GetUpStreamRtmp UpStream
func (s *server) GetUpStreamRtmp(ctx context.Context, req *v1.GetUpStreamRtmpReq) (*v1.GetUpStreamRtmpReply, error) {
	resp := &v1.GetUpStreamRtmpReply{}
	if req.RoomId == 0 || req.Platform == "" {
		st, _ := ecode.Error(ecode.RequestErr, "some fields are empty").WithDetails(resp)
		return nil, st
	}

	if req.Ip == "" {
		if cmd, ok := nmd.FromContext(ctx); ok {
			if ip, ok := cmd[nmd.RemoteIP].(string); ok {
				req.Ip = ip
			}
		}
	}

	info, err := s.svr.GetUpStreamRtmp(ctx, req.RoomId, req.FreeFlow, req.Ip, req.AreaId, int(req.Attentions), 0, req.Platform)

	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	if info != nil {
		resp.Up = &v1.UpStreamRtmp{
			Addr:    info.Addr,
			Code:    info.Code,
			NewLink: info.NewLink,
		}
	}
	return resp, nil
}

// StreamCut 切流的房间和时间, 内部调用
func (s *server) StreamCut(ctx context.Context, req *v1.StreamCutReq) (*v1.EmptyStruct, error) {
	roomID := req.RoomId
	cutTime := req.CutTime

	resp := &v1.EmptyStruct{}
	if roomID <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_ids is not right").WithDetails(resp)
		return nil, st
	}

	if cutTime == 0 {
		cutTime = 1
	}

	s.svr.StreamCut(ctx, roomID, cutTime, 0)

	return &v1.EmptyStruct{}, nil
}

// Ping Service
func (s *server) Ping(ctx context.Context, req *v1.PingReq) (*v1.PingReply, error) {
	return &v1.PingReply{}, nil
}

// Close Service
func (s *server) Close(ctx context.Context, req *v1.CloseReq) (*v1.CloseReply, error) {
	return &v1.CloseReply{}, nil
}

// ClearStreamStatus
func (s *server) ClearStreamStatus(ctx context.Context, req *v1.ClearStreamStatusReq) (*v1.EmptyStruct, error) {
	rid := req.RoomId

	resp := &v1.EmptyStruct{}
	if rid <= 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_ids is not right").WithDetails(resp)
		return nil, st
	}

	err := s.svr.ClearStreamStatus(ctx, rid)
	if err != nil {
		st, _ := ecode.Error(ecode.RequestErr, err.Error()).WithDetails(resp)
		return nil, st
	}

	return &v1.EmptyStruct{}, nil
}

// CheckLiveStreamList
func (s *server) CheckLiveStreamList(ctx context.Context, req *v1.CheckLiveStreamReq) (*v1.CheckLiveStreamReply, error) {
	resp := &v1.CheckLiveStreamReply{}

	rids := req.RoomId
	if len(rids) == 0 {
		st, _ := ecode.Error(ecode.RequestErr, "room_ids is empty").WithDetails(resp)
		return nil, st
	}

	res := s.svr.CheckLiveStreamList(ctx, rids)

	resp.List = res
	return resp, nil
}

// mockStream 模拟生成的流名
func mockStreamName(uid string) string {
	num := rand.Int63n(88888888)
	return fmt.Sprintf("live_%s_%d", uid, num+1111111)
}

// mockStreamKey 模拟生成的key
func mockStreamKey(uid string) string {
	str := fmt.Sprintf("nvijqwopW1%s%d", uid, time.Now().Unix())
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	md5Str := hex.EncodeToString(cipherStr)
	return md5Str
}

func middleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		out := ""
		if err != nil {
			out = err.Error()
		} else {
			jo, _ := json.Marshal(resp)
			out = string(jo)
		}

		// 记录调用方法
		log.Infov(ctx,
			log.KV("path", info.FullMethod),
			log.KV("caller", nmd.String(ctx, nmd.Caller)),
			log.KV("input_params", fmt.Sprintf("%s", req)),
			log.KV("output_data", out),
		)
		return
	}
}
