-- Migration: done_becomes_status
-- Created at: 2017-12-15 20:04:52
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `done`,
    DROP COLUMN `done_at`,
    DROP COLUMN `done_description`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `done` BOOLEAN NOT NULL AFTER `rank`,
    ADD COLUMN `done_description` TEXT DEFAULT NULL AFTER `done`,
    ADD COLUMN `done_at` DATETIME(6) DEFAULT NULL AFTER `done_description`;

COMMIT;
