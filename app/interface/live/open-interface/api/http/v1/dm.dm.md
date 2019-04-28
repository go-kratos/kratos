<!-- package=live.openinterface.v1 -->
- [/xlive/open-interface/v1/dm/sendmsg](#xliveopen-interfacev1dmsendmsg) 
- [/xlive/open-interface/v1/dm/getConf](#xliveopen-interfacev1dmgetConf) 

## /xlive/open-interface/v1/dm/sendmsg
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|Msg|是|string||
|Ts|是|string||
|RoomID|是|integer||
|Group|是|string||
|Sign|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/open-interface/v1/dm/getConf
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|Ts|是|string||
|Sign|是|string||
|Group|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "ws_port": [
            0
        ],
        "wss_port": [
            0
        ],
        "tcp_port": [
            0
        ],
        "ip_list": [
            ""
        ],
        "domain_list": [
            ""
        ]
    }
}
```

