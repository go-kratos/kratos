package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/polarismesh/polaris-go"
)

var (
	namespace     = "default"
	fileGroup     = "test"
	originContent = `server:
		port: 8080`
	updatedContent = `server:
		port: 8090`
	configCenterURL = "http://127.0.0.1:8090"
)

func makeJSONRequest(uri string, data string, method string, headers map[string]string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, uri, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
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
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	return &configClient{
		token: token,
	}, nil
}

func getToken() (string, error) {
	data, err := json.Marshal(map[string]string{
		"name":     "polaris",
		"password": "polaris",
	})
	if err != nil {
		return "", err
	}
	// login use default user
	res, err := makeJSONRequest(fmt.Sprintf("%s/core/v1/user/login", configCenterURL), string(data), http.MethodPost, map[string]string{})
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
		"namespace": namespace,
		"group":     fileGroup,
		"content":   originContent,
		"modifyBy":  "polaris",
		"format":    "yaml",
	})
	if err != nil {
		return err
	}
	res, err := makeJSONRequest(fmt.Sprintf("%s/config/v1/configfiles", configCenterURL), string(data), http.MethodPost, map[string]string{
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
		return errors.New("create error")
	}
	return nil
}

func (client *configClient) updateConfigFile(name string) error {
	data, err := json.Marshal(map[string]string{
		"name":      name,
		"namespace": namespace,
		"group":     fileGroup,
		"content":   updatedContent,
		"modifyBy":  "polaris",
		"format":    "yaml",
	})
	if err != nil {
		return err
	}
	res, err := makeJSONRequest(fmt.Sprintf("%s/config/v1/configfiles", configCenterURL), string(data), http.MethodPut, map[string]string{
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
		return errors.New("update error")
	}
	return nil
}

func (client *configClient) deleteConfigFile(name string) error {
	data, err := json.Marshal(map[string]string{})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/config/v1/configfiles?namespace=%s&group=%s&name=%s", configCenterURL, namespace, fileGroup, name)
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
		return errors.New("delete error")
	}
	return nil
}

func (client *configClient) publishConfigFile(name string) error {
	data, err := json.Marshal(map[string]string{
		"namespace": namespace,
		"group":     fileGroup,
		"fileName":  name,
		"name":      name,
	})
	if err != nil {
		return err
	}
	res, err := makeJSONRequest(fmt.Sprintf("%s/config/v1/configfiles/release", configCenterURL), string(data), http.MethodPost, map[string]string{
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
		return errors.New("publish error")
	}
	return nil
}

func TestConfig(t *testing.T) {
	name := "test.yaml"
	client, err := newConfigClient()
	if err != nil {
		t.Fatal(err)
	}
	if err = client.createConfigFile(name); err != nil {
		t.Fatal(err)
	}
	if err = client.publishConfigFile(name); err != nil {
		t.Fatal(err)
	}

	// Always remember clear test resource
	configAPI, err := polaris.NewConfigAPI()
	if err != nil {
		t.Fatal(err)
	}
	config, err := New(configAPI, WithNamespace(namespace), WithFileGroup(fileGroup), WithFileName(name))
	if err != nil {
		t.Fatal(err)
	}
	kv, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(kv) != 1 || kv[0].Key != name || string(kv[0].Value) != originContent {
		t.Fatal("config error")
	}

	w, err := config.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = client.deleteConfigFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = w.Next(); err != nil {
			t.Fatal(err)
		}
		if err = w.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	if err = client.updateConfigFile(name); err != nil {
		t.Fatal(err)
	}

	if err = client.publishConfigFile(name); err != nil {
		t.Fatal(err)
	}

	if kv, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	if len(kv) != 1 || kv[0].Key != name || string(kv[0].Value) != updatedContent {
		t.Fatal("config error")
	}
}

func TestExtToFormat(t *testing.T) {
	name := "ext.yaml"
	client, err := newConfigClient()
	if err != nil {
		t.Fatal(err)
	}
	if err = client.createConfigFile(name); err != nil {
		t.Fatal(err)
	}
	if err = client.publishConfigFile(name); err != nil {
		t.Fatal(err)
	}

	// Always remember clear test resource
	defer func() {
		if err = client.deleteConfigFile(name); err != nil {
			t.Fatal(err)
		}
	}()

	configAPI, err := polaris.NewConfigAPI()
	if err != nil {
		t.Fatal(err)
	}

	config, err := New(configAPI, WithNamespace(namespace), WithFileGroup(fileGroup), WithFileName(name))
	if err != nil {
		t.Fatal(err)
	}
	kv, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(len(kv), 1) {
		t.Errorf("len(kvs) = %d", len(kv))
	}
	if !reflect.DeepEqual(name, kv[0].Key) {
		t.Errorf("kvs[0].Key is %s", kv[0].Key)
	}
	if !reflect.DeepEqual(originContent, string(kv[0].Value)) {
		t.Errorf("kvs[0].Value is %s", kv[0].Value)
	}
	if !reflect.DeepEqual("yaml", kv[0].Format) {
		t.Errorf("kvs[0].Format is %s", kv[0].Format)
	}
}
