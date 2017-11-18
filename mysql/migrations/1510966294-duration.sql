-- Migration: duration
-- Created at: 2017-11-18 01:51:34
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `duration` VARCHAR(64) NOT NULL AFTER `description`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `duration`;

COMMIT;
