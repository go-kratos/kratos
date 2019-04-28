<!-- package=live.livedemo -->
- [/xlive/live-demo/foo2/hello](#xlivelive-demofoo2hello) 

## /xlive/live-demo/foo2/hello
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid|

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

