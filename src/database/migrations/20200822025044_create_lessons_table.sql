CREATE TABLE `lessons` (
  `id`          BINARY(16)	                              NOT NULL,
  `course_id`   BINARY(16)	                              NOT NULL,
  `name`        VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL UNIQUE,
  `ordering`    INT                                       NOT NULL DEFAULT 0,
  `lesson_data` JSON,
  `created_at`  TIMESTAMP                                 NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`  TIMESTAMP                                 NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`  TIMESTAMP                                 NULL DEFAULT NULL,
  FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`),
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = `utf8mb4`
  COLLATE = `utf8mb4_unicode_ci`
