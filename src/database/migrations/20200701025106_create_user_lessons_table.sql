CREATE TABLE `user_lessons` (
  `user_id`    BINARY(16)	NOT NULL,
  `lesson_id`  BINARY(16)	NOT NULL,
  `created_at` TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP  NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP  NULL DEFAULT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  FOREIGN KEY (`lesson_id`) REFERENCES `lessons` (`id`),
  UNIQUE KEY (`user_id`, `lesson_id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
