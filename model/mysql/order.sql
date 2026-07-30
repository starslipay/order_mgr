
CREATE TABLE `t_order` (
  `transaction_id` VARCHAR(64) COMMENT '交易ID',
  `spid` VARCHAR(64) NOT NULL COMMENT '商户ID',
  `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
  `uid` BIGINT NOT NULL COMMENT '用户uid',
  `amount` BIGINT NOT NULL COMMENT '金额',
  `cur_type` TINYINT NOT NULL COMMENT '货币类型',
  `trade_state` TINYINT NOT NULL COMMENT '交易状态',
  `pay_type` TINYINT NOT NULL COMMENT '支付类型',
  `pay_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '支付时间',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`transaction_id`),
  INDEX `idx_create_time` (`create_time`),
  INDEX `idx_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- goctl model mysql ddl -src order.sql -dir .
-- -c：开启缓存（redis，可选，不加则无缓存）

