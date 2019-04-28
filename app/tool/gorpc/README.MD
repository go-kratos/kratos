### gorpc 

#### 说明
1. 根据service 方法生成rpc client 以及rpc model层 代码
2. service方法格式暂限定为以下格式

```
func (s *Receiver) FuncName(c context.Contex,args ...interface{}) (err error) {

}

func (s *Receiver) FuncName(c context.Contex,args ...interface{}) (resp interface{},err error) {
	
}
```
3. args 参数应为基础类型，如 int，string,slice 等

#### 使用
1. 进入到对应目录下，执行gorpc 即可生成rpc client代码，生成的代码默认放在 project/rpc/client/ 目录下

2. 使用-model 参数生成rpc 参数 model ，生成的代码默认放在 project/model/ 目录下


#### exmaple 

```
func (s *Service) Example(c context.Contex,i int,b string)(res *Resp,err error) {
	return
}
```
以上代码片段将自动生成如下代码
```model
type ArgExample struct {
	I int
	B string 
}
```

``` client
func (s *Service) Example(c context.Context, arg *model.ArgExample) (res *Resp, err error) {
	res = new(Resp)
	err = s.client.Call(c, _jury, arg, res)
	return
}
```