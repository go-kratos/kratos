<!-- package=live.livedemo.v1 -->
- [/xlive/demo/v1/foo/uname_by_uid_custom_route](#xlivedemov1foouname_by_uid_custom_route)  根据uid得到uname
- [/xlive/live-demo/v1/foo/get_info](#xlivelive-demov1fooget_info)  获取房间信息
- [/xlive/live-demo/v1/foo/uname_by_uid3](#xlivelive-demov1foouname_by_uid3)  根据uid得到uname v3
- [/xlive/internal/live-demo/v1/foo/uname_by_uid4](#xliveinternallive-demov1foouname_by_uid4)  test comment
- [/xlive/live-demo/v1/foo/get_dynamic](#xlivelive-demov1fooget_dynamic) 
- [/xlive/live-demo/v1/foo/nointerface](#xlivelive-demov1foonointerface) 

## /xlive/demo/v1/foo/uname_by_uid_custom_route
### 根据uid得到uname

 这是详细说明

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid aaa|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  用户名
        "uname": "hello",
        //  idshaha
        "ids": [
            343242
        ],
        "list": [
            {
                "hello": "\"withquote",
                "world": ""
            }
        ],
        "alist": {
            "hello": "\"withquote",
            "world": ""
        },
        "amap": {
            "mapKey": {
                "hello": "\"withquote",
                "world": ""
            }
        }
    }
}
```


## /xlive/live-demo/v1/foo/get_info
### 获取房间信息

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_id|是|integer| 房间id `mock:"123"|
|many_ids|否|多个integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  房间id 注释貌似只有放在前面才能被识别，放到字段声明后面是没用的
        "roomid": 0,
        //  用户名
        "uname": "",
        //  开播时间
        "live_time": "",
        "amap": {
            "1": ""
        },
        "rate": 6.02214129e23,
        //  用户mid
        "mid": 0
    }
}
```


## /xlive/live-demo/v1/foo/uname_by_uid3
### 根据uid得到uname v3

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid aaa|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  用户名
        "uname": "hello",
        //  idshaha
        "ids": [
            343242
        ],
        "list": [
            {
                "hello": "\"withquote",
                "world": ""
            }
        ],
        "alist": {
            "hello": "\"withquote",
            "world": ""
        },
        "amap": {
            "mapKey": {
                "hello": "\"withquote",
                "world": ""
            }
        }
    }
}
```


## /xlive/internal/live-demo/v1/foo/uname_by_uid4
### test comment

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid aaa|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  用户名
        "uname": "hello",
        //  idshaha
        "ids": [
            343242
        ],
        "list": [
            {
                "hello": "\"withquote",
                "world": ""
            }
        ],
        "alist": {
            "hello": "\"withquote",
            "world": ""
        },
        "amap": {
            "mapKey": {
                "hello": "\"withquote",
                "world": ""
            }
        }
    }
}
```


## /xlive/live-demo/v1/foo/get_dynamic
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid aaa|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  用户名
        "uname": "hello",
        //  idshaha
        "ids": [
            343242
        ],
        "list": [
            {
                "hello": "\"withquote",
                "world": ""
            }
        ],
        "alist": {
            "hello": "\"withquote",
            "world": ""
        },
        "amap": {
            "mapKey": {
                "hello": "\"withquote",
                "world": ""
            }
        }
    }
}
```


## /xlive/live-demo/v1/foo/nointerface
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid aaa|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  用户名
        "uname": "hello",
        //  idshaha
        "ids": [
            343242
        ],
        "list": [
            {
                "hello": "\"withquote",
                "world": ""
            }
        ],
        "alist": {
            "hello": "\"withquote",
            "world": ""
        },
        "amap": {
            "mapKey": {
                "hello": "\"withquote",
                "world": ""
            }
        }
    }
}
```

