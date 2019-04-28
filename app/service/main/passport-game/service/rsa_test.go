package service

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	s := "Qiu+Iyqi/URq8fbd/FQAtN8CPJ/LT9AAG5/dSd+RSRjH4OgGp3tJMvOk7uqzzgWDi0EYDsLp2xme61NdjhITz1eRmN11oyu2cYqsk9Rqp0o3LzNHeMuUL4GMdWlhRNBQMoNJdzmSaPZoIlgL3vxK+ghg/zF8clqvMSYryuy+rZI="
	key, err := privKey([]byte(_originPrivKey))
	fmt.Println(err)
	d, _ := base64.StdEncoding.DecodeString(s)
	b, err := rsaDecryptPKCS8(key, d)
	fmt.Println(err)
	fmt.Println(string(b))
}

func TestRSADecryptPKCS8(t *testing.T) {
	plain := []byte("hello")
	pub, err := pubKey([]byte(_originPubKey))
	if err != nil {
		t.Errorf("failed to init pub key, error(%v)", err)
		t.FailNow()
	}
	d, err := rsaEncryptPKCS8(pub, plain)
	if err != nil {
		t.Errorf("failed to rsa encrypt, error(%v)", err)
		t.FailNow()
	}
	t.Logf("d: %s", base64.StdEncoding.EncodeToString(d))
	priv, err := privKey([]byte(_originPrivKey))
	if err != nil {
		t.Errorf("failed to init priv key, error(%v)", err)
		t.FailNow()
	}
	fn := "/tmp/pwd.txt"
	f, err := os.Open(fn)
	if err != nil {
		t.Errorf("failed to open file %s, error(%v)", fn, err)
		t.FailNow()
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	i := 0
	for {
		line, err := rd.ReadString('\n')

		if err != nil || io.EOF == err {
			break
		}
		d, _ = base64.StdEncoding.DecodeString(line)
		p, err := rsaDecryptPKCS8(priv, d)
		if err != nil {
			t.Errorf("failed to rsa decrypt, error(%v)", err)
			t.FailNow()
		}
		i++
		t.Logf("%d: %s\n", i, p[16:])
	}
}
