-- mcn签约表
drop table if exists mcn_sign;
create table mcn_sign (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT 0 COMMENT 'mcn的mid',
  company_name varchar(32) NOT NULL DEFAULT '' COMMENT '企业名称',
  company_license_id varchar(32) NOT NULL DEFAULT '' COMMENT '营业执照注册号',
  company_license_link varchar(255) NOT NULL DEFAULT '' COMMENT '营业执照链接',
  contract_link varchar(255) NOT NULL DEFAULT '' COMMENT '合同链接',
  contact_name varchar(16) NOT NULL DEFAULT '' COMMENT '对接人姓名',
  contact_title varchar(16) NOT NULL DEFAULT '' COMMENT '对接人职务',
  contact_idcard varchar(32) NOT NULL DEFAULT '' COMMENT '对接人身份证号',
  contact_phone varchar(16) NOT NULL DEFAULT '' COMMENT '对接人手机号',
  begin_date date NOT NULL DEFAULT '0000-00-00' COMMENT '合同开始时间',
  end_date date NOT NULL DEFAULT '0000-00-00' COMMENT '合同结束时间',
  reject_reason varchar(255) NOT NULL DEFAULT '' COMMENT '驳回理由',
	`reject_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '驳回时间',
	`pay_expire_state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '付款到期状态:1:未到期 2:即将到期',
  state tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态,0未申请，1待审核，2已驳回，10已签约，11冷却中，12已到期，13封禁，14清退, 15待开启，100移除',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_mcn_mid` (`mcn_mid`),
	KEY `ix_mtime` (`mtime`),
	KEY `ix_state` (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn签约表';

-- mcn付款表
drop table if exists mcn_sign_pay;
CREATE TABLE mcn_sign_pay (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  due_date date NOT NULL DEFAULT '0000-00-00' COMMENT '付款时间',
  pay_value BIGINT(20) NOT NULL DEFAULT '0' COMMENT '金额',
  state tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0未支付，1已支付，100删除',
  note varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (id),
  KEY ix_signid (sign_id),
  KEY ix_mcn_mid (mcn_mid),
  KEY ix_generate_date (due_date),
  KEY ix_mtime (mtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='付款表';


-- mcn up绑定表
drop table if exists mcn_up;
CREATE TABLE mcn_up (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT 0 COMMENT 'mcn的mid',
  up_mid int(11) unsigned NOT NULL DEFAULT 0 COMMENT '绑定up的mid',
  begin_date date NOT NULL DEFAULT '0000-00-00' COMMENT '合同开始时间',
  end_date date NOT NULL DEFAULT '0000-00-00' COMMENT '合同结束时间',
  contract_link varchar(255) NOT NULL DEFAULT '' COMMENT '与up合同链接',
  up_auth_link varchar(255) NOT NULL DEFAULT '' COMMENT 'up授权协议链接',
  reject_reason varchar(255) NOT NULL DEFAULT '' COMMENT '驳回理由',
  reject_time timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '驳回时间',
  state tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0未授权，1已拒绝，2审核中，3已驳回，10已签约，11已冻结，12已到期，13封禁，14已解约，100删除',
  state_change_time timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '状态变化时间',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_sign_id_mcn_mid_up_mid (sign_id, mcn_mid, up_mid),
  KEY ix_up_mid(up_mid),
  KEY ix_mtime (mtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn绑定up表';


-- 数据相关表
-- 1。mcn整体数据表
-- 2。mcn下各up主数据表
-- 3。Top稿件表
-- 1。mcn整体数据表
drop table if exists mcn_data_summary;
create table mcn_data_summary (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  up_count int(11) unsigned NOT NULL DEFAULT '0' COMMENT '签约UP主数',
  fans_count_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计粉丝量',
  fans_count_online int(11) unsigned NOT NULL DEFAULT '0' COMMENT '线上涨粉量',
  fans_count_real int(11) unsigned NOT NULL DEFAULT '0' COMMENT '实际涨粉量',
  fans_count_cheat_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计作弊粉丝',
  fans_count_increase_day int(11) unsigned NOT NULL DEFAULT '0' COMMENT '当日新增粉丝数',
  play_count_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计播放数',
  play_count_increase_day int(11) unsigned NOT NULL DEFAULT '0' COMMENT '当日新增播放数',
  archive_count_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计投稿量',
  active_tid smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '分区,表示某个分区',
  generate_date date NOT NULL DEFAULT '0000-00-00' COMMENT '计算日',
  data_type tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据类型，1按天，2按月',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sign_id_generate_date_active_tid_data_type` (sign_id, generate_date, active_tid, data_type),
  KEY ix_mcn_mid (mcn_mid),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_generate_date` (`generate_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn整体数据';

drop table if exists mcn_data_up_detail;
create table mcn_data_up_detail (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  up_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT '签约UP主ID',
  fans_count_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计粉丝量',
  fans_count_online int(11) unsigned NOT NULL DEFAULT '0' COMMENT '线上涨粉量',
  fans_count_real int(11) unsigned NOT NULL DEFAULT '0' COMMENT '实际涨粉量',
  fans_count_cheat_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计作弊粉丝',
  fans_count_increase_day int(11) unsigned NOT NULL DEFAULT '0' COMMENT '当日新增粉丝数',
  play_count_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计播放数',
  play_count_increase_day int(11) unsigned NOT NULL DEFAULT '0' COMMENT '当日新增播放数',
  archive_count_accumulate int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计投稿量',
  active_tid smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT 'Up所属分区',
  generate_date date NOT NULL DEFAULT '0000-00-00' COMMENT '计算日',
  data_type tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据类型，1按天，2按月',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sign_id_generate_date_data_type_up_mid` (sign_id, generate_date, data_type, up_mid),
  KEY ix_mcn_mid (mcn_mid),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_generate_date` (`generate_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn up整体数据';

-- 2。mcn下各up主数据表
drop table if exists mcn_data_up;
create table mcn_data_up (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  up_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up的mid',
  data_type tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据类型，1累计，2昨日，3上周，4上月',
  fans_increase_accumulate int(11) unsigned NOT NULL default '0' COMMENT '粉丝数增涨量',
  archive_count int(11) unsigned NOT NULL default '0' COMMENT '投搞量',
  play_count int(11) unsigned NOT NULL default '0' COMMENT '播放量',
  fans_increase_month int(11) unsigned NOT NULL default '0' COMMENT '近一个月涨粉量',
  fans_count int(11) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝总量',
  fans_count_active int(11) unsigned NOT NULL DEFAULT '0' COMMENT '活跃粉丝总量',
  generate_date date NOT NULL DEFAULT '0000-00-00' COMMENT '计算日',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sign_id_generate_date_data_type` (sign_id, generate_date, data_type),
  KEY ix_mcn_mid (mcn_mid),
  KEY ix_up_mid (up_mid),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_generate_date` (`generate_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn下up数据';



-- alter table mcn_up_test add column state_change_time timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '状态变化时间';

-- 增加字段, fat 1, uat 1, prod 1
alter table mcn_data_summary add column fans_count_real_accumulate bigint(20)  NOT NULL DEFAULT '0' COMMENT '累计实际涨粉量';
alter table mcn_data_summary add column fans_count_online_accumulate bigint(20)  NOT NULL DEFAULT '0' COMMENT '累计线上涨粉量';
alter table mcn_data_summary add column archive_count_day int(11) NOT NULL DEFAULT '0' COMMENT '当日新增投稿量';

alter table mcn_data_up_detail add column fans_count_real_accumulate bigint(20) NOT NULL DEFAULT '0' COMMENT '累计实际涨粉量';
alter table mcn_data_up_detail add column fans_count_online_accumulate bigint(20) NOT NULL DEFAULT '0' COMMENT '累计线上涨粉量';
alter table mcn_data_up_detail add column archive_count_day int(11) NOT NULL DEFAULT '0' COMMENT '当日新增投稿量';

-- 修改数据字段类型，去掉unsigned， 修改播放相关的为bigint, fat 1, uat 1, prod 1
alter table mcn_data_summary modify column fans_count_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计粉丝量';
alter table mcn_data_summary modify column fans_count_online int(11) NOT NULL DEFAULT '0' COMMENT '线上涨粉量';
alter table mcn_data_summary modify column fans_count_real int(11) NOT NULL DEFAULT '0' COMMENT '实际涨粉量';
alter table mcn_data_summary modify column fans_count_cheat_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计作弊粉丝';
alter table mcn_data_summary modify column fans_count_increase_day int(11) NOT NULL DEFAULT '0' COMMENT '当日新增粉丝数';
alter table mcn_data_summary modify column play_count_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计播放数';
alter table mcn_data_summary modify column play_count_increase_day int(11) NOT NULL DEFAULT '0' COMMENT '当日新增播放数';
alter table mcn_data_summary modify column archive_count_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计投稿量';

alter table mcn_data_up_detail modify column fans_count_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计粉丝量';
alter table mcn_data_up_detail modify column fans_count_online int(11) NOT NULL DEFAULT '0' COMMENT '线上涨粉量';
alter table mcn_data_up_detail modify column fans_count_real int(11) NOT NULL DEFAULT '0' COMMENT '实际涨粉量';
alter table mcn_data_up_detail modify column fans_count_cheat_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计作弊粉丝';
alter table mcn_data_up_detail modify column fans_count_increase_day int(11) NOT NULL DEFAULT '0' COMMENT '当日新增粉丝数';
alter table mcn_data_up_detail modify column play_count_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计播放数';
alter table mcn_data_up_detail modify column play_count_increase_day int(11) NOT NULL DEFAULT '0' COMMENT '当日新增播放数';
alter table mcn_data_up_detail modify column archive_count_accumulate int(11) NOT NULL DEFAULT '0' COMMENT '累计投稿量';

alter table mcn_data_up modify column fans_increase_accumulate int(11) NOT NULL default '0' COMMENT '粉丝数增涨量';
alter table mcn_data_up modify column archive_count int(11) NOT NULL default '0' COMMENT '投搞量';
alter table mcn_data_up modify column play_count bigint(20) NOT NULL default '0' COMMENT '播放量';
alter table mcn_data_up modify column fans_increase_month int(11) NOT NULL default '0' COMMENT '近一个月涨粉量';
alter table mcn_data_up modify column fans_count int(11) NOT NULL DEFAULT '0' COMMENT '粉丝总量';
alter table mcn_data_up modify column fans_count_active int(11) NOT NULL DEFAULT '0' COMMENT '活跃粉丝总量';

-- fat 1, uat 1, prod 1
alter table mcn_data_summary modify column play_count_accumulate bigint(20) NOT NULL DEFAULT '0' COMMENT '累计播放数';
alter table mcn_data_summary modify column play_count_increase_day bigint(20) NOT NULL DEFAULT '0' COMMENT '当日/月新增播放数';
alter table mcn_data_up_detail modify column play_count_accumulate bigint(20) NOT NULL DEFAULT '0' COMMENT '累计播放数';
alter table mcn_data_up_detail modify column play_count_increase_day bigint(20) NOT NULL DEFAULT '0' COMMENT '当日/月新增播放数';
alter table mcn_sign modify column reject_time datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '驳回时间';

-- fat 1, uat 1, prod 1, 增加时间上的索引
create index ix_end_date on mcn_sign (end_date);
create index ix_begin_date on mcn_sign (begin_date);
create index ix_end_date on mcn_up (end_date);
create index ix_begin_date on mcn_up (begin_date);

--------- 2期
-- fat 1, uat 1, prod 1, 增加表
alter table mcn_up add column up_type tinyint(4) not null default '0' comment '用户类型，0为站内，1为站外';
alter table mcn_up add column site_link varchar(255) not null default '' comment 'up主站外账号链接';

-- mcn_data_import_up: table
CREATE TABLE `mcn_data_import_up` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `mcn_mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  `sign_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  `up_mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'UP主 mid',
  `standard_fans_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '达标粉丝数类型, 1: 1w粉丝',
  `standard_fans_date` int(11) NOT NULL DEFAULT '0' COMMENT '达到粉丝数门槛花费的时间，秒',
  `standard_archive_count` int(11) NOT NULL DEFAULT '0' COMMENT '达标时投稿量',
  `standard_fans_count` int(11) NOT NULL DEFAULT '0' COMMENT '达标时粉丝数',
  `is_reward` int(11) NOT NULL DEFAULT '0' COMMENT '奖励情况 0:未奖励 1:已奖励',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sign_id_mid_type` (`sign_id`,`up_mid`,`standard_fans_type`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='引入账号数据表';

-- mcn_up_recommend_pool: table
CREATE TABLE `mcn_up_recommend_pool` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `up_mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up mid',
  `fans_count` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝量',
  `fans_count_increase_month` int(11) NOT NULL DEFAULT '0' COMMENT '本月粉丝增长量',
  `archive_count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '总稿件数',
  `play_count_accumulate` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '累积播放量',
  `play_count_average` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '稿均播放量',
  `active_tid` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '分区,表示某个分区',
  `last_archive_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT '最近投稿时间',
  `state` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '推荐池状态: 1:未推荐 2:推荐 3:禁止推荐 100:移除',
  `source` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '推荐池来源: 1:自动添加(大数据)  2:手动添加',
  `generate_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT '大数据更新时间',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_up_mid` (`up_mid`),
  KEY `ix_state` (`state`),
  KEY `ix_active_tid` (`active_tid`),
  KEY `ix_fans_count` (`fans_count`),
  KEY `ix_play_count_accumulate` (`play_count_accumulate`),
  KEY `ix_play_count_average` (`play_count_average`),
  KEY `ix_fans_count_increase_month` (`fans_count_increase_month`),
  KEY `ix_source` (`source`),
  KEY `ix_generate_time` (`generate_time`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn-up主推荐池';

-- mcn_up_recommend_source: table
CREATE TABLE `mcn_up_recommend_source` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `up_mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up mid',
  `fans_count` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝量',
  `fans_count_increase_month` int(11) NOT NULL DEFAULT '0' COMMENT '本月粉丝增长量',
  `archive_count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '总稿件数',
  `play_count_accumulate` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '累积播放量',
  `play_count_average` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '稿均播放量',
  `active_tid` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '分区,表示某个分区',
  `last_archive_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT '最近投稿时间',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_up_mid` (`up_mid`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn-up主推荐池的来源(大数据提供)';

-- fat 1, uat 1, prod 0
-- 1。涨粉量排名
create table mcn_rank_up_fans (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  up_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up的mid',
  value1 int(11) NOT NULL default '0' COMMENT '数据1',
  value2 int(11) NOT NULL default '0' COMMENT '数据2',
  active_tid smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '分区,表示某个分区',
  data_type tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据类型，1累计（总榜），2昨日，3上周，4上月，5活跃粉丝(累计)',
  generate_date date NOT NULL DEFAULT '0000-00-00' COMMENT '计算日',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sign_id_generate_date_data_type_up_mid` (sign_id, generate_date, data_type, up_mid),
  KEY ix_mcn_mid (mcn_mid),
  KEY ix_up_mid (up_mid),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_generate_date` (`generate_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn下up涨粉量排名';

-- 2。Top稿件表
create table mcn_rank_archive_likes (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  mcn_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  up_mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up的mid',
  archive_id bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '稿件id',
  like_count bigint(20) NOT NULL DEFAULT '0' COMMENT '日/周/月新增点赞数，根据data_type统计',
  data_type tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据类型，1累计，2昨日，3上周，4上月',
  tid smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '分区ID',
  generate_date date NOT NULL DEFAULT '0000-00-00' COMMENT '计算日',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY ix_mcn_mid (mcn_mid),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_generate_date` (`generate_date`),
  UNIQUE KEY `uk_sign_id_generate_date_data_type_archive_id` (sign_id, generate_date, data_type, archive_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='mcn下top稿件表';



-- mcn_data_up_cheat: table
CREATE TABLE `mcn_data_up_cheat` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `mcn_mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn mid',
  `sign_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'mcn签约ID',
  `up_mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主 mid',
  `generate_date` date NOT NULL DEFAULT '0000-00-00' COMMENT '计算日',
  `fans_count_cheat_increase_day` int(11) NOT NULL DEFAULT '0' COMMENT '新增作弊粉丝量',
  `fans_count_cheat_cleaned_accumulate` int(11) NOT NULL DEFAULT '0' COMMENT '已清除粉丝量',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_up_mid_sign_id_generate_date` (`up_mid`,`sign_id`,`generate_date`),
  KEY `ix_mcn_mid` (`mcn_mid`),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_generate_date` (`generate_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='作弊筛选详情表'
;

-- fat 1, uat 1, prod 0
ALTER TABLE mcn_rank_archive_likes CHANGE archive_id avid bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '稿件id';

ALTER TABLE `bilibili_upcrm`.`mcn_data_up_cheat`
  ADD COLUMN `fans_count_cheat_accumulate` int(11) NOT NULL DEFAULT '0' COMMENT '累计作弊粉丝',
  ADD COLUMN `fans_count_accumulate` int(11) NOT NULL DEFAULT '0' COMMENT '实际粉丝量';

-- fat 1, uat 1, prod 0
ALTER TABLE mcn_up ADD COLUMN confirm_time timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT 'up确认时间';


--------四期-------
alter table mcn_sign add permission int(11) unsigned default '1' not null comment '权限列表-属性位';
alter table mcn_up
add permission int(11) unsigned default '1' not null comment '权限列表-属性位',
add publication_price bigint default '0' not null comment '刊例价(千分位*1000)';