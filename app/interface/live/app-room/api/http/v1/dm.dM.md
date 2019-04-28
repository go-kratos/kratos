<!-- package=live.approom.v1 -->
- [/xlive/app-room/v1/dM/SendMsg](#xliveapp-roomv1dMSendMsg) 
- [/xlive/app-room/v1/dM/GetHistory](#xliveapp-roomv1dMGetHistory) 

## /xlive/app-room/v1/dM/SendMsg
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer||
|msg|是|string||
|rnd|是|string||
|fontsize|是|integer||
|mode|否|integer||
|color|是|integer||
|bubble|否|integer||
|build|否|integer||
|anti|否|string||
|platform|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/app-room/v1/dM/GetHistory
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "Room": [
            ""
        ],
        "Admin": [
            ""
        ]
    }
}
```

