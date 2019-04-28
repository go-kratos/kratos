package model

// DiscountInfo DiscountInfo
type DiscountInfo struct {
	Id         int64                            `json:"id"`
	SceneKey   int64                            `json:"scene_key"`
	SceneValue []int64                          `json:"scene_value"`
	Platform   int64                            `json:"platform"`
	List       map[int64]map[string]interface{} `json:"list"`
}

// BagGiftStatus BagGiftStatus
type BagGiftStatus struct {
	Status int64       `json:"status"`
	Gift   []*GiftInfo `json:"gift"`
}

// GiftInfo GiftInfo
type GiftInfo struct {
	GiftID   int64  `json:"gift_id"`
	GiftNum  int64  `json:"gift_num"`
	ExpireAt string `json:"expireat"`
}

//DayGiftInfo DayGiftInfo
type DayGiftInfo struct {
	ID      int64  `json:"id"`
	UID     int64  `json:"uid"`
	Day     string `json:"day"`
	DayInfo string `json:"day_info"`
}

//BagInfo BagInfo
type BagInfo struct {
	ID      int64 `json:"id"`
	GiftNum int64 `json:"gift_num"`
}

//WeekGiftInfo WeekGiftInfo
type WeekGiftInfo struct {
	ID       int64  `json:"id"`
	UID      int64  `json:"uid"`
	Week     int    `json:"week"`
	Level    int64  `json:"level"`
	WeekInfo string `json:"week_info"`
}

//BagGiftList BagGiftList
type BagGiftList struct {
	ID       int64 `json:"id"`
	UID      int64 `json:"uid"`
	GiftID   int64 `json:"gift_id"`
	GiftNum  int64 `json:"gift_num"`
	ExpireAt int64 `json:"expireat"`
}

// GiftOnline .
type GiftOnline struct {
	Id                        int64  `json:"id"`
	GiftId                    int64  `json:"gift_id"`
	Name                      string `json:"name"`
	Price                     int64  `json:"price"`
	CoinType                  int64  `json:"coin_type"`
	Type                      int64  `json:"type"`
	Effect                    int64  `json:"effect"`
	CornerMark                string `json:"corner_mark"`
	Broadcast                 int64  `json:"broadcast"`
	Draw                      int64  `json:"draw"`
	AssetImgBasic             string `json:"asset_img_basic"`
	AssetImgDynamic           string `json:"asset_img_dynamic"`
	AssetFrameAnimation       string `json:"asset_frame_animation"`
	AnimationFrameNum         int64  `json:"animation_frame_num"`
	AssetGif                  string `json:"asset_gif"`
	AssetWebp                 string `json:"asset_webp"`
	AssetFullScWeb            string `json:"asset_full_sc_web"`
	AssetFullScHorizontal     string `json:"asset_full_sc_horizontal"`
	AssetFullScVertical       string `json:"asset_full_sc_vertical"`
	AssetFullScHorizontalSvga string `json:"asset_full_sc_horizontal_svga"`
	AssetFullScVerticalSvga   string `json:"asset_full_sc_vertical_svga"`
	AssetBulletHead           string `json:"asset_bullet_head"`
	AssetBulletTail           string `json:"asset_bullet_tail"`
	Desc                      string `json:"desc"`
	Rights                    string `json:"rights"`
	Rule                      string `json:"rule"`
	PrivilegeRequired         int64  `json:"privilege_required"`
	LimitInterval             int64  `json:"limit_interval"`
}

//GiftPlan .
type GiftPlan struct {
	Id         int64  `json:"id"`
	List       string `json:"list"`
	SilverList string `json:"silver_list"`
	SceneKey   int64  `json:"scene_key"`
	SceneValue int64  `json:"scene_value"`
	Mtime      string `json:"mtime"`
	Platform   int64  `json:"platform"`
}

//DiscountPlan .
type DiscountPlan struct {
	Id         int64  `json:"id"`
	SceneKey   int64  `json:"scene_key"`
	SceneValue string `json:"scene_value"`
	Platform   int64  `json:"platform"`
}

//DiscountGift .
type DiscountGift struct {
	Id             int64  `json:"id"`
	DiscountId     int64  `json:"discount_id"`
	GiftId         int64  `json:"gift_id"`
	UserType       int64  `json:"user_type"`
	DiscountPrice  int64  `json:"discount_price"`
	CornerMark     string `json:"corner_mark"`
	CornerPosition int64  `json:"corner_position"`
}
