-- Migration: admin-user
-- Created at: 2018-05-08 10:54:30
-- ====  UP  ====

BEGIN;

ALTER TABLE `users`
    ADD COLUMN `is_admin` BOOLEAN DEFAULT FALSE AFTER `username`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `users`
    DROP COLUMN `is_admin`;

COMMIT;
