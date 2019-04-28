#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @Date    : 2017-12-25 下午2:43
# @Author  : Hedan (hedan@bilibili.com)
# @file    : databusTest



import json
import redis


# [databus]
#     key = "4ba46ba31f9a44ef"
#     secret = "99985eb4451cfb1b899ca0fbe3c4bdc8"
#     group = "AccountLog-MainAccount-P"
#     topic = "AccountLog-T"
#     action = "pub"
#     name = "member-service/databus"
#     proto = "tcp"
#     addr = "172.16.33.158:6205"
#     idle = 100
#     active = 100
#     dialTimeout = "1s"
#     readTimeout = "60s"
#     writeTimeout = "1s"
#     idleTimeout = "10s"



"""参考文档  http://info.bilibili.co/pages/viewpage.action?pageId=3670491
key          服务配置中databus的key
value        服务配置中databus的value
host和port   配置服务中databus的host和port
group-topic  一个group对应一个topic(开发申请)
"""

#每日登录 passport 获得5经验值
# data_passport ={
# 	 'mid':4780461,
#      'loginip':1726481463,
#      "timestamp":1516517576,
#       }
# auth_passport = '4ba46ba31f9a44ef:99985eb4451cfb1b899ca0fbe3c4bdc8@PassportLog-MainAccount-P/topic=PassportLog-T&role=pub&offset=new'
# rc = redis.Redis(host='172.16.33.158', port=6205, socket_keepalive=True)
# print rc.execute_command('auth',auth_passport)

# msg = json.dumps(data_passport, ensure_ascii=False)
# rc.set("1023", msg)


# 每日登录 account 获得5经验值
# data_passport ={
# 	 'mid':110000092,
#      'loginip':1726481463,
#      "timestamp":1516517576,
#       }
# auth_passport = '4ba46ba31f9a44ef:99985eb4451cfb1b899ca0fbe3c4bdc8@AccountLoginAward-MainAccount-P/topic=AccountLoginAward-T&role=pub&offset=new'
# rc = redis.Redis(host='172.16.33.158', port=6205, socket_keepalive=True)
# print rc.execute_command('auth',auth_passport)

# msg = json.dumps(data_passport, ensure_ascii=False)
# rc.set("1023", msg)




# 每日首次分享视频 archive 获得5经验值
# data_archive ={
#      'event':"share",
# 	   'mid':110000092,
#      'ip':"127.0.0.1",
#      "ts":111,
#       }
# auth_archive = '4ba46ba31f9a44ef:99985eb4451cfb1b899ca0fbe3c4bdc8@AccountExp-MainAccount-P/topic=AccountExp-T&role=pub&offset=new'
# rc = redis.Redis(host='172.16.33.158', port=6205, socket_keepalive=True)
# print rc.execute_command('auth',auth_archive)
# msg = json.dumps(data_archive, ensure_ascii=False)
# rc.set("1023", msg)
#
#

# 每日首次分享视频 archive 获得5经验值
data_archive ={
     'event':"view",
	#  'mid':110000092,
      'mid':4780461,
     'ip':"127.0.0.1",
     "ts":1521745000,
      }
auth_archive = '4ba46ba31f9a44ef:99985eb4451cfb1b899ca0fbe3c4bdc8@AccountExp-MainAccount-P/topic=AccountExp-T&role=pub&offset=new'
rc = redis.Redis(host='172.16.33.158', port=6205, socket_keepalive=True)
print rc.execute_command('auth',auth_archive)
msg = json.dumps(data_archive, ensure_ascii=False)
rc.set("1023", msg)

# 每日首次观看 history 获得5经验值
# data_history ={
#     "action":"insert",
# 	"table":"aso_account",
# 	"data":{"mid":110000193,
# 	"userid":"test_dan06",
# 	"uname":"test_dan08",
# 	"pwd":"6cfbf96b8f0eb2e0a82b46a4236e8883",
# 	"salt":"D8fd30Kj",
# 	"email":"169d9106a74d5e95de71be6cf373af04",
# 	"tel":"218cb4bf8762354eae473b3b612f707e",
# 	"country_id":1,
# 	"mobile_verified":0,
# 	"isleak":0,
# 	"mtime":"2018-01-02 17:18:58"},
# 	"flag":0
#       }
# auth_history = '0QEO9F8JuuIxZzNDvklH:0QEO9F8JuuIxZzNDvklI@PassportGameTrans-ENCRYPT-P/topic=PassportGameTrans-T&role=pub&offset=new'
# rc = redis.Redis(host='172.16.33.158', port=6205, socket_keepalive=True)
# print rc.execute_command('auth',auth_history)

# msg = json.dumps(data_history, ensure_ascii=False)
# rc.set("110000193", msg)



# data_history ={
#     "action":"insert",
# 	"table":"aso_account",
# 	"data":{"mid":110000194,
# 	"userid":"test_dan07",
# 	"uname":"test_dan07",
# 	"pwd":"d2c9d4acdfe9942979d7b4d3e711d499",
# 	"salt":"5Brw3JuP",
# 	"email":"169d9106a74d5e95de71be6cf373af04",
# 	"tel":"b65919197178b5db6ed5f85a229dfaf9",
# 	"country_id":1,
# 	"mobile_verified":0,
# 	"isleak":0,
# 	"mtime":"2018-01-02 17:18:58"},
# 	"flag":0
#       }
# auth_history = '0QEO9F8JuuIxZzNDvklH:0QEO9F8JuuIxZzNDvklI@PassportGameTrans-ENCRYPT-P/topic=PassportGameTrans-T&role=pub&offset=new'
# rc = redis.Redis(host='172.16.33.158', port=6205, socket_keepalive=True)
# print rc.execute_command('auth',auth_history)

# msg = json.dumps(data_history, ensure_ascii=False)
# rc.set("110000193", msg)


