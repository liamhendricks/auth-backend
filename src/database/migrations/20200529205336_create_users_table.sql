CREATE TABLE `users` (
  `id`         BINARY(16)		                             NOT NULL,
  `name`       VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL UNIQUE,
  `email`      VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL UNIQUE,
  `password`   VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `user_type`  VARCHAR(20)  COLLATE `utf8mb4_unicode_ci` NOT NULL DEFAULT "Free",
  `created_at` TIMESTAMP                                 NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP                                 NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP                                 NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_email_unique` (`email`),
  KEY `users_email_index` (`email`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
