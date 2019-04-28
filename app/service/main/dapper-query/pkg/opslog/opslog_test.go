package opslog

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestOpsLog(t *testing.T) {
	testSessionID := "c860e25e5360fc08888a3aaf8c7a0bec"
	testResponse := `{"responses":[{"took":4,"timed_out":false,"_shards":{"total":8,"successful":8,"failed":0},"hits":{"total":2,"max_score":null,"hits":[{"_index":"billions-main.web-svr.web-interface-@2018.09.10-uat-1","_type":"logs","_id":"AWXBhTtRe-NhC44S955A","_version":1,"_score":null,"_source":{"@timestamp":"2018-09-10T03:27:34.42933Z","app_id":"main.web-svr.web-interface","args":"","env":"uat","error":"","instance_id":"web-interface-32096-758958f64f-k4gnc","ip":"172.22.35.133:9000","level":"INFO","level_value":1,"path":"/passport.service.identify.v1.Identify/GetCookieInfo","ret":0,"source":"go-common/library/net/rpc/warden.logging:195","stack":"\u003cnil\u003e","traceid":"2406767965117552819","ts":0.001041696,"user":"","zone":"sh001"},"fields":{"@timestamp":[1536550054429]},"highlight":{"traceid":["@kibana-highlighted-field@2406767965117552819@/kibana-highlighted-field@"]},"sort":[1536550054429]},{"_index":"billions-main.web-svr.web-interface-@2018.09.10-uat-1","_type":"logs","_id":"AWXBhTfFS1y0J6vacgAH","_version":1,"_score":null,"_source":{"@timestamp":"2018-09-10T03:27:34.429376Z","app_id":"main.web-svr.web-interface","env":"uat","err":"-101","instance_id":"web-interface-32096-758958f64f-k4gnc","ip":"10.23.50.21","level":"ERROR","level_value":3,"method":"GET","mid":null,"msg":"账号未登录","params":"","path":"/x/web-interface/nav","ret":-101,"source":"go-common/library/net/http/blademaster.Logger.func1:46","stack":"-101","traceid":"2406767965117552819","ts":0.001167299,"user":"no_user","zone":"sh001"},"fields":{"@timestamp":[1536550054429]},"highlight":{"traceid":["@kibana-highlighted-field@2406767965117552819@/kibana-highlighted-field@"]},"sort":[1536550054429]}]},"aggregations":{"2":{"buckets":[{"key_as_string":"2018-09-10T09:00:00.000+08:00","key":1536541200000,"doc_count":2}]}},"status":200}]}`
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie(_ajsSessioID)
		if err != nil || session.Value != testSessionID {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid session id: %s", session.Value)
			return
		}
		bufReader := bufio.NewReader(r.Body)
		first, err := bufReader.ReadString('\n')
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !strings.Contains(first, "billions-main.web-svr.web-interface*") {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid familys: %s", first)
			return
		}
		w.Write([]byte(testResponse))
	}))
	defer svr.Close()
	client := New(svr.URL, nil)
	end := time.Now().Unix()
	start := end - 3600
	familys := []string{"main.web-svr.web-interface"}
	records, err := client.Query(context.Background(), familys, 8111326167741382285, testSessionID, start, end)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		t.Logf("record: %v", record)
	}
}

func TestOpsLogReal(t *testing.T) {
	sessionID := os.Getenv("TEST_SESSION_ID")
	if sessionID == "" {
		t.Skipf("miss sessionID skip test")
	}

	traceID, _ := strconv.ParseUint("7b91b9a72f87c13", 16, 64)
	client := New("http://uat-ops-log.bilibili.co/elasticsearch/_msearch", nil)
	records, err := client.Query(context.Background(), []string{"main.community.tag"}, traceID, sessionID, 1545296000, 1545296286)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		t.Logf("record: %v", record)
	}
}
