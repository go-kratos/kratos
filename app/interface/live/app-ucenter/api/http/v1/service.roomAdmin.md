<!-- package=live.appucenter.v1 -->
- [/xlive/app-ucenter/v1/roomAdmin/is_any](#xliveapp-ucenterv1roomAdminis_any)  根据登录态获取功能入口是否显示, 需要登录态
- [/xlive/app-ucenter/v1/roomAdmin/get_by_uid](#xliveapp-ucenterv1roomAdminget_by_uid)  获取用户拥有的的所有房管身份
- [/xlive/app-ucenter/v1/roomAdmin/resign](#xliveapp-ucenterv1roomAdminresign)  辞职房管
- [/xlive/app-ucenter/v1/roomAdmin/search_for_admin](#xliveapp-ucenterv1roomAdminsearch_for_admin)  查询需要添加的房管
- [/xlive/app-ucenter/v1/roomAdmin/get_by_anchor](#xliveapp-ucenterv1roomAdminget_by_anchor)  获取主播拥有的的所有房管
- [/xlive/app-ucenter/v1/roomAdmin/dismiss](#xliveapp-ucenterv1roomAdmindismiss)  撤销房管
- [/xlive/app-ucenter/v1/roomAdmin/appoint](#xliveapp-ucenterv1roomAdminappoint)  任命房管

## /xlive/app-ucenter/v1/roomAdmin/is_any
### 根据登录态获取功能入口是否显示, 需要登录态

> 需要登录

#### 方法：GET

#### 请求参数


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否有房管
        "has_admin": 0
    }
}
```


## /xlive/app-ucenter/v1/roomAdmin/get_by_uid
### 获取用户拥有的的所有房管身份

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer| 页数|

#### 响应

```javascript
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


## /xlive/app-ucenter/v1/roomAdmin/resign
### 辞职房管

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer| 房间号|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/app-ucenter/v1/roomAdmin/search_for_admin
### 查询需要添加的房管

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|key_word|是|string||

#### 响应

```javascript
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


## /xlive/app-ucenter/v1/roomAdmin/get_by_anchor
### 获取主播拥有的的所有房管

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer| 页数|

#### 响应

```javascript
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
                "level": 0
            }
        ]
    }
}
```


## /xlive/app-ucenter/v1/roomAdmin/dismiss
### 撤销房管

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 房管的用户uid|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/app-ucenter/v1/roomAdmin/appoint
### 任命房管

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 房管的uid|

#### 响应

```javascript
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
        "roomid": 0,
        //  创建时间　"2017-07-26 17:12:51"
        "ctime": ""
    }
}
```

