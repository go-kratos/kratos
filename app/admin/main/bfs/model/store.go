package model

// zookeeper save the store meta data.
//
//                                 /rack -- rack root path
//									 |
//                   /rack-a -------- --------- /rack-b -- rack node path
//                     |
//         /store-a -------- /store-b -- store node path (data: {"stat":"localhost:6061","admin":"localhost:6063","api":"localhost:6062","status":0})
//            |                 |
// /volume-1 -                   - /volume-4 volume node path (data: /tmp/block_1,/tmp/block_1.idx,1)
// /volume-2 -                   - /volume-5
// /volume-3 -                   - /volume-6

// Store status const.
const (
	// bit
	StoreStatusEnableBit = 31
	StoreStatusReadBit   = 0
	StoreStatusWriteBit  = 1
	StoreStatusSyncBit   = 2
	// status
	StoreStatusInit   = 0                                              // 0
	StoreStatusEnable = (1 << StoreStatusEnableBit)                    // 2147483648
	StoreStatusRead   = StoreStatusEnable | (1 << StoreStatusReadBit)  // 2147483649
	StoreStatusWrite  = StoreStatusEnable | (1 << StoreStatusWriteBit) // 2147483650
	StoreStatusHealth = StoreStatusRead | StoreStatusWrite             // 2147483651
	StoreStatusSync   = StoreStatusEnable | (1 << StoreStatusSyncBit)  // 2147483652
	StoreStatusFail   = StoreStatusEnable                              // 2147483648
)

// Rack get all store and volume.
type Rack struct {
	Stores map[string]*Store `json:"stores"`
}

// Store meta data.
type Store struct {
	Stat   string `json:"stat"`
	Admin  string `json:"admin"`
	API    string `json:"api"`
	ID     string `json:"id"`
	Rack   string `json:"rack"`
	Status int    `json:"status"`
	States struct {
		Init   bool `json:"init"`
		Enable bool `json:"enable"`
		Read   bool `json:"read"`
		Write  bool `json:"write"`
		Sync   bool `json:"sync"`
		Fail   bool `json:"fail"`
	} `json:"states"`
	Volumes []string `json:"volumes"`
}

// ParseStates parse states.
func (s *Store) ParseStates() {
	s.States.Init = s.Status == StoreStatusInit
	s.States.Enable = s.Status&StoreStatusEnable == StoreStatusEnable
	s.States.Read = s.Status&StoreStatusRead == StoreStatusRead
	s.States.Write = s.Status&StoreStatusWrite == StoreStatusWrite
	s.States.Sync = s.Status&StoreStatusSync == StoreStatusSync
	s.States.Fail = s.Status == StoreStatusFail
}
