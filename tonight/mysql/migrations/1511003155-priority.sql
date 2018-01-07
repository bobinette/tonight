-- Migration: priority
-- Created at: 2017-11-18 12:05:55
-- ====  UP  ====

ALTER TABLE `tasks`
    ADD COLUMN `priority` INT NOT NULL AFTER `description`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `priority`;

COMMIT;
