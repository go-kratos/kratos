<!-- package=live.webucenter -->
- [/xlive/web-ucenter/user/get_user_info](#xliveweb-ucenteruserget_user_info)  根据uid查询用户信息

## /xlive/web-ucenter/user/get_user_info
### 根据uid查询用户信息

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string| platform in url|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  用户uid
        "uid": 0,
        //  用户名
        "uname": "",
        //  头像
        "face": "",
        //  主站硬币
        "billCoin": 0.1,
        //  用户银瓜子
        "silver": 0,
        //  用户金瓜子
        "gold": 0,
        //  用户成就点
        "achieve": 0,
        //  月费姥爷
        "vip": 0,
        //  年费姥爷
        "svip": 0,
        //  用户等级
        "user_level": 0,
        //  用户下一等级
        "user_next_level": 0,
        //  用户在当前等级已经获得的经验
        "user_intimacy": 0,
        //  用户从当前等级升级到下一级所需总经验
        "user_next_intimacy": 0,
        //  新增字段，判断用户是否达到满级 0:没有1:满级
        "is_level_top": 0,
        //  用户等级排名
        "user_level_rank": "",
        //  年返逻辑，已失效
        "user_charged": 0
    }
}
```

