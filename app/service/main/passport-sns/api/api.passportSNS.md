<!-- package=passport.service.sns -->
- [/x/internal/passport-sns/authorize/url](#xinternalpassport-snsauthorizeurl)  GetAuthorizeURL get authorize url
- [/x/internal/passport-sns/bind](#xinternalpassport-snsbind)  Bind bind sns account
- [/x/internal/passport-sns/unbind](#xinternalpassport-snsunbind)  Unbind unbind sns account
- [/x/internal/passport-sns/info](#xinternalpassport-snsinfo)  GetInfo get info by mid
- [/x/internal/passport-sns/info/code](#xinternalpassport-snsinfocode)  GetInfoByCode get info by authorize code
- [/x/internal/passport-sns/info/update](#xinternalpassport-snsinfoupdate)  UpdateInfo update info

## /x/internal/passport-sns/authorize/url
### GetAuthorizeURL get authorize url

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|app_id|是|string||
|platform|是|string||
|redirect_url|是|string||
|display|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "url": ""
    }
}
```


## /x/internal/passport-sns/bind
### Bind bind sns account

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|mid|是|integer||
|code|是|string||
|app_id|是|string||
|platform|是|string||
|redirect_url|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /x/internal/passport-sns/unbind
### Unbind unbind sns account

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|mid|是|integer||
|app_id|否|string||
|platform|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /x/internal/passport-sns/info
### GetInfo get info by mid

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|mid|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "infos": [
            {
                "mid": 0,
                "platform": "",
                "unionid": "",
                "expires": 0
            }
        ]
    }
}
```


## /x/internal/passport-sns/info/code
### GetInfoByCode get info by authorize code

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|code|是|string||
|app_id|是|string||
|platform|是|string||
|redirect_url|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "mid": 0,
        "unionid": "",
        "openid": "",
        "expires": 0,
        "token": ""
    }
}
```


## /x/internal/passport-sns/info/update
### UpdateInfo update info

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|app_id|是|string||
|mid|是|integer||
|open_id|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

