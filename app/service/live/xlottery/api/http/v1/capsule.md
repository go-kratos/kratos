##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/get_coin_list`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer||
|pageSize|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "total": 0,
        "totalPage": 0,
        "list": [
            {
                "id": 0,
                "title": "",
                "changeNum": 0,
                "startTime": "",
                "endTime": "",
                "status": 0,
                "giftType": 0,
                "giftConfig": "",
                "areaIds": [
                    {
                        "parentId": 0,
                        "isAll": 0,
                        "sonIds": [
                            0
                        ]
                    }
                ]
            }
        ]
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/update_coin_config`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|title|否|string||
|changeNum|否|integer||
|startTime|否|string||
|endTime|否|string||
|status|否|integer||
|giftType|否|integer||
|giftConfig|否|string||
|areaIds|否|多个unknown||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/update_coin_status`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|status|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/delete_coin`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/get_pool_config`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer||
|pageSize|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "total": 0,
        "totalPage": 0,
        "list": [
            {
                "id": 0,
                "coinTitle": 0,
                "title": "",
                "start_time": "",
                "endTime": "",
                "status": 0,
                "rule": ""
            }
        ]
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/update_pool_config`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|coinTitle|否|integer||
|title|否|string||
|startTime|否|string||
|endTime|否|string||
|rule|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/delete_pool`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/update_pool_status`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|status|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/get_pool_detail`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|poolId|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "ID": 0,
                "poolID": 0,
                "type": 0,
                "num": 0,
                "objectID": 0,
                "webUrl": "",
                "mobileUrl": "",
                "description": "",
                "jumpUrl": "",
                "proType": 0,
                "chance": "",
                "loop": 0,
                "limit": 0
            }
        ]
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/update_pool_detail`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|poolId|否|integer||
|type|否|integer||
|num|否|integer||
|objectId|否|integer||
|webUrl|否|string||
|mobileUrl|否|string||
|description|否|string||
|jumpUrl|否|string||
|proType|否|integer||
|chance|否|string||
|loop|否|integer||
|limit|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/lottery/v1/capsule/delete_pool_detail`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "retStatus": true
    }
}
```

