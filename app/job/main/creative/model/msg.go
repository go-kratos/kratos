package model

import "encoding/json"

//Msg for databus.
type Msg struct {
	MID       int64 `json:"mid"`
	From      int   `json:"from"`
	IsAuthor  int   `json:"is_author"`
	TimeStamp int64 `json:"timestamp"`
}

// CanalMsg canal databus msg.
type CanalMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

//TaskMsg for task notify.
type TaskMsg struct {
	MID       int64 `json:"mid"`
	Count     int64 `json:"count"`
	From      int   `json:"from"`
	TimeStamp int64 `json:"timestamp"`
}

// ShareMsg share databus msg.
type ShareMsg struct {
	OID  int64 `json:"oid"`
	MID  int64 `json:"mid"`
	TP   int   `json:"tp"`
	Time int64 `json:"time"`
}

// StatLike archive like count
type StatLike struct {
	MID          int64  `json:"mid"`
	Type         string `json:"type"`
	ID           int64  `json:"id"`
	Count        int64  `json:"count"`
	DislikeCount int64  `json:"dislike_count"`
	TimeStamp    int64  `json:"timestamp"`
}

// StatView ViewMsg archive view count
type StatView struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// StatDM archive DM count
type StatDM struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// StatReply archive reply count
type StatReply struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// StatFav archive collection count
type StatFav struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// StatCoin archive coin count
type StatCoin struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// StatShare archive share count
type StatShare struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// StatRank archive rank
type StatRank struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// RelaMessage Message define relation binlog databus message.
type RelaMessage struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Relation user_relation_mid_0~user_relation_mid_49
type Relation struct {
	MID       int64  `json:"mid,omitempty"`
	FID       int64  `json:"fid,omitempty"`
	Attribute uint32 `json:"attribute"`
	Status    int    `json:"status"`
	MTime     string `json:"mtime"`
	CTime     string `json:"ctime"`
}

// Stat user_relation_stat
type Stat struct {
	MID       int64 `json:"mid,omitempty"`
	Following int64 `json:"following"`
	Whisper   int64 `json:"whisper"`
	Black     int64 `json:"black"`
	Follower  int64 `json:"follower"`
}
