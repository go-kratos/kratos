<!-- package=live.xroom.v1 -->
- [/live.xroom.v1.Room/getMultiple](#live.xroom.v1.RoomgetMultiple)  批量根据room_ids获取房间信息
- [/live.xroom.v1.Room/getMultipleByUids](#live.xroom.v1.RoomgetMultipleByUids)  批量根据uids获取房间信息
- [/live.xroom.v1.Room/isAnchor](#live.xroom.v1.RoomisAnchor)  批量根据uids判断是否是主播，如果是返回主播的room_id，否则返回0

## /live.xroom.v1.Room/getMultiple
### 批量根据room_ids获取房间信息

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_ids|是|多个integer| room_ids数组，长号|
|attrs|是|多个string| 要获取的房间信息维度 status:状态相关 show:展示相关 area:分区相关 anchor:主播相关|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  主播room_id => 房间维度信息
        "list": {
            "1": {
                //  room_id 房间长号
                "room_id": 0,
                //  uid 主播uid
                "uid": 0,
                //  Model1：房间信息（状态相关）
                "status": {
                    //  直播间状态 0未开播，1直播中；2轮播中；
                    "live_status": 0,
                    //  横竖屏方向 0横屏，1竖屏
                    "live_screen_type": 0,
                    //  是否开播过标识
                    "live_mark": 0,
                    //  封禁状态：0未封禁；1审核封禁; 2全网封禁
                    "lock_status": 0,
                    //  封禁时间戳
                    "lock_time": 0,
                    //  隐藏状态 0不隐藏，1隐藏
                    "hidden_status": 0,
                    //  隐藏时间戳
                    "hidden_time": 0,
                    //  直播类型 0默认 1摄像头直播 2录屏直播 3语音直播
                    "live_type": 0
                },
                //  Model2：房间信息（展示相关）
                "show": {
                    //  short_id 短号
                    "short_id": 0,
                    //  直播间标题
                    "title": "",
                    //  直播间封面
                    "cover": "",
                    //  直播间标签
                    "tags": "",
                    //  直播间背景图
                    "background": "",
                    //  直播间简介
                    "description": "",
                    //  关键帧
                    "keyframe": "",
                    //  人气值
                    "popularity_count": 0,
                    //  房间tag（角标）
                    "tag_list": [
                        {
                            "tag_id": 0,
                            "tag_sub_id": 0,
                            "tag_value": 0,
                            "tag_ext": ""
                        }
                    ],
                    //  最近一次开播时间戳
                    "live_start_time": 0
                },
                //  Model3：房间信息（分区相关）
                "area": {
                    //  直播间分区id
                    "area_id": 0,
                    //  直播间分区名称
                    "area_name": "",
                    //  直播间父分区id
                    "parent_area_id": 0,
                    //  直播间父分区名称
                    "parent_area_name": ""
                },
                //  Model4：房间信息（主播相关）
                "anchor": {
                    //  主播类型
                    "anchor_profile_type": 0,
                    //  主播等级
                    "anchor_level": {
                        //  等级
                        "level": 0,
                        //  当前等级颜色
                        "color": 0,
                        //  当前积分
                        "score": 0,
                        //  当前等级最小积分
                        "left": 0,
                        //  下一等级起始积分
                        "right": 0,
                        //  下一个经验值
                        "max_level": 0
                    }
                }
            }
        }
    }
}
```


## /live.xroom.v1.Room/getMultipleByUids
### 批量根据uids获取房间信息

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uids|是|多个integer| 主播uids|
|attrs|是|多个string| 要获取的房间信息维度 status:状态相关 show:展示相关 area:分区相关 anchor:主播相关|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  主播UID => 房间维度信息
        "list": {
            "1": {
                //  room_id 房间长号
                "room_id": 0,
                //  uid 主播uid
                "uid": 0,
                //  Model1：房间信息（状态相关）
                "status": {
                    //  直播间状态 0未开播，1直播中；2轮播中；
                    "live_status": 0,
                    //  横竖屏方向 0横屏，1竖屏
                    "live_screen_type": 0,
                    //  是否开播过标识
                    "live_mark": 0,
                    //  封禁状态：0未封禁；1审核封禁; 2全网封禁
                    "lock_status": 0,
                    //  封禁时间戳
                    "lock_time": 0,
                    //  隐藏状态 0不隐藏，1隐藏
                    "hidden_status": 0,
                    //  隐藏时间戳
                    "hidden_time": 0,
                    //  直播类型 0默认 1摄像头直播 2录屏直播 3语音直播
                    "live_type": 0
                },
                //  Model2：房间信息（展示相关）
                "show": {
                    //  short_id 短号
                    "short_id": 0,
                    //  直播间标题
                    "title": "",
                    //  直播间封面
                    "cover": "",
                    //  直播间标签
                    "tags": "",
                    //  直播间背景图
                    "background": "",
                    //  直播间简介
                    "description": "",
                    //  关键帧
                    "keyframe": "",
                    //  人气值
                    "popularity_count": 0,
                    //  房间tag（角标）
                    "tag_list": [
                        {
                            "tag_id": 0,
                            "tag_sub_id": 0,
                            "tag_value": 0,
                            "tag_ext": ""
                        }
                    ],
                    //  最近一次开播时间戳
                    "live_start_time": 0
                },
                //  Model3：房间信息（分区相关）
                "area": {
                    //  直播间分区id
                    "area_id": 0,
                    //  直播间分区名称
                    "area_name": "",
                    //  直播间父分区id
                    "parent_area_id": 0,
                    //  直播间父分区名称
                    "parent_area_name": ""
                },
                //  Model4：房间信息（主播相关）
                "anchor": {
                    //  主播类型
                    "anchor_profile_type": 0,
                    //  主播等级
                    "anchor_level": {
                        //  等级
                        "level": 0,
                        //  当前等级颜色
                        "color": 0,
                        //  当前积分
                        "score": 0,
                        //  当前等级最小积分
                        "left": 0,
                        //  下一等级起始积分
                        "right": 0,
                        //  下一个经验值
                        "max_level": 0
                    }
                }
            }
        }
    }
}
```


## /live.xroom.v1.Room/isAnchor
### 批量根据uids判断是否是主播，如果是返回主播的room_id，否则返回0

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uids|是|多个integer| 主播uids|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  uid => room_id(长号)，room_id=0表示没有创建房间
        "list": {
            "1": 0
        }
    }
}
```

