package model

import "time"

type RtcMediaSource struct {
	SourceID      uint32
	ChannelID     uint64
	UserID        uint64
	Type          uint8
	Codec         string
	MediaSpecific string
	Status        uint8
}

type RtcChannel struct {
	ChannelID   uint64
	OwnerUserID uint64
	Type        uint8
	Status      uint8
	Cluster     string
}

type RtcCall struct {
	CallID    uint32
	UserID    uint64
	ChannelID uint64
	Version   uint32
	Token     string
	Status    uint8
	JoinTime  time.Time
	LeaveTime time.Time
}

type RtcMediaPublish struct {
	UserID       uint64
	CallID       uint32
	ChannelID    uint64
	Switch       uint8
	Width        uint32
	Height       uint32
	FrameRate    uint8
	VideoCodec   string
	VideoProfile string
	Channel      uint8
	SampleRate   uint32
	AudioCodec   string
	Bitrate      uint32
	MixConfig    string
}
