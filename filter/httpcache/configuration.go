package httpcache

import (
	"regexp"
	"time"

	"github.com/darkweak/souin/configurationtypes"
	"github.com/darkweak/souin/plugins"
	"github.com/go-kratos/kratos/v2/config"
)

const (
	configurationKey = "httpcache"
	path             = "path"
	url              = "url"
	configurationPK  = "configuration"
)

func parseRecursively(values map[string]config.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range values {
		if v, e := value.Bool(); e == nil {
			result[key] = v
			continue
		}
		switch value.Load().(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			if v, e := value.Int(); e == nil {
				result[key] = v
				continue
			}
		case float32, float64:
			if v, e := value.Float(); e == nil {
				result[key] = v
				continue
			}
		}
		svalue, _ := value.String()
		if v, e := time.ParseDuration(svalue); e == nil {
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

func parseAPI(apiConfiguration map[string]config.Value) configurationtypes.API {
	var a configurationtypes.API
	var prometheusConfiguration, souinConfiguration map[string]config.Value

	for apiK, apiV := range apiConfiguration {
		switch apiK {
		case "prometheus":
			prometheusConfiguration, _ = apiV.Map()
		case "souin":
			souinConfiguration, _ = apiV.Map()
		}
	}
	if prometheusConfiguration != nil {
		a.Prometheus = configurationtypes.APIEndpoint{}
		a.Prometheus.Enable = true
		if prometheusConfiguration["basepath"] != nil {
			a.Prometheus.BasePath, _ = prometheusConfiguration["basepath"].String()
		}
	}
	if souinConfiguration != nil {
		a.Souin = configurationtypes.APIEndpoint{}
		a.Souin.Enable = true
		if souinConfiguration["basepath"] != nil {
			a.Souin.BasePath, _ = souinConfiguration["basepath"].String()
		}
	}

	return a
}

func parseCacheKeys(ccConfiguration map[string]config.Value) map[configurationtypes.RegValue]configurationtypes.Key {
	cacheKeys := make(map[configurationtypes.RegValue]configurationtypes.Key)
	for cacheKeysConfigurationK, cacheKeysConfigurationV := range ccConfiguration {
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

	return cacheKeys
}

func parseDefaultCache(dcConfiguration map[string]config.Value) *configurationtypes.DefaultCache {
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
	for defaultCacheK, defaultCacheV := range dcConfiguration {
		switch defaultCacheK {
		case "allowed_http_verbs":
			headers, _ := defaultCacheV.Slice()
			dc.AllowedHTTPVerbs = make([]string, 0)
			for _, header := range headers {
				h, _ := header.String()
				dc.AllowedHTTPVerbs = append(dc.AllowedHTTPVerbs, h)
			}
		case "badger":
			provider := configurationtypes.CacheProvider{}
			badgerConfiguration, _ := defaultCacheV.Map()
			for badgerConfigurationK, badgerConfigurationV := range badgerConfiguration {
				switch badgerConfigurationK {
				case url:
					provider.URL, _ = badgerConfigurationV.String()
				case path:
					provider.Path, _ = badgerConfigurationV.String()
				case configurationPK:
					configMap, e := badgerConfigurationV.Map()
					if e == nil {
						provider.Configuration = parseRecursively(configMap)
					}
				}
			}
			dc.Badger = provider
		case "cdn":
			cdn := configurationtypes.CDN{}
			cdnConfiguration, _ := defaultCacheV.Map()
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
			dc.CDN = cdn
		case "etcd":
			provider := configurationtypes.CacheProvider{}
			etcdConfiguration, _ := defaultCacheV.Map()
			for etcdConfigurationK, etcdConfigurationV := range etcdConfiguration {
				switch etcdConfigurationK {
				case url:
					provider.URL, _ = etcdConfigurationV.String()
				case path:
					provider.Path, _ = etcdConfigurationV.String()
				case configurationPK:
					configMap, e := etcdConfigurationV.Map()
					if e == nil {
						provider.Configuration = parseRecursively(configMap)
					}
				}
			}
			dc.Etcd = provider
		case "headers":
			headers, _ := defaultCacheV.Slice()
			dc.Headers = make([]string, 0)
			for _, header := range headers {
				h, _ := header.String()
				dc.Headers = append(dc.Headers, h)
			}
		case "nuts":
			provider := configurationtypes.CacheProvider{}
			nutsConfiguration, _ := defaultCacheV.Map()
			for nutsConfigurationK, nutsConfigurationV := range nutsConfiguration {
				switch nutsConfigurationK {
				case url:
					provider.URL, _ = nutsConfigurationV.String()
				case path:
					provider.Path, _ = nutsConfigurationV.String()
				case configurationPK:
					configMap, e := nutsConfigurationV.Map()
					if e == nil {
						provider.Configuration = parseRecursively(configMap)
					}
				}
			}
			dc.Nuts = provider
		case "olric":
			provider := configurationtypes.CacheProvider{}
			olricConfiguration, _ := defaultCacheV.Map()
			for olricConfigurationK, olricConfigurationV := range olricConfiguration {
				switch olricConfigurationK {
				case url:
					provider.URL, _ = olricConfigurationV.String()
				case path:
					provider.Path, _ = olricConfigurationV.String()
				case configurationPK:
					configMap, e := olricConfigurationV.Map()
					if e == nil {
						provider.Configuration = parseRecursively(configMap)
					}
				}
			}
			dc.Distributed = true
			dc.Olric = provider
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

	return &dc
}

func parseURLs(urls map[string]config.Value) map[string]configurationtypes.URL {
	u := make(map[string]configurationtypes.URL)

	for urlK, urlV := range urls {
		currentURL := configurationtypes.URL{
			TTL:     configurationtypes.Duration{},
			Headers: nil,
		}
		currentValue, _ := urlV.Map()
		if currentValue["headers"] != nil {
			currentURL.Headers = make([]string, 0)
			headers, _ := currentValue["headers"].Slice()
			for _, header := range headers {
				h, _ := header.String()
				currentURL.Headers = append(currentURL.Headers, h)
			}
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

	return u
}

func parseSurrogateKeys(surrogates map[string]config.Value) map[string]configurationtypes.SurrogateKeys {
	u := make(map[string]configurationtypes.SurrogateKeys)

	for surrogateK, surrogateV := range surrogates {
		surrogate := configurationtypes.SurrogateKeys{}
		currentValue, _ := surrogateV.Map()
		for key, value := range currentValue {
			switch key {
			case "headers":
				surrogate.Headers = map[string]string{}
				headers, e := value.Map()
				if e == nil {
					for hKey, hValue := range headers {
						v, _ := hValue.String()
						surrogate.Headers[hKey] = v
					}
				}
			case "url":
				surl, _ := currentValue["url"].String()
				surrogate.URL = surl
			}
		}
		u[surrogateK] = surrogate
	}

	return u
}

// ParseConfiguration parse the Kratos configuration into a valid HTTP
// cache configuration object.
func ParseConfiguration(c config.Config) plugins.BaseConfiguration {
	var configuration plugins.BaseConfiguration

	values, _ := c.Value(configurationKey).Map()
	for key, v := range values {
		switch key {
		case "api":
			apiConfiguration, _ := v.Map()
			configuration.API = parseAPI(apiConfiguration)
		case "cache_keys":
			cacheKeysConfiguration, _ := v.Map()
			configuration.CacheKeys = parseCacheKeys(cacheKeysConfiguration)
		case "default_cache":
			defaultCache, _ := v.Map()
			configuration.DefaultCache = parseDefaultCache(defaultCache)
		case "log_level":
			configuration.LogLevel, _ = v.String()
		case "urls":
			urls, _ := v.Map()
			configuration.URLs = parseURLs(urls)
		case "ykeys":
			ykeys, _ := v.Map()
			configuration.Ykeys = parseSurrogateKeys(ykeys)
		case "surrogate_keys":
			surrogates, _ := v.Map()
			configuration.SurrogateKeys = parseSurrogateKeys(surrogates)
		}
	}

	return configuration
}
