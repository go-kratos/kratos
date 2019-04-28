package dispatch

import (
	"fmt"
	"testing"
)

func TestMatcher(t *testing.T) {
	config := []byte(`{
        "ip_max_limit": 2,
        "default_domain" : "broadcastlv.chat.bilibili.com",
        "danmaku_common_dispatch": {
            "china" :{
                "china_telecom": {
                    "master": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24
                    },
                    "slave": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24,
                        "aliyun": 30
                    }
                },
                "china_unicom": {
                    "master": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24
                    },
                    "slave": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24,
                        "aliyun": 30
                    }
                },
                "cmcc": {
                    "master": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24
                    },
                    "slave": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24,
                        "aliyun": 30
                    }
                },
                "other": {
                    "master": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24
                    },
                    "slave": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24,
                        "aliyun": 30
                    }
                }
            },
            "oversea": [
                {
                    "rule":"($lng >= -20) && ($lng <= 160)",
                    "master": {
                        "tencent_siliconvalley": 10
                    },
                    "slave": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24,
                        "aliyun": 30
                    }
                },
                {
                    "master": {
                        "tencent_siliconvalley": 10
                    },
                    "slave": {
                        "tencent_shanghai": 10,
                        "tencent_guangzhou": 6,
                        "kingsoft": 24,
                        "aliyun": 30
                    }
                }
            ],
            "unknown" : {
                "master": {
                    "tencent_shanghai": 10,
                    "tencent_guangzhou": 6,
                    "kingsoft": 24
                },
                "slave": {
                    "tencent_shanghai": 10,
                    "tencent_guangzhou": 6,
                    "kingsoft": 24,
                    "aliyun": 30
                }
            }
        },
        "danmaku_vip_dispatch" : [
                {
                    "rule":"$uid==120497668",
                    "ip": ["118.89.14.174"]
                },
                {
                    "rule":"$uid % 10 == 1",
                    "group": ["tencent_guangzhou"]
                },
                {
                    "rule":"$uid == 221122111"
                }
        ],
        "danmaku_comet_group": {
            "tencent_shanghai": [
                "118.89.14.174",
                "118.89.14.115",
                "118.89.14.103",
                "118.89.14.206",
                "118.89.13.229"
            ],
            "tencent_guangzhou": [
                "211.159.194.41",
                "211.159.194.115",
                "211.159.194.105"
            ],
            "tencent_hongkong": [
                "119.28.56.183"
            ],
            "tencent_siliconvalley": [
                "49.51.37.200"
            ],
            "kingsoft": [
                "120.92.78.57",
                "120.92.158.137",
                "120.92.112.150"
            ],
            "aliyun": [
                "101.132.195.89",
                "47.104.64.120",
                "59.110.167.237",
                "47.92.112.162",
                "47.96.139.69",
                "119.23.41.85"
            ]
        }
    }`)
	m, err := NewMatcher(config, nil, nil, nil)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	fmt.Println(m)
}
