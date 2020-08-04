CREATE TABLE `lessons` (
  `id`          BINARY(16)	                              NOT NULL,
  `name`        VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `lesson_type` VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `created_at`  TIMESTAMP                                 NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`  TIMESTAMP                                 NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`  TIMESTAMP                                 NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
