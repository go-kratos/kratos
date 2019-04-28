package model

import (
	"encoding/json"
	"sync"
)

//Item .
type Item struct {
	PercentMovementSpeedMod interface{} `json:"percent_movement_speed_mod"`
	PercentLifeStealMod     interface{} `json:"percent_life_steal_mod"`
	PercentAttackSpeedMod   interface{} `json:"percent_attack_speed_mod"`
	Name                    string      `json:"name"`
	ImageURL                string      `json:"image_url"`
	ID                      int64       `json:"id"`
	GoldTotal               interface{} `json:"gold_total"`
	GoldSell                interface{} `json:"gold_sell"`
	GoldPurchasable         bool        `json:"gold_purchasable"`
	GoldBase                interface{} `json:"gold_base"`
	FlatSpellBlockMod       interface{} `json:"flat_spell_block_mod"`
	FlatPhysicalDamageMod   interface{} `json:"flat_physical_damage_mod"`
	FlatMpRegenMod          interface{} `json:"flat_mp_regen_mod"`
	FlatMpPoolMod           interface{} `json:"flat_mp_pool_mod"`
	FlatMovementSpeedMod    interface{} `json:"flat_movement_speed_mod"`
	FlatMagicDamageMod      interface{} `json:"flat_magic_damage_mod"`
	FlatHpRegenMod          interface{} `json:"flat_hp_regen_mod"`
	FlatHpPoolMod           interface{} `json:"flat_hp_pool_mod"`
	FlatCritChanceMod       interface{} `json:"flat_crit_chance_mod"`
	FlatArmorMod            interface{} `json:"flat_armor_mod"`
}

//Game .
type Game struct {
	WinnerType string          `json:"winner_type"`
	Winner     json.RawMessage `json:"winner"`
	Teams      json.RawMessage `json:"teams"`
	Position   int64           `json:"position"`
	Players    json.RawMessage `json:"players"`
	MatchID    int64           `json:"match_id"`
	Match      json.RawMessage `json:"match"`
	Length     int64           `json:"length"`
	ID         int64           `json:"id"`
	Finished   interface{}     `json:"finished"`
	BeginAt    interface{}     `json:"begin_at"`
}

//Champion .
type Champion struct {
	VideogameVersions    []string    `json:"videogame_versions"`
	Spellblockperlevel   float64     `json:"spellblockperlevel"`
	Spellblock           float64     `json:"spellblock"`
	Name                 string      `json:"name"`
	Mpregenperlevel      float64     `json:"mpregenperlevel"`
	Mpregen              float64     `json:"mpregen"`
	Mpperlevel           interface{} `json:"mpperlevel"`
	Mp                   float64     `json:"mp"`
	Movespeed            interface{} `json:"movespeed"`
	ImageURL             string      `json:"image_url"`
	ID                   int64       `json:"id"`
	Hpregenperlevel      float64     `json:"hpregenperlevel"`
	Hpregen              interface{} `json:"hpregen"`
	Hpperlevel           interface{} `json:"hpperlevel"`
	Hp                   interface{} `json:"hp"`
	Critperlevel         interface{} `json:"critperlevel"`
	Crit                 interface{} `json:"crit"`
	BigImageURL          string      `json:"big_image_url"`
	Attackspeedperlevel  float64     `json:"attackspeedperlevel"`
	Attackspeedoffset    interface{} `json:"attackspeedoffset"`
	Attackrange          interface{} `json:"attackrange"`
	Attackdamageperlevel float64     `json:"attackdamageperlevel"`
	Attackdamage         float64     `json:"attackdamage"`
	Armorperlevel        float64     `json:"armorperlevel"`
	Armor                interface{} `json:"armor"`
}

//Hero .
type Hero struct {
	*LdInfo
	LocalizedName string `json:"localized_name"`
}

// LdInfo .
type LdInfo struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	ID       int64  `json:"id"`
}

// SyncGame store leida game list
type SyncGame struct {
	Data map[int64][]*Game
	sync.Mutex
}

// SyncItem store item list
type SyncItem struct {
	Data map[int64]*Item
	sync.Mutex
}

// SyncChampion store champion list
type SyncChampion struct {
	Data map[int64]*Champion
	sync.Mutex
}

// SyncHero store hero list
type SyncHero struct {
	Data map[int64]*Hero
	sync.Mutex
}

// SyncInfo store leida base info list
type SyncInfo struct {
	Data map[int64]*LdInfo
	sync.Mutex
}
