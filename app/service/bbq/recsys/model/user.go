package model

import (
	"github.com/Dai0522/go-hash/bloomfilter"
)

//Tuple ...
type Tuple struct {
	Timestamp int64
	Count     int64
}

//UserProfile 用户画像数据 包括历史画像和实时日志
type UserProfile struct {
	Mid    int64  `json:"Mid,omitempty"`
	Buvid  string `json:"Buvid,omitempty"`
	Name   string `json:"Name,omitempty"`
	Gender int8   `json:"Gender,omitempty"`

	ViewVideos []int64 `json:"ViewVideos,omitempty"`

	//bbq user profile
	//key:up mid, value: timestamp
	BBQFollowAction map[int64]int64 `json:"BBQFollowAction,omitempty"`
	//key:up mid, value: 1
	BBQFollow map[int64]int64 `json:"BBQFollow,omitempty"`
	BBQBlack  map[int64]int64 `json:"BBQBlack,omitempty"`

	BBQTags    map[string]float64 `json:"BBQTags,omitempty"`
	BBQZones   map[string]float64 `json:"BBQZones,omitempty"`
	BBQPrefUps map[int64]int64    `json:"BBQPrefUps,omitempty"`

	//bili user profile
	BiliTags  map[string]float64 `json:"BiliTags,omitempty"`
	Zones1    map[string]float64 `json:"Zones1,omitempty"`
	Zones2    map[string]float64 `json:"Zones2,omitempty"`
	FollowUps map[int64]int64    `json:"FollowUps,omitempty"`

	//bbq实时数据
	//key: SVID, value: timestamp
	PosVideos  map[int64]int64 `json:"PosVideos,omitempty"`
	NegVideos  map[int64]int64 `json:"NegVideos,omitempty"`
	LikeVideos map[int64]int64 `json:"LikeVideos,omitempty"`

	//key: tagID, value: count
	LikeTagIDs map[int64]int64 `json:"LikeTagIDs,omitempty"`
	PosTagIDs  map[int64]int64 `json:"PosTagIDs,omitempty"`
	NegTagIDs  map[int64]int64 `json:"NegTagIDs,omitempty"`

	//key: UP MID, value: timestamp
	LikeUPs map[int64]int64 `json:"LikeUPs,omitempty"`

	//for old retrieve function
	LikeTags map[string]float64 `json:"LikeTags,omitempty"`
	PosTags  map[string]float64 `json:"PosTags,omitempty"`
	NegTags  map[string]float64 `json:"NegTags,omitempty"`

	//DedupVideos 根据ID去重
	DedupVideos    []int64      `json:"DedupVideos,omitempty"`
	LastRecords    []Record4Dup `json:"LastRecords,omitempty"`
	LastUpsRecords []Record4Dup `json:"LastRecords,omitempty"`

	//BloomFilter 去重用到 SVID
	BloomFilter *bloomfilter.BloomFilter `json:"BloomFilter,omitempty"`
}
