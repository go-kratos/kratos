package model

import "time"

const (
	_ int8 = iota
	// TargetType001 该UID下开放浏览的稿件≥1
	TargetType001
	// TargetType002 该UID分享自己视频的次数≥1
	TargetType002
	// TargetType003 该UID在创作学院的观看记录≥1
	TargetType003
	// TargetType004 该UID下所有avid的获得评论数≥3
	TargetType004
	// TargetType005 该UID下所有avid获得分享数≥3
	TargetType005
	// TargetType006 该UID的所有avid的获得收藏数≥5
	TargetType006
	// TargetType007 该UID下所有avid的获得硬币数≥5
	TargetType007
	// TargetType008 该UID下所有avid获得点赞数≥5
	TargetType008
	// TargetType009 该UID下所有avid的获得弹幕数≥5
	TargetType009
	// TargetType010 该UID的粉丝数≥10
	TargetType010
	// TargetType011 任务完成期间该UID的水印开关为打开状态
	TargetType011
	// TargetType012 该UID的关注列表含有“哔哩哔哩创作中心”
	TargetType012
	// TargetType013 用手机投稿上传视频
	TargetType013
	// TargetType014 该UID下开放浏览的稿件≥5
	TargetType014
	// TargetType015 该UID下任意avid的获得点击量≥1000
	TargetType015
	// TargetType016 该UID下任意avid的评论≥30
	TargetType016
	// TargetType017 该UID下任意avid的获得分享数≥10
	TargetType017
	// TargetType018 该UID下任意avid的获得收藏数≥30
	TargetType018
	// TargetType019 该UID下任意avid的获得硬币数≥50
	TargetType019
	// TargetType020 该UID下任意avid的获得点赞数≥50
	TargetType020
	// TargetType021 该UID下任意avid的获得弹幕数≥50
	TargetType021
	// TargetType022 该UID的粉丝数≥1000
	TargetType022
	// TargetType023 该UID的激励计划状态为已开通
	TargetType023
	// TargetType024 该UID粉丝勋章为开启状态
	TargetType024
)

const (
	//TaskIncomplete  任务未完成
	TaskIncomplete = -1
	//TaskCompleted   任务已完成
	TaskCompleted = 0

	//MsgForWaterMark 发送用户设置水印消息
	MsgForWaterMark = 1
	//MsgForAcademyFavVideo 发送用户已在创作学院观看过自己喜欢的视频的消息
	MsgForAcademyFavVideo = 2
	//MsgForGrowAccount 发送用户已在参加激励计划的消息
	MsgForGrowAccount = 3
	//MsgForOpenFansMedal 成功开通粉丝勋章
	MsgForOpenFansMedal = 4
)

// UserTask for def user task struct.
type UserTask struct {
	ID           int64     `json:"id"`
	MID          int64     `json:"mid"`
	TaskID       int64     `json:"task_id"`
	TaskGroupID  int64     `json:"task_group_id"`
	TaskType     int8      `json:"task_type"`
	State        int8      `json:"state"`
	TaskBindTime time.Time `json:"task_bind_time"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"mtime"`
}

// Task for def task struct.
type Task struct {
	ID           int64     `json:"id"`
	GroupID      int64     `json:"-"`
	Type         int8      `json:"type"`
	State        int8      `json:"-"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
	Comment      string    `json:"-"`
	TargetType   int8      `json:"-"`
	TargetValue  int64     `json:"-"`
	CompleteSate int8      `json:"complete_state"`
	CTime        time.Time `json:"-"`
	MTime        time.Time `json:"-"`
}

// GiftReward for gift reward
type GiftReward struct {
	ID       int64     `json:"id"`
	TaskType int8      `json:"task_type"`
	RewardID int64     `json:"reward_id"`
	State    int8      `json:"state"`
	Comment  string    `json:"comment"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// Up for up new arc.
type Up struct {
	AID int64
	MID int64
}
