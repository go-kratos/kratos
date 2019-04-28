drop table if exists offline_activity_info;
create table offline_activity_info (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID,活动ID',
  title varchar(32) NOT NULL DEFAULT '' COMMENT '标题',
  link varchar(255) NOT NULL DEFAULT '' COMMENT '活动链接',
  bonus_type tinyint(4) NOT NULL DEFAULT 0 COMMENT '0,奖品；1,奖金',
  memo varchar(32) NOT NULL DEFAULT '' COMMENT '备注',
  creator varchar(16) NOT NULL DEFAULT '' COMMENT '创建者',
  state tinyint(4) NOT NULL DEFAULT 0 COMMENT '0初始；1发送中；2等待结果；10处理完成；100删除、无效',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='线下活动信息';

drop table if exists offline_activity_bonus;
create table offline_activity_bonus (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID, bonusID',
  activity_id int(11) unsigned NOT NULL DEFAULT 0 COMMENT '活动ID',
  total_money bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '奖项总金额，1=1/1000元',
  member_count int(11) unsigned NOT NULL DEFAULT 0 COMMENT '总人数',
  state tinyint(4) NOT NULL DEFAULT 0 COMMENT '0正常；100删除、无效',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='线下奖项信息';

drop table if exists offline_activity_result;
create table offline_activity_result (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  activity_id int(11) unsigned NOT NULL DEFAULT 0 COMMENT '活动ID',
  bonus_id int(11) unsigned NOT NULL DEFAULT 0 COMMENT '奖励ID',
  bonus_type tinyint(4) NOT NULL DEFAULT 0 COMMENT '0,奖品；1,奖金',
  mid int(11) unsigned NOT NULL DEFAULT 0 COMMENT 'memberID',
  state tinyint(4) NOT NULL DEFAULT 0 COMMENT '0初始；1审核中；10成功；11失败',
  bonus_money bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '奖励，1=1/1000元',
  order_id varchar(32) NOT NULL DEFAULT '' COMMENT '交易id,需要唯一',

  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_aid_bid_mid` (activity_id, bonus_id, mid),
  KEY `ix_mid` (mid),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='线下活动奖励发放信息';

drop table if exists offline_activity_shell_order;
create table offline_activity_shell_order (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  result_id int(11) unsigned NOT NULL COMMENT '对应的result ID',
  order_id varchar(32) NOT NULL DEFAULT '' COMMENT '交易id,需要唯一',
  order_status varchar(16) NOT NULL DEFAULT '' COMMENT '订单状态',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_id` (order_id),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='贝壳订单';

-- fat 1, uat 1, prod 1
alter table offline_activity_info add key ix_ctime(ctime);
alter table offline_activity_bonus add key ix_activity_id(activity_id);