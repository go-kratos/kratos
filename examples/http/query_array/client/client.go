package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	tests := []string{"names=1,2,3,4", "names[]=kratos,go-kratos", "names=1,3&&age=15&names=2,4", "names[]=1,2&names=3,4"}
	for _, test := range tests {
		httpGet(test)
	}
}

func httpGet(v string) {
	resp, err := http.Get("http://127.0.0.1:8080/hello?" + v)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v + " ---> " + string(body))
}
