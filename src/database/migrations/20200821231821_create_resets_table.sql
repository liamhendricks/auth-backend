CREATE TABLE `resets` (
  `id`         BINARY(16) NOT NULL,
  `user_id`    BINARY(16) NOT NULL,
  `token`      BINARY(16) NOT NULL,
  `expiration` TIMESTAMP  NOT NULL,
  `created_at` TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP  NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP  NULL DEFAULT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
