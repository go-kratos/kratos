<!-- package=live.xuser.v1 -->
- [/live.xuser.v1.Guard/Buy](#live.xuser.v1.GuardBuy)  Buy 购买大航海
- [/live.xuser.v1.Guard/GetByUIDTargetID](#live.xuser.v1.GuardGetByUIDTargetID)  GetByUIDTargetID 获取我与目标用户守护关系,不支持批量(P0级)
- [/live.xuser.v1.Guard/GetByTargetIdsBatch](#live.xuser.v1.GuardGetByTargetIdsBatch)  GetByTargetIdsBatch 获取我与目标用户守护关系,支持批量(P2级,必要时刻降级)
- [/live.xuser.v1.Guard/GetByUIDTargetIds](#live.xuser.v1.GuardGetByUIDTargetIds)  GetByUIDTargetIds 根据uids批量获取所有守护关系,粉丝勋章使用
- [/live.xuser.v1.Guard/GetByUIDForGift](#live.xuser.v1.GuardGetByUIDForGift)  GetByUID 获取我所有的守护,不支持批量(P0级)
- [/live.xuser.v1.Guard/GetByUIDBatch](#live.xuser.v1.GuardGetByUIDBatch)  GetByUIDBatch 根据uids获取所有的守护,支持批量(P2级)
- [/live.xuser.v1.Guard/GetAnchorRecentTopGuard](#live.xuser.v1.GuardGetAnchorRecentTopGuard)  GetAnchorRecentTopGuard 获取最近的提督弹窗提醒
- [/live.xuser.v1.Guard/GetTopListGuard](#live.xuser.v1.GuardGetTopListGuard)  GetTopListGuard 获取某个up主的守护排行榜
- [/live.xuser.v1.Guard/GetTopListGuardNum](#live.xuser.v1.GuardGetTopListGuardNum)  GetTopListGuardNum 获取某个up主所有的守护数量,和GetTopListGuard接口的区别是此接口用于房间页首屏,逻辑比较简单,因此拆分开来
- [/live.xuser.v1.Guard/ClearUIDCache](#live.xuser.v1.GuardClearUIDCache)  ClearUIDCache 清除cache

## /live.xuser.v1.Guard/Buy
### Buy 购买大航海

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|order_id|是|string||
|uid|是|integer||
|ruid|是|integer||
|guard_level|是|integer||
|num|是|integer||
|platform|是|integer||
|source|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": 0
    }
}
```


## /live.xuser.v1.Guard/GetByUIDTargetID
### GetByUIDTargetID 获取我与目标用户守护关系,不支持批量(P0级)

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||
|target_id|是|integer||
|sort_type|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": {
            "1": {
                //  主键
                "id": 0,
                //  uid
                "uid": 0,
                //  target_id
                "target_id": 0,
                //  守护类型 1为总督，2为提督，3为舰长
                "privilege_type": 0,
                //  start_time
                "start_time": "",
                //  expired_time
                "expired_time": "",
                //  ctime
                "ctime": "",
                //  utime
                "utime": ""
            }
        }
    }
}
```


## /live.xuser.v1.Guard/GetByTargetIdsBatch
### GetByTargetIdsBatch 获取我与目标用户守护关系,支持批量(P2级,必要时刻降级)

#### 方法：GET

#### 请求参数

```javascript
{
    "targetIDs": [
        {
            "target_id": 0,
            "sort_type": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.xuser.v1.Guard/GetByUIDTargetIds
### GetByUIDTargetIds 根据uids批量获取所有守护关系,粉丝勋章使用

#### 方法：GET

#### 请求参数

```javascript
{
    "uid": 0,
    "targetIDs": [
        {
            "target_id": 0,
            "sort_type": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": {
            "1": {
                //  主键
                "id": 0,
                //  uid
                "uid": 0,
                //  target_id
                "target_id": 0,
                //  守护类型 1为总督，2为提督，3为舰长
                "privilege_type": 0,
                //  start_time
                "start_time": "",
                //  expired_time
                "expired_time": "",
                //  ctime
                "ctime": "",
                //  utime
                "utime": ""
            }
        }
    }
}
```


## /live.xuser.v1.Guard/GetByUIDForGift
### GetByUID 获取我所有的守护,不支持批量(P0级)

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": {
            "1": {
                //  主键
                "id": 0,
                //  uid
                "uid": 0,
                //  target_id
                "target_id": 0,
                //  守护类型 1为总督，2为提督，3为舰长
                "privilege_type": 0,
                //  start_time
                "start_time": "",
                //  expired_time
                "expired_time": "",
                //  ctime
                "ctime": "",
                //  utime
                "utime": ""
            }
        }
    }
}
```


## /live.xuser.v1.Guard/GetByUIDBatch
### GetByUIDBatch 根据uids获取所有的守护,支持批量(P2级)

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uids|是|多个integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": {
            "1": {
                "list": [
                    {
                        //  主键
                        "id": 0,
                        //  uid
                        "uid": 0,
                        //  target_id
                        "target_id": 0,
                        //  守护类型 1为总督，2为提督，3为舰长
                        "privilege_type": 0,
                        //  start_time
                        "start_time": "",
                        //  expired_time
                        "expired_time": "",
                        //  ctime
                        "ctime": "",
                        //  utime
                        "utime": ""
                    }
                ]
            }
        }
    }
}
```


## /live.xuser.v1.Guard/GetAnchorRecentTopGuard
### GetAnchorRecentTopGuard 获取最近的提督弹窗提醒

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  主键
        "cnt": 0,
        "list": [
            {
                "uid": 0,
                "end_time": 0,
                "is_open": 0
            }
        ]
    }
}
```


## /live.xuser.v1.Guard/GetTopListGuard
### GetTopListGuard 获取某个up主的守护排行榜

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||
|page|否|integer||
|page_size|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  守护总数量
        "num": 0,
        "page": 0,
        "now": 0,
        "list": [
            {
                "uid": 0,
                "ruid": 0,
                "rank": 0,
                "guard_level": 0
            }
        ],
        "top3": [
            {
                "uid": 0,
                "ruid": 0,
                "rank": 0,
                "guard_level": 0
            }
        ]
    }
}
```


## /live.xuser.v1.Guard/GetTopListGuardNum
### GetTopListGuardNum 获取某个up主所有的守护数量,和GetTopListGuard接口的区别是此接口用于房间页首屏,逻辑比较简单,因此拆分开来

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "total_count": 0
    }
}
```


## /live.xuser.v1.Guard/ClearUIDCache
### ClearUIDCache 清除cache

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||
|magic_key|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

