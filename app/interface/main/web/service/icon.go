package service

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/main/web/model"
	resmdl "go-common/app/service/main/resource/model"
	"go-common/library/log"
)

const (
	_iconFixType = "fix"
)

// IndexIcon get index icons
func (s *Service) IndexIcon() (res *model.IndexIcon) {
	return s.indexIcon
}

func fmtIndexIcon(icons []*resmdl.IndexIcon) {
	for _, v := range icons {
		v.Icon = strings.Replace(v.Icon, "http://", "//", 1)
	}
}

func (s *Service) randomIndexIcon(icons []*resmdl.IndexIcon) (icon *model.IndexIcon) {
	var (
		item          *resmdl.IndexIcon
		total, weight int
	)
	length := len(icons)
	if length == 0 {
		return new(model.IndexIcon)
	}
	for _, v := range icons {
		if v.Weight == 0 {
			total++
		} else {
			total += v.Weight
		}
	}
	if total == length {
		item = icons[s.r.Intn(length)]
		return &model.IndexIcon{ID: item.ID, Title: item.Title, Links: item.Links, Icon: item.Icon, Weight: item.Weight}
	}
	randWeight := s.r.Intn(total)
	for _, v := range icons {
		if v.Weight == 0 {
			weight++
		} else {
			weight += v.Weight
		}
		if weight > randWeight {
			item = v
			break
		}
	}
	return &model.IndexIcon{ID: item.ID, Title: item.Title, Links: item.Links, Icon: item.Icon, Weight: item.Weight}
}

func (s *Service) indexIconproc() {
	var (
		data  map[string][]*resmdl.IndexIcon
		icons []*resmdl.IndexIcon
		ok    bool
		err   error
		mutex = sync.RWMutex{}
	)
	go func() {
		for {
			if data, err = s.res.IndexIcon(context.Background()); err != nil {
				log.Error("s.res.IndexIcon error(%v)", err)
				time.Sleep(time.Second)
				continue
			}
			mutex.Lock()
			icons, ok = data[_iconFixType]
			mutex.Unlock()
			if ok {
				if len(icons) == 0 {
					log.Error("s.res.IndexIcon data error")
					time.Sleep(time.Second)
					continue
				} else {
					sort.Slice(icons, func(i, j int) bool { return icons[i].Weight > icons[j].Weight })
					fmtIndexIcon(icons)
				}
			} else {
				log.Error("s.res.IndexIcon data error")
				time.Sleep(time.Second)
				continue
			}
			time.Sleep(time.Duration(s.c.WEB.PullOnlineInterval))
		}
	}()
	for {
		mutex.RLock()
		tempIcons := make([]*resmdl.IndexIcon, len(icons))
		copy(tempIcons, icons)
		mutex.RUnlock()
		s.indexIcon = s.randomIndexIcon(tempIcons)
		time.Sleep(500 * time.Millisecond)
	}
}
