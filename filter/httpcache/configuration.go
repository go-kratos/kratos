package kratos

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/darkweak/souin/configurationtypes"
	"github.com/darkweak/souin/plugins"
	"github.com/go-kratos/kratos/v2/config"
)

const configuration_key = "httpcache"

func parseRecursively(values map[string]config.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range values {
		if v, e := value.Bool(); e == nil {
			result[key] = v
			continue
		}
		if v, e := value.Duration(); e == nil {
			result[key] = v
			continue
		}
		if v, e := value.Float(); e == nil {
			result[key] = v
			continue
		}
		if v, e := value.Int(); e == nil {
			result[key] = v
			continue
		}
		if v, e := value.Map(); e == nil {
			result[key] = parseRecursively(v)
			continue
		}
	}

	return result
}

func ParseConfiguration(c config.Config) plugins.BaseConfiguration {
	var configuration plugins.BaseConfiguration

	values, _ := c.Value(configuration_key).Map()
	for key, v := range values {
		switch key {
		case "api":
			var a configurationtypes.API
			var prometheusConfiguration, souinConfiguration, securityConfiguration map[string]config.Value
			apiConfiguration, _ := v.Map()
			for apiK, apiV := range apiConfiguration {
				switch apiK {
				case "prometheus":
					prometheusConfiguration, _ = apiV.Map()
				case "souin":
					souinConfiguration, _ = apiV.Map()
				case "security":
					securityConfiguration, _ = apiV.Map()
				}

			}
			if prometheusConfiguration != nil {
				a.Prometheus = configurationtypes.APIEndpoint{}
				a.Prometheus.Enable = true
				if prometheusConfiguration["basepath"] != nil {
					a.Prometheus.BasePath, _ = prometheusConfiguration["basepath"].String()
				}
				if prometheusConfiguration["security"] != nil {
					a.Prometheus.Security, _ = prometheusConfiguration["security"].Bool()
				}
			}
			if souinConfiguration != nil {
				a.Souin = configurationtypes.APIEndpoint{}
				a.Souin.Enable = true
				if souinConfiguration["basepath"] != nil {
					a.Souin.BasePath, _ = souinConfiguration["basepath"].String()
				}
				if souinConfiguration["security"] != nil {
					a.Souin.Security, _ = souinConfiguration["security"].Bool()
				}
			}
			if securityConfiguration != nil {
				a.Security = configurationtypes.SecurityAPI{}
				a.Security.Enable = true
				if securityConfiguration["basepath"] != nil {
					a.Security.BasePath, _ = securityConfiguration["basepath"].String()
				}
				if securityConfiguration["users"] != nil {
					users, _ := securityConfiguration["users"].Slice()
					a.Security.Users = make([]configurationtypes.User, 0)
					for _, user := range users {
						currentUser, _ := user.Map()
						username, _ := currentUser["username"].String()
						password, _ := currentUser["password"].String()
						a.Security.Users = append(a.Security.Users, configurationtypes.User{
							Username: username,
							Password: password,
						})
					}
				}
			}
			configuration.API = a
		case "cache_keys":
			cacheKeys := make(map[configurationtypes.RegValue]configurationtypes.Key)
			cacheKeysConfiguration, _ := v.Map()
			for cacheKeysConfigurationK, cacheKeysConfigurationV := range cacheKeysConfiguration {
				ck := configurationtypes.Key{}
				cacheKeysConfigurationVMap, _ := cacheKeysConfigurationV.Map()
				for cacheKeysConfigurationVMapK := range cacheKeysConfigurationVMap {
					switch cacheKeysConfigurationVMapK {
					case "disable_body":
						ck.DisableBody = true
					case "disable_host":
						ck.DisableHost = true
					case "disable_method":
						ck.DisableMethod = true
					}
				}
				rg := regexp.MustCompile(cacheKeysConfigurationK)
				cacheKeys[configurationtypes.RegValue{Regexp: rg}] = ck
			}
			configuration.CacheKeys = cacheKeys
		case "default_cache":
			dc := configurationtypes.DefaultCache{
				Distributed: false,
				Headers:     []string{},
				Olric: configurationtypes.CacheProvider{
					URL:           "",
					Path:          "",
					Configuration: nil,
				},
				Regex:               configurationtypes.Regex{},
				TTL:                 configurationtypes.Duration{},
				DefaultCacheControl: "",
			}
			defaultCache, _ := v.Map()
			for defaultCacheK, defaultCacheV := range defaultCache {
				switch defaultCacheK {
				case "badger":
					provider := configurationtypes.CacheProvider{}
					badgerConfiguration, _ := v.Map()
					for badgerConfigurationK, badgerConfigurationV := range badgerConfiguration {
						switch badgerConfigurationK {
						case "url":
							provider.URL, _ = badgerConfigurationV.String()
						case "path":
							provider.Path, _ = badgerConfigurationV.String()
						case "configuration":
							configMap, e := badgerConfigurationV.Map()
							if e == nil {
								provider.Configuration = parseRecursively(configMap)
							}
						}
					}
					configuration.DefaultCache.Badger = provider
				case "cdn":
					cdn := configurationtypes.CDN{}
					cdnConfiguration, _ := v.Map()
					for cdnConfigurationK, cdnConfigurationV := range cdnConfiguration {
						switch cdnConfigurationK {
						case "api_key":
							cdn.APIKey, _ = cdnConfigurationV.String()
						case "dynamic":
							cdn.Dynamic = true
						case "hostname":
							cdn.Hostname, _ = cdnConfigurationV.String()
						case "network":
							cdn.Network, _ = cdnConfigurationV.String()
						case "provider":
							cdn.Provider, _ = cdnConfigurationV.String()
						case "strategy":
							cdn.Strategy, _ = cdnConfigurationV.String()
						}
					}
					configuration.DefaultCache.CDN = cdn
				case "etcd":
					provider := configurationtypes.CacheProvider{}
					etcdConfiguration, _ := v.Map()
					for etcdConfigurationK, etcdConfigurationV := range etcdConfiguration {
						switch etcdConfigurationK {
						case "url":
							provider.URL, _ = etcdConfigurationV.String()
						case "path":
							provider.Path, _ = etcdConfigurationV.String()
						case "configuration":
							configMap, e := etcdConfigurationV.Map()
							if e == nil {
								provider.Configuration = parseRecursively(configMap)
							}
						}
					}
					configuration.DefaultCache.Etcd = provider
				case "headers":
					headers, _ := defaultCacheV.Slice()
					dc.Headers = make([]string, 0)
					for _, header := range headers {
						h, _ := header.String()
						dc.Headers = append(dc.Headers, h)
					}
				case "nuts":
					provider := configurationtypes.CacheProvider{}
					nutsConfiguration, _ := v.Map()
					for nutsConfigurationK, nutsConfigurationV := range nutsConfiguration {
						switch nutsConfigurationK {
						case "url":
							provider.URL, _ = nutsConfigurationV.String()
						case "path":
							provider.Path, _ = nutsConfigurationV.String()
						case "configuration":
							configMap, e := nutsConfigurationV.Map()
							if e == nil {
								provider.Configuration = parseRecursively(configMap)
							}
						}
					}
					configuration.DefaultCache.Nuts = provider
				case "olric":
					provider := configurationtypes.CacheProvider{}
					olricConfiguration, _ := v.Map()
					for olricConfigurationK, olricConfigurationV := range olricConfiguration {
						switch olricConfigurationK {
						case "url":
							provider.URL, _ = olricConfigurationV.String()
						case "path":
							provider.Path, _ = olricConfigurationV.String()
						case "configuration":
							configMap, e := olricConfigurationV.Map()
							if e == nil {
								provider.Configuration = parseRecursively(configMap)
							}
						}
					}
					configuration.DefaultCache.Distributed = true
					configuration.DefaultCache.Olric = provider
				case "regex":
					regex, _ := defaultCacheV.Map()
					exclude, _ := regex["exclude"].String()
					if exclude != "" {
						dc.Regex = configurationtypes.Regex{Exclude: exclude}
					}
				case "ttl":
					sttl, err := defaultCacheV.String()
					ttl, _ := time.ParseDuration(sttl)
					if err == nil {
						dc.TTL = configurationtypes.Duration{Duration: ttl}
					}
				case "stale":
					sstale, err := defaultCacheV.String()
					stale, _ := time.ParseDuration(sstale)
					if err == nil {
						dc.Stale = configurationtypes.Duration{Duration: stale}
					}
				case "default_cache_control":
					dc.DefaultCacheControl, _ = defaultCacheV.String()
				}
			}
			configuration.DefaultCache = &dc
		case "log_level":
			configuration.LogLevel, _ = v.String()
		case "urls":
			u := make(map[string]configurationtypes.URL)
			urls, _ := v.Map()

			for urlK, urlV := range urls {
				currentURL := configurationtypes.URL{
					TTL:     configurationtypes.Duration{},
					Headers: nil,
				}
				currentValue, _ := urlV.Map()
				currentURL.Headers = make([]string, 0)
				headers, _ := currentValue["headers"].Slice()
				for _, header := range headers {
					h, _ := header.String()
					currentURL.Headers = append(currentURL.Headers, h)
				}
				sttl, err := currentValue["ttl"].String()
				ttl, _ := time.ParseDuration(sttl)
				if err == nil {
					currentURL.TTL = configurationtypes.Duration{Duration: ttl}
				}
				if _, exists := currentValue["default_cache_control"]; exists {
					currentURL.DefaultCacheControl, _ = currentValue["default_cache_control"].String()
				}
				u[urlK] = currentURL
			}
			configuration.URLs = u
		case "ykeys":
			ykeys := make(map[string]configurationtypes.SurrogateKeys)
			d, _ := json.Marshal(v)
			_ = json.Unmarshal(d, &ykeys)
			configuration.Ykeys = ykeys
		case "surrogate_keys":
			ykeys := make(map[string]configurationtypes.SurrogateKeys)
			d, _ := json.Marshal(v)
			_ = json.Unmarshal(d, &ykeys)
			configuration.Ykeys = ykeys
		}
	}

	return configuration
}
