-- Migration: done_description
-- Created at: 2017-11-18 10:17:09
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    CHANGE `doneAt` `done_at` DATETIME(6) DEFAULT NULL AFTER `done`,
    ADD COLUMN `done_description` TEXT DEFAULT NULL AFTER `done`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    CHANGE `done_at` `doneAt` DATETIME(6) DEFAULT NULL AFTER `done`,
    DROP COLUMN `done_description`;

COMMIT;
