package polaris

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/config"

	"github.com/polarismesh/polaris-go"
)

var (
	testNamespace     = "default"
	testFileGroup     = "test"
	testOriginContent = `server:
		port: 8080`
	testUpdatedContent = `server:
		port: 8090`
	testCenterURL = "http://127.0.0.1:8090"
)

func makeJSONRequest(uri string, data string, method string, headers map[string]string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, uri, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

type commonRes struct {
	Code int32 `json:"code"`
}

type LoginRes struct {
	Code          int32 `json:"code"`
	LoginResponse struct {
		Token string `json:"token"`
	} `json:"loginResponse"`
}

type configClient struct {
	token string
}

func newConfigClient() (*configClient, error) {
	token, err := getToken(testCenterURL)
	if err != nil {
		return nil, err
	}
	return &configClient{
		token: token,
	}, nil
}

func getToken(testCenterURL string) (string, error) {
	data, err := json.Marshal(map[string]string{
		"name":     "polaris",
		"password": "polaris",
	})
	if err != nil {
		return "", err
	}
	// login use default user
	res, err := makeJSONRequest(fmt.Sprintf("%s/core/v1/user/login", testCenterURL), string(data), http.MethodPost, map[string]string{})
	if err != nil {
		return "", nil
	}
	var loginRes LoginRes
	if err = json.Unmarshal(res, &loginRes); err != nil {
		return "", err
	}
	return loginRes.LoginResponse.Token, nil
}

func (client *configClient) createConfigFile(name string) error {
	data, err := json.Marshal(map[string]string{
		"name":      name,
		"namespace": testNamespace,
		"group":     testFileGroup,
		"content":   testOriginContent,
		"modifyBy":  "polaris",
		"format":    "yaml",
	})
	if err != nil {
		return err
	}
	res, err := makeJSONRequest(fmt.Sprintf("%s/config/v1/configfiles", testCenterURL), string(data), http.MethodPost, map[string]string{
		"X-Polaris-Token": client.token,
	})
	if err != nil {
		return err
	}

	var resJSON commonRes
	err = json.Unmarshal(res, &resJSON)
	if err != nil {
		return err
	}
	if resJSON.Code != 200000 {
		return fmt.Errorf("create error, res: %s", string(res))
	}
	return nil
}

func (client *configClient) updateConfigFile(name string) error {
	data, err := json.Marshal(map[string]string{
		"name":      name,
		"namespace": testNamespace,
		"group":     testFileGroup,
		"content":   testUpdatedContent,
		"modifyBy":  "polaris",
		"format":    "yaml",
	})
	if err != nil {
		return err
	}
	res, err := makeJSONRequest(fmt.Sprintf("%s/config/v1/configfiles", testCenterURL), string(data), http.MethodPut, map[string]string{
		"X-Polaris-Token": client.token,
	})
	if err != nil {
		return err
	}
	var resJSON commonRes
	err = json.Unmarshal(res, &resJSON)
	if err != nil {
		return err
	}
	if resJSON.Code != 200000 {
		return fmt.Errorf("update error, res: %s", string(res))
	}
	return nil
}

func (client *configClient) deleteConfigFile(name string) error {
	data, err := json.Marshal(map[string]string{})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/config/v1/configfiles?namespace=%s&group=%s&name=%s", testCenterURL, testNamespace, testFileGroup, name)
	res, err := makeJSONRequest(url, string(data), http.MethodDelete, map[string]string{
		"X-Polaris-Token": client.token,
	})
	if err != nil {
		return err
	}
	var resJSON commonRes
	err = json.Unmarshal(res, &resJSON)
	if err != nil {
		return err
	}
	if resJSON.Code != 200000 {
		return fmt.Errorf("delete error, res: %s", string(res))
	}
	return nil
}

func (client *configClient) publishConfigFile(name string) error {
	data, err := json.Marshal(map[string]string{
		"namespace": testNamespace,
		"group":     testFileGroup,
		"fileName":  name,
		"name":      name,
	})
	if err != nil {
		return err
	}
	res, err := makeJSONRequest(fmt.Sprintf("%s/config/v1/configfiles/release", testCenterURL), string(data), http.MethodPost, map[string]string{
		"X-Polaris-Token": client.token,
	})
	if err != nil {
		return err
	}
	var resJSON commonRes
	err = json.Unmarshal(res, &resJSON)
	if err != nil {
		return err
	}
	if resJSON.Code != 200000 {
		return fmt.Errorf("publish error, res: %s", string(res))
	}
	return nil
}

func TestConfig(t *testing.T) {
	name := "kratos-polaris-test.yaml"
	client, err := newConfigClient()
	if err != nil {
		t.Fatal(err)
	}
	_ = client.deleteConfigFile(name)
	if err = client.createConfigFile(name); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	if err = client.publishConfigFile(name); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	// Always remember clear test resource
	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}
	p := New(sdk)
	config, err := p.Config(WithConfigFile(File{Name: name, Group: testFileGroup}))
	if err != nil {
		t.Fatal(err)
	}
	kv, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}

	for _, value := range kv {
		t.Logf("key: %s, value: %s", value.Key, value.Value)
	}
	if len(kv) != 1 || kv[0].Key != name || string(kv[0].Value) != testOriginContent {
		t.Fatal("config error")
	}

	w, err := config.Watch()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = client.deleteConfigFile(name)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err = client.updateConfigFile(name); err != nil {
		t.Fatal(err)
	}

	if err = client.publishConfigFile(name); err != nil {
		t.Fatal(err)
	}

	if kv, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	for _, value := range kv {
		t.Log(value.Key, string(value.Value))
	}

	if len(kv) != 1 || kv[0].Key != name || string(kv[0].Value) != testUpdatedContent {
		t.Fatal("config error")
	}
}

func TestExtToFormat(t *testing.T) {
	name := "kratos-polaris-ext.yaml"
	client, err := newConfigClient()
	if err != nil {
		t.Fatal(err)
	}
	_ = client.deleteConfigFile(name)
	if err = client.createConfigFile(name); err != nil {
		t.Fatal(err)
	}
	if err = client.publishConfigFile(name); err != nil {
		t.Fatal(err)
	}

	// Always remember clear test resource
	t.Cleanup(func() {
		if err = client.deleteConfigFile(name); err != nil {
			t.Fatal(err)
		}
	})

	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}
	p := New(sdk)

	cfg, err := p.Config(WithConfigFile(File{Name: name, Group: testFileGroup}))
	if err != nil {
		t.Fatal(err)
	}

	kv, err := cfg.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(len(kv), 1) {
		t.Errorf("len(kvs) = %d", len(kv))
	}
	if !reflect.DeepEqual(name, kv[0].Key) {
		t.Errorf("kvs[0].Key is %s", kv[0].Key)
	}
	if !reflect.DeepEqual(testOriginContent, string(kv[0].Value)) {
		t.Errorf("kvs[0].Value is %s", kv[0].Value)
	}
	if !reflect.DeepEqual("yaml", kv[0].Format) {
		t.Errorf("kvs[0].Format is %s", kv[0].Format)
	}
}

func TestGetMultipleConfig(t *testing.T) {
	client, err := newConfigClient()
	files := make([]File, 0, 3)
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("kratos-polaris-test-%d.yaml", i)
		if err != nil {
			t.Fatal(err)
		}
		_ = client.deleteConfigFile(name)
		if err = client.createConfigFile(name); err != nil {
			t.Fatal(err)
		}
		if err = client.publishConfigFile(name); err != nil {
			t.Fatal(err)
		}
		files = append(files, File{Name: name, Group: testFileGroup})
	}

	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}
	p := New(sdk, WithNamespace("default"))

	cfg, err := p.Config(WithConfigFile(files...))
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := cfg.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, kv := range kvs {
		t.Logf("key: %s, value: %s", kv.Key, kv.Value)
	}

	w, err := cfg.Watch()
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if err = client.publishConfigFile(file.Name); err != nil {
			t.Fatal(err)
		}
		kvs, err := w.Next()
		if err != nil {
			t.Fatal(err)
		}
		m := make(map[string]*config.KeyValue)
		for _, kv := range kvs {
			m[kv.Key] = kv
		}
		if !reflect.DeepEqual(file.Name, m[file.Name].Key) {
			t.Errorf("m[file.Name].Key is %s", m[file.Name].Key)
		}
		if !reflect.DeepEqual(testOriginContent, string(m[file.Name].Value)) {
			t.Errorf("m[file.Name].Value is %s", m[file.Name].Value)
		}
		if !reflect.DeepEqual("yaml", m[file.Name].Format) {
			t.Errorf("m[file.Name].Format is %s", m[file.Name].Format)
		}
	}
}
