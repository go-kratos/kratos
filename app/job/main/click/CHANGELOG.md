# click-job

### v1.8.5
> 1. 冷门稿件集中落地的时间由2:00-5:00改为2:00-6:30，保证可以全部处理完

### v1.8.4
> 1. 修复rtype=2的问题

### v1.8.3
> 1. 将inline_play_heartbeat改为inline_play_to_view，且暂时兼容老逻辑

### v1.8.2
> 1. inline播放第一次上报由判断plat为6/7/8/9改为判断UA或者plat（兼容老逻辑），并且上报rtype =2
> 2. 增加played_time_enough逻辑（inline播放十秒计数），inline播放第二次上报由判断UA中的auto_play改为判断UA中的inline_play_heartbeat/played_time_enough/auto_play（兼容老逻辑）

### v1.8.1
> 1.增加拜年祭日志

### v1.8.0
> 1.2019年拜年祭的单品播放数加进主视频的播放数

### v1.7.7
> 1.lv=-2时，未登录用户也不计算播放数
> 2.支持黑名单列表

### v1.7.6
> 1.多线程一起pub

### v1.7.5
> 1.记录自动播放的buvidToDid的信息，并在之后的上报信息中根据它计算vv

### v1.7.4
> 1.消费report-click merge后的databus

### v1.7.3
> 1.自动播放判断

### v1.7.2
> 1.增加autoplay不计数的逻辑

### v1.7.1
> 1.迁移BM

### v1.7.0
> 1. update infoc sdk

### v1.6.1
> 1.支持plat=5  

### v1.6.0
> 1.迁移到主站目录下  

### v1.5.0
> 1.删除消费kafka的所有代码  

### v1.4.1
> 1.增加trace  

### v1.4.0
> 1.消费从kafka迁移到databus  

### v1.3.0
> 1.plat:5增加tv版  
> 2.forbid增加-2，全面封禁  
> 3.ugc&pgc计算方式改版 https://www.tapd.cn/20095661/prong/stories/view/1120095661001049183 

### v1.2.0
> 1.PGC根据epid来计算点击数，提高播放数  

### v1.1.3
> 1.增加释放内存的时间配置  

### v1.1.2
> 1.更新计数的时间配置化  

### v1.1.1
> 1.补充单元测试  

### v1.1.0
> 1.双写新的databus  

### v1.0.3
> 1.在凌晨2~4点之间释放内存中的数据  

### v1.0.2
> 1.处理maxAID逻辑  

### v1.0.1
> 1.简化程序日志  

### v1.0.0
> 1.项目初始化  
