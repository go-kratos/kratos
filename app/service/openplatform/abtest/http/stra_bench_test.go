package http

import (
	"context"
	"fmt"
	"net/url"
	"testing"
)

func BenchmarkVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		params := url.Values{}
		for k, v := range guscs[0].testData {
			params.Set(k, v)
		}
		req, _ := client.NewRequest("GET", _getVersionURL, "127.0.0.1", params)
		var res struct {
			Code int `json:"code"`
			Data struct {
				V int         `json:"v"`
				D interface{} `json:"d"`
			} `json:"data"`
		}

		if err := client.Do(context.TODO(), req, &res); err != nil {
			fmt.Println(err)
		}
	}
}
