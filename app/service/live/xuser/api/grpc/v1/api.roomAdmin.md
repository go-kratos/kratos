## 根据登录态获取功能入口是否显示, 需要登录态

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/is_any`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否有房管
        "has_admin": 0
    }
}
```

## 获取用户拥有的的所有房管身份

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/get_by_uid`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 用户uid|
|page|否|integer| 页数|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "page": {
            //  当前页码
            "page": 0,
            //  每页大小
            "page_size": 0,
            //  总页数
            "total_page": 0,
            //  总记录数
            "total_count": 0
        },
        "data": [
            {
                //  用户id
                "uid": 0,
                //  房间号
                "roomid": 0,
                //  主播的用户id
                "anchor_id": 0,
                //  主播用户名
                "uname": "",
                //  主播封面
                "anchor_cover": "",
                //  上任时间
                "ctime": ""
            }
        ]
    }
}
```

## 辞职房管

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/resign`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer| 房间号|
|uid|是|integer| 用户uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

## 查询需要添加的房管

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/search_for_admin`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||
|key_word|是|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": [
            {
                //  用户id
                "uid": 0,
                //  是否房管
                "is_admin": 0,
                //  用户名
                "uname": "",
                //  用户头像
                "face": "",
                //  粉丝勋章名称
                "medal_name": "",
                //  粉丝勋章等级
                "level": 0
            }
        ]
    }
}
```

## 获取主播拥有的的所有房管

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/get_by_anchor`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer| 页数|
|uid|是|integer| 用户uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "page": {
            //  当前页码
            "page": 0,
            //  每页大小
            "page_size": 0,
            //  总页数
            "total_page": 0,
            //  总记录数
            "total_count": 0
        },
        "data": [
            {
                //  用户id
                "uid": 0,
                //  用户名
                "uname": "",
                //  用户头像
                "face": "",
                //  上任时间
                "ctime": "",
                //  粉丝勋章名称
                "medal_name": "",
                //  粉丝勋章等级
                "level": 0,
                //  房间号
                "roomid": 0
            }
        ]
    }
}
```

## 获取主播拥有的的所有房管,房间号维度

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/get_by_room`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer| 房间号|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": [
            {
                //  上任时间
                "ctime": "",
                //  房管的用户id
                "uid": 0,
                //  房间号
                "roomid": 0
            }
        ]
    }
}
```

## 撤销房管

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/dismiss`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 房管的用户uid|
|anchor_id|是|integer| 主播uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

## 任命房管

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/appoint`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 房管的uid|
|anchor_id|是|integer| 主播uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  banner
        "userinfo": {
            //  用户id
            "uid": 0,
            //  用户名
            "uname": ""
        },
        //  房管的用户id
        "uid": 0,
        //  房间号
        "roomid": 0
    }
}
```

## 是否房管

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/is_admin`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 房管的uid|
|anchor_id|是|integer| 主播uid|
|roomid|是|integer| 房间号|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  banner
        "userinfo": {
            //  用户id
            "uid": 0,
            //  用户名
            "uname": ""
        },
        //  房管的用户id
        "uid": 0,
        //  房间号
        "roomid": 0
    }
}
```

## 是否房管, 不额外返回用户信息, 不判断是否主播自己

`GET http://api.live.bilibili.com/xlive/xuser/v1/roomAdmin/is_admin_short`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 房管的uid|
|roomid|是|integer| 房间号|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否房管 0:不是,1:是
        "result": 0
    }
}
```

