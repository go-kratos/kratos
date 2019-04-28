package model

import (
	"net/url"
	"reflect"
)

const _query = "query"

// TreeResponse Tree Response.
type TreeResponse struct {
	Code int      `json:"code"`
	Data UserTree `json:"data"`
}

// UserTree User Tree.
type UserTree struct {
	Bilibili map[string]interface{} `json:"bilibili"`
}

// TreeMachinesResponse Tree Machines Response.
type TreeMachinesResponse struct {
	Code int      `json:"code"`
	Data []string `json:"data"`
}

// TreeSonResponse Tree Son Response.
type TreeSonResponse struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

// TreeRoleResponse Tree Role Response.
type TreeRoleResponse struct {
	Code int         `json:"code"`
	Data []*TreeRole `json:"data"`
}

// TreeInstancesResponse Tree Instance Response.
type TreeInstancesResponse struct {
	Code int                      `json:"code"`
	Data map[string]*TreeInstance `json:"data"`
}

// TreeInstance Tree Instance.
type TreeInstance struct {
	HostName     string `json:"hostname"`
	IP           string `json:"ip"`
	InstanceType string `json:"instance_type"`
	InternalIP   string `json:"internal_ip"`
	ServiceIP    string `json:"service_ip"`
	ExtendIP     string `json:"extend_ip"`
}

// TreeAppInstanceRequest Tree App Instance Request.
type TreeAppInstanceRequest struct {
	Paths []string `json:"paths"`
}

// TreeAppInstanceResponse Tree App Instance Response.
type TreeAppInstanceResponse struct {
	Code int                           `json:"code"`
	Data map[string][]*TreeAppInstance `json:"data"`
}

// TreeAppInstance Tree App Instance.
type TreeAppInstance struct {
	HostName string `json:"hostname"`
}

// TreePlatformTokenRequest Tree Platform Token Request.
type TreePlatformTokenRequest struct {
	UserName   string `json:"user_name"`
	PlatformID string `json:"platform_id"`
}

// TreeRole Tree Role.
type TreeRole struct {
	UserName string `json:"user_name"`
	Role     int    `json:"role"`
	OldRole  int    `json:"old_role"`
	RdSre    bool   `json:"rd_sre"`
}

// TreeConf tree conf.
type TreeConf struct {
	Host   string
	Key    string
	Secret string
}

// TreeInstanceRequest request for hostname.
type TreeInstanceRequest struct {
	Path          string `query:"path"`
	PathFuzzy     string `query:"path_fuzzy"`
	Hostname      string `query:"hostname"`
	HostnameFuzzy string `query:"hostname_fuzzy"`
	HostnameRegex string `query:"hostname_regex"`
}

// ToQueryURI convert field to uri.
func (tir TreeInstanceRequest) ToQueryURI() string {
	var (
		params = &url.Values{}
		t      = reflect.TypeOf(tir)
		v      = reflect.ValueOf(tir)
		fv     string
	)
	for i := 0; i < t.NumField(); i++ {
		fv = v.Field(i).Interface().(string)
		if fv != "" {
			params.Set(t.Field(i).Tag.Get(_query), fv)
		}
	}
	return params.Encode()
}
