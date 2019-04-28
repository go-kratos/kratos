package dao

import (
	"encoding/json"
	"fmt"
	"path"

	"go-common/app/admin/main/bfs/model"
	"go-common/library/log"
)

// Racks get rack infos.
func (d *Dao) Racks(cluster string) (racks map[string]*model.Rack, err error) {
	zkConn, ok := d.zkcs[cluster]
	if !ok {
		err = fmt.Errorf("not support cluster:%s", cluster)
		return
	}
	root := d.c.Zookeepers[cluster].RackRoot
	// get all rocks
	rackNames, _, err := zkConn.Children(root)
	if err != nil {
		log.Error("get from:%s children error(%v)", root, err)
		return
	}
	racks = make(map[string]*model.Rack)
	for _, rackName := range rackNames {
		// get all the store
		rackPath := path.Join(root, rackName)
		storeNames, _, err := zkConn.Children(rackPath)
		if err != nil {
			log.Error("cannot get /rack/%s info error(%v)", rackName, err)
			continue
		}
		log.Info("rack:%s got store:%v", rackName, storeNames)
		rack := &model.Rack{Stores: make(map[string]*model.Store)}
		racks[rackName] = rack
		// get the volume of store
		for _, storeName := range storeNames {
			storePath := path.Join(rackPath, storeName)
			storeMeta, _, err := zkConn.Get(storePath)
			if err != nil {
				log.Error("cannot get from:%s info error(%v)", storePath, err)
				continue
			}
			var data model.Store
			if err = json.Unmarshal(storeMeta, &data); err != nil {
				log.Error("cannot Unmarshal %s info get from:%s error(%v)", storeMeta, storePath, err)
				continue
			}
			volumes, _, err := zkConn.Children(storePath)
			if err != nil {
				log.Error("cannot list %s info error(%v)", storePath, err)
				continue
			}
			data.Volumes = volumes
			log.Info("rack:%s store:%s got volume:%v", rackName, storeName, volumes)
			rack.Stores[storeName] = &data
		}
	}
	return
}

// Volumes get volume infos.
func (d *Dao) Volumes(cluster string) (volumes map[string]*model.VolumeState, err error) {
	zkConn, ok := d.zkcs[cluster]
	if !ok {
		err = fmt.Errorf("not support cluster:%s", cluster)
		return
	}
	root := d.c.Zookeepers[cluster].VolumeRoot
	// get all volumes
	volumeIDs, _, err := zkConn.Children(root)
	if err != nil {
		log.Error("cannot get from:%s info error(%v)", root, err)
		return
	}
	log.Info("got volumes: %s", volumeIDs)
	volumes = make(map[string]*model.VolumeState)
	for _, id := range volumeIDs {
		// get all the volumes
		volumePath := path.Join(root, id)
		volumeState, _, err := zkConn.Get(volumePath)
		if err != nil {
			log.Error("cannot get %s info error(%v)", volumePath, err)
			continue
		}
		var data model.VolumeState
		if err := json.Unmarshal(volumeState, &data); err != nil {
			log.Error("cannot Unmarshal %s info get from:%s error(%v)", volumeState, volumePath, err)
			continue
		}
		volumes[id] = &data
	}
	return
}

// Groups get group infos.
func (d *Dao) Groups(cluster string) (groups map[string]*model.Group, err error) {
	zkConn, ok := d.zkcs[cluster]
	if !ok {
		err = fmt.Errorf("not support cluster:%s", cluster)
		return
	}
	root := d.c.Zookeepers[cluster].GroupRoot
	// get all groups
	groupNames, _, err := zkConn.Children(root)
	if err != nil {
		log.Error("cannot get from:%s info error(%v)", root, err)
		return
	}
	log.Info("got groups: %s", groupNames)
	groups = make(map[string]*model.Group)
	for _, groupName := range groupNames {
		// get all the volumes
		groupPath := path.Join(root, groupName)
		storeNames, _, err := zkConn.Children(groupPath)
		if err != nil {
			log.Error("cannot get from:%s info error(%v)", groupPath, err)
			continue
		}
		groups[groupName] = &model.Group{
			Stores: storeNames,
		}
	}
	return
}
