/*
 Navicat Premium Data Transfer

 Source Server         : test2
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 192.168.75.132:3306
 Source Schema         : flight

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 30/03/2022 10:41:28
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for payment
-- ----------------------------
DROP TABLE IF EXISTS `payment`;
CREATE TABLE `payment`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `sn` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '流水单号',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_state` tinyint(1) NOT NULL DEFAULT 0,
  `version` bigint NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '用户id',
  `pay_mode` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支付方式 1:微信支付 2:支付宝 3:钱包余额',
  `trade_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '第三方支付类型',
  `trade_state` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '第三方交易状态',
  `pay_total` bigint NOT NULL DEFAULT 0 COMMENT '支付总金额(分)',
  `transaction_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '第三方支付单号',
  `trade_state_desc` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '支付状态描述',
  `order_sn` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '业务单号',
  `pay_status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '平台内交易状态(第三方支付这条也要有)   -2：回滚废弃 -1:支付失败 0:未支付 1:支付成功 2:已退款',
  `pay_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT '支付成功时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sn`(`sn`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '第三方支付流水记录' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
