CREATE TABLE `lesson_data` (
  `id`               BINARY(16)                                NOT NULL,
  `lesson_id`        BINARY(16)                                NOT NULL,
  `header_image`     VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `main_header`      TEXT COLLATE `utf8mb4_unicode_ci`,
  `main_description` TEXT COLLATE `utf8mb4_unicode_ci`,
  `main_body`        TEXT COLLATE `utf8mb4_unicode_ci`,
  `advice`           TEXT COLLATE `utf8mb4_unicode_ci`,
  `main_image`       VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `lesson_objective` TEXT COLLATE `utf8mb4_unicode_ci`,
  `first_obj_image`  VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `second_obj_image` VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL,
  `created_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`       TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`       TIMESTAMP NULL DEFAULT NULL,
  FOREIGN KEY (`lesson_id`) REFERENCES `lessons` (`id`),
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
