/*
 Navicat Premium Data Transfer

 Source Server         : uranus
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 127.0.0.1:33069
 Source Schema         : flight

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 10/04/2022 20:18:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for flight_order
-- ----------------------------
DROP TABLE IF EXISTS `flight_order`;
CREATE TABLE `flight_order`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_state` tinyint NOT NULL DEFAULT 0,
  `version` bigint NOT NULL DEFAULT 0 COMMENT '版本号',
  `sn` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '订单号',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '下单用户id',
  `ticket_id` bigint NOT NULL DEFAULT 0 COMMENT '票id',
  `depart_position` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '起飞地点',
  `depart_time` datetime NOT NULL COMMENT '起飞时间',
  `arrive_position` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '降落地点',
  `arrive_time` datetime NOT NULL COMMENT '降落时间',
  `ticket_price` bigint NOT NULL COMMENT '票价(分)',
  `discount` bigint NOT NULL DEFAULT 0 COMMENT '折扣(-n%)',
  `trade_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '-1: 已取消 0:待支付 1:未使用 2:已使用  3:已退款 4:已过期',
  `trade_code` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '确认码',
  `order_total_price` bigint NOT NULL DEFAULT 0 COMMENT '订单总价格(分)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sn`(`sn`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
