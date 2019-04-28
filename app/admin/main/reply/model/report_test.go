package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestReport(t *testing.T) {
	s := `{"code":0,"order":"[{\"ctime\":{\"order\":\"asc\"}}]","page":1,"pagesize":50,"pagecount":1,"total":1,"result":[{"id":121002129,"oid":10098721,"type":1,"report_mid":27515245,"reason":8,"content":"","count":1,"state":0,"score":0,"ctime":"2017-12-14 10:22:03","mtime":"2017-12-14 10:22:03","adminid":null,"opresult":null,"opremark":null,"opctime":null,"reply_mid":27515232,"root":0,"parent":0,"floor":2,"like":0,"reply_state":0,"message":"(=\u30fb\u03c9\u30fb=)","typeid":76,"arc_mid":27515615,"reporter":"Test000011","replier":"\u53ee\u5f53\u732b3333777","doc_id":"121002129_10098721_1"}]}`
	var data SearchReportResult
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
