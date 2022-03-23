/*
 Navicat Premium Data Transfer

 Source Server         : putpp.com
 Source Server Type    : MySQL
 Source Server Version : 50736
 Source Host           : putpp.com:13306
 Source Schema         : zero_gorm

 Target Server Type    : MySQL
 Target Server Version : 50736
 File Encoding         : 65001

 Date: 23/03/2022 19:03:17
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for flight_infos
-- ----------------------------
DROP TABLE IF EXISTS `flight_infos`;
CREATE TABLE `flight_infos`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `del_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已经删除',
  `version` bigint(20) NOT NULL DEFAULT 0 COMMENT '版本号',
  `flight_number` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '对应的航班号',
  `set_out_date` datetime(3) NULL DEFAULT NULL COMMENT '出发日期',
  `punctuality` tinyint(3) UNSIGNED NULL DEFAULT NULL COMMENT '准点率',
  `start_position` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '起飞地点',
  `start_time` datetime(3) NULL DEFAULT NULL COMMENT '起飞时间',
  `end_position` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '降落地点',
  `end_time` datetime(3) NULL DEFAULT NULL COMMENT '降落时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_flight_infos_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
