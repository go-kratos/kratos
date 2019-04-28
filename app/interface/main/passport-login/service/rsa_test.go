package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecode(t *testing.T) {

	Convey("decode", t, func() {
		once.Do(startService)
		s := "b9hmhL481b4+1yPLFIejdELcCcYrA/Fqq/aZAtqUGPz0ljbSQGIgUei0cNWqE5PFQakOrf6+IF/OwlUnyfp44+Ocj6KKNbl1QGPFPgGaBb8coPCanHSpm4Ap6iszh/cok56qsAWSOLgW/7fiO3RFzI6JMj3trgViscxUWkK4Ofg="
		key, err := priKey([]byte(privateKey))
		fmt.Println(err, key)
		d, _ := base64.StdEncoding.DecodeString(s)
		b, err := rsaDecryptPKCS8(key, d)
		pwd := string(b)
		fmt.Println(err, pwd)
		timestamp := pwd[:16]
		pwdStr := pwd[16:]
		fmt.Println(timestamp)
		tt, _ := Hash2TsSeconds(timestamp)
		fmt.Println(tt, pwdStr)
	})
}
func TestService_RSAKey(t *testing.T) {
	once.Do(startService)
	Convey("Test check user", t, func() {
		var c = context.TODO()
		res := s.RSAKey(c)
		So(res, ShouldNotBeNil)

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}
