#初始table
CREATE TABLE `up_rank` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  `type` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '排行榜类型',
  `value` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排行榜数值,根据type不同，代表的含义不同',
  `generate_date` date NOT NULL DEFAULT '0000-00-00' COMMENT '排行榜日',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `value2` int(11) NOT NULL DEFAULT '0' COMMENT '分数2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_type_date` (`generate_date`,`type`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Up蹿升榜'


# fat uat表增加字段
# 增加默认值
alter table up_base_info alter column active_tid  set default 0;
alter table up_base_info alter column attr  set default 0;
alter table up_base_info alter column mid  set default 0;
alter table up_base_info alter column mid  set default 0;
alter table up_play_info alter column mid  set default 0;
alter table up_play_info alter column business_type  set default 0;
alter table up_play_info alter column article_count  set default 0;
alter table up_play_info alter column play_count_90day  set default 0;
alter table up_play_info alter column play_count_7day  set default 0;
alter table up_play_info alter column play_count_30day  set default 0;
alter table up_play_info alter column play_count_accumulate  set default 0;

alter table up_stats_history alter column type  set default 0;
alter table up_stats_history alter column sub_type  set default 0;
alter table up_stats_history alter column generate_date  set default '0000-00-00';

alter table up_rank alter column mid  set default 0;
alter table up_rank alter column type  set default 0;
alter table up_rank alter column value  set default 0;
alter table up_rank alter column generate_date set default '0000-00-00';

alter table task_info alter column generate_date  set default '0000-00-00';
alter table task_info alter column task_type  set default 0;

alter table up_base_info add column active_tid smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '最多稿件分区';
alter table up_base_info add column attr int(11) NOT NULL COMMENT '属性，以位区分';
alter table up_rank add column value2 int(11) NOT NULL DEFAULT 0 COMMENT '分数2';
ALTER TABLE up_base_info MODIFY active_tid smallint(6) unsigned NOT NULL DEFAULT 0 COMMENT '最多稿件分区';

#FAT
alter table up_play_info drop column play_count_avg;
alter table up_play_info drop column play_count_avg_90day;
alter table up_play_info add column `article_count` int(11) NOT NULL DEFAULT '0' COMMENT '总稿件数';
alter table up_play_info add column `play_count_90day` int(11) NOT NULL DEFAULT '0' COMMENT '90天内稿件总播放次数';
alter table up_play_info add column `play_count_30day` int(11) NOT NULL DEFAULT '0' COMMENT '30天内稿件总播放次数';
alter table up_play_info add column `play_count_7day` int(11) NOT NULL DEFAULT '0' COMMENT '7天内稿件总播放次数';

#FAT UAT 
DROP INDEX ix_mid ON up_base_info;
alter table up_base_info add unique key uk_mid_type (`mid`,`business_type`);

#每天一条记录 增加分数段表， (uat 1, fat 1, prod 1)
create table score_section_history(
    id int(11) unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '自增ID',
    generate_date date NOT NULL COMMENT '生成日期',
    score_type SMALLINT(6) NOT NULL DEFAULT 0 COMMENT '类型, 1质量分，2影响分，3信用分',
    section_0 int(11) NOT NULL DEFAULT 0 COMMENT '0~100的人数',
    section_1 int(11) NOT NULL DEFAULT 0 COMMENT '101~200',
    section_2 int(11) NOT NULL DEFAULT 0 COMMENT '201~300',
    section_3 int(11) NOT NULL DEFAULT 0 COMMENT '301~400',
    section_4 int(11) NOT NULL DEFAULT 0 COMMENT '401~500',
    section_5 int(11) NOT NULL DEFAULT 0 COMMENT '501~600',
    section_6 int(11) NOT NULL DEFAULT 0 COMMENT '601~700',
    section_7 int(11) NOT NULL DEFAULT 0 COMMENT '701~800',
    section_8 int(11) NOT NULL DEFAULT 0 COMMENT '801~900',
    section_9 int(11) NOT NULL DEFAULT 0 COMMENT '901~1000',
    ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	  mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    UNIQUE uk_date_type(generate_date, score_type),
    KEY ix_mtime (mtime)
) engine=innodb DEFAULT charset=utf8 comment='up分数段人数分布表';

#增加信用分、影响分、质量分字段 (uat 1, fat 1, prod 1)
alter table up_base_info add COLUMN credit_score INT NOT NULL DEFAULT 500 COMMENT '信用分';
alter table up_base_info add COLUMN pr_score INT NOT NULL DEFAULT 0 COMMENT '影响分';
alter table up_base_info add COLUMN quality_score INT NOT NULL DEFAULT 0 COMMENT '质量分';

#修复key错误 (uat 1, fat 1, prod 1)
DROP INDEX uk_type_date ON up_rank;
alter table up_rank add unique key uk_date_type_mid (`generate_date`,`type`, `mid`);

#增加生日、地域等字段 (uat 1, fat 1, prod 1)
alter table up_base_info add COLUMN birthday DATE NOT NULL DEFAULT '0000-00-00' COMMENT '生日';
alter table up_base_info add COLUMN active_province varchar(32) NOT NULL DEFAULT '' COMMENT '省份';
alter table up_base_info add COLUMN active_city varchar(32) NOT NULL DEFAULT '' COMMENT '城市';

#增加task info的unique key(uat 1, fat 1, prod 1)
alter table task_info add unique key uk_date_type (`generate_date`,`task_type`);