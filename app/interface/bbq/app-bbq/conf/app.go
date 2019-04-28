package conf

import "github.com/BurntSushi/toml"

// AppSetting .
type AppSetting map[string]interface{}

// Set .
func (a *AppSetting) Set(text string) error {
	if _, err := toml.Decode(text, a); err != nil {
		panic(err)
	}
	return nil
}
