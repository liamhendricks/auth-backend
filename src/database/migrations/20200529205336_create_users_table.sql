CREATE TABLE `users` (
  `id`         BINARY(16)		                             NOT NULL,
  `name`       VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `email`      VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL UNIQUE,
  `password`   VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
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
