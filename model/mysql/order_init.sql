-- Drop database if exists
DROP DATABASE IF EXISTS `order_db`;

-- Create database
CREATE DATABASE IF NOT EXISTS `order_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Use database
USE `order_db`;

DROP TABLE IF EXISTS `t_order`;
CREATE TABLE `t_order` (
  `transaction_id` VARCHAR(64) NOT NULL COMMENT '交易ID',
  `out_order_no` VARCHAR(64) NOT NULL COMMENT '商户订单号',
  `merchant_id` VARCHAR(64) NOT NULL COMMENT '商户ID',
  `merchant_uid` BIGINT NOT NULL COMMENT '商户uid',
  `merchant_name` VARCHAR(64) NOT NULL COMMENT '商户名称',
  `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
  `uid` BIGINT NOT NULL COMMENT '用户uid',
  `trade_state` TINYINT NOT NULL COMMENT '交易状态',
  `amount` BIGINT NOT NULL COMMENT '金额',
  `cur_type` TINYINT NOT NULL COMMENT '货币类型',
  `pay_type` TINYINT NOT NULL COMMENT '支付类型',
  `pay_time` datetime NOT NULL COMMENT '支付时间',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`transaction_id`),
  INDEX `idx_create_time` (`create_time`),
  INDEX `idx_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `t_recharge_withdraw_order`;
CREATE TABLE `t_recharge_withdraw_order` (
  `transaction_id` VARCHAR(64) NOT NULL COMMENT '充值提现单号(trade_itg Bank2CPre/C2BankPre生成)',
  `channel_transaction_id` VARCHAR(64) NOT NULL COMMENT '银行渠道单号(bank_channel生成)',
  `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
  `uid` BIGINT NOT NULL COMMENT '用户uid',
  `order_type` TINYINT NOT NULL COMMENT '单据类型: 1=充值(Bank2C) 2=提现(C2Bank)',
  `bank_type` TINYINT NOT NULL COMMENT '银行类型',
  `amount` BIGINT NOT NULL COMMENT '金额',
  `cur_type` TINYINT NOT NULL COMMENT '货币类型',
  `trade_state` TINYINT NOT NULL COMMENT '交易状态: 1=初始 2=成功 3=失败',
  `memo` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`transaction_id`),
  UNIQUE INDEX `uk_channel_transaction_id` (`channel_transaction_id`),
  INDEX `idx_create_time` (`create_time`),
  INDEX `idx_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- linux:  mysql -h 127.0.0.1 -P 3306 -u root -proot123456 < order_init.sql
-- windows: Get-Content -Encoding UTF8 order_init.sql | mysql -h 127.0.0.1 -P 3306 -u root -proot123456
-- 只读权限 multipass exec master1 -- sudo kubectl exec -it -n pay-ns mysql-0 -- mysql -ustarslipay -ppayClipayA2026
-- root权限 multipass exec master1 -- sudo kubectl exec -it -n pay-ns mysql-0 -- mysql -uroot -proot123456


-- select * from user_db.t_uid_segment limit 2;