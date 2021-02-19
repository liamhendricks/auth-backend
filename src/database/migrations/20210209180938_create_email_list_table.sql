CREATE TABLE `email_list` (
  `id`    BINARY(16)                                NOT NULL,
  `email` VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL UNIQUE,
  `name`  VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL UNIQUE,
  `created_at` TIMESTAMP                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP                            NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP                            NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
