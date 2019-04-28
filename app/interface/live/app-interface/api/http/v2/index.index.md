## 首页大接口

`GET http://api.live.bilibili.com/xlive/app-interface/v2/index/getAllList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string| 平台|
|device|是|string| 设备|
|scale|是|string| 分辨率|
|build|是|integer| 版本号|
|relation_page|是|integer|关注页码|
|rec_page|否|integer|推荐页码 当前推荐页（用于天马强推），不传默认按1处理|
|quality|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        // 刷新重新请求间隔
        "interval": 0,
        // 是否命中天马灰度
        "is_sky_horse_gray": 0,
        // banner类型
        "banner": [
            {
                //  模块信息
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                //  模块列表
                "list": [
                    {
                        // 唯一标识id
                        "id": 0,
                        // 跳转url
                        "link": "",
                        // 图片url
                        "pic": "",
                        // 标题
                        "title": "",
                        // 内容
                        "content": ""
                    }
                ]
            }
        ],
        // 常用标签类型
        "my_tag": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                "extra_info": {
                    // 是否命中常用标签灰度
                    "is_gray": 0,
                    // offline已下线标签
                    "offline": [
                        {
                            "id": 0,
                            "area_v2_name": ""
                        }
                    ]
                },
                "list": [
                    {
                        "area_v2_id": 0,
                        "area_v2_name": "",
                        "area_v2_parent_id": 0,
                        "area_v2_parent_name": "",
                        "pic": "",
                        "link": "",
                        "is_advice": 0
                    }
                ]
            }
        ],
        // 分区入口类型
        "area_entrance": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                "list": [
                    {
                        // 唯一标识id
                        "id": 0,
                        // 跳转url
                        "link": "",
                        // 图片url
                        "pic": "",
                        // 标题
                        "title": "",
                        // 内容
                        "content": ""
                    }
                ]
            }
        ],
        // 大航海提示类型
        "sea_patrol": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                "extra_info": {
                    // 唯一标识id
                    "id": 0,
                    // 跳转url
                    "link": "",
                    // 图片url
                    "pic": "",
                    // 标题
                    "title": "",
                    // 内容
                    "content": ""
                }
            }
        ],
        // 我的关注类型
        "my_idol": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                "extra_info": {
                    "total_count": 0,
                    "time_desc": "",
                    "uname_desc": "",
                    "tags_desc": "",
                    "card_type": 0,
                    "relation_page": 0
                },
                "list": [
                    {
                        "roomid": 0,
                        "uid": 0,
                        "uname": "",
                        "face": "",
                        "cover": "",
                        "title": "",
                        "area": 0,
                        "live_time": 0,
                        "area_name": "",
                        "area_v2_id": 0,
                        "area_v2_name": "",
                        "area_v2_parent_name": "",
                        "area_v2_parent_id": 0,
                        "live_tag_name": "",
                        "online": 0,
                        "play_url": "",
                        "play_url_h265": "",
                        "accept_quality": [
                            0
                        ],
                        "current_quality": 0,
                        "pk_id": 0,
                        "link": "",
                        "special_attention": 0,
                        "broadcast_type": 0,
                        "pendent_ru": "",
                        "pendent_ru_color": "",
                        "pendent_ru_pic": "",
                        "official_verify": 0
                    }
                ]
            }
        ],
        // 通用房间列表类型
        "room_list": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                "list": [
                    {
                        // 当前拥有清晰度列表
                        "accept_quality": [
                            0
                        ],
                        // 二级分区id
                        "area_v2_id": 0,
                        // 一级分区id
                        "area_v2_parent_id": 0,
                        // 二级分区名称
                        "area_v2_name": "",
                        // 一级分区名称
                        "area_v2_parent_name": "",
                        // 横竖屏  0:横屏 1:竖屏 -1:异常情况
                        "broadcast_type": 0,
                        // 封面，封面现在有3种：关键帧、封面图、秀场封面（正方形的），返回哪个由后端决定
                        "cover": "",
                        // 当前清晰度,清晰度((0)) 0:默认码率, 2:800 3:1500 4:原画
                        "current_quality": 0,
                        // 主播头像
                        "face": "",
                        // 跳转链接
                        "link": "",
                        // 人气值
                        "online": 0,
                        // 新版角标-右上 默认为空 只能是文字！！！@古月 【5.29显示更新】：服务端还是吐右上（兼容老版），5.29显示在左上
                        "pendent_ru": "",
                        // 【5.29显示更新】：服务端还是吐右上，5.29客户端显示在左上,对应的背景图片
                        "pendent_ru_color": "",
                        // 新版移动端角标色值-右上
                        "pendent_ru_pic": "",
                        // pk_id
                        "pk_id": 0,
                        // 秒开播放串 h264
                        "play_url": "",
                        // 推荐类型 1：人气 2：营收 3：运营强推 4：天马推荐（暂定）用于客户端打点
                        "rec_type": 0,
                        // 房间id
                        "roomid": 0,
                        // 房间标题
                        "title": "",
                        // 主播uname
                        "uname": "",
                        // 秒开播放串 h265
                        "play_url_h265": ""
                    }
                ]
            }
        ],
        // 小时榜类型
        "hour_rank": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                "extra_info": {
                    // 14:00-15:00榜单
                    "sub_title": ""
                },
                "list": [
                    {
                        // 排名
                        "rank": 0,
                        // 用户uid
                        "uid": 0,
                        // 房间id
                        "roomid": 0,
                        // 用户昵称
                        "uname": "",
                        // 用户头像
                        "face": "",
                        // 直播状态
                        "live_status": 0,
                        // 父分区id
                        "area_v2_parent_id": 0,
                        // 父分区名称
                        "area_v2_parent_name": "",
                        // 子分区id
                        "area_v2_id": 0,
                        // 子分区名称
                        "area_v2_name": ""
                    }
                ]
            }
        ],
        // 活动卡
        "activity_card": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  模块跳转链接
                    "link": "",
                    //  模块图标
                    "pic": "",
                    //  模块标题
                    "title": "",
                    //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  模块排序值
                    "sort": 0,
                    //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
                    "count": 0
                },
                "list": [
                    {
                        "card": {
                            // 活动id
                            "aid": 0,
                            // 活动图片
                            "pic": "",
                            // 活动标题
                            "title": "",
                            // 活动文案
                            "text": "",
                            // 图片链接
                            "pic_link": "",
                            // 围观链接
                            "go_link": "",
                            // 三种：去围观，预约，已预约
                            "button_text": "",
                            // 代表卡片所处于的状态 0可以去围观,1用户可以点击去预约,2用户可以点击取消预约
                            "status": 0,
                            // card,room和av排序值
                            "sort": 0
                        },
                        "room": [
                            {
                                // 是否开播
                                "is_live": 0,
                                // 房间id
                                "room_id": 0,
                                // 房间标题
                                "title": "",
                                // 主播名
                                "u_name": "",
                                // 人气值
                                "online": 0,
                                // 封面
                                "cover": "",
                                // 父一级分区id
                                "area_v2_parent_id": 0,
                                // 二级分区id
                                "area_v2_id": 0,
                                // card,room和av排序值
                                "sort": 0
                            }
                        ],
                        "av": [
                            {
                                // 视频
                                "avid": 0,
                                // avid
                                "title": "",
                                // 视频标题
                                "view_count": 0,
                                // 浏览
                                "dan_ma_ku": 0,
                                // 弹幕
                                "duration": 0,
                                // 时长
                                "cover": "",
                                // card,room和av排序值
                                "sort": 0
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
```

## 换一换接口

`GET http://api.live.bilibili.com/xlive/app-interface/v2/index/change`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|module_id|是|integer|模块id|
|attention_room_id|否|string|关注的room ids|
|page|是|integer|换一换的当前页数|
|platform|是|string||
|build|是|integer||
|device|是|string||
|quality|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "module_info": {
            //  模块id
            "id": 0,
            //  模块跳转链接
            "link": "",
            //  模块图标
            "pic": "",
            //  模块标题
            "title": "",
            //  模块类型 1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
            "type": 0,
            //  模块排序值
            "sort": 0,
            //  模块数据源数量，按需、目前只有推荐有，其它模块都是默认值0
            "count": 0
        },
        "is_sky_horse_gray": 0,
        "list": [
            {
                // 当前拥有清晰度列表
                "accept_quality": [
                    0
                ],
                // 二级分区id
                "area_v2_id": 0,
                // 一级分区id
                "area_v2_parent_id": 0,
                // 二级分区名称
                "area_v2_name": "",
                // 一级分区名称
                "area_v2_parent_name": "",
                // 横竖屏  0:横屏 1:竖屏 -1:异常情况
                "broadcast_type": 0,
                // 封面，封面现在有3种：关键帧、封面图、秀场封面（正方形的），返回哪个由后端决定
                "cover": "",
                // 当前清晰度,清晰度((0)) 0:默认码率, 2:800 3:1500 4:原画
                "current_quality": 0,
                // 主播头像
                "face": "",
                // 跳转链接
                "link": "",
                // 人气值
                "online": 0,
                // 新版角标-右上 默认为空 只能是文字！！！@古月 【5.29显示更新】：服务端还是吐右上（兼容老版），5.29显示在左上
                "pendent_ru": "",
                // 【5.29显示更新】：服务端还是吐右上，5.29客户端显示在左上,对应的背景图片
                "pendent_ru_color": "",
                // 新版移动端角标色值-右上
                "pendent_ru_pic": "",
                // pk_id
                "pk_id": 0,
                // 秒开播放串 h264
                "play_url": "",
                // 推荐类型 1：人气 2：营收 3：运营强推 4：天马推荐（暂定）用于客户端打点
                "rec_type": 0,
                // 房间id
                "roomid": 0,
                // 房间标题
                "title": "",
                // 主播uname
                "uname": "",
                // 秒开播放串 h265
                "play_url_h265": ""
            }
        ]
    }
}
```

