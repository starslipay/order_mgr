
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

-- goctl model mysql ddl -src order.sql -dir .
-- -c：开启缓存（redis，可选，不加则无缓存）

