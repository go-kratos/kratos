package area

import (
	"context"
	"sync"

	"go-common/app/service/main/filter/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

type AreaList func(context.Context) ([]*model.Area, error)

func New() (a *Area) {
	return &Area{
		areaMap: make(map[string]*model.Area),
	}
}

type Area struct {
	areaMap map[string]*model.Area // map[areaname] area
	sync.RWMutex
}

func (a *Area) Load(ctx context.Context, loader AreaList) (err error) {
	if loader == nil {
		err = errors.New("filter service area load failed , loader is nil")
		return
	}
	var (
		areaMap  = make(map[string]*model.Area)
		areaList []*model.Area
	)
	if areaList, err = loader(ctx); err != nil {
		return
	}
	for _, area := range areaList {
		areaMap[area.Name] = area
	}
	a.Lock()
	a.areaMap = areaMap
	a.Unlock()
	return
}

func (a *Area) CheckArea(areas []string) bool {
	a.RLock()
	defer a.RUnlock()
	for _, area := range areas {
		if _, ok := a.areaMap[area]; !ok {
			log.Error("invalid area[%s]", area)
			return false
		}
	}
	return true
}

func (a *Area) AreaNames() (names []string) {
	a.RLock()
	defer a.RUnlock()
	for name := range a.areaMap {
		names = append(names, name)
	}
	return
}

func (a *Area) Area(name string) (area *model.Area) {
	a.RLock()
	defer a.RUnlock()
	area = a.areaMap[name]
	return
}
