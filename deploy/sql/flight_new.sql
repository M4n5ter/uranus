/*
 Navicat Premium Data Transfer

 Source Server         : putpp.com
 Source Server Type    : MySQL
 Source Server Version : 50736
 Source Host           : putpp.com:13306
 Source Schema         : test

 Target Server Type    : MySQL
 Target Server Version : 50736
 File Encoding         : 65001

 Date: 24/03/2022 14:20:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for flight_infos
-- ----------------------------
DROP TABLE IF EXISTS `flight_infos`;
CREATE TABLE `flight_infos`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime,
  `del_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已经删除',
  `version` bigint(20) NOT NULL DEFAULT 0 COMMENT '版本号',
  `flight_number` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'unknown' COMMENT '对应的航班号',
  `set_out_date` datetime(6) NOT NULL COMMENT '出发日期',
  `punctuality` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '准点率(%)',
  `depart_position` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'unknown' COMMENT '起飞地点',
  `depart_time` datetime(6) NOT NULL COMMENT '起飞时间',
  `arrive_position` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'unknown' COMMENT '降落地点',
  `arrive_time` datetime(6) NOT NULL COMMENT '降落时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_flight_infos_delete_time`(`delete_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for flights
-- ----------------------------
DROP TABLE IF EXISTS `flights`;
CREATE TABLE `flights`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime,
  `del_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已经删除',
  `version` bigint(20) NOT NULL DEFAULT 0 COMMENT '版本号',
  `number` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'unknown' COMMENT '航班号(YT1234)',
  `flt_type` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'unknown' COMMENT '机型',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `number`(`number`) USING BTREE,
  INDEX `idx_flights_delete_time`(`delete_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for refund_and_change_infos
-- ----------------------------
DROP TABLE IF EXISTS `refund_and_change_infos`;
CREATE TABLE `refund_and_change_infos`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime,
  `del_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已经删除',
  `version` bigint(20) NOT NULL DEFAULT 0 COMMENT '版本号',
  `ticket_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '对应的票ID',
  `is_refund` tinyint(1) NOT NULL COMMENT '1为退订，0为改票',
  `time1` datetime(6) NOT NULL COMMENT '时间1',
  `fee1` bigint(20) UNSIGNED NOT NULL DEFAULT 99999 COMMENT '符合时间1时需要的费用(￥/人)',
  `time2` datetime(6) NOT NULL COMMENT '时间2',
  `fee2` bigint(20) UNSIGNED NOT NULL DEFAULT 99999 COMMENT '符合时间2时需要的费用(￥/人)',
  `time3` datetime(6) NOT NULL COMMENT '时间3',
  `fee3` bigint(20) UNSIGNED NOT NULL DEFAULT 99999 COMMENT '符合时间3时需要的费用(￥/人)',
  `time4` datetime(6) NOT NULL COMMENT '时间4',
  `fee4` bigint(20) UNSIGNED NOT NULL DEFAULT 99999 COMMENT '符合时间4时需要的费用(￥/人)',
  `time5` datetime(6) NULL DEFAULT NULL COMMENT '时间5',
  `fee5` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '符合时间5时需要的费用(￥/人)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_refund_and_change_infos_delete_time`(`delete_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for spaces
-- ----------------------------
DROP TABLE IF EXISTS `spaces`;
CREATE TABLE `spaces`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime,
  `del_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已经删除',
  `version` bigint(20) NOT NULL DEFAULT 0 COMMENT '版本号',
  `flight_info_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '对应的航班信息id',
  `is_first_class` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是头等舱/商务舱',
  `total` bigint(20) NOT NULL DEFAULT 0 COMMENT '总量',
  `surplus` bigint(20) NOT NULL DEFAULT 0 COMMENT '剩余量',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_spaces_delete_time`(`delete_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tickets
-- ----------------------------
DROP TABLE IF EXISTS `tickets`;
CREATE TABLE `tickets`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime,
  `del_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已经删除',
  `version` bigint(20) NOT NULL DEFAULT 0 COMMENT '版本号',
  `space_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '对应舱位ID',
  `price` bigint(20) UNSIGNED NOT NULL DEFAULT 999999 COMMENT '价格(￥)',
  `discount` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '折扣(-n%)',
  `cba` tinyint(3) UNSIGNED NOT NULL DEFAULT 20 COMMENT '托运行李额(KG)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tickets_delete_time`(`delete_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
