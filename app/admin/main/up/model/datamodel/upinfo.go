package datamodel

//UpArchiveData hbase up archive info table
type UpArchiveData struct {
	Avs        int64 `json:"avs" qualifier:"avs"`          // 稿件数
	Play       int64 `json:"play" qualifier:"play"`        // 播放数
	Likes      int64 `json:"likes" qualifier:"likes"`      // 点赞数
	Danmu      int64 `json:"danmu" qualifier:"danmu"`      //弹幕
	Reply      int64 `json:"reply" qualifier:"reply"`      //评论
	Share      int64 `json:"share" qualifier:"share"`      //分享
	Fav        int64 `json:"fav" qualifier:"fav"`          //收藏
	Coin       int64 `json:"coin" qualifier:"coin"`        //硬币
	MaxPlayAid int64 `json:"max_play_aid" qualifier:"aid"` //最高播放的aid
	MaxPlay    int64 `json:"max_play" qualifier:"playest"` //最高播放次数
}

//UpArchiveTagData tag data
type UpArchiveTagData struct {
	// key-> index, value-> tagid
	TagMap map[string]int64 `family:"tag"`
}

//UpArchiveTypeData type data
type UpArchiveTypeData struct {
	Tid        int64  `json:"tid" qualifier:"tid"`           // 一级分区
	TidName    string `json:"tid_name"`                      // 分区名字
	Count      int64  `json:"count" qualifier:"cnt"`         // 稿件数量
	SubTid     int64  `json:"sub_tid" qualifier:"sub_tid"`   //二级分区
	SubCount   int64  `json:"sub_count" qualifier:"sub_cnt"` //稿件数量
	SubTidName string `json:"sub_tid_name"`                  // 分区名字
}
