<!-- package=live.appblink.v1 -->
- [/xlive/app-blink/v1/room/GetInfo](#xliveapp-blinkv1roomGetInfo) 获取房间基本信息
- [/xlive/app-blink/v1/room/Create](#xliveapp-blinkv1roomCreate) 创建房间

## /xlive/app-blink/v1/room/GetInfo
###获取房间基本信息

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "room_id": 0,
        "uid": 0,
        "uname": "",
        "title": "",
        "face": "",
        "try_time": "",
        "live_status": 0,
        "area_v2_name": "",
        "area_v2_id": 0,
        "master_level": 0,
        "master_level_color": 0,
        "master_score": 0,
        "master_next_level": 0,
        "max_level": 0,
        "fc_num": 0,
        "rcost": 0,
        "medal_status": 0,
        "medal_name": "",
        "medal_rename_status": 0,
        "is_medal": 0,
        "full_text": "",
        "identify_status": 0,
        "lock_status": 0,
        "lock_time": "",
        "open_medal_level": 0,
        "master_next_level_score": 0,
        "parent_id": 0,
        "parent_name": ""
    }
}
```


## /xlive/app-blink/v1/room/Create
###创建房间

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "room_id": ""
    }
}
```

