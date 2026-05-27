CREATE DATABASE IF NOT EXISTS hertz_demo DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE hertz_demo;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL COMMENT '用户名',
    email VARCHAR(128) NOT NULL COMMENT '邮箱',
    password VARCHAR(256) NOT NULL COMMENT '密码',
    nickname VARCHAR(64) DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(256) DEFAULT '' COMMENT '头像',
    status TINYINT DEFAULT 1 COMMENT '状态: 1正常 2禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY idx_username (username),
    UNIQUE KEY idx_email (email),
    KEY idx_status (status),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

