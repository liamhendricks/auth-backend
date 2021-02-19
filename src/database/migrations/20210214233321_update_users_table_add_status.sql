ALTER TABLE `users`
  ADD COLUMN `status` VARCHAR(120) COLLATE `utf8mb4_unicode_ci` NOT NULL DEFAULT "OK" AFTER `user_type`;
