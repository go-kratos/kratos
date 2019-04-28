package anchorReward

import (
	"context"
	model "go-common/app/service/live/xrewardcenter/model/anchorTask"
	"go-common/library/log"
)

// RawRewardConf return reward config from db.
func (d *Dao) RawRewardConf(c context.Context, id int64) (res *model.AnchorRewardConf, err error) {
	rewards := []*model.AnchorRewardConf{}

	if err := d.orm.Model(&model.AnchorReward{}).Find(&rewards, "id=?", id).Error; err != nil {
		log.Error("getRewardById (%v) error(%v)", id, err)
		return nil, err
	}
	if len(rewards) != 0 {
		res = rewards[0]
	}

	return
}
