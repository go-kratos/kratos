package memcache

import (
	"github.com/gogo/protobuf/proto"
)

// RawItem item with FlagRAW flag.
//
// Expiration is the cache expiration time, in seconds: either a relative
// time from now (up to 1 month), or an absolute Unix epoch time.
// Zero means the Item has no expiration time.
func RawItem(key string, data []byte, flags uint32, expiration int32) *Item {
	return &Item{Key: key, Flags: flags | FlagRAW, Value: data, Expiration: expiration}
}

// JSONItem item with FlagJSON flag.
//
// Expiration is the cache expiration time, in seconds: either a relative
// time from now (up to 1 month), or an absolute Unix epoch time.
// Zero means the Item has no expiration time.
func JSONItem(key string, v interface{}, flags uint32, expiration int32) *Item {
	return &Item{Key: key, Flags: flags | FlagJSON, Object: v, Expiration: expiration}
}

// ProtobufItem item with FlagProtobuf flag.
//
// Expiration is the cache expiration time, in seconds: either a relative
// time from now (up to 1 month), or an absolute Unix epoch time.
// Zero means the Item has no expiration time.
func ProtobufItem(key string, message proto.Message, flags uint32, expiration int32) *Item {
	return &Item{Key: key, Flags: flags | FlagProtobuf, Object: message, Expiration: expiration}
}
