package model

// type Volumes struct {
// 	Vol map[string]*VolumeState
// }

// VolumeState  for zk /volume stat
type VolumeState struct {
	TotalWriteProcessed uint64 `json:"total_write_processed"`
	TotalWriteDelay     uint64 `json:"total_write_delay"`
	FreeSpace           uint32 `json:"free_space"`
	Dir                 string `json:"dir"`
}

// func (v *Volumes) GetVolumeState(id string) *VolumeState {
// 	vs, ok := v.Vol[id]
// 	if !ok {
// 		return nil
// 	}
// 	return vs
// }
