CREATE TABLE `lesson_block` (
  `id`             BINARY(16)                                NOT NULL,
  `lesson_data_id` BINARY(16)                                NOT NULL,
  `name`           VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `header`         TEXT COLLATE `utf8mb4_unicode_ci`,
  `created_at`     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`     TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`     TIMESTAMP NULL DEFAULT NULL,
  FOREIGN KEY (`lesson_data_id`) REFERENCES `lesson_data` (`id`),
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
