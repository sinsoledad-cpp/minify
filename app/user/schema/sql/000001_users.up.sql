CREATE TABLE `users` (
     `id` BIGINT UNSIGNED AUTO_INCREMENT ,
     `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名，用于登录',
     `email` VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱，用于登录或找回密码',
     `password_hash` VARCHAR(255) NOT NULL COMMENT 'bcrypt 哈希后的密码',
     `role` VARCHAR(20) NOT NULL DEFAULT 'user' COMMENT '用户角色 (e.g., user, admin)，用于 Casbin',
     `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间 (微秒精度)',
     `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间 (微秒精度)',
      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
