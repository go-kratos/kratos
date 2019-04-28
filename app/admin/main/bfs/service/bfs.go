package service

import (
	"context"

	"go-common/app/admin/main/bfs/model"
	"go-common/library/ecode"
)

// Total .
func (s *Service) Total(c context.Context, arg *model.ArgCluster) (resp *model.RespTotal, err error) {
	volumes, err := s.Volumes(c, arg)
	if err != nil {
		return
	}
	num := int64(len(volumes.Volumes))
	var fs int64
	for _, volume := range volumes.Volumes {
		fs += int64(volume.FreeSpace)
	}
	fs = (fs * 8) / 1024 / 1024 / 1024 // GB
	// groups
	groups, err := s.d.Groups(arg.Cluster)
	if err != nil {
		return
	}
	// stores
	racks, err := s.d.Racks(arg.Cluster)
	if err != nil {
		return
	}
	var stores int64
	for _, rack := range racks {
		stores += int64(len(rack.Stores))
	}
	resp = &model.RespTotal{
		Space:     32 * num,
		FreeSpace: fs,
		Groups:    int64(len(groups)),
		Stores:    stores,
		Volumes:   num,
	}
	return
}

// Racks .
func (s *Service) Racks(c context.Context, arg *model.ArgCluster) (resp *model.RespRack, err error) {
	racks, err := s.d.Racks(arg.Cluster)
	if err != nil {
		return
	}
	for _, rack := range racks {
		for _, store := range rack.Stores {
			store.ParseStates()
		}
	}
	resp = &model.RespRack{Racks: racks}
	return
}

// Groups .
func (s *Service) Groups(c context.Context, arg *model.ArgCluster) (resp *model.RespGroup, err error) {
	groups, err := s.d.Groups(arg.Cluster)
	if err != nil {
		return
	}
	racks, err := s.Racks(c, arg)
	if err != nil {
		return
	}
	volumes, err := s.Volumes(c, arg)
	if err != nil {
		return
	}
	for _, group := range groups {
		group.StoreDatas = make(map[string]*model.Store)
		for _, rack := range racks.Racks {
			for sname, store := range rack.Stores {
				for _, name := range group.Stores {
					if sname == name {
						group.StoreDatas[name] = store
					}
				}
			}
		}
		for _, store := range group.StoreDatas {
			num := int64(len(store.Volumes))
			var fs int64
			for _, volume := range store.Volumes {
				fs += int64(volumes.Volumes[volume].FreeSpace)
			}
			fs = (fs * 8) / 1024 / 1024 / 1024 // GB
			group.Total.Space = 32 * num
			group.Total.FreeSpace = fs
			group.Total.Volumes = num
			break
		}
	}
	resp = &model.RespGroup{Groups: groups}
	return
}

// Volumes .
func (s *Service) Volumes(c context.Context, arg *model.ArgCluster) (resp *model.RespVolume, err error) {
	volumes, err := s.d.Volumes(arg.Cluster)
	if err != nil {
		return
	}
	resp = &model.RespVolume{Volumes: volumes}
	return
}

// AddVolume add volume.
func (s *Service) AddVolume(c context.Context, arg *model.ArgAddVolume) (err error) {
	return s.d.AddVolume(c, arg.Group, arg.Num)
}

// AddFreeVolume add free volume.
func (s *Service) AddFreeVolume(c context.Context, arg *model.ArgAddFreeVolume) (err error) {
	return s.d.AddFreeVolume(c, arg.Group, arg.Dir, arg.Num)
}

// Compact compact store.
func (s *Service) Compact(c context.Context, arg *model.ArgCompact) (err error) {
	return s.d.Compact(c, arg.Group, arg.Vid)
}

// SetGroupStatus set group status(read,write,sync,health).
func (s *Service) SetGroupStatus(c context.Context, arg *model.ArgGroupStatus) (err error) {
	if arg.Status != "read" && arg.Status != "write" && arg.Status != "sync" && arg.Status != "health" {
		err = ecode.RequestErr
		return
	}
	return s.d.SetGroupStatus(c, arg.Group, arg.Status)
}
