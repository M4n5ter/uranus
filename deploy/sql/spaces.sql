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

 Date: 23/03/2022 19:03:32
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for spaces
-- ----------------------------
DROP TABLE IF EXISTS `spaces`;
CREATE TABLE `spaces`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `del_state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已经删除',
  `version` bigint(20) NOT NULL DEFAULT 0 COMMENT '版本号',
  `flight_info_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '对应的航班信息id',
  `is_first_class` tinyint(1) NULL DEFAULT NULL COMMENT '是否是头等舱/商务舱',
  `price` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '价格',
  `surplus` bigint(20) NULL DEFAULT NULL COMMENT '剩余量',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_spaces_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
