package vault

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/vault/api"
)

const testPath = "secret/data/db/auth"

func TestVaultRead(t *testing.T) {
	vaultValue := map[string]interface{}{"user": "myuser", "password": "mypassword2"}

	client, err := api.NewClient(&api.Config{
		Address: "http://127.0.0.1:8200",
	})
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("12ff741b-4438-6cb5-63f0-0b6a11a3f4cc")

	if _, err := client.Logical().Write(testPath, vaultValue); err != nil {
		t.Fatal(err)
	}

	readSecret, err := client.Logical().Read(testPath)
	if !reflect.DeepEqual(readSecret.Data, vaultValue) {
		t.Fatal("Not read the correct value")
	}

	source, err := New(client, WithPath(testPath))
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(kvs)

	if reflect.DeepEqual(kvs, vaultValue) {
		t.Fatal("config error")
	}

	if _, err := client.Logical().Delete(testPath); err != nil {
		t.Error(err)
	}
}

func TestConfigWithVault(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "http://127.0.0.1:8200",
	})
	if err != nil {
		t.Fatal(err)
	}

	vaultSrc, err := New(client, WithPath(testPath))
	if err != nil {
		t.Fatal(err)
	}
	c := config.New(config.WithSource(vaultSrc))

	vaultValue := map[string]interface{}{"user": "myuser", "password": "mypassword2"}

	if _, err := client.Logical().Write(testPath, vaultValue); err != nil {
		t.Fatal(err)
	}

	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	name, err := c.Value("user").String()
	fmt.Println(name)

	if _, err := client.Logical().Delete(testPath); err != nil {
		t.Error(err)
	}
}
