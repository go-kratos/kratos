package module

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
)

// Setting ModuleConf
type Setting struct {
	Timestamp int64           `json:"timestamp"`
	CheckSum  int64           `json:"check_sum"`
	Data      json.RawMessage `json:"data"`
}

const (
	_playerID     = 1
	_playerKey    = "player"
	_webPlayerID  = 2
	_webPlayerKey = "web_player"
)

var _moduleMap = map[string]int{
	_playerKey:    _playerID,
	_webPlayerKey: _webPlayerID,
}

// VerifyModuleKey verify key
func VerifyModuleKey(key string) int {
	return _moduleMap[key]
}

// Result get module message
func Result(moduleKeyID int, data string) (rm json.RawMessage, checkSum int64, err error) {
	var (
		bs []byte
	)
	switch moduleKeyID {
	case _playerID:
		player := &Player{}
		playerSha1 := &PlayerSha1{}
		err = json.Unmarshal([]byte(data), player)
		if err != nil {
			return
		}
		bs, err = json.Marshal(player)
		if err != nil {
			return
		}
		rm = json.RawMessage(bs)
		// check_sum
		err = json.Unmarshal([]byte(data), playerSha1)
		if err != nil {
			return
		}
		bs, err = json.Marshal(playerSha1)
		if err != nil {
			return
		}
		checkSum = int64(crc32.ChecksumIEEE(bs))
		return
	case _webPlayerID:
		player := &WebPlayer{}
		playerSha1 := &WebPlayerSha1{}
		err = json.Unmarshal([]byte(data), player)
		if err != nil {
			return
		}
		bs, err = json.Marshal(player)
		if err != nil {
			return
		}
		rm = json.RawMessage(bs)
		// check_sum
		err = json.Unmarshal([]byte(data), playerSha1)
		if err != nil {
			return
		}
		bs, err = json.Marshal(playerSha1)
		if err != nil {
			return
		}
		checkSum = int64(crc32.ChecksumIEEE(bs))
		return
	}
	err = fmt.Errorf("module_key_id not found: %v", moduleKeyID)
	return
}
