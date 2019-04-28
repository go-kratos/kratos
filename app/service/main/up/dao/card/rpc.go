package card

import (
	"context"
	"fmt"
	"sync"

	arcgrpc "go-common/app/service/main/archive/api"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const _avURL = "http://www.bilibili.com/video/av%d"

//ListVideoArchive list videos by avids
func (d *Dao) ListVideoArchive(ctx context.Context, avids []int64) (videos []*model.UpCardVideo, err error) {

	archives, err := d.listArchives(ctx, avids)
	if err != nil {
		log.Error("d.listArchives error(%v) arg(%v)", err, avids)
		err = ecode.CreativeArcServiceErr
		return
	}
	for avid, data := range archives {
		video := transfer(avid, data)
		videos = append(videos, video)
	}

	return
}

//AvidVideoMap get <avid, video> map by avids
func (d *Dao) AvidVideoMap(ctx context.Context, avids []int64) (avidVideoMap map[int64]*model.UpCardVideo, err error) {
	avidVideoMap = make(map[int64]*model.UpCardVideo)
	archives, err := d.listArchives(ctx, avids)
	if err != nil {
		log.Error("d.listArchives error(%v) arg(%v)", err, avids)
		err = ecode.CreativeArcServiceErr
		return
	}

	for avid, data := range archives {
		video := transfer(avid, data)
		avidVideoMap[avid] = video
	}
	return
}

func (d *Dao) listArchives(ctx context.Context, avids []int64) (archives map[int64]*arcgrpc.Arc, err error) {
	archives = make(map[int64]*arcgrpc.Arc)

	var (
		g errgroup.Group
		m sync.Mutex
	)

	dao.Split(0, len(avids), 300, func(start int, end int) {
		g.Go(func() (err error) {
			var (
				arg = &arcgrpc.ArcsRequest{
					Aids: avids[start:end],
				}
				res *arcgrpc.ArcsReply
			)
			if res, err = global.GetArcClient().Arcs(ctx, arg); err != nil {
				log.Error("d.acc.Archives3 arg(%v) error(%v)", arg, err)
				err = nil
			} else {
				for k, v := range res.Arcs {
					m.Lock()
					archives[k] = v
					m.Unlock()
				}
			}
			return
		})
	})

	if err = g.Wait(); err != nil {
		log.Error("g.Wait() error(%v)", err)
	}
	return
}

func transfer(avid int64, archive *arcgrpc.Arc) (video *model.UpCardVideo) {
	return &model.UpCardVideo{
		URL:      fmt.Sprintf(_avURL, avid),
		Title:    archive.Title,
		Picture:  archive.Pic,
		Duration: archive.Duration,
		CTime:    archive.Ctime,
	}
}
