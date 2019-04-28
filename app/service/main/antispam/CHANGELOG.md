#### antispa

##### Version 1.3.0
> 1.unit test

##### Version 1.2.6
> 1. grpc identify

##### Version 1.2.5
> 1. print ping err

##### Version 1.2.4
> 1. 迁移至BM框架

##### Version 1.2.3
> 1. 增加register接口

##### Version 1.2.2
> 1. redis migrate folder

##### Version 1.2.1
> 1. err judgement fix panic

##### Version 1.2.0
> 1. move in main path

##### Version 1.1.16
> 1. delete statsd

##### Version 1.1.15
> 1. limit max regexps a area can have (in conf file)

##### Version 1.1.14
> 1. optimize code

##### Version 1.1.13
> 1. temporary fix nil keyword pointer panic

##### Version 1.1.12
> 1. return deleted regexps where get all regexps

##### Version 1.1.11
> 1. fix danmu share reply's regexps

##### Version 1.1.10
> 1. return precious keyword hit counts when incr
> 2. make auto white strategy configurable.

##### Version 1.1.9
> 1. fix 'fetch rules without getting area'

##### Version 1.1.8
> 1. fix autoWhite deviation bug

##### Version 1.1.6
> 1. add turning keyword into white automatically strategy

##### Version 1.1.5
> 1. danmu has its own regexps list 

##### Version 1.1.4
> 1. add danmu

##### Version 1.1.3
> 1. add unit test

##### Version 1.1.2
> 1. change `Id` to `ID`
> 2. remove unused columun keywords.senderId and keywords.regexp_content
> 3. add some comments on exported method/function
> 4. add new area "main_site_dm"

##### Version 1.1.1
> 1. fix missed "return err" bug
> 2. remove unused configuration "RefreshTrieDBSizePerQuery"

##### Version 1.1.0
> 1. use 'ctime BETWEEN ... AND XXX' instead of 'ctime < XXX'

##### Version 1.0.9
> 1. asynchronous incr count & persisit senderId in Filter.Check
> 2. add ruleDefaultExpireSec and regexpDefaultExpireSec configure options

##### Version 1.0.8
> 1. change recycle keyword sql

##### Version 1.0.7
> 1. limit durationSec and allowedCount max value

##### Version 1.0.6
> 1. refresh trie instead of building one frequently
> 2. add area `live_dm`

##### Version 1.0.5
> 1. record sender_id only if sender_id > 0
> 2. change logic to expire total_count cache
> 3. add err log when fail to update cache

##### Version 1.0.4
> 1. fix rate_limit_rule sql bug

##### Version 1.0.3
> 1. change sql to avoid slow query

##### Version 1.0.2
> 1. add log when ping error

##### Version 1.0.1
> 1. 修复ping error

##### Version 1.0.0
> 1. init commit  
> 2. fix ZRANGEBYSCORE params  
