package v1

import (
	"context"
	"strings"
	"time"

	v1pb "go-common/app/service/live/rtc/api/v1"
	"go-common/app/service/live/rtc/internal/conf"
	"go-common/app/service/live/rtc/internal/dao"
	"go-common/app/service/live/rtc/internal/model"
	"go-common/app/service/live/rtc/internal/service"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

// RtcService struct
type RtcService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
	redis.Conn
	tool     *service.Tool
	dispatch *service.Dispatcher
}

//NewRtcService init
func NewRtcService(c *conf.Config) (s *RtcService) {
	s = &RtcService{
		conf: c,
		dao:  dao.New(c),
		tool: service.NewTool(),
	}
	return s
}

// JoinChannel implementation
func (s *RtcService) JoinChannel(ctx context.Context, req *v1pb.JoinChannelRequest) (resp *v1pb.JoinChannelResponse, err error) {
	resp = &v1pb.JoinChannelResponse{}

	call := &model.RtcCall{
		UserID:    req.UserId,
		ChannelID: req.ChannelId,
		Version:   req.ProtoVersion,
		Token:     s.tool.RandomString(4),
		Status:    0,
		JoinTime:  time.Now(),
		LeaveTime: time.Unix(0, 0),
	}

	callID, err := s.dao.CreateCall(ctx, call)
	if err != nil {
		return nil, err
	}
	resp.CallId = callID
	resp.Token = call.Token

	for _, mediaSource := range req.Source {
		ms := &model.RtcMediaSource{
			ChannelID:     req.ChannelId,
			UserID:        req.UserId,
			Type:          uint8(mediaSource.Type),
			Codec:         mediaSource.Codec,
			MediaSpecific: mediaSource.MediaSpecific,
			Status:        0,
		}
		ssrc, err := s.dao.CreateMediaSource(ctx, ms)
		if err != nil {
			return nil, err
		}
		mediaSource.Ssrc = ssrc
		mediaSource.UserId = req.UserId
		resp.Source = append(resp.Source, mediaSource)
	}
	return resp, nil
}

// LeaveChannel implementation
func (s *RtcService) LeaveChannel(ctx context.Context, req *v1pb.LeaveChannelRequest) (resp *v1pb.LeaveChannelResponse, err error) {
	resp = &v1pb.LeaveChannelResponse{}
	if e := s.dao.UpdateCallStatus(context.Background(), req.ChannelId, req.CallId, req.UserId, time.Now(), 1); e != nil {
		log.Error("[LeaveChannel]UpdateCallStatus ChannelID:%d,CallID:%d,UserID:%d,error:%v", req.ChannelId, req.CallId, req.UserId, e)
	}
	if e := s.dao.UpdateMediaSourceStatus(context.Background(), req.ChannelId, req.CallId, req.UserId, 1); e != nil {
		log.Error("[LeaveChannel]UpdateMediaSourceStatus ChannelID:%d,CallID:%d,UserID:%d,error:%v", req.ChannelId, req.CallId, req.UserId, e)
	}
	return resp, nil
}

// PublishStream implementation
func (s *RtcService) PublishStream(ctx context.Context, req *v1pb.PublishStreamRequest) (resp *v1pb.PublishStreamResponse, err error) {
	resp = &v1pb.PublishStreamResponse{}
	if e := s.dao.CreateMediaPublish(ctx, &model.RtcMediaPublish{
		UserID:       req.UserId,
		CallID:       req.CallId,
		ChannelID:    req.ChannelId,
		Switch:       1,
		Width:        req.EncoderConfig.Width,
		Height:       req.EncoderConfig.Height,
		FrameRate:    uint8(req.EncoderConfig.FrameRate),
		VideoCodec:   req.EncoderConfig.VideoCodec,
		VideoProfile: req.EncoderConfig.VideoProfile,
		Channel:      uint8(req.EncoderConfig.Channel),
		SampleRate:   req.EncoderConfig.SampleRate,
		AudioCodec:   req.EncoderConfig.AudioCodec,
		Bitrate:      req.EncoderConfig.Bitrate,
		MixConfig:    req.MixConfig,
	}); e != nil {
		return nil, e
	}
	return resp, nil
}

// TerminateStream implementation
func (s *RtcService) TerminateStream(ctx context.Context, req *v1pb.TerminateStreamRequest) (resp *v1pb.TerminateStreamResponse, err error) {
	resp = &v1pb.TerminateStreamResponse{}
	if e := s.dao.TerminateStream(ctx, req.ChannelId, req.CallId); e != nil {
		return nil, e
	}
	return resp, nil
}

// Channel implementation
func (s *RtcService) Channel(ctx context.Context, req *v1pb.ChannelRequest) (resp *v1pb.ChannelResponse, err error) {
	resp = &v1pb.ChannelResponse{}
	mediaSource, err := s.dao.GetMediaSource(ctx, req.ChannelId)
	if err != nil {
		return nil, err
	}
	for _, s := range mediaSource {
		var mediaType v1pb.MediaSource_MediaType
		switch s.Type {
		case 1:
			mediaType = v1pb.MediaSource_VIDEO
		case 2:
			mediaType = v1pb.MediaSource_AUDIO
		case 3:
			mediaType = v1pb.MediaSource_DATA
		case 4:
			mediaType = v1pb.MediaSource_SMALL_VIDEO
		default:
			mediaType = v1pb.MediaSource_OTHER
		}
		resp.MediaSource = append(resp.MediaSource, &v1pb.MediaSource{
			Type:          mediaType,
			Codec:         s.Codec,
			MediaSpecific: s.MediaSpecific,
			Ssrc:          s.SourceID,
			UserId:        s.UserID,
		})
	}
	resp.Server, err = s.dispatch.AccessNode(req.ChannelId)
	if err != nil {
		return nil, err
	}
	//TODO: Read this value from Config
	resp.TcpPort = 2247
	resp.UdpPort = 2248
	return resp, nil
}

// Stream implementation
func (s *RtcService) Stream(ctx context.Context, req *v1pb.StreamRequest) (resp *v1pb.StreamResponse, err error) {
	resp = &v1pb.StreamResponse{}
	publish, err := s.dao.GetMediaPublishConfig(ctx, req.ChannelId, req.CallId)
	if err != nil {
		return nil, err
	}
	resp.MixConfig = publish.MixConfig
	resp.EncoderConfig = &v1pb.EncoderConfig{
		Width:        publish.Width,
		Height:       publish.Height,
		Bitrate:      publish.Bitrate,
		FrameRate:    uint32(publish.FrameRate),
		VideoCodec:   publish.VideoCodec,
		VideoProfile: publish.VideoProfile,
		Channel:      uint32(publish.Channel),
		SampleRate:   publish.SampleRate,
		AudioCodec:   publish.AudioCodec,
	}
	return resp, nil
}

// SetRtcConfig implementation
// `method:"POST"`
func (s *RtcService) SetRtcConfig(ctx context.Context, req *v1pb.SetRtcConfigRequest) (resp *v1pb.SetRtcConfigResponse, err error) {
	resp = &v1pb.SetRtcConfigResponse{}
	if e := s.dao.UpdateMediaPublishConfig(ctx, req.ChannelId, req.CallId, req.Config); e != nil {
		return nil, e
	}
	return resp, nil
}

// VerifyToken implementation
// `method:"GET"`
func (s *RtcService) VerifyToken(ctx context.Context, req *v1pb.VerifyTokenRequest) (resp *v1pb.VerifyTokenResponse, err error) {
	resp = &v1pb.VerifyTokenResponse{}
	if strings.Compare(req.Token, "") == 0 {
		resp.Pass = false
		return resp, nil
	}
	token, e := s.dao.GetToken(ctx, req.ChannelId, req.CallId)
	if e != nil {
		return nil, e
	}
	if strings.Compare(token, req.Token) == 0 {
		resp.Pass = true
	} else {
		resp.Pass = false
	}
	return resp, nil
}
