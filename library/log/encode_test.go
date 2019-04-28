package log

import (
	"fmt"
	"testing"
	"time"

	"go-common/library/log/internal"
)

func TestJsonEncode(t *testing.T) {
	enc := core.NewJSONEncoder(core.EncoderConfig{
		EncodeTime:     core.EpochTimeEncoder,
		EncodeDuration: core.SecondsDurationEncoder,
	}, core.NewBuffer(0))
	KV("constant", "constant").AddTo(enc)
	for i := 0; i < 3; i++ {
		b := core.GetPool()
		err := enc.Encode(b, KV("no", i), KV("cat", "is cat"), KV("dog", time.Now()))
		if err != nil {
			t.Fatalf("enc.Encode error(%v)", err)
		}
		fmt.Println(string(b.Bytes()))
		b.Free()
	}
}
