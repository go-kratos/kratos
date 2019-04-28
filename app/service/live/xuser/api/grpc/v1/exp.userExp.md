## GetUserExpMulti 获取用户经验与等级信息,支持批量

`GET http://api.live.bilibili.com/xlive/xuser/v1/userExp/GetUserExp`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uids|是|多个integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": {
            "1": {
                "uid": 0,
                "userLevel": {
                    //  当前用户等级
                    "level": 0,
                    //  下一等级
                    "nextLevel": 0,
                    //  当前等级对应的经验
                    "userExpLeft": 0,
                    //  下一等级对应的经验
                    "userExpRight": 0,
                    //  用户当前经验
                    "userExp": 0,
                    //  升级到下一等级对应的经验
                    "userExpNextLevel": 0,
                    //  当前等级颜色
                    "color": 0,
                    //  下一等级左侧对应的经验
                    "userExpNextLeft": 0,
                    //  下一等级右侧对应的经验
                    "userExpNextRight": 0,
                    "isLevelTop": 0
                },
                "anchorLevel": {
                    //  当前用户等级
                    "level": 0,
                    //  下一等级
                    "nextLevel": 0,
                    //  当前等级对应的经验
                    "userExpLeft": 0,
                    //  下一等级对应的经验
                    "userExpRight": 0,
                    //  用户当前经验
                    "userExp": 0,
                    //  升级到下一等级对应的经验
                    "userExpNextLevel": 0,
                    //  当前等级颜色
                    "color": 0,
                    //  下一等级左侧对应的经验
                    "userExpNextLeft": 0,
                    //  下一等级右侧对应的经验
                    "userExpNextRight": 0,
                    //  主播积分,userExp/100
                    "anchorScore": 0,
                    "isLevelTop": 0
                }
            }
        }
    }
}
```

## AddUserExp 增加用户经验,不支持批量

`GET http://api.live.bilibili.com/xlive/xuser/v1/userExp/AddUserExp`

### 请求参数

```json
{
    "userInfo": {
        "uid": 0,
        "req_biz": 0,
        "type": 0,
        "num": 0
    }
}
```

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

