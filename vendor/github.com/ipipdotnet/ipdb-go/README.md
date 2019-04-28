# ipdb-go
IPIP.net officially supported IP database ipdb format parsing library

# Installing
<code>
    go get github.com/ipipdotnet/ipdb-go
</code>

# Example
<pre>
<code>
package main

import (
	"github.com/ipipdotnet/ipdb-go"
	"fmt"
	"log"
)

func main() {

    // 支持IPDB格式地级市精度IP离线库(免费版，每周高级版，每日标准版，每日高级版，每日专业版，每日旗舰版)
	db, err := ipdb.NewCity("/path/to/city.ipv4.ipdb")
	if err != nil {
		log.Fatal(err)
	}

	db.Reload("/path/to/city.ipv4.ipdb") // 更新 ipdb 文件后可调用 Reload 方法重新加载内容

	fmt.Println(db.IsIPv4()) // check database support ip type
	fmt.Println(db.IsIPv6()) // check database support ip type
	fmt.Println(db.BuildTime()) // database build time
	fmt.Println(db.Languages()) // database support language
	fmt.Println(db.Fields()) // database support fields

	fmt.Println(db.FindInfo("2001:250:200::", "CN")) // return CityInfo
	fmt.Println(db.Find("1.1.1.1", "CN")) // return []string
	fmt.Println(db.FindMap("118.28.8.8", "CN")) // return map[string]string
	fmt.Println(db.FindInfo("127.0.0.1", "CN")) // return CityInfo

	fmt.Println()
}
</code>
</pre>