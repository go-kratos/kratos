package apollo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/conf/paladin/apollo/internal/mockserver"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	go func() {
		if err := mockserver.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	// wait for mock server to run
	time.Sleep(time.Millisecond * 500)
}

func teardown() {
	mockserver.Close()
}

func TestApollo(t *testing.T) {
	var (
		testAppYAML           = "app.yml"
		testAppYAMLContent1   = "test: test12234\ntest2: test333"
		testAppYAMLContent2   = "test: 1111"
		testClientJSON        = "client.json"
		testClientJSONContent = `{"name":"agollo"}`
	)
	os.Setenv("APOLLO_APP_ID", "SampleApp")
	os.Setenv("APOLLO_CLUSTER", "default")
	os.Setenv("APOLLO_CACHE_DIR", "/tmp")
	os.Setenv("APOLLO_META_ADDR", "localhost:8010")
	os.Setenv("APOLLO_NAMESPACES", fmt.Sprintf("%s,%s", testAppYAML, testClientJSON))
	mockserver.Set(testAppYAML, "content", testAppYAMLContent1)
	mockserver.Set(testClientJSON, "content", testClientJSONContent)
	ad := &apolloDriver{}
	apollo, err := ad.New()
	if err != nil {
		t.Fatalf("new apollo error, %v", err)
	}
	value := apollo.Get(testAppYAML)
	if content, _ := value.String(); content != testAppYAMLContent1 {
		t.Fatalf("got app.yml unexpected value %s", content)
	}
	value = apollo.Get(testClientJSON)
	if content, _ := value.String(); content != testClientJSONContent {
		t.Fatalf("got app.yml unexpected value %s", content)
	}
	mockserver.Set(testAppYAML, "content", testAppYAMLContent2)
	updates := apollo.WatchEvent(context.TODO(), testAppYAML)
	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}
	value = apollo.Get(testAppYAML)
	if content, _ := value.String(); content != testAppYAMLContent2 {
		t.Fatalf("got app.yml unexpected updated value %s", content)
	}
}
