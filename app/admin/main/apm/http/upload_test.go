package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	convey.Convey("upload", t, func() {
		url := "http://127.0.0.1:7331/x/admin/apm/ut/upload"
		payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; " +
			"name=\"html_file\"\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; " +
			"name=\"mid\"\r\n\r\n6\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; " +
			"name=\"username\"\r\n\r\nchengxing\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: " +
			"form-data; name=\"report_file\"\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
		req.Header.Add("Cookie", "_AJSESSIONID=cf81d40c0e9d960a0ce89ceeb05c5670; username=chengxing; "+
			"sven-apm=4104f6b8cb1d967a0dd45d6934638ba2bfc86cd239bf7bab095b8a1cc3f85b65")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("Postman-Token", "69b1317b-0b01-43a0-a85d-c899f64ae34e")

		res, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		fmt.Println(string(body))

	})
}
