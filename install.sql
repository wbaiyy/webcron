/*
 Navicat MySQL Data Transfer

 Source Server         : 10.40.6.26
 Source Server Type    : MySQL
 Source Server Version : 50623
 Source Host           : 10.40.6.26:3306
 Source Schema         : web_cron

 Target Server Type    : MySQL
 Target Server Version : 50623
 File Encoding         : 65001

 Date: 15/05/2020 17:48:25
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_action
-- ----------------------------
DROP TABLE IF EXISTS `t_action`;
CREATE TABLE `t_action` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `action` varchar(20) NOT NULL DEFAULT '',
  `actor` varchar(20) NOT NULL DEFAULT '',
  `object_type` varchar(20) NOT NULL DEFAULT '',
  `object_id` int(11) NOT NULL DEFAULT '0',
  `extra` varchar(1000) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_env
-- ----------------------------
DROP TABLE IF EXISTS `t_env`;
CREATE TABLE `t_env` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` int(11) NOT NULL DEFAULT '0',
  `name` varchar(20) NOT NULL DEFAULT '',
  `ssh_user` varchar(20) NOT NULL DEFAULT '',
  `ssh_port` varchar(10) NOT NULL DEFAULT '',
  `ssh_key` varchar(100) NOT NULL DEFAULT '',
  `pub_dir` varchar(100) NOT NULL DEFAULT '',
  `before_shell` longtext NOT NULL,
  `after_shell` longtext NOT NULL,
  `server_count` int(11) NOT NULL DEFAULT '0',
  `send_mail` int(11) NOT NULL DEFAULT '0',
  `mail_tpl_id` int(11) NOT NULL DEFAULT '0',
  `mail_to` varchar(1000) NOT NULL DEFAULT '',
  `mail_cc` varchar(1000) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `t_env_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_env_server
-- ----------------------------
DROP TABLE IF EXISTS `t_env_server`;
CREATE TABLE `t_env_server` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` int(11) NOT NULL DEFAULT '0',
  `env_id` int(11) NOT NULL DEFAULT '0',
  `server_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `t_env_server_env_id` (`env_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_mail_tpl
-- ----------------------------
DROP TABLE IF EXISTS `t_mail_tpl`;
CREATE TABLE `t_mail_tpl` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '',
  `subject` varchar(200) NOT NULL DEFAULT '',
  `content` longtext NOT NULL,
  `mail_to` varchar(1000) NOT NULL DEFAULT '',
  `mail_cc` varchar(1000) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_perm
-- ----------------------------
DROP TABLE IF EXISTS `t_perm`;
CREATE TABLE `t_perm` (
  `module` varchar(20) NOT NULL DEFAULT '' COMMENT '模块名',
  `action` varchar(20) NOT NULL DEFAULT '' COMMENT '操作名',
  UNIQUE KEY `module` (`module`,`action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_project
-- ----------------------------
DROP TABLE IF EXISTS `t_project`;
CREATE TABLE `t_project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `domain` varchar(100) NOT NULL DEFAULT '',
  `version` varchar(20) NOT NULL DEFAULT '',
  `version_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `repo_url` varchar(100) NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '0',
  `error_msg` longtext NOT NULL,
  `agent_id` int(11) NOT NULL DEFAULT '0' COMMENT '跳板机ID',
  `ignore_list` longtext NOT NULL,
  `before_shell` longtext NOT NULL,
  `after_shell` longtext NOT NULL,
  `create_verfile` int(11) NOT NULL DEFAULT '0',
  `verfile_path` varchar(50) NOT NULL DEFAULT '',
  `task_review` tinyint(4) NOT NULL DEFAULT '0' COMMENT '发布是否需要审批',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_ptask_log
-- ----------------------------
DROP TABLE IF EXISTS `t_ptask_log`;
CREATE TABLE `t_ptask_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `ptask_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext NOT NULL COMMENT '任务输出',
  `error` text NOT NULL COMMENT '错误信息',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `process_time` int(11) NOT NULL DEFAULT '0' COMMENT '消耗时间/毫秒',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_ptask_id` (`ptask_id`,`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=2719 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for t_ptasks
-- ----------------------------
DROP TABLE IF EXISTS `t_ptasks`;
CREATE TABLE `t_ptasks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '任务名称',
  `command` varchar(255) NOT NULL DEFAULT '' COMMENT '任务命令脚本',
  `retry_times` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '重试次数',
  `interval_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '重启间隔时间',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `group_id` tinyint(2) NOT NULL DEFAULT '0' COMMENT '分组ID',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `run_status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '运行状态，1:未开启,2运行中,3暂停中',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '任务状态1:可用，2关闭',
  `output_file` varchar(255) NOT NULL DEFAULT '' COMMENT '任务输出文件',
  `notify_users` varchar(255) NOT NULL DEFAULT '' COMMENT '任务结束通知人',
  `num` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '常驻任务进程数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_ptasks
-- ----------------------------
BEGIN;
INSERT INTO `t_ptasks` VALUES (6, 'auditMessage-test-1', '/usr/local/audit_engine/audit-engine -c /usr/local/audit_engine/environment1/config_prod.json -q auditMessage_OBS,auditMessageCoupon_OBS --consume', 10, 300, '消息审核-测试环境-第一套', 3, 1575949953, 1, 1, '/var/log/audit_engine/environment1/test_auditMessage_OBS.log', '608395,607084', 1);
INSERT INTO `t_ptasks` VALUES (10, 'auditMessage-test-2', '/usr/local/audit_engine/audit-engine -c /usr/local/audit_engine/environment2/config_prod.json -q auditMessage_OBS,auditMessageCoupon_OBS --consume', 3, 10, '消息审核-测试环境-第二套', 3, 1575949925, 1, 1, '/var/log/audit_engine/environment2/test_auditMessage_OBS.log', '608395,607084', 1);
INSERT INTO `t_ptasks` VALUES (11, 'resultAndRevoke_test_1', '/usr/local/audit_engine/audit-engine -c /usr/local/audit_engine/environment1/config_prod.json -q obsAuditResult_OBS,auditRevoke_OBS --consume', 10, 100, '审核结果通知消息撤销服务--测试环境-第一套', 3, 1575949920, 1, 1, '/var/log/audit_engine/environment1/test_obsAuditResult_OBS.log', '608395,607084', 1);
INSERT INTO `t_ptasks` VALUES (12, 'resultAndRevoke_test_2', '/usr/local/audit_engine/audit-engine -c /usr/local/audit_engine/environment2/config_prod.json -q obsAuditResult_OBS,auditRevoke_OBS --consume', 3, 10, '审核结果通知消息撤销服务--测试环境-第二套', 3, 1575949913, 1, 1, '/var/log/audit_engine/environment2/test_obsAuditResult_OBS.log', '608395,607084', 1);
INSERT INTO `t_ptasks` VALUES (14, 'orderEvent', 'php /opt/htdocs/sunzhiming/gearbest-task/artisan task:event --queue=orderEvent_GB', 3, 30, 'task系统订单事件任务', 4, 1579062102, 3, 1, '', '', 1);
INSERT INTO `t_ptasks` VALUES (15, '会员事件', 'php /opt/htdocs/sunzhiming/gearbest-task/artisan task:event --queue=userEvent_GB', 3, 10, '会员事件', 4, 1583302728, 1, 1, '', '', 1);
INSERT INTO `t_ptasks` VALUES (16, '评论事件邮件', 'php /opt/htdocs/sunzhiming/gearbest-task/artisan task:event --queue=reviewEvent_GB', 3, 10, '', 4, 1562060351, 3, 1, '/data/logs/supervisor/task_event_review.log', '', 1);
INSERT INTO `t_ptasks` VALUES (17, '到货通知发送邮件', 'php /opt/htdocs/sunzhiming/gearbest-task/artisan11 task:event --queue=goodsEvent_GB1', 3, 10, '', 4, 1568018813, 1, 1, '', '', 1);
INSERT INTO `t_ptasks` VALUES (18, '支付汇率变更通知', 'php /opt/htdocs/liuhua1/gearbest-task/artisan task:event --queue=noticeSync_GB', 10, 30, '支付汇率变更vv、邮件通知', 4, 1566886039, 1, 1, '', '', 1);
INSERT INTO `t_ptasks` VALUES (19, 'ejob客户端', 'sudo /usr/local/ejob/EJobGo1.0.0_rc.06/ejob-server start', 3, 10, 'ejob客户端脚本', 4, 1575950815, 1, 1, '', '607084', 1);
INSERT INTO `t_ptasks` VALUES (25, '发起coupon审核', 'php /opt/htdocs/wangbei/gearbest-task/artisan task:event --queue=promotionCouponMsgForApproveCenter_GB', 10, 120, '（第二套）soa->task->审核这东西', 4, 1579158966, 1, 1, '', '', 1);
INSERT INTO `t_ptasks` VALUES (26, 'coupon审核结果通知', 'php /opt/htdocs/wangbei/gearbest-task/artisan task:event --queue=auditResultNotify_GB', 10, 120, '（第二套）审核中心->task->soa', 4, 1579158947, 1, 1, '', '', 1);
COMMIT;

-- ----------------------------
-- Table structure for t_role
-- ----------------------------
DROP TABLE IF EXISTS `t_role`;
CREATE TABLE `t_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(20) NOT NULL DEFAULT '',
  `project_ids` varchar(1000) NOT NULL DEFAULT '',
  `description` varchar(200) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_role_perm
-- ----------------------------
DROP TABLE IF EXISTS `t_role_perm`;
CREATE TABLE `t_role_perm` (
  `role_id` int(11) unsigned NOT NULL,
  `perm` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`role_id`,`perm`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_server
-- ----------------------------
DROP TABLE IF EXISTS `t_server`;
CREATE TABLE `t_server` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type_id` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0:普通服务器, 1:跳板机',
  `ip` varchar(20) NOT NULL DEFAULT '' COMMENT '服务器IP',
  `area` varchar(20) NOT NULL DEFAULT '' COMMENT '机房',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '描述',
  `ssh_port` int(11) NOT NULL COMMENT 'ssh端口',
  `ssh_user` varchar(50) NOT NULL DEFAULT '' COMMENT 'ssh帐号',
  `ssh_pwd` varchar(100) NOT NULL DEFAULT '' COMMENT 'ssh密码',
  `ssh_key` varchar(100) NOT NULL DEFAULT '' COMMENT 'sshkey路径',
  `work_dir` varchar(100) NOT NULL DEFAULT '' COMMENT '工作目录',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_task
-- ----------------------------
DROP TABLE IF EXISTS `t_task`;
CREATE TABLE `t_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '分组ID',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `task_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(100) NOT NULL DEFAULT '' COMMENT '时间表达式',
  `concurrent` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否只允许一个实例',
  `command` text NOT NULL COMMENT '命令详情',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0停用 1启用',
  `notify` tinyint(4) NOT NULL DEFAULT '0' COMMENT '通知设置',
  `notify_email` text NOT NULL COMMENT '通知人列表',
  `timeout` smallint(6) NOT NULL DEFAULT '0' COMMENT '超时设置',
  `execute_times` int(11) NOT NULL DEFAULT '0' COMMENT '累计执行次数',
  `prev_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上次执行时间',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_task
-- ----------------------------
BEGIN;
INSERT INTO `t_task` VALUES (1, 1, 3, '重启测试审核中心', 0, '每隔1个小时重启一次', '0 1-16/5 * * * ?', 0, 'curl  http://10.40.2.132:8811/api/restart?ids=6,10,11,12', 0, 0, '', 10, 2685, 1579155360, 1550491333);
INSERT INTO `t_task` VALUES (2, 2, 3, '清理审核中心日志', 0, '定期清理审核中心日志', '0 0 0  15/30 * ?', 0, '/usr/local/audit_engine/cleanLog.sh', 1, 0, '', 0, 4, 1589472251, 1579169990);
COMMIT;

-- ----------------------------
-- Table structure for t_task_group
-- ----------------------------
DROP TABLE IF EXISTS `t_task_group`;
CREATE TABLE `t_task_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_name` varchar(50) NOT NULL DEFAULT '' COMMENT '组名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_task_group
-- ----------------------------
BEGIN;
INSERT INTO `t_task_group` VALUES (1, 1, 'GB网站', 'GB网站', 0);
INSERT INTO `t_task_group` VALUES (2, 1, 'OBS任务', 'OBS任务', 0);
INSERT INTO `t_task_group` VALUES (3, 1, '审核中心', '审核中心任务', 0);
INSERT INTO `t_task_group` VALUES (4, 1, 'task任务', 'gearbest-task任务', 0);
COMMIT;

-- ----------------------------
-- Table structure for t_task_log
-- ----------------------------
DROP TABLE IF EXISTS `t_task_log`;
CREATE TABLE `t_task_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext NOT NULL COMMENT '任务输出',
  `error` text NOT NULL COMMENT '错误信息',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `process_time` int(11) NOT NULL DEFAULT '0' COMMENT '消耗时间/毫秒',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`,`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=2890 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(10) NOT NULL DEFAULT '' COMMENT '密码盐',
  `last_login` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常 -1禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_user
-- ----------------------------
BEGIN;
INSERT INTO `t_user` VALUES (1, 'admin', 'admin@example.com', '7fef6171469e80d32c0559f88b377245', '', 1562645153, '10.37.5.235', 0);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
