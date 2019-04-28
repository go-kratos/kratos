package model

import (
	xtime "go-common/library/time"
)

// type and states
const (
	StateBlank   = 0
	StateLike    = 1
	StateDislike = 2

	TypeLike           = 1
	TypeCancelLike     = 2
	TypeDislike        = 3
	TypeCancelDislike  = 4
	TypeLikeReverse    = 5
	TypeDislikeReverse = 6

	ItemListLike    = 1
	ItemListDislike = 2
	ItemListAll     = 3

	UserListLike    = 1
	UserListDislike = 2
	UserListAll     = 3
)

// Business .
type Business struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	MessageListType   uint8  `json:"message_list_type"`
	UserListType      uint8  `json:"user_list_type"`
	UserLikesLimit    int    `json:"user_likes_limit"`
	MessageLikesLimit int    `json:"message_likes_limit"`
	EnableOriginID    int    `json:"enable_origin_id"`
}

// EnableItemLikeList .
func (b *Business) EnableItemLikeList() bool {
	return (b.MessageListType == ItemListLike) || (b.MessageListType == ItemListAll)
}

// EnableItemDislikeList .
func (b *Business) EnableItemDislikeList() bool {
	return (b.MessageListType == ItemListDislike) || (b.MessageListType == ItemListAll)
}

// EnableUserLikeList .
func (b *Business) EnableUserLikeList() bool {
	return (b.UserListType == UserListLike) || (b.UserListType == UserListAll)
}

// EnableUserDislikeList .
func (b *Business) EnableUserDislikeList() bool {
	return (b.UserListType == UserListDislike) || (b.UserListType == UserListAll)
}

// UserLikeRecord .
type UserLikeRecord struct {
	Mid  int64      `json:"mid"`
	Time xtime.Time `json:"time"`
}

// LikeItem .
type LikeItem struct {
	Business  string
	OriginID  int64
	MessageID int64
}

// LikeCounts .
type LikeCounts struct {
	Like    int64
	Dislike int64
	UpMid   int64
}

// Stats .
type Stats struct {
	OriginID int64 `json:"origin_id"`
	ID       int64 `json:"id"`
	Likes    int64 `json:"likes"`
	Dislikes int64 `json:"dislikes"`
}

// ItemLikeRecord .
type ItemLikeRecord struct {
	MessageID int64      `json:"message_id"`
	Time      xtime.Time `json:"time"`
}

// StatMsg .
type StatMsg struct {
	Type         string `json:"type"`
	ID           int64  `json:"id"`
	Count        int64  `json:"count"`
	Timestamp    int64  `json:"timestamp"`
	OriginID     int64  `json:"origin_id,omitempty"`
	DislikeCount int64  `json:"dislike_count,omitempty"`
	Mid          int64  `json:"mid,omitempty"`
	UpMid        int64  `json:"up_mid,omitempty"`
}
