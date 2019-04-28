package http

import (
	"strconv"

	"go-common/app/interface/main/activity/model/like"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func subject(c *bm.Context) {
	params := c.Request.Form
	sidStr := params.Get("sid")
	sid, err := strconv.ParseInt(sidStr, 10, 32)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", sidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(likeSvc.Subject(c, sid))
}

func vote(c *bm.Context) {
	var (
		mid int64
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	params := c.Request.Form
	voteStr := params.Get("vote")
	vote, err := strconv.ParseInt(voteStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	stageStr := params.Get("stage")
	stage, err := strconv.ParseInt(stageStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if strRe, _ := likeSvc.OnlineVote(c, mid, vote, stage, aid); !strRe {
		c.JSON(nil, ecode.NotModified)
		return
	}
	c.JSON("ok", nil)
}

func ltime(c *bm.Context) {
	params := c.Request.Form
	sidStr := params.Get("sid")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	index, err := likeSvc.Ltime(c, sid)
	if err != nil {
		log.Error("error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if index == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(index, nil)
}

func likeAct(c *bm.Context) {
	p := new(like.ParamAddLikeAct)
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.LikeAct(c, p, mid))
}

func storyKingAct(c *bm.Context) {
	p := new(like.ParamStoryKingAct)
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.StoryKingAct(c, p, mid))
}

func storyKingLeft(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.StoryKingLeftTime(c, p.Sid, mid))

}

func upList(c *bm.Context) {
	p := new(like.ParamList)
	if err := c.Bind(p); err != nil {
		return
	}
	mid := int64(0)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	c.JSON(likeSvc.UpList(c, p, mid))
}

func missionLike(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.MissionLike(c, p.Sid, mid))
}

func missionLikeAct(c *bm.Context) {
	p := new(like.ParamMissionLikeAct)
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.MissionLikeAct(c, p, mid))
}

func missionInfo(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
		Lid int64 `form:"lid" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.MissionInfo(c, p.Sid, p.Lid, mid))

}
func missionTops(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
		Num int   `form:"num" validate:"min=1,max=200"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(likeSvc.MissionTops(c, p.Sid, p.Num))
}

func missionUser(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
		Lid int64 `form:"lid" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(likeSvc.MissionUser(c, p.Sid, p.Lid))
}

func missionRank(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.MissionRank(c, p.Sid, mid))
}

func missionFriends(c *bm.Context) {
	p := new(like.ParamMissionFriends)
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.MissionFriendsList(c, p, mid))
}

func missionAward(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.MissionAward(c, p.Sid, mid))
}

func missionAchieve(c *bm.Context) {
	p := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
		ID  int64 `form:"id" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(likeSvc.MissionAchieve(c, p.Sid, p.ID, mid))
}

func likeActList(c *bm.Context) {
	v := new(struct {
		Sid  int64   `form:"sid" validate:"min=1"`
		Lids []int64 `form:"lids,split" validate:"min=1,max=50,dive,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	mid := int64(0)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	c.JSON(likeSvc.LikeActList(c, v.Sid, mid, v.Lids))
}

func subjectInit(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.SubjectInitialize(c, v.Sid-1))
}

func likeInit(c *bm.Context) {
	v := new(struct {
		Lid int64 `form:"lid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.LikeInitialize(c, v.Lid-1))
}

func subjectLikeListInit(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.SubjectLikeListInitialize(c, v.Sid))
}

func likeActCountInit(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.LikeActCountInitialize(c, v.Sid))
}

func tagList(c *bm.Context) {
	var (
		err  error
		cnt  int
		list []*like.Like
	)
	v := new(struct {
		Sid   int64  `form:"sid" validate:"min=1"`
		TagID int64  `form:"tag_id" validate:"min=1"`
		Type  string `form:"type" default:"ctime"`
		Pn    int    `form:"pn" default:"1" validate:"min=1"`
		Ps    int    `form:"ps" default:"30" validate:"min=1,max=30"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Type != "ctime" && v.Type != "random" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if list, cnt, err = likeSvc.TagArcList(c, v.Sid, v.TagID, v.Pn, v.Ps, v.Type, metadata.String(c, metadata.RemoteIP)); err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"page": map[string]int{
			"num":   v.Pn,
			"size":  v.Ps,
			"total": cnt,
		},
		"list": list,
	}
	c.JSON(data, nil)
}

func regionList(c *bm.Context) {
	var (
		err  error
		cnt  int
		list []*like.Like
	)
	v := new(struct {
		Sid  int64  `form:"sid" validate:"min=1"`
		Rid  int16  `form:"rid" validate:"min=1"`
		Type string `form:"type" default:"ctime"`
		Pn   int    `form:"pn" default:"1" validate:"min=1"`
		Ps   int    `form:"ps" default:"30" validate:"min=1,max=30"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Type != "ctime" && v.Type != "random" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if list, cnt, err = likeSvc.RegionArcList(c, v.Sid, v.Rid, v.Pn, v.Ps, v.Type, metadata.String(c, metadata.RemoteIP)); err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"page": map[string]int{
			"num":   v.Pn,
			"size":  v.Ps,
			"total": cnt,
		},
		"list": list,
	}
	c.JSON(data, nil)
}

func tagStats(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(likeSvc.TagLikeCounts(c, v.Sid))
}

func subjectStat(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(likeSvc.SubjectStat(c, v.Sid))
}

func setSubjectStat(c *bm.Context) {
	v := new(like.SubjectStat)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.SetSubjectStat(c, v))
}

func viewRank(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
		Pn  int   `form:"pn" default:"1" validate:"min=1"`
		Ps  int   `form:"ps" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	list, count, err := likeSvc.ViewRank(c, v.Sid, v.Pn, v.Ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["list"] = list
	data["page"] = map[string]int{
		"pn":    v.Pn,
		"ps":    v.Ps,
		"count": count,
	}
	c.JSON(data, err)
}

func setViewRank(c *bm.Context) {
	v := new(struct {
		Sid  int64   `form:"sid" validate:"min=1"`
		Aids []int64 `form:"aids,split" validate:"min=1,dive,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.SetViewRank(c, v.Sid, v.Aids))
}

func groupData(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	ck := c.Request.Header.Get("cookie")
	c.JSON(likeSvc.ObjectGroup(c, v.Sid, ck))
}

func setLikeContent(c *bm.Context) {
	v := new(struct {
		Lid int64 `form:"lid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.SetLikeContent(c, v.Lid))
}

func addLikeAct(c *bm.Context) {
	v := new(struct {
		Sid   int64 `form:"sid" validate:"min=1"`
		Lid   int64 `form:"lid" validate:"min=1"`
		Score int64 `form:"score" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, likeSvc.AddLikeActCache(c, v.Sid, v.Lid, v.Score))
}

func likeActCache(c *bm.Context) {
	v := new(struct {
		Sid int64 `form:"sid" validate:"min=1"`
		Lid int64 `form:"lid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(likeSvc.LikeActCache(c, v.Sid, v.Lid))
}

func likeOidsInfo(c *bm.Context) {
	v := new(struct {
		Type int     `form:"type" validate:"min=1"`
		Oids []int64 `form:"oids,split" validate:"required,min=1,max=50,dive,gt=0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(likeSvc.LikeOidsInfo(c, v.Type, v.Oids))
}
