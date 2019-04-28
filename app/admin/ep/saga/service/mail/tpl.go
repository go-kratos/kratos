package mail

var (
	mailTPL = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>saga page</title>
</head>
<body>
<form action="" method="post" class="basic-grey" style="margin-left:auto;
margin-right:auto;
max-width: 500px;
background: #F7F7F7;
padding: 25px 15px 25px 10px;
font: 12px Georgia, 'Times New Roman', Times, serif;
color: #888;
text-shadow: 1px 1px 1px #FFF;
border:1px solid #E4E4E4;">
<h1 style="font-size: 25px;
padding: 0px 0px 10px 40px;
display: block;
border-bottom:1px solid #E4E4E4;
margin: -10px -15px 30px -10px;;
color: #888;">
    Saga 
<span style="display: block;font-size: 11px;">
    Merge Request 事件通知</span>
</h1>
<label style="display: block;margin: 0px;">
<span style="float: left;
width:100%;
text-align: left;
padding-right: 10px;
padding-left: 30px;
margin-top: 10px;
color: #888;">
申请人 : {{.UserName}}
<br />
来源分支 : {{.SourceBranch}}
<br />
目标分支 : {{.TargetBranch}}
<br />
修改标题 : {{.Title}}
<br />
修改说明 : {{.Description}}
<br />
<a href="{{.URL}}">点击查看..</a>
<br />
<br />
<br />
<h3>额外信息: </h3>
{{.Info}}
<br />
<br />
<br />
</span>
</label>
</form>
</body>
</html>`
	mailTPL3 = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>saga page</title>
</head>
<body>
<form action="" method="post" class="basic-grey" style="
margin-right:auto;
max-width: 500px;
font: 12px Georgia, 'Times New Roman', Times, serif;
color: #888;
text-shadow: 1px 1px 1px #FFF;">
<h1 style="font-size: 25px;
padding: 10px 0px 10px 40px;
display: block;
border-bottom:1px solid #E4E4E4;
margin: -10px -15px 5px -10px;;
color: #03A9F4;">
    Saga 
<span style="display: block;font-size: 11px;">
	事件通知</span>
</h1>
</form>
<label style="display: block;margin: 0px;">
<span style="float: left;
width:100%;
text-align: left;
padding-right: 10px;
padding-left: 30px;
margin-top: 10px;
color: #888;">
执行状态 :
<font class="{{.PipelineStatus}}" >
 {{.PipeStatus}}
 </font>
<br />
Pipeline信息:
<font class="{{.PipelineStatus}}" >
<a href="{{.URL}}">{{.URL}}</a>
 </font>
<br />
来源分支 : {{.SourceBranch}}
<br />
修改说明 : {{.Description}}
<br />
额外信息:  {{.Info}}
<br />
<br />
<br />
</span>
</label>
</body>
<style type="text/css">
.failed {
	color: #f21303;
}
.failed a{
	color: #f21303;
}
.success {
	color: #1aaa55;
}
.success a{
	color: #1aaa55;
}
</style>
</html>`
)
