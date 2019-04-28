package http

import (
	"fmt"
	"math/rand"
	"time"

	bm "go-common/library/net/http/blademaster"

	"github.com/Dai0522/go-hash/murmur3"
)

// SessionID .
func SessionID(ctx *bm.Context) string {
	ts := time.Now().Unix()
	dev, _ := ctx.Get("device")
	buvid := dev.(*bm.Device).Buvid
	rnum := rand.Uint64()

	str := fmt.Sprintf("%s:%d:%d", buvid, rnum, ts)
	hc := murmur3.New().Murmur3_128([]byte(str))

	result := ""
	pattern := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for _, v := range hc {
		tmp := int(v & 15)
		result += pattern[tmp]
		tmp = int((v >> 4) & 15)
		result += pattern[tmp]
	}

	return result
}
