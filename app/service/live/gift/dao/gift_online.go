package dao

import (
	"context"
	"go-common/app/service/live/gift/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_getAllGift = "SELECT id,gift_id,name,price,coin_type,type,effect,corner_mark,broadcast,draw,asset_img_basic,asset_img_dynamic,asset_frame_animation,animation_frame_num,asset_gif,asset_webp,asset_full_sc_web,asset_full_sc_horizontal,asset_full_sc_vertical,asset_full_sc_horizontal_svga,asset_full_sc_vertical_svga,asset_bullet_head,asset_bullet_tail,`desc`,rights,rule,limit_interval FROM gift_online"
)

// GetAllGift GetAllGift
func (d *Dao) GetAllGift(ctx context.Context) (gifts []*model.GiftOnline, err error) {
	log.Info("GetAllGift")
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getAllGift); err != nil {
		log.Error("query getAllGift error,err %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		g := &model.GiftOnline{}
		if err = rows.Scan(&g.Id, &g.GiftId, &g.Name, &g.Price, &g.CoinType, &g.Type, &g.Effect, &g.CornerMark, &g.Broadcast,
			&g.Draw, &g.AssetImgBasic, &g.AssetImgDynamic, &g.AssetFrameAnimation, &g.AnimationFrameNum, &g.AssetGif, &g.AssetWebp,
			&g.AssetFullScWeb, &g.AssetFullScHorizontal, &g.AssetFullScVertical, &g.AssetFullScHorizontalSvga, &g.AssetFullScVerticalSvga,
			&g.AssetBulletHead, &g.AssetBulletTail, &g.Desc, &g.Rights, &g.Rule, &g.LimitInterval); err != nil {
			log.Error("getAllGift scan error,err %v", err)
			return
		}
		gifts = append(gifts, g)
	}
	return
}
